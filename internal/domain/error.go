package domain

import "errors"

var (
	ErrSaveToFileData        = errors.New("err save to file data")
	ErrFileNotFound          = errors.New("file not found")
	ErrShutdownAlreadyCalled = errors.New("shutdown already called")
	ErrInvalidfileID         = errors.New("invalid file ID, must be between 1 and 10")
	ErrFailedReadFile        = errors.New("failed to read file")
)
