/*
 *  Copyright (c) 2024, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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

package v1alpha2

import (
	"strings"

	constants "github.com/wso2/apk/common-go-libs/constants"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// SetupWebhookWithManager creates a new webhook builder for Authentication CRD
func (r *Authentication) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-dp-wso2-com-v1alpha2-authentication,mutating=true,failurePolicy=fail,sideEffects=None,groups=dp.wso2.com,resources=authentications,verbs=create;update,versions=v1alpha2,name=mauthentication.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &Authentication{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Authentication) Default() {
	// TODO(user): fill in your defaulting logic.
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-dp-wso2-com-v1alpha2-authentication,mutating=false,failurePolicy=fail,sideEffects=None,groups=dp.wso2.com,resources=authentications,verbs=create;update,versions=v1alpha2,name=vauthentication.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &Authentication{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Authentication) ValidateCreate() (admission.Warnings, error) {
	// TODO(user): fill in your validation logic upon object creation.
	return nil, r.ValidateAuthentication()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Authentication) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	return nil, r.ValidateAuthentication()
}

// ValidateAuthentication validates the Authentication
func (r *Authentication) ValidateAuthentication() error {
	var allErrs field.ErrorList

	if r.Spec.TargetRef.Name == "" {
		allErrs = append(allErrs, field.Required(field.NewPath("spec").Child("targetRef").Child("name"), "Name is required"))
	}
	if !(r.Spec.TargetRef.Kind == constants.KindAPI || r.Spec.TargetRef.Kind == constants.KindResource) {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("targetRef").Child("kind"), r.Spec.TargetRef.Kind,
			"Invalid Kind is provided"))
	}

	var mutualSSL *MutualSSLConfig
	var oauth2Auth Oauth2Auth

	isOAuthDisabled := false
	isOAuthOptional := false

	isMTLSDisabled := false
	isMTLSMandatory := false

	if r.Spec.Default != nil && r.Spec.Default.AuthTypes != nil && r.Spec.Default.AuthTypes.MutualSSL != nil {
		oauth2Auth = r.Spec.Default.AuthTypes.Oauth2
		mutualSSL = r.Spec.Default.AuthTypes.MutualSSL

		isOAuthDisabled = oauth2Auth.Disabled
		isOAuthOptional = oauth2Auth.Required == "optional"

		isMTLSMandatory = strings.ToLower(mutualSSL.Required) == "mandatory"
		isMTLSDisabled = mutualSSL.Disabled

		if mutualSSL != nil && r.Spec.TargetRef.Kind != constants.KindAPI {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("default").Child("authTypes").Child("oauth2"), r.Spec.Default.AuthTypes.MutualSSL,
				"invalid authentication - mTLS can currently only be added for APIs"))
		}

		if (mutualSSL == nil || !isMTLSMandatory) && isOAuthDisabled {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("default").Child("authTypes").Child("mtls"), r.Spec.Default.AuthTypes,
				"invalid authentication configuration - security not enforced for API"))
		}

		if isMTLSDisabled && isOAuthOptional {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("default").Child("authTypes").Child("mtls"), r.Spec.Default.AuthTypes,
				"invalid authentication configuration - security not enforced for API"))
		}

		if mutualSSL != nil && len(mutualSSL.CertificatesInline) == 0 && len(mutualSSL.ConfigMapRefs) == 0 && len(mutualSSL.SecretRefs) == 0 {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("default").Child("authTypes").Child("mtls"), r.Spec.Default.AuthTypes.MutualSSL,
				"invalid mTLS configuration - certificates not provided"))
		}
	} else if r.Spec.Override != nil && r.Spec.Override.AuthTypes != nil && r.Spec.Override.AuthTypes.MutualSSL != nil {
		oauth2Auth := r.Spec.Override.AuthTypes.Oauth2
		mutualSSL = r.Spec.Override.AuthTypes.MutualSSL

		isOAuthDisabled = r.Spec.Override.AuthTypes.Oauth2.Disabled
		isOAuthOptional = oauth2Auth.Required == "optional"

		isMTLSMandatory = strings.ToLower(mutualSSL.Required) == "mandatory"
		isMTLSDisabled = mutualSSL.Disabled

		if mutualSSL != nil && r.Spec.TargetRef.Kind != constants.KindAPI {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("override").Child("authTypes").Child("oauth2"), r.Spec.Override.AuthTypes.MutualSSL,
				"invalid authentication - mTLS can currently only be added for APIs"))
		}

		if (mutualSSL == nil || !isMTLSMandatory) && isOAuthDisabled {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("override").Child("authTypes").Child("mtls"), r.Spec.Override.AuthTypes,
				"invalid authentication configuration - security not enforced for API"))
		}

		if isMTLSDisabled && isOAuthOptional {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("override").Child("authTypes").Child("mtls"), r.Spec.Override.AuthTypes,
				"invalid authentication configuration - security not enforced for API"))
		}

		if mutualSSL != nil && len(mutualSSL.CertificatesInline) == 0 && len(mutualSSL.ConfigMapRefs) == 0 && len(mutualSSL.SecretRefs) == 0 {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("override").Child("authTypes").Child("mtls"), r.Spec.Override.AuthTypes.MutualSSL,
				"invalid mTLS configuration - certificates not provided"))
		}
	}

	if len(allErrs) > 0 {
		return apierrors.NewInvalid(
			schema.GroupKind{Group: "dp.wso2.com", Kind: "Authentication"},
			r.Name, allErrs)
	}
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Authentication) ValidateDelete() (admission.Warnings, error) {
	// TODO(user): fill in your validation logic upon object deletion.
	return nil, nil
}
