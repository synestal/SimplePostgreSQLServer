package main

import (
	// The sql go library is needed to interact with the database
	"database/sql"
	"fmt"
)

type Store interface {
	CreatePerson(person *Person) error
	GetPerson() ([]*Person, error)
	DeletePerson(idBox int) error
}

type dbStore struct {
	db *sql.DB
}

var store Store

func (store *dbStore) CreatePerson(person *Person) error {
	_, err := store.db.Query(
		"INSERT INTO peopleinfo(name,birthday,occupation) VALUES ($1,$2,$3)",
		person.Name, person.Birthday, person.Occupation)
	return err
}

func (store *dbStore) GetPerson() ([]*Person, error) {
	rows, err := store.db.Query("SELECT id, name, birthday, occupation FROM peopleinfo ORDER BY id ASC ")
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)
	var personList []*Person
	for rows.Next() {
		person := &Person{}
		if err := rows.Scan(&person.Id, &person.Name, &person.Birthday, &person.Occupation); err != nil {
			return nil, err
		}
		personList = append(personList, person)
	}
	return personList, nil
}

func (store *dbStore) DeletePerson(idBox int) error {
	maxIndex, err := store.db.Query("SELECT MAX(id) AS a FROM peopleinfo")
	if err != nil {
		return err
	}
	defer func(maxIndex *sql.Rows) {
		err := maxIndex.Close()
		if err != nil {

		}
	}(maxIndex)
	var maxId int
	for maxIndex.Next() {
		err := maxIndex.Scan(&maxId)
		if err != nil {
			return err
		}
	}
	if maxId < idBox {
		fmt.Println("Index out of range")
		return err
	}
	_, err = store.db.Exec("DELETE FROM peopleinfo WHERE id=$1", idBox)
	return err
}
