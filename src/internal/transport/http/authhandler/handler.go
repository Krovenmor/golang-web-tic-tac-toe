package authhandler

import (
	"WebTic-tac-toe2/internal/service/auth"
	"WebTic-tac-toe2/internal/service/jwtauth"
	"WebTic-tac-toe2/internal/transport/http/utils"
	"fmt"
	"net/http"
)

type AuthHandler struct {
	service auth.UserService
	jwt     jwtauth.JWTService
}

func NewAuthHandler(s auth.UserService, j jwtauth.JWTService) *AuthHandler {
	return &AuthHandler{service: s, jwt: j}
}

var successMap = map[string]string{"status": "success"}

func (a *AuthHandler) Register(rw http.ResponseWriter, req *http.Request) {
	var sr JwtRequest
	err := utils.Recv(req, &sr)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	if len(sr.Login) <= 1 || len(sr.Password) <= 1 {
		http.Error(rw, "Login or password too short", http.StatusBadRequest)
		return
	}

	ctx := req.Context()
	err = a.service.Register(ctx, sr.Login, sr.Password)
	if err != nil {
		http.Error(rw, fmt.Sprintf("Trouble with register new user: %q", err.Error()), http.StatusBadRequest)
		return
	}

	utils.Send(rw, &successMap)
}

func (a *AuthHandler) LogIn(rw http.ResponseWriter, req *http.Request) {

	var jreq JwtRequest
	err := utils.Recv(req, &jreq)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := req.Context()
	uUUID, err := a.service.LogIn(ctx, jreq.Login, jreq.Password)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	aToken, err := a.jwt.GenAccessToken(uUUID)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rToken, err := a.jwt.GenRefreshToken(req.Context(), uUUID, "")
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := ToJWTResponce(aToken, rToken)
	utils.Send(rw, resp)
}

func (a *AuthHandler) UpdateAccess(rw http.ResponseWriter, req *http.Request) {

	var reqBody RefreshJwtRequest
	err := utils.Recv(req, &reqBody)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err = a.jwt.IsValidRefresh(req.Context(), reqBody.RefreshToken)
	if err != nil {
		http.Error(rw, fmt.Sprintf("Refresh token not valid: %q", err.Error()), http.StatusBadRequest)
		return
	}

	uUUID, err := a.jwt.GetUUIDfromRefreshToken(reqBody.RefreshToken)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	aToken, err := a.jwt.GenAccessToken(uUUID)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := ToJWTResponce(aToken, reqBody.RefreshToken)
	utils.Send(rw, resp)
}

func (a *AuthHandler) UpdateRefresh(rw http.ResponseWriter, req *http.Request) {

	var reqBody RefreshJwtRequest
	err := utils.Recv(req, &reqBody)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err = a.jwt.IsValidRefresh(req.Context(), reqBody.RefreshToken)
	if err != nil {
		http.Error(rw, fmt.Sprintf("Refresh token not valid: %q", err.Error()), http.StatusBadRequest)
		return
	}

	uUUID, err := a.jwt.GetUUIDfromRefreshToken(reqBody.RefreshToken)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	aToken, err := a.jwt.GenAccessToken(uUUID)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rToken, err := a.jwt.GenRefreshToken(req.Context(), uUUID, reqBody.RefreshToken)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := ToJWTResponce(aToken, rToken)
	utils.Send(rw, resp)
}
