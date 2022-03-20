package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const dbFilePath = "./diet.db"
const (
	Barathi = iota
	Durga
	Indira
	Sakthi
)

func DateToIntDate(t time.Time) (int, error) {
	dateStr := t.Format("20060102")
	dateInt, err := strconv.Atoi(dateStr)
	if err != nil {
		return 0, err
	}
	return dateInt, nil
}

func DateToDbIntDate(dateStr string) (int, error) {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return 0, err
	}
	return DateToIntDate(date)
}

func DbIntDateToDate(dateInt int) (string, error) {
	dateStr := strconv.Itoa(dateInt)
	date, err := time.Parse("20060102", dateStr)
	if err != nil {
		return "", err
	}
	dateStr = date.Format("2006-01-02")
	return dateStr, nil
}

func NameToID(name string) (int, error) {
	switch name {
	case "Barathi":
		return Barathi, nil
	case "Durga":
		return Durga, nil
	case "Indira":
		return Indira, nil
	case "Sakthi":
		return Sakthi, nil
	default:
		return -1, errors.New("Invalid name: " + name)
	}
}

func IDToName(ID int) (string, error) {
	switch ID {
	case Barathi:
		return "Barathi", nil
	case Durga:
		return "Durga", nil
	case Indira:
		return "Indira", nil
	case Sakthi:
		return "Sakthi", nil
	default:
		return "", errors.New("Invalid ID: " + strconv.Itoa(ID))
	}
}

func IsDatabaseCreated(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

func CreateDatabase(filePath string) (*sql.DB, error) {
	fmt.Println("Creating database:", filePath)
	db, err := sql.Open("sqlite3", filePath)
	if err != nil {
		return nil, err
	}

	ddlStatements := `CREATE TABLE diet (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		date INTEGER NOT NULL,
		morning TEXT,
		pre_breakfast TEXT,
		breakfast TEXT,
		noon TEXT,
		lunch TEXT,
		evening TEXT,
		dinner TEXT,
		post_dinner TEXT,
		night TEXT,
		UNIQUE(user_id,date)
	);`
	_, err = db.Exec(ddlStatements)

	if err != nil {
		return nil, err
	}
	return db, nil
}

func OpenDatabase() (*sql.DB, error) {
	if IsDatabaseCreated(dbFilePath) {
		fmt.Println("Opening database:", dbFilePath)
		return sql.Open("sqlite3", dbFilePath)
	} else {
		return CreateDatabase(dbFilePath)
	}
}

func CloseDatabase(db *sql.DB) {
	err := db.Close()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Closing database:", dbFilePath)
}

func ReadRecords(nextNdays int) ([]Record, error) {
	now := time.Now()
	till := now.Add(time.Hour * 24 * time.Duration(nextNdays))
	from, err := DateToIntDate(now)
	if err != nil {
		return nil, err
	}
	to, err := DateToIntDate(till)
	if err != nil {
		return nil, err
	}
	stmt, err := db.Prepare("SELECT * FROM diet WHERE date >= ? AND date <= ?")
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(from, to)
	if err != nil {
		return nil, err
	}
	records := make([]Record, 0, 100)
	// user_id,date,morning,pre_breakfast,breakfast,noon,lunch,evening,dinner,post_dinner,night
	for rows.Next() {
		var userID int
		var date int
		r := Record{}
		err = rows.Scan(&r.ID, &userID, &date, &r.Morning, &r.PreBreakfast, &r.Breakfast,
			&r.Noon, &r.Lunch, &r.Evening, &r.Dinner, &r.PostDinner, &r.Night)
		if err != nil {
			return nil, err
		}
		r.Name, err = IDToName(userID)
		if err != nil {
			return nil, err
		}
		r.Date, err = DbIntDateToDate(date)
		if err != nil {
			return nil, err
		}
		records = append(records, r)
	}
	rows.Close()
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return records, nil
}

func CreateRecords(db *sql.DB, records []Record) error {
	values := make([]string, 0, len(records))
	args := make([]interface{}, 0, len(records)*11)
	for _, r := range records {
		values = append(values, "(?,?,?,?,?,?,?,?,?,?,?)")
		nameID, err := NameToID(r.Name)
		if err != nil {
			return err
		}
		dbIntDate, err := DateToDbIntDate(r.Date)
		if err != nil {
			return err
		}
		args = append(args, nameID, dbIntDate, r.Morning, r.PreBreakfast,
			r.Breakfast, r.Noon, r.Lunch, r.Evening, r.Dinner, r.PostDinner, r.Night)
	}

	stmt := fmt.Sprintf(
		`INSERT INTO diet (user_id,date,morning,pre_breakfast,breakfast,noon,lunch,evening,dinner,post_dinner,night)
		VALUES %s`, strings.Join(values, ","))
	_, err := db.Exec(stmt, args...)
	return err
}

func DeleteRecord(db *sql.DB, id int) error {
	stmt, err := db.Prepare("DELETE FROM diet WHERE id = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(id)
	return err
}

func GetRecord(db *sql.DB, id int) (*Record, error) {
	stmt, err := db.Prepare("SELECT * FROM diet WHERE id = ?")
	if err != nil {
		return nil, err
	}
	row := stmt.QueryRow(id)
	var userID int
	var date int
	r := Record{}
	err = row.Scan(&r.ID, &userID, &date, &r.Morning, &r.PreBreakfast, &r.Breakfast,
		&r.Noon, &r.Lunch, &r.Evening, &r.Dinner, &r.PostDinner, &r.Night)
	switch err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		r.Name, err = IDToName(userID)
		if err != nil {
			return nil, err
		}
		r.Date, err = DbIntDateToDate(date)
		if err != nil {
			return nil, err
		}
		return &r, nil
	default:
		return nil, err
	}
}

func UpdateRecord(db *sql.DB, r Record) (*Record, error) {
	nameID, err := NameToID(r.Name)
	if err != nil {
		return nil, err
	}
	dbIntDate, err := DateToDbIntDate(r.Date)
	if err != nil {
		return nil, err
	}
	stmt := `UPDATE diet SET
		user_id = ?, date = ?, morning = ?, pre_breakfast = ?, breakfast = ?, 
		noon = ?, lunch = ?, evening = ?, dinner = ? ,post_dinner = ?,night = ?
		WHERE id = ?;`
	res, err := db.Exec(stmt, nameID, dbIntDate, r.Morning, r.PreBreakfast,
		r.Breakfast, r.Noon, r.Lunch, r.Evening, r.Dinner, r.PostDinner, r.Night, r.ID)
	if err != nil {
		return nil, err
	}
	c, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	switch c {
	case 0:
		return nil, nil
	case 1:
		return &r, nil
	default:
		return nil, errors.New(fmt.Sprintf("Number of rows affected %v > 1", c))
	}
}
