package org.wso2.apk.runtime;

import org.wso2.apk.runtime.api.APIDefinition;
import org.wso2.apk.runtime.definitions.AsyncApiParser;
import org.wso2.apk.runtime.definitions.OAS3Parser;
import org.wso2.apk.runtime.model.API;

import java.util.ArrayList;
import java.util.List;

public class DefinitionParserFactory {
    private static final List<APIDefinition> parsers = new ArrayList<>();

    private DefinitionParserFactory() {
    }

    static {
        parsers.add(new AsyncApiParser());
        parsers.add(new OAS3Parser());

    }

    public static APIDefinition getParser(API api) {
        if ("HTTP".equals(api.getType())) {
            return new OAS3Parser();
        } else if ("ASYNC".equals(api.getType())){
            return new AsyncApiParser();
        }
        return null;
    }
    public static APIDefinition getParser(String apiType) {
        if ("HTTP".equals(apiType)) {
            return new OAS3Parser();
        } else if ("ASYNC".equals(apiType)){
            return new AsyncApiParser();
        }
        return null;
    }

    public static APIDefinition getValidatedParser(String definition) {
        for (APIDefinition parser : parsers) {
            if (parser.canHandleDefinition(definition)){
                return parser;
            }
        }
        return null;
    }
}
