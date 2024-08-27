package org.wso2.apk.config;

import org.wso2.apk.config.api.APIDefinition;
import org.wso2.apk.config.definitions.AsyncApiParser;
import org.wso2.apk.config.definitions.OAS3Parser;
import org.wso2.apk.config.definitions.ProtoParser;
import org.wso2.apk.config.model.API;

import java.util.ArrayList;
import java.util.List;

public class DefinitionParserFactory {
    private static final List<APIDefinition> parsers = new ArrayList<>();

    private DefinitionParserFactory() {
    }

    static {
        parsers.add(new AsyncApiParser());
        parsers.add(new OAS3Parser());
        parsers.add(new ProtoParser());
    }

    public static APIDefinition getParser(API api) {
        if (APIConstants.ParserType.REST.name().equals(api.getType())
                || APIConstants.ParserType.GRAPHQL.name().equals(api.getType())) {
            return new OAS3Parser();
        } else if (APIConstants.ParserType.GRPC.name().equals(api.getType().toUpperCase())) {
            return new ProtoParser();
        } else if (APIConstants.ParserType.ASYNC.name().equals(api.getType())) {
            return new AsyncApiParser();
        }
        return null;
    }

    public static APIDefinition getParser(String apiType) {
        if (APIConstants.ParserType.REST.name().equals(apiType)
                || APIConstants.ParserType.GRAPHQL.name().equals(apiType)) {
            return new OAS3Parser();
        } else if (APIConstants.ParserType.GRPC.name().equals(apiType.toUpperCase())) {
            return new ProtoParser();
        } else if ("ASYNC".equals(apiType)) {
            return new AsyncApiParser();
        }
        return null;
    }

    public static APIDefinition getValidatedParser(String definition) {
        for (APIDefinition parser : parsers) {
            if (parser.canHandleDefinition(definition)) {
                return parser;
            }
        }
        return null;
    }
}
