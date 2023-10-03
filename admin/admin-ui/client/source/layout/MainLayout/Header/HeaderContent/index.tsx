
import React from 'react';
// material-ui
import { Box, Theme, useMediaQuery } from '@mui/material';

// project import
import Profile from './Profile';
import MobileSection from './MobileSection';
import { Grid } from '@mui/material';

// ==============================|| HEADER - CONTENT ||============================== //

const HeaderContent = () => {
    const matchesXs = useMediaQuery((theme: Theme) => theme.breakpoints.down('md'));

    return (
        <>
            <Grid container direction='row' justifyContent='space-between'>
                <Grid item sx={{ mt: 2, mb: 2 }}>
                    {matchesXs && <Box sx={{ width: '100%', ml: 1 }} />}
                </Grid>
                <Grid item display='grid'>
                    {!matchesXs && <Profile />}
                    {matchesXs && <MobileSection />}
                </Grid>
            </Grid>
        </>
    );
};

export default HeaderContent;
