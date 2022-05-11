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
	parent_did text DEFAULT '',
	create_date DATETIME DEFAULT CURRENT_TIMESTAMP,
	primary key (owner_id, dirname, parent_did)
);
-- 文件
create table if not exists files(
	fid text,
	owner_id text not null,
	filename text not null,
	file_size integer not null,
	parent_did text DEFAULT '',
	file_real_path text not null,
	create_date DATETIME DEFAULT CURRENT_TIMESTAMP,
	primary key (owner_id, filename, parent_did)
);

-- 管理员账号
create table if not exists admin(
	uid text primary key
);

-- 用户资源组，同时会作为用户进入资源组页面的顶层文件夹
create table if not exists user_group(
	gid varchar(10) NOT NULL, -- 资源组id
	name text NOT NULL, -- 资源组名称
	user_ids text DEFAULT '', -- 该资源组包含用户id，以分号分割
	create_date DATETIME DEFAULT CURRENT_TIMESTAMP,
	group_type INTEGER DEFAULT 0 not null, -- 资源组类型；0 为公共可访问，1为user_ids可访问，其他未定义
	primary key (name)
);

-- 用户资源组文件
create table if not exists user_group_resource(
	gid varchar(10) NOT NULL, -- 所属资源组id
	rid varchar(10) NOT NULL, -- 资源id
	fid text DEFAULT '', -- 实际文件id，如果是文件的话
	did text DEFAULT '', -- 实际目录id，如果是目录的话
	name text NOT NULL, -- 资源名称
	parent_did text DEFAULT '', -- 父目录，顶层时，为空
	rtype integer NOT NULL, -- 资源类型；1是文件夹；2是文件
	uid text NOT NULL, -- 拥有者
	create_date DATETIME DEFAULT CURRENT_TIMESTAMP,
	expire_date DATETIME, -- 过期时间，需要加一个定时任务
	primary key (gid, parent_did, name)
);
	`)
	shouldInsertDefaultAdmin()
	shouldInsertDefaultGroup()
	utils.CheckErr(err)
}
func shouldInsertDefaultGroup() {
	rows, err := DB.Query("select * from user_group")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	hasRow := false
	for rows.Next() {
		hasRow = true
		break
	}
	if !hasRow { // 注入默认
		CreateGroup("默认", GroupTypeCommon)
	}
}

func shouldInsertDefaultAdmin() {
	rows, err := DB.Query("select * from admin")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	hasRow := false
	for rows.Next() {
		hasRow = true
		break
	}
	if !hasRow { // 注入默认
		uid := utils.GenerateUid()
		username := "admin"
		password := utils.GeneratePassword()
		stmt, _ := DB.Prepare("insert into admin (uid) values(?)")
		defer stmt.Close()
		stmt.Exec(uid)
		AddUserWithId(uid, username, password)
		log.Print("插入默认管理员账号：")
		log.Print("\n账号：", username)
		log.Print("\n密码：", password)
	}
}
