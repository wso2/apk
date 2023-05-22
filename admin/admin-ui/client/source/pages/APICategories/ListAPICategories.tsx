/*
 *  Copyright (c) 2023, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 */

import { Grid, Typography } from '@mui/material';
import { default as Alert, default as MuiAlert } from '@mui/material/Alert';
import AlertTitle from '@mui/material/AlertTitle';
import Card from '@mui/material/Card';
import Snackbar from '@mui/material/Snackbar';
import Stack from '@mui/material/Stack';
import axios from 'axios';
import Loader from 'components/Loader';
import PaginatedClientSide from 'components/data-table/PaginatedClientSide';
import React, { useEffect, useState } from 'react';
import { FormattedMessage, useIntl } from 'react-intl';
import AddUpdateAPICategory from './AddUpdateAPICategory';
import DeleteAPICategory from './DeleteAPICategory';

export default function ListAPICategories() {

  const [data, setData] = useState<{ count: number; list: [{ id: string; name: string; description: string; numberOfAPIs: number; }]; } | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string>("");
  const [snackbarOpen, setSnackbarOpen] = useState(false);
  const [formattedMessage, setFormattedMessage] = useState({
    id: '',
    defaultMessage: '',
  });

  const intl = useIntl();

  const fetchData = () => {
    setLoading(true);
    axios('/api/am/admin/api-categories', {
      method: 'GET',
      withCredentials: true,
    }).then((res) => {
      setData(res.data);
    }).catch((err) => {
      setError(err);
    }).finally(() => {
      setLoading(false);
    });
  };

  useEffect(() => {
    fetchData();
  }, []);

  const columns = React.useMemo(
    () => [
      {
        Header: 'Category Name',
        accessor: 'name',
      },
      {
        Header: 'Description',
        accessor: 'description',
      },
      {
        Header: 'Number Of APIs',
        accessor: 'numberOfAPIs',
      },
      {
        Header: 'Actions',
        accessor: 'actions',
        Cell: (e) => {
          return (
            <Stack direction="row" spacing={1}>
              <AddUpdateAPICategory
                id={e.row.original.id}
                nameProp={e.row.original.name}
                descriptionProp={e.row.original.description}
                updateList={fetchData}
              />
              <DeleteAPICategory
                id={e.row.original.id}
                updateList={fetchData}
              />
            </Stack>
          );
        },
      },
    ],
    []
  )

  const searchProps = {
    searchPlaceholder: intl.formatMessage({
      id: 'AdminPages.APICategories.List.search.default',
      defaultMessage: 'Search by Category name',
    }),
  };

  if (loading) {
    return <Loader />;
  }
  if (error || data === null || data === undefined) {
    return (
      <>
        <Typography variant='h3'>API Categories</Typography>
        <br />
        <Alert severity='error'>
          <AlertTitle>Error</AlertTitle>
          There's an error when fetching API Categories — <strong>check it out!</strong>
        </Alert>
      </>
    );
  }
  if (data && data.count === 0) {
    return (
      <div>
        <div>
          <Grid container direction='row' justifyContent='space-between'>
            <Grid item sx={{ mt: 2, mb: 2 }}>
              <Typography variant='h3'>API Categories</Typography>
            </Grid>
            <Grid item display='grid'>
              <AddUpdateAPICategory
                id={undefined}
                nameProp={undefined}
                descriptionProp={undefined}
                updateList={fetchData}
              />
            </Grid>
          </Grid>
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
  return (
    <div>
      <div>
        <Grid container direction='row' justifyContent='space-between'>
          <Grid item sx={{ mt: 2, mb: 2 }}>
            <Typography variant='h3'>API Categories</Typography>
          </Grid>
          <Grid item display='grid'>
            <AddUpdateAPICategory
              id={undefined}
              nameProp={undefined}
              descriptionProp={undefined}
              updateList={fetchData}
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
  )
}
