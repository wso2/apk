import { useEffect, useState } from "react";
import axios from "axios";
import { settings, tenant, apiCategories, applicationThrottlePolicies } from './dummyPayloads';
// eslint-disable-next-line @typescript-eslint/no-var-requires
// const Settings = require('Settings');
// const restApi = Settings.server.restApi;

type JSONValue =
    | string
    | number
    | boolean
    | { [x: string]: JSONValue }
    | Array<JSONValue>;
interface useFetchProps {
    url: string,
}

/** =================== Remove the following logics when the actual impl is there */
const getDummyData = (url) => {
    return new Promise(function (resolve, reject) {
        const resFunction = () => {
            if (url === '/settings') {
                resolve(settings)
            }
            if (url === '/tenant-info/YWRtaW5AY2FyYm9uLnN1cGVy') {
                resolve(tenant)
            }
            if (url === '/api-categories') {
                resolve(apiCategories)
            }
            if (url === '/throttling/policies/application') {
                // Duplicate applicationThrottlePolicies.list to test pagination 100 times
                const oneItem = applicationThrottlePolicies.list[0];
                for (let i = 0; i < 100; i++){
                    applicationThrottlePolicies.list.push(oneItem);
                }
                resolve(applicationThrottlePolicies)
            }
        }
        setTimeout(resFunction, 0);
    });
}
/** ========================= End of sleep function =========================== */

const useAxios = ({ url }: useFetchProps) => {
    const [data, setData] = useState<any | null>(null);
    const [loading, setLoading] = useState<JSONValue | null>(true);
    const [error, setError] = useState<string>("");
    // let fetchUrl = '';
    // if (restApi.endsWith('/') && url.startsWith('/')) {
    //     fetchUrl = restApi.slice(0, -1) + url;
    // } else if (!restApi.endsWith('/') && !url.startsWith('/')) {
    //     fetchUrl = restApi + '/' + url;
    // } else {
    //     fetchUrl = restApi + url;
    // }

    useEffect(() => {
        /** =================== Remove the following logics when the actual impl is there */
        getDummyData(url)
            .then((data) => {
                setData(data);
            })
            .finally(() => {
                setLoading(false);
            });
        /** ========================= End of dummy data return =========================== */
        // axios.get(fetchUrl)
        //     .then((res) => {
        //         setData(res.data);
        //     })
        //     .catch((e) => {
        //         setError(e);
        //     })
        //     .finally(() => {
        //         setLoading(false);
        //     })
    }, [url]);

    return { data, loading, error };
};

export default useAxios;