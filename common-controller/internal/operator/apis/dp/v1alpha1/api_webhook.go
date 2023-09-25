/*
 *  Copyright (c) 2022, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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

package v1alpha1

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/wso2/apk/adapter/pkg/logging"
	"github.com/wso2/apk/common-controller/internal/config"
	"github.com/wso2/apk/common-controller/internal/loggers"
	"github.com/wso2/apk/common-controller/internal/operator/utils"
	"golang.org/x/exp/slices"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var c client.Client

// SetupWebhookWithManager creates a new webhook builder for API
func (r *API) SetupWebhookWithManager(mgr ctrl.Manager) error {
	c = mgr.GetClient()
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-dp-wso2-com-v1alpha1-api,mutating=true,failurePolicy=fail,sideEffects=None,groups=dp.wso2.com,resources=apis,verbs=create;update,versions=v1alpha1,name=mapi.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &API{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *API) Default() {
	// TODO: Add any defaulting logic here
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-dp-wso2-com-v1alpha1-api,mutating=false,failurePolicy=fail,sideEffects=None,groups=dp.wso2.com,resources=apis,verbs=create;update,versions=v1alpha1,name=vapi.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &API{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *API) ValidateCreate() (admission.Warnings, error) {
	return nil, r.validateAPI()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *API) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	return nil, r.validateAPI()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *API) ValidateDelete() (admission.Warnings, error) {

	// TODO(user): fill in your validation logic upon object deletion.
	return nil, nil
}

// validateAPI validate api crd fields
func (r *API) validateAPI() error {
	var allErrs field.ErrorList
	conf := config.ReadConfigs()
	namespaces := conf.CommonController.Operator.Namespaces
	if len(namespaces) > 0 {
		if !slices.Contains(namespaces, r.Namespace) {
			loggers.LoggerAPK.Debugf("API validation Skipped for namespace: %v", r.Namespace)
			return nil
		}
	}

	if r.Spec.BasePath == "" {
		allErrs = append(allErrs, field.Required(field.NewPath("spec").Child("basePath"), "API basePath is required"))
	} else if errMsg := validateAPIBasePathFormat(r.Spec.BasePath, r.Spec.APIVersion); errMsg != "" {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("basePath"), r.Spec.BasePath, errMsg))
	} else if err := r.validateAPIBasePathExistsAndDefaultVersion(); err != nil {
		allErrs = append(allErrs, err)
	}

	// Organization value should not be empty as it required when applying ratelimit policy
	if r.Spec.Organization == "" {
		allErrs = append(allErrs, field.Required(field.NewPath("spec").Child("organization"), "Organization can not be empty"))
	}

	if !(len(r.Spec.Production) > 0 && r.Spec.Production[0].HTTPRouteRefs != nil && len(r.Spec.Production[0].HTTPRouteRefs) > 0) && !(len(r.Spec.Sandbox) > 0 && r.Spec.Sandbox[0].HTTPRouteRefs != nil && len(r.Spec.Sandbox[0].HTTPRouteRefs) > 0) {
		allErrs = append(allErrs, field.Required(field.NewPath("spec"),
			"both API production and sandbox endpoint references cannot be empty"))
	}

	var prodHTTPRoute1, sandHTTPRoute1 []string
	if len(r.Spec.Production) > 0 {
		prodHTTPRoute1 = r.Spec.Production[0].HTTPRouteRefs
	}
	if len(r.Spec.Sandbox) > 0 {
		sandHTTPRoute1 = r.Spec.Sandbox[0].HTTPRouteRefs
	}

	if isEmptyStringsInArray(prodHTTPRoute1) {
		allErrs = append(allErrs, field.Required(field.NewPath("spec").Child("production").Child("httpRouteRefs"),
			"API production endpoint reference cannot be empty"))
	}

	if isEmptyStringsInArray(sandHTTPRoute1) {
		allErrs = append(allErrs, field.Required(field.NewPath("spec").Child("sandbox").Child("httpRouteRefs"),
			"API sandbox endpoint reference cannot be empty"))
	}

	if len(allErrs) > 0 {
		return apierrors.NewInvalid(
			schema.GroupKind{Group: "dp.wso2.com", Kind: "API"},
			r.Name, allErrs)
	}
	return nil
}

func isEmptyStringsInArray(strings []string) bool {
	for _, str := range strings {
		if str == "" {
			return true
		}
	}
	return false
}

func (r *API) validateAPIBasePathExistsAndDefaultVersion() *field.Error {

	apiList, err := retrieveAPIList()
	if err != nil {
		return field.InternalError(field.NewPath("spec").Child("basePath"),
			errors.New("unable to list APIs for API basePath validation"))

	}
	currentAPIBasePathWithoutVersion := getBasePathWithoutVersion(r.Spec.BasePath)
	incomingAPIEnvironment := utils.GetEnvironment(r.Spec.Environment)
	for _, api := range apiList {
		if (types.NamespacedName{Namespace: r.Namespace, Name: r.Name} !=
			types.NamespacedName{Namespace: api.Namespace, Name: api.Name}) {

			existingAPIEnvironment := utils.GetEnvironment(api.Spec.Environment)
			if api.Spec.Organization == r.Spec.Organization && api.Spec.BasePath == r.Spec.BasePath &&
				incomingAPIEnvironment == existingAPIEnvironment {
				return &field.Error{
					Type:     field.ErrorTypeDuplicate,
					Field:    field.NewPath("spec").Child("basePath").String(),
					BadValue: r.Spec.BasePath,
					Detail:   "an API has been already created for the basePath"}
			}
			if r.Spec.IsDefaultVersion {
				targetAPIBasePathWithoutVersion := getBasePathWithoutVersion(api.Spec.BasePath)
				targetAPIBasePathWithVersion := api.Spec.BasePath
				if api.Spec.IsDefaultVersion {
					if targetAPIBasePathWithoutVersion == currentAPIBasePathWithoutVersion {
						return &field.Error{
							Type:     field.ErrorTypeForbidden,
							Field:    field.NewPath("spec").Child("isDefaultVersion").String(),
							BadValue: r.Spec.BasePath,
							Detail:   "this API already has a default version"}
					}

				}
				if targetAPIBasePathWithVersion == currentAPIBasePathWithoutVersion {
					return &field.Error{
						Type:     field.ErrorTypeForbidden,
						Field:    field.NewPath("spec").Child("isDefaultVersion").String(),
						BadValue: r.Spec.BasePath,
						Detail:   fmt.Sprintf("api: %s's basePath path is colliding with default path", r.Name)}
				}
			}
		}
	}
	return nil
}

func retrieveAPIList() ([]API, error) {
	ctx := context.Background()
	conf := config.ReadConfigs()
	namespaces := conf.CommonController.Operator.Namespaces
	var apis []API
	if namespaces == nil {
		apiList := &APIList{}
		if err := c.List(ctx, apiList, &client.ListOptions{}); err != nil {
			loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2605, logging.CRITICAL, "Unable to list APIs: %v", err.Error()))
			return nil, err
		}
		apis = make([]API, len(apiList.Items))
		copy(apis[:], apiList.Items[:])
	} else {
		for _, namespace := range namespaces {
			apiList := &APIList{}
			if err := c.List(ctx, apiList, &client.ListOptions{Namespace: namespace}); err != nil {
				loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2605, logging.CRITICAL, "Unable to list APIs: %v", err.Error()))
				return nil, err
			}
			apis = append(apis, apiList.Items...)
		}
	}
	return apis, nil
}

func validateAPIBasePathFormat(basePath string, apiVersion string) string {
	if !strings.HasSuffix("/"+basePath, apiVersion) {
		return "API basePath value should contain the /{APIVersion} at end."
	}
	return ""
}

// getBasePathWithoutVersion returns the basePath without version
func getBasePathWithoutVersion(basePath string) string {
	lastIndex := strings.LastIndex(basePath, "/")
	if lastIndex != -1 {
		return basePath[:lastIndex]
	}
	return basePath
}
