package module

import (
	"context"
	"fmt"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"os"
	"time"
)

type Postgres struct {
	Client *pgxpool.Pool
	Ctx    context.Context
}

var UserActionLogPool [][]interface{}
var UserActionLogPoolSize = 200
var UserActionLogPoolFull = make(chan bool)

func CreatePostgresClient(e *Environments) *Postgres {
	ctx := context.Background()
	dsn := fmt.Sprintf("host=/var/run/postgresql port=%s user=%s password=%s dbname=%s sslmode=disable", e.DbPort, e.DbUser, e.DbPass, e.DbName)
	db, err := pgxpool.Connect(ctx, dsn)
	db.Config().MaxConns = 50
	db.Config().MaxConnIdleTime = 10 * time.Second
	if err != nil {
		fmt.Printf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	//defer db.Close(ctx)
	return &Postgres{Client: db, Ctx: ctx}
}

func ScanRakutanInfo(rows pgx.Rows) RakutanInfos {
	var rakutanInfos RakutanInfos
	defer rows.Close()
	for rows.Next() {
		var id int
		var facultyName, lectureName string
		var register, passed pgtype.Int2Array

		err := rows.Scan(&id, &facultyName, &lectureName, &register, &passed)
		if err != nil {
			log.Println(err)
		}
		rakutanInfos = append(rakutanInfos, RakutanInfo{
			ID:          id,
			FacultyName: facultyName,
			LectureName: lectureName,
			Register:    register,
			Passed:      passed,
		})
	}
	return rakutanInfos
}

func (p *Postgres) InsertUser(uid string) bool {
	result, err := p.Client.Exec(p.Ctx, "INSERT INTO users (uid, is_verified, registered_at) VALUES ($1, $2)", uid, false, time.Now())
	if err != nil {
		log.Println(err)
		return false
	}
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return false
	}
	return true
}

func (p *Postgres) InsertUserAction(userID string, action UserAction) error {
	_, err := p.Client.Exec(p.Ctx, "INSERT INTO user_logs (uid, action, timestamp) VALUES ($1, $2, $3)", userID, action, time.Now())
	if err != nil {
		return err
	}
	return nil
}

func (p *Postgres) BulkInsertUserAction() error {
	copyCount, err := p.Client.CopyFrom(
		p.Ctx,
		pgx.Identifier{"user_logs"},
		[]string{"uid", "action", "timestamp"},
		pgx.CopyFromRows(UserActionLogPool),
	)
	if err != nil {
		return err
	}
	if copyCount != 0 {
		log.Printf("[Bulk] Copied %d rows\n", copyCount)
	}
	ClearUserActionLogPool()
	return nil
}

func (p *Postgres) UpsertVerificationToken(uid string, token string) error {
	_, err := p.Client.Exec(p.Ctx, "INSERT INTO verification_tokens (uid, token, created_at) VALUES ($1, $2, $3) ON CONFLICT (uid) DO UPDATE SET token = $2, created_at = $3", uid, token, time.Now())
	if err != nil {
		return err
	}
	return nil
}

func (p *Postgres) CheckVerificationToken(uid, token string) (bool, error) {
	var i int
	err := p.Client.QueryRow(p.Ctx, "SELECT count(*) FROM verification_tokens WHERE uid = $1 AND token = $2", uid, token).Scan(&i)
	if err != nil {
		return false, err
	}
	return i > 0, nil
}

func (p *Postgres) DeleteVerificationToken(uid string) error {
	_, err := p.Client.Exec(p.Ctx, "DELETE FROM verification_tokens WHERE uid = $1", uid)
	if err != nil {
		return err
	}
	return nil
}

func (p *Postgres) UpdateUserVerification(uid string) error {
	_, err := p.Client.Exec(p.Ctx, "UPDATE users SET is_verified = true, verified_at = $1 WHERE uid = $2", time.Now(), uid)
	if err != nil {
		return err
	}
	err = p.DeleteVerificationToken(uid)
	if err != nil {
		return err
	}
	return nil
}

func (p *Postgres) IsRegistered(uid string) error {
	_, err := p.Client.Exec(p.Ctx, "INSERT INTO users(uid) SELECT $1 WHERE NOT EXISTS ( SELECT 1 FROM users WHERE uid = $2)", uid, uid)
	if err != nil {
		return err
	}
	return nil
}

func (p *Postgres) IsVerified(uid string) (bool, error) {
	var isVerified bool
	err := p.Client.QueryRow(p.Ctx, "SELECT is_verified FROM users WHERE uid = $1", uid).Scan(&isVerified)
	if err != nil {
		log.Println(err)
		return false, err
	}
	return isVerified, nil
}

func (p *Postgres) GetRakutanInfoByID(id int) (ExecStatus[RakutanInfos], bool) {
	var status ExecStatus[RakutanInfos]
	rows, err := p.Client.Query(p.Ctx, "SELECT * FROM rakutan WHERE id = $1", id)
	if err != nil {
		log.Println(err)
		status.Err = ErrorMessageGetRakutanInfoByIDError
		return status, false
	}
	status.Result = ScanRakutanInfo(rows)
	return status, true
}

func (p *Postgres) GetRakutanInfoByIds(ids []int) (ExecStatus[RakutanInfos], bool) {
	var status ExecStatus[RakutanInfos]
	rows, err := p.Client.Query(p.Ctx, "SELECT * FROM rakutan WHERE id = ANY($1)", ids)
	if err != nil {
		log.Println(err)
		status.Err = ErrorMessageGetRakutanInfoByIDError
		return status, false
	}
	status.Result = ScanRakutanInfo(rows)
	return status, true
}

func (p *Postgres) GetRakutanInfoByLectureName(lectureName string, subStringSearch bool) (ExecStatus[RakutanInfos], bool) {
	var status ExecStatus[RakutanInfos]
	var rows pgx.Rows
	var err error
	if subStringSearch {
		rows, err = p.Client.Query(p.Ctx, "SELECT * FROM rakutan WHERE lecture_name ILIKE CONCAT('%%', $1::text,'%%')", lectureName)
	} else {
		rows, err = p.Client.Query(p.Ctx, "SELECT * FROM rakutan WHERE lecture_name ILIKE CONCAT($1::text,'%%')", lectureName)
	}

	if err != nil {
		log.Println(err)
		status.Err = ErrorMessageGetRakutanInfoByNameError
		return status, false
	}
	status.Result = ScanRakutanInfo(rows)
	return status, true
}

func (p *Postgres) GetAllIDByOmikuji10() (RakutanInfoIDs, bool) {
	var ids RakutanInfoIDs
	rows, err := p.Client.Query(p.Ctx, "SELECT id, faculty_name, lecture_name, register, passed FROM mat_view_omikuji10")
	if err != nil {
		log.Println(err)
		return ids, false
	}
	rakutanInfos := ScanRakutanInfo(rows)
	for _, rakutanInfo := range rakutanInfos {
		ids = append(ids, rakutanInfo.ID)
	}
	return ids, true
}

func (p *Postgres) GetAllIDByOmikuji(types OmikujiType) (RakutanInfoIDs, bool) {
	var ids RakutanInfoIDs
	var query string
	switch types {
	case Rakutan:
		query = "SELECT id, faculty_name, lecture_name, register, passed FROM mat_view_rakutan"
	case Onitan:
		query = "SELECT id, faculty_name, lecture_name, register, passed FROM mat_view_onitan"
	}
	rows, err := p.Client.Query(p.Ctx, query)
	if err != nil {
		log.Println(err)
		return ids, false
	}
	rakutanInfos := ScanRakutanInfo(rows)
	for _, rakutanInfo := range rakutanInfos {
		ids = append(ids, rakutanInfo.ID)
	}
	return ids, true
}

func (p *Postgres) GetRakutanInfoByOmikuji(types OmikujiType) (ExecStatus[RakutanInfos], bool) {
	var status ExecStatus[RakutanInfos]
	var err error
	var rows pgx.Rows
	switch types {
	case Rakutan:
		rows, err = p.Client.Query(p.Ctx, "SELECT id, faculty_name, lecture_name, register, passed FROM mat_view_rakutan ORDER BY random() LIMIT 1")
	case Onitan:
		rows, err = p.Client.Query(p.Ctx, "SELECT id, faculty_name, lecture_name, register, passed FROM mat_view_onitan ORDER BY random() LIMIT 1")
	}
	if err != nil {
		log.Println(err)
		status.Err = ErrorMessageGetRakutanInfoByOmikujiError
		return status, false
	}
	status.Result = ScanRakutanInfo(rows)
	return status, true
}

func (p *Postgres) GetFavorites(uid string) (ExecStatus[RakutanInfos], bool) {
	var status ExecStatus[RakutanInfos]
	rows, err := p.Client.Query(p.Ctx, "SELECT r.* FROM favorites as f INNER JOIN rakutan as r ON f.id = r.id WHERE f.uid = $1", uid)
	if err != nil {
		log.Println(err)
		status.Err = ErrorMessageGetFavError
		return status, false
	}
	status.Result = ScanRakutanInfo(rows)
	return status, true
}

func (p *Postgres) GetFavoriteByID(uid string, id int) (ExecStatus[RakutanInfos], bool) {
	var status ExecStatus[RakutanInfos]
	rows, err := p.Client.Query(p.Ctx, "SELECT r.* FROM favorites as f INNER JOIN rakutan as r ON f.id = r.id WHERE f.uid = $1 AND f.id = $2", uid, id)
	if err != nil {
		log.Println(err)
		status.Err = ErrorMessageGetFavError
		return status, false
	}
	status.Result = ScanRakutanInfo(rows)
	return status, true
}

func (p *Postgres) ToggleFavorite(uid string, id int) (string, bool) {
	var favoriteIDs []int
	rows, err := p.Client.Query(p.Ctx, "SELECT id FROM favorites WHERE uid = $1", uid)
	if err != nil {
		log.Println(err)
		return ErrorMessageGetFavError, false
	}
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			log.Println(err)
			return ErrorMessageGetFavError, false
		}
		favoriteIDs = append(favoriteIDs, id)
	}
	for _, favoriteID := range favoriteIDs {
		if favoriteID == id {
			return p.UnsetFavorite(uid, id)
		}
	}
	if len(favoriteIDs) >= 50 {
		return ErrorMessageFavLimitError, false
	}
	_, err = p.Client.Exec(p.Ctx, "INSERT INTO favorites (uid, id) VALUES ($1, $2)", uid, id)
	// TODO: Duplicate key errorをチェックする
	if err != nil {
		log.Println(err)
		return ErrorMessageInsertFavError, false
	}
	return SuccessMessageInsertFav, true
}

func (p *Postgres) UnsetFavorite(uid string, id int) (string, bool) {
	_, err := p.Client.Exec(p.Ctx, "DELETE FROM favorites WHERE uid = $1 AND id = $2", uid, id)
	if err != nil {
		log.Println(err)
		return ErrorMessageDeleteFavError, false
	}
	return SuccessMessageDeleteFav, true
}

type FindByMethod int

const (
	Name FindByMethod = iota
	ID
	Omikuji
	Omikuji10
)

func GetRakutanInfo(c Clients, e *Environments, uid string, method FindByMethod, value interface{}) (ExecStatus[RakutanInfos], bool) {
	var ok bool
	var status ExecStatus[RakutanInfos]

	switch method {
	case ID:
		status, ok = GetRakutanInfoByID(c, value.(int))
	case Name:
		var subStringSearch bool
		searchWord := value.(string)
		if search := []rune(value.(string)); string(search[:1]) == "%" || string(search[:1]) == "％" {
			subStringSearch = true
			searchWord = string(search[1:])
		}
		status, ok = c.Postgres.GetRakutanInfoByLectureName(searchWord, subStringSearch)
	case Omikuji:
		var ids []int
		ids, ok = c.Redis.GetRandomOmikuji(value.(OmikujiType), 1)
		if !ok {
			status, ok = c.Postgres.GetRakutanInfoByOmikuji(value.(OmikujiType))
		} else {
			status, ok = GetRakutanInfoByID(c, ids[0])
		}

	case Omikuji10:
		var ids []int
		ids, ok = c.Redis.GetRandomOmikuji(value.(OmikujiType), 10)
		if ok {
			status, ok = GetRakutanInfoByIDs(c, ids)
		} else {
			status.Err = ErrorMessageGetRakutanInfoIDsByOmikujiError
		}
	}

	// Set isVerified, isFavorite and kakomonURL
	if ok && len(status.Result) == 1 {
		if faforites, ok := c.Postgres.GetFavoriteByID(uid, status.Result[0].ID); ok && len(faforites.Result) == 1 {
			status.Result[0].IsFavorite = true
		}
		if isVerified, err := c.Postgres.IsVerified(uid); err == nil && isVerified {
			status.Result[0].IsVerified = true
			if kakomonStatus, found := GetKakomonURL(c, e, status.Result[0].ID, status.Result[0].LectureName); found {
				status.Result[0].KakomonURL = string(any(kakomonStatus.Result).(KUWikiKakomon))
			}
		}
	}

	return status, ok
}

func GetRakutanInfoByID(c Clients, id int) (ExecStatus[RakutanInfos], bool) {
	var status ExecStatus[RakutanInfos]
	var ok bool
	status, ok = c.Redis.GetRakutanInfoByID(id)
	if !ok {
		status, ok = c.Postgres.GetRakutanInfoByID(id)
		if ok {
			redisKey := fmt.Sprintf("rinfo:%d", status.Result[0].ID)
			go c.Redis.SetRedis(redisKey, status.Result[0], 0)
		}
	}
	return status, ok
}

func GetRakutanInfoByIDs(c Clients, ids []int) (ExecStatus[RakutanInfos], bool) {
	var status ExecStatus[RakutanInfos]
	var ok bool
	var missingIDs []int
	status, missingIDs = c.Redis.GetRakutanInfoByIDs(ids)
	if len(missingIDs) == 0 {
		return status, true
	}
	var pStatus ExecStatus[RakutanInfos]
	pStatus, ok = c.Postgres.GetRakutanInfoByIds(missingIDs)
	if ok {
		status.Result = append(status.Result, pStatus.Result...)
		for _, v := range pStatus.Result {
			redisKey := fmt.Sprintf("rinfo:%d", v.ID)
			go c.Redis.SetRedis(redisKey, v, 0)
		}
	}

	return status, ok
}

func AppendUserActionLogPool(uid string, action UserAction) {
	UserActionLogPool = append(UserActionLogPool, []interface{}{uid, action, time.Now()})
	if len(UserActionLogPool) >= UserActionLogPoolSize {
		UserActionLogPoolFull <- true
	}
}

func ClearUserActionLogPool() {
	UserActionLogPool = UserActionLogPool[:0]
}
