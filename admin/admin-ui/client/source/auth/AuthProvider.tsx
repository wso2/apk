/* eslint-disable @typescript-eslint/no-var-requires */
/* eslint-disable @typescript-eslint/no-empty-function */
import React, { createContext, useContext, useEffect, useState } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { OIDCRequestParamsInterface } from 'auth/types/oidc-request-params';
const Settings = require('Settings');
import { REQUESTED_PATH, USER } from 'auth/constants/token';
import { SessionUser } from 'types/SessionUser';

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
  const [user, setUser] = useState<SessionUser | null>(null);
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();
  const { pathname } = useLocation();

  // Define a function to handle user login
  const handleLogin = () => {
    sessionStorage.setItem(REQUESTED_PATH, pathname);
    document.location.href = `${Settings.idp.authorization_endpoint}?` +
      `response_type=code` +
      `&client_id=${Settings.idp.client_id}` +
      `&redirect_uri=${Settings.idp.redirect_uri}/token`;
    // Navigate to the login page
  };

  // Define a function to handle user logout
  const handleLogout = () => {
    // Clear the user authentication state and navigate to the home page
    setIsAuthenticated(false);
    setUser(null);
    sessionStorage.removeItem(USER);
    navigate('/');
  };

  /**
   * Define a function to handle routing after authentication
   */
  const route = async () => {
    // Set the user authentication state and navigate to the requested page
    setIsAuthenticated(true);
    const requestedPath = sessionStorage.getItem(REQUESTED_PATH);
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

  const updateUserFromRefreshToken = async () => {
    // Call the nodejs backend to get the access token from the refresh token ( this needs to implement )          
    try {
      // const newUser = getAccessTokenFromRefreshToken() as any;
      // sessionStorage.setItem(USER, `{user: ${newUser}}`);
      throw new Error('Refresh token is expired');
    } catch (e) {
      // Either the refresh token is expired or the refresh token is not valid
      console.log(e);
      handleLogin();
    }
  }
  // Define a useEffect hook to check for an authenticated user on mount
  useEffect(() => {
    const fetchUser = async () => {
      // Check for an authenticated user by parsing the URL for an access token
      const query = new URLSearchParams(window.location.search);
      const user = query.get('user');
      const exp = query.get('exp');
      if (user) {
        const userObject = { user, exp } as SessionUser;
        sessionStorage.setItem(USER, JSON.stringify(userObject));
        setUser(userObject);
        route();
      } else {
        // If there is no authorization code, check for an access token in the session storage
        const userFromSession = sessionStorage.getItem(USER);

        if (userFromSession === null) {
          // If there is no token stored in the session storage, it means the session is cleared,
          // we also assume the refresh token is expired, so we need to re-login
          handleLogin();
        } else if (userFromSession !== null) {
          // If there is a token stored in the session storage, it means the session is not cleared,
          // we also assume the refresh token is not expired, so we can get a new access token from the refresh token
          const userFromSessionObject = JSON.parse(userFromSession);
          if (userFromSessionObject.exp * 1000 < Date.now()) {
            // The access token is expired, we need to get a new access token from the refresh token
            updateUserFromRefreshToken()
          } else {
            // The access token is not expired, we can use the access token to get the user info
            setUser(userFromSessionObject);
            route();
          }
          setUser(userFromSessionObject);
          route();
        } else {
          updateUserFromRefreshToken();
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
