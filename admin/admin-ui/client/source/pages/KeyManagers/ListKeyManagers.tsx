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
import { useIntl } from 'react-intl';
import { Grid, Typography } from '@mui/material';
import Button from '@mui/material/Button';
import AddCircleOutlineIcon from '@mui/icons-material/AddCircleOutline';

import PaginatedClientSide from 'components/data-table/PaginatedClientSide';
import GetKeyManagers from 'components/hooks/getKeyManagers';
import Loader from 'components/Loader';
import Switch from '@mui/material/Switch';
import Delete from './DeleteKeyManager';
import IconButton from '@mui/material/IconButton';
import EditIcon from '@mui/icons-material/Edit';

export default function ListKeyManagers() {
    const [trigger, setTrigger] = useState<boolean>(false);
    const { data, loading, error } = GetKeyManagers({ trigger: trigger, setTrigger: setTrigger });

    const intl = useIntl();

    const searchProps = {
        searchPlaceholder: intl.formatMessage({
            id: 'AdminPages.Organizations.List.search.default',
            defaultMessage: 'Search by KeyManager name',
        }),
    };

    const fetchData = () => {
        setTrigger(true);
    };

    const columns = [
        {
            Header: 'Name',
            accessor: 'name',
        },
        {
            Header: 'Type',
            accessor: 'type',
        },
        {
            Header: 'Issuer',
            accessor: 'issuer',
        },
        {
            Header: 'Enabled',
            accessor: 'enabled',
            Cell: ({ row }: { row: any }) => (
                <Switch
                    checked={row.original.enabled}
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
                    <IconButton color='primary' size='small'>
                        <EditIcon fontSize='small' />
                    </IconButton>
                    <Delete keyManagerId={row.original.id} updateList={fetchData} />
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
                        <Typography variant='h3'>Key Managers</Typography>
                    </Grid>
                    <Grid item display='grid'>
                        <Button variant='contained' style={{ margin: 'auto' }} startIcon={<AddCircleOutlineIcon />}>Add Key Managers</Button>
                    </Grid>
                </Grid>
            </div>
            <div>
                <PaginatedClientSide data={data.list} columns={columns} searchProps={searchProps} />
            </div>
        </div>
    );
}