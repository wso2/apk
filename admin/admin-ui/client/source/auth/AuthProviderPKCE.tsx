/* eslint-disable @typescript-eslint/no-var-requires */
/* eslint-disable @typescript-eslint/no-empty-function */
import React, { createContext, useContext, useEffect, useState } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { getUserInfoEndpoint, initOPConfiguration } from 'auth/op-config';
import { sendAuthorizationRequest, sendTokenRequest } from 'auth/sign-in';
import { OIDCRequestParamsInterface } from 'auth/types/oidc-request-params';
import { getSessionParameter, setSessionParameter, initUserSession, getAccessTokenFromRefreshToken } from "auth/session";
const Settings = require('Settings');
import { ACCESS_TOKEN, REQUESTED_PATH } from 'auth/constants/token';

// Define a context for storing the user authentication state and related functions
type AuthContextType = {
  isAuthenticated: boolean;
  user: any;
  loading: boolean;
  login: () => void;
  logout: () => void;
};

const AuthContext = createContext<AuthContextType>({
  isAuthenticated: false,
  user: null,
  loading: true,
  login: () => { },
  logout: () => { },
});

const requestParams: OIDCRequestParamsInterface = {
  clientId: Settings.idp.client_id,
  scope: Settings.idp.scope,
  state: Settings.idp.state,
  serverOrigin: Settings.idp.server_origin
};

interface AuthProviderProps {
  children: React.ReactNode;
}

// Define a higher-order component for wrapping the app with the AuthProvider context
export const AuthProvider = ({ children }: AuthProviderProps) => {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [user, setUser] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();
  const { pathname } = useLocation();

  // Define a function to handle user login
  const handleLogin = () => {
    initOPConfiguration(Settings.idp.well_known).then(() => {
      sendAuthorizationRequest(requestParams);
      setSessionParameter(REQUESTED_PATH, pathname);
    })
    // Navigate to the login page
  };

  // Define a function to handle user logout
  const handleLogout = () => {
    // Clear the user authentication state and navigate to the home page
    setIsAuthenticated(false);
    setUser(null);
    navigate('/');
  };

  /**
   * Call user info endpoint to get user details
   */
  const getUserInfo = async (token) => {
    // If there is a token, set the user authentication state and navigate to the home page
    setIsAuthenticated(true);
    const userInfoResponse = await fetch(getUserInfoEndpoint(), {
      headers: { Authorization: `Bearer ${token}` },
    });
    setUser(userInfoResponse);
    const requestedPath = getSessionParameter(REQUESTED_PATH);
    if (requestedPath) {
      navigate(requestedPath);
      sessionStorage.removeItem(REQUESTED_PATH);
    } else if (pathname && pathname !== 'undefined' && pathname !== '') {
      navigate(pathname)
    } else {
      navigate('/');
    }
    setLoading(false);
  }


  // Define a useEffect hook to check for an authenticated user on mount
  useEffect(() => {
    const fetchUser = async () => {
      // Check for an authenticated user by parsing the URL for an access token
      const query = new URLSearchParams(window.location.search);
      const code = query.get('code');
      let response;
      if (code) {
        try {
          // Exchange the authorization code for an access token
          response = await sendTokenRequest(requestParams);
        } catch (error) {
          if (error.response.status === 400) {
            sendAuthorizationRequest(requestParams);
            setSessionParameter(REQUESTED_PATH, pathname);
          }
        }
        initUserSession(response);
        const accessToken = response.accessToken;
        getUserInfo(accessToken);
      } else {
        // If there is no authorization code, check for an access token in the session storage
        let token = getSessionParameter(ACCESS_TOKEN);
        if (!token) {
          // If there is no token stored in the session storage, it means the session is cleared,
          // we also assume the refresh token is expired, so we need to re-login
          handleLogin();
        } else {
          // Note: getAccessToken() will return existing token if the access token is not expired
          token = getAccessTokenFromRefreshToken() as any;
          if (token) {
            // If there is a token, set the user authentication state and navigate to the home page
            getUserInfo(token);
          } else {
            // If there is no token, call the auth endpoint
            handleLogin();
          }
        }
      }
    };
    fetchUser();
  }, [navigate]);

  if (loading) {
    return <div>Loading...</div>
  }

  // Return the AuthProvider context provider with the user authentication state and related functions
  return (
    <AuthContext.Provider value={{ isAuthenticated, user, loading, login: handleLogin, logout: handleLogout }}>
      {children}
    </AuthContext.Provider>
  );
};

// Define a custom hook for accessing the AuthProvider context
export const useAuth = () => useContext(AuthContext);
