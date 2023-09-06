<p>Packages:</p>
<ul>
<li>
<a href="#dp.wso2.com%2fv1alpha1">dp.wso2.com/v1alpha1</a>
</li>
</ul>
<h2 id="dp.wso2.com/v1alpha1">dp.wso2.com/v1alpha1</h2>
<p>
<p>Package v1alpha1 contains the API Schema definitions for WSO2 APK.</p>
</p>
Resource Types:
<ul><li>
<a href="#dp.wso2.com/v1alpha1.API">API</a>
</li><li>
<a href="#dp.wso2.com/v1alpha1.APIPolicy">APIPolicy</a>
</li><li>
<a href="#dp.wso2.com/v1alpha1.Authentication">Authentication</a>
</li><li>
<a href="#dp.wso2.com/v1alpha1.Backend">Backend</a>
</li><li>
<a href="#dp.wso2.com/v1alpha1.BackendJWT">BackendJWT</a>
</li><li>
<a href="#dp.wso2.com/v1alpha1.InterceptorService">InterceptorService</a>
</li><li>
<a href="#dp.wso2.com/v1alpha1.JWTIssuer">JWTIssuer</a>
</li><li>
<a href="#dp.wso2.com/v1alpha1.RateLimitPolicy">RateLimitPolicy</a>
</li><li>
<a href="#dp.wso2.com/v1alpha1.Scope">Scope</a>
</li></ul>
<h3 id="dp.wso2.com/v1alpha1.API">API
<a class="headerlink" href="#dp.wso2.com%2fv1alpha1.API" title="Permanent link">¶</a>
</h3>
<p>
<p>API is the Schema for the apis API</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>apiVersion</code></br>
string</td>
<td>
<code>
dp.wso2.com/v1alpha1
</code>
</td>
</tr>
<tr>
<td>
<code>kind</code></br>
string
</td>
<td><code>API</code></td>
</tr>
<tr>
<td>
<code>metadata</code></br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.23/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.APISpec">
APISpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>apiDisplayName</code></br>
<em>
string
</em>
</td>
<td>
<p>APIDisplayName is the unique name of the API in
the namespace defined. &ldquo;Namespace/APIDisplayName&rdquo; can
be used to uniquely identify an API.</p>
</td>
</tr>
<tr>
<td>
<code>apiVersion</code></br>
<em>
string
</em>
</td>
<td>
<p>APIVersion is the version number of the API.</p>
</td>
</tr>
<tr>
<td>
<code>isDefaultVersion</code></br>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>IsDefaultVersion indicates whether this API version should be used as a default API</p>
</td>
</tr>
<tr>
<td>
<code>definitionFileRef</code></br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>DefinitionFileRef contains the OpenAPI 3 or Swagger
definition of the API in a ConfigMap.</p>
</td>
</tr>
<tr>
<td>
<code>definitionPath</code></br>
<em>
string
</em>
</td>
<td>
<p>DefinitionPath contains the path to expose the API definition.</p>
</td>
</tr>
<tr>
<td>
<code>production</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.EnvConfig">
[]EnvConfig
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Production contains a list of references to HttpRoutes
of type HttpRoute.
xref: <a href="https://github.com/kubernetes-sigs/gateway-api/blob/main/apis/v1beta1/httproute_types.go">https://github.com/kubernetes-sigs/gateway-api/blob/main/apis/v1beta1/httproute_types.go</a></p>
</td>
</tr>
<tr>
<td>
<code>sandbox</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.EnvConfig">
[]EnvConfig
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Sandbox contains a list of references to HttpRoutes
of type HttpRoute.
xref: <a href="https://github.com/kubernetes-sigs/gateway-api/blob/main/apis/v1beta1/httproute_types.go">https://github.com/kubernetes-sigs/gateway-api/blob/main/apis/v1beta1/httproute_types.go</a></p>
</td>
</tr>
<tr>
<td>
<code>apiType</code></br>
<em>
string
</em>
</td>
<td>
<p>APIType denotes the type of the API.
Possible values could be REST, GraphQL, Async</p>
</td>
</tr>
<tr>
<td>
<code>context</code></br>
<em>
string
</em>
</td>
<td>
<p>Context denotes the context of the API.
e.g: /pet-store-api/1.0.6</p>
</td>
</tr>
<tr>
<td>
<code>organization</code></br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Organization denotes the organization.
related to the API</p>
</td>
</tr>
<tr>
<td>
<code>systemAPI</code></br>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>SystemAPI denotes if it is an internal system API.</p>
</td>
</tr>
<tr>
<td>
<code>apiProperties</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.Property">
[]Property
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>APIProperties denotes the custom properties of the API.</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.APIStatus">
APIStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="dp.wso2.com/v1alpha1.APIPolicy">APIPolicy
<a class="headerlink" href="#dp.wso2.com%2fv1alpha1.APIPolicy" title="Permanent link">¶</a>
</h3>
<p>
<p>APIPolicy is the Schema for the apipolicies API</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>apiVersion</code></br>
string</td>
<td>
<code>
dp.wso2.com/v1alpha1
</code>
</td>
</tr>
<tr>
<td>
<code>kind</code></br>
string
</td>
<td><code>APIPolicy</code></td>
</tr>
<tr>
<td>
<code>metadata</code></br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.23/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.APIPolicySpec">
APIPolicySpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>default</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.PolicySpec">
PolicySpec
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>override</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.PolicySpec">
PolicySpec
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>targetRef</code></br>
<em>
sigs.k8s.io/gateway-api/apis/v1alpha2.PolicyTargetReference
</em>
</td>
<td>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.APIPolicyStatus">
APIPolicyStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="dp.wso2.com/v1alpha1.Authentication">Authentication
<a class="headerlink" href="#dp.wso2.com%2fv1alpha1.Authentication" title="Permanent link">¶</a>
</h3>
<p>
<p>Authentication is the Schema for the authentications API</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>apiVersion</code></br>
string</td>
<td>
<code>
dp.wso2.com/v1alpha1
</code>
</td>
</tr>
<tr>
<td>
<code>kind</code></br>
string
</td>
<td><code>Authentication</code></td>
</tr>
<tr>
<td>
<code>metadata</code></br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.23/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.AuthenticationSpec">
AuthenticationSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>default</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.AuthSpec">
AuthSpec
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>override</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.AuthSpec">
AuthSpec
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>targetRef</code></br>
<em>
sigs.k8s.io/gateway-api/apis/v1alpha2.PolicyTargetReference
</em>
</td>
<td>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.AuthenticationStatus">
AuthenticationStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="dp.wso2.com/v1alpha1.Backend">Backend
<a class="headerlink" href="#dp.wso2.com%2fv1alpha1.Backend" title="Permanent link">¶</a>
</h3>
<p>
(<em>Appears on:</em>
<a href="#dp.wso2.com/v1alpha1.ResolvedBackend">ResolvedBackend</a>)
</p>
<p>
<p>Backend is the Schema for the backends API</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>apiVersion</code></br>
string</td>
<td>
<code>
dp.wso2.com/v1alpha1
</code>
</td>
</tr>
<tr>
<td>
<code>kind</code></br>
string
</td>
<td><code>Backend</code></td>
</tr>
<tr>
<td>
<code>metadata</code></br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.23/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.BackendSpec">
BackendSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>services</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.Service">
[]Service
</a>
</em>
</td>
<td>
<p>Services holds hosts and ports</p>
</td>
</tr>
<tr>
<td>
<code>protocol</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.BackendProtocolType">
BackendProtocolType
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Protocol defines the backend protocol</p>
</td>
</tr>
<tr>
<td>
<code>basePath</code></br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>BasePath defines the base path of the backend</p>
</td>
</tr>
<tr>
<td>
<code>tls</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.TLSConfig">
TLSConfig
</a>
</em>
</td>
<td>
<p>TLS defines the TLS configurations of the backend</p>
</td>
</tr>
<tr>
<td>
<code>security</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.SecurityConfig">
SecurityConfig
</a>
</em>
</td>
<td>
<p>Security defines the security configurations of the backend</p>
</td>
</tr>
<tr>
<td>
<code>circuitBreaker</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.CircuitBreaker">
CircuitBreaker
</a>
</em>
</td>
<td>
<p>CircuitBreaker defines the circuit breaker configurations</p>
</td>
</tr>
<tr>
<td>
<code>timeout</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.Timeout">
Timeout
</a>
</em>
</td>
<td>
<p>Timeout configuration for the backend</p>
</td>
</tr>
<tr>
<td>
<code>retry</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.RetryConfig">
RetryConfig
</a>
</em>
</td>
<td>
<p>Retry configuration for the backend</p>
</td>
</tr>
<tr>
<td>
<code>healthCheck</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.HealthCheck">
HealthCheck
</a>
</em>
</td>
<td>
<p>HealthCheck configuration for the backend tcp health check</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.BackendStatus">
BackendStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="dp.wso2.com/v1alpha1.BackendJWT">BackendJWT
<a class="headerlink" href="#dp.wso2.com%2fv1alpha1.BackendJWT" title="Permanent link">¶</a>
</h3>
<p>
<p>BackendJWT is the Schema for the backendjwts API</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>apiVersion</code></br>
string</td>
<td>
<code>
dp.wso2.com/v1alpha1
</code>
</td>
</tr>
<tr>
<td>
<code>kind</code></br>
string
</td>
<td><code>BackendJWT</code></td>
</tr>
<tr>
<td>
<code>metadata</code></br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.23/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.BackendJWTSpec">
BackendJWTSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>encoding</code></br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Encoding of the JWT token</p>
</td>
</tr>
<tr>
<td>
<code>header</code></br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Header of the JWT token</p>
</td>
</tr>
<tr>
<td>
<code>signingAlgorithm</code></br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Signing algorithm of the JWT token</p>
</td>
</tr>
<tr>
<td>
<code>tokenTTL</code></br>
<em>
uint32
</em>
</td>
<td>
<em>(Optional)</em>
<p>TokenTTL time to live for the backend JWT token in seconds</p>
</td>
</tr>
<tr>
<td>
<code>customClaims</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.CustomClaim">
[]CustomClaim
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>CustomClaims holds custom claims that needs to be added to the jwt</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.BackendJWTStatus">
BackendJWTStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="dp.wso2.com/v1alpha1.InterceptorService">InterceptorService
<a class="headerlink" href="#dp.wso2.com%2fv1alpha1.InterceptorService" title="Permanent link">¶</a>
</h3>
<p>
<p>InterceptorService is the Schema for the interceptorservices API</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>apiVersion</code></br>
string</td>
<td>
<code>
dp.wso2.com/v1alpha1
</code>
</td>
</tr>
<tr>
<td>
<code>kind</code></br>
string
</td>
<td><code>InterceptorService</code></td>
</tr>
<tr>
<td>
<code>metadata</code></br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.23/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.InterceptorServiceSpec">
InterceptorServiceSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>backendRef</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.BackendReference">
BackendReference
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>includes</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.InterceptorInclusion">
[]InterceptorInclusion
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Includes defines the types of data which should be included when calling the interceptor service</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.InterceptorServiceStatus">
InterceptorServiceStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="dp.wso2.com/v1alpha1.JWTIssuer">JWTIssuer
<a class="headerlink" href="#dp.wso2.com%2fv1alpha1.JWTIssuer" title="Permanent link">¶</a>
</h3>
<p>
<p>JWTIssuer is the Schema for the jwtissuers API</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>apiVersion</code></br>
string</td>
<td>
<code>
dp.wso2.com/v1alpha1
</code>
</td>
</tr>
<tr>
<td>
<code>kind</code></br>
string
</td>
<td><code>JWTIssuer</code></td>
</tr>
<tr>
<td>
<code>metadata</code></br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.23/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.JWTIssuerSpec">
JWTIssuerSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>name</code></br>
<em>
string
</em>
</td>
<td>
<p>Name is the unique name of the JWT Issuer in
the Organization defined . &ldquo;Organization/Name&rdquo; can
be used to uniquely identify an Issuer.</p>
</td>
</tr>
<tr>
<td>
<code>organization</code></br>
<em>
string
</em>
</td>
<td>
<p>Organization denotes the organization of the JWT Issuer.</p>
</td>
</tr>
<tr>
<td>
<code>issuer</code></br>
<em>
string
</em>
</td>
<td>
<p>Issuer denotes the issuer of the JWT Issuer.</p>
</td>
</tr>
<tr>
<td>
<code>consumerKeyClaim</code></br>
<em>
string
</em>
</td>
<td>
<p>ConsumerKeyClaim denotes the claim key of the consumer key.</p>
</td>
</tr>
<tr>
<td>
<code>scopesClaim</code></br>
<em>
string
</em>
</td>
<td>
<p>ScopesClaim denotes the claim key of the scopes.</p>
</td>
</tr>
<tr>
<td>
<code>signatureValidation</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.SignatureValidation">
SignatureValidation
</a>
</em>
</td>
<td>
<p>SignatureValidation denotes the signature validation method of jwt</p>
</td>
</tr>
<tr>
<td>
<code>claimMappings</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.[]github.com/wso2/apk/adapter/internal/operator/apis/dp/v1alpha1.ClaimMapping">
[]github.com/wso2/apk/adapter/internal/operator/apis/dp/v1alpha1.ClaimMapping
</a>
</em>
</td>
<td>
<p>ClaimMappings denotes the claim mappings of the jwt</p>
</td>
</tr>
<tr>
<td>
<code>targetRef</code></br>
<em>
sigs.k8s.io/gateway-api/apis/v1alpha2.PolicyTargetReference
</em>
</td>
<td>
<p>TargetRef denotes the reference to the which gateway it applies to</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.JWTIssuerStatus">
JWTIssuerStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="dp.wso2.com/v1alpha1.RateLimitPolicy">RateLimitPolicy
<a class="headerlink" href="#dp.wso2.com%2fv1alpha1.RateLimitPolicy" title="Permanent link">¶</a>
</h3>
<p>
<p>RateLimitPolicy is the Schema for the ratelimitpolicies API</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>apiVersion</code></br>
string</td>
<td>
<code>
dp.wso2.com/v1alpha1
</code>
</td>
</tr>
<tr>
<td>
<code>kind</code></br>
string
</td>
<td><code>RateLimitPolicy</code></td>
</tr>
<tr>
<td>
<code>metadata</code></br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.23/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.RateLimitPolicySpec">
RateLimitPolicySpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>default</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.RateLimitAPIPolicy">
RateLimitAPIPolicy
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>override</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.RateLimitAPIPolicy">
RateLimitAPIPolicy
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>targetRef</code></br>
<em>
sigs.k8s.io/gateway-api/apis/v1alpha2.PolicyTargetReference
</em>
</td>
<td>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.RateLimitPolicyStatus">
RateLimitPolicyStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="dp.wso2.com/v1alpha1.Scope">Scope
<a class="headerlink" href="#dp.wso2.com%2fv1alpha1.Scope" title="Permanent link">¶</a>
</h3>
<p>
<p>Scope is the Schema for the scopes API</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>apiVersion</code></br>
string</td>
<td>
<code>
dp.wso2.com/v1alpha1
</code>
</td>
</tr>
<tr>
<td>
<code>kind</code></br>
string
</td>
<td><code>Scope</code></td>
</tr>
<tr>
<td>
<code>metadata</code></br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.23/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.ScopeSpec">
ScopeSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>names</code></br>
<em>
[]string
</em>
</td>
<td>
<p>Name scope name</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.ScopeStatus">
ScopeStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="dp.wso2.com/v1alpha1.APIAuth">APIAuth
<a class="headerlink" href="#dp.wso2.com%2fv1alpha1.APIAuth" title="Permanent link">¶</a>
</h3>
<p>
(<em>Appears on:</em>
<a href="#dp.wso2.com/v1alpha1.AuthSpec">AuthSpec</a>)
</p>
<p>
<p>APIAuth Authentication scheme type and details</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>jwt</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.JWTAuth">
JWTAuth
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>JWT is to specify the JWT authentication scheme details</p>
</td>
</tr>
<tr>
<td>
<code>apiKey</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.APIKeyAuth">
[]APIKeyAuth
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>APIKey is to specify the APIKey authentication scheme details</p>
</td>
</tr>
<tr>
<td>
<code>testConsoleKey</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.TestConsoleKeyAuth">
TestConsoleKeyAuth
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>TestConsoleKey is to specify the Test Console Key authentication scheme details</p>
</td>
</tr>
</tbody>
</table>
<h3 id="dp.wso2.com/v1alpha1.APIKeyAuth">APIKeyAuth
<a class="headerlink" href="#dp.wso2.com%2fv1alpha1.APIKeyAuth" title="Permanent link">¶</a>
</h3>
<p>
(<em>Appears on:</em>
<a href="#dp.wso2.com/v1alpha1.APIAuth">APIAuth</a>)
</p>
<p>
<p>APIKeyAuth APIKey Authentication scheme details</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>in</code></br>
<em>
string
</em>
</td>
<td>
<pre><code>In is to specify how the APIKey is passed to the request
</code></pre>
</td>
</tr>
<tr>
<td>
<code>name</code></br>
<em>
string
</em>
</td>
<td>
<p>Name is the name of the header or query parameter to be used</p>
</td>
</tr>
<tr>
<td>
<code>sendTokenToUpstream</code></br>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>SendTokenToUpstream is to specify whether the APIKey should be sent to the upstream</p>
</td>
</tr>
</tbody>
</table>
<h3 id="dp.wso2.com/v1alpha1.APIPolicySpec">APIPolicySpec
<a class="headerlink" href="#dp.wso2.com%2fv1alpha1.APIPolicySpec" title="Permanent link">¶</a>
</h3>
<p>
(<em>Appears on:</em>
<a href="#dp.wso2.com/v1alpha1.APIPolicy">APIPolicy</a>)
</p>
<p>
<p>APIPolicySpec defines the desired state of APIPolicy</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>default</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.PolicySpec">
PolicySpec
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>override</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.PolicySpec">
PolicySpec
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>targetRef</code></br>
<em>
sigs.k8s.io/gateway-api/apis/v1alpha2.PolicyTargetReference
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="dp.wso2.com/v1alpha1.APIPolicyStatus">APIPolicyStatus
<a class="headerlink" href="#dp.wso2.com%2fv1alpha1.APIPolicyStatus" title="Permanent link">¶</a>
</h3>
<p>
(<em>Appears on:</em>
<a href="#dp.wso2.com/v1alpha1.APIPolicy">APIPolicy</a>)
</p>
<p>
<p>APIPolicyStatus defines the observed state of APIPolicy</p>
</p>
<h3 id="dp.wso2.com/v1alpha1.APIRateLimitPolicy">APIRateLimitPolicy
<a class="headerlink" href="#dp.wso2.com%2fv1alpha1.APIRateLimitPolicy" title="Permanent link">¶</a>
</h3>
<p>
(<em>Appears on:</em>
<a href="#dp.wso2.com/v1alpha1.RateLimitAPIPolicy">RateLimitAPIPolicy</a>)
</p>
<p>
<p>APIRateLimitPolicy defines the desired state of APIPolicy</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>requestsPerUnit</code></br>
<em>
uint32
</em>
</td>
<td>
<p>RequestPerUnit is the number of requests allowed per unit time</p>
</td>
</tr>
<tr>
<td>
<code>unit</code></br>
<em>
string
</em>
</td>
<td>
<p>Unit is the unit of the requestsPerUnit</p>
</td>
</tr>
</tbody>
</table>
<h3 id="dp.wso2.com/v1alpha1.APISpec">APISpec
<a class="headerlink" href="#dp.wso2.com%2fv1alpha1.APISpec" title="Permanent link">¶</a>
</h3>
<p>
(<em>Appears on:</em>
<a href="#dp.wso2.com/v1alpha1.API">API</a>)
</p>
<p>
<p>APISpec defines the desired state of API</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>apiDisplayName</code></br>
<em>
string
</em>
</td>
<td>
<p>APIDisplayName is the unique name of the API in
the namespace defined. &ldquo;Namespace/APIDisplayName&rdquo; can
be used to uniquely identify an API.</p>
</td>
</tr>
<tr>
<td>
<code>apiVersion</code></br>
<em>
string
</em>
</td>
<td>
<p>APIVersion is the version number of the API.</p>
</td>
</tr>
<tr>
<td>
<code>isDefaultVersion</code></br>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>IsDefaultVersion indicates whether this API version should be used as a default API</p>
</td>
</tr>
<tr>
<td>
<code>definitionFileRef</code></br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>DefinitionFileRef contains the OpenAPI 3 or Swagger
definition of the API in a ConfigMap.</p>
</td>
</tr>
<tr>
<td>
<code>definitionPath</code></br>
<em>
string
</em>
</td>
<td>
<p>DefinitionPath contains the path to expose the API definition.</p>
</td>
</tr>
<tr>
<td>
<code>production</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.EnvConfig">
[]EnvConfig
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Production contains a list of references to HttpRoutes
of type HttpRoute.
xref: <a href="https://github.com/kubernetes-sigs/gateway-api/blob/main/apis/v1beta1/httproute_types.go">https://github.com/kubernetes-sigs/gateway-api/blob/main/apis/v1beta1/httproute_types.go</a></p>
</td>
</tr>
<tr>
<td>
<code>sandbox</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.EnvConfig">
[]EnvConfig
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Sandbox contains a list of references to HttpRoutes
of type HttpRoute.
xref: <a href="https://github.com/kubernetes-sigs/gateway-api/blob/main/apis/v1beta1/httproute_types.go">https://github.com/kubernetes-sigs/gateway-api/blob/main/apis/v1beta1/httproute_types.go</a></p>
</td>
</tr>
<tr>
<td>
<code>apiType</code></br>
<em>
string
</em>
</td>
<td>
<p>APIType denotes the type of the API.
Possible values could be REST, GraphQL, Async</p>
</td>
</tr>
<tr>
<td>
<code>context</code></br>
<em>
string
</em>
</td>
<td>
<p>Context denotes the context of the API.
e.g: /pet-store-api/1.0.6</p>
</td>
</tr>
<tr>
<td>
<code>organization</code></br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Organization denotes the organization.
related to the API</p>
</td>
</tr>
<tr>
<td>
<code>systemAPI</code></br>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>SystemAPI denotes if it is an internal system API.</p>
</td>
</tr>
<tr>
<td>
<code>apiProperties</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.Property">
[]Property
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>APIProperties denotes the custom properties of the API.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="dp.wso2.com/v1alpha1.APIStatus">APIStatus
<a class="headerlink" href="#dp.wso2.com%2fv1alpha1.APIStatus" title="Permanent link">¶</a>
</h3>
<p>
(<em>Appears on:</em>
<a href="#dp.wso2.com/v1alpha1.API">API</a>)
</p>
<p>
<p>APIStatus defines the observed state of API</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>deploymentStatus</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.DeploymentStatus">
DeploymentStatus
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>DeploymentStatus denotes the deployment status of the API</p>
</td>
</tr>
</tbody>
</table>
<h3 id="dp.wso2.com/v1alpha1.AuthSpec">AuthSpec
<a class="headerlink" href="#dp.wso2.com%2fv1alpha1.AuthSpec" title="Permanent link">¶</a>
</h3>
<p>
(<em>Appears on:</em>
<a href="#dp.wso2.com/v1alpha1.AuthenticationSpec">AuthenticationSpec</a>)
</p>
<p>
<p>AuthSpec specification of the authentication service</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>disabled</code></br>
<em>
bool
</em>
</td>
<td>
<p>Disabled is to disable all authentications</p>
</td>
</tr>
<tr>
<td>
<code>authTypes</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.APIAuth">
APIAuth
</a>
</em>
</td>
<td>
<p>AuthTypes is to specify the authentication scheme types and details</p>
</td>
</tr>
</tbody>
</table>
<h3 id="dp.wso2.com/v1alpha1.AuthenticationSpec">AuthenticationSpec
<a class="headerlink" href="#dp.wso2.com%2fv1alpha1.AuthenticationSpec" title="Permanent link">¶</a>
</h3>
<p>
(<em>Appears on:</em>
<a href="#dp.wso2.com/v1alpha1.Authentication">Authentication</a>)
</p>
<p>
<p>AuthenticationSpec defines the desired state of Authentication</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>default</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.AuthSpec">
AuthSpec
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>override</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.AuthSpec">
AuthSpec
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>targetRef</code></br>
<em>
sigs.k8s.io/gateway-api/apis/v1alpha2.PolicyTargetReference
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="dp.wso2.com/v1alpha1.AuthenticationStatus">AuthenticationStatus
<a class="headerlink" href="#dp.wso2.com%2fv1alpha1.AuthenticationStatus" title="Permanent link">¶</a>
</h3>
<p>
(<em>Appears on:</em>
<a href="#dp.wso2.com/v1alpha1.Authentication">Authentication</a>)
</p>
<p>
<p>AuthenticationStatus defines the observed state of Authentication</p>
</p>
<h3 id="dp.wso2.com/v1alpha1.BackendJWTSpec">BackendJWTSpec
<a class="headerlink" href="#dp.wso2.com%2fv1alpha1.BackendJWTSpec" title="Permanent link">¶</a>
</h3>
<p>
(<em>Appears on:</em>
<a href="#dp.wso2.com/v1alpha1.BackendJWT">BackendJWT</a>)
</p>
<p>
<p>BackendJWTSpec defines the desired state of BackendJWT</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>encoding</code></br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Encoding of the JWT token</p>
</td>
</tr>
<tr>
<td>
<code>header</code></br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Header of the JWT token</p>
</td>
</tr>
<tr>
<td>
<code>signingAlgorithm</code></br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Signing algorithm of the JWT token</p>
</td>
</tr>
<tr>
<td>
<code>tokenTTL</code></br>
<em>
uint32
</em>
</td>
<td>
<em>(Optional)</em>
<p>TokenTTL time to live for the backend JWT token in seconds</p>
</td>
</tr>
<tr>
<td>
<code>customClaims</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.CustomClaim">
[]CustomClaim
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>CustomClaims holds custom claims that needs to be added to the jwt</p>
</td>
</tr>
</tbody>
</table>
<h3 id="dp.wso2.com/v1alpha1.BackendJWTStatus">BackendJWTStatus
<a class="headerlink" href="#dp.wso2.com%2fv1alpha1.BackendJWTStatus" title="Permanent link">¶</a>
</h3>
<p>
(<em>Appears on:</em>
<a href="#dp.wso2.com/v1alpha1.BackendJWT">BackendJWT</a>)
</p>
<p>
<p>BackendJWTStatus defines the observed state of BackendJWT</p>
</p>
<h3 id="dp.wso2.com/v1alpha1.BackendJWTToken">BackendJWTToken
<a class="headerlink" href="#dp.wso2.com%2fv1alpha1.BackendJWTToken" title="Permanent link">¶</a>
</h3>
<p>
(<em>Appears on:</em>
<a href="#dp.wso2.com/v1alpha1.PolicySpec">PolicySpec</a>)
</p>
<p>
<p>BackendJWTToken holds backend JWT token information</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code></br>
<em>
string
</em>
</td>
<td>
<p>Name holds the name of the BackendJWT resource.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="dp.wso2.com/v1alpha1.BackendProtocolType">BackendProtocolType
(<code>string</code> alias)</p><a class="headerlink" href="#dp.wso2.com%2fv1alpha1.BackendProtocolType" title="Permanent link">¶</a>
</h3>
<p>
(<em>Appears on:</em>
<a href="#dp.wso2.com/v1alpha1.BackendSpec">BackendSpec</a>, 
<a href="#dp.wso2.com/v1alpha1.ResolvedBackend">ResolvedBackend</a>)
</p>
<p>
<p>BackendProtocolType defines the backend protocol type.</p>
</p>
<table>
<thead>
<tr>
<th>Value</th>
<th>Description</th>
</tr>
</thead>
<tbody><tr><td><p>&#34;http&#34;</p></td>
<td><p>HTTPProtocol is the http protocol</p>
</td>
</tr><tr><td><p>&#34;https&#34;</p></td>
<td><p>HTTPSProtocol is the https protocol</p>
</td>
</tr><tr><td><p>&#34;ws&#34;</p></td>
<td><p>WSProtocol is the ws protocol</p>
</td>
</tr><tr><td><p>&#34;wss&#34;</p></td>
<td><p>WSSProtocol is the wss protocol</p>
</td>
</tr></tbody>
</table>
<h3 id="dp.wso2.com/v1alpha1.BackendReference">BackendReference
<a class="headerlink" href="#dp.wso2.com%2fv1alpha1.BackendReference" title="Permanent link">¶</a>
</h3>
<p>
(<em>Appears on:</em>
<a href="#dp.wso2.com/v1alpha1.InterceptorServiceSpec">InterceptorServiceSpec</a>)
</p>
<p>
<p>BackendReference refers to a Backend resource as the interceptor service.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code></br>
<em>
string
</em>
</td>
<td>
<p>Name is the name of the Backend resource.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="dp.wso2.com/v1alpha1.BackendSpec">BackendSpec
<a class="headerlink" href="#dp.wso2.com%2fv1alpha1.BackendSpec" title="Permanent link">¶</a>
</h3>
<p>
(<em>Appears on:</em>
<a href="#dp.wso2.com/v1alpha1.Backend">Backend</a>)
</p>
<p>
<p>BackendSpec defines the desired state of Backend</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>services</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.Service">
[]Service
</a>
</em>
</td>
<td>
<p>Services holds hosts and ports</p>
</td>
</tr>
<tr>
<td>
<code>protocol</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.BackendProtocolType">
BackendProtocolType
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Protocol defines the backend protocol</p>
</td>
</tr>
<tr>
<td>
<code>basePath</code></br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>BasePath defines the base path of the backend</p>
</td>
</tr>
<tr>
<td>
<code>tls</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.TLSConfig">
TLSConfig
</a>
</em>
</td>
<td>
<p>TLS defines the TLS configurations of the backend</p>
</td>
</tr>
<tr>
<td>
<code>security</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.SecurityConfig">
SecurityConfig
</a>
</em>
</td>
<td>
<p>Security defines the security configurations of the backend</p>
</td>
</tr>
<tr>
<td>
<code>circuitBreaker</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.CircuitBreaker">
CircuitBreaker
</a>
</em>
</td>
<td>
<p>CircuitBreaker defines the circuit breaker configurations</p>
</td>
</tr>
<tr>
<td>
<code>timeout</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.Timeout">
Timeout
</a>
</em>
</td>
<td>
<p>Timeout configuration for the backend</p>
</td>
</tr>
<tr>
<td>
<code>retry</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.RetryConfig">
RetryConfig
</a>
</em>
</td>
<td>
<p>Retry configuration for the backend</p>
</td>
</tr>
<tr>
<td>
<code>healthCheck</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.HealthCheck">
HealthCheck
</a>
</em>
</td>
<td>
<p>HealthCheck configuration for the backend tcp health check</p>
</td>
</tr>
</tbody>
</table>
<h3 id="dp.wso2.com/v1alpha1.BackendStatus">BackendStatus
<a class="headerlink" href="#dp.wso2.com%2fv1alpha1.BackendStatus" title="Permanent link">¶</a>
</h3>
<p>
(<em>Appears on:</em>
<a href="#dp.wso2.com/v1alpha1.Backend">Backend</a>)
</p>
<p>
<p>BackendStatus defines the observed state of Backend</p>
</p>
<h3 id="dp.wso2.com/v1alpha1.BasicSecurityConfig">BasicSecurityConfig
<a class="headerlink" href="#dp.wso2.com%2fv1alpha1.BasicSecurityConfig" title="Permanent link">¶</a>
</h3>
<p>
(<em>Appears on:</em>
<a href="#dp.wso2.com/v1alpha1.SecurityConfig">SecurityConfig</a>)
</p>
<p>
<p>BasicSecurityConfig defines basic security configurations</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>secretRef</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.SecretRef">
SecretRef
</a>
</em>
</td>
<td>
<p>SecretRef to credentials</p>
</td>
</tr>
</tbody>
</table>
<h3 id="dp.wso2.com/v1alpha1.CERTConfig">CERTConfig
<a class="headerlink" href="#dp.wso2.com%2fv1alpha1.CERTConfig" title="Permanent link">¶</a>
</h3>
<p>
(<em>Appears on:</em>
<a href="#dp.wso2.com/v1alpha1.JWKS">JWKS</a>, 
<a href="#dp.wso2.com/v1alpha1.SignatureValidation">SignatureValidation</a>)
</p>
<p>
<p>CERTConfig defines the certificate configuration</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>certificateInline</code></br>
<em>
string
</em>
</td>
<td>
<p>CertificateInline is the Inline Certificate entry</p>
</td>
</tr>
<tr>
<td>
<code>secretRef</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.RefConfig">
RefConfig
</a>
</em>
</td>
<td>
<p>SecretRef denotes the reference to the Secret that contains the Certificate</p>
</td>
</tr>
<tr>
<td>
<code>configMapRef</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.RefConfig">
RefConfig
</a>
</em>
</td>
<td>
<p>ConfigMapRef denotes the reference to the ConfigMap that contains the Certificate</p>
</td>
</tr>
</tbody>
</table>
<h3 id="dp.wso2.com/v1alpha1.CORSPolicy">CORSPolicy
<a class="headerlink" href="#dp.wso2.com%2fv1alpha1.CORSPolicy" title="Permanent link">¶</a>
</h3>
<p>
(<em>Appears on:</em>
<a href="#dp.wso2.com/v1alpha1.PolicySpec">PolicySpec</a>)
</p>
<p>
<p>CORSPolicy holds CORS policy information</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>accessControlAllowCredentials</code></br>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>AllowCredentials indicates whether the request can include user credentials like
cookies, HTTP authentication or client side SSL certificates.</p>
</td>
</tr>
<tr>
<td>
<code>accessControlAllowHeaders</code></br>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>AccessControlAllowHeaders indicates which headers can be used
during the actual request.</p>
</td>
</tr>
<tr>
<td>
<code>accessControlAllowMethods</code></br>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>AccessControlAllowMethods indicates which methods can be used
during the actual request.</p>
</td>
</tr>
<tr>
<td>
<code>accessControlAllowOrigins</code></br>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>AccessControlAllowOrigins indicates which origins can be used
during the actual request.</p>
</td>
</tr>
<tr>
<td>
<code>accessControlExposeHeaders</code></br>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>AccessControlExposeHeaders indicates which headers can be exposed
as part of the response by listing their names.</p>
</td>
</tr>
<tr>
<td>
<code>accessControlMaxAge</code></br>
<em>
int
</em>
</td>
<td>
<em>(Optional)</em>
<p>AccessControlMaxAge indicates how long the results of a preflight request
can be cached in a preflight result cache.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="dp.wso2.com/v1alpha1.CircuitBreaker">CircuitBreaker
<a class="headerlink" href="#dp.wso2.com%2fv1alpha1.CircuitBreaker" title="Permanent link">¶</a>
</h3>
<p>
(<em>Appears on:</em>
<a href="#dp.wso2.com/v1alpha1.BackendSpec">BackendSpec</a>, 
<a href="#dp.wso2.com/v1alpha1.ResolvedBackend">ResolvedBackend</a>)
</p>
<p>
<p>CircuitBreaker defines the circuit breaker configurations</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>maxConnections</code></br>
<em>
uint32
</em>
</td>
<td>
<em>(Optional)</em>
<p>MaxConnections is the maximum number of connections that will make to the upstream cluster.</p>
</td>
</tr>
<tr>
<td>
<code>maxPendingRequests</code></br>
<em>
uint32
</em>
</td>
<td>
<em>(Optional)</em>
<p>MaxPendingRequests is the maximum number of pending requests that will allow to the upstream cluster.</p>
</td>
</tr>
<tr>
<td>
<code>maxRequests</code></br>
<em>
uint32
</em>
</td>
<td>
<em>(Optional)</em>
<p>MaxRequests is the maximum number of parallel requests that will make to the upstream cluster.</p>
</td>
</tr>
<tr>
<td>
<code>maxRetries</code></br>
<em>
uint32
</em>
</td>
<td>
<em>(Optional)</em>
<p>MaxRetries is the maximum number of parallel retries that will allow to the upstream cluster.</p>
</td>
</tr>
<tr>
<td>
<code>maxConnectionPools</code></br>
<em>
uint32
</em>
</td>
<td>
<em>(Optional)</em>
<p>MaxConnectionPools is the maximum number of parallel connection pools that will allow to the upstream cluster.
If not specified, the default is unlimited.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="dp.wso2.com/v1alpha1.ClaimMapping">ClaimMapping
<a class="headerlink" href="#dp.wso2.com%2fv1alpha1.ClaimMapping" title="Permanent link">¶</a>
</h3>
<p>
<p>ClaimMapping defines the reference configuration</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>remoteClaim</code></br>
<em>
string
</em>
</td>
<td>
<p>RemoteClaim denotes the remote claim</p>
</td>
</tr>
<tr>
<td>
<code>localClaim</code></br>
<em>
string
</em>
</td>
<td>
<p>LocalClaim denotes the local claim</p>
</td>
</tr>
</tbody>
</table>
<h3 id="dp.wso2.com/v1alpha1.CustomClaim">CustomClaim
<a class="headerlink" href="#dp.wso2.com%2fv1alpha1.CustomClaim" title="Permanent link">¶</a>
</h3>
<p>
(<em>Appears on:</em>
<a href="#dp.wso2.com/v1alpha1.BackendJWTSpec">BackendJWTSpec</a>)
</p>
<p>
<p>CustomClaim holds custom claim information</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>claim</code></br>
<em>
string
</em>
</td>
<td>
<p>Claim name</p>
</td>
</tr>
<tr>
<td>
<code>value</code></br>
<em>
string
</em>
</td>
<td>
<p>Claim value</p>
</td>
</tr>
<tr>
<td>
<code>type</code></br>
<em>
string
</em>
</td>
<td>
<p>Claim type</p>
</td>
</tr>
</tbody>
</table>
<h3 id="dp.wso2.com/v1alpha1.CustomRateLimitPolicy">CustomRateLimitPolicy
<a class="headerlink" href="#dp.wso2.com%2fv1alpha1.CustomRateLimitPolicy" title="Permanent link">¶</a>
</h3>
<p>
(<em>Appears on:</em>
<a href="#dp.wso2.com/v1alpha1.RateLimitAPIPolicy">RateLimitAPIPolicy</a>)
</p>
<p>
<p>CustomRateLimitPolicy defines the desired state of CustomPolicy</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>requestsPerUnit</code></br>
<em>
uint32
</em>
</td>
<td>
<p>RequestPerUnit is the number of requests allowed per unit time</p>
</td>
</tr>
<tr>
<td>
<code>unit</code></br>
<em>
string
</em>
</td>
<td>
<p>Unit is the unit of the requestsPerUnit</p>
</td>
</tr>
<tr>
<td>
<code>key</code></br>
<em>
string
</em>
</td>
<td>
<p>Key is the key of the custom policy</p>
</td>
</tr>
<tr>
<td>
<code>value</code></br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Value is the value of the custom policy</p>
</td>
</tr>
<tr>
<td>
<code>organization</code></br>
<em>
string
</em>
</td>
<td>
<p>Organization is the organization of the policy</p>
</td>
</tr>
</tbody>
</table>
<h3 id="dp.wso2.com/v1alpha1.DeploymentStatus">DeploymentStatus
<a class="headerlink" href="#dp.wso2.com%2fv1alpha1.DeploymentStatus" title="Permanent link">¶</a>
</h3>
<p>
(<em>Appears on:</em>
<a href="#dp.wso2.com/v1alpha1.APIStatus">APIStatus</a>)
</p>
<p>
<p>DeploymentStatus contains the status of the API deployment</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>status</code></br>
<em>
string
</em>
</td>
<td>
<p>Status denotes the state of the API in its lifecycle.
Possible values could be Accepted, Invalid, Deploy etc.</p>
</td>
</tr>
<tr>
<td>
<code>message</code></br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Message represents a user friendly message that explains the
current state of the API.</p>
</td>
</tr>
<tr>
<td>
<code>accepted</code></br>
<em>
bool
</em>
</td>
<td>
<p>Accepted represents whether the API is accepted or not.</p>
</td>
</tr>
<tr>
<td>
<code>transitionTime</code></br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.23/#time-v1-meta">
Kubernetes meta/v1.Time
</a>
</em>
</td>
<td>
<p>TransitionTime represents the last known transition timestamp.</p>
</td>
</tr>
<tr>
<td>
<code>events</code></br>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Events contains a list of events related to the API.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="dp.wso2.com/v1alpha1.EnvConfig">EnvConfig
<a class="headerlink" href="#dp.wso2.com%2fv1alpha1.EnvConfig" title="Permanent link">¶</a>
</h3>
<p>
(<em>Appears on:</em>
<a href="#dp.wso2.com/v1alpha1.APISpec">APISpec</a>)
</p>
<p>
<p>EnvConfig contains the environment specific configuration</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>httpRouteRefs</code></br>
<em>
[]string
</em>
</td>
<td>
<p>HTTPRouteRefs denotes the environment of the API.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="dp.wso2.com/v1alpha1.HealthCheck">HealthCheck
<a class="headerlink" href="#dp.wso2.com%2fv1alpha1.HealthCheck" title="Permanent link">¶</a>
</h3>
<p>
(<em>Appears on:</em>
<a href="#dp.wso2.com/v1alpha1.BackendSpec">BackendSpec</a>, 
<a href="#dp.wso2.com/v1alpha1.ResolvedBackend">ResolvedBackend</a>)
</p>
<p>
<p>HealthCheck defines the health check configurations</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>timeout</code></br>
<em>
uint32
</em>
</td>
<td>
<em>(Optional)</em>
<p>Timeout is the time to wait for a health check response.
If the timeout is reached the health check attempt will be considered a failure.</p>
</td>
</tr>
<tr>
<td>
<code>interval</code></br>
<em>
uint32
</em>
</td>
<td>
<em>(Optional)</em>
<p>Interval is the time between health check attempts in seconds.</p>
</td>
</tr>
<tr>
<td>
<code>unhealthyThreshold</code></br>
<em>
uint32
</em>
</td>
<td>
<em>(Optional)</em>
<p>UnhealthyThreshold is the number of consecutive health check failures required
before a backend is marked unhealthy.</p>
</td>
</tr>
<tr>
<td>
<code>healthyThreshold</code></br>
<em>
uint32
</em>
</td>
<td>
<em>(Optional)</em>
<p>HealthyThreshold is the number of healthy health checks required before a host is marked healthy.
Note that during startup, only a single successful health check is required to mark a host healthy.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="dp.wso2.com/v1alpha1.InterceptorInclusion">InterceptorInclusion
(<code>string</code> alias)</p><a class="headerlink" href="#dp.wso2.com%2fv1alpha1.InterceptorInclusion" title="Permanent link">¶</a>
</h3>
<p>
(<em>Appears on:</em>
<a href="#dp.wso2.com/v1alpha1.InterceptorServiceSpec">InterceptorServiceSpec</a>)
</p>
<p>
<p>InterceptorInclusion defines the type of data which can be included in the interceptor request/response path</p>
</p>
<table>
<thead>
<tr>
<th>Value</th>
<th>Description</th>
</tr>
</thead>
<tbody><tr><td><p>&#34;invocation_context&#34;</p></td>
<td><p>InterceptorInclusionInvocationContext is the type to include invocation context</p>
</td>
</tr><tr><td><p>&#34;request_body&#34;</p></td>
<td><p>InterceptorInclusionRequestBody is the type to include request body</p>
</td>
</tr><tr><td><p>&#34;request_headers&#34;</p></td>
<td><p>InterceptorInclusionRequestHeaders is the type to include request headers</p>
</td>
</tr><tr><td><p>&#34;request_trailers&#34;</p></td>
<td><p>InterceptorInclusionRequestTrailers is the type to include request trailers</p>
</td>
</tr><tr><td><p>&#34;response_body&#34;</p></td>
<td><p>InterceptorInclusionResponseBody is the type to include response body</p>
</td>
</tr><tr><td><p>&#34;response_headers&#34;</p></td>
<td><p>InterceptorInclusionResponseHeaders is the type to include response headers</p>
</td>
</tr><tr><td><p>&#34;response_trailers&#34;</p></td>
<td><p>InterceptorInclusionResponseTrailers is the type to include response trailers</p>
</td>
</tr></tbody>
</table>
<h3 id="dp.wso2.com/v1alpha1.InterceptorReference">InterceptorReference
<a class="headerlink" href="#dp.wso2.com%2fv1alpha1.InterceptorReference" title="Permanent link">¶</a>
</h3>
<p>
(<em>Appears on:</em>
<a href="#dp.wso2.com/v1alpha1.PolicySpec">PolicySpec</a>)
</p>
<p>
<p>InterceptorReference holds InterceptorService reference using name and namespace</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code></br>
<em>
string
</em>
</td>
<td>
<p>Name is the referced CR&rsquo;s name of InterceptorService resource.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="dp.wso2.com/v1alpha1.InterceptorServiceSpec">InterceptorServiceSpec
<a class="headerlink" href="#dp.wso2.com%2fv1alpha1.InterceptorServiceSpec" title="Permanent link">¶</a>
</h3>
<p>
(<em>Appears on:</em>
<a href="#dp.wso2.com/v1alpha1.InterceptorService">InterceptorService</a>)
</p>
<p>
<p>InterceptorServiceSpec defines the desired state of InterceptorService</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>backendRef</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.BackendReference">
BackendReference
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>includes</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.InterceptorInclusion">
[]InterceptorInclusion
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Includes defines the types of data which should be included when calling the interceptor service</p>
</td>
</tr>
</tbody>
</table>
<h3 id="dp.wso2.com/v1alpha1.InterceptorServiceStatus">InterceptorServiceStatus
<a class="headerlink" href="#dp.wso2.com%2fv1alpha1.InterceptorServiceStatus" title="Permanent link">¶</a>
</h3>
<p>
(<em>Appears on:</em>
<a href="#dp.wso2.com/v1alpha1.InterceptorService">InterceptorService</a>)
</p>
<p>
<p>InterceptorServiceStatus defines the observed state of InterceptorService</p>
</p>
<h3 id="dp.wso2.com/v1alpha1.JWKS">JWKS
<a class="headerlink" href="#dp.wso2.com%2fv1alpha1.JWKS" title="Permanent link">¶</a>
</h3>
<p>
(<em>Appears on:</em>
<a href="#dp.wso2.com/v1alpha1.SignatureValidation">SignatureValidation</a>)
</p>
<p>
<p>JWKS defines the JWKS endpoint</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>url</code></br>
<em>
string
</em>
</td>
<td>
<p>URL is the URL of the JWKS endpoint</p>
</td>
</tr>
<tr>
<td>
<code>tls</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.CERTConfig">
CERTConfig
</a>
</em>
</td>
<td>
<p>TLS denotes the TLS configuration of the JWKS endpoint</p>
</td>
</tr>
</tbody>
</table>
<h3 id="dp.wso2.com/v1alpha1.JWTAuth">JWTAuth
<a class="headerlink" href="#dp.wso2.com%2fv1alpha1.JWTAuth" title="Permanent link">¶</a>
</h3>
<p>
(<em>Appears on:</em>
<a href="#dp.wso2.com/v1alpha1.APIAuth">APIAuth</a>)
</p>
<p>
<p>JWTAuth JWT Authentication scheme details</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>disabled</code></br>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>Disabled is to disable JWT authentication</p>
</td>
</tr>
<tr>
<td>
<code>header</code></br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Header is the header name used to pass the JWT token</p>
</td>
</tr>
<tr>
<td>
<code>sendTokenToUpstream</code></br>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>SendTokenToUpstream is to specify whether the JWT token should be sent to the upstream</p>
</td>
</tr>
</tbody>
</table>
<h3 id="dp.wso2.com/v1alpha1.JWTIssuerMapping">JWTIssuerMapping
(<code>map[k8s.io/apimachinery/pkg/types.NamespacedName]*github.com/wso2/apk/adapter/internal/operator/apis/dp/v1alpha1.ResolvedJWTIssuer</code> alias)</p><a class="headerlink" href="#dp.wso2.com%2fv1alpha1.JWTIssuerMapping" title="Permanent link">¶</a>
</h3>
<p>
<p>JWTIssuerMapping maps read reconciled Backend and resolve properties into ResolvedJWTIssuer struct</p>
</p>
<h3 id="dp.wso2.com/v1alpha1.JWTIssuerSpec">JWTIssuerSpec
<a class="headerlink" href="#dp.wso2.com%2fv1alpha1.JWTIssuerSpec" title="Permanent link">¶</a>
</h3>
<p>
(<em>Appears on:</em>
<a href="#dp.wso2.com/v1alpha1.JWTIssuer">JWTIssuer</a>)
</p>
<p>
<p>JWTIssuerSpec defines the desired state of JWTIssuer</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code></br>
<em>
string
</em>
</td>
<td>
<p>Name is the unique name of the JWT Issuer in
the Organization defined . &ldquo;Organization/Name&rdquo; can
be used to uniquely identify an Issuer.</p>
</td>
</tr>
<tr>
<td>
<code>organization</code></br>
<em>
string
</em>
</td>
<td>
<p>Organization denotes the organization of the JWT Issuer.</p>
</td>
</tr>
<tr>
<td>
<code>issuer</code></br>
<em>
string
</em>
</td>
<td>
<p>Issuer denotes the issuer of the JWT Issuer.</p>
</td>
</tr>
<tr>
<td>
<code>consumerKeyClaim</code></br>
<em>
string
</em>
</td>
<td>
<p>ConsumerKeyClaim denotes the claim key of the consumer key.</p>
</td>
</tr>
<tr>
<td>
<code>scopesClaim</code></br>
<em>
string
</em>
</td>
<td>
<p>ScopesClaim denotes the claim key of the scopes.</p>
</td>
</tr>
<tr>
<td>
<code>signatureValidation</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.SignatureValidation">
SignatureValidation
</a>
</em>
</td>
<td>
<p>SignatureValidation denotes the signature validation method of jwt</p>
</td>
</tr>
<tr>
<td>
<code>claimMappings</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.[]github.com/wso2/apk/adapter/internal/operator/apis/dp/v1alpha1.ClaimMapping">
[]github.com/wso2/apk/adapter/internal/operator/apis/dp/v1alpha1.ClaimMapping
</a>
</em>
</td>
<td>
<p>ClaimMappings denotes the claim mappings of the jwt</p>
</td>
</tr>
<tr>
<td>
<code>targetRef</code></br>
<em>
sigs.k8s.io/gateway-api/apis/v1alpha2.PolicyTargetReference
</em>
</td>
<td>
<p>TargetRef denotes the reference to the which gateway it applies to</p>
</td>
</tr>
</tbody>
</table>
<h3 id="dp.wso2.com/v1alpha1.JWTIssuerStatus">JWTIssuerStatus
<a class="headerlink" href="#dp.wso2.com%2fv1alpha1.JWTIssuerStatus" title="Permanent link">¶</a>
</h3>
<p>
(<em>Appears on:</em>
<a href="#dp.wso2.com/v1alpha1.JWTIssuer">JWTIssuer</a>)
</p>
<p>
<p>JWTIssuerStatus defines the observed state of JWTIssuer</p>
</p>
<h3 id="dp.wso2.com/v1alpha1.PolicySpec">PolicySpec
<a class="headerlink" href="#dp.wso2.com%2fv1alpha1.PolicySpec" title="Permanent link">¶</a>
</h3>
<p>
(<em>Appears on:</em>
<a href="#dp.wso2.com/v1alpha1.APIPolicySpec">APIPolicySpec</a>)
</p>
<p>
<p>PolicySpec contains API policies</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>requestInterceptors</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.InterceptorReference">
[]InterceptorReference
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>RequestInterceptors referenced to intercetor services to be applied
to the request flow.</p>
</td>
</tr>
<tr>
<td>
<code>responseInterceptors</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.InterceptorReference">
[]InterceptorReference
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>ResponseInterceptors referenced to intercetor services to be applied
to the response flow.</p>
</td>
</tr>
<tr>
<td>
<code>backendJwtPolicy</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.BackendJWTToken">
BackendJWTToken
</a>
</em>
</td>
<td>
<p>BackendJWTPolicy holds reference to backendJWT policy configurations</p>
</td>
</tr>
<tr>
<td>
<code>cORSPolicy</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.CORSPolicy">
CORSPolicy
</a>
</em>
</td>
<td>
<p>CORS policy to be applied to the API.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="dp.wso2.com/v1alpha1.Property">Property
<a class="headerlink" href="#dp.wso2.com%2fv1alpha1.Property" title="Permanent link">¶</a>
</h3>
<p>
(<em>Appears on:</em>
<a href="#dp.wso2.com/v1alpha1.APISpec">APISpec</a>)
</p>
<p>
<p>Property holds key value pair of APIProperties</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code></br>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>value</code></br>
<em>
string
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="dp.wso2.com/v1alpha1.RateLimitAPIPolicy">RateLimitAPIPolicy
<a class="headerlink" href="#dp.wso2.com%2fv1alpha1.RateLimitAPIPolicy" title="Permanent link">¶</a>
</h3>
<p>
(<em>Appears on:</em>
<a href="#dp.wso2.com/v1alpha1.RateLimitPolicySpec">RateLimitPolicySpec</a>)
</p>
<p>
<p>RateLimitAPIPolicy defines the desired state of Policy</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>api</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.APIRateLimitPolicy">
APIRateLimitPolicy
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>API level ratelimit policy</p>
</td>
</tr>
<tr>
<td>
<code>custom</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.CustomRateLimitPolicy">
CustomRateLimitPolicy
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Custom ratelimit policy</p>
</td>
</tr>
</tbody>
</table>
<h3 id="dp.wso2.com/v1alpha1.RateLimitPolicySpec">RateLimitPolicySpec
<a class="headerlink" href="#dp.wso2.com%2fv1alpha1.RateLimitPolicySpec" title="Permanent link">¶</a>
</h3>
<p>
(<em>Appears on:</em>
<a href="#dp.wso2.com/v1alpha1.RateLimitPolicy">RateLimitPolicy</a>)
</p>
<p>
<p>RateLimitPolicySpec defines the desired state of RateLimitPolicy</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>default</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.RateLimitAPIPolicy">
RateLimitAPIPolicy
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>override</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.RateLimitAPIPolicy">
RateLimitAPIPolicy
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>targetRef</code></br>
<em>
sigs.k8s.io/gateway-api/apis/v1alpha2.PolicyTargetReference
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="dp.wso2.com/v1alpha1.RateLimitPolicyStatus">RateLimitPolicyStatus
<a class="headerlink" href="#dp.wso2.com%2fv1alpha1.RateLimitPolicyStatus" title="Permanent link">¶</a>
</h3>
<p>
(<em>Appears on:</em>
<a href="#dp.wso2.com/v1alpha1.RateLimitPolicy">RateLimitPolicy</a>)
</p>
<p>
<p>RateLimitPolicyStatus defines the observed state of RateLimitPolicy</p>
</p>
<h3 id="dp.wso2.com/v1alpha1.RefConfig">RefConfig
<a class="headerlink" href="#dp.wso2.com%2fv1alpha1.RefConfig" title="Permanent link">¶</a>
</h3>
<p>
(<em>Appears on:</em>
<a href="#dp.wso2.com/v1alpha1.CERTConfig">CERTConfig</a>, 
<a href="#dp.wso2.com/v1alpha1.TLSConfig">TLSConfig</a>)
</p>
<p>
<p>RefConfig holds a config for a secret or a configmap</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code></br>
<em>
string
</em>
</td>
<td>
<p>Name of the secret or configmap</p>
</td>
</tr>
<tr>
<td>
<code>key</code></br>
<em>
string
</em>
</td>
<td>
<p>Key of the secret or configmap</p>
</td>
</tr>
</tbody>
</table>
<h3 id="dp.wso2.com/v1alpha1.ResolvedBackend">ResolvedBackend
<a class="headerlink" href="#dp.wso2.com%2fv1alpha1.ResolvedBackend" title="Permanent link">¶</a>
</h3>
<p>
<p>ResolvedBackend holds backend properties</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>Backend</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.Backend">
Backend
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>Services</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.Service">
[]Service
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>Protocol</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.BackendProtocolType">
BackendProtocolType
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>TLS</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.ResolvedTLSConfig">
ResolvedTLSConfig
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>Security</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.ResolvedSecurityConfig">
ResolvedSecurityConfig
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>CircuitBreaker</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.CircuitBreaker">
CircuitBreaker
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>Timeout</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.Timeout">
Timeout
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>Retry</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.RetryConfig">
RetryConfig
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>basePath</code></br>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>HealthCheck</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.HealthCheck">
HealthCheck
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="dp.wso2.com/v1alpha1.ResolvedBasicSecurityConfig">ResolvedBasicSecurityConfig
<a class="headerlink" href="#dp.wso2.com%2fv1alpha1.ResolvedBasicSecurityConfig" title="Permanent link">¶</a>
</h3>
<p>
(<em>Appears on:</em>
<a href="#dp.wso2.com/v1alpha1.ResolvedSecurityConfig">ResolvedSecurityConfig</a>)
</p>
<p>
<p>ResolvedBasicSecurityConfig defines resolved basic security configuration</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>Username</code></br>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>Password</code></br>
<em>
string
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="dp.wso2.com/v1alpha1.ResolvedJWKS">ResolvedJWKS
<a class="headerlink" href="#dp.wso2.com%2fv1alpha1.ResolvedJWKS" title="Permanent link">¶</a>
</h3>
<p>
(<em>Appears on:</em>
<a href="#dp.wso2.com/v1alpha1.ResolvedSignatureValidation">ResolvedSignatureValidation</a>)
</p>
<p>
<p>ResolvedJWKS holds the resolved properties of JWKS</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>URL</code></br>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>TLS</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.ResolvedTLSConfig">
ResolvedTLSConfig
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="dp.wso2.com/v1alpha1.ResolvedJWTIssuer">ResolvedJWTIssuer
<a class="headerlink" href="#dp.wso2.com%2fv1alpha1.ResolvedJWTIssuer" title="Permanent link">¶</a>
</h3>
<p>
<p>ResolvedJWTIssuer holds the resolved properties of JWTIssuer</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>Name</code></br>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>Organization</code></br>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>Issuer</code></br>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>ConsumerKeyClaim</code></br>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>ScopesClaim</code></br>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>SignatureValidation</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.ResolvedSignatureValidation">
ResolvedSignatureValidation
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>ClaimMappings</code></br>
<em>
map[string]string
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="dp.wso2.com/v1alpha1.ResolvedSecurityConfig">ResolvedSecurityConfig
<a class="headerlink" href="#dp.wso2.com%2fv1alpha1.ResolvedSecurityConfig" title="Permanent link">¶</a>
</h3>
<p>
(<em>Appears on:</em>
<a href="#dp.wso2.com/v1alpha1.ResolvedBackend">ResolvedBackend</a>)
</p>
<p>
<p>ResolvedSecurityConfig defines enpoint resolved security configurations</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>Type</code></br>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>Basic</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.ResolvedBasicSecurityConfig">
ResolvedBasicSecurityConfig
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="dp.wso2.com/v1alpha1.ResolvedSignatureValidation">ResolvedSignatureValidation
<a class="headerlink" href="#dp.wso2.com%2fv1alpha1.ResolvedSignatureValidation" title="Permanent link">¶</a>
</h3>
<p>
(<em>Appears on:</em>
<a href="#dp.wso2.com/v1alpha1.ResolvedJWTIssuer">ResolvedJWTIssuer</a>)
</p>
<p>
<p>ResolvedSignatureValidation holds the resolved properties of SignatureValidation</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>JWKS</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.ResolvedJWKS">
ResolvedJWKS
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>Certificate</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.ResolvedTLSConfig">
ResolvedTLSConfig
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="dp.wso2.com/v1alpha1.ResolvedTLSConfig">ResolvedTLSConfig
<a class="headerlink" href="#dp.wso2.com%2fv1alpha1.ResolvedTLSConfig" title="Permanent link">¶</a>
</h3>
<p>
(<em>Appears on:</em>
<a href="#dp.wso2.com/v1alpha1.ResolvedBackend">ResolvedBackend</a>, 
<a href="#dp.wso2.com/v1alpha1.ResolvedJWKS">ResolvedJWKS</a>, 
<a href="#dp.wso2.com/v1alpha1.ResolvedSignatureValidation">ResolvedSignatureValidation</a>)
</p>
<p>
<p>ResolvedTLSConfig defines enpoint TLS configurations</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>ResolvedCertificate</code></br>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>AllowedSANs</code></br>
<em>
[]string
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="dp.wso2.com/v1alpha1.RetryConfig">RetryConfig
<a class="headerlink" href="#dp.wso2.com%2fv1alpha1.RetryConfig" title="Permanent link">¶</a>
</h3>
<p>
(<em>Appears on:</em>
<a href="#dp.wso2.com/v1alpha1.BackendSpec">BackendSpec</a>, 
<a href="#dp.wso2.com/v1alpha1.ResolvedBackend">ResolvedBackend</a>)
</p>
<p>
<p>RetryConfig defines retry configurations</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>count</code></br>
<em>
uint32
</em>
</td>
<td>
<p>Count defines the number of retries.
If exceeded, TooEarly(425 response code) response will be sent to the client.</p>
</td>
</tr>
<tr>
<td>
<code>baseIntervalMillis</code></br>
<em>
uint32
</em>
</td>
<td>
<em>(Optional)</em>
<p>BaseIntervalMillis is exponential retry back off and it defines the base interval between retries in milliseconds.
maximum interval is 10 times of the BaseIntervalMillis</p>
</td>
</tr>
<tr>
<td>
<code>statusCodes</code></br>
<em>
[]uint32
</em>
</td>
<td>
<em>(Optional)</em>
<p>StatusCodes defines the list of status codes to retry</p>
</td>
</tr>
</tbody>
</table>
<h3 id="dp.wso2.com/v1alpha1.ScopeSpec">ScopeSpec
<a class="headerlink" href="#dp.wso2.com%2fv1alpha1.ScopeSpec" title="Permanent link">¶</a>
</h3>
<p>
(<em>Appears on:</em>
<a href="#dp.wso2.com/v1alpha1.Scope">Scope</a>)
</p>
<p>
<p>ScopeSpec defines the desired state of Scope</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>names</code></br>
<em>
[]string
</em>
</td>
<td>
<p>Name scope name</p>
</td>
</tr>
</tbody>
</table>
<h3 id="dp.wso2.com/v1alpha1.ScopeStatus">ScopeStatus
<a class="headerlink" href="#dp.wso2.com%2fv1alpha1.ScopeStatus" title="Permanent link">¶</a>
</h3>
<p>
(<em>Appears on:</em>
<a href="#dp.wso2.com/v1alpha1.Scope">Scope</a>)
</p>
<p>
<p>ScopeStatus defines the observed state of Scope</p>
</p>
<h3 id="dp.wso2.com/v1alpha1.SecretRef">SecretRef
<a class="headerlink" href="#dp.wso2.com%2fv1alpha1.SecretRef" title="Permanent link">¶</a>
</h3>
<p>
(<em>Appears on:</em>
<a href="#dp.wso2.com/v1alpha1.BasicSecurityConfig">BasicSecurityConfig</a>)
</p>
<p>
<p>SecretRef to credentials</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code></br>
<em>
string
</em>
</td>
<td>
<p>Name of the secret</p>
</td>
</tr>
<tr>
<td>
<code>usernameKey</code></br>
<em>
string
</em>
</td>
<td>
<p>Namespace of the secret</p>
</td>
</tr>
<tr>
<td>
<code>passwordKey</code></br>
<em>
string
</em>
</td>
<td>
<p>Key of the secret</p>
</td>
</tr>
</tbody>
</table>
<h3 id="dp.wso2.com/v1alpha1.SecurityConfig">SecurityConfig
<a class="headerlink" href="#dp.wso2.com%2fv1alpha1.SecurityConfig" title="Permanent link">¶</a>
</h3>
<p>
(<em>Appears on:</em>
<a href="#dp.wso2.com/v1alpha1.BackendSpec">BackendSpec</a>)
</p>
<p>
<p>SecurityConfig defines enpoint security configurations</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>basic</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.BasicSecurityConfig">
BasicSecurityConfig
</a>
</em>
</td>
<td>
<p>Basic security configuration</p>
</td>
</tr>
</tbody>
</table>
<h3 id="dp.wso2.com/v1alpha1.Service">Service
<a class="headerlink" href="#dp.wso2.com%2fv1alpha1.Service" title="Permanent link">¶</a>
</h3>
<p>
(<em>Appears on:</em>
<a href="#dp.wso2.com/v1alpha1.BackendSpec">BackendSpec</a>, 
<a href="#dp.wso2.com/v1alpha1.ResolvedBackend">ResolvedBackend</a>)
</p>
<p>
<p>Service holds host and port information for the service</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>host</code></br>
<em>
string
</em>
</td>
<td>
<p>Host is the hostname of the service</p>
</td>
</tr>
<tr>
<td>
<code>port</code></br>
<em>
uint32
</em>
</td>
<td>
<p>Port of the service</p>
</td>
</tr>
</tbody>
</table>
<h3 id="dp.wso2.com/v1alpha1.SignatureValidation">SignatureValidation
<a class="headerlink" href="#dp.wso2.com%2fv1alpha1.SignatureValidation" title="Permanent link">¶</a>
</h3>
<p>
(<em>Appears on:</em>
<a href="#dp.wso2.com/v1alpha1.JWTIssuerSpec">JWTIssuerSpec</a>)
</p>
<p>
<p>SignatureValidation defines the signature validation method</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>jwks</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.JWKS">
JWKS
</a>
</em>
</td>
<td>
<p>JWKS denotes the JWKS endpoint information</p>
</td>
</tr>
<tr>
<td>
<code>certificate</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.CERTConfig">
CERTConfig
</a>
</em>
</td>
<td>
<p>Certificate denotes the certificate information</p>
</td>
</tr>
</tbody>
</table>
<h3 id="dp.wso2.com/v1alpha1.TLSConfig">TLSConfig
<a class="headerlink" href="#dp.wso2.com%2fv1alpha1.TLSConfig" title="Permanent link">¶</a>
</h3>
<p>
(<em>Appears on:</em>
<a href="#dp.wso2.com/v1alpha1.BackendSpec">BackendSpec</a>)
</p>
<p>
<p>TLSConfig defines enpoint TLS configurations</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>certificateInline</code></br>
<em>
string
</em>
</td>
<td>
<p>CertificateInline is the Inline Certificate entry</p>
</td>
</tr>
<tr>
<td>
<code>secretRef</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.RefConfig">
RefConfig
</a>
</em>
</td>
<td>
<p>SecretRef denotes the reference to the Secret that contains the Certificate</p>
</td>
</tr>
<tr>
<td>
<code>configMapRef</code></br>
<em>
<a href="#dp.wso2.com/v1alpha1.RefConfig">
RefConfig
</a>
</em>
</td>
<td>
<p>ConfigMapRef denotes the reference to the ConfigMap that contains the Certificate</p>
</td>
</tr>
<tr>
<td>
<code>allowedSANs</code></br>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>AllowedCNs is the list of allowed Subject Alternative Names (SANs)</p>
</td>
</tr>
</tbody>
</table>
<h3 id="dp.wso2.com/v1alpha1.TestConsoleKeyAuth">TestConsoleKeyAuth
<a class="headerlink" href="#dp.wso2.com%2fv1alpha1.TestConsoleKeyAuth" title="Permanent link">¶</a>
</h3>
<p>
(<em>Appears on:</em>
<a href="#dp.wso2.com/v1alpha1.APIAuth">APIAuth</a>)
</p>
<p>
<p>TestConsoleKeyAuth Test Console Key Authentication scheme details</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>header</code></br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Header is the header name used to pass the Test Console Key</p>
</td>
</tr>
<tr>
<td>
<code>sendTokenToUpstream</code></br>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>SendTokenToUpstream is to specify whether the Test Console Key should be sent to the upstream</p>
</td>
</tr>
</tbody>
</table>
<h3 id="dp.wso2.com/v1alpha1.Timeout">Timeout
<a class="headerlink" href="#dp.wso2.com%2fv1alpha1.Timeout" title="Permanent link">¶</a>
</h3>
<p>
(<em>Appears on:</em>
<a href="#dp.wso2.com/v1alpha1.BackendSpec">BackendSpec</a>, 
<a href="#dp.wso2.com/v1alpha1.ResolvedBackend">ResolvedBackend</a>)
</p>
<p>
<p>Timeout defines the timeout configurations</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>upstreamResponseTimeout</code></br>
<em>
uint32
</em>
</td>
<td>
<p>UpstreamResponseTimeout spans between the point at which the entire downstream request (i.e. end-of-stream) has been processed and
when the upstream response has been completely processed.
A value of 0 will disable the route’s timeout.</p>
</td>
</tr>
<tr>
<td>
<code>downstreamRequestIdleTimeout</code></br>
<em>
uint32
</em>
</td>
<td>
<em>(Optional)</em>
<p>DownstreamRequestIdleTimeout bounds the amount of time the request&rsquo;s stream may be idle.
A value of 0 will completely disable the route&rsquo;s idle timeout.</p>
</td>
</tr>
</tbody>
</table>
<hr/>
<p><em>
Generated with <code>gen-crd-api-reference-docs</code>.
</em></p>
