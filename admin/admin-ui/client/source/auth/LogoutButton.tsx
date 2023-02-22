import React from "react";
import { sendSignOutRequest } from "./sign-out";
import { endAuthenticatedSession } from "./session";
import { resetOPConfiguration } from "./op-config";
import Settings from '../../public/conf/Settings';

const LogoutButton = () => {

  return (
    <button onClick={() => sendSignOutRequest(Settings.loginUri, () => {
      endAuthenticatedSession();
      resetOPConfiguration();
    })}>
      Log Out
    </button>
  );
};

export default LogoutButton;