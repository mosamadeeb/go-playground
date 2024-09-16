package chirpydb

import (
	"errors"
)

type User struct {
	Id    int    `json:"id"`
	Email string `json:"email"`

	// We should not be marshalling passwords, but our database is in JSON so we have to lol
	Password string `json:"password"`
}

// Creates a new user and saves it to disk
func (db *DB) CreateUser(email string, password string) (User, error) {
	_, err := db.GetUserByEmail(email)
	if err == nil {
		// Entity already exists
		return User{}, ErrExists
	} else if !errors.Is(err, ErrNotExist) {
		return User{}, err
	}

	dbStruct, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user := User{
		dbStruct.Users.IdCount,
		email,
		password,
	}

	dbStruct.Users.IdCount++
	dbStruct.Users.Items[user.Id] = user

	if err := db.writeDB(dbStruct); err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *DB) GetUserByEmail(email string) (User, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	for _, v := range dbStruct.Users.Items {
		if v.Email == email {
			return v, nil
		}
	}

	return User{}, ErrNotExist
}

func (db *DB) UpdateUser(id int, email, password string) (User, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user, ok := dbStruct.Users.Items[id]
	if !ok {
		return User{}, ErrNotExist
	}

	user.Email = email
	user.Password = password
	dbStruct.Users.Items[id] = user

	if err := db.writeDB(dbStruct); err != nil {
		return User{}, err
	}

	return user, nil
}
