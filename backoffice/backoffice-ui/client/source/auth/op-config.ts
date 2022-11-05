import axios from "axios";
import {
    ACCESS_TOKEN,
    AUTHORIZATION_ENDPOINT,
    TOKEN_ENDPOINT,
    END_SESSION_ENDPOINT,
    JWKS_ENDPOINT,
    ISSUER,
    OP_CONFIG_INITIATED
} from './constants/token';
import { getSessionParameter, removeSessionParameter, setSessionParameter } from "./session";
import Settings from '../../public/conf/Settings';

/**
 * Checks whether openid configuration initiated.
 *
 * @returns {boolean}
 */
 export const isOPConfigInitiated = (): boolean => {
    return getSessionParameter(OP_CONFIG_INITIATED) === "true";
};

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
 * Set openid configuration initiated.
 */
 export const setOPConfigInitiated = (): void => {
    setSessionParameter(OP_CONFIG_INITIATED, "true");
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
 * Initialize openid provider configuration.
 *
 * @param {string} wellKnownEndpoint openid provider configuration.
 * @param {boolean} forceInit whether to initialize the configuration again.
 * @returns {Promise<any>} promise.
 */
/* eslint-disable @typescript-eslint/no-explicit-any */
export const initOPConfiguration = (
    wellKnownEndpoint: string,
    forceInit: boolean
): Promise<any> => {

    if (!forceInit && isOPConfigInitiated()) {
        Promise.resolve("success");
    }

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
            setEndSessionEndpoint(Settings.logoutEndpoint);
            setJwksUri(response.data.jwks_uri);
            setIssuer(response.data.issuer);
            setOPConfigInitiated();

            return Promise.resolve("success");
        }).catch((error) => {
            return Promise.reject(error);
        });
};

/**
 * Reset openid provider configuration.
 */
export const resetOPConfiguration = (): void => {
    removeSessionParameter(AUTHORIZATION_ENDPOINT);
    removeSessionParameter(TOKEN_ENDPOINT);
    removeSessionParameter(END_SESSION_ENDPOINT);
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
 export const getJwksUri = (): string|null => {
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