

const Settings = {
    API_PORT: 9443,
    API_HOST: 'localhost',
    API_TRANSPORT: 'http',
    IDP_CLIENT_ID: 'FbCSH23HybQMV9UlXJfeKHogAEHojzCO',
    wellKnown: 'https://dev-kw-oeodk.us.auth0.com/.well-known/openid-configuration',
    serverOrigin: 'https://dev-kw-oeodk.us.auth0.com/',
    loginUri: 'https://localhost:4000',
    logoutEndpoint: 'https://dev-kw-oeodk.us.auth0.com/v2/logout',
    scope: 'openid offline_access',
    state: 'RlZyVjlqYUpHTzltWC42c2FNRDRJT1JPfk1+TUFEa0RLb04yZldwYkpxVA==',
    restApi: 'http://localhost:81/chanaka3d/wso-2_api_manager_admin/v3/',
};

if (typeof module !== 'undefined') {
    module.exports = Settings; // For Jest unit tests
}
