/*
 *  Copyright (c) 2023, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var backendlog = logf.Log.WithName("backend-resource")

// SetupWebhookWithManager sets up and registers the backend webhook with the manager.
func (r *Backend) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-dp-wso2-com-v1alpha1-backend,mutating=true,failurePolicy=fail,sideEffects=None,groups=dp.wso2.com,resources=backends,verbs=create;update,versions=v1alpha1,name=mbackend.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &Backend{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Backend) Default() {
	backendlog.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-dp-wso2-com-v1alpha1-backend,mutating=false,failurePolicy=fail,sideEffects=None,groups=dp.wso2.com,resources=backends,verbs=create;update,versions=v1alpha1,name=vbackend.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &Backend{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Backend) ValidateCreate() (admission.Warnings, error) {
	return nil, r.validateBackendSpec()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Backend) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	return nil, r.validateBackendSpec()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Backend) ValidateDelete() (admission.Warnings, error) {
	backendlog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil, nil
}

func (r *Backend) validateBackendSpec() error {
	var allErrs field.ErrorList
	retryConfig := r.Spec.Retry
	if retryConfig != nil {
		for _, statusCode := range retryConfig.StatusCodes {
			if statusCode > 598 || statusCode < 401 {
				allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("retry").Child("statusCodes"),
					retryConfig.StatusCodes, "status code should be between 401 and 598"))
			}
		}
	}

	if len(allErrs) > 0 {
		return apierrors.NewInvalid(
			schema.GroupKind{Group: "dp.wso2.com", Kind: "Backend"},
			r.Name, allErrs)
	}
	return nil
}
