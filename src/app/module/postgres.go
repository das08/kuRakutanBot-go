package module

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
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
	ID          int           `db:"id"`
	FacultyName string        `db:"faculty_name"`
	LectureName string        `db:"lecture_name"`
	Register    pq.Int32Array `db:"register"`
	Passed      pq.Int32Array `db:"passed"`
}

type RakutanType interface {
	[]RakutanInfo2
}

type QueryStatus2[T RakutanType] struct {
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

func (p *Postgres) GetRakutanInfoByID(id int) (QueryStatus2, error) {
	var status QueryStatus2
	var rakutanInfos []RakutanInfo2
	// TODO: do not hard code table name
	err := p.Client.Select(&rakutanInfos, "SELECT * FROM rakutan2021 WHERE id = $1", id)
	if err != nil {
		log.Println(err)
		status.Err = err.Error()
		return status, err
	}
	status.Result = rakutanInfos
	return status, nil
}

func (p *Postgres) GetRakutanInfoByLectureName(lectureName string) (QueryStatus2, error) {
	var status QueryStatus2
	var rakutanInfos []RakutanInfo2
	// TODO: consider LIKE search
	err := p.Client.Select(&rakutanInfos, "SELECT * FROM rakutan2021 WHERE lecture_name = $1", lectureName)
	if err != nil {
		log.Println(err)
		status.Err = err.Error()
		return status, err
	}
	status.Result = rakutanInfos
	return status, nil
}

func (p *Postgres) GetFavorites(uid string) (QueryStatus2, error) {
	var status QueryStatus2
	var rakutanInfos []RakutanInfo2
	err := p.Client.Select(&rakutanInfos, "SELECT * FROM favorites as f INNER JOIN rakutan2021 as r WHERE f.id = r.id AND f.uid = $1", uid)
	if err != nil {
		log.Println(err)
		status.Err = err.Error()
		return status, err
	}
	status.Result = rakutanInfos
	return status, nil
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
