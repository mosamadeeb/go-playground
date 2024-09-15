package chirpydb

type User struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

// Creates a new user and saves it to disk
func (db *DB) CreateUser(email string) (User, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user := User{
		dbStruct.Users.IdCount,
		email,
	}

	dbStruct.Users.IdCount++
	dbStruct.Users.Items[user.Id] = user

	if err := db.writeDB(dbStruct); err != nil {
		return User{}, err
	}

	return user, nil
}
