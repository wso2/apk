package authenticator

const (
	// ExpiredToken is the error code for an invalid token.
	ExpiredToken = 900901
	// ExpiredTokenMessage is the error code for an invalid token.
	ExpiredTokenMessage = "Token is Expired"
	// MissingCredentials is the error code for missing credentials.
	MissingCredentials = 900902
	// MissingCredentialsMesage is the error message for missing credentials.
	MissingCredentialsMesage = "Missing Credentials"
	// InvalidCredentials is the error code for invalid credentials.
	InvalidCredentials = 900903
	// InvalidCredentialsMessage is the error message for invalid credentials.
	InvalidCredentialsMessage = "Invalid Credentials"
	// APIAuthGeneralError is the error code for an unclassified authentication failure.
	APIAuthGeneralError = 900900
	// APIAuthGeneralErrorMessage is the error message for an unclassified authentication failure.
	APIAuthGeneralErrorMessage = "Unclassified Authentication Failure"
)
