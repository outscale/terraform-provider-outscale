package oks

import "errors"

// Resource errors
var ErrResourceEmpty = errors.New("empty")

// Global errors
const errSetTerraformState = "Unable to reconcile Terraform state from API response"
