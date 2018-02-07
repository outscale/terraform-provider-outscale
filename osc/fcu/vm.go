package fcu

import "github.com/terraform-providers/terraform-provider-outscale/osc"

//VMOperations defines all the operations needed for FCU VMs
type VMOperations struct {
	client *osc.Client
}

//VMService all the necessary actions for them VM service
type VMService interface{}
