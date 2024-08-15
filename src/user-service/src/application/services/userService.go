package services

import (
	"io"

	"github.com/google/uuid"
)

type UserService interface {
	ChangeEmail(userID uuid.UUID, newMail string) error
	ChangeAvatar(userID uuid.UUID, avatar io.Reader) error
}
