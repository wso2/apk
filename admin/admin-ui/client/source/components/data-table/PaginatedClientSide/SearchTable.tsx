/*
 * Copyright (c) 2023, WSO2 LLC. (https://www.wso2.com) All Rights Reserved.
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

import React, { useState } from 'react';
import TextField from '@mui/material/TextField';
import Autocomplete from '@mui/material/Autocomplete';
import AppBar from '@mui/material/AppBar';
import Toolbar from '@mui/material/Toolbar';
import Grid from '@mui/material/Grid';
import SearchIcon from '@mui/icons-material/Search';

export default function SearchTable({filter, setFilter, searchProps: {searchPlaceholder} }) {

    return (
        <AppBar position='static' color='default' elevation={0}>
            <Toolbar>
                <Grid container spacing={2} alignItems='center' style={{
                    width: '25%'
                }}>
                    <Grid item>
                        {<SearchIcon color='inherit' />}
                    </Grid>
                    <Grid item xs>
                        {
                            <TextField
                                fullWidth
                                placeholder={searchPlaceholder}
                                onChange={(e) => setFilter(e.target.value)}
                                value={filter || ''}
                            />
                        }
                    </Grid>
                </Grid>
            </Toolbar>
        </AppBar>
    );
}
