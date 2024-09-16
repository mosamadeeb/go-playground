package chirpydb

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
)

var (
	ErrExists   = errors.New("entity already exists")
	ErrNotExist = errors.New("entity does not exist")
)

type DBMap[T any] struct {
	IdCount int       `json:"id_count"`
	Items   map[int]T `json:"items"`
}

type DBStructure struct {
	Chirps DBMap[Chirp] `json:"chirps"`
	Users  DBMap[User]  `json:"users"`
}

type DB struct {
	path string
	mux  *sync.RWMutex
}

// Creates a new database connection and creates the database file if it doesn't exist
func NewDB(path string, debug bool) (*DB, error) {
	if debug {
		if err := os.Remove(path); err != nil {
			// Log the error and continue
			log.Printf("could not remove test database: %v", err)
		}
	}

	db := &DB{
		path,
		&sync.RWMutex{},
	}

	if err := db.ensureDB(); err != nil {
		return nil, err
	}

	return db, nil
}

// Creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	file, err := os.Open(db.path)
	if err == nil {
		file.Close()
		return nil
	} else if errors.Is(err, os.ErrNotExist) {
		data, err := json.Marshal(DBStructure{
			DBMap[Chirp]{1, map[int]Chirp{}},
			DBMap[User]{1, map[int]User{}},
		})
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
