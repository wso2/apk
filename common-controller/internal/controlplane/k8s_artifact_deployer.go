package controlplane

import (
	"context"

	"github.com/wso2/apk/common-controller/internal/server"
	cpv1alpha2 "github.com/wso2/apk/common-go-libs/apis/cp/v1alpha2"
	"github.com/wso2/apk/common-go-libs/utils"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// K8sArtifactDeployer is a struct that implements ArtifactDeployer interface
type K8sArtifactDeployer struct {
	client client.Client
}

// NewK8sArtifactDeployer creates a new K8sArtifactDeployer
func NewK8sArtifactDeployer(mgr manager.Manager) K8sArtifactDeployer {
	return K8sArtifactDeployer{client: mgr.GetClient()}
}

// DeployApplication deploys an application
func (k8sArtifactDeployer K8sArtifactDeployer) DeployApplication(application server.Application) error {
	crApplication := cpv1alpha2.Application{ObjectMeta: v1.ObjectMeta{Name: application.UUID, Namespace: utils.GetOperatorPodNamespace()},
		Spec: cpv1alpha2.ApplicationSpec{Name: application.Name, Owner: application.Owner, Organization: application.OrganizationID, Attributes: application.Attributes}}
	return k8sArtifactDeployer.client.Create(context.Background(), &crApplication)
}

// DeploySubscription deploys a subscription
func (k8sArtifactDeployer K8sArtifactDeployer) DeploySubscription(subscription server.Subscription) error {
	crSubscription := cpv1alpha2.Subscription{ObjectMeta: v1.ObjectMeta{Name: subscription.UUID, Namespace: utils.GetOperatorPodNamespace()},
		Spec: cpv1alpha2.SubscriptionSpec{Organization: subscription.Organization, API: cpv1alpha2.API{Name: subscription.SubscribedAPI.Name, Version: subscription.SubscribedAPI.Version}, SubscriptionStatus: subscription.SubStatus}}
	return k8sArtifactDeployer.client.Create(context.Background(), &crSubscription)
}

// DeployApplicationMappings deploys an application mapping
func (k8sArtifactDeployer K8sArtifactDeployer) DeployApplicationMappings(applicationMapping server.ApplicationMapping) error {
	crApplicationMapping := cpv1alpha2.ApplicationMapping{ObjectMeta: v1.ObjectMeta{Name: applicationMapping.UUID, Namespace: utils.GetOperatorPodNamespace()},
		Spec: cpv1alpha2.ApplicationMappingSpec{ApplicationRef: applicationMapping.ApplicationRef, SubscriptionRef: applicationMapping.SubscriptionRef}}
	return k8sArtifactDeployer.client.Create(context.Background(), &crApplicationMapping)
}

// DeployKeyMappings deploys a key mapping
func (k8sArtifactDeployer K8sArtifactDeployer) DeployKeyMappings(keyMapping server.ApplicationKeyMapping) error {
	var crApplication cpv1alpha2.Application
	err := k8sArtifactDeployer.client.Get(context.Background(), client.ObjectKey{Name: keyMapping.ApplicationUUID, Namespace: utils.GetOperatorPodNamespace()}, &crApplication)
	if err != nil {
		return err
	}
	securitySchemes := cpv1alpha2.SecuritySchemes{}
	if crApplication.Spec.SecuritySchemes != nil {
		securitySchemes = *crApplication.Spec.SecuritySchemes
	}
	if keyMapping.SecurityScheme == "OAuth2" {
		if securitySchemes.OAuth2 == nil {
			securitySchemes.OAuth2 = &cpv1alpha2.SecurityScheme{Environments: []cpv1alpha2.Environment{generateSecurityScheme(keyMapping)}}
		} else {
			environments := make([]cpv1alpha2.Environment, 0)
			for _, environment := range securitySchemes.OAuth2.Environments {
				if environment.EnvID != keyMapping.EnvID || environment.AppID != keyMapping.ApplicationIdentifier || environment.KeyType != keyMapping.KeyType {
					environments = append(environments, environment)
				}
			}
			securitySchemes.OAuth2.Environments = append(environments, generateSecurityScheme(keyMapping))
		}
	}
	return k8sArtifactDeployer.client.Update(context.Background(), &crApplication)
}

// DeleteAllApplicationMappings deletes all application mappings
func (k8sArtifactDeployer K8sArtifactDeployer) DeleteAllApplicationMappings() error {
	return nil
}

// DeleteAllApplications deletes all applications
func (k8sArtifactDeployer K8sArtifactDeployer) DeleteAllApplications() error {
	return nil
}

// DeleteAllKeyMappings deletes all key mappings
func (k8sArtifactDeployer K8sArtifactDeployer) DeleteAllKeyMappings() error {
	return nil
}

// DeleteAllSubscriptions deletes all subscriptions
func (k8sArtifactDeployer K8sArtifactDeployer) DeleteAllSubscriptions() error {
	return nil
}

// DeleteApplication deletes an application
func (k8sArtifactDeployer K8sArtifactDeployer) DeleteApplication(applicationID string) error {
	return nil
}

// DeleteApplicationMappings deletes an application mapping
func (k8sArtifactDeployer K8sArtifactDeployer) DeleteApplicationMappings(applicationID string) error {
	return nil
}

// DeleteKeyMappings deletes a key mapping
func (k8sArtifactDeployer K8sArtifactDeployer) DeleteKeyMappings(applicationID string) error {
	return nil
}

// DeleteSubscription deletes a subscription
func (k8sArtifactDeployer K8sArtifactDeployer) DeleteSubscription(subscriptionID string) error {
	return nil
}

// DeployAllApplicationMappings deploys all application mappings
func (k8sArtifactDeployer K8sArtifactDeployer) DeployAllApplicationMappings(applicationMappings server.ApplicationMappingList) error {
	return nil
}

// DeployAllApplications deploys all key mappings
func (k8sArtifactDeployer K8sArtifactDeployer) DeployAllApplications(applications server.ApplicationList) error {
	return nil
}

// DeployAllKeyMappings deploys all key mappings
func (k8sArtifactDeployer K8sArtifactDeployer) DeployAllKeyMappings(keyMappings server.ApplicationKeyMappingList) error {
	return nil
}

// DeployAllSubscriptions deploys all subscriptions
func (k8sArtifactDeployer K8sArtifactDeployer) DeployAllSubscriptions(subscriptions server.SubscriptionList) error {
	return nil
}

// GetAllApplicationMappings returns all application mappings
func (k8sArtifactDeployer K8sArtifactDeployer) GetAllApplicationMappings() (server.ApplicationMappingList, error) {
	return server.ApplicationMappingList{}, nil
}

// GetAllApplications returns all applications
func (k8sArtifactDeployer K8sArtifactDeployer) GetAllApplications() (server.ApplicationList, error) {
	return server.ApplicationList{}, nil
}

// GetAllKeyMappings returns all key mappings
func (k8sArtifactDeployer K8sArtifactDeployer) GetAllKeyMappings() (server.ApplicationKeyMappingList, error) {
	return server.ApplicationKeyMappingList{}, nil
}

// GetAllSubscriptions returns all subscriptions
func (k8sArtifactDeployer K8sArtifactDeployer) GetAllSubscriptions() (server.SubscriptionList, error) {
	return server.SubscriptionList{}, nil
}

// GetApplication returns an application
func (k8sArtifactDeployer K8sArtifactDeployer) GetApplication(applicationID string) (server.Application, error) {
	return server.Application{}, nil
}

// GetApplicationMappings returns an application mapping
func (k8sArtifactDeployer K8sArtifactDeployer) GetApplicationMappings(applicationID string) (server.ApplicationMapping, error) {
	return server.ApplicationMapping{}, nil
}

// GetKeyMappings returns a key mapping
func (k8sArtifactDeployer K8sArtifactDeployer) GetKeyMappings(applicationID string) (server.ApplicationKeyMapping, error) {
	return server.ApplicationKeyMapping{}, nil
}

// GetSubscription returns a subscription
func (k8sArtifactDeployer K8sArtifactDeployer) GetSubscription(subscriptionID string) (server.Subscription, error) {
	return server.Subscription{}, nil
}

// GenerateSecurityScheme generates a security scheme
func generateSecurityScheme(keyMapping server.ApplicationKeyMapping) cpv1alpha2.Environment {
	return cpv1alpha2.Environment{EnvID: keyMapping.EnvID, AppID: keyMapping.ApplicationIdentifier, KeyType: keyMapping.KeyType}
}
