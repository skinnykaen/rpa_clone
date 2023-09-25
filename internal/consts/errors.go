package consts

// http code 400
const (
	ErrEmailAlreadyInUse        = "email already in use"
	ErrAtoi                     = "string to int error"
	ErrIncorrectPasswordOrEmail = "incorrect password or email"
	ErrNotFoundInDB             = "not found"
	ErrShortPassword            = "please input password, at least 6 symbols"
)

// http code 401
const (
	ErrTokenExpired     = "token expired"
	ErrNotStandardToken = "token claims are not of type *StandardClaims"
)

// http code 403
const (
	ErrUserIsNotActive     = "user is not active. please check your email"
	ErrProjectPageIsBanned = "the projectPage is banned. no access"
	ErrAccessDenied        = "access denied"
	ErrMessagingToYourself = "sending a message to yourself"
	ErrChatWithYourself    = "creating chat with yourself"
	ErrNotFoundAuthToken   = "authToken not found in transport payload"
	ErrEmptyDataWithClaims = "empty data with claims"
)

// http code 500
const (
	ErrThereIsNoObservers = "there is no connection to the observer"
)

// ErrActivationLinkUnavailable have http code 503
const (
	ErrActivationLinkUnavailable = "activation link is currently unavailable"
)

// ErrIncorrectInputParam edx errors
const (
	ErrIncorrectInputParam = "error incorrect input params"
)
