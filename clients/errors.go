package clients

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// IsUnimplemented returns true if an error indicates that the underlying grpc call
// was unimplemented on the server side.
func IsUnimplemented(err error) bool {
	st := status.Convert(err)
	return st != nil && st.Code() == codes.Unimplemented
}
