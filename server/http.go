package main

type Error struct {
	Message string
}

func MakeError(msg string) *Error {
	err := Error{
		Message: msg,
	}
	return &err
}
