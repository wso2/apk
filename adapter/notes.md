func getHTTPFilters(globalLuaScript string) []*hcmv3.HttpFilter {

 r.ProviderResources.GatewayAPIResources.Subscribe(ctx),
 
  buildXdsTCPListener(name, address string, port uint32, k
  
func (t *Translator) addXdsHTTPFilterChain(xdsListener *listenerv3.Listener,


func (t *Translator) buildExtAuth(policy *egv1a1.SecurityPolicy, resources *Resources, envoyProxy *egv1a1.EnvoyProxy) (*ir.ExtAuth, error) {
	// create a function like this and hard code the enforcer cluster and add this to the ir.httproute
	
irRoute := &ir.HTTPRoute{
			Name: irRouteName(httpRoute, ruleIdx, matchIdx),
		}
		
		
Create a method name processAPI 
in this function pass 
securityPolicies []*API,
	gateways []*GatewayContext,
	routes []RouteContext,
	resources *Resources,
	xdsIR XdsIRMap,
	

loop through all targeted httproute refs in each API. find the correct irRoutes using this approach 

irListener := xdsIR[irKey].GetHTTPListener(irListenerName(listener))
			if irListener != nil {
				for _, r := range irListener.Routes {
					if strings.HasPrefix(r.Name, prefix) {

now add the extAuth to all the matching irRoutes


func enableFilterOnRoute(route *routev3.Route, filterName string) error {

k delete httproute apk-test-wso2-apk-authentication-endpoint-ds-httproute apk-test-wso2-apk-commonoauth-ds-httproute apk-test-wso2-apk-config-deploy-api-route apk-test-wso2-apk-jwks-endpoint-ds-httproute apk-test-wso2-apk-dcr-ds-httproute apk-test-wso2-apk-oauth-ds-httproute -n apk
