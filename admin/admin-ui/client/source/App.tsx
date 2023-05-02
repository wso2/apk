/* eslint-disable react/no-unescaped-entities */
import React from "react";
import MainRoutes from 'routes/MainRoutes';
import { AuthProvider } from 'auth/AuthProvider';
import ThemeCustomization from 'themes';
import ScrollTop from 'components/ScrollTop';
import { IntlProvider } from 'react-intl'

export default function AppTmp() {
  return (
    <ThemeCustomization>
      <ScrollTop>
        <IntlProvider locale={'en'} messages={{}}>
          <AuthProvider>
            <MainRoutes />
          </AuthProvider>
        </IntlProvider>
      </ScrollTop>
    </ThemeCustomization>
  );
}

