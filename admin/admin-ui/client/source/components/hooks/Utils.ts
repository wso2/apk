// eslint-disable-next-line @typescript-eslint/no-var-requires
const Settings = require('Settings');
const restApi = Settings.app.rest_api;

// Function to return formatted URL
export const getFormattedUrl = (url: string) => {
    let fetchUrl = '';
    if (restApi.endsWith('/') && url.startsWith('/')) {
        fetchUrl = restApi.slice(0, -1) + url;
    } else if (!restApi.endsWith('/') && !url.startsWith('/')) {
        fetchUrl = restApi + '/' + url;
    } else {
        fetchUrl = restApi + url;
    }
    return fetchUrl;
}