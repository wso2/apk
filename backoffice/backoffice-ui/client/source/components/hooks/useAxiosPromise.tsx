import axios from "axios";
import { restApi } from "../../../public/conf/Settings";

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