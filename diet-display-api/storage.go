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
		night TEXT
	);
	CREATE INDEX diet_user_id_date_idx ON diet (date,user_id);`
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
