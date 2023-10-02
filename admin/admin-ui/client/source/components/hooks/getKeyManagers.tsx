import axios from "axios";
import { useEffect, useState } from "react";

type JSONValue =
    | string
    | number
    | boolean
    | { [x: string]: JSONValue }
    | Array<JSONValue>;

interface orgProps { trigger: boolean, setTrigger: (boolean) => void }
 
const getKeyManagers =  ({trigger, setTrigger} : orgProps ) => {
    const [data, setData] = useState<any | null>(null);
    const [loading, setLoading] = useState<JSONValue | null>(true);
    const [error, setError] = useState<string>("");

    const viewData = () => {
        axios('/api/am/admin/key-managers', {
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
                setTrigger(false);
            });
    }
    
    useEffect(() => {
        viewData();
    }, []);

    useEffect(() => {
        if (trigger) {
            viewData();
        }
    }, [trigger]);

    return { data, loading, error };
};

export default getKeyManagers;
