package main

import (
	"context"
	"database/sql"
	"io/fs"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

func SetUserVersion(conn *sql.Conn, v int64) {
	conn.ExecContext(context.Background(), "PRAGMA user_version = "+strconv.Itoa(int(v)))
}

func GetUserVersion(conn *sql.Conn) int64 {
	rows, err := conn.QueryContext(context.Background(), "PRAGMA user_version")
	if err != nil {
		log.Fatalln("failed to get user_version", err)
	}
	var ver int64
	rows.Next()
	err = rows.Scan(&ver)
	if err != nil {
		log.Fatalln("failed to read user_version", err)
	}
	return ver
}

func num(a fs.DirEntry) int {
	i, err := strconv.Atoi(strings.Split(a.Name(), "-")[0])
	if err != nil {
		log.Fatalln("invalid migration found", a.Name())
	}
	return i
}

func RunMigrations(conn *sql.Conn) {
	GetUserVersion(conn)

	migs, err := fs.ReadDir(os.DirFS("./migrations"), ".")
	if err != nil {
		log.Fatalln("failed to find any migrations", err)
	}

	slices.SortFunc(migs, func(a, b fs.DirEntry) int {
		return num(a) - num(b)
	})

	max := num(migs[len(migs)-1])

	for i := GetUserVersion(conn) + 1; i <= int64(max); i++ {
		mig := migs[i-1]
		log.Printf("migrating %v --> %v using %v", i-1, i, mig.Name())

		mig_sql, err := fs.ReadFile(os.DirFS("./migrations"), mig.Name())
		if err != nil {
			log.Fatalln("failed to read migration", err)
		}

		_, err = conn.ExecContext(context.Background(), string(mig_sql))
		if err != nil {
			log.Fatalf("error in %v migration: %v", mig.Name(), err)
		}

		SetUserVersion(conn, i)
	}

	log.Println("done migrations")
}
