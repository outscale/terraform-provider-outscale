package dl

import "time"

// CreateConnectionInput ...
type CreateConnectionInput struct {
	_              struct{} `type:"structure"`
	Bandwidth      *string  `locationName:"bandwidth" type:"string" required:"true"`
	ConnectionName *string  `locationName:"connectionName" type:"string" required:"true"`
	LagID          *string  `locationName:"lagId" type:"string"`
	Location       *string  `locationName:"location" type:"string" required:"true"`
}

// Connection ...
type Connection struct {
	_               struct{}   `type:"structure"`
	AwsDevice       *string    `locationName:"awsDevice" type:"string"`
	Bandwidth       *string    `locationName:"bandwidth" type:"string"`
	ConnectionID    *string    `locationName:"connectionId" type:"string"`
	ConnectionName  *string    `locationName:"connectionName" type:"string"`
	ConnectionState *string    `locationName:"connectionState" type:"string" enum:"ConnectionState"`
	LagID           *string    `locationName:"lagId" type:"string"`
	LoaIssueTime    *time.Time `locationName:"loaIssueTime" type:"timestamp" timestampFormat:"unix"`
	Location        *string    `locationName:"location" type:"string"`
	OwnerAccount    *string    `locationName:"ownerAccount" type:"string"`
	PartnerName     *string    `locationName:"partnerName" type:"string"`
	Region          *string    `locationName:"region" type:"string"`
	Vlan            *int64     `locationName:"vlan" type:"integer"`
}

// DescribeConnectionsInput ...
type DescribeConnectionsInput struct {
	_            struct{} `type:"structure"`
	ConnectionID *string  `locationName:"connectionId" type:"string"`
}

// Connections ...
type Connections struct {
	_           struct{}      `type:"structure"`
	Connections []*Connection `locationName:"connections" type:"list"`
	RequestID   *string       `locationName:"requestId" type:"string"`
}

// DeleteConnectionInput ...
type DeleteConnectionInput struct {
	_            struct{} `type:"structure"`
	ConnectionID *string  `locationName:"connectionId" type:"string" required:"true"`
}
