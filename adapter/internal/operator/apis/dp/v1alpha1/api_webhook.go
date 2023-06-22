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
	"regexp"
	"strings"

	"github.com/wso2/apk/adapter/config"
	"github.com/wso2/apk/adapter/internal/loggers"
	"github.com/wso2/apk/adapter/pkg/logging"
	"golang.org/x/exp/slices"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
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
func (r *API) ValidateCreate() error {
	return r.validateAPI()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *API) ValidateUpdate(old runtime.Object) error {
	return r.validateAPI()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *API) ValidateDelete() error {

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}

// validateAPI validate api crd fields
func (r *API) validateAPI() error {
	var allErrs field.ErrorList
	conf := config.ReadConfigs()
	namespaces := conf.Adapter.Operator.Namespaces
	if len(namespaces) > 0 {
		if !slices.Contains(namespaces, r.Namespace) {
			loggers.LoggerAPK.Debugf("API validation Skipped for namespace: %v", r.Namespace)
			return nil
		}
	}
	if r.Spec.APIDisplayName == "" {
		allErrs = append(allErrs, field.Required(field.NewPath("spec").Child("apiDisplayName"), "API display name is required"))
	} else if errMsg := validateAPIDisplayNameFormat(r.Spec.APIDisplayName); errMsg != "" {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("apiDisplayName"), r.Spec.APIDisplayName, errMsg))
	}

	if r.Spec.APIVersion == "" {
		allErrs = append(allErrs, field.Required(field.NewPath("spec").Child("apiVersion"), "API version is required"))
	} else if errMsg := validateAPIVersionFormat(r.Spec.APIVersion); errMsg != "" {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("apiVersion"), r.Spec.APIVersion, errMsg))
	}

	if r.Spec.Context == "" {
		allErrs = append(allErrs, field.Required(field.NewPath("spec").Child("context"), "API context is required"))
	} else if errMsg := validateAPIContextFormat(r.Spec.Context, r.Spec.APIVersion); errMsg != "" {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("context"), r.Spec.Context, errMsg))
	} else if err := r.validateAPIContextExists(); err != nil {
		allErrs = append(allErrs, err)
	}

	if errMsg := validateAPITypeFormat(r.Spec.APIType); errMsg != "" {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("apiType"), r.Spec.APIType, errMsg))
	} else {
		r.Spec.APIType = "REST"
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

func (r *API) validateAPIContextExists() *field.Error {

	apiList, err := retrieveAPIList()
	if err != nil {
		return field.InternalError(field.NewPath("spec").Child("context"),
			errors.New("unable to list APIs for API context validation"))

	}
	for _, api := range apiList {
		if (types.NamespacedName{Namespace: r.Namespace, Name: r.Name} !=
			types.NamespacedName{Namespace: api.Namespace, Name: api.Name}) && api.Spec.Organization == r.Spec.Organization && api.Spec.Context == r.Spec.Context {
			return &field.Error{
				Type:     field.ErrorTypeDuplicate,
				Field:    field.NewPath("spec").Child("context").String(),
				BadValue: r.Spec.Context,
				Detail:   "an API has been already created for the context"}
		}
	}
	return nil
}

func retrieveAPIList() ([]API, error) {
	ctx := context.Background()
	conf := config.ReadConfigs()
	namespaces := conf.Adapter.Operator.Namespaces
	var apis []API
	if namespaces == nil {
		apiList := &APIList{}
		if err := c.List(ctx, apiList, &client.ListOptions{}); err != nil {
			loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2605, err.Error()))
			return nil, err
		}
		apis = make([]API, len(apiList.Items))
		copy(apis[:], apiList.Items[:])
	} else {
		for _, namespace := range namespaces {
			apiList := &APIList{}
			if err := c.List(ctx, apiList, &client.ListOptions{Namespace: namespace}); err != nil {
				loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2605, err.Error()))
				return nil, err
			}
			apis = append(apis, apiList.Items...)
		}
	}
	return apis, nil
}

func validateAPIContextFormat(context string, apiVersion string) string {
	if len(context) > 232 {
		return "API context character length should not exceed 232."
	}
	if match, _ := regexp.MatchString("^[/][a-zA-Z0-9~/_.-]*$", context); !match {
		return "invalid API context. Does not start with / or includes invalid characters."
	}
	if !strings.HasSuffix("/"+context, apiVersion) {
		return "API context value should contain the /{APIVersion} at end."
	}
	return ""
}

func validateAPIDisplayNameFormat(apiName string) string {
	if len(apiName) > 60 {
		return "API display name character length should not exceed 60."
	}
	if match, _ := regexp.MatchString("^[^~!@#;:%^*()+={}|\\<>\"'',&$\\[\\]\\/]*$", apiName); !match {
		return "invalid API display name. Includes invalid characters."
	}
	return ""
}

func validateAPIVersionFormat(version string) string {
	if len(version) > 30 {
		return "API version length should not exceed 30."
	}
	if match, _ := regexp.MatchString("^[^~!@#;:%^*()+={}|\\<>\"'',&/$\\[\\]\\s+\\/]+$", version); !match {
		return "invalid API version. Includes invalid characters."
	}
	return ""
}

func validateAPITypeFormat(apiType string) string {
	if apiType != "" && strings.ToUpper(apiType) != "REST" {
		return "invalid API type. Only REST is supported"
	}
	return ""
}
