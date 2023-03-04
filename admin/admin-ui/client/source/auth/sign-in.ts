import axios from 'axios';
import { getCodeVerifier, getCodeChallenge, getJWKForTheIdToken, isValidIdToken } from './crypto';
import { OIDCRequestParamsInterface } from './types/oidc-request-params';
import { TokenResponseInterface } from './types/token-response';
import {
    AUTHORIZATION_CODE,
    PKCE_CODE_VERIFIER,
    SERVICE_RESOURCES,
    REQUEST_PARAMS,
} from './constants/token';
import { getSessionParameter, removeSessionParameter, setSessionParameter } from "./session";
import { getAuthorizeEndpoint, getTokenEndpoint, getJwksUri, getIssuer, getToken } from "./op-config";
// eslint-disable-next-line @typescript-eslint/no-var-requires
const Settings = require('Settings');
/**
 * Send authorization request.
 * @param requestParams  Request parameters.
 * @returns  Promise.
 */
export const sendAuthorizationRequest = (requestParams: OIDCRequestParamsInterface): Promise<never> | any => {
    const authorizeEndpoint = getAuthorizeEndpoint();
    if (!authorizeEndpoint || authorizeEndpoint.trim().length === 0) {
        return Promise.reject(new Error("Invalid authorize endpoint found."));
    }
    // Generate code verifier and code challenge.

    const codeVerifier = getCodeVerifier();
    const codeChallenge = getCodeChallenge(codeVerifier);
    setSessionParameter(PKCE_CODE_VERIFIER, codeVerifier);
    const authorizeRequest = `${authorizeEndpoint}?` +
        `response_type=code` +
        `&client_id=${requestParams.clientId}` +
        `&scope=${requestParams.scope}` +
        `&state=${requestParams.state}` +
        `&code_challenge_method=S256` +
        `&code_challenge=${codeChallenge}` +
        `&redirect_uri=${Settings.idp.redirect_uri}`;
    
    document.location.href = authorizeRequest;
    return false;
};

/**
 *  
This function is used to validate an ID token obtained from the authorization server after exchanging an authorization code.

The ID token is a JSON Web Token (JWT) that contains information about the authenticated user and the authorization transaction. 
The ID token is signed by the authorization server using a private key and can be verified using the public key obtained from the
 server's JSON Web Key Set (JWKS) endpoint.

The purpose of the validateIdToken function is to retrieve the public key from the JWKS endpoint and use it to verify the 
signature on the ID token. The function also checks that the iss (issuer) claim in the ID token matches the expected issuer 
(i.e., the authorization server), that the aud (audience) claim in the ID token matches the client ID, and that the nonce 
value passed during the authentication request matches the nonce value in the ID token.

By validating the ID token, the function ensures that the token was issued by a trusted authority, 
that it has not been tampered with, and that it contains the expected claims. This helps prevent impersonation and other 
security threats in the SPA.
 *
 * @param {string} clientId client ID.
 * @param {string} idToken id_token received from the IdP.
 * @returns {Promise<boolean>} whether token is valid.
 */
/* eslint-disable @typescript-eslint/no-explicit-any */
const validateIdToken = (clientId: string, idToken: string, serverOrigin: string): Promise<any> => {
    // Get the JWKS endpoint from the OP configuration.
    const jwksEndpoint = getJwksUri();

    // If the JWKS endpoint is not available, return an error.
    if (!jwksEndpoint || jwksEndpoint.trim().length === 0) {
        return Promise.reject("Invalid JWKS URI found.");
    }
    // Get the public key from the JWKS endpoint.
    return axios.get(jwksEndpoint)
        .then((response: any) => {
            if (response.status !== 200) {
                return Promise.reject(new Error("Failed to load public keys from JWKS URI: "
                    + jwksEndpoint));
            }
            // Get the public key from the JWKS endpoint.
            const jwk = getJWKForTheIdToken(idToken.split(".")[0], response.data.keys);
            let issuer = getIssuer();
            // If the issuer is not available, use the server origin.
            if (!issuer || issuer.trim().length === 0) {
                issuer = serverOrigin + SERVICE_RESOURCES.token;
            }
            // Validate the ID token.
            return Promise.resolve(isValidIdToken(idToken, jwk, clientId, issuer));
        }).catch((error: any) => {
            return Promise.reject(error);
        });
};

/**
 * Send token request.
 *
 * @param {OIDCRequestParamsInterface} requestParams request parameters required for token request.
 * @returns {Promise<TokenResponseInterface>} token response data or error.
 */
export const sendTokenRequest = (
    requestParams: OIDCRequestParamsInterface
): Promise<TokenResponseInterface> => {

    const tokenEndpoint = getTokenEndpoint();
    // const stsEndoint = 'https://da59-203-94-95-4.in.ngrok.io/api/am/sts/v1/oauth2/token';

    if (!tokenEndpoint || tokenEndpoint.trim().length === 0) {
        return Promise.reject(new Error("Invalid token endpoint found."));
    }

    const code = new URL(window.location.href).searchParams.get(AUTHORIZATION_CODE);

    const body = [
        `client_id=${requestParams.clientId}`,
        `code=${code}`,
        "grant_type=authorization_code",
        `redirect_uri=${Settings.idp.redirect_uri}`];

    if (Settings.idp.pkce) {
        body.push(`code_verifier=${getSessionParameter(PKCE_CODE_VERIFIER)}`);
    }

    return axios.post(tokenEndpoint, body.join("&"))
        .then((response: any) => {
            if (response.status !== 200) {
                return Promise.reject(new Error("Invalid status code received in the token response: "
                    + response.status));
            }
            const tokenResponse: TokenResponseInterface = {
                accessToken: response.data.access_token,
                expiresIn: response.data.expires_in,
                idToken: response.data.id_token,
                refreshToken: response.data.refresh_token,
                scope: response.data.scope,
                tokenType: response.data.token_type
            };
            if (Settings.idp.pkce) {
                removeSessionParameter(PKCE_CODE_VERIFIER);
                return validateIdToken(requestParams.clientId, response.data.id_token, requestParams.serverOrigin)
                    .then((valid) => {
                        if (valid) {
                            setSessionParameter(REQUEST_PARAMS, JSON.stringify(requestParams));

                            return Promise.resolve(tokenResponse)
                        }
                        return Promise.reject(new Error("Invalid id_token in the token response: " + response.data.id_token));
                    });
            } else {
                // If PKCE is disabled, set the request parameters in the session storage without the verification.
                return Promise.resolve(tokenResponse)
            }

        }).catch((error: any) => {
            return Promise.reject(error);
        });
};

/**
 * Send refresh token request.
 *
 * @param {OIDCRequestParamsInterface} requestParams request parameters required for token request.
 * @param {string} refreshToken
 * @returns {Promise<TokenResponseInterface>} refresh token response data or error.
 */
export const sendRefreshTokenRequest = (
    requestParams: OIDCRequestParamsInterface,
    refreshToken: string
): Promise<TokenResponseInterface> => {

    const tokenEndpoint = getTokenEndpoint();

    if (!tokenEndpoint || tokenEndpoint.trim().length === 0) {
        return Promise.reject("Invalid token endpoint found.");
    }

    const body = [
        `client_id=${requestParams.clientId}`,
        `refresh_token=${refreshToken}`,
        "grant_type=refresh_token"];

    return axios.post(tokenEndpoint, body.join("&"))
        .then((response: any) => {
            if (response.status !== 200) {
                return Promise.reject(new Error("Invalid status code received in the refresh token response: "
                    + response.status));
            }

            return validateIdToken(requestParams.clientId, response.data.id_token, requestParams.serverOrigin)
                .then((valid) => {
                    if (valid) {
                        const tokenResponse: TokenResponseInterface = {
                            accessToken: response.data.access_token,
                            expiresIn: response.data.expires_in,
                            idToken: response.data.id_token,
                            refreshToken: response.data.refresh_token,
                            scope: response.data.scope,
                            tokenType: response.data.token_type
                        };

                        return Promise.resolve(tokenResponse);
                    }
                    return Promise.reject(new Error("Invalid id_token in the token response: " +
                        response.data.id_token));
                });
        }).catch((error: any) => {
            return Promise.reject(error);
        });
};
