package app

import (
	"time"

	"github.com/google/uuid"
)

const (
	CELL_EMPTY = iota
	CELL_PLAYER
	CELL_COMPUTER
)

type GameField struct {
	Field [3][3]int
}

type CurrentGame struct {
	GField GameField
	UUID   uuid.UUID
	State  int
}

type CurrentGamePair struct {
	GField GameField
	Guuid  uuid.UUID
	Fuuid  uuid.UUID
	Suuid  uuid.UUID
	State  int
}

type GamePairInfo struct {
	Guuid   uuid.UUID
	Wuuid   uuid.UUID
	State   int
	Created time.Time
}

type GameLeaderBoardEntry struct {
	PlayerUUID  uuid.UUID
	PlayerLogin string
	Ratio       float64
}
