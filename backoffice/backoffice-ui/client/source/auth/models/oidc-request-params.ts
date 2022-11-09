/**
 * OIDC request parameters.
 */
 export interface OIDCRequestParamsInterface {
    clientId: string;
    // redirectUri: string;
    scope: string;
    state: string;
    serverOrigin: string;
}