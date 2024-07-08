package webgit

import (
	"time"
)

type Note interface {
	GetId() string
	GetBody() string

	UpdateBody(body string) error
	GetCreatedAt() time.Time
}

type MergeRequestInterface interface {
	HasApproved() (bool, error)
	Approve() error
	Unapprove() error

	CreateMergeRequestNote(body string) (Note, error)
	GetMergeRequestNote(noteId string) (Note, error)
	ListMergeRequestNotes() ([]Note, error)
	ListMergeRequestNotesAfter(t time.Time) ([]Note, error)
}
