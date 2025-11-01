package main

type SendSignInLinkRequest struct {
	RequestType     string `json:"requestType"` // must be "EMAIL_SIGNIN"
	Email           string `json:"email"`
	ContinueURL     string `json:"continueUrl"`
	HandleCodeInApp bool   `json:"handleCodeInApp,omitempty"`
}

type SendSignInLinkResponse struct {
	Kind  string `json:"kind,omitempty"`
	Email string `json:"email"`
}

type SignInWithEmailLinkRequest struct {
	Email   string `json:"email"`
	OobCode string `json:"oobCode"`
}

type SignInWithEmailLinkResponse struct {
	IDToken      string `json:"idToken"`
	Email        string `json:"email"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    string `json:"expiresIn"`
	LocalID      string `json:"localId"`
}
