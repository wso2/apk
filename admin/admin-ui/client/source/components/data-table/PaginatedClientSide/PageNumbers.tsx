import React from 'react'
import TextField from '@mui/material/TextField';
import Button from '@mui/material/Button';
import Box from '@mui/material/Box';
import { faAngleRight, faAngleLeft, faForwardFast, faBackwardFast } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import InputLabel from '@mui/material/InputLabel';
import MenuItem from '@mui/material/MenuItem';
import FormControl from '@mui/material/FormControl';
import Select, { SelectChangeEvent } from '@mui/material/Select';

interface PageNumbersProps {
    canPreviousPage: boolean;
    canNextPage: boolean;
    pageOptions: any[];
    pageCount: number;
    gotoPage: (page: number) => void;
    nextPage: () => void;
    previousPage: () => void;
    setPageSize: (size: number) => void;
    pageIndex: number;
    pageSize: number;
}
export default function PageNumbers({
    canPreviousPage,
    canNextPage,
    pageOptions,
    pageCount,
    gotoPage,
    nextPage,
    previousPage,
    setPageSize,
    pageIndex,
    pageSize }: PageNumbersProps) {
    return (
        <Box display='flex' alignItems='center' flexDirection='row' justifyContent='end' py={2}>
            <Button sx={{ minWidth: 40 }} onClick={() => gotoPage(0)} disabled={!canPreviousPage}>
                <FontAwesomeIcon icon={faBackwardFast} />
            </Button>
            <Button sx={{ minWidth: 40 }} onClick={() => previousPage()} disabled={!canPreviousPage}>
                <FontAwesomeIcon icon={faAngleLeft} />
            </Button>
            <Button sx={{ minWidth: 40 }} onClick={() => nextPage()} disabled={!canNextPage}>
                <FontAwesomeIcon icon={faAngleRight} />
            </Button>
            <Button sx={{ minWidth: 40 }} onClick={() => gotoPage(pageCount - 1)} disabled={!canNextPage}>
                <FontAwesomeIcon icon={faForwardFast} />
            </Button>
            <span>
                Page{' '}
                <strong>
                    {pageIndex + 1} of {pageOptions.length}
                </strong>{' '}
            </span>
            <Box px={1}>
                <TextField
                    label="Go to page"
                    defaultValue={pageIndex + 1}
                    type="number"
                    onChange={e => {
                        const page = e.target.value ? Number(e.target.value) - 1 : 0
                        gotoPage(page)
                    }}
                    sx={{ width: 100, textAlign: 'center', backgroundColor: 'white', margin: 0, padding: 0 }}
                />
            </Box>
            <FormControl sx={{ m: 1, minWidth: 120 }} size="small">
                <InputLabel id="demo-select-small">Age</InputLabel>
                <Select
                    labelId="demo-select-small"
                    id="demo-select-small"
                    value={pageSize}
                    label="Age"
                    onChange={e => {
                        setPageSize(Number(e.target.value))
                    }}
                >
                    {[10, 20, 30, 40, 50].map(pageSize => (
                        <MenuItem key={pageSize} value={pageSize}>
                            Show {pageSize}
                        </MenuItem>
                    ))}
                    <MenuItem value="">
                        <em>None</em>
                    </MenuItem>
                </Select>
            </FormControl>
        </Box>
    )
}
