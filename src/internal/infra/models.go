package infra

import (
	"time"

	"github.com/google/uuid"
)

type GameFieldPG struct {
	Field [3][3]int `json:"field"`
}

type GamePG struct {
	Field GameFieldPG
	UUID  uuid.UUID
}

type GamePairPG struct {
	Field GameFieldPG
	Guuid uuid.UUID
	Fuuid uuid.UUID
	Suuid uuid.UUID
	State int
}

type GamePairInfoPG struct {
	Guuid   uuid.UUID `db:"uuid"`
	Wuuid   uuid.UUID `db:"winner_uuid"`
	State   int       `db:"state"`
	Created time.Time `db:"created"`
}

type GameLeaderBoardPG struct {
	PlayerUUID  uuid.UUID `db:"player"`
	PlayerLogin string    `db:"login"`
	Ratio       float64   `db:"ratio"`
}
