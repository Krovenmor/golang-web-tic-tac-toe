package authhandler

func ToJWTResponce(aToken, rToken string) *JwtResponse {
	return &JwtResponse{
		AccessToken:  aToken,
		RefreshToken: rToken,
	}
}
