function getAPIList() returns string?|APIList|error {
    API[]|error? apis = getAPIsDAO();
    if apis is API[] {
        int count = apis.length();
        APIList apisList = {count: count, list: apis};
        return apisList;
    } else {
        return apis;
    }
}
