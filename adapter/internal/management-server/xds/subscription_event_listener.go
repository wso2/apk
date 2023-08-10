package xds

// import (
// 	"context"
// 	"errors"

// 	"github.com/wso2/apk/adapter/internal/loggers"
// 	logging "github.com/wso2/apk/adapter/internal/logging"
// 	cpv1alpha1 "github.com/wso2/apk/adapter/internal/operator/apis/cp/v1alpha1"

// 	apierrors "k8s.io/apimachinery/pkg/api/errors"
// 	"k8s.io/apimachinery/pkg/types"
// 	ctrlcache "sigs.k8s.io/controller-runtime/pkg/cache"
// 	"sigs.k8s.io/controller-runtime/pkg/client"
// )

// // HandleSubscriptionEventsFromMgtServer handles the Subscription events
// func HandleSubscriptionEventsFromMgtServer(c client.Client, cReader client.Reader) {
// 	for subscriptionEvent := range subscriptionChannel {
// 		switch subscriptionEvent.Type {
// 		case SubscriptionCreate:
// 			if found, _, err := checkSubscriptionExists(subscriptionEvent.Subscription, c, cReader); err == nil && !found {
// 				if err := c.Create(context.Background(), *&subscriptionEvent.Subscription); err != nil {
// 					loggers.LoggerXds.ErrorC(logging.PrintError(logging.Error1721, logging.CRITICAL, "Error creating subscription: %v", err.Error()))
// 				} else {
// 					loggers.LoggerXds.Info("Subscription created: " + subscriptionEvent.Subscription.Name)
// 				}
// 			}
// 			break
// 		case SubscriptionUpdate:
// 			if found, subscription, err := checkSubscriptionExists(subscriptionEvent.Subscription, c, cReader); err == nil && found {
// 				subscription.Spec = subscriptionEvent.Subscription.Spec
// 				err := c.Update(context.Background(), subscription)
// 				if err != nil {
// 					loggers.LoggerXds.ErrorC(logging.PrintError(logging.Error1722, logging.CRITICAL, "Error updating subscription: %v", err.Error()))
// 				} else {
// 					loggers.LoggerXds.Info("Subscription updated: " + subscriptionEvent.Subscription.Name)
// 				}
// 			}
// 			break
// 		case SubscriptionDelete:
// 			err := c.Delete(context.Background(), *&subscriptionEvent.Subscription)
// 			if err != nil {
// 				loggers.LoggerXds.ErrorC(logging.PrintError(logging.Error1723, logging.CRITICAL, "Error deleting subscription: %v", err.Error()))
// 			} else {
// 				loggers.LoggerXds.Info("Subscription deleted: " + subscriptionEvent.Subscription.Name)
// 			}
// 			break
// 		default:
// 			loggers.LoggerXds.Info("Unknown Subscription Event Type")
// 		}
// 	}
// }

// func checkSubscriptionExists(subscription *cpv1alpha1.Subscription, c client.Client, cReader client.Reader) (bool, *cpv1alpha1.Subscription, error) {
// 	var retrivedSubscription = new(cpv1alpha1.Subscription)
// 	// Try reading from cache
// 	if err := c.Get(context.Background(), types.NamespacedName{
// 		Name:      subscription.Name,
// 		Namespace: subscription.Namespace}, retrivedSubscription); err != nil {

// 		target := &ctrlcache.ErrCacheNotStarted{}
// 		if errors.As(err, &target) {
// 			// Try reading from api server directly
// 			if err := cReader.Get(context.Background(), types.NamespacedName{
// 				Name:      subscription.Name,
// 				Namespace: subscription.Namespace}, retrivedSubscription); err != nil {

// 				if !apierrors.IsNotFound(err) {
// 					loggers.LoggerXds.ErrorC(logging.PrintError(logging.Error1724, logging.CRITICAL, "Error retrieving subscription: %v", err.Error()))
// 					return false, nil, err
// 				}
// 				return false, nil, nil
// 			}
// 		} else if !apierrors.IsNotFound(err) {
// 			loggers.LoggerXds.ErrorC(logging.PrintError(logging.Error1724, logging.CRITICAL, "Error retrieving subscription: %v", err.Error()))
// 			return false, nil, err
// 		} else {
// 			return false, nil, nil
// 		}
// 	}
// 	return true, retrivedSubscription, nil
// }
