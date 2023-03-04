

const Settings = {
    // oath0 config
    // idp: {
    //     client_id: 'wqCauPZqUn6UigcAZbU9Z6jxwNTfzOAb',
    //     well_known: 'https://construct.auth0.com/.well-known/openid-configuration',
    //     serverOrigin: 'https://construct.auth0.com/',
    //     redirect_uri: 'https://localhost:4000',
    //     logout_endpoint: 'https://construct.auth0.com/v2/logout',
    //     scope: 'openid offline_access',
    //     state: 'RlZyVjlqYUpHTzltWC42c2FNRDRJT1JPfk1+TUFEa0RLb04yZldwYkpxVA==',
    //     pkce: true,
    // },
    idp: {
        client_id: '01edb99d-d4d3-1ae8-b8b9-f09c5a557009',
        client_secret: '01edb99d-d4d3-1ae8-97f9-55a7256d84e4',
        host: 'idp.am.wso2.com:9095',
        server_origin: 'https://idp.am.wso2.com:9095/',
        redirect_uri: 'https://localhost:4000',
        logout_endpoint: 'https://idp.am.wso2.com:9095/logout',
        scope: 'openid offline_access',
        state: 'RlZyVjlqYUpHTzltWC42c2FNRDRJT1JPfk1+TUFEa0RLb04yZldwYkpxVA==',
        authorization_endpoint: 'https://idp.am.wso2.com:9095/oauth2/authorize',
        token_endpoint: 'https://idp.am.wso2.com:9095/oauth2/token',
        jwks_uri: 'https://idp.am.wso2.com:9095/oauth2/jwks',
        issuer: 'https://idp.am.wso2.com:9095/oauth2/token',
        userinfo_endpoint: 'https://idp.am.wso2.com:9095/oauth2/userinfo',
        pkce: false,
    },
    app: {
        rest_api: 'https://api.am.wso2.com:9095/api/am/admin',
    }
};

if (typeof module !== 'undefined') {
    module.exports = Settings; // For Jest unit tests
}
