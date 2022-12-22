package main

type DB struct {
	Host     string
	Password string
	Port     string
	Name     string
}

func CreateConnection(db DB) (DB, error) {

	return db, nil
}
