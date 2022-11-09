import { getCodeVerifier, getCodeChallenge, getJWKForTheIdToken, isValidIdToken } from './crypto';
import { OIDCRequestParamsInterface } from './models/oidc-request-params';
import { TokenResponseInterface } from './models/token-response';
import {
    AUTHORIZATION_CODE,
    PKCE_CODE_VERIFIER,
    SERVICE_RESOURCES,
    REQUEST_PARAMS,
    REQUEST_STATUS,
} from './constants/token';
import { getSessionParameter, removeSessionParameter, setSessionParameter } from "./session";
import { getAuthorizeEndpoint, getTokenEndpoint, getJwksUri, getIssuer, getToken } from "./op-config";
import Settings from '../../public/conf/Settings'
const axios = require('axios');

/**
 * Checks whether authorization code present in the request.
 *
 * @returns {boolean} true if authorization code is present.
 */
export const hasAuthorizationCode = (): boolean => {
    return !!new URL(window.location.href).searchParams.get(AUTHORIZATION_CODE);
};

export const hasValidToken = (): boolean => {
    return !!getToken();
};

export const sendAuthorizationRequest = (requestParams: OIDCRequestParamsInterface): Promise<never> | any => {
    const authorizeEndpoint = getAuthorizeEndpoint();

    if (!authorizeEndpoint || authorizeEndpoint.trim().length === 0) {
        return Promise.reject(new Error("Invalid authorize endpoint found."));
    }

    let authorizeRequest = authorizeEndpoint + "?response_type=code&client_id="
        + requestParams.clientId;

    authorizeRequest += "&scope=" + requestParams.scope;

    authorizeRequest += "&state=" + requestParams.state;

    const codeVerifier = getCodeVerifier();
    const codeChallenge = getCodeChallenge(codeVerifier);
    setSessionParameter(PKCE_CODE_VERIFIER, codeVerifier);

    authorizeRequest += "&code_challenge_method=S256&code_challenge=" + codeChallenge;
    authorizeRequest += `&redirect_uri=${Settings.loginUri}/token`;

    document.location.href = authorizeRequest;

    return false;
};

/**
 * Validate id_token.
 *
 * @param {string} clientId client ID.
 * @param {string} idToken id_token received from the IdP.
 * @returns {Promise<boolean>} whether token is valid.
 */
/* eslint-disable @typescript-eslint/no-explicit-any */
const validateIdToken = (clientId: string, idToken: string, serverOrigin: string): Promise<any> => {
    const jwksEndpoint = getJwksUri();

    if (!jwksEndpoint || jwksEndpoint.trim().length === 0) {
        return Promise.reject("Invalid JWKS URI found.");
    }

    return axios.get(jwksEndpoint)
        .then((response: any) => {
            if (response.status !== 200) {
                return Promise.reject(new Error("Failed to load public keys from JWKS URI: "
                    + jwksEndpoint));
            }

            const jwk = getJWKForTheIdToken(idToken.split(".")[0], response.data.keys);
            let issuer = getIssuer();
            if (!issuer || issuer.trim().length === 0) {
                issuer = serverOrigin + SERVICE_RESOURCES.token;
            }

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
    removeSessionParameter(REQUEST_STATUS);

    const body = [
        `client_id=${requestParams.clientId}`,
        `code=${code}`,
        "grant_type=authorization_code",
        `redirect_uri=${Settings.loginUri}/users`,
        `code_verifier=${getSessionParameter(PKCE_CODE_VERIFIER)}`];

    return axios.post(tokenEndpoint, body.join("&"))
        .then((response: any) => {
            if (response.status !== 200) {
                return Promise.reject(new Error("Invalid status code received in the token response: "
                    + response.status));
            }
            removeSessionParameter(PKCE_CODE_VERIFIER);

            return validateIdToken(requestParams.clientId, response.data.id_token, requestParams.serverOrigin)
                .then((valid) => {
                    if (valid) {
                        setSessionParameter(REQUEST_PARAMS, JSON.stringify(requestParams));
                        const tokenResponse: TokenResponseInterface = {
                            accessToken: response.data.access_token,
                            expiresIn: response.data.expires_in,
                            idToken: response.data.id_token,
                            refreshToken: response.data.refresh_token,
                            scope: response.data.scope,
                            tokenType: response.data.token_type
                        };
                        return Promise.resolve(tokenResponse)
                    }
                    return Promise.reject(new Error("Invalid id_token in the token response: " + response.data.id_token));
                });

            // const body = [];
            // body.push(`client_id=${requestParams.clientId}`);
            // body.push(`scope=apim:api_manage apim:subscription_manage apim:tier_manage apim:admin`);
            // body.push("grant_type=urn:ietf:params:oauth:grant-type:token-exchange");
            // body.push(`subject_token_type=urn:ietf:params:oauth:token-type:jwt`);
            // body.push(`requested_token_type=urn:ietf:params:oauth:token-type:jwt`);
            // body.push(`subject_token=${response.data.access_token}`);
            // body.push(`org_handle=organization`);
            // return axios.post(stsEndoint, body.join("&"))
            //     .then((response: any) => {
            //         window.sessionStorage.setItem("exchanged_token", response.data.access_token);
            //         window.location.href = "https://localhost:4000/users";
            //     });
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
