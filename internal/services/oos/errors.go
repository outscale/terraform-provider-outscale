package oos

import (
	"errors"
	"net/http"

	"github.com/aws/smithy-go"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// Resource errors
var ErrResourceEmpty = errors.New("empty")

// Global errors
const errSetTerraformState = "Unable to reconcile Terraform state from API response"

// Error helpers
func IsNotFound(err error) bool {
	var responseErr *smithyhttp.ResponseError

	return errors.As(err, &responseErr) && responseErr.HTTPStatusCode() == http.StatusNotFound
}

func IsErrorCode(err error, code string) bool {
	apiErr, ok := errors.AsType[smithy.APIError](err)

	return ok && apiErr.ErrorCode() == code
}

func IsErrorMessage(err error, code, message string) bool {
	apiErr, ok := errors.AsType[smithy.APIError](err)

	return ok && apiErr.ErrorCode() == code && apiErr.ErrorMessage() == message
}
