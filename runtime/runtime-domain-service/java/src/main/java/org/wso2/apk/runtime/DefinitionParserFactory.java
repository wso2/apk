package org.wso2.apk.runtime;

import org.wso2.apk.apimgt.api.APIDefinition;
import org.wso2.apk.apimgt.impl.definitions.AsyncApiParser;
import org.wso2.apk.apimgt.impl.definitions.OAS2Parser;
import org.wso2.apk.apimgt.impl.definitions.OAS3Parser;
import org.wso2.apk.runtime.model.API;

import java.util.ArrayList;
import java.util.List;

public class DefinitionParserFactory {
    private static List<APIDefinition> parsers = new ArrayList<>();

    private DefinitionParserFactory() {
    }

    static {
        parsers.add(new AsyncApiParser());
        parsers.add(new OAS3Parser());
        parsers.add(new OAS2Parser());

    }

    public static APIDefinition getParser(API api) {
        if ("HTTP".equals(api.getType())) {
            return new OAS3Parser();
        } else {
            return new AsyncApiParser();
        }
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
