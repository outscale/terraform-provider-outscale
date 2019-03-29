package dl

import (
	"time"
)

// CreateConnectionInput ...
type CreateConnectionInput struct {
	_              struct{} `type:"structure"`
	Bandwidth      *string  `json:"bandwidth,omitempty" type:"string" required:"true"`
	ConnectionName *string  `json:"connectionName,omitempty" type:"string" required:"true"`
	LagID          *string  `json:"lagId,omitempty" type:"string"`
	Location       *string  `json:"location,omitempty" type:"string" required:"true"`
}

// Connection ...
type Connection struct {
	_               struct{}   `type:"structure"`
	AwsDevice       *string    `json:"awsDevice" type:"string"`
	Bandwidth       *string    `json:"bandwidth" type:"string"`
	ConnectionID    *string    `json:"connectionId" type:"string"`
	ConnectionName  *string    `json:"connectionName" type:"string"`
	ConnectionState *string    `json:"connectionState" type:"string" enum:"ConnectionState"`
	LagID           *string    `json:"lagId" type:"string"`
	LoaIssueTime    *time.Time `json:"loaIssueTime" type:"timestamp" timestampFormat:"unix"`
	Location        *string    `json:"location" type:"string"`
	OwnerAccount    *string    `json:"ownerAccount" type:"string"`
	PartnerName     *string    `json:"partnerName" type:"string"`
	Region          *string    `json:"region" type:"string"`
	Vlan            *int64     `json:"vlan" type:"integer"`
}

// DescribeConnectionsInput ...
type DescribeConnectionsInput struct {
	_            struct{} `type:"structure"`
	ConnectionID *string  `json:"connectionId" type:"string"`
}

// Connections ...
type Connections struct {
	_           struct{}      `type:"structure"`
	Connections []*Connection `json:"connections" type:"list"`
	RequestID   *string       `json:"requestId" type:"string"`
}

// DeleteConnectionInput ...
type DeleteConnectionInput struct {
	_            struct{} `type:"structure"`
	ConnectionID *string  `json:"connectionId" type:"string" required:"true"`
}

// CreatePrivateVirtualInterfaceInput ...
type CreatePrivateVirtualInterfaceInput struct {
	_                          struct{}                    `type:"structure"`
	ConnectionID               *string                     `locationName:"connectionId" type:"string" required:"true"`
	NewPrivateVirtualInterface *NewPrivateVirtualInterface `locationName:"newPrivateVirtualInterface" type:"structure" required:"true"`
}

// NewPrivateVirtualInterface ...
type NewPrivateVirtualInterface struct {
	_                      struct{} `type:"structure"`
	AddressFamily          *string  `locationName:"addressFamily" type:"string" enum:"AddressFamily"`
	AmazonAddress          *string  `locationName:"amazonAddress" type:"string"`
	Asn                    *int64   `locationName:"asn" type:"integer" required:"true"`
	AuthKey                *string  `locationName:"authKey" type:"string"`
	CustomerAddress        *string  `locationName:"customerAddress" type:"string"`
	DirectConnectGatewayID *string  `locationName:"directConnectGatewayId" type:"string"`
	VirtualGatewayID       *string  `locationName:"virtualGatewayId" type:"string"`
	VirtualInterfaceName   *string  `locationName:"virtualInterfaceName" type:"string" required:"true"`
	Vlan                   *int64   `locationName:"vlan" type:"integer" required:"true"`
}

// VirtualInterface ...
type VirtualInterface struct {
	_                      struct{}             `type:"structure"`
	AddressFamily          *string              `locationName:"addressFamily" type:"string" enum:"AddressFamily"`
	AmazonAddress          *string              `locationName:"amazonAddress" type:"string"`
	AmazonSideAsn          *int64               `locationName:"amazonSideAsn" type:"long"`
	Asn                    *int64               `locationName:"asn" type:"integer"`
	AuthKey                *string              `locationName:"authKey" type:"string"`
	BgpPeers               []*BGPPeer           `locationName:"bgpPeers" type:"list"`
	ConnectionID           *string              `locationName:"connectionId" type:"string"`
	CustomerAddress        *string              `locationName:"customerAddress" type:"string"`
	CustomerRouterConfig   *string              `locationName:"customerRouterConfig" type:"string"`
	DirectConnectGatewayID *string              `locationName:"directConnectGatewayId" type:"string"`
	Location               *string              `locationName:"location" type:"string"`
	OwnerAccount           *string              `locationName:"ownerAccount" type:"string"`
	RouteFilterPrefixes    []*RouteFilterPrefix `locationName:"routeFilterPrefixes" type:"list"`
	VirtualGatewayID       *string              `locationName:"virtualGatewayId" type:"string"`
	VirtualInterfaceID     *string              `locationName:"virtualInterfaceId" type:"string"`
	VirtualInterfaceName   *string              `locationName:"virtualInterfaceName" type:"string"`
	VirtualInterfaceState  *string              `locationName:"virtualInterfaceState" type:"string" enum:"VirtualInterfaceState"`
	VirtualInterfaceType   *string              `locationName:"virtualInterfaceType" type:"string"`
	Vlan                   *int64               `locationName:"vlan" type:"integer"`
}

// RouteFilterPrefix ...
type RouteFilterPrefix struct {
	_    struct{} `type:"structure"`
	Cidr *string  `locationName:"cidr" type:"string"`
}

// BGPPeer ...
type BGPPeer struct {
	_               struct{} `type:"structure"`
	AddressFamily   *string  `locationName:"addressFamily" type:"string" enum:"AddressFamily"`
	AmazonAddress   *string  `locationName:"amazonAddress" type:"string"`
	Asn             *int64   `locationName:"asn" type:"integer"`
	AuthKey         *string  `locationName:"authKey" type:"string"`
	BgpPeerState    *string  `locationName:"bgpPeerState" type:"string" enum:"BGPPeerState"`
	BgpStatus       *string  `locationName:"bgpStatus" type:"string" enum:"BGPStatus"`
	CustomerAddress *string  `locationName:"customerAddress" type:"string"`
}

// DescribeVirtualInterfacesInput ...
type DescribeVirtualInterfacesInput struct {
	_                  struct{} `type:"structure"`
	ConnectionID       *string  `locationName:"connectionId" type:"string"`
	VirtualInterfaceID *string  `locationName:"virtualInterfaceId" type:"string"`
}

// DescribeVirtualInterfacesOutput ...
type DescribeVirtualInterfacesOutput struct {
	_                 struct{}            `type:"structure"`
	VirtualInterfaces []*VirtualInterface `locationName:"virtualInterfaces" type:"list"`
	RequestID         *string             `json:"requestId" type:"string"`
}

// DeleteVirtualInterfaceInput ...
type DeleteVirtualInterfaceInput struct {
	_                  struct{} `type:"structure"`
	VirtualInterfaceID *string  `locationName:"virtualInterfaceId" type:"string" required:"true"`
}

// DeleteVirtualInterfaceOutput ...
type DeleteVirtualInterfaceOutput struct {
	_                     struct{} `type:"structure"`
	VirtualInterfaceState *string  `locationName:"virtualInterfaceState" type:"string" enum:"VirtualInterfaceState"`
}

//DescribeLocationsInput ...
type DescribeLocationsInput struct {
	_ struct{} `type:"structure"`
}

// DescribeLocationsOutput ...
type DescribeLocationsOutput struct {
	_         struct{}    `type:"structure"`
	Locations []*Location `locationName:"locations" type:"list"` // A list of colocation hubs where network providers have equipment. Most regions have multiple locations available.
	RequestID *string     `json:"RequestId" type:"string"`
}

//Location ...
type Location struct {
	_            struct{} `type:"structure"`
	LocationCode *string  `locationName:"locationCode" type:"string"` // The code used to indicate the Outscale Direct Connect location.
	LocationName *string  `locationName:"locationName" type:"string"` // The name of the Outscale Direct Connect location. The name includes the colocation partner name and the physical site of the lit building.
}

//DescribeVirtualGatewaysInput ...
type DescribeVirtualGatewaysInput struct {
	_ struct{} `type:"structure"`
}

//DescribeVirtualGatewaysOutput ... A structure containing a list of virtual private gateways.
type DescribeVirtualGatewaysOutput struct {
	_ struct{} `type:"structure"`
	// A list of virtual private gateways.
	VirtualGateways []*VirtualGateway `locationName:"virtualGateways" type:"list"`
	RequestID       *string           `json:"RequestId" type:"string"`
}

// VirtualGateway ...
type VirtualGateway struct {
	_                   struct{} `type:"structure"`
	VirtualGatewayID    *string  `json:"VirtualGatewayId" locationName:"virtualGatewayId" type:"string"`
	VirtualGatewayState *string  `locationName:"virtualGatewayState" type:"string"`
}
