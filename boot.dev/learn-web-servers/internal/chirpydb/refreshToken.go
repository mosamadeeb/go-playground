package chirpydb

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

var ErrExpired = errors.New("token expired")

type RefreshToken struct {
	Token     string    `json:"token"`
	UserId    int       `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (db *DB) AddRefreshToken(userId int, expiresAt time.Time) (RefreshToken, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return RefreshToken{}, err
	}

	randBytes := make([]byte, 32)
	if _, err := rand.Read(randBytes); err != nil {
		return RefreshToken{}, err
	}

	tokenString := hex.EncodeToString(randBytes)
	refreshToken := RefreshToken{tokenString, userId, expiresAt}

	dbStruct.RefreshTokens[tokenString] = refreshToken

	if err := db.writeDB(dbStruct); err != nil {
		return RefreshToken{}, err
	}

	return refreshToken, nil
}

func (db *DB) RevokeRefreshToken(tokenString string) error {
	dbStruct, err := db.loadDB()
	if err != nil {
		return err
	}

	_, ok := dbStruct.RefreshTokens[tokenString]
	if !ok {
		return ErrNotExist
	}

	delete(dbStruct.RefreshTokens, tokenString)

	if err := db.writeDB(dbStruct); err != nil {
		return err
	}

	return nil
}

func (db *DB) CheckRefreshToken(tokenString string) (userId int, err error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return 0, err
	}

	refreshToken, ok := dbStruct.RefreshTokens[tokenString]
	if !ok {
		return 0, ErrNotExist
	}

	if time.Now().UTC().After(refreshToken.ExpiresAt.UTC()) {
		delete(dbStruct.RefreshTokens, tokenString)

		if err := db.writeDB(dbStruct); err != nil {
			return 0, fmt.Errorf("error deleting expired token: %w", err)
		}

		return 0, ErrExpired
	}

	return refreshToken.UserId, nil
}
