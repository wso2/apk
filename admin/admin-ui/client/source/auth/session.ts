import { Semaphore } from "await-semaphore";
import {
    ACCESS_TOKEN,
    ID_TOKEN,
    REQUEST_PARAMS,
    REFRESH_TOKEN,
    SCOPE,
    TOKEN_TYPE,
    ACCESS_TOKEN_EXPIRE_IN,
    ACCESS_TOKEN_ISSUED_AT
} from "./constants/token";
import { sendRefreshTokenRequest } from "./sign-in";
import { TokenResponseInterface } from "./models/token-response";

/**
 * Semaphore used for synchronizing the refresh token requests.
 */
const semaphore = new Semaphore(1);

/**
 * Remove parameter from session storage.
 *
 * @param {string} key.
 */
export const removeSessionParameter = (key: string): void => {
    sessionStorage.removeItem(key);
};

/**
 * Set parameter to session storage.
 *
 * @param {string} key.
 * @param value value.
 */
export const setSessionParameter = (key: string, value: string): void => {
    sessionStorage.setItem(key, value);
};

/**
 * Get parameter from session storage.
 *
 * @param {string} key.
 * @returns {string | null} parameter value or null.
 */
export const getSessionParameter = (key: string): string => {
    return sessionStorage.getItem(key) || '';
};



/**
 * End authenticated user session.
 */
 export const endAuthenticatedSession = (): void => {
    removeSessionParameter(ACCESS_TOKEN);
    removeSessionParameter(ID_TOKEN);
    removeSessionParameter(ACCESS_TOKEN_EXPIRE_IN);
    removeSessionParameter(ACCESS_TOKEN_ISSUED_AT);
    removeSessionParameter(SCOPE);
    removeSessionParameter(REFRESH_TOKEN);
    removeSessionParameter(TOKEN_TYPE);
};

/**
 * Initialize authenticated user session.
 *
 * @param {TokenResponseInterface} tokenResponse.
 * @param authenticatedUser authenticated user.
 */
export const initUserSession = (tokenResponse: TokenResponseInterface): void => {
    endAuthenticatedSession();
    setSessionParameter(ACCESS_TOKEN, tokenResponse.accessToken);
    setSessionParameter(ID_TOKEN, tokenResponse.idToken);
    setSessionParameter(ACCESS_TOKEN_EXPIRE_IN, tokenResponse.expiresIn);
    setSessionParameter(ACCESS_TOKEN_ISSUED_AT, (Date.now() / 1000).toString());
    setSessionParameter(SCOPE, tokenResponse.scope);
    setSessionParameter(REFRESH_TOKEN, tokenResponse.refreshToken);
    setSessionParameter(TOKEN_TYPE, tokenResponse.tokenType);
};

/**
 * Use this method when need to get a token using the refresh token.
 * Get access token.
 *
 * @returns {Promise<string>} access token.
 */
export const getAccessToken = (): Promise<string> => {
    const accessToken = getSessionParameter(ACCESS_TOKEN);
    const expiresIn = getSessionParameter(ACCESS_TOKEN_EXPIRE_IN);
    const issuedAt = getSessionParameter(ACCESS_TOKEN_ISSUED_AT);

    if (!accessToken || accessToken.trim().length === 0 || !expiresIn || expiresIn.length === 0 || !issuedAt
        || issuedAt.length === 0) {
        endAuthenticatedSession();

        return Promise.reject(new Error("Invalid user session."));
    }

    function getValidityPeriod(): number {
        const currentExpiresIn = getSessionParameter(ACCESS_TOKEN_EXPIRE_IN);
        const currentIssuedAt = getSessionParameter(ACCESS_TOKEN_ISSUED_AT);

        return (parseInt(currentIssuedAt, 10) + parseInt(currentExpiresIn, 10)) - Math.floor(Date.now() / 1000);
    }

    let validityPeriod = getValidityPeriod();

    if (validityPeriod <= 300) {

        return semaphore.use(() => {
            validityPeriod = getValidityPeriod();
            if (validityPeriod <= 300) {
                const requestParams = JSON.parse(getSessionParameter(REQUEST_PARAMS));
                return sendRefreshTokenRequest(requestParams, getSessionParameter(REFRESH_TOKEN))
                    .then((tokenResponse: any) => {
                        initUserSession(tokenResponse);
                        return Promise.resolve(tokenResponse.accessToken);
                    }).catch((error: any) => {
                        return Promise.reject(error);
                    });
            } else {
                return Promise.resolve(getSessionParameter(ACCESS_TOKEN));
            }
        });
    } else {
        return Promise.resolve(accessToken);
    }
};