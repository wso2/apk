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
	"strconv"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// SetupWebhookWithManager creates a new webhook builder for BackendJWT
func (r *BackendJWT) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-dp-wso2-com-v1alpha1-backendjwt,mutating=true,failurePolicy=fail,sideEffects=None,groups=dp.wso2.com,resources=backendjwts,verbs=create;update,versions=v1alpha1,name=mbackendjwt.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &BackendJWT{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *BackendJWT) Default() {
	// TODO(user): fill in your defaulting logic.
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-dp-wso2-com-v1alpha1-backendjwt,mutating=false,failurePolicy=fail,sideEffects=None,groups=dp.wso2.com,resources=backendjwts,verbs=create;update,versions=v1alpha1,name=vbackendjwt.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &BackendJWT{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *BackendJWT) ValidateCreate() (admission.Warnings, error) {
	// TODO(user): fill in your validation logic upon object creation.
	return nil, r.validateJWTClaims()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *BackendJWT) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	// TODO(user): fill in your validation logic upon object update.
	return nil, r.validateJWTClaims()
}

func (r *BackendJWT) validateJWTClaims() error {
	var allErrs field.ErrorList
	claims := r.Spec.CustomClaims
	for _, claim := range claims {
		valType := claim.Type
		switch valType {
		case "int":
			if _, err := strconv.Atoi(claim.Value); err != nil {
				allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("customClaims"), claim, "Provided value is not an integer"))
			}
		case "float":
			if _, err := strconv.ParseFloat(claim.Value, 64); err != nil {
				allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("customClaims"), claim, "Provided value is not a float"))
			}
		case "bool":
			if _, err := strconv.ParseBool(claim.Value); err != nil {
				allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("customClaims"), claim, "Provided value is not a boolean"))
			}
		case "date":
			if _, err := time.Parse(time.RFC3339, claim.Value); err != nil {
				allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("customClaims"), claim, "Provided value is not a date"))
			}
		case "long":
			if _, err := strconv.ParseInt(claim.Value, 10, 64); err != nil {
				allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("customClaims"), claim, "Provided value is not a long"))
			}
		}
	}

	if len(allErrs) > 0 {
		return apierrors.NewInvalid(
			schema.GroupKind{Group: "dp.wso2.com", Kind: "BackendJWT"},
			r.Name, allErrs)
	}
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *BackendJWT) ValidateDelete() (admission.Warnings, error) {
	// TODO(user): fill in your validation logic upon object deletion.
	return nil, nil
}
