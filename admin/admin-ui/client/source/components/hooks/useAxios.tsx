import { useEffect, useState } from "react";
import axios from "axios";
import { restApi } from "../../../public/conf/Settings";

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