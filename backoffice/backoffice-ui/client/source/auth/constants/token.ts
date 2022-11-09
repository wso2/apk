interface ServiceResourcesType {
    jwks: string;
    token: string;
}

export const SERVICE_RESOURCES: ServiceResourcesType = {
    jwks: "/oauth2/jwks",
    token: "/oauth2/token"
};

export const AUTHORIZATION_CODE = "code";
export const ID_TOKEN = "id_token";
export const REFRESH_TOKEN = "refresh_token";
export const ACCESS_TOKEN = "access_token";
export const PKCE_CODE_VERIFIER = "pkce_code_verifier";
export const AUTHORIZATION_ENDPOINT = "authorization_endpoint";
export const TOKEN_ENDPOINT = "token_endpoint";
export const END_SESSION_ENDPOINT = "end_session_endpoint";
export const JWKS_ENDPOINT = "jwks_uri";
export const ISSUER = "issuer";
export const OP_CONFIG_INITIATED = "op_config_initiated";
export const REQUEST_PARAMS = "request_params";
export const REQUEST_STATUS = "request_status";
export const ACCESS_TOKEN_EXPIRE_IN = "expires_in";
export const ACCESS_TOKEN_ISSUED_AT = "issued_at";
export const SCOPE = "scope";
export const TOKEN_TYPE = "token_type";