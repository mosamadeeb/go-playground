// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package database

import (
	"time"
)

type User struct {
	ID        int32
	CreatedAt time.Time
	UpdatedAt time.Time
	Email     string
}
