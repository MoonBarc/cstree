package main

import (
	"context"
	"database/sql"
	"log"
)

type TreeDB struct {
	addProgram      *sql.Stmt
	deleteProgram   *sql.Stmt
	setDisabled     *sql.Stmt
	attachLog       *sql.Stmt
	authorsPrograms *sql.Stmt
	allPrograms     *sql.Stmt
}

func prepareOrDie(conn *sql.Conn, stmt string) *sql.Stmt {
	prep, err := conn.PrepareContext(context.Background(), stmt)
	if err != nil {
		log.Fatalln("failed to prepare statement!", err)
	}

	return prep
}

func Prepare(conn *sql.Conn) TreeDB {
	return TreeDB{
		addProgram:      prepareOrDie(conn, "INSERT INTO programs (program, name, author) VALUES (?, ?, ?)"),
		deleteProgram:   prepareOrDie(conn, "DELETE FROM programs WHERE id = ? AND author = ?"),
		setDisabled:     prepareOrDie(conn, "UPDATE programs SET disabled = ? WHERE id = ?"),
		attachLog:       nil, // TODO <--
		authorsPrograms: prepareOrDie(conn, "SELECT id, name, author, disabled FROM programs WHERE author = ?"),
		allPrograms:     prepareOrDie(conn, "SELECT id, name, author, disabled FROM programs"),
	}
}

type Program = struct {
	id                    int64
	name, program, author string
	disabled              bool
}

func (db *TreeDB) AddProgram(program, name, author string) error {
	_, err := db.addProgram.Exec(program, name, author)
	return err
}

func (db *TreeDB) DeleteProgram(id int64, author string) error {
	_, err := db.deleteProgram.Exec(id, author)
	return err
}

func (db *TreeDB) SetDisabled(id int64, disabled bool) error {
	_, err := db.setDisabled.Exec(disabled, id)
	return err
}

func (db *TreeDB) AttachLog() {
	panic("todo")
}

func extractProgram(row *sql.Rows, full bool) Program {
	var id int64
	var name, author, program string
	var disabled bool
	if full {
		row.Scan(&id, &name, &author, &disabled)
	} else {
		row.Scan(&id, &name, &author, &disabled, &program)
	}

	return Program{
		id:       id,
		name:     name,
		author:   author,
		disabled: disabled,
		program:  program,
	}
}

func extractPrograms(rows *sql.Rows) []Program {
	progs := make([]Program, 5)

	for rows.Next() {
		progs = append(progs, extractProgram(rows, false))
	}

	return progs
}

func (db *TreeDB) AuthorsPrograms(author string) ([]Program, error) {
	rows, err := db.authorsPrograms.Query(author)
	if err != nil {
		return nil, err
	}
	return extractPrograms(rows), nil
}

func (db *TreeDB) AllPrograms() ([]Program, error) {
	rows, err := db.allPrograms.Query()
	if err != nil {
		return nil, err
	}

	return extractPrograms(rows), nil
}
