// Import your page components
import React, { lazy } from 'react';
import { Route, Routes } from "react-router-dom";

// project import
import Loadable from 'components/Loadable';
import MainLayout from 'layout/MainLayout';

// render - dashboard
// const DashboardDefault = Loadable(lazy(() => import('pages/Dashboard/Dashboard')));

// render - sample page
const SamplePage = Loadable(lazy(() => import('pages/sample/SamplePage')));
const ApplicationRatePlans = Loadable(lazy(() => import('pages/ApplicationRatePlans/ListApplicationRatePlans')));

// ==============================|| MAIN ROUTING ||============================== //
export default function MainRoutes() {
    return (
        <Routes>
            <Route path="/" element={<MainLayout />}>
                <Route path="/" element={<SamplePage />} />
                <Route path="dashboard" element={<SamplePage />} />
                <Route path="advanced-policies" element={<SamplePage />} />
                <Route path="application-rate-plans" element={<ApplicationRatePlans />} />
                <Route path="business-plans" element={<SamplePage />} />
                <Route path="custom-policies" element={<SamplePage />} />
                {/* <Route path="deny-policies" element={<SamplePage />} /> */}
                {/* <Route path="gateways" element={<SamplePage />} /> */}
                <Route path="api-categories" element={<SamplePage />} />
                {/* <Route path="key-managers" element={<SamplePage />} /> */}
                <Route path="advanced" element={<SamplePage />} />
            </Route>
        </Routes>
    );
}
