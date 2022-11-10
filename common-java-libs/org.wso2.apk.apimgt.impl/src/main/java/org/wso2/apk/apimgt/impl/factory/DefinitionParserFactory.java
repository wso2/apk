/*
 * Copyright (c) 2022, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
 *
 * WSO2 Inc. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.wso2.apk.apimgt.impl.factory;

import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import org.wso2.apk.apimgt.api.APIDefinition;
import org.wso2.apk.apimgt.impl.APIConstants;
import org.wso2.apk.apimgt.impl.definitions.AsyncApiParser;
import org.wso2.apk.apimgt.impl.definitions.OAS2Parser;
import org.wso2.apk.apimgt.impl.definitions.OAS3Parser;


/**
 * Factory for getting definition parser instances.
 */
public class DefinitionParserFactory {
    private static final Log log = LogFactory.getLog(DefinitionParserFactory.class);

    //use getShape method to get object of type shape
    public APIDefinition getAPIDefinitionParser(APIConstants.ParserType parserType){
        switch (parserType) {
            case OAS3:
                return new OAS3Parser();
            case OAS2:
                return new OAS2Parser();
            case ASYNC:
                return new AsyncApiParser();
            default:
                return null;
        }
    }
}
