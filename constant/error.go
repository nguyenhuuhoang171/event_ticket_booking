package constant

import "errors"

const (
	INTERNAL_SERVER_ERROR = "Lỗi hệ thống"
)

var ErrMissingID = errors.New("entity id is required")
