/* eslint-disable @typescript-eslint/ban-types */
/* eslint-disable no-empty-pattern */
import React from 'react';
import PaginatedClientSide from 'components/data-table/PaginatedClientSide'
import useApplicationRatePlans from "components/hooks/useApplicationRatePlans";
import Loader from "components/Loader";
// import { components, paths, operations } from 'types/Types';

type Props = {}

export default function ListApplicationRatePlans({ }: Props) {
  const { data, loading, error } = useApplicationRatePlans();
  const columns = React.useMemo(
    () => [
      {
        Header: 'Policy Name',
        accessor: 'policyName',
      },
      {
        Header: 'Display Name',
        accessor: 'displayName',
      },
      {
        Header: 'policyId',
        accessor: 'policyId',
      },
      {
        Header: 'isDeployed',
        accessor: 'isDeployed',
      },

    ],
    []
  )
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
    <PaginatedClientSide data={data.list} columns={columns} />
  )
}