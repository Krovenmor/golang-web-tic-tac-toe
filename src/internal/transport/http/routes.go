package http

import (
	"WebTic-tac-toe2/internal/app"
	"WebTic-tac-toe2/internal/service/auth"
	"WebTic-tac-toe2/internal/service/jwtauth"
	"WebTic-tac-toe2/internal/transport/http/authhandler"
	"WebTic-tac-toe2/internal/transport/http/gamehandler"
	"net/http"
)

type MainHandler struct {
	gameHandler       gamehandler.GameHandler
	authHandler       authhandler.AuthHandler
	userAuthenticator authhandler.UserAuthenticator
}

func NewMainHandler(gs app.GameService, us auth.UserService, jwt jwtauth.JWTService) *MainHandler {
	mh := &MainHandler{}
	mh.gameHandler = *gamehandler.NewGameHandler(gs)
	mh.authHandler = *authhandler.NewAuthHandler(us, jwt)
	mh.userAuthenticator = *authhandler.NewUserAuthenticator(jwt)
	return mh
}

func (mh *MainHandler) RegisterRoutes(m *http.ServeMux) http.Handler {
	m.HandleFunc("POST /api/game/{uuid}", mh.gameHandler.HandleMove)
	m.HandleFunc("POST /api/game/{uuid}/join", mh.gameHandler.HandleJoinPairGame)
	m.HandleFunc("POST /api/game/create/solo", mh.gameHandler.HandleNewGame)
	m.HandleFunc("POST /api/game/create/pair", mh.gameHandler.HandleNewGamePair)

	m.HandleFunc("GET /api/game/{uuid}/shortinfo", mh.gameHandler.SendGameShortInfo)
	m.HandleFunc("GET /api/game/{uuid}/fullinfo", mh.gameHandler.SendGameFullInfo)
	m.HandleFunc("GET /api/game/player/{uuid}/info", mh.gameHandler.SendPlayerInfo)
	m.HandleFunc("GET /api/game/player/{uuid}/games", mh.gameHandler.HandleListCompletedPairGames)
	m.HandleFunc("GET /api/game/pair", mh.gameHandler.HandleListPairGames)
	m.HandleFunc("GET /api/game/leaderboard", mh.gameHandler.HandleLeaderBoard)

	m.HandleFunc("POST /api/auth/register", mh.authHandler.Register)
	m.HandleFunc("POST /api/auth/login", mh.authHandler.LogIn)

	m.HandleFunc("POST /api/jwt/update/access", mh.authHandler.UpdateAccess)
	m.HandleFunc("POST /api/jwt/update/refresh", mh.authHandler.UpdateRefresh)

	m.HandleFunc("GET /", mh.gameHandler.HandleIndex)

	return mh.userAuthenticator.Middleware(m)
}
