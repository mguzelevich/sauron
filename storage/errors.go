package storage

import ()

type Error string

func (e Error) Error() string { return string(e) }

const (
	ErrEntityNotFound     = Error("Not found")
	ErrCollectionNotFound = Error("Collection not found")
)
