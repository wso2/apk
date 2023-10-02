/* eslint-disable react/jsx-key */
import React, { useState } from 'react'
import { useTable, usePagination, useGlobalFilter, useSortBy } from 'react-table'
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import Paper from '@mui/material/Paper';
import PageNumbers from './PageNumbers';
import SearchTable from './SearchTable';
import KeyboardArrowUpIcon from '@mui/icons-material/KeyboardArrowUp';
import KeyboardArrowDownwardIcon from '@mui/icons-material/KeyboardArrowDown';
import Box from '@mui/material/Box';
import Grid from '@mui/material/Grid';

// import '../Tmp.css';

interface TableProps {
  columns: any,
  data: any,
  searchProps: any,
}
interface PaginatedClientSideProps {
  data: any,
  columns: any,
  searchProps: any,
}

function ApplicationPoliciesTable({ columns, data, searchProps }: TableProps) {
  // Use the state and functions returned from useTable to build your UI
  const {
    getTableProps,
    getTableBodyProps,
    headerGroups,
    prepareRow,
    state,
    setGlobalFilter,
    page, // Instead of using 'rows', we'll use page,
    // which has only the rows for the active page

    // The rest of these things are super handy, too ;)
    canPreviousPage,
    canNextPage,
    pageOptions,
    pageCount,
    gotoPage,
    nextPage,
    previousPage,
    setPageSize,
    toggleSortBy,
    state: { pageIndex, pageSize },
  } = useTable(
    {
      columns,
      data,
      initialState: { pageIndex: 0, },
      disableSortRemove: true,
    },
    useGlobalFilter,
    useSortBy,
    usePagination,
  )

  const { globalFilter } = state;

  // Render the UI for your table
  return (
    <>
      <Paper>
        <SearchTable filter={globalFilter} setFilter={setGlobalFilter} searchProps={searchProps} />
        <Table {...getTableProps()} >
          <TableHead>
            {headerGroups.map(headerGroup => (
              <TableRow {...headerGroup.getHeaderGroupProps()}>
                {headerGroup.headers.map(column => (
                  <TableCell {...column.getHeaderProps()}>
                    <Grid container direction="row" alignItems={'center'} columnSpacing={2}>
                      <Grid item>
                        {column.render('Header')}
                      </Grid>
                      {column.sortable && (
                        <Grid item>
                          <Box sx={{ display: 'flex', flexDirection: 'column' }}>
                            <KeyboardArrowUpIcon fontSize="small" onClick={() => toggleSortBy(column.id, false, false)} color={column.isSorted && !column.isSortedDesc ? 'disabled' : 'primary'}/>
                            <KeyboardArrowDownwardIcon fontSize="small" onClick={() => toggleSortBy(column.id, true, false)} color={column.isSorted && column.isSortedDesc ? 'disabled' : 'primary'}/>
                          </Box>
                        </Grid>
                      )}
                    </Grid>
                  </TableCell>
                ))}
              </TableRow>
            ))}
          </TableHead>
          <TableBody {...getTableBodyProps()}>
            {page.map((row, i) => {
              prepareRow(row)
              return (
                <TableRow {...row.getRowProps()}>
                  {row.cells.map(cell => {
                    return <TableCell {...cell.getCellProps()}>{cell.render('Cell')}</TableCell>
                  })}
                </TableRow>
              )
            })}
          </TableBody>
        </Table>
      </Paper>
      <PageNumbers
        canPreviousPage={canPreviousPage}
        canNextPage={canNextPage}
        pageOptions={pageOptions}
        pageCount={pageCount}
        gotoPage={gotoPage}
        nextPage={nextPage}
        previousPage={previousPage}
        setPageSize={setPageSize}
        pageIndex={pageIndex}
        pageSize={pageSize}
      />
    </>
  )
}

function PaginatedClientSide({ data, columns, searchProps }: PaginatedClientSideProps) {
  return (
    <div className='table'>
      <ApplicationPoliciesTable columns={columns} data={data} searchProps={searchProps} />
    </div>
  )
}

export default PaginatedClientSide