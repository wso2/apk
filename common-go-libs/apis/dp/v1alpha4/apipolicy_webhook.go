/*
 *  Copyright (c) 2025, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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

package v1alpha4

import (
	constants "github.com/wso2/apk/common-go-libs/constants"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"sigs.k8s.io/gateway-api/apis/v1beta1"
)

// SetupWebhookWithManager creates a new webhook builder for APIPolicy
func (r *APIPolicy) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-dp-wso2-com-v1alpha4-apipolicy,mutating=true,failurePolicy=fail,sideEffects=None,groups=dp.wso2.com,resources=apipolicies,verbs=create;update,versions=v1alpha4,name=mapipolicy.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &APIPolicy{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *APIPolicy) Default() {}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-dp-wso2-com-v1alpha4-apipolicy,mutating=false,failurePolicy=fail,sideEffects=None,groups=dp.wso2.com,resources=apipolicies,verbs=create;update,versions=v1alpha4,name=vapipolicy.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &APIPolicy{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *APIPolicy) ValidateCreate() (admission.Warnings, error) {
	return nil, r.ValidatePolicy()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *APIPolicy) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	return nil, r.ValidatePolicy()
}

// ValidatePolicy validates the APIPolicy
func (r *APIPolicy) ValidatePolicy() error {
	var allErrs field.ErrorList

	if r.Spec.TargetRef.Name == "" {
		allErrs = append(allErrs, field.Required(field.NewPath("spec").Child("targetRef").Child("name"), "Name is required"))
	}
	if !(r.Spec.TargetRef.Kind == constants.KindAPI || r.Spec.TargetRef.Kind == constants.KindResource ||
		r.Spec.TargetRef.Kind == constants.KindGateway) {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("targetRef").Child("kind"), r.Spec.TargetRef.Kind,
			"Invalid Kind is provided"))
	}
	if r.Spec.TargetRef.Namespace != nil && r.Spec.TargetRef.Namespace != (*v1beta1.Namespace)(&r.Namespace) {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("targetRef").Child("namespace"), r.Spec.TargetRef.Namespace,
			"namespace cross reference is not allowed"))
	}
	if len(allErrs) > 0 {
		return apierrors.NewInvalid(
			schema.GroupKind{Group: "dp.wso2.com", Kind: "APIPolicy"},
			r.Name, allErrs)
	}
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *APIPolicy) ValidateDelete() (admission.Warnings, error) {
	// TODO(user): fill in your validation logic upon object deletion.
	return nil, nil
}
