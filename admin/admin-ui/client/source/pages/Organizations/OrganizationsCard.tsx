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
import PaginatedClientSide from 'components/data-table/PaginatedClientSide';
import Switch from '@mui/material/Switch';
import Stack from '@mui/material/Stack';
import IconButton from '@mui/material/IconButton';
import DeleteForeverIcon from '@mui/icons-material/DeleteForever';
import VisibilityOutlinedIcon from '@mui/icons-material/VisibilityOutlined';
import EditIcon from '@mui/icons-material/Edit';
import FormDialogBase from '../../components/@extended/FormDialogBase';
import DialogContentText from '@mui/material/DialogContentText';
import TextField from '@mui/material/TextField';
import { FormattedMessage } from 'react-intl';

export default function ListOrganizations() {

    const data1 = {
        "count": 1,
        "list": [
            {
                "name": "Business",
                "displayName": "Business",
                "organizationClaimValue": "0123-456-789-101",
                "enabled": true,
                "serviceNamespaces": [
                    "Test"
                ]
            },
            {
                "name": "Finance",
                "displayName": "Finance",
                "organizationClaimValue": "0123-456-789-103",
                "enabled": true,
                "serviceNamespaces": [
                    "Test1"
                ]
            },
            {
                "name": "Organization 123",
                "displayName": "Organization 123",
                "organizationClaimValue": "0123-456-789-104",
                "enabled": false,
                "serviceNamespaces": [
                    "Test1"
                ]
            }
        ]
    }
    const formSaveCallback = () => {
        console.log("Hello Dulith");
    };

    const [viewing, setViewing] = useState(false);

    const handleViewClickOpen = () => {
        setViewing(true);
    };

    const columns = React.useMemo(
        () => [
            {
                Header: 'Organization Name',
                accessor: 'displayName',
            },
            {
                Header: 'Organization Claim Value',
                accessor: 'organizationClaimValue',
            },
            {
                Header: 'Enabled',
                accessor: 'enabled',
                Cell: (e) => {
                    const [state, setState] = React.useState({
                        checked: e.row.original.enabled,
                    })
                    const handleChange = (event) => {
                        setState({ ...state, [event.target.name]: event.target.checked });
                    };
                    return (
                        <Switch
                            checked={state.checked}
                            onChange={handleChange}
                            color="primary"
                            name="checkedB"
                            key={e.row.id}
                        />
                    );
                },
            },
            {
                Header: 'Actions',
                accessor: 'actions',
                Cell: (e) => {
                    const [state, setState] = React.useState({
                        name: e.row.original.name,
                        displayName: e.row.original.displayName,
                        organizationClaimValue: e.row.original.organizationClaimValue,
                        serviceNamespaces: e.row.original.serviceNamespaces,
                    })
                    return (
                        <Stack direction="row" spacing={1}>
                            <FormDialogBase
                                title=''
                                saveButtonText=''
                                icon={<VisibilityOutlinedIcon />}
                                formSaveCallback={formSaveCallback}
                            >
                                <div>
                                    <span>
                                        <FormattedMessage id='AdminPages.ApiCategories.AddEdit.form.name' defaultMessage={state.name} />
                                    </span>
                                </div>
                            </FormDialogBase>
                            <FormDialogBase
                                title="Edit Organization"
                                saveButtonText='Save'
                                icon={<EditIcon />}
                                formSaveCallback={formSaveCallback}
                            >
                                <>
                                    <TextField
                                        autoFocus
                                        margin='dense'
                                        name='name'
                                        value={state.name}
                                        fullWidth
                                        variant='outlined'
                                        label={(
                                            <span>
                                                <FormattedMessage id='AdminPages.ApiCategories.AddEdit.form.name' defaultMessage='Name' />
                                                <span>*</span>
                                            </span>
                                        )}
                                        helperText={'Name of the Organization'}
                                    />
                                    <TextField
                                        margin='dense'
                                        name='displayName'
                                        value={state.displayName}
                                        fullWidth
                                        multiline
                                        label={(
                                            <span>
                                                <FormattedMessage id='AdminPages.ApiCategories.AddEdit.form.displayName' defaultMessage='Display Name' />
                                                <span>*</span>
                                            </span>
                                        )}
                                        helperText={'Display Name of the Organization'}
                                        variant='outlined'
                                    />
                                    <TextField
                                        margin='dense'
                                        name='organizationClaimValue'
                                        value={state.organizationClaimValue}
                                        fullWidth
                                        multiline
                                        label={(
                                            <span>
                                                <FormattedMessage id='AdminPages.ApiCategories.AddEdit.form.organizationClaimValue' defaultMessage='Organization Claim Value' />
                                                <span>*</span>
                                            </span>
                                        )}
                                        helperText={'Claim Value of the Organization'}
                                        variant='outlined'
                                    />
                                    <TextField
                                        margin='dense'
                                        name='serviceNamespaces'
                                        value={state.serviceNamespaces}
                                        fullWidth
                                        multiline
                                        label={(
                                            <span>
                                                <FormattedMessage id='AdminPages.ApiCategories.AddEdit.form.serviceNamespaces' defaultMessage='Service Namespaces' />
                                                <span>*</span>
                                            </span>
                                        )}
                                        helperText={'Service Namespaces of the Organization'}
                                        variant='outlined'
                                    />
                                </>

                            </FormDialogBase>
                            <FormDialogBase
                                title='Delete Organization?'
                                saveButtonText='Delete'
                                icon={<DeleteForeverIcon />}
                                formSaveCallback={formSaveCallback}                           
                            >
                                <DialogContentText>Are you sure you want to delete this Organization?</DialogContentText>
                            </FormDialogBase>
                        </Stack>
                    );
                },
            },

        ],
        []
    )
    return (
        <PaginatedClientSide data={data1.list} columns={columns} />
    )
}
