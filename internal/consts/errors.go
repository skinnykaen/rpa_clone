package consts

const (
	ErrEmailAlreadyInUse         = "email already in use"                         // http code 400
	ErrAtoi                      = "string to int error"                          // 400
	ErrActivationLinkUnavailable = "activation link is currently unavailable"     // 503
	ErrIncorrectPasswordOrEmail  = "incorrect password or email"                  // 400
	ErrUserIsNotActive           = "user is not active. please check your email"  // 403
	ErrShortPassword             = "please input password, at least 6 symbols"    // 400
	ErrTokenExpired              = "token expired"                                // 401
	ErrNotStandardToken          = "token claims are not of type *StandardClaims" // 401
	ErrProjectPageIsBanned       = "the projectPage is banned. no access"         // 403
	ErrAccessDenied              = "access denied"                                // 403
	ErrNotFoundInDB              = "not found"                                    // 400
)
