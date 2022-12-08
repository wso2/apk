package org.wso2.apk.runtime.definitions;

import org.wso2.apk.runtime.model.Scope;

import java.util.Set;

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
    public static Scope findScopeByKey(Set<Scope> scopes, String key) {

        for (Scope scope : scopes) {
            if (scope.getKey().equals(key)) {
                return scope;
            }
        }
        return null;
    }


}
