package module

import (
	"github.com/jackc/pgtype"
)

type Clients struct {
	Postgres *Postgres
	Redis    *Redis
}

type OmikujiType string

const (
	Normal  OmikujiType = "normal"
	Rakutan OmikujiType = "rakutan"
	Onitan  OmikujiType = "onitan"
	Ten     OmikujiType = "all"
)

type OmikujiText struct {
	Text  string
	Color string
}

type UserAction string

const (
	UserActionSearch    UserAction = "search"
	UserActionRakutan   UserAction = "rakutan"
	UserActionOnitan    UserAction = "onitan"
	UserActionOmikuji10 UserAction = "omikuji10"
	UserActionSetFav    UserAction = "set_fav"
	UserActionUnsetFav  UserAction = "unset_fav"
	UserActionGetFav    UserAction = "get_fav"
	UserActionInfo      UserAction = "info"
	UserActionHelp      UserAction = "help"
	UserActionEmail     UserAction = "email"
	UserActionVerify    UserAction = "verify"
)

type RakutanInfo struct {
	ID          int              `db:"id" json:"id"`
	FacultyName string           `db:"faculty_name" json:"fn"`
	LectureName string           `db:"lecture_name" json:"ln"`
	Register    pgtype.Int2Array `db:"register" json:"r"`
	Passed      pgtype.Int2Array `db:"passed" json:"p"`
	KakomonURL  string           `db:"kakomon_url" json:"k"`
	IsFavorite  bool             `json:"-"`
	IsVerified  bool             `json:"-"`
}

func (r *RakutanInfo) GetLatestDetail() (int, int) {
	passed, register := 0, 0
	for i := 0; i < len(r.Register.Elements); i++ {
		if r.Register.Elements[i].Status == pgtype.Present {
			passed = int(r.Passed.Elements[i].Int)
			register = int(r.Register.Elements[i].Int)
			break
		}
	}
	return passed, register
}

type RakutanInfos []RakutanInfo

type FlexMessages []FlexMessage

type RakutanInfoIDs []interface{}

type KUWikiKakomon string

type ReturnType interface {
	RakutanInfos | FlexMessages | RakutanInfoIDs | KUWikiKakomon
}

type ExecStatus[T ReturnType] struct {
	Result T
	Err    string
}
