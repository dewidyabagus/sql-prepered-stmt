package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

const (
	host     = "172.17.0.1"
	port     = 3307
	username = "root"
	password = "root123"
	database = "testdb"
	loc      = "Asia%2FJakarta"
)

var mx = new(sync.Mutex)
var statements = new(sync.Map)

func main() {
	// MySQL Format: user:password@tcp(127.0.0.1:3306)/hello
	// Driver      : https://github.com/go-sql-driver/mysql
	conn, err := sqlx.Connect(
		"mysql", fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4,utf8&loc=%s", username, password, host, port, database, loc,
		),
	)
	if err != nil {
		log.Fatal("Open database connection:", err.Error())
	}
	defer conn.Close()

	if err := conn.Ping(); err != nil {
		log.Fatal("Ping database:", err.Error())
	}

	wg := new(sync.WaitGroup)

	for i := 1; i <= 4; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()

			FindEmpoyeeByID(context.TODO(), conn, uint64(n), nil)
		}(i)
	}

	wg.Wait()
	log.Println("Success ping database")
}

func FindEmpoyeeByID(ctx context.Context, conn *sqlx.DB, id uint64, employee *Empoyee) (err error) {
	mx.Lock()

	stmt, ok := statements.Load("FindEmpoyeeByID")
	if !ok {
		if stmt, err = conn.PrepareContext(ctx, "SELECT * FROM employees WHERE id = ?"); err != nil {
			mx.Unlock()
			return errors.Wrap(err, "Prepare statement")
		}
		statements.Store("FindEmpoyeeByID", stmt)
	}
	mx.Unlock()

	return stmt.(*sql.Stmt).QueryRowContext(ctx, id).Scan(employee)
}
