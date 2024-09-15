package chirpydb

import (
	"maps"
	"slices"
)

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

// Creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	chirp := Chirp{
		dbStruct.Chirps.IdCount,
		body,
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

	chirps := slices.SortedFunc(maps.Values(dbStruct.Chirps.Items), func(a, b Chirp) int {
		return a.Id - b.Id
	})

	return chirps, nil
}
