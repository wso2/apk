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
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var authenticationlog = logf.Log.WithName("authentication-resource")

// SetupWebhookWithManager sets up and registers the webhook with the manager.
func (r *Authentication) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-dp-wso2-com-v1alpha1-authentication,mutating=true,failurePolicy=fail,sideEffects=None,groups=dp.wso2.com,resources=authentications,verbs=create;update,versions=v1alpha1,name=mauthentication.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &Authentication{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Authentication) Default() {
	authenticationlog.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-dp-wso2-com-v1alpha1-authentication,mutating=false,failurePolicy=fail,sideEffects=None,groups=dp.wso2.com,resources=authentications,verbs=create;update,versions=v1alpha1,name=vauthentication.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &Authentication{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Authentication) ValidateCreate() error {
	fmt.Println("Auth validate create")
	return r.validateAuthentication()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Authentication) ValidateUpdate(old runtime.Object) error {
	return r.validateAuthentication()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Authentication) ValidateDelete() error {
	authenticationlog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}

func (r *Authentication) validateAuthentication() error {
	var allErrs field.ErrorList
	var isMtlsDefault bool
	var mtlsOverride string
	if r.Spec.Default != nil && r.Spec.Default.ExternalService.AuthTypes != nil {
		mtlsDefault := r.Spec.Default.ExternalService.AuthTypes.MutualSSL
		if mtlsDefault != "" {
			isMtlsDefault = true
		}
	}

	if r.Spec.Override != nil && r.Spec.Override.ExternalService.AuthTypes != nil {
		mtlsOverride = r.Spec.Override.ExternalService.AuthTypes.MutualSSL
	}

	if r.Spec.Override != nil && r.Spec.Override.ExternalService.AuthTypes != nil && mtlsOverride == "" && !isMtlsDefault {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("override").Child("ext").Child("authTypes").Child("mutualSSL"),
			r.Spec.Override.ExternalService.AuthTypes.MutualSSL, "mutualSSL is mandatory when default is not set"))
	}

	if len(allErrs) > 0 {
		return apierrors.NewInvalid(
			schema.GroupKind{Group: "dp.wso2.com", Kind: "Authentication"},
			r.Name, allErrs)
	}
	return nil
}
