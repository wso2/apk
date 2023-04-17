import React, { ReactNode } from 'react';
import PropTypes from 'prop-types';
import { forwardRef } from 'react';

// material-ui
import { useTheme, SxProps } from '@mui/material/styles';
import { Card, CardContent, CardHeader, CardProps, Divider, Typography } from '@mui/material';

// header style
const headerSX = {
    p: 2.5,
    '& .MuiCardHeader-action': { m: '0px auto', alignSelf: 'center' },
    backgroundColor: 'transparent',//(theme) => (theme.palette.mode === 'dark' ? theme.palette.dark.main : theme.palette.grey[50])
};

// ==============================|| CUSTOM - MAIN CARD ||============================== //
export type CardSXProps = {
    sx: SxProps;
}

interface MainCardProps extends Omit<CardProps, 'sx'> {
    border?: boolean,
    boxShadow?: boolean,
    contentSX?: SxProps,
    darkTitle?: boolean,
    divider?: boolean,
    secondary?: ReactNode,
    shadow?: string,
    sx?: CardSXProps,
    title?: string,
    codeHighlight?: boolean,
    content?: boolean,
    children?: ReactNode
}

type Ref = HTMLDivElement | null;

const MainCard = forwardRef<Ref, MainCardProps>(
    (
        {
            border = true,
            boxShadow,
            children,
            content = true,
            contentSX = { backgroundColor: '#fff'},
            darkTitle,
            divider = true,
            elevation,
            secondary,
            shadow,
            sx = {},
            title,
            codeHighlight,
            ...others
        },
        ref
    ): JSX.Element => {
        const theme = useTheme();
        boxShadow = theme.palette.mode === 'dark' ? boxShadow || true : boxShadow;

        return (
            <Card
                elevation={elevation || 0}
                ref={ref}
                {...others}
                sx={{
                    ...sx,
                    border: border ? '1px solid' : 'none',
                    borderRadius: 2,
                    borderColor: theme.palette.mode === 'dark' ? theme.palette.divider : theme.palette.grey.A200,
                    boxShadow: boxShadow && (!border || theme.palette.mode === 'dark') ? shadow || theme.customShadows.z1 : 'inherit',
                    ':hover': {
                        boxShadow: boxShadow ? shadow || theme.customShadows.z1 : 'inherit'
                    },
                    '& pre': {
                        m: 0,
                        p: '16px !important',
                        fontFamily: theme.typography.fontFamily,
                        fontSize: '0.75rem'
                    },
                    backgroundColor: 'transparent',
                }}
                {...others}

            >
                {/* card header and action */}
                {!darkTitle && title && (
                    <CardHeader sx={headerSX} titleTypographyProps={{ variant: 'subtitle1' }} title={title} action={secondary} />
                )}
                {darkTitle && title && (
                    <CardHeader sx={headerSX} title={<Typography variant="h3">{title}</Typography>} action={secondary} />
                )}

                {/* content & header divider */}
                {title && divider && <Divider />}

                {/* card content */}
                {content && <CardContent sx={contentSX}>{children}</CardContent>}
                {!content && children}

            </Card>
        );
    }
);
MainCard.displayName = 'MainCard';
export default MainCard;
