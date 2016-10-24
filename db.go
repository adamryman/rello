package rello

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

var db *sql.DB

func InitDatabase(conn string) error {
	d, err := sql.Open("sqlite3", conn)
	if err != nil {
		return errors.Wrapf(err, "cannot connect to %s", conn)
	}

	db = d

	return nil
}

func createCheck(in Check) (int64, error) {
	const query = "INSERT INTO checks(checkItemId, datetime) values(?, ?)"

	stmt, err := db.Prepare(query)
	if err != nil {
		return -1, err
	}

	res, err := stmt.Exec(in.CheckItemId, in.Date)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()

	return id, nil
}

func createCheckItem(in CheckItem) (int64, error) {
	const query = "INSERT INTO checkItems(name, trelloId, userId) values(?, ?, ?)"

	stmt, err := db.Prepare(query)
	if err != nil {
		return -1, err
	}

	res, err := stmt.Exec(in.Name, in.Id, in.UserId)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()

	return id, nil
}

func readCheckItemByTrelloId(tId string) (*CheckItem, error) {
	const query = "SELECT * FROM checkItems WHERE trelloId=?"

	row := db.QueryRow(query, tId)

	var (
		id     int64
		name   string
		userId string
	)
	err := row.Scan(&id, &name, &tId, &userId)
	if err != nil {
		return nil, err
	}

	return &CheckItem{
		Id:      tId,
		Name:    name,
		RelloId: id,
		UserId:  userId,
	}, nil
}

//func CreateCheckByTrelloId() (int64, error) {

//}
