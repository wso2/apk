import { ID_TOKEN } from "./constants/token";
import { getEndSessionEndpoint } from "./op-config";
import { getSessionParameter } from "./session";
// eslint-disable-next-line @typescript-eslint/no-var-requires
const Settings = require('Settings');
/**
 * Handle user sign out.
 *
 * @returns {}
 */
/* eslint-disable @typescript-eslint/no-explicit-any */
export const sendSignOutRequest =  (redirectUri: string, sessionClearCallback: any): Promise<any> | undefined => {
    const logoutEndpoint = getEndSessionEndpoint();
    const logoutRequest = logoutEndpoint + "?client_id=" + Settings.idp.client_id + "&returnTo=" + redirectUri
    if (!logoutEndpoint || logoutEndpoint.trim().length === 0) {
        return Promise.reject(new Error("Invalid logout endpoint found."));
    }

    const idToken = getSessionParameter(ID_TOKEN);

    if (!idToken || idToken.trim().length === 0) {
        return Promise.reject(new Error("Invalid id_token found."));
    }

    sessionClearCallback();
    Promise.resolve("Logout success!");

    document.location.href = logoutRequest;
};