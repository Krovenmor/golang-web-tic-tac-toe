package gamehandler

import (
	"WebTic-tac-toe2/internal/app"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func send[T any](rw http.ResponseWriter, to_send *T) {
	rw.Header().Set("Content-type", "application/json")
	rw.WriteHeader(http.StatusOK)
	err := json.NewEncoder(rw).Encode(*to_send)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

func sendShortGameSolo(rw http.ResponseWriter, game *app.CurrentGame) {
	resp, err := ToWebShortSolo(game)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	send(rw, resp)
}

func sendShortGamePair(rw http.ResponseWriter, game *app.CurrentGamePair) {
	resp, err := ToWebShortPair(game)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	send(rw, resp)
}

func sendFullGameSolo(rw http.ResponseWriter, game *app.CurrentGame) {
	resp, err := ToWeb(game)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	send(rw, resp)
}

func sendFullGamePair(rw http.ResponseWriter, game *app.CurrentGamePair) {
	resp, err := ToWebPair(game)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	send(rw, resp)
}

func sendPairList(rw http.ResponseWriter, games []uuid.UUID) {
	pl, err := ToWebPairList(games)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	send(rw, pl)
}

func sendPairCompletedList(rw http.ResponseWriter, games []app.GamePairInfo) {
	pl, err := ToWebPairCompletedList(games)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	send(rw, &pl)
}

func sendLeaderBoard(rw http.ResponseWriter, games []app.GameLeaderBoardEntry) {
	pl, err := ToWebLeaderBoard(games)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	send(rw, &pl)
}
