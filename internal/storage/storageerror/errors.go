package storageerror

import "errors"

var (
	ErrBookAlredyExist = errors.New("book alredy exist")
	ErrEmptyStorage    = errors.New("book storage is empty")
	ErrBookNoFound     = errors.New("book not found")
)
