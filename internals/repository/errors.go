package repository

import "errors"

var ErrDuplicatePhoneNumber = errors.New("model:duplicate phone number")
var ErrInvalidCredentials = errors.New("models: invalid credentials")
