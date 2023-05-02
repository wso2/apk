import axios from "axios";
// eslint-disable-next-line @typescript-eslint/no-var-requires
const Settings = require('Settings');
const restApi = Settings.server.restApi;
interface useFetchProps {
    url: string,
}

const useAxiosPromise = ({ url }: useFetchProps) => {
    let fetchUrl = '';
    if (restApi.endsWith('/') && url.startsWith('/')) {
        fetchUrl = restApi.slice(0, -1) + url;
    } else if (!restApi.endsWith('/') && !url.startsWith('/')) {
        fetchUrl = restApi + '/' + url;
    } else {
        fetchUrl = restApi + url;
    }
    const promise = axios.get(fetchUrl)
        .then((res) => {
            return (res.data);
        })


    return promise;
};

export default useAxiosPromise;