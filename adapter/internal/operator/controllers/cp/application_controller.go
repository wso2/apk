/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cp

import (
	"context"
	"fmt"

	"github.com/wso2/apk/adapter/internal/discovery/xds"
	"github.com/wso2/apk/adapter/internal/loggers"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	cpv1alpha1 "github.com/wso2/apk/adapter/internal/operator/apis/cp/v1alpha1"
	"github.com/wso2/apk/adapter/pkg/discovery/api/wso2/discovery/subscription"
)

// ApplicationReconciler reconciles a Application object
type ApplicationReconciler struct {
	client           client.Client
	Scheme           *runtime.Scheme
	applicationCache applicationCache
}

type applicationCache map[types.NamespacedName]*cpv1alpha1.Application

// NewApplicationController creates a new Application controller instance.
func NewApplicationController(mgr manager.Manager) error {
	r := &ApplicationReconciler{
		client:           mgr.GetClient(),
		applicationCache: applicationCache{},
	}
	c, err := controller.New("Application", mgr, controller.Options{Reconciler: r})
	if err != nil {
		loggers.LoggerAPKOperator.Errorf("Error creating Application controller: %v", err)
		return err
	}

	if err := c.Watch(&source.Kind{Type: &cpv1alpha1.Application{}}, &handler.EnqueueRequestForObject{}); err != nil {
		loggers.LoggerAPKOperator.Errorf("Error watching Application resources: %v", err)
		return err
	}

	loggers.LoggerAPKOperator.Info("Application Controller successfully started. Watching Application Objects....")
	return nil
}

//+kubebuilder:rbac:groups=cp.wso2.com,resources=applications,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cp.wso2.com,resources=applications/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=cp.wso2.com,resources=applications/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// Modify the Reconcile function to compare the state specified by
// the Application object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *ApplicationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	loggers.LoggerAPKOperator.Info("Reconciling application: " + req.NamespacedName.String())

	applicationKey := req.NamespacedName
	var application = new(cpv1alpha1.Application)
	if err := r.client.Get(ctx, applicationKey, application); err != nil {
		if errors.IsNotFound(err) {
			// The application doesn't exist in the applicationCache, remove it
			delete(r.applicationCache, applicationKey)
			loggers.LoggerAPKOperator.Info("deleted application from application cache")

			sendUpdates(r.applicationCache)
			return ctrl.Result{}, nil
		}
		return reconcile.Result{}, fmt.Errorf("failed to get application %s/%s",
			applicationKey.Namespace, applicationKey.Name)
	}
	// The application created / updated, add to the map
	if r.applicationCache == nil {
		r.applicationCache = applicationCache{}
	}
	r.applicationCache[applicationKey] = application

	sendUpdates(r.applicationCache)
	return ctrl.Result{}, nil
}

func sendUpdates(cache applicationCache) {
	var applications []cpv1alpha1.ApplicationSpec
	for _, application := range cache {
		applications = append(applications, application.Spec)
	}
	appList := marshalApplicationList(applications)
	xds.UpdateEnforcerApplications(appList)

	subList := marshalSubscriptionList(applications)
	xds.UpdateEnforcerSubscriptions(subList)

	appKeyMappingList := marshalApplicationKeyMapping(applications)
	xds.UpdateEnforcerApplicationKeyMappings(appKeyMappingList)
}

func marshalApplicationList(applicationList []cpv1alpha1.ApplicationSpec) *subscription.ApplicationList {
	applications := []*subscription.Application{}
	for _, appInternal := range applicationList {
		app := &subscription.Application{
			Uuid:       appInternal.UUID,
			Name:       appInternal.Name,
			Policy:     appInternal.Policy,
			Attributes: appInternal.Attributes,
		}
		applications = append(applications, app)
	}
	return &subscription.ApplicationList{
		List: applications,
	}
}

func marshalSubscriptionList(applicationList []cpv1alpha1.ApplicationSpec) *subscription.SubscriptionList {
	subscriptions := []*subscription.Subscription{}
	for _, appInternal := range applicationList {
		for _, subInternal := range appInternal.Subscriptions {
			sub := &subscription.Subscription{
				SubscriptionUUID:  subInternal.UUID,
				PolicyId:          subInternal.PolicyID,
				SubscriptionState: subInternal.SubscriptionStatus,
				AppUUID:           appInternal.UUID,
			}
			subscriptions = append(subscriptions, sub)
		}
	}
	return &subscription.SubscriptionList{
		List: subscriptions,
	}
}

func marshalApplicationKeyMapping(applicationList []cpv1alpha1.ApplicationSpec) *subscription.ApplicationKeyMappingList {
	applicationKeyMappings := []*subscription.ApplicationKeyMapping{}
	for _, appInternal := range applicationList {
		for _, consumerKeyInternal := range appInternal.ConsumerKeys {
			consumerKey := &subscription.ApplicationKeyMapping{
				ConsumerKey:     consumerKeyInternal.Key,
				KeyManager:      consumerKeyInternal.KeyManager,
				ApplicationUUID: appInternal.UUID,
			}
			applicationKeyMappings = append(applicationKeyMappings, consumerKey)
		}
	}
	return &subscription.ApplicationKeyMappingList{
		List: applicationKeyMappings,
	}
}
