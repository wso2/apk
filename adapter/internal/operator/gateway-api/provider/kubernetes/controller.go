/*
 *  Copyright (c) 2024, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 */

package kubernetes

// todo(amali) verify gateway api latest imports
// todo(amali) fix formatting in logs
import (
	"context"
	"fmt"

	"github.com/wso2/apk/adapter/config"
	"github.com/wso2/apk/adapter/internal/loggers"
	"github.com/wso2/apk/adapter/pkg/logging"
	"github.com/wso2/apk/adapter/pkg/utils/stringutils"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"

	"github.com/wso2/apk/adapter/internal/operator/constants"
	gatewayapi "github.com/wso2/apk/adapter/internal/operator/gateway-api"
	"github.com/wso2/apk/adapter/internal/operator/message"
	"github.com/wso2/apk/adapter/internal/operator/status"
	"github.com/wso2/apk/adapter/internal/operator/utils"
	dpv1alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha1"
	dpv1alpha2 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha2"
	discoveryv1 "k8s.io/api/discovery/v1"
	k8errors "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
	gwapiv1b1 "sigs.k8s.io/gateway-api/apis/v1beta1"
)

type gatewayReconcilerNew struct {
	client        client.Client
	statusUpdater *status.UpdateHandler
	mgr           manager.Manager
	resources     *message.ProviderResources
	extGVKs       []schema.GroupVersionKind
	namespace     string
	store         *kubernetesProviderStore
}

type resourceMappings struct {
	// Map for storing namespaces for Route, Service and Gateway objects.
	allAssociatedNamespaces map[string]struct{}
	// Map for storing backendRefs' NamespaceNames referred by various Route objects.
	allAssociatedBackendRefs map[gwapiv1.BackendObjectReference]struct{}
	// Map for storing APIs NamespaceNames
	allAssociatedAPIs map[*dpv1alpha1.API]struct{}
	// extensionRefFilters is a map of filters managed by an extension.
	// The key is the namespaced name of the filter and the value is the
	// unstructured form of the resource.
	extensionRefFilters map[types.NamespacedName]unstructured.Unstructured
}

const (
	gatewayClassControllerName = "wso2.com/apk-gateway-default"
	gatewayClassFinalizer      = gwapiv1.GatewayClassFinalizerGatewaysExist
)

// InitGatewayController creates a new Gateway controller instance. Gateway Controllers watches for gwapiv1.Gateway.
func InitGatewayController(mgr manager.Manager, pResourses *message.ProviderResources, statusUpdater *status.UpdateHandler) error {
	conf := config.ReadConfigs()
	r := &gatewayReconcilerNew{
		client:        mgr.GetClient(),
		statusUpdater: statusUpdater,
		mgr:           mgr,
		resources:     pResourses,
		namespace:     conf.Deployment.Gateway.Namespace,
		store:         newProviderStore(),
	}
	c, err := controller.New(constants.GatewayControllerNew, mgr, controller.Options{Reconciler: r})
	if err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error3119, logging.BLOCKER, "Error creating API controller, error: %v", err))
		return err
	}

	ctx := context.Background()
	// Watch resources
	if err := r.watchResources(ctx, mgr, c); err != nil {
		return err
	}
	// Subscribe to status updates
	r.subscribeAndUpdateStatus(ctx)
	loggers.LoggerAPKOperator.Info("New Controller successfully started. Watching Gateway Objects....")
	return nil
}

// Reconcile handles reconciling all resources in a single call. Any resource event should enqueue the
// same reconcile.Request containing the gateway controller name. This allows multiple resource updates to
// be handled by a single call to Reconcile. The reconcile.Request DOES NOT map to a specific resource.
func (r *gatewayReconcilerNew) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var (
		managedGCs []*gwapiv1.GatewayClass
		err        error
	)

	loggers.LoggerAPKOperator.Infof("Reconciling all gateways due to change in %s:%s", req.Namespace, req.Name)

	// Get the GatewayClasses managed by the Envoy Gateway Controller.
	managedGCs, err = r.managedGatewayClasses(ctx)
	if err != nil {
		return reconcile.Result{}, err
	}

	// The gatewayclass was already deleted/finalized and there are stale queue entries.
	if managedGCs == nil {
		r.resources.GatewayAPIResources.Delete(string(gatewayClassControllerName))
		loggers.LoggerAPKOperator.Info("no accepted gatewayclass")
		return reconcile.Result{}, nil
	}

	// we only manage one gatewayclass
	managedGC := managedGCs[0]

	// Collect all the Gateway API resources, Envoy Gateway customized resources,
	// and their referenced resources for the managed GatewayClasses, and store
	// them per GatewayClass.
	// For example:
	// - Gateway API resources: Gateways, xRoutes ...
	// - Referenced resources: Services, ServiceImports, EndpointSlices, Secrets, ConfigMaps ...
	gwcResource := gatewayapi.NewResources()
	gwcResource.GatewayClass = managedGC
	gwcResources := gatewayapi.ControllerResources{gwcResource}
	resourceMappings := newResourceMapping()

	// Add all Gateways, their associated Routes, and referenced resources to the resourceTree
	if err = r.processGateways(ctx, managedGC, resourceMappings, gwcResource); err != nil {
		return reconcile.Result{}, err
	}

	// Add the referenced services, ServiceImports, and EndpointSlices in
	// the collected BackendRefs to the resourceTree.
	// BackendRefs are referred by various Route objects and the ExtAuth in SecurityPolicies.
	r.processBackendRefs(ctx, gwcResource, resourceMappings)

	// For this particular Gateway, and all associated objects, check whether the
	// namespace exists. Add to the resourceTree.
	for ns := range resourceMappings.allAssociatedNamespaces {
		namespace, err := r.getNamespace(ctx, ns)
		if err != nil {
			loggers.LoggerAPKOperator.Error("Unable to find the namespace ", err)
			if k8errors.IsNotFound(err) {
				return reconcile.Result{}, nil
			}
			return reconcile.Result{}, err
		}

		gwcResource.Namespaces = append(gwcResource.Namespaces, namespace)
	}

	if err := r.updateStatusForGatewayClass(ctx, managedGC, true, string(gwapiv1.GatewayClassReasonAccepted), status.MsgValidGatewayClass); err != nil {
		loggers.LoggerAPKOperator.Errorf("Unable to update GatewayClass %s status %v", managedGC.Name, err)
		return reconcile.Result{}, err
	}
	if len(gwcResource.Gateways) == 0 {
		loggers.LoggerAPKOperator.Info("No gateways found for accepted gatewayclass ", managedGC.Name)

		// If needed, remove the finalizer from the accepted GatewayClass.
		if err := r.removeFinalizer(ctx, managedGC); err != nil {
			loggers.LoggerAPKOperator.Errorf("Failed to remove finalizer from gatewayclass %s, %v",
				managedGC.Name, err)
			return reconcile.Result{}, err
		}
	} else {
		// finalize the accepted GatewayClass.
		if err := r.addFinalizer(ctx, managedGC); err != nil {
			loggers.LoggerAPKOperator.Errorf("Failed adding finalizer to gatewayclass %s, %v",
				managedGC.Name, err)
			return reconcile.Result{}, err
		}
	}

	// Store the Gateway Resources for the GatewayClass.
	// The Store is triggered even when there are no Gateways associated to the
	// GatewayClass. This would happen in case the last Gateway is removed and the
	// Store will be required to trigger a cleanup of envoy infra resources.
	r.resources.GatewayAPIResources.Store(gatewayClassControllerName, &gwcResources)

	return ctrl.Result{}, nil
}

// managedGatewayClasses returns a list of GatewayClass objects that are managed by the APK Gateway Controller.
func (r *gatewayReconcilerNew) managedGatewayClasses(ctx context.Context) ([]*gwapiv1.GatewayClass, error) {
	var gatewayClasses gwapiv1.GatewayClassList
	if err := r.client.List(ctx, &gatewayClasses); err != nil {
		return nil, fmt.Errorf("error listing gatewayclasses: %w", err)
	}

	var cc controlledClasses

	for _, gwClass := range gatewayClasses.Items {
		gwClass := gwClass
		if gwClass.Spec.ControllerName == gatewayClassControllerName {
			// The gatewayclass was marked for deletion and the finalizer removed,
			// so clean-up dependents.
			if !gwClass.DeletionTimestamp.IsZero() &&
				!stringutils.StringInSlice(gatewayClassFinalizer, gwClass.Finalizers) {
				loggers.LoggerAPKOperator.Info("Gatewayclass marked for deletion")
				cc.removeMatch(&gwClass)
				continue
			}

			cc.addMatch(&gwClass)
		}
	}

	return cc.matchedClasses, nil
}

func newResourceMapping() *resourceMappings {
	return &resourceMappings{
		allAssociatedNamespaces:  map[string]struct{}{},
		allAssociatedBackendRefs: map[gwapiv1.BackendObjectReference]struct{}{},
		extensionRefFilters:      map[types.NamespacedName]unstructured.Unstructured{},
	}
}

func (r *gatewayReconcilerNew) processGateways(ctx context.Context, managedGC *gwapiv1.GatewayClass, resourceMap *resourceMappings, resourceTree *gatewayapi.Resources) error {
	// Find gateways for the managedGC
	// Find the Gateways that reference this Class.
	gatewayList := &gwapiv1.GatewayList{}
	if err := r.client.List(ctx, gatewayList, &client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(classGatewayIndex, managedGC.Name),
	}); err != nil || len(gatewayList.Items) < 1 {
		loggers.LoggerAPKOperator.Infof("No associated Gateways found for GatewayClass name %s, %v", managedGC.Name, err)
		return err
	}

	for _, gtw := range gatewayList.Items {
		loggers.LoggerAPKOperator.Infof("Processing Gateway namespace %s: Name %s", gtw.Namespace, gtw.Name)
		resourceMap.allAssociatedNamespaces[gtw.Namespace] = struct{}{}

		for _, listener := range gtw.Spec.Listeners {
			// Get Secret for gateway if it exists.
			if terminatesTLS(&listener) {
				for _, certRef := range listener.TLS.CertificateRefs {
					certRef := certRef
					if refsSecret(&certRef) {
						if err := r.processSecretRef(ctx, resourceMap, resourceTree, gatewayapi.KindGateway, gtw.Namespace,
							gtw.Name, certRef); err != nil {
							loggers.LoggerAPKOperator.Errorf("Failed to process TLS SecretRef for gateway %+v, secretRef %+v, Error: %v", gtw, certRef, err)
						}
					}
				}
			}
		}

		// Get HTTPRoute objects and check if it exists.
		if err := r.processHTTPRoutes(ctx, utils.NamespacedName(&gtw).String(), resourceMap, resourceTree); err != nil {
			return err
		}

		// Discard Status to reduce memory consumption in watchable
		// It will be recomputed by the gateway-api layer
		gtw.Status = gwapiv1.GatewayStatus{}
		resourceTree.Gateways = append(resourceTree.Gateways, &gtw)
	}

	return nil
}

// processSecretRef adds the referenced Secret to the resourceTree if it's valid.
// - If it exists in the same namespace as the owner.
// - If it exists in a different namespace, and there is a ReferenceGrant.
func (r *gatewayReconcilerNew) processSecretRef(
	ctx context.Context,
	resourceMap *resourceMappings,
	resourceTree *gatewayapi.Resources,
	ownerKind string,
	ownerNS string,
	ownerName string,
	secretRef gwapiv1.SecretObjectReference,
) error {
	secret := new(corev1.Secret)
	secretNS := gatewayapi.NamespaceDerefOr(secretRef.Namespace, ownerNS)
	err := r.client.Get(ctx,
		types.NamespacedName{Namespace: secretNS, Name: string(secretRef.Name)},
		secret,
	)
	if err != nil {
		return fmt.Errorf("unable to find the Secret: %s/%s, Error: %v", secretNS, string(secretRef.Name), err)
	}
	if secretNS != ownerNS {
		from := ObjectKindNamespacedName{
			kind:      ownerKind,
			namespace: ownerNS,
			name:      ownerName,
		}
		to := ObjectKindNamespacedName{
			kind:      gatewayapi.KindSecret,
			namespace: secretNS,
			name:      secret.Name,
		}
		refGrant, err := r.findReferenceGrant(ctx, from, to)
		switch {
		case err != nil:
			return fmt.Errorf("failed to find ReferenceGrant: %w", err)
		case refGrant == nil:
			return fmt.Errorf(
				"no matching ReferenceGrants found: from %s/%s to %s/%s",
				from.kind, from.namespace, to.kind, to.namespace)
		default:
			// RefGrant found
			resourceTree.ReferenceGrants = append(resourceTree.ReferenceGrants, refGrant)
			loggers.LoggerAPKOperator.Infof("Added ReferenceGrant to resource map namespace %s, name %s",
				refGrant.Namespace, refGrant.Name)
		}
	}
	resourceMap.allAssociatedNamespaces[secretNS] = struct{}{} // TODO Zhaohuabing do we need this line?
	resourceTree.Secrets = append(resourceTree.Secrets, secret)
	loggers.LoggerAPKOperator.Infof("Processing Secret namespace %s name %s", secretNS, string(secretRef.Name))
	return nil
}

func (r *gatewayReconcilerNew) findReferenceGrant(ctx context.Context, from, to ObjectKindNamespacedName) (*gwapiv1b1.ReferenceGrant, error) {
	refGrantList := new(gwapiv1b1.ReferenceGrantList)
	opts := &client.ListOptions{FieldSelector: fields.OneTermEqualSelector(targetRefGrantRouteIndex, to.kind)}
	if err := r.client.List(ctx, refGrantList, opts); err != nil {
		return nil, fmt.Errorf("failed to list ReferenceGrants: %w", err)
	}

	refGrants := refGrantList.Items

	for _, refGrant := range refGrants {
		if refGrant.Namespace == to.namespace {
			for _, src := range refGrant.Spec.From {
				if src.Kind == gwapiv1.Kind(from.kind) && string(src.Namespace) == from.namespace {
					return &refGrant, nil
				}
			}
		}
	}

	// No ReferenceGrant found.
	return nil, nil
}

func (r *gatewayReconcilerNew) enqueueClass(_ context.Context, _ client.Object) []reconcile.Request {
	return []reconcile.Request{{NamespacedName: types.NamespacedName{
		Name: string(gatewayClassControllerName),
	}}}
}

// watchResources watches gateway api resources.
func (r *gatewayReconcilerNew) watchResources(ctx context.Context, mgr manager.Manager, c controller.Controller) error {
	if err := c.Watch(
		source.Kind(mgr.GetCache(), &gwapiv1.GatewayClass{}),
		handler.EnqueueRequestsFromMapFunc(r.enqueueClass),
		predicate.GenerationChangedPredicate{},
		predicate.NewPredicateFuncs(r.hasMatchingController),
	); err != nil {
		return err
	}

	// Watch Gateway CRUDs and reconcile affected GatewayClass.
	gPredicates := []predicate.Predicate{
		predicate.GenerationChangedPredicate{},
		predicate.NewPredicateFuncs(r.validateGatewayForReconcile),
	}

	if err := c.Watch(
		source.Kind(mgr.GetCache(), &gwapiv1.Gateway{}),
		handler.EnqueueRequestsFromMapFunc(r.enqueueClass),
		gPredicates...,
	); err != nil {
		return err
	}
	if err := addGatewayIndexers(ctx, mgr); err != nil {
		return err
	}

	// Watch HTTPRoute CRUDs and process affected Gateways.
	httprPredicates := []predicate.Predicate{predicate.GenerationChangedPredicate{}}
	if err := c.Watch(
		source.Kind(mgr.GetCache(), &gwapiv1.HTTPRoute{}),
		handler.EnqueueRequestsFromMapFunc(r.enqueueClass),
		httprPredicates...,
	); err != nil {
		return err
	}
	if err := addHTTPRouteIndexers(ctx, mgr); err != nil {
		return err
	}

	// Watch API CRUDs and process affected Gateways.
	apiPredicates := []predicate.Predicate{predicate.GenerationChangedPredicate{}}
	if err := c.Watch(
		source.Kind(mgr.GetCache(), &dpv1alpha2.API{}),
		handler.EnqueueRequestsFromMapFunc(r.enqueueClass),
		apiPredicates...,
	); err != nil {
		return err
	}
	if err := addAPIIndexers(ctx, mgr); err != nil {
		return err
	}

	// // Watch GRPCRoute CRUDs and process affected Gateways.
	// grpcrPredicates := []predicate.Predicate{predicate.GenerationChangedPredicate{}}
	// if err := c.Watch(
	// 	source.Kind(mgr.GetCache(), &gwapiv1a2.GRPCRoute{}),
	// 	handler.EnqueueRequestsFromMapFunc(r.enqueueClass),
	// 	grpcrPredicates...,
	// ); err != nil {
	// 	return err
	// }
	// if err := addGRPCRouteIndexers(ctx, mgr); err != nil {
	// 	return err
	// }

	// // Watch TLSRoute CRUDs and process affected Gateways.
	// tlsrPredicates := []predicate.Predicate{predicate.GenerationChangedPredicate{}}
	// if r.namespaceLabel != nil {
	// 	tlsrPredicates = append(tlsrPredicates, predicate.NewPredicateFuncs(r.hasMatchingNamespaceLabels))
	// }
	// if err := c.Watch(
	// 	source.Kind(mgr.GetCache(), &gwapiv1a2.TLSRoute{}),
	// 	handler.EnqueueRequestsFromMapFunc(r.enqueueClass),
	// 	tlsrPredicates...,
	// ); err != nil {
	// 	return err
	// }
	// if err := addTLSRouteIndexers(ctx, mgr); err != nil {
	// 	return err
	// }

	// // Watch UDPRoute CRUDs and process affected Gateways.
	// udprPredicates := []predicate.Predicate{predicate.GenerationChangedPredicate{}}
	// if r.namespaceLabel != nil {
	// 	udprPredicates = append(udprPredicates, predicate.NewPredicateFuncs(r.hasMatchingNamespaceLabels))
	// }
	// if err := c.Watch(
	// 	source.Kind(mgr.GetCache(), &gwapiv1a2.UDPRoute{}),
	// 	handler.EnqueueRequestsFromMapFunc(r.enqueueClass),
	// 	udprPredicates...,
	// ); err != nil {
	// 	return err
	// }
	// if err := addUDPRouteIndexers(ctx, mgr); err != nil {
	// 	return err
	// }

	// // Watch TCPRoute CRUDs and process affected Gateways.
	// tcprPredicates := []predicate.Predicate{predicate.GenerationChangedPredicate{}}
	// if r.namespaceLabel != nil {
	// 	tcprPredicates = append(tcprPredicates, predicate.NewPredicateFuncs(r.hasMatchingNamespaceLabels))
	// }
	// if err := c.Watch(
	// 	source.Kind(mgr.GetCache(), &gwapiv1a2.TCPRoute{}),
	// 	handler.EnqueueRequestsFromMapFunc(r.enqueueClass),
	// 	tcprPredicates...,
	// ); err != nil {
	// 	return err
	// }
	// if err := addTCPRouteIndexers(ctx, mgr); err != nil {
	// 	return err
	// }

	// Watch Service CRUDs and process affected *Route objects and services belongs to gateways
	servicePredicates := []predicate.Predicate{predicate.NewPredicateFuncs(r.validateServiceForReconcile)}
	if err := c.Watch(
		source.Kind(mgr.GetCache(), &corev1.Service{}),
		handler.EnqueueRequestsFromMapFunc(r.enqueueClass),
		servicePredicates...,
	); err != nil {
		return err
	}

	// Watch Service CRUDs and process affected *Route objects and services belongs to gateways
	backendPredicates := []predicate.Predicate{predicate.NewPredicateFuncs(r.validateBackendForReconcile)}
	if err := c.Watch(
		source.Kind(mgr.GetCache(), &dpv1alpha2.Backend{}),
		handler.EnqueueRequestsFromMapFunc(r.enqueueClass),
		backendPredicates...,
	); err != nil {
		return err
	}

	// serviceImportCRDExists := r.serviceImportCRDExists(mgr)
	// if !serviceImportCRDExists {
	// 	loggers.LoggerAPKOperator.Info("ServiceImport CRD not found, skipping ServiceImport watch")
	// }

	// // Watch ServiceImport CRUDs and process affected *Route objects.
	// if serviceImportCRDExists {
	// 	if err := c.Watch(
	// 		source.Kind(mgr.GetCache(), &mcsapi.ServiceImport{}),
	// 		handler.EnqueueRequestsFromMapFunc(r.enqueueClass),
	// 		predicate.GenerationChangedPredicate{},
	// 		predicate.NewPredicateFuncs(r.validateServiceImportForReconcile)); err != nil {
	// 		// ServiceImport is not available in the cluster, skip the watch and not throw error.
	// 		loggers.LoggerAPKOperator.Info("unable to watch ServiceImport: %s", err.Error())
	// 	}
	// }

	// // Watch EndpointSlice CRUDs and process affected *Route objects.
	// esPredicates := []predicate.Predicate{
	// 	predicate.GenerationChangedPredicate{},
	// 	predicate.NewPredicateFuncs(r.validateEndpointSliceForReconcile),
	// }
	// if r.namespaceLabel != nil {
	// 	esPredicates = append(esPredicates, predicate.NewPredicateFuncs(r.hasMatchingNamespaceLabels))
	// }
	// if err := c.Watch(
	// 	source.Kind(mgr.GetCache(), &discoveryv1.EndpointSlice{}),
	// 	handler.EnqueueRequestsFromMapFunc(r.enqueueClass),
	// 	esPredicates...,
	// ); err != nil {
	// 	return err
	// }

	// Watch Node CRUDs to update Gateway Address exposed by Service of type NodePort.
	// Node creation/deletion and ExternalIP updates would require update in the Gateway
	nPredicates := []predicate.Predicate{
		predicate.GenerationChangedPredicate{},
		predicate.NewPredicateFuncs(r.handleNode),
	}
	// resource address.
	if err := c.Watch(
		source.Kind(mgr.GetCache(), &corev1.Node{}),
		handler.EnqueueRequestsFromMapFunc(r.enqueueClass),
		nPredicates...,
	); err != nil {
		return err
	}

	// Watch Secret CRUDs and process affected EG CRs (Gateway, SecurityPolicy, more in the future).
	secretPredicates := []predicate.Predicate{
		predicate.GenerationChangedPredicate{},
		predicate.NewPredicateFuncs(r.validateSecretForReconcile),
	}
	if err := c.Watch(
		source.Kind(mgr.GetCache(), &corev1.Secret{}),
		handler.EnqueueRequestsFromMapFunc(r.enqueueClass),
		secretPredicates...,
	); err != nil {
		return err
	}

	// // Watch ConfigMap CRUDs and process affected ClienTraffiPolicies and BackendTLSPolicies.
	// configMapPredicates := []predicate.Predicate{
	// 	predicate.GenerationChangedPredicate{},
	// 	predicate.NewPredicateFuncs(r.validateConfigMapForReconcile),
	// }
	// if r.namespaceLabel != nil {
	// 	configMapPredicates = append(configMapPredicates, predicate.NewPredicateFuncs(r.hasMatchingNamespaceLabels))
	// }
	// if err := c.Watch(
	// 	source.Kind(mgr.GetCache(), &corev1.ConfigMap{}),
	// 	handler.EnqueueRequestsFromMapFunc(r.enqueueClass),
	// 	configMapPredicates...,
	// ); err != nil {
	// 	return err
	// }

	// Watch ReferenceGrant CRUDs and process affected Gateways.
	rgPredicates := []predicate.Predicate{predicate.GenerationChangedPredicate{}}
	if err := c.Watch(
		source.Kind(mgr.GetCache(), &gwapiv1b1.ReferenceGrant{}),
		handler.EnqueueRequestsFromMapFunc(r.enqueueClass),
		rgPredicates...,
	); err != nil {
		return err
	}
	if err := addReferenceGrantIndexers(ctx, mgr); err != nil {
		return err
	}

	// Watch Deployment CRUDs and process affected Gateways.
	dPredicates := []predicate.Predicate{predicate.NewPredicateFuncs(r.validateDeploymentForReconcile)}
	if err := c.Watch(
		source.Kind(mgr.GetCache(), &appsv1.Deployment{}),
		handler.EnqueueRequestsFromMapFunc(r.enqueueClass),
		dPredicates...,
	); err != nil {
		return err
	}

	// // Watch BackendTLSPolicy
	// btlsPredicates := []predicate.Predicate{predicate.GenerationChangedPredicate{}}
	// if r.namespaceLabel != nil {
	// 	btlsPredicates = append(btlsPredicates, predicate.NewPredicateFuncs(r.hasMatchingNamespaceLabels))
	// }

	// if err := c.Watch(
	// 	source.Kind(mgr.GetCache(), &gwapiv1a2.BackendTLSPolicy{}),
	// 	handler.EnqueueRequestsFromMapFunc(r.enqueueClass),
	// 	btlsPredicates...,
	// ); err != nil {
	// 	return err
	// }

	// if err := addBtlsIndexers(ctx, mgr); err != nil {
	// 	return err
	// }

	// loggers.LoggerAPKOperator.Info("Watching gatewayAPI related objects")

	// // Watch any additional GVKs from the registered extension.
	// uPredicates := []predicate.Predicate{predicate.GenerationChangedPredicate{}}
	// if r.namespaceLabel != nil {
	// 	uPredicates = append(uPredicates, predicate.NewPredicateFuncs(r.hasMatchingNamespaceLabels))
	// }
	// for _, gvk := range r.extGVKs {
	// 	u := &unstructured.Unstructured{}
	// 	u.SetGroupVersionKind(gvk)
	// 	if err := c.Watch(source.Kind(mgr.GetCache(), u),
	// 		handler.EnqueueRequestsFromMapFunc(r.enqueueClass),
	// 		uPredicates...,
	// 	); err != nil {
	// 		return err
	// 	}
	// 	loggers.LoggerAPKOperator.Info("Watching additional resource", "resource", gvk.String())
	// }

	return nil
}

// processBackendRefs adds the referenced resources in BackendRefs to the resourceTree, including:
// - Services
// - ServiceImports
// - EndpointSlices
func (r *gatewayReconcilerNew) processBackendRefs(ctx context.Context, gwcResource *gatewayapi.Resources, resourceMappings *resourceMappings) {
	for backendRef := range resourceMappings.allAssociatedBackendRefs {
		backendRefKind := gatewayapi.KindDerefOr(backendRef.Kind, gatewayapi.KindService)
		loggers.LoggerAPKOperator.Infof("Processing Backend kind %s, namespace %s, name %s", backendRefKind,
			string(*backendRef.Namespace), string(backendRef.Name))

		var endpointSliceLabelKey string
		switch backendRefKind {
		case gatewayapi.KindService:
			service := new(corev1.Service)
			err := r.client.Get(ctx, types.NamespacedName{Namespace: string(*backendRef.Namespace), Name: string(backendRef.Name)}, service)
			if err != nil {
				loggers.LoggerAPKOperator.Errorf("Failed to get Service namespace %s, name: %s, Error: %v",
					string(*backendRef.Namespace), string(backendRef.Name), err)
			} else {
				resourceMappings.allAssociatedNamespaces[service.Namespace] = struct{}{}
				gwcResource.Services = append(gwcResource.Services, service)
				loggers.LoggerAPKOperator.Infof("Added Service to resource tree namespace %s, name: %s",
					string(*backendRef.Namespace), string(backendRef.Name))
			}
			endpointSliceLabelKey = discoveryv1.LabelServiceName

		case gatewayapi.KindBackend:
			backend := new(dpv1alpha2.Backend)
			err := r.client.Get(ctx, types.NamespacedName{Namespace: string(*backendRef.Namespace), Name: string(backendRef.Name)}, backend)
			if err != nil {
				loggers.LoggerAPKOperator.Errorf("Failed to get Backend namespace %s, name %s, Error: %v", string(*backendRef.Namespace),
					string(backendRef.Name), err)
			} else {
				resourceMappings.allAssociatedNamespaces[backend.Namespace] = struct{}{}
				gwcResource.Backends = append(gwcResource.Backends, backend)
				loggers.LoggerAPKOperator.Infof("Added Backend to resource tree namespace %s, name %s",
					string(*backendRef.Namespace), string(backendRef.Name))
			}
		}

		if endpointSliceLabelKey != "" {
			// Retrieve the EndpointSlices associated with the service
			endpointSliceList := new(discoveryv1.EndpointSliceList)
			opts := []client.ListOption{
				client.MatchingLabels(map[string]string{
					endpointSliceLabelKey: string(backendRef.Name),
				}),
				client.InNamespace(string(*backendRef.Namespace)),
			}
			if err := r.client.List(ctx, endpointSliceList, opts...); err != nil {
				loggers.LoggerAPKOperator.Errorf("Failed to get EndpointSlices namespace %s, kind %s, name %s, Error %v",
					string(*backendRef.Namespace), backendRefKind, string(backendRef.Name), err)
			} else {
				for _, endpointSlice := range endpointSliceList.Items {
					endpointSlice := endpointSlice
					loggers.LoggerAPKOperator.Infof("Added EndpointSlice to resource tree namespace %s, name %s",
						endpointSlice.Namespace, endpointSlice.Name)
					gwcResource.EndpointSlices = append(gwcResource.EndpointSlices, &endpointSlice)
				}
			}
		}
	}
}

func (r *gatewayReconcilerNew) getNamespace(ctx context.Context, name string) (*corev1.Namespace, error) {
	nsKey := types.NamespacedName{Name: name}
	ns := new(corev1.Namespace)
	if err := r.client.Get(ctx, nsKey, ns); err != nil {
		loggers.LoggerAPKOperator.Error("Unable to get Namespace ", err)
		return nil, err
	}
	return ns, nil
}

func (r *gatewayReconcilerNew) updateStatusForGatewayClass(
	ctx context.Context,
	gc *gwapiv1.GatewayClass,
	accepted bool,
	reason,
	msg string) error {
	if r.statusUpdater != nil {
		r.statusUpdater.Send(status.Update{
			NamespacedName: types.NamespacedName{Name: gc.Name},
			Resource:       &gwapiv1.GatewayClass{},
			UpdateStatus: func(obj client.Object) client.Object {
				gc, ok := obj.(*gwapiv1.GatewayClass)
				if !ok {
					panic(fmt.Sprintf("unsupported object type %T", obj))
				}

				return status.SetGatewayClassAccepted(gc.DeepCopy(), accepted, reason, msg)
			},
		})
	} else {
		// this branch makes testing easier by not going through the status.Updater.
		duplicate := status.SetGatewayClassAccepted(gc.DeepCopy(), accepted, reason, msg)

		if err := r.client.Status().Update(ctx, duplicate); err != nil && !k8errors.IsNotFound(err) {
			return fmt.Errorf("error updating status of gatewayclass %s: %w", duplicate.Name, err)
		}
	}
	return nil
}

// removeFinalizer removes the gatewayclass finalizer from the provided gc, if it exists.
func (r *gatewayReconcilerNew) removeFinalizer(ctx context.Context, gc *gwapiv1.GatewayClass) error {
	if stringutils.StringInSlice(gatewayClassFinalizer, gc.Finalizers) {
		base := client.MergeFrom(gc.DeepCopy())
		gc.Finalizers = stringutils.RemoveString(gc.Finalizers, gatewayClassFinalizer)
		if err := r.client.Patch(ctx, gc, base); err != nil {
			return fmt.Errorf("failed to remove finalizer from gatewayclass %s: %w", gc.Name, err)
		}
	}
	return nil
}

// addFinalizer adds the gatewayclass finalizer to the provided gc, if it doesn't exist.
func (r *gatewayReconcilerNew) addFinalizer(ctx context.Context, gc *gwapiv1.GatewayClass) error {
	if !stringutils.StringInSlice(gatewayClassFinalizer, gc.Finalizers) {
		base := client.MergeFrom(gc.DeepCopy())
		gc.Finalizers = append(gc.Finalizers, gatewayClassFinalizer)
		if err := r.client.Patch(ctx, gc, base); err != nil {
			return fmt.Errorf("failed to add finalizer to gatewayclass %s: %w", gc.Name, err)
		}
	}
	return nil
}
