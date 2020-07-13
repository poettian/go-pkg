package go_sql_driver

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func QueryMultiRows() error {
	db, err := sql.Open("mysql", "docker:secret@tcp(127.0.0.1:3306)/default")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var (
		id int
		name string
	)

	rows, err := db.Query("select id, username from blog_auth where id = ?", 1)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(id, name)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return err
}

func PrepareQuery() error {
	db, err := sql.Open("mysql", "docker:secret@tcp(127.0.0.1:3306)/default")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var (
		id int
		name string
	)

	stmt, err := db.Prepare("select id, username from blog_auth where id = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(1)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(id, name)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return err
}

func QueryOneRow() error {
	db, err := sql.Open("mysql", "docker:secret@tcp(127.0.0.1:3306)/default")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var name string
	err = db.QueryRow("select username from blog_auth where id = ?", 1).Scan(&name)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(name)

	return err
}

func InsertRow() error {
	db, err := sql.Open("mysql", "docker:secret@tcp(127.0.0.1:3306)/default")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	stmt, err := db.Prepare("insert into blog_auth(username,password) values(?,?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	res, err := stmt.Exec("phper", "phpbest")
	if err != nil {
		log.Fatal(err)
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("ID = %d, affected = %d\n", lastId, rowCnt)

	res, err = stmt.Exec("gopher", "gobest")
	if err != nil {
		log.Fatal(err)
	}
	lastId, err = res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	rowCnt, err = res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("ID = %d, affected = %d\n", lastId, rowCnt)

	return nil
}
