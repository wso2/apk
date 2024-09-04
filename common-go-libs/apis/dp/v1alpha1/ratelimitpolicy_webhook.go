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

// SetupWebhookWithManager creates a new webhook builder for RateLimitPolicy
func (r *RateLimitPolicy) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-dp-wso2-com-v1alpha1-ratelimitpolicy,mutating=true,failurePolicy=fail,sideEffects=None,groups=dp.wso2.com,resources=ratelimitpolicies,verbs=create;update,versions=v1alpha1,name=mratelimitpolicy.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &RateLimitPolicy{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *RateLimitPolicy) Default() {
	// TODO(user): fill in your defaulting logic.
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-dp-wso2-com-v1alpha1-ratelimitpolicy,mutating=false,failurePolicy=fail,sideEffects=None,groups=dp.wso2.com,resources=ratelimitpolicies,verbs=create;update,versions=v1alpha1,name=vratelimitpolicy.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &RateLimitPolicy{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *RateLimitPolicy) ValidateCreate() (admission.Warnings, error) {

	// TODO(user): fill in your validation logic upon object creation.
	return nil, r.ValidatePolicies()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *RateLimitPolicy) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {

	// TODO(user): fill in your validation logic upon object update.
	return nil, r.ValidatePolicies()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *RateLimitPolicy) ValidateDelete() (admission.Warnings, error) {

	// TODO(user): fill in your validation logic upon object deletion.
	return nil, nil
}

// ValidatePolicies validates the policies in the RateLimitPolicy
func (r *RateLimitPolicy) ValidatePolicies() error {
	var allErrs field.ErrorList
	if r.Spec.TargetRef.Name == "" {
		allErrs = append(allErrs, field.Required(field.NewPath("spec").Child("targetRef").Child("name"),
			"Name is required"))
	}
	if !(r.Spec.TargetRef.Kind == constants.KindAPI || r.Spec.TargetRef.Kind == constants.KindResource ||
		r.Spec.TargetRef.Kind == constants.KindGateway || r.Spec.TargetRef.Kind == "Subscription") {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("targetRef").Child("kind"),
			r.Spec.TargetRef.Kind, "Invalid Kind is provided"))
	}
	if r.Spec.TargetRef.Namespace != nil && r.Spec.TargetRef.Namespace != (*v1beta1.Namespace)(&r.Namespace) {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("targetRef").Child("namespace"),
			r.Spec.TargetRef.Namespace, "namespace cross reference is not allowed"))
	}

	if len(allErrs) > 0 {
		return apierrors.NewInvalid(
			schema.GroupKind{Group: "dp.wso2.com", Kind: "RateLimitPolicy"},
			r.Name, allErrs)
	}
	return nil
}
