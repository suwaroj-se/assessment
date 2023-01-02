package expense

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type conDB struct {
	DB *sql.DB
}

func GetDB() *sql.DB {
	return db
}

func Connection(db *sql.DB) *conDB {
	return &conDB{db}
}

var db *sql.DB

func InitDB() {
	var err error
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Connect to database error", err)
	}

	createTB := `CREATE TABLE IF NOT EXISTS expenses ( id SERIAL PRIMARY KEY, title TEXT, amount FLOAT, note TEXT, tags TEXT[]);`

	_, err = db.Exec(createTB)
	if err != nil {
		log.Fatal("can't create table", err)
	}
}

func CloseDB() {
	db.Close()
}
