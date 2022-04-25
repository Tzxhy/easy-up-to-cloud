package models

import (
	"database/sql"
	"log"
	"os"
	"path"

	"gitee.com/tzxhy/web/utils"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitSqlite3() {
	if DB != nil {
		return
	}
	dir, _ := os.Getwd()
	_DB, err := sql.Open("sqlite3", path.Join(dir, "test.db"))
	if err != nil {
		log.Fatal(err)
	}
	DB = _DB
	InitTables()
}

func InitTables() {
	_, err := DB.Exec(`
-- 用户
create table if not exists users(
	uid text primary key,
	name varchar(64) not null,
	password varchar(64) not null,
	create_date DATETIME DEFAULT CURRENT_TIMESTAMP
);
-- 目录
create table if not exists dirs(
	did text,
	owner_id text not null,
	dirname text not null,
	parent_did text,
	create_date DATETIME DEFAULT CURRENT_TIMESTAMP,
	primary key (owner_id, dirname, parent_did)
);
-- 文件
create table if not exists files(
	fid text,
	owner_id text not null,
	filename text not null,
	file_size integer not null,
	parent_did text,
	file_real_path text not null,
	create_date DATETIME DEFAULT CURRENT_TIMESTAMP,
	primary key (owner_id, filename, parent_did)
);
	`)
	utils.CheckErr(err)
}
