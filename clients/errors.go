package clients

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// IsUnimplemented returns true if an error indicates that the underlying grpc call
// was unimplemented on the server side.
func IsUnimplemented(err error) bool {
	if err == nil {
		return false
	}
	st := status.Convert(err)
	if st.Code() == codes.Unimplemented {
		return true
	}
	return IsUnimplemented(errors.Unwrap(err))
}
