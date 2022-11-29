/*
 * Copyright (c) 2021, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
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

import React, {
    useReducer, useState, Suspense, lazy, useEffect,
} from 'react';
import { useNavigate } from 'react-router-dom';
import PropTypes from 'prop-types';
import { FormattedMessage } from 'react-intl';
import { makeStyles } from '@material-ui/core/styles';
import {
    Box, Grid, Button,
} from '@material-ui/core';
import { green } from '@material-ui/core/colors';
import API from 'AppData/api';
import Alert from 'AppComponents/Shared/Alert';
import { Progress } from 'AppComponents/Shared';
import ContentBase from 'AppComponents/AdminPages/Addons/ContentBase';

const MonacoEditor = lazy(() => import('react-monaco-editor' /* webpackChunkName: "TeantConfAddMonacoEditor" */));

const useStyles = makeStyles((theme) => ({
    root: {
        marginBottom: theme.spacing(15),
    },
    error: {
        color: theme.palette.error.dark,
    },
    dialog: {
        minWidth: theme.spacing(150),

    },
    tenantConfHeading: {
        marginBottom: theme.spacing(1),
    },
    showSampleButton: {
        marginTop: theme.spacing(2),
    },
    helperText: {
        color: green[600],
        fontSize: theme.spacing(1.6),
        marginLeft: theme.spacing(1),
    },
    infoBox: {
        marginTop: theme.spacing(1),
        marginBottom: theme.spacing(2),
    },
    buttonBox: {
        marginTop: theme.spacing(2),
        marginBottom: theme.spacing(2),
        marginLeft: theme.spacing(1),
    },
    saveButton: {
        marginRight: theme.spacing(2),
    },
}));

/**
 * Reducer
 * @param {JSON} state The second number.
 * @returns {Promise}
 */
function reducer(state, { field, value }) {
    switch (field) {
        case 'tenantConf':
            return { ...state, [field]: value };
        case 'tenantConfSchema':
            return { ...state, [field]: value };
        case 'editDetails':
            return value;
        case 'reset':
            return { tenantConf: '', tenantConfSchema: '' };
        default:
            return state;
    }
}

/**
 * Render a list
 * @returns {JSX} Header AppBar components.
 */
function TenantConfSave() {
    const classes = useStyles();
    const navigate = useNavigate();
    const [initialState] = useState({
        tenantConf: '',
        tenantConfSchema: '',
    });
    const [state, dispatch] = useReducer(reducer, initialState);
    const {
        tenantConf, tenantConfSchema,
    } = state;
    const restApi = new API();

    useEffect(() => {
        let tenantConfVal;
        let tenantConfSchemaVal;
        restApi.tenantConfSchemaGet().then((result) => {
            tenantConfSchemaVal = result.body;
            dispatch({ field: 'tenantConfSchema', value: tenantConfSchemaVal });
        });
        restApi.tenantConfGet().then((result) => {
            tenantConfVal = JSON.stringify(result.body, null, '\t');
            dispatch({ field: 'tenantConf', value: tenantConfVal });
        });
    }, []);

    const editorWillMount = (monaco) => {
        const schemaVal = state.tenantConfSchema;
        monaco.languages.json.jsonDefaults.setDiagnosticsOptions({
            completion: true,
            validate: true,
            format: true,
            schemas: [{
                uri: 'http://myserver/foo-schema.json',
                fileMatch: ['*'],
                schema: schemaVal,
            }],
        });
    };

    const tenantConfOnChange = (newValue) => {
        dispatch({ field: 'tenantConf', value: newValue });
    };

    const cancelUpdate = () => {
        let tenantConfVal;
        restApi.tenantConfGet().then((result) => {
            tenantConfVal = JSON.stringify(result.body, null, '\t');
            const editState = {
                tenantConf: tenantConfVal,
            };
            dispatch({ field: 'editDetails', value: editState });
        });
    };

    const formSaveCallback = () => {
        const tenantConfiguration = state.tenantConf;
        const promisedUpdateTenantConf = restApi.updateTenantConf(tenantConfiguration);
        return promisedUpdateTenantConf
            .then(() => {
                Alert.success(
                    <FormattedMessage
                        id='Settings.Advanced.TenantConf.edit.success'
                        defaultMessage='Advanced Configuration saved successfully'
                    />,
                );
                navigate('/settings/advanced');
            })
            .catch((error) => {
                const { response } = error;
                if (response.body) {
                    Alert.error(response.body.description);
                }
                console.error(error);
            });
    };

    return (
        <ContentBase
            pageStyle='full'
            title={(
                <FormattedMessage
                    id='Settings.Advanced.TenantConfSave.title.save'
                    defaultMessage='Advanced Configurations'
                />
            )}
        >
            <Box component='div' m={2} className={classes.root} name={tenantConfSchema}>
                <Grid container spacing={3}>
                    <Grid item xs={12} md={12} lg={12}>
                        <Suspense fallback={<Progress />}>
                            <MonacoEditor
                                language='json'
                                height='615px'
                                theme='vs-dark'
                                value={tenantConf}
                                onChange={tenantConfOnChange}
                                editorWillMount={editorWillMount}
                            />
                        </Suspense>
                    </Grid>
                    <Box component='span' className={classes.buttonBox}>
                        <Button
                            variant='contained'
                            color='primary'
                            className={classes.saveButton}
                            onClick={formSaveCallback}
                            data-testid='monaco-editor-save'
                        >
                            <FormattedMessage
                                id='Settings.Advanced.TenantConfSave.form.save'
                                defaultMessage='Save'
                            />
                        </Button>
                        <Button
                            variant='contained'
                            onClick={cancelUpdate}
                        >
                            <FormattedMessage
                                id='Settings.Advanced.TenantConfSave.form.cancel'
                                defaultMessage='Cancel'
                            />
                        </Button>
                    </Box>
                </Grid>
            </Box>
        </ContentBase>
    );
}

TenantConfSave.propTypes = {
    title: PropTypes.shape({}).isRequired,
};

export default TenantConfSave;
