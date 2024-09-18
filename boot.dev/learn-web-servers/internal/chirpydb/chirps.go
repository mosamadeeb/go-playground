package chirpydb

import (
	"maps"
	"slices"
)

type Chirp struct {
	Id       int    `json:"id"`
	Body     string `json:"body"`
	AuthorId int    `json:"author_id"`
}

// Creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string, authorId int) (Chirp, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	chirp := Chirp{
		dbStruct.Chirps.IdCount,
		body,
		authorId,
	}

	dbStruct.Chirps.IdCount++
	dbStruct.Chirps.Items[chirp.Id] = chirp

	if err := db.writeDB(dbStruct); err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

// Returns all chirps in the database, sorted by ID
func (db *DB) GetChirps() ([]Chirp, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return []Chirp{}, err
	}

	return slices.Collect(maps.Values(dbStruct.Chirps.Items)), nil
}

func (db *DB) GetChirp(id int) (Chirp, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	chirp, ok := dbStruct.Chirps.Items[id]
	if !ok {
		return Chirp{}, ErrNotExist
	}

	return chirp, nil
}

func (db *DB) DeleteChirp(id int) error {
	dbStruct, err := db.loadDB()
	if err != nil {
		return err
	}

	_, ok := dbStruct.Chirps.Items[id]
	if !ok {
		return ErrNotExist
	}

	delete(dbStruct.Chirps.Items, id)

	if err := db.writeDB(dbStruct); err != nil {
		return err
	}

	return nil
}
