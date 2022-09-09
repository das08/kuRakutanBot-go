package module

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type Postgres struct {
	Client *sqlx.DB
	Ctx    context.Context
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
)

type RakutanInfo2 struct {
	ID          int    `db:"id"`
	FacultyName string `db:"faculty_name"`
	LectureName string `db:"lecture_name"`
	Register    []int  `db:"register"`
	Passed      []int  `db:"passed"`
	KakomonURL  string `db:"kakomon_url"`
	IsFavorite  bool
}

func (r *RakutanInfo2) GetLatestDetail() (int, int) {
	passed, register := 0, 0
	for i := 0; i < len(r.Register); i++ {
		if r.Register[i] != 0 {
			passed = r.Passed[i]
			register = r.Register[i]
			break
		}
	}
	return passed, register
}

type ReturnType interface {
	[]RakutanInfo2 | []FlexMessage
}

type QueryStatus2[T ReturnType] struct {
	Result T
	Err    string
}

func CreatePostgresClient(e *Environments) *Postgres {
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", e.DB_USER, e.DB_PASS, e.DB_NAME)
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	//defer db.Close()
	return &Postgres{Client: db, Ctx: context.Background()}
}

func (p *Postgres) InsertUser(uid string) bool {
	result, err := p.Client.Exec("INSERT INTO users (uid, is_verified, registered_at) VALUES ($1, $2)", uid, false, time.Now())
	if err != nil {
		log.Println(err)
		return false
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return false
	}
	return true
}

func (p *Postgres) InsertUserAction(userID string, action UserAction) error {
	_, err := p.Client.Exec("INSERT INTO user_logs (uid, action, timestamp) VALUES ($1, $2, $3)", userID, action, time.Now())
	if err != nil {
		return err
	}
	return nil
}

func (p *Postgres) InsertVerificationToken(uid string, token string) error {
	_, err := p.Client.Exec("INSERT INTO verification_tokens (uid, token, created_at) VALUES ($1, $2, $3)", uid, token, time.Now())
	if err != nil {
		return err
	}
	return nil
}

func (p *Postgres) CheckVerificationToken(uid, token string) (bool, error) {
	var i int
	err := p.Client.Get(&i, "SELECT count(*) FROM verification_tokens WHERE uid = $1 AND token = $2", uid, token)
	if err != nil {
		return false, err
	}
	return i > 0, nil
}

func (p *Postgres) UpdateUserVerification(uid string) error {
	_, err := p.Client.Exec("UPDATE users SET is_verified = true WHERE uid = $1", uid)
	if err != nil {
		return err
	}
	return nil
}

func (p *Postgres) IsVerified(uid string) (bool, error) {
	var isVerified bool
	err := p.Client.Get(&isVerified, "SELECT is_verified FROM users WHERE uid = $1", uid)
	if err != nil {
		log.Println(err)
		return false, err
	}
	return isVerified, nil
}

func (p *Postgres) GetRakutanInfoByID(id int) (QueryStatus2[[]RakutanInfo2], bool) {
	var status QueryStatus2[[]RakutanInfo2]
	var rakutanInfos []RakutanInfo2
	err := p.Client.Select(&rakutanInfos, "SELECT * FROM rakutan WHERE id = $1", id)
	if err != nil {
		log.Println(err)
		status.Err = ErrorMessageGetRakutanInfoByIDError
		return status, false
	}
	status.Result = rakutanInfos
	return status, true
}

func (p *Postgres) GetRakutanInfoByLectureName(lectureName string) (QueryStatus2[[]RakutanInfo2], bool) {
	var status QueryStatus2[[]RakutanInfo2]
	var rakutanInfos []RakutanInfo2
	// TODO: consider LIKE search
	err := p.Client.Select(&rakutanInfos, "SELECT * FROM rakutan WHERE lecture_name = $1", lectureName)
	if err != nil {
		log.Println(err)
		status.Err = ErrorMessageGetRakutanInfoByNameError
		return status, false
	}
	status.Result = rakutanInfos
	return status, true
}

func (p *Postgres) GetRakutanInfoByOmikuji(types OmikujiType) (QueryStatus2[[]RakutanInfo2], bool) {
	var status QueryStatus2[[]RakutanInfo2]
	var rakutanInfos []RakutanInfo2
	var err error
	switch types {
	case Rakutan:
		err = p.Client.Select(&rakutanInfos, "SELECT * FROM mat_view_rakutan ORDER BY random() LIMIT 1")
	case Onitan:
		err = p.Client.Select(&rakutanInfos, "SELECT * FROM mat_view_onitan ORDER BY random() LIMIT 1")
	}
	if err != nil {
		log.Println(err)
		status.Err = ErrorMessageGetRakutanInfoByOmikujiError
		return status, false
	}
	status.Result = rakutanInfos
	return status, true
}

func (p *Postgres) GetFavorites(uid string) (QueryStatus2[[]RakutanInfo2], bool) {
	var status QueryStatus2[[]RakutanInfo2]
	var rakutanInfos []RakutanInfo2
	err := p.Client.Select(&rakutanInfos, "SELECT r.* FROM favorites as f INNER JOIN rakutan as r WHERE f.id = r.id AND f.uid = $1", uid)
	if err != nil {
		log.Println(err)
		status.Err = ErrorMessageGetFavError
		return status, false
	}
	status.Result = rakutanInfos
	return status, true
}

func (p *Postgres) GetFavoriteByID(uid string, id int) (QueryStatus2[[]RakutanInfo2], bool) {
	var status QueryStatus2[[]RakutanInfo2]
	var rakutanInfos []RakutanInfo2
	err := p.Client.Select(&rakutanInfos, "SELECT r.* FROM favorites as f INNER JOIN rakutan as r WHERE f.id = r.id AND f.uid = $1 AND f.id = $2", uid, id)
	if err != nil {
		log.Println(err)
		status.Err = ErrorMessageGetFavError
		return status, false
	}
	status.Result = rakutanInfos
	return status, true
}

func (p *Postgres) SetFavorite(uid string, id int) (string, bool) {
	var favoriteIDs []int
	err := p.Client.Select(&favoriteIDs, "SELECT id FROM favorites WHERE uid = $1", uid)
	if err != nil {
		log.Println(err)
		return ErrorMessageGetFavError, false
	}
	for _, favoriteID := range favoriteIDs {
		if favoriteID == id {
			return ErrorMessageAlreadyFavError, false
		}
	}
	if len(favoriteIDs) >= 50 {
		return ErrorMessageFavLimitError, false
	}

	_, err = p.Client.Exec("INSERT INTO favorites (uid, id) VALUES ($1, $2)", uid, id)
	// TODO: Duplicate key errorをチェックする
	if err != nil {
		log.Println(err)
		return ErrorMessageInsertFavError, false
	}

	// TODO: 講義名を取得する
	return fmt.Sprintf(SuccessMessageInsertFav, ""), true
}

func (p *Postgres) UnsetFavorite(uid string, id int) (string, bool) {
	_, err := p.Client.Exec("DELETE FROM favorites WHERE uid = $1 AND id = $2", uid, id)
	if err != nil {
		log.Println(err)
		return ErrorMessageDeleteFavError, false
	}
	// TODO: 講義名を取得する
	return fmt.Sprintf(SuccessMessageDeleteFav, ""), true
}
