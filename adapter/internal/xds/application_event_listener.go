package xds

import (
	"context"
	"fmt"

	"github.com/wso2/apk/adapter/internal/loggers"
	cpv1alpha1 "github.com/wso2/apk/adapter/internal/operator/apis/cp/v1alpha1"
	"github.com/wso2/apk/adapter/pkg/logging"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// HandleApplicationEventsFromMgtServer handles the Application events
func HandleApplicationEventsFromMgtServer(c client.Client) {
	for applicationEvent := range applicationChannel {
		switch applicationEvent.Type {
		case APPLICATION_CREATE:
			err := c.Create(context.Background(), *&applicationEvent.Application)
			if err != nil {
				loggers.LoggerXds.ErrorC(logging.ErrorDetails{
					Message:   fmt.Sprint("Error creating application: ", err.Error()),
					Severity:  logging.CRITICAL,
					ErrorCode: 1707,
				})
			} else {
				loggers.LoggerXds.Info("Application created: " + applicationEvent.Application.Name)
			}
			break
		case APPLICATION_UPDATE:
			var application = new(cpv1alpha1.Application)
			if err := c.Get(context.Background(), types.NamespacedName{
				Name:      applicationEvent.Application.Name,
				Namespace: applicationEvent.Application.Namespace}, application); err != nil {
				loggers.LoggerXds.ErrorC(logging.ErrorDetails{
					Message:   fmt.Sprint("Error retrieving application: ", err.Error()),
					Severity:  logging.CRITICAL,
					ErrorCode: 1708,
				})
				break
			}
			application.Spec = applicationEvent.Application.Spec
			err := c.Update(context.Background(), application)
			if err != nil {
				loggers.LoggerXds.ErrorC(logging.ErrorDetails{
					Message:   fmt.Sprint("Error updating application: ", err.Error()),
					Severity:  logging.CRITICAL,
					ErrorCode: 1709,
				})
			} else {
				loggers.LoggerXds.Info("Application updated: " + applicationEvent.Application.Name)
			}
			break
		case APPLICATION_DELETE:
			err := c.Delete(context.Background(), *&applicationEvent.Application)
			if err != nil {
				loggers.LoggerXds.ErrorC(logging.ErrorDetails{
					Message:   fmt.Sprint("Error deleting application: ", err.Error()),
					Severity:  logging.CRITICAL,
					ErrorCode: 1710,
				})
			} else {
				loggers.LoggerXds.Info("Application deleted: " + applicationEvent.Application.Name)
			}
			break
		default:
			loggers.LoggerXds.Info("Unknown Application Event Type")
		}
	}
}
