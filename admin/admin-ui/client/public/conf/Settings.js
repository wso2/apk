

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
        client_id: '01edf918-4b86-1a56-93bc-ddc78be49440',
        client_secret: '01edf918-4b86-1a56-9d87-fb96b32c8d5e',
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
        rest_api: 'https://127.0.0.1:9095/api/am/admin',
    }
};

if (typeof module !== 'undefined') {
    module.exports = Settings; // For Jest unit tests
}
