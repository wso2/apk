import axios from "axios";
import {
    ACCESS_TOKEN,
    AUTHORIZATION_ENDPOINT,
    TOKEN_ENDPOINT,
    END_SESSION_ENDPOINT,
    JWKS_ENDPOINT,
    ISSUER,
    USERINFO_ENDPOINT,
} from './constants/token';
import { getSessionParameter, removeSessionParameter, setSessionParameter } from "./session";
// eslint-disable-next-line @typescript-eslint/no-var-requires
const Settings = require('Settings');

/**
 * Set OAuth2 authorize endpoint.
 *
 * @param {string} authorizationEndpoint
 */
export const setAuthorizeEndpoint = (authorizationEndpoint: string): void => {
    setSessionParameter(AUTHORIZATION_ENDPOINT, authorizationEndpoint);
};

/**
 * Set OAuth2 token endpoint.
 *
 * @param {string} tokenEndpoint
 */
export const setTokenEndpoint = (tokenEndpoint: string): void => {
    setSessionParameter(TOKEN_ENDPOINT, tokenEndpoint);
};

/**
 * Set OIDC end session endpoint.
 *
 * @param {string} endSessionEndpoint
 */
export const setEndSessionEndpoint = (endSessionEndpoint: string): void => {
    setSessionParameter(END_SESSION_ENDPOINT, endSessionEndpoint);
};

/**
 * Set JWKS URI.
 *
 * @param jwksEndpoint
 */
export const setJwksUri = (jwksEndpoint: string): void => {
    setSessionParameter(JWKS_ENDPOINT, jwksEndpoint);
};

/**
 * Set id_token issuer.
 *
 * @param issuer id_token issuer.
 */
export const setIssuer = (issuer: string): void => {
    setSessionParameter(ISSUER, issuer);
};

/**
 * Set userinfo endpoint.
 * @param userinfoEndpoint 
 */
export const setUserinfoEndpoint = (userinfoEndpoint: string): void => {
    setSessionParameter(USERINFO_ENDPOINT, userinfoEndpoint);
};
/**
 * Initialize openid provider configuration.
 *
 * @param {string} wellKnownEndpoint openid provider configuration.
 * @returns {Promise<any>} promise.
 */
/* eslint-disable @typescript-eslint/no-explicit-any */
export const initOPConfiguration = (
    wellKnownEndpoint: string,
): Promise<any> => {
    if (wellKnownEndpoint && wellKnownEndpoint.trim().length > 0) {
        if (!wellKnownEndpoint || wellKnownEndpoint.trim().length === 0) {
            return Promise.reject(new Error("OpenID provider configuration endpoint is not defined."));
        }
        return axios.get(wellKnownEndpoint)
            .then((response) => {
                if (response.status !== 200) {
                    return Promise.reject(new Error("Failed to load OpenID provider configuration from: "
                        + wellKnownEndpoint));
                }
                setAuthorizeEndpoint(response.data.authorization_endpoint);
                setTokenEndpoint(response.data.token_endpoint);
                setEndSessionEndpoint(Settings.idp.logout_endpoint);
                setJwksUri(response.data.jwks_uri);
                setIssuer(response.data.issuer);
                setUserinfoEndpoint(response.data.userinfo_endpoint);
                return Promise.resolve("success");
            }).catch((error) => {
                return Promise.reject(error);
            });
    } else {
        setAuthorizeEndpoint(Settings.idp.authorization_endpoint);
        setTokenEndpoint(Settings.idp.token_endpoint);
        setEndSessionEndpoint(Settings.idp.logout_endpoint);
        setJwksUri(Settings.idp.jwks_uri);
        setIssuer(Settings.idp.issuer);
        setUserinfoEndpoint(Settings.idp.userinfo_endpoint);
        return Promise.resolve("success");
    }
};

/**
 * Reset openid provider configuration.
 */
export const resetOPConfiguration = (): void => {
    removeSessionParameter(AUTHORIZATION_ENDPOINT);
    removeSessionParameter(TOKEN_ENDPOINT);
    removeSessionParameter(END_SESSION_ENDPOINT);
    removeSessionParameter(JWKS_ENDPOINT);
    removeSessionParameter(ISSUER);
    removeSessionParameter(USERINFO_ENDPOINT);
};

/**
 * Get OAuth2 authorize endpoint.
 *
 * @returns {string|null}
 */
export const getAuthorizeEndpoint = (): string | null => {
    return getSessionParameter(AUTHORIZATION_ENDPOINT);
};

/**
 * Get OAuth2 token endpoint.
 *
 * @returns {string|null}
 */
export const getTokenEndpoint = (): string | null => {
    return getSessionParameter(TOKEN_ENDPOINT);
};

/**
 * Get OIDC end session endpoint.
 *
 * @returns {string|null}
 */
export const getEndSessionEndpoint = (): string | null => {
    return getSessionParameter(END_SESSION_ENDPOINT);
};

/**
 * Get JWKS URI.
 *
 * @returns {string|null}
 */
export const getJwksUri = (): string | null => {
    return getSessionParameter(JWKS_ENDPOINT);
};

/**
 * Get id_token issuer.
 *
 * @returns {any}
 */
export const getIssuer = (): string => {
    return getSessionParameter(ISSUER);
};

export const getToken = (): string => {
    return getSessionParameter(ACCESS_TOKEN);
};

export const getUserInfoEndpoint = (): string => {
    return getSessionParameter(USERINFO_ENDPOINT);
}