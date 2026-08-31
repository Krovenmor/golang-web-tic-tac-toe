package gamehandler

import (
	"WebTic-tac-toe2/internal/app"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

const (
	defLeaderBoardLimit = 10
)

type GameHandler struct {
	service app.GameService
}

func NewGameHandler(s app.GameService) *GameHandler {
	return &GameHandler{service: s}
}

func (gh *GameHandler) HandleIndex(rw http.ResponseWriter, req *http.Request) {
	http.ServeFile(rw, req, "web/index.html")
}

func (gh *GameHandler) HandleNewGame(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	gUUID, err := getUuidFromAuth(req)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	game, err := gh.service.NewGame(ctx, gUUID)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	sendFullGameSolo(rw, game)
}

func (gh *GameHandler) HandleNewGamePair(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	fUUID, err := getUuidFromAuth(req)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	game, err := gh.service.NewPairGame(ctx, fUUID)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	sendFullGamePair(rw, game)
}

func (gh *GameHandler) HandleListPairGames(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	games, err := gh.service.GetAllAvailablePairGames(ctx)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	sendPairList(rw, games)
}

func (gh *GameHandler) HandleListCompletedPairGames(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	uUUID, err := getUuidFromAuth(req)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	games, err := gh.service.GetAllCompletedPairGames(ctx, uUUID)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	sendPairCompletedList(rw, games)
}

func (gh *GameHandler) HandleLeaderBoard(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	limit := defLeaderBoardLimit
	limitStr := req.URL.Query().Get("limit")
	if limitStr != "" {
		limitCnv, err := strconv.Atoi(limitStr)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
		if limitCnv <= 0 {
			http.Error(rw, "Limit can't be <= 0", http.StatusBadRequest)
			return
		}
		limit = limitCnv
	}
	games, err := gh.service.GetLeaderBoard(ctx, uint(limit))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	sendLeaderBoard(rw, games)
}

func (gh *GameHandler) HandleJoinPairGame(rw http.ResponseWriter, req *http.Request) {
	gUUID, err := getUuidFromPath(req)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	sUUID, err := getUuidFromAuth(req)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	ctx := req.Context()
	game, err := gh.service.JoinPairGame(ctx, gUUID, sUUID)
	if err != nil {
		http.Error(rw, fmt.Sprintf("Trouble with join to game: %q", err.Error()), http.StatusBadRequest)
		return
	}
	sendFullGamePair(rw, game)
}

func (gh *GameHandler) handleSoloGame(rw http.ResponseWriter, req *http.Request, gameFieldIncome GameFieldWeb, gUUID uuid.UUID) {
	ctx := req.Context()

	gameSaved, err := gh.getSoloGame(ctx, gUUID)
	if err != nil {
		http.Error(rw, fmt.Sprintf("handleSoloGame trouble, err: %q", err), http.StatusInternalServerError)
		return
	}

	game, err := FromWebShortSolo(gameFieldIncome, gameSaved)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	is_valid, err := gh.service.IsValid(ctx, game)
	if !is_valid || err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	is_over, err := gh.service.IsOver(ctx, game)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if !is_over {
		game, err = gh.service.MakeMove(ctx, game)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
		_, err := gh.service.IsOver(ctx, game)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	sendShortGameSolo(rw, game)
}

func (gh *GameHandler) handlePairGame(rw http.ResponseWriter, req *http.Request, gameFieldIncome GameFieldWeb, gUUID uuid.UUID) {

	ctx := req.Context()
	game, err := gh.service.GetPairGame(ctx, gUUID)

	if err != nil {
		http.Error(rw, fmt.Sprintf("handlePairGame trouble, err: %q", err), http.StatusInternalServerError)
		return
	}
	if game == nil {
		http.Error(rw, "handlePairGame trouble: game == nil", http.StatusInternalServerError)
		return
	}

	gameR, err := FromWebShortPair(gameFieldIncome, game)
	if err != nil {
		http.Error(rw, fmt.Sprintf("handlePairGame, FromWebPair err: %q", err), http.StatusInternalServerError)
		return
	}

	isValid, err := gh.service.IsValidPair(ctx, gameR)
	if !isValid {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(rw, fmt.Sprintf("handlePairGame, IsValidPair err: %q", err), http.StatusInternalServerError)
		return
	}

	is_over, err := gh.service.IsOverPair(ctx, gameR)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if !is_over {
		game, err = gh.service.MakeMovePair(ctx, gameR)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
		_, err := gh.service.IsOverPair(ctx, game)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	sendShortGamePair(rw, game)
}

func (gh *GameHandler) HandleMove(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	gUUID, err := getUuidFromPath(req)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	var gameFieldIncome GameFieldWeb
	err = json.NewDecoder(req.Body).Decode(&gameFieldIncome)
	if err != nil {
		http.Error(rw, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	isPair, err := gh.service.IsPairGame(ctx, gUUID)
	if err != nil {
		http.Error(rw, fmt.Sprintf("IsPairGame trouble, err: %q", err), http.StatusInternalServerError)
		return
	}

	if isPair {
		gh.handlePairGame(rw, req, gameFieldIncome, gUUID)
	} else {
		gh.handleSoloGame(rw, req, gameFieldIncome, gUUID)
	}
}

func (gh *GameHandler) SendGameShortInfo(rw http.ResponseWriter, req *http.Request) {
	gUUID, err := getUuidFromPath(req)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	ctx := req.Context()
	isPair, err := gh.service.IsPairGame(ctx, gUUID)
	if err != nil {
		http.Error(rw, fmt.Sprintf("IsPairGame trouble, err: %q", err), http.StatusInternalServerError)
		return
	}
	if isPair {
		game, err := gh.getPairGame(ctx, gUUID)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		sendShortGamePair(rw, game)
	} else {
		game, err := gh.getSoloGame(ctx, gUUID)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		sendShortGameSolo(rw, game)
	}
}

func (gh *GameHandler) SendGameFullInfo(rw http.ResponseWriter, req *http.Request) {
	gUUID, err := getUuidFromPath(req)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	ctx := req.Context()
	isPair, err := gh.service.IsPairGame(ctx, gUUID)
	if err != nil {
		http.Error(rw, fmt.Sprintf("IsPairGame trouble, err: %q", err), http.StatusInternalServerError)
		return
	}
	if isPair {
		game, err := gh.getPairGame(ctx, gUUID)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		sendFullGamePair(rw, game)
	} else {
		game, err := gh.getSoloGame(ctx, gUUID)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		sendFullGameSolo(rw, game)
	}
}

func (gh *GameHandler) SendPlayerInfo(rw http.ResponseWriter, req *http.Request) {
	pUUID, err := getUuidFromPath(req)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	ctx := req.Context()
	name, err := gh.service.GetPlayerName(ctx, pUUID)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	pi := ToPlayerInfo(name)
	send(rw, pi)
}
