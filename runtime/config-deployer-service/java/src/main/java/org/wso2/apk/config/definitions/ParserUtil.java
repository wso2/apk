package org.wso2.apk.config.definitions;


public class ParserUtil {
    private ParserUtil() {
    }
    /**
     * Find scope object in a set based on the key
     *
     * @param scopes - Set of scopes
     * @param key    - Key to search with
     * @return Scope - scope object
     */
    public static String findScopeByKey(String[] scopes, String key) {

        for (String scope : scopes) {
            if (scope.equals(key)) {
                return scope;
            }
        }
        return null;
    }


}
