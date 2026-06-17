package model

type HttpError struct {
	StatusCode  int
	ErrorDetail error
}
