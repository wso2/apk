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

/**
 * Render organization data table.
 * @param {JSON} props component props.
 * @returns {TSX} Loading animation.
 */
import React, { useState } from 'react';
import Stack from '@mui/material/Stack';
import { useIntl, FormattedMessage } from 'react-intl';
import { Grid, Typography } from '@mui/material';

import axios from 'axios';
import PaginatedClientSide from 'components/data-table/PaginatedClientSide';
import GetOrganizations from 'components/hooks/getOrganizations';
import Loader from 'components/Loader';
import Delete from './DeleteOrganization';
import AddEdit from './AddEditOrganization';
import TooltipPopup from './TooltipPopup';
import Switch from '@mui/material/Switch';
import Snackbar from '@mui/material/Snackbar';

import MuiAlert from '@mui/material/Alert';

export default function ListOrganizations() {
    const [trigger, setTrigger] = useState<boolean>(false);
    const { data, loading, error } = GetOrganizations({ trigger: trigger, setTrigger: setTrigger });
    const [snackbarOpen, setSnackbarOpen] = useState(false);
    const [formattedMessage, setFormattedMessage] = useState({
        id: '',
        defaultMessage: '',
    });
    const intl = useIntl();

    const fetchData = () => {
        setTrigger(true);
    };

    const createOrganization = (organization) => {
        axios
            .post('/api/am/admin/organizations/', {
                ...organization
            }, {
                withCredentials: true,
            })
            .then(() => {
                setFormattedMessage({
                    id: 'AdminPages.Organizations.AddEdit.form.add.successful',
                    defaultMessage: 'Organization added successfully'
                });
                setSnackbarOpen(true);
            })
            .catch((error) => {
                throw error.response.body.description;
            })
            .finally(() => {
                fetchData();
            });
    };

    const updateOrganization = (organization) => {
        axios
            .put('/api/am/admin/organizations/' + organization.id, {
                ...organization
            }, {
                withCredentials: true,
            })
            .then(() => {
                setFormattedMessage({
                    id: 'AdminPages.Organizations.AddEdit.form.edit.successful',
                    defaultMessage: 'Organization edited successfully'
                });
                setSnackbarOpen(true);
            })
            .catch((error) => {
                throw error.response.body.description;
            })
            .finally(() => {
                fetchData();
            });
    };

    const handleEnabledChange = (e: React.ChangeEvent<HTMLInputElement>, organization) => {
        updateOrganization({ ...organization, [e.target.name]: e.target.checked} );
    };

    const searchProps = {
        searchPlaceholder: intl.formatMessage({
            id: 'AdminPages.Organizations.List.search.default',
            defaultMessage: 'Search by Organization name',
        }),
    };

    const columns = [
        {
            Header: 'Organization Name',
            accessor: 'displayName',
            sortable: true,
            Cell: ({ row }: { row: any }) => (
                <Stack direction="row" spacing={1} style={{ alignItems: 'center' }}>
                    <TooltipPopup org={row.original} />
                </Stack>
            ),
        },
        {
            Header: 'Organization Claim Value',
            accessor: 'organizationClaimValue',
        },
        {
            Header: 'Enabled',
            accessor: 'enabled',
            Cell: ({ row }: { row: any }) => (
                <Switch
                    checked={row.original.enabled}
                    onChange={(e) => handleEnabledChange(e, row.original)}
                    color='primary'
                    name='enabled'
                    key={row.original.id}
                    size='small'
                />
            ),
        },
        {
            Header: 'Actions',
            accessor: 'actions',
            Cell: ({ row }: { row: any }) => (
                <Stack direction="row" spacing={1}>
                    <AddEdit
                        datarow={row.original}
                        updateList={(updatedOrganization) => updateOrganization(updatedOrganization)}
                    />
                    <Delete orgId={row.original.id} updateList={fetchData} />
                </Stack>
            ),
        },
    ];

    if (error) {
        return <div>Error</div>;
    }
    if (loading) {
        return <Loader />;
    }
    if (data && data.length === 0) {
        return <div>No data</div>;
    }
    return (
        <div>
            <div>
                <Grid container direction='row' justifyContent='space-between'>
                    <Grid item sx={{ mt: 2, mb: 2 }}>
                        <Typography variant='h3'>Organizations</Typography>
                    </Grid>
                    <Grid item display='grid'>
                        <AddEdit
                            datarow={null}
                            updateList={(newOrganization) => createOrganization(newOrganization)}
                        />
                    </Grid>
                </Grid>
            </div>
            <div>
                <PaginatedClientSide data={data.list} columns={columns} searchProps={searchProps} />
            </div>
            <Snackbar open={snackbarOpen} autoHideDuration={10000} onClose={() => setSnackbarOpen(false)}>
                <MuiAlert onClose={() => setSnackbarOpen(false)} severity='success' sx={{ width: '100%' }}>
                    {
                        <FormattedMessage
                            id={formattedMessage.id}
                            defaultMessage={formattedMessage.defaultMessage}
                        /> 
                    }
                </MuiAlert>
            </Snackbar>
        </div>
    );
}
