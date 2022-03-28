package db

import (
	"awesomeProjectRentaTeam/pkg/erx"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

const (
	CREATE_TABLE_POSTS = `CREATE TABLE posts(
		id            	INTEGER PRIMARY KEY AUTOINCREMENT,
		title         	varchar(255) not null,
		text          	text not null,
		date_created	integer not null); `
	CREATE_TABLE_TAGS = `CREATE TABLE tags(
		id            	INTEGER PRIMARY KEY AUTOINCREMENT,
		tag         	type UNIQUE not null);`
	CREATE_TABLE_POST_TAGS = `CREATE TABLE post_tags(
		id            	INTEGER PRIMARY KEY AUTOINCREMENT,
		post_id         integer not null,
		tag_id			integer not null);`
)

func NewSqliteDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func InitBlogDb(path string) error {

	pathDb := path + "/blog.db"

	f, err := os.Stat(pathDb)
	if err != nil && !os.IsNotExist(err) {
		return erx.New(err)
	}
	if f != nil {
		return erx.NewError(605, "Data Base is already exist")
	}
	if os.IsNotExist(err) {
		err = ForceInitBlogDb(path)
	}

	return err
}

func ForceInitBlogDb(path string) error {
	pathDb := path + "/blog.db"

	err := os.MkdirAll(path, 0755)
	if err != nil {
		return erx.New(err)
	}

	_, err = os.Create(pathDb)
	if err != nil {
		return erx.New(err)
	}

	db, err := sql.Open("sqlite3", pathDb)
	if err != nil {
		return erx.New(err)
	}

	defer db.Close()

	_, err = db.Exec(CREATE_TABLE_POSTS)
	if err != nil {
		return erx.New(err)
	}

	_, err = db.Exec(CREATE_TABLE_TAGS)
	if err != nil {
		return erx.New(err)
	}

	_, err = db.Exec(CREATE_TABLE_POST_TAGS)
	if err != nil {
		return erx.New(err)
	}
	return nil
}
