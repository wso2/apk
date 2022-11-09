import React, { ReactNode } from 'react';
import PropTypes from 'prop-types';
import { forwardRef } from 'react';

// material-ui
import { useTheme, SxProps } from '@mui/material/styles';
import { Card, CardContent, CardHeader, CardProps, Divider, Typography } from '@mui/material';

// header style
const headerSX = {
    p: 2.5,
    '& .MuiCardHeader-action': { m: '0px auto', alignSelf: 'center' }
};

// ==============================|| CUSTOM - MAIN CARD ||============================== //
export type CardSXProps = {
    sx: SxProps;
}

interface MainCardProps {
    border?: boolean,
    boxShadow?: boolean,
    contentSX?: object,
    darkTitle?: boolean,
    divider?: boolean,
    elevation?: number,
    secondary?: ReactNode,
    shadow?: string,
    sx?: CardSXProps,
    title?: string,
    codeHighlight?: boolean,
    content?: boolean,
    children?: ReactNode
}

export type Ref = HTMLDivElement;

const MainCard = forwardRef<Ref, MainCardProps>(
    (
        {
            border = true,
            boxShadow,
            children,
            content = true,
            contentSX = {},
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
    ) => {
        const theme = useTheme();
        boxShadow = theme.palette.mode === 'dark' ? boxShadow || true : boxShadow;

        return (
            <Card
                elevation={elevation || 0}
                ref={ref}
                sx={{
                    ...sx,
                    border: border ? '1px solid' : 'none',
                    borderRadius: 2,
                    borderColor: theme.palette.mode === 'dark' ? theme.palette.divider : theme.palette.grey[100],
                    '& pre': {
                        m: 0,
                        p: '16px !important',
                        fontFamily: theme.typography.fontFamily,
                        fontSize: '0.75rem'
                    }
                }}
                // sx={{
                //     ...sx,
                //     border: border ? '1px solid' : 'none',
                //     borderRadius: 2,
                //     borderColor: theme.palette.mode === 'dark' ? theme.palette.divider : theme.palette.grey[100],
                //     boxShadow: boxShadow && (!border || theme.palette.mode === 'dark') ? shadow || theme.customShadows.z1 : 'inherit',
                //     ':hover': {
                //         boxShadow: boxShadow //? shadow || theme.customShadows.z1 : 'inherit'
                //     },
                //     '& pre': {
                //         m: 0,
                //         p: '16px !important',
                //         fontFamily: theme.typography.fontFamily,
                //         fontSize: '0.75rem'
                //     }
                // }}
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

export default MainCard;
