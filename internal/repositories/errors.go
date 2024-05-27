package repositories

import "errors"

var ErrUserAlreadyExists = errors.New("user already exist")

var ErrInvalidCredentials = errors.New("invalid credantials")

var ErrOrderConflict = errors.New("the order number has already been uploaded by another user")

var ErrOrderAlreadyUpload = errors.New("the order number has already been uploaded by this user")

var ErrInvalidAccural = errors.New("invali processing accural")

var ErrManyRequests = errors.New("to many requests")
