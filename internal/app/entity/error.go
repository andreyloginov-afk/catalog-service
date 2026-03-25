package entity

import "errors"

var (
	ErrNotFound            = errors.New("not found")
	ErrAlreadyExists       = errors.New("alrady exists")
	ErrCategoryHasProducts = errors.New("category has linked products")
	ErrIncorrectParametrs  = errors.New("incorrects parameters")
)
