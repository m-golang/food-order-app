package main

import "errors"

// ErrInvalidFullName is an error used when a user's full name is invalid.
var ErrInvalidFullName = errors.New("model: Invalid full name")

// ErrInvalidPhoneNumber is an error used when a user's phone number is invalid.
var ErrInvalidPhoneNumber = errors.New("model: Invalid phone number")

// ErrWeakPassword is an error used when a password doesn't meet the required complexity.
var ErrWeakPassword = errors.New("model: password must be at least 8 characters long and contain at least one uppercase letter, one lowercase letter, one digit, and one special character")
