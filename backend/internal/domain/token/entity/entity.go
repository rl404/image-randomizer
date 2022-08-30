package entity

// CreateAccessTokenRequest is request model for create access token.
type CreateAccessTokenRequest struct {
	UserID     int64
	AccessUUID string
}

// CreateRefreshTokenRequest is request model for create refresh token.
type CreateRefreshTokenRequest struct {
	UserID      int64
	RefreshUUID string
}

// Token is entity for token.
type Token struct {
	AccessToken  string
	RefreshToken string
}
