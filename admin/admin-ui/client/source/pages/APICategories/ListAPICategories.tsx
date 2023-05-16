import Alert from '@mui/material/Alert';
import AlertTitle from '@mui/material/AlertTitle';
import Card from '@mui/material/Card';
import Stack from '@mui/material/Stack';
import axios from 'axios';
import Loader from 'components/Loader';
import PaginatedClientSide from 'components/data-table/PaginatedClientSide';
import React, { useEffect, useState } from 'react';
import { FormattedMessage } from 'react-intl';
import AddUpdateAPICategory from './AddUpdateAPICategory';
import DeleteAPICategory from './DeleteAPICategory';

export default function ListAPICategories() {

  const [data, setData] = useState<{ count: number; list: [{ id: string; name: string; description: string; numberOfAPIs: number; }]; } | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string>("");

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
  if (loading) {
    return <Loader />;
  }
  if (error || data === null || data === undefined) {
    return (
      <Alert severity='error'>
        <AlertTitle>Error</AlertTitle>
        There's an error when fetching API Categories — <strong>check it out!</strong>
      </Alert>
    );
  }
  if (data && data.count === 0) {
    return (
      <Card>
        <AddUpdateAPICategory
          id={undefined}
          nameProp={undefined}
          descriptionProp={undefined}
          updateList={fetchData}
        />
        <FormattedMessage
          id='AdminPages.ApiCategories.List.empty.content.apicategories'
          defaultMessage='Add API Category'
        />
      </Card>
    );
  }
  return (
    <>
      <AddUpdateAPICategory
        id={undefined}
        nameProp={undefined}
        descriptionProp={undefined}
        updateList={fetchData}
      />
      <PaginatedClientSide data={data.list} columns={columns} />
    </>
  )
}
