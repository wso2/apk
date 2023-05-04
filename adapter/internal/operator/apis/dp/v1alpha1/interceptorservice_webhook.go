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
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var interceptorservicelog = logf.Log.WithName("interceptorservice-resource")

// SetupWebhookWithManager creates a new webhook builder for InterceptorService
func (r *InterceptorService) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-dp-wso2-com-v1alpha1-interceptorservice,mutating=true,failurePolicy=fail,sideEffects=None,groups=dp.wso2.com,resources=interceptorservices,verbs=create;update,versions=v1alpha1,name=minterceptorservice.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &InterceptorService{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *InterceptorService) Default() {
	if r.Spec.BackendRef.Namespace == "" {
		r.Spec.BackendRef.Namespace = r.Namespace
	}
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-dp-wso2-com-v1alpha1-interceptorservice,mutating=false,failurePolicy=fail,sideEffects=None,groups=dp.wso2.com,resources=interceptorservices,verbs=create;update,versions=v1alpha1,name=vinterceptorservice.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &InterceptorService{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *InterceptorService) ValidateCreate() error {
	interceptorservicelog.Info("validate create", "name", r.Name)

	// TODO(user): fill in your validation logic upon object creation.
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *InterceptorService) ValidateUpdate(old runtime.Object) error {
	interceptorservicelog.Info("validate update", "name", r.Name)

	// TODO(user): fill in your validation logic upon object update.
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *InterceptorService) ValidateDelete() error {
	interceptorservicelog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}
