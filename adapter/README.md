# WSO2 APK - Adapter

This guide has information to setup adapter component for development and tasks for k8s operator development.

## Prerequisites

The following should be installed in your dev machine.

- [Gradle](https://gradle.org/install/) >= 7.5.1 version
- [Docker](https://docs.docker.com/engine/install/ubuntu/) >= 17.03 version
- [Golang](https://go.dev/doc/install) >= 1.23 version
- [Revive](https://github.com/mgechev/revive#installation) latest version
- [Kubebuilder](https://book.kubebuilder.io/quick-start.html#installation)

## Setting up

1. Clone the `wso2/apk` repository and navigate into adapter directory in home directory of the cloned tree.

    ```bash
    git clone https://github.com/wso2/apk.git
    cd adapter
    ```

2. Check whether you can build the project without running into any issues.

    ```
    gradle build
    ```

    This will build go binary and pack into a docker image named as `adapter:1.0.0-SNAPSHOT`.

3. If you ran into any issue first check whether the [prerequisites](#prerequisites) are satisfied.

## Operator

Since the adapter component uses Kubebuilder framework for operator development, when doing some tasks not listed below then first place to check is the [Kubebuilder documentation](https://book.kubebuilder.io/).

Code for the operator lies in `{PROJECT_HOME}/adapter/internal/operator`. This will be referred as `OPERATOR_HOME` in the upcoming sections.  

Following are some tasks with the steps that a developer might do in operator development:

- [Adding a new Kind](#adding-a-new-kind)
- [Adding a new property to an existing Kind](#adding-a-new-property-to-an-existing-kind)
- [Adding validating and defaulting logic](#adding-validating-and-defaulting-logic)

### Adding a new Kind

1. Decide what the k8s resource group will be depending on whether the CRD is for the control-plane or for the data-plane.

    | Plane  | k8s group |
    | ------------- | ------------- |
    | Data Plane  | dp  |
    | Control Plane  | cp  |

2. Decide the version for the CRD. Current version for all the CRDs are used as `v1alpha1`. 

3. Change directory into `OPERATOR_HOME`.

4. Rename OPERATOR_HOME/operator.go to main.go temperorily.

5. Let's say we want a new Kind called `APIPolicy` for data plane then run the following Kubebuilder command to scaffold out the new kind.

    ```bash
    kubebuilder create api --group dp --version v1alpha1 --kind APIPolicy
    ```

6. This will prompt for creating the resource. Respond yes for that since we need to generate the CRD for it.

    ```bash
    Create Resource [y/n]
    y
    ```

7. Next it will prompt for generating the boilerplate code for a controller, respond yes to it. As we are using a single controller in the current architecture. If your CR changes can be mapped to a `API` kind change event then you can delete the controller file. But, there might be cases you want a separate controller, then keep the generated controller file and add the code there.

    ```bash
    Create Controller [y/n]
    y
    ```

    Now new scaffold files/changes should be available in following directory structure:

    ```bash
    {OPERATOR_HOME}
    ├── PROJECT
    ├── apis
    │   ├── cp
    │   │   └── v1alpha1
    │   │       ├── ...
    │   │       └── ...
    │   └── dp
    │       └── v1alpha1
    |           ├── ...
    │           ├── apipolicy_types.go
    │           └── zz_generated.deepcopy.go
    .
    .
    .
    ├── controllers
    │   ├── cp
    │   │   ├── ...
    │   │   └── ...
    │   └── dp
    |       ├── ...
    │       └── apipolicy_controller.go
    ```
3. Revert the naming change did to internal/operator/operator.go in step 4.

9. The `apipolicy_types.go` contains the go struct representing the our example `APIPolicy` kind. You need to fill in the `APIPolicySpec` and `APIPolicyStatus` structs as per the needs.

    ```go
    // APIPolicySpec defines the desired state of APIPolicy
    type APIPolicySpec struct {
        // +kubebuilder:validation:MinLength=4
        Type      string                          `json:"type,omitempty"`
        ...
        ...
        TargetRef gwapiv1a2.PolicyTargetReference `json:"targetRef,omitempty"`
    }
    ```

    Here we have set the `Type` property to be required by adding `// +kubebuilder:validation:MinLength=4` [marker](https://book.kubebuilder.io/reference/markers/crd-validation.htm).

10. Since this example `APIPolicy` kind related to `API` kind, we can delete the `apipolicy_controller.go` file.

11. Adding the indexers: To filter out events to reconciliation loop, we need to index the `APIPolicy` resources in the operator in memory cache. Let's say we want `APIPolicy` resources to create a index using the targetRef section when the kind is `HTTPRoute` then the code for that will be as below. 

    > **NOTE**
    >
    > For this example this index is not used inside `getAPIsForAPIPolicy()` method and added here as an example.

    i. Declare the index name:
    ```go
    const httpRouteAPIPolicyIndex = "httpRouteAPIPolicyIndex"
    ```

    ii. Add the indexer code snippet inside the `addIndexers` function:

    ```go
    if err := mgr.GetFieldIndexer().IndexField(ctx, &dpv1alpha1.APIPolicy{}, httpRouteAPIPolicyIndex,
		func(rawObj k8client.Object) []string {
			apiPolicy := rawObj.(*dpv1alpha1.APIPolicy)
			var httpRoutes []string
			if apiPolicy.Spec.TargetRef.Kind == constants.KindAPI {
				httpRoutes = append(httpRoutes,
					types.NamespacedName{
						Namespace: apiPolicy.Namespace,
						Name:      string(apiPolicy.Spec.TargetRef.Name),
					}.String())
			}
			return httpRoutes
		}); err != nil {
		return err
	}
    ```

12. Adding event filtering handler: We need to implement `getAPIsForAPIPolicy()` function to filter out the `APIPolicy` changes as below.

    ```go
    func (apiReconciler *APIReconciler) getAPIsForAPIPolicy(obj k8client.Object) []reconcile.Request {
        ctx := context.Background()
        apiPolicy, ok := obj.(*dpv1alpha1.APIPolicy)
        if !ok {
            loggers.LoggerAPKOperator.ErrorC(logging.ErrorDetails{
                Message:   fmt.Sprintf("Unexpected object type, bypassing reconciliation: %v", apiPolicy),
                Severity:  logging.TRIVIAL,
                ErrorCode: 2670,
            })
            return []reconcile.Request{}
        }

        httpRoute := &gwapiv1.HTTPRoute{}
        if err := apiReconciler.client.Get(ctx, types.NamespacedName{
            Name: string(apiPolicy.Spec.TargetRef.Name),
            Namespace: utils.GetNamespace((*gwapiv1.Namespace)(apiPolicy.Spec.TargetRef.Namespace),
                apiPolicy.Namespace),
        }, httpRoute); err != nil {
            loggers.LoggerAPKOperator.ErrorC(logging.ErrorDetails{
                Message:   fmt.Sprintf("Unable to find associated HTTPRoutes for APIPolicy: %s", utils.NamespacedName(apiPolicy).String()),
                Severity:  logging.CRITICAL,
                ErrorCode: 2671,
            })
            return []reconcile.Request{}
        }

        requests := []reconcile.Request{}
        requests = append(requests, apiReconciler.getAPIForHTTPRoute(httpRoute)...)
        return requests
    }
    ```

13. Adding the watchers: Since the `APIPolicy` kind resource changes are fed into the Reconcile loop of `api_controller.go`, Add following code snippet under at the end of `NewAPIController` function.

    ```go
    if err := c.Watch(&source.Kind{Type: &dpv1alpha1.APIPolicy{}}, handler.EnqueueRequestsFromMapFunc(r.getAPIsForAPIPolicy),
        predicates...); err != nil {
        loggers.LoggerAPKOperator.ErrorC(logging.ErrorDetails{
            Message:   fmt.Sprintf("Error watching APIPolicy resources: %v", err),
            Severity:  logging.BLOCKER,
            ErrorCode: <give-new-error-code-here>,
        })
        return err
    }
    ```

14. Generating CRD and other resource yamls:

    ```bash
    gradle build
    ```

    This will generate artefacts inside `{OPERATOR_HOME}/config` directory.

15. To make the CRD and other resource changes affect, you need to move the k8s resources to the helm chart in `PROJECT_HOME/helm-charts` directory:

    - Append new rules to the `ClusterRole` in `helm-charts/templates/serviceAccount/apk-cluster-role.yaml`.

        ```yaml
        - apiGroups: ["dp.wso2.com"]
          resources: ["apipolicies"]
          verbs: ["get","list","watch","update","delete","create"]
        - apiGroups: ["dp.wso2.com"]
          resources: ["apipolicies/finalizers"]
          verbs: ["update"]
        - apiGroups: ["dp.wso2.com"]
          resources: ["apipolicies/status"]
          verbs: ["get","patch","update"]
        ```

### Adding a new property to an existing Kind

1. Add the new property in spec or status of the existing resource in `<resource>_types.go` file.

2. Add the logic inside the `<resource>_controller.go` file.

3. Follow the step `12` and step `13` to generate and move the changes of the CRDs and other resources. 

### Adding validating and defaulting logic

Other than the basic validations we can add using [kubebuilder markers](https://book.kubebuilder.io/reference/markers/crd-validation.html) (which are finally getting added in openapi schema section CRD yaml file). In some cases we need other validation cannot achieve using the markers. For example, cross resource validations like context property in `API` kind.

In that case, we can write the validating and defaulting logic by generating more scaffold code as described in [Implementing defaulting/validating webhooks](https://book.kubebuilder.io/cronjob-tutorial/webhook-implementation.html) section in kubebuilder docs.

Refer to this example [PR](https://github.com/wso2/apk/pull/370) for more information.
1. Rename internal/operator/operator.go to main.go

2. Create webhook resources. Example command would be similar to;

    ```bash
    kubebuilder create webhook --group dp --version v1alpha1 --kind APIPolicy --defaulting --programmatic-validation
    ```
3. Revert the naming change did to internal/operator/operator.go in step 1.

4. Copy `manifests.yaml` new entries to helm chart.

5. Add webhook setup to `operator.go`.

    ```go
    (&dpv1alpha1.APIPolicy{}).SetupWebhookWithManager(mgr)
    ```
