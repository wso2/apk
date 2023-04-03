import axios from "axios";
import { useEffect, useState } from "react";
import { getFormattedUrl } from './Utils';

type JSONValue =
    | string
    | number
    | boolean
    | { [x: string]: JSONValue }
    | Array<JSONValue>;


const useApplicationRatePlans = () => {
    const [data, setData] = useState<any | null>(null);
    const [loading, setLoading] = useState<JSONValue | null>(true);
    const [error, setError] = useState<string>("");

    //const fetchUrl = getFormattedUrl('application-rate-plans');
    useEffect(() => {
        axios('/api', {
            method: 'GET',
            withCredentials: true,
        })
        .then((res) => {
            setData(res.data);
        })
        .catch((err) => {
            setError(err);
        })
        .finally(() => {
            setLoading(false);
        });
    }, []);


    return { data, loading, error };
};

export default useApplicationRatePlans;