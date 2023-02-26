

const Settings = {
    idp: {
        IDP_CLIENT_ID: 'wqCauPZqUn6UigcAZbU9Z6jxwNTfzOAb',
        wellKnown: 'https://construct.auth0.com/.well-known/openid-configuration',
        serverOrigin: 'https://construct.auth0.com/',
        loginUri: 'https://localhost:4000',
        logoutEndpoint: 'https://construct.auth0.com/v2/logout',
        scope: 'openid offline_access',
        state: 'RlZyVjlqYUpHTzltWC42c2FNRDRJT1JPfk1+TUFEa0RLb04yZldwYkpxVA==',
    },
    server: {
        API_PORT: 9443,
        API_HOST: 'localhost',
        API_TRANSPORT: 'http',
        restApi: 'http://localhost:9445/api/am/admin/v3/',
    },
    theme: {
        defaultPath: '/dashboard/default',
        fontFamily: `'Montserrat', sans-serif;`,
        i18n: 'en',
        miniDrawer: false,
        container: true,
        mode: 'light',
        presetColor: 'default',
        themeDirection: 'ltr',
        docUrl: 'https://apim.docs.wso2.com/en/',
        drawerWidth: 260,
        twitterColor: '#1DA1F2',
        facebookColor: '#3b5998',
        linkedInColor: '#0e76a8',
    }
};

if (typeof module !== 'undefined') {
    module.exports = Settings; // For Jest unit tests
}
