package constants

const (
	AuthStatus             = "Auth Status"
	AuthAuthenticated      = "Authenticated"
	AuthBadRequest         = "Bad Request"         // Request does not follow the required struct
	AuthNoCookie           = "No Cookie"           // No cookie in the request header
	AuthNoSession          = "No Session"          // No session found in the session store
	AuthInvalidCredentials = "Invalid Credentials" // wrong username/ password (either username does not exist OR password is wrong)
	AuthServerErr          = "Server Error"        // Server unabel to carry out the operation (eg write to db or session store)
)
