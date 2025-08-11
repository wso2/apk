package crbuilder

import (
	"fmt"
	"github.com/wso2/apk/config-deployer-service-go/internal/dto"
	"sigs.k8s.io/controller-runtime/pkg/client"
	dpv2alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v2alpha1"
)

func CreateResources(apiResourceBundle *dto.APIResourceBundle) ([]client.Object, error) {
	var errors []error
	if apiResourceBundle == nil || apiResourceBundle.APKConf == nil || apiResourceBundle.APKConf.EndpointConfigurations == nil {
		return nil, fmt.Errorf("invalid APIResourceBundle")
	}

	objects := make([]client.Object, 0)

	// Create RouteMetadata
	if len(apiResourceBundle.APKConf.EndpointConfigurations.Production) > 0 {
		// Create the RouteMetadata object
		routeMetadata := &dpv2alpha1.RouteMetadata{
			Name:      apiResourceBundle.Name,
			Namespace: apiResourceBundle.Namespace,
		}
		objects = append(objects, routeMetadata)
	}
	if len(apiResourceBundle.APKConf.EndpointConfigurations.Sandbox) > 0 {
		// Create the RouteMetadata object for sandbox
		routeMetadata := &dpv2alpha1.RouteMetadata{
			Name:      apiResourceBundle.Name,
			Namespace: apiResourceBundle.Namespace,
		}
		objects = append(objects, routeMetadata)
	}

	return objects, errors

}