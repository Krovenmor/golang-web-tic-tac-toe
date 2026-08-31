package gamehandler

import (
	"time"

	"github.com/google/uuid"
)

type GameFieldWeb struct {
	Field [3][3]int `json:"field"`
}

type GameCreateResponseBody struct {
	GUUID     string       `json:"gameuuid"`
	GameField GameFieldWeb `json:"gamefield"`
	State     int          `json:"state"`
}

type GamePairCreateResponseBody struct {
	GUUID     uuid.UUID    `json:"gameuuid"`
	FUUID     uuid.UUID    `json:"fPlayerUuid"`
	SUUID     uuid.UUID    `json:"sPlayerUuid"`
	GameField GameFieldWeb `json:"gamefield"`
	State     int          `json:"state"`
}

type GameShortInfoResponseBody struct {
	GameField GameFieldWeb `json:"gamefield"`
	State     int          `json:"state"`
}

type GamePairListResponseBody struct {
	Games []uuid.UUID `json:"games"`
}

type GamePairInfoResponseBody struct {
	GUUID   uuid.UUID `json:"gameuuid"`
	WUUID   uuid.UUID `json:"winneruuid"`
	Created time.Time `json:"datecreated"`
	State   int       `json:"state"`
}

type GameLeaderBoardEntryResponseBody struct {
	PlayerUUID  uuid.UUID `json:"uuid"`
	PlayerLogin string    `json:"login"`
	Ratio       float64   `json:"w_ratio"`
}

type PlayerInfoResponseBody struct {
	PlayerName string `json:"playername"`
}

type UuidIncomeBody struct {
	UUID uuid.UUID `json:"uuid"`
}
