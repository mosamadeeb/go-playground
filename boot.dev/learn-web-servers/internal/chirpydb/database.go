package chirpydb

import (
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"os"
	"slices"
	"sync"
)

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

type DBStructure struct {
	IdCount int           `json:"id_count"`
	Chirps  map[int]Chirp `json:"chirps"`
}

type DB struct {
	path string
	mux  *sync.RWMutex
}

// Creates a new database connection and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {
	db := &DB{
		path,
		&sync.RWMutex{},
	}

	if err := db.ensureDB(); err != nil {
		return nil, err
	}

	return db, nil
}

// Creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	chirp := Chirp{
		dbStruct.IdCount,
		body,
	}

	dbStruct.IdCount++
	dbStruct.Chirps[chirp.Id] = chirp

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

	chirps := slices.SortedFunc(maps.Values(dbStruct.Chirps), func(a, b Chirp) int {
		return a.Id - b.Id
	})

	return chirps, nil
}

// Creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	file, err := os.Open(db.path)
	if err == nil {
		file.Close()
		return nil
	} else if errors.Is(err, os.ErrNotExist) {
		data, err := json.Marshal(DBStructure{1, map[int]Chirp{}})
		if err != nil {
			return fmt.Errorf("error marshalling database json: %w", err)
		}

		os.WriteFile(db.path, data, 0o666)

		return nil
	} else {
		return fmt.Errorf("could not create database file: %w", err)
	}
}

// Reads the database file into memory
func (db *DB) loadDB() (DBStructure, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	if err := db.ensureDB(); err != nil {
		return DBStructure{}, err
	}

	data, err := os.ReadFile(db.path)
	if err != nil {
		return DBStructure{}, fmt.Errorf("could read write database file: %w", err)
	}

	var dbStruct DBStructure

	if err := json.Unmarshal(data, &dbStruct); err != nil {
		return DBStructure{}, fmt.Errorf("error unmarshalling database json: %w", err)
	}

	return dbStruct, nil
}

// Writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	if err := db.ensureDB(); err != nil {
		return err
	}

	data, err := json.Marshal(dbStructure)
	if err != nil {
		return fmt.Errorf("error marshalling database json: %w", err)
	}

	if err := os.WriteFile(db.path, data, 0o666); err != nil {
		return fmt.Errorf("could not write database file: %w", err)
	}

	return nil
}
