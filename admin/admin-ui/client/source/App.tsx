import React, { useEffect, useState } from "react";
import { hasAuthorizationCode, sendAuthorizationRequest, sendTokenRequest, hasValidToken } from 'auth/sign-in';
import { initOPConfiguration } from 'auth/op-config';
import { OIDCRequestParamsInterface } from 'auth/models/oidc-request-params';
import { getSessionParameter, setSessionParameter, initUserSession } from "auth/session";
import Settings from '../public/conf/Settings';
import { REQUEST_STATUS } from 'auth/constants/token';
import ThemeCustomization from 'themes';
import ScrollTop from 'components/ScrollTop';
import { AppContextProvider } from "context/AppContext";
import { getUser } from "auth/AuthManager";
import LoggedInUser from 'types/LoggedInUser';
import Error from "types/Error";
import Routes from 'routes';
import useAxios from "components/hooks/useAxios";
import { IntlProvider } from 'react-intl'

export default function App() {
  const [errors, setErrors] = useState<Array<Error>>([]);
  const [selectedRoute, setSelectedRoute] = useState<string>('dashboard');

  const alterError = (error: Error, add: boolean | null) => {
    const newErrors = [...errors].filter(err => err.key !== error.key);
    if (add) {
      newErrors.push(error);
      setErrors(newErrors);
    } else {
      setErrors(newErrors.filter(err => err !== error));
    }
  }
  const requestParams: OIDCRequestParamsInterface = {
    clientId: Settings.IDP_CLIENT_ID,
    scope: Settings.scope,
    state: Settings.state,
    serverOrigin: Settings.serverOrigin
  };

  let hasToken = hasValidToken();
  let initial_req_status = getSessionParameter(REQUEST_STATUS);

  useEffect(() => {
    if (initial_req_status !== 'sent') {
      initOPConfiguration(Settings.wellKnown, false).then(() => {
        sendAuthorizationRequest(requestParams);
        setSessionParameter(REQUEST_STATUS, 'sent');
      })
    }
  }, [initial_req_status]);

  useEffect(() => {
    if (!hasToken && hasAuthorizationCode()) {
      sendTokenRequest(requestParams)
        .then((response) => {
          initUserSession(
            response
          );
        })
        .catch((error) => {
          if (error.response.status === 400) {
            sendAuthorizationRequest(requestParams);
          }
          /* ================= TODO ================= */
          // Put the proper error from response
          alterError({ key: 'userError', value: 'Error while fetching token' }, true);
        })
        .then(() => {
          //window.location.href = `${Settings.loginUri}/users`;
        })
    }
  }, [hasToken]);
  if (errors.length > 0) {
    return <div>show errors</div>
  }

  if (hasToken) {
    /* ================= TODO ================= */
    // We are getting the user with scopes ( hard coded.) This logic need to be change after the proper STS config is done.
    const user = getUser();
    const { data: settings, loading: loadingSettings, error: errorSettings } = useAxios({ url: '/settings' });
    const { data: tenantConfig, loading: loadingTenantConfig, error: errorTenantConfig } = useAxios({ url: '/tenant-info/YWRtaW5AY2FyYm9uLnN1cGVy' });

    if (loadingSettings || loadingTenantConfig) {
      return <div>Loading ..</div>
    }

    if (errorSettings || errorTenantConfig) {
      return <div>Error</div>
    }

    const { tenantDomain } = tenantConfig;
    const isSuperTenant = (tenantDomain === 'carbon.super');
    return (
      <AppContextProvider value={{ settings, user, isSuperTenant, selectedRoute, setSelectedRoute }}>
        <ThemeCustomization>
          <ScrollTop>
            <IntlProvider locale={'en'} messages={{}}>
              <Routes />
            </IntlProvider>
          </ScrollTop>
        </ThemeCustomization>
      </AppContextProvider>
    );
  } else {
    return <>Redirecting to login page</>
  }
}
