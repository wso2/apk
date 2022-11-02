import Base64 from "crypto-js/enc-base64";
import WordArray from "crypto-js/lib-typedarrays";
import sha256 from "crypto-js/sha256";
import { KEYUTIL, KJUR } from "jsrsasign";
import { JWKInterface } from "./models/crypto";
/**
 * Get URL encoded string.
 *
 * @param {any} value.
 * @returns {string} base 64 url encoded value.
 */
export const base64URLEncode = (value: any): string => {
    return Base64.stringify(value)
        .replace(/\+/g, "-")
        .replace(/\//g, "_")
        .replace(/=/g, "");
};

/**
 * Generate code verifier.
 *
 * @returns {string} code verifier.
 */
 export const getCodeVerifier = (): string => {
    return base64URLEncode(WordArray.random(32));
};

/**
 * Derive code challenge from the code verifier.
 *
 * @param {string} verifier.
 * @returns {string} code challenge.
 */
 export const getCodeChallenge = (verifier: string): string => {
    return base64URLEncode(sha256(verifier));
};

/**
 * Get the supported signing algorithms for the id_token.
 *
 * @returns {string[]} array of supported algorithms.
 */
 export const getSupportedSignatureAlgorithms = (): string[] => {
    return ["RS256", "RS512", "RS384", "PS256", "HS256"];
};

/**
 * Get JWK used for the id_token
 *
 * @param {string} jwtHeader header of the id_token.
 * @param {JWKInterface[]} keys jwks response.
 * @returns {any} public key.
 */
/* eslint-disable @typescript-eslint/no-explicit-any */
export const getJWKForTheIdToken = (jwtHeader: string, keys: JWKInterface[]): Error|any => {
    const headerJSON = JSON.parse(atob(jwtHeader));

    for (const key of keys) {
        if (headerJSON.kid === key.kid) {
            return KEYUTIL.getKey({ kty: key.kty, e: key.e, n: key.n });
        }
    }

    throw new Error("Failed to find the 'kid' specified in the id_token. 'kid' found in the header : "
        + headerJSON.kid + ", Expected values: " + keys.map((key) => key.kid).join(", "));
};

/**
 * Verify id token.
 *
 * @param idToken id_token received from the IdP.
 * @param jwk public key used for signing.
 * @param {string} clientID app identification.
 * @param {string} issuer id_token issuer.
 * @returns {any} whether the id_token is valid.
 */
/* eslint-disable @typescript-eslint/no-explicit-any */
export const isValidIdToken = (idToken: any, jwk: any, clientID: string, issuer: string): any => {
    return KJUR.jws.JWS.verifyJWT(idToken, jwk, {
        alg: getSupportedSignatureAlgorithms(),
        aud: [clientID],
        gracePeriod: 3600,
        iss: [issuer]
    });
};