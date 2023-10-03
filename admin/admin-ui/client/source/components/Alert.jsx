/*
 * Copyright (c) 2019, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
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
import toast from 'react-hot-toast';
import React from 'react';

export default {
    info: toast,
    success: toast.success,
    error: toast.error,
    warning: (message, options) => toast(message, {
        style: { background: '#ffd891' },
        icon: (
            <span
                style={{
                    fontSize: '21px',
                    color: '#c17e03',
                    fontWeight: 'bold',
                }}
            >
                &#9888;
            </span>
        ),
        ...options,
    }),
    loading: toast.promise,
};
