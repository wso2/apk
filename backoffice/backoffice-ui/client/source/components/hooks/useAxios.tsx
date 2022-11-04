import { useEffect, useState } from "react";
import axios from "axios";
import { restApi } from "../../../public/conf/Settings";
import { settings, tenantConfig } from "../../Mock";

type JSONValue =
    | string
    | number
    | boolean
    | { [x: string]: JSONValue }
    | Array<JSONValue>;
interface useFetchProps {
    url: string,
}

const useAxios = ({ url }: useFetchProps) => {

    // mock
    const isLoading = false;
    const isError = false;
    if (url == '/settings') {
        return { data : settings, loading : isLoading, error : isError};
    }
    if (url == '/tenant-info/YWRtaW5AY2FyYm9uLnN1cGVy') {
        return { data : tenantConfig, loading : isLoading, error : isError};
    }
    // mock end

    const [data, setData] = useState<any | null>(null);
    const [loading, setLoading] = useState<JSONValue | null>(true);
    const [error, setError] = useState<string>("");
    let fetchUrl = '';
    if (restApi.endsWith('/') && url.startsWith('/')) {
        fetchUrl = restApi.slice(0, -1) + url;
    } else if (!restApi.endsWith('/') && !url.startsWith('/')) {
        fetchUrl = restApi + '/' + url;
    } else {
        fetchUrl = restApi + url;
    }

    useEffect(() => {
        axios.get(fetchUrl)
            .then((res) => {
                setData(res.data);
            })
            .catch((e) => {
                setError(e);
            })
            .finally(() => {
                setLoading(false);
            })
    }, [url]);

    return { data, loading, error };
};

export default useAxios;