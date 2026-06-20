package model

type Response struct {
	Error *Error      `json:"error,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}

type Error struct {
	Status  int    `json:"-"`
	Message string `json:"message,omitempty"`
}

func (e *Error) Error() string {
	return e.Message
}

func NewError(status int, message string) *Error {
	return &Error{Status: status, Message: message}
}
