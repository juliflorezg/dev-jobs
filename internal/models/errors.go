package models

import "errors"

var (
	ErrNoRecord                = errors.New("models: no matching record found")
	ErrInvalidCredentials      = errors.New("models: invalid credentials")
	ErrDuplicateEmail          = errors.New("models: duplicate email")
	ErrDuplicateCompanyName    = errors.New("models: duplicate company name")
	ErrDuplicateCompanyWebsite = errors.New("models: duplicate company website")
	ErrNoCompany               = errors.New("models: there is no company for that company_id")
	ErrCouldNotConvertToJSON   = errors.New("models: could not convert value to JSON")
	ErrUseAlreadyApplied       = errors.New("models: user has already applied to the job post")
)
