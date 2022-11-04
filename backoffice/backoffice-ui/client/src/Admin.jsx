import React, { Suspense, lazy, useState, useEffect } from 'react';
import {
    BrowserRouter as Router, Redirect, Route, Routes,
} from 'react-router-dom';
import AuthManager from 'AppData/AuthManager';
import CONSTS from 'AppData/Constants';
import qs from 'qs';
import Utils from 'AppData/Utils';
import Logout from 'AppComponents/Logout';
import Progress from 'AppComponents/Shared/Progress';
import AdminRootErrorBoundary from 'AppComponents/Shared/AdminRootErrorBoundary';
import Configurations from 'Config';
import { IntlProvider } from 'react-intl';
import RedirectToLogin from 'AppComponents/Shared/RedirectToLogin';
// Localization
import UnexpectedError from 'AppComponents/Base/Errors/UnexpectedError';
// import LoginDenied from './app/LoginDenied';

const ProtectedApp = lazy(() => import('./ProtectedApp' /* webpackChunkName: "ProtectedApps" */));

/**
 * Define base routes for the application
 * @returns {React.Component} base routes for the application
 */
const Admin = () => {
    /**
     * Language.
     * @type {string}
     */
    const language = (navigator.languages && navigator.languages[0]) || navigator.language || navigator.userLanguage;

    /**
     * Language without region code.
     */
    const languageWithoutRegionCode = language.toLowerCase().split(/[_-]+/)[0];
    const { search } = window.location;
    const queryString = search.replace(/^\?/, '');
    /* With QS version up we can directly use {ignoreQueryPrefix: true} option */
    const queryParams = qs.parse(queryString);
    const { environment = Utils.getCurrentEnvironment().label } = queryParams;
    const [userResolved, setUserResolved] = useState(false);
    const [user, setUser] = useState(AuthManager.getUser(environment));
    const [messages, setMessages] = useState({});
    const [notEnoughPermission, setNotEnoughPermission] = useState(false);
    const [unexpectedServerError, setUnexpectedServerError] = useState(false);

    /**
     * Initialize i18n.
     */
    useEffect(() => {
        const locale = languageWithoutRegionCode || language;
        loadLocale(locale);
        const user = AuthManager.getUser();
        if (user) {
            setUser(user);
            setUserResolved(true);
        } else {
            // If no user data available , Get the user info from existing token information
            // This could happen when OAuth code authentication took place and could send
            // user information via redirection
            const userPromise = AuthManager.getUserFromToken();
            userPromise
                .then((loggedUser) => {
                    if (loggedUser != null) {
                        setUser(loggedUser);
                        setUserResolved(true);
                    } else {
                        console.log('User returned with null, redirect to login page');
                        setUserResolved(false);
                    }
                })
                .catch((error) => {
                    if (error && error.message === CONSTS.errorCodes.INSUFFICIENT_PREVILEGES) {
                        setUserResolved(true);
                        setNotEnoughPermission(true);
                    } else if (error && error.message === CONSTS.errorCodes.UNEXPECTED_SERVER_ERROR) {
                        setUserResolved(true);
                        setUnexpectedServerError(true);
                    } else {
                        console.log('Error: ' + error + ',redirecting to login page');
                        setUserResolved(true);
                    }
                });
        }
    }, []);

    /**
     * Load locale file.
     *
     * @param {string} locale Locale name
     */
    const loadLocale = (locale) => {
        // Skip loading the locale file if the language code is english,
        // Because we have used english defaultMessage in the FormattedText component
        // and en.json is generated from those default messages, Hence no point of fetching it
        fetch(`${Configurations.app.context}/site/public/locales/${locale}.json`)
            .then((resp) => resp.json())
            .then((messages) => setMessages(messages));
    }

    /**
     *
     *
     * @returns {React.Component} Render complete app component
     * @memberof Admin
     */

    const locale = languageWithoutRegionCode || language;
    if (!userResolved) {
        return <Progress message='Resolving user ...' />;
    }
    return (
        <div>
            <IntlProvider locale={locale} messages={messages}>
                <AdminRootErrorBoundary appName='Admin Application'>
                    <Router basename={Configurations.app.context}>
                        <Routes>
                            {/* <Redirect exact from='/login' to='/' /> */}
                            <Route path='/logout' component={Logout} />
                            <Route
                                render={() => {
                                    if (notEnoughPermission) {
                                        return <div>Login denied</div>;
                                    } else if (unexpectedServerError) {
                                        return <UnexpectedError />;
                                    } else if (!user) {
                                        return <RedirectToLogin />;
                                    }
                                    return (
                                        <Suspense fallback={<Progress message='Loading app ...' />}>
                                            <ProtectedApp user={user} />
                                        </Suspense>
                                    );
                                }}
                            />
                        </Routes>
                    </Router>
                </AdminRootErrorBoundary>
            </IntlProvider>
        </div>
    );
}

export default Admin;
