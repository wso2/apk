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

	"github.com/wso2/apk/adapter/internal/loggers"
	"github.com/wso2/apk/adapter/pkg/logging"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
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
	loggers.LoggerAPKOperator.Infof("Validate API create: %s", r.Name)
	return r.validateAPI()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *API) ValidateUpdate(old runtime.Object) error {
	loggers.LoggerAPKOperator.Infof("Validate API update: %s", r.Name)
	return r.validateAPI()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *API) ValidateDelete() error {
	loggers.LoggerAPKOperator.Infof("Validate API delete: %s", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}

func (r *API) validateAPI() error {
	var allErrs field.ErrorList
	if err := r.validateAPIContext(); err != nil {
		allErrs = append(allErrs, err)
	}
	if len(allErrs) == 0 {
		return nil
	}

	return apierrors.NewInvalid(
		schema.GroupKind{Group: "dp.wso2.com", Kind: "API"},
		r.Name, allErrs)
}

func (r *API) validateAPIContext() *field.Error {
	ctx := context.Background()
	apiList := &APIList{}
	if err := c.List(ctx, apiList); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.ErrorDetails{
			Message:   fmt.Sprintf("Unable to list APIs for API context validation"),
			Severity:  logging.CRITICAL,
			ErrorCode: 2900,
		})
		return field.InternalError(field.NewPath("spec").Child("context"),
			errors.New("Unable to list APIs for API context validation"))
	}
	for _, api := range apiList.Items {
		if api.Spec.Context == r.Spec.Context {
			return field.Invalid(field.NewPath("spec").Child("context"),
				r.Spec.Context, fmt.Sprintf("an API has been already created for the context: %s", r.Spec.Context))
		}
	}
	return nil
}
