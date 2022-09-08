package module

import (
	"context"
	"fmt"
	rakutan "github.com/das08/kuRakutanBot-go/models/rakutan"
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

type RakutanType interface {
	[]rakutan.RakutanInfo2
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

func (p *Postgres) InsertUserAction(userID string, action UserAction) error {
	_, err := p.Client.Exec("INSERT INTO user_logs (uid, action, timestamp) VALUES ($1, $2, $3)", userID, action, time.Now())
	if err != nil {
		return err
	}
	return nil
}

func (p *Postgres) getRakutanInfoByID(id int) (QueryStatus2, error) {
	var status QueryStatus2
	var rakutanInfos []rakutan.RakutanInfo2
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

func (p *Postgres) getRakutanInfoByLectureName(lectureName string) (QueryStatus2, error) {
	var status QueryStatus2
	var rakutanInfos []rakutan.RakutanInfo2
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

func (p *Postgres) getFavorites(uid string) (QueryStatus2, error) {
	var status QueryStatus2
	var rakutanInfos []rakutan.RakutanInfo2
	err := p.Client.Select(&rakutanInfos, "SELECT * FROM favorites as f INNER JOIN rakutan2021 as r WHERE f.id = r.id AND f.uid = $1", uid)
	if err != nil {
		log.Println(err)
		status.Err = err.Error()
		return status, err
	}
	status.Result = rakutanInfos
	return status, nil
}
