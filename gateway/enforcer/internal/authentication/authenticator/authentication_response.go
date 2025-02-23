package authenticator

// AuthenticationResponse contains the authentication status of the request when it passed through an authenticator.
type AuthenticationResponse struct {
	Authenticated               bool
	MandatoryAuthentication     bool
	ContinueToNextAuthenticator bool
	ErrorCode                   int64
	ErrorMessage                string
}
