package crbuilder

import (
	"fmt"

	dpv2alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v2alpha1"
	"github.com/wso2/apk/config-deployer-service-go/internal/constants"
	"github.com/wso2/apk/config-deployer-service-go/internal/dto"
	util "github.com/wso2/apk/config-deployer-service-go/internal/util"
	"sigs.k8s.io/controller-runtime/pkg/client"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateResources(apiResourceBundle *dto.APIResourceBundle) ([]client.Object, error) {
	var error error
	if apiResourceBundle == nil || apiResourceBundle.APKConf == nil || apiResourceBundle.APKConf.EndpointConfigurations == nil {
		return nil, fmt.Errorf("invalid APIResourceBundle")
	}

	objects := make([]client.Object, 0)

	// Create RouteMetadata
	if len(apiResourceBundle.APKConf.EndpointConfigurations.Production) > 0 {
		// Create the RouteMetadata object
		routeMetadata := &dpv2alpha1.RouteMetadata{
			ObjectMeta: metav1.ObjectMeta{
				Name:      util.GenerateRouteMetadataName(apiResourceBundle.APKConf.Name, constants.SANDBOX_TYPE, apiResourceBundle.APKConf.Version, apiResourceBundle.Organization),
				Namespace: apiResourceBundle.Namespace,
			},
		}
		objects = append(objects, routeMetadata)
	}
	if len(apiResourceBundle.APKConf.EndpointConfigurations.Sandbox) > 0 {
		// Create the RouteMetadata object for sandbox
		routeMetadata := &dpv2alpha1.RouteMetadata{
			ObjectMeta: metav1.ObjectMeta{
				Name:      util.GenerateRouteMetadataName(apiResourceBundle.APKConf.Name, constants.SANDBOX_TYPE, apiResourceBundle.APKConf.Version, apiResourceBundle.Organization),
				Namespace: apiResourceBundle.Namespace,
			},
		}
		objects = append(objects, routeMetadata)
	}

	return objects, error

}