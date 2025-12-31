package oapi

// This file exports private functions and constants for testing purposes.
// Items exported here are only available to *_test.go files.

var (
	// Export hasOAPILaunchPermission for testing
	HasOAPILaunchPermission = hasOAPILaunchPermission

	// Export ResourceOutscaleVPNConnectionRouteParseID for testing
	ParseVPNConnectionRouteID = ResourceOutscaleVPNConnectionRouteParseID
)

const (
	// Export test-only constants
	TestAccVmType = testAccVmType
)
