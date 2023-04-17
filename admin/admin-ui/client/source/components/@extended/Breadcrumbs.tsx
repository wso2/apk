import React from 'react';
import PropTypes from 'prop-types';
import { useEffect, useState } from 'react';
import { Link, useLocation } from 'react-router-dom';

// material-ui
import MuiBreadcrumbs from '@mui/material/Breadcrumbs';
import { Grid, Typography } from '@mui/material';

// project imports
import MainCard from '../MainCard';

// ==============================|| BREADCRUMBS ||============================== //

interface BreadcrumbsProps {
    navigation: any;
    title: boolean;
    divider?: boolean;
    others?: any;
}

const Breadcrumbs = ({ navigation, title, ...others }: BreadcrumbsProps) => {
    const location = useLocation();
    const [main, setMain] = useState<{ type?: string; title?: string; url?: string, breadcrumbs?: any}>();
    const [item, setItem] = useState<{ type?: string; title?: string; url?: string, breadcrumbs?: any}>();

    // set active item state
    const getCollapse = (menu: any) => {
        if (menu.children) {
            menu.children.filter((collapse) => {
                if (collapse.type && collapse.type === 'collapse') {
                    getCollapse(collapse);
                } else if (collapse.type && collapse.type === 'item') {
                    if (location.pathname === collapse.url) {
                        setMain(menu);
                        setItem(collapse);
                    }
                }
                return false;
            });
        }
    };

    useEffect(() => {
        navigation?.items?.map((menu) => {
            if (menu.type && menu.type === 'group') {
                getCollapse(menu);
            }
            return false;
        });
    });

    // only used for component demo breadcrumbs
    if (location.pathname === '/breadcrumbs') {
        location.pathname = '/dashboard/analytics';
    }

    let mainContent;
    let itemContent;
    let breadcrumbContent = <Typography />;

    // collapse item
    if (main && main.type === 'collapse') {
        mainContent = (
            <Typography component={Link} to={document.location.pathname} variant="h6" sx={{ textDecoration: 'none' }} color="textSecondary">
                {main.title}
            </Typography>
        );
    }

    // items
    if (item && item.type === 'item') {
        itemContent = (
            <Typography variant="subtitle1" color="textPrimary">
                {item.title}
            </Typography>
        );

        // main
        if (item.breadcrumbs !== false) {
            breadcrumbContent = (
                <MainCard border={false} {...others} content={false}>
                    <Grid container direction="column" justifyContent="flex-start" alignItems="flex-start" spacing={1}>
                        <Grid item>
                            <MuiBreadcrumbs aria-label="breadcrumb">
                                <Typography component={Link} to="/" color="textSecondary" variant="h6" sx={{ textDecoration: 'none' }}>
                                    Home
                                </Typography>
                                {mainContent}
                                {itemContent}
                            </MuiBreadcrumbs>
                        </Grid>
                        {title && (
                            <Grid item sx={{ mt: 2, mb: 2 }}>
                                <Typography variant="h3">{item.title}</Typography>
                            </Grid>
                        )}
                    </Grid>
                </MainCard>
            );
        }
    }

    return breadcrumbContent;
};

Breadcrumbs.propTypes = {
    navigation: PropTypes.object,
    title: PropTypes.bool
};

export default Breadcrumbs;
