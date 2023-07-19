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
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var apipolicylog = logf.Log.WithName("apipolicy-resource")

// SetupWebhookWithManager creates a new webhook builder for APIPolicy
func (r *APIPolicy) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-dp-wso2-com-v1alpha1-apipolicy,mutating=true,failurePolicy=fail,sideEffects=None,groups=dp.wso2.com,resources=apipolicies,verbs=create;update,versions=v1alpha1,name=mapipolicy.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &APIPolicy{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *APIPolicy) Default() {
	if r.Spec.Override != nil {
		if len(r.Spec.Override.RequestInterceptors) > 0 {
			for i := range r.Spec.Override.RequestInterceptors {
				if len(r.Spec.Override.RequestInterceptors[i].Ref) == 0 {
					r.Spec.Override.RequestInterceptors[i].Ref = r.Name
				}
			}
		}
		if len(r.Spec.Override.ResponseInterceptors) > 0 {
			for i := range r.Spec.Override.ResponseInterceptors {
				if len(r.Spec.Override.ResponseInterceptors[i].Ref) == 0 {
					r.Spec.Override.ResponseInterceptors[i].Ref = r.Name
				}
			}
		}
	}
	if r.Spec.Default != nil {
		if len(r.Spec.Default.RequestInterceptors) > 0 {
			for i := range r.Spec.Default.RequestInterceptors {
				if len(r.Spec.Default.RequestInterceptors[i].Ref) == 0 {
					r.Spec.Default.RequestInterceptors[i].Ref = r.Name
				}
			}
		}
		if r.Spec.Default.ResponseInterceptors != nil {
			for i := range r.Spec.Default.ResponseInterceptors {
				if len(r.Spec.Default.ResponseInterceptors[i].Ref) == 0 {
					r.Spec.Default.ResponseInterceptors[i].Ref = r.Name
				}
			}
		}
	}
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-dp-wso2-com-v1alpha1-apipolicy,mutating=false,failurePolicy=fail,sideEffects=None,groups=dp.wso2.com,resources=apipolicies,verbs=create;update,versions=v1alpha1,name=vapipolicy.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &APIPolicy{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *APIPolicy) ValidateCreate() error {
	return r.validateJWTClaims()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *APIPolicy) ValidateUpdate(old runtime.Object) error {
	return r.validateJWTClaims()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *APIPolicy) ValidateDelete() error {
	apipolicylog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}

func (r *APIPolicy) validateJWTClaims() error {
	var allErrs field.ErrorList
	if r.Spec.Default != nil && r.Spec.Default.BackendJWTToken != nil {
		claims := r.Spec.Default.BackendJWTToken.CustomClaims
		for _, claim := range claims {
			valType := claim.Type
			switch valType {
			case "int":
				if _, err := strconv.Atoi(claim.Value); err != nil {
					allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("override").Child("backendJWTToken").Child("customClaims"), claim, "Provided value is not an integer"))
				}
			case "float":
				if _, err := strconv.ParseFloat(claim.Value, 64); err != nil {
					allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("override").Child("backendJWTToken").Child("customClaims"), claim, "Provided value is not a float"))
				}
			case "bool":
				if _, err := strconv.ParseBool(claim.Value); err != nil {
					allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("override").Child("backendJWTToken").Child("customClaims"), claim, "Provided value is not a boolean"))
				}
			case "date":
				if _, err := time.Parse(time.RFC3339, claim.Value); err != nil {
					allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("override").Child("backendJWTToken").Child("customClaims"), claim, "Provided value is not a date"))
				}
			case "long":
				if _, err := strconv.ParseInt(claim.Value, 10, 64); err != nil {
					allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("override").Child("backendJWTToken").Child("customClaims"), claim, "Provided value is not a long"))
				}
			}
		}
	}
	if r.Spec.Override != nil && r.Spec.Override.BackendJWTToken != nil {
		claims := r.Spec.Override.BackendJWTToken.CustomClaims
		for _, claim := range claims {
			valType := claim.Type
			switch valType {
			case "int":
				if _, err := strconv.Atoi(claim.Value); err != nil {
					allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("override").Child("backendJWTToken").Child("customClaims"), claim, "Provided value is not an integer"))
				}
			case "float":
				if _, err := strconv.ParseFloat(claim.Value, 64); err != nil {
					allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("override").Child("backendJWTToken").Child("customClaims"), claim, "Provided value is not a float"))
				}
			case "bool":
				if _, err := strconv.ParseBool(claim.Value); err != nil {
					allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("override").Child("backendJWTToken").Child("customClaims"), claim, "Provided value is not a boolean"))
				}
			case "date":
				if _, err := time.Parse(time.RFC3339, claim.Value); err != nil {
					allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("override").Child("backendJWTToken").Child("customClaims"), claim, "Provided value is not a date"))
				}
			case "long":
				if _, err := strconv.ParseInt(claim.Value, 10, 64); err != nil {
					allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("override").Child("backendJWTToken").Child("customClaims"), claim, "Provided value is not a long"))
				}
			}
		}
	}
	if len(allErrs) > 0 {
		return apierrors.NewInvalid(
			schema.GroupKind{Group: "dp.wso2.com", Kind: "APIPolicy"},
			r.Name, allErrs)
	}
	return nil
}
