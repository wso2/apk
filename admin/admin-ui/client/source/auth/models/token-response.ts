/**
 * Interface of the OAuth2/OIDC tokens.
 */
 export interface TokenResponseInterface {
    accessToken: string;
    idToken: string;
    refreshToken: string;
    expiresIn: string;
    scope: string;
    tokenType: string;
}