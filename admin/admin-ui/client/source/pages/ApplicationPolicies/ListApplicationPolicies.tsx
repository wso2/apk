/* eslint-disable @typescript-eslint/ban-types */
/* eslint-disable no-empty-pattern */
import React from 'react';
import PaginatedClientSide from 'components/data-table/PaginatedClientSide'
import useAxios from "components/hooks/useAxios";
import Loader from "components/Loader";

type Props = {}

export default function ListApplicationPolicies({ }: Props) {
  const { data, loading, error } = useAxios({ url: '/throttling/policies/application' });
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
  return (
    <PaginatedClientSide data={data.list} columns={columns} />
  )
}