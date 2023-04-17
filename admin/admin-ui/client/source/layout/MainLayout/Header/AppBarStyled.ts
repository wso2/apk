import React from 'react';
// material-ui
import { styled, Theme } from '@mui/material/styles';
import AppBar, { AppBarProps } from '@mui/material/AppBar';

// eslint-disable-next-line @typescript-eslint/no-var-requires, no-undef
const Settings = require('Settings');

interface AppBarStyledProps extends AppBarProps {
    open: boolean,
}

// ==============================|| HEADER - APP BAR STYLED ||============================== //

const AppBarStyled = styled(AppBar, { shouldForwardProp: (prop) => prop !== 'open' })<AppBarStyledProps>(({ theme, open }) => ({
    zIndex: theme.zIndex.drawer + 1,
    transition: theme.transitions.create(['width', 'margin'], {
        easing: theme.transitions.easing.sharp,
        duration: theme.transitions.duration.leavingScreen
    }),
    ...(open && {
        marginLeft: 260,
        width: `calc(100% - 260px)`,
        transition: theme.transitions.create(['width', 'margin'], {
            easing: theme.transitions.easing.sharp,
            duration: theme.transitions.duration.enteringScreen
        })
    })
}));

export default AppBarStyled;
