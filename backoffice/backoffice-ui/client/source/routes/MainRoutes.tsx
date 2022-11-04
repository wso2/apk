import React, { lazy } from 'react';

// project import
import Loadable from 'components/Loadable';
import MainLayout from 'layout/MainLayout';

// render - dashboard
const DashboardDefault = Loadable(lazy(() => import('pages/Dashboard/Dashboard')));

// render - sample page
const SamplePage = Loadable(lazy(() => import('pages/sample/SamplePage')));
// ==============================|| MAIN ROUTING ||============================== //

const MainRoutes = {
    path: '/',
    element: <MainLayout />,
    children: [
        {
            path: '/',
            element: <DashboardDefault />
        },
        {
            path: 'dashboard',
            element: <DashboardDefault />
        },
        {
            path: 'advanced-policies',
            element: <SamplePage />
        },
        {
            path: 'application-policies',
            element: <SamplePage />
        },
        {
            path: 'subscription-policies',
            element: <SamplePage />
        },
        {
            path: 'custom-policies',
            element: <SamplePage />
        },
        {
            path: 'deny-policies',
            element: <SamplePage />
        },
        {
            path: 'gateways',
            element: <SamplePage />
        },
        {
            path: 'api-categories',
            element: <SamplePage />
        },
        {
            path: 'key-managers',
            element: <SamplePage />
        },
        {
            path: 'user-creation',
            element: <SamplePage />
        },
        {
            path: 'application-creation',
            element: <SamplePage />
        },
        {
            path: 'application-deletion',
            element: <SamplePage />
        },
        {
            path: 'subscription-creation',
            element: <SamplePage />
        },
        {
            path: 'subscription-deletion',
            element: <SamplePage />
        },
        {
            path: 'subscription-update',
            element: <SamplePage />
        },
        {
            path: 'application-registration',
            element: <SamplePage />
        },
        {
            path: 'api-state-change',
            element: <SamplePage />
        },
        {
            path: 'applications',
            element: <SamplePage />
        }
        ,
        {
            path: 'scope-assignments',
            element: <SamplePage />
        },
        {
            path: 'advanced',
            element: <SamplePage />
        }
    ]
};

export default MainRoutes;
