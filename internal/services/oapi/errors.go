package oapi

import "errors"

// Resource errors
var (
	ErrResourceEmpty = errors.New("empty")

	ErrResourceInvalidIOPS = errors.New(`The "iops" parameter can only be set when creating an "io1" volume.
Check Outscale API documentation for more details:
https://docs.outscale.com/en/userguide/About-Volumes.html#_volume_types_and_iops`)

	ErrResourceInvalidEIMPath = errors.New(`The path name must begin and end with a slash (/), and contain between 1 and 512 alphanumeric characters and/or slashes (/), or underscores (_).`)
)

// Data source query errors
var (
	ErrNoResults = errors.New("your query returned no results, change your search criteria and try again")

	ErrMultipleResults = errors.New("your query returned multiple results, use more specific search criteria")

	ErrFilterRequired = errors.New("filters must be assigned")
)
