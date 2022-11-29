/*
 * Copyright (c) 2020, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
 *
 * WSO2 Inc. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
import React from 'react';
import { Routes, Route } from 'react-router-dom';
import List from 'AppComponents/Throttling/Advanced/List';
import AddEdit from 'AppComponents/Throttling/Advanced/AddEdit';
import ResourceNotFound from 'AppComponents/Base/Errors/ResourceNotFound';

/**
 * Render a list
 * @returns {JSX} Header AppBar components.
 */
function AdvancedThrottlePolicies() {
    return (
        <Routes>
            <Route exact path='/throttling/advanced' component={List} />
            <Route exact path='/throttling/advanced/create' component={AddEdit} />
            <Route path='/throttling/advanced/:id/' component={AddEdit} />
            <Route component={ResourceNotFound} />
        </Routes>
    );
}

export default AdvancedThrottlePolicies;
