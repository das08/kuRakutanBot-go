package module

import (
	"fmt"
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
	Ten     OmikujiType = "omikuji10"
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

func (r *RakutanInfo) GetRakutanPercent() float32 {
	p, reg := r.GetLatestDetail()
	return getPercentage(p, reg)
}

func (r *RakutanInfo) GetRakutanPercentBreakdown() []string {
	var rakutanPercent []string
	for i := 0; i < len(r.Passed.Elements); i++ {
		p := int(r.Passed.Elements[i].Int)
		reg := int(r.Register.Elements[i].Int)
		breakdown := fmt.Sprintf("(%d/%d)", p, reg)
		if r.Register.Elements[i].Status == pgtype.Null {
			rakutanPercent = append(rakutanPercent, "---% "+breakdown)
		} else {
			rakutanPercent = append(rakutanPercent, fmt.Sprintf("%.1f%% ", getPercentage(p, reg))+breakdown)
		}
	}
	return rakutanPercent
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
