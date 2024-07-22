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
	// var oauth2Auth OAuth2Auth
	var authTypes *APIAuth

	isOAuthEnabled := true
	isOAuthMandatory := true

	isMTLSEnabled := false
	isMTLSMandatory := false

	isAPIKeyEnabled := false
	isAPIKeyMandatory := false
	errorType := "default"

	if r.Spec.Default != nil && r.Spec.Default.AuthTypes != nil {
		authTypes = r.Spec.Default.AuthTypes
	}

	if r.Spec.Override != nil && r.Spec.Override.AuthTypes != nil {
		authTypes = r.Spec.Override.AuthTypes
		errorType = "override"
	}

	isOAuthEnabled = !authTypes.OAuth2.Disabled
	isOAuthMandatory = authTypes.OAuth2.Required == "mandatory"

	if authTypes.MutualSSL != nil {
		mutualSSL = authTypes.MutualSSL
		isMTLSEnabled = !authTypes.MutualSSL.Disabled
		isMTLSMandatory = authTypes.MutualSSL.Required == "mandatory"
	}

	if authTypes.APIKey != nil {
		isAPIKeyEnabled = true
		isAPIKeyMandatory = authTypes.APIKey.Required == "mandatory"
	}

	if mutualSSL != nil && r.Spec.TargetRef.Kind != constants.KindAPI {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("default").Child("authTypes").Child("oauth2"), r.Spec.Default.AuthTypes.MutualSSL,
			"invalid authentication - mTLS can currently only be added for APIs"))
	}

	isMTLSMandatory = isMTLSEnabled && isMTLSMandatory       // false
	isOAuthMandatory = isOAuthEnabled && isOAuthMandatory    // true && true
	isAPIKeyMandatory = isAPIKeyEnabled && isAPIKeyMandatory // true && false = false

	isMTLSOptional := isMTLSEnabled && !isMTLSMandatory
	isOAuthOptional := isOAuthEnabled && !isOAuthMandatory
	isAPIKeyOptional := isAPIKeyEnabled && !isAPIKeyMandatory

	if !(
	// at least one must be enabled and mandatory
	(isMTLSMandatory || isOAuthMandatory || isAPIKeyMandatory) ||
		// mTLS is enabled and one of OAuth2 or APIKey is optional
		(isMTLSOptional && (isOAuthOptional || isAPIKeyOptional))) {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child(errorType).Child("authTypes"), r.Spec.Default.AuthTypes,
			"invalid authtypes provided: one of mTLS, APIKey, OAuth2 has to be enabled and mandatory "+
				"OR mTLS and one of OAuth2 or APIKey need to be optional "+
				"OR all three can be optional"))
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
