package testutils

type TestCaseUsers struct {
	Description   string
	ApiCreator    Credentials
	ApiPublisher  Credentials
	ApiSubscriber Credentials
	Admin         Credentials
	CtlUser       Credentials
}

type Credentials struct {
	Username string
	Password string
}
