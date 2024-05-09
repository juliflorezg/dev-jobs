package models

import "errors"

var (
	ErrNoRecord                = errors.New("models: no matching record found")
	ErrInvalidCredentials      = errors.New("models: invalid credentials")
	ErrDuplicateEmail          = errors.New("models: duplicate email")
	ErrDuplicateCompanyName    = errors.New("models: duplicate company name")
	ErrDuplicateCompanyWebsite = errors.New("models: duplicate company website")
)
