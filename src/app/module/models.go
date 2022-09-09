package module

import "github.com/jackc/pgtype"

type Clients struct {
	Postgres *Postgres
	Redis    *Redis
}

type OmikujiType string

const (
	Normal  OmikujiType = "normal"
	Rakutan OmikujiType = "rakutan"
	Onitan  OmikujiType = "onitan"
)

type OmikujiText struct {
	Text  string
	Color string
}

type UserAction string

const (
	UserActionSearch   UserAction = "search"
	UserActionRakutan  UserAction = "rakutan"
	UserActionOnitan   UserAction = "onitan"
	UserActionSetFav   UserAction = "set_fav"
	UserActionUnsetFav UserAction = "unset_fav"
	UserActionGetFav   UserAction = "get_fav"
	UserActionInfo     UserAction = "info"
	UserActionHelp     UserAction = "help"
	UserActionEmail    UserAction = "email"
	UserActionVerify   UserAction = "verify"
)

type RakutanInfo struct {
	ID          int              `db:"id"`
	FacultyName string           `db:"faculty_name"`
	LectureName string           `db:"lecture_name"`
	Register    pgtype.Int2Array `db:"register"`
	Passed      pgtype.Int2Array `db:"passed"`
	KakomonURL  string           `db:"kakomon_url"`
	IsFavorite  bool
	IsVerified  bool
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

type KUWikiKakomon string

type ReturnType interface {
	RakutanInfos | FlexMessages | KUWikiKakomon
}

type ExecStatus[T ReturnType] struct {
	Result T
	Err    string
}
