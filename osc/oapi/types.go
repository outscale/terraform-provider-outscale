// GENERATED FILE: DO NOT EDIT!

package oapi

// Types used by the API.
// implements the service definition of AcceptNetPeeringRequest
type AcceptNetPeeringRequest struct {
	DryRun       bool   `json:"DryRun,omitempty"`
	NetPeeringId string `json:"NetPeeringId,omitempty"`
}

// implements the service definition of AcceptNetPeeringResponse
type AcceptNetPeeringResponse struct {
	NetPeering      NetPeering      `json:"NetPeering,omitempty"`
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of AccepterNet
type AccepterNet struct {
	AccountId string   `json:"AccountId,omitempty"`
	IpRanges  []string `json:"IpRanges,omitempty"`
	NetId     string   `json:"NetId,omitempty"`
}

// implements the service definition of Addition
type Addition struct {
	AccountIds       []string `json:"AccountIds,omitempty"`
	GlobalPermission bool     `json:"GlobalPermission,omitempty"`
}

// implements the service definition of CancelExportTaskRequest
type CancelExportTaskRequest struct {
	DryRun       bool   `json:"DryRun,omitempty"`
	ExportTaskId string `json:"ExportTaskId,omitempty"`
}

// implements the service definition of CancelExportTaskResponse
type CancelExportTaskResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of CopyImageRequest
type CopyImageRequest struct {
	ClientToken      string `json:"ClientToken,omitempty"`
	Description      string `json:"Description,omitempty"`
	DryRun           bool   `json:"DryRun,omitempty"`
	Name             string `json:"Name,omitempty"`
	SourceImageId    string `json:"SourceImageId,omitempty"`
	SourceRegionName string `json:"SourceRegionName,omitempty"`
}

// implements the service definition of CopyImageResponse
type CopyImageResponse struct {
	Image           Image           `json:"Image,omitempty"`
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of CopyImage_BlockDeviceMappings
type CopyImage_BlockDeviceMappings struct {
	Bsu               CopyImage_Bsu `json:"Bsu,omitempty"`
	DeviceName        string        `json:"DeviceName,omitempty"`
	VirtualDeviceName string        `json:"VirtualDeviceName,omitempty"`
}

// implements the service definition of CopyImage_Bsu
type CopyImage_Bsu struct {
	DeleteOnVmDeletion bool   `json:"DeleteOnVmDeletion,omitempty"`
	Iops               int64  `json:"Iops,omitempty"`
	SnapshotId         string `json:"SnapshotId,omitempty"`
	VolumeSize         int64  `json:"VolumeSize,omitempty"`
	VolumeType         string `json:"VolumeType,omitempty"`
}

// implements the service definition of CopyImage_PermissionToLaunch
type CopyImage_PermissionToLaunch struct {
	AccountIds       []string `json:"AccountIds,omitempty"`
	GlobalPermission bool     `json:"GlobalPermission,omitempty"`
}

// implements the service definition of CopySnapshotRequest
type CopySnapshotRequest struct {
	Description      string `json:"Description,omitempty"`
	DryRun           bool   `json:"DryRun,omitempty"`
	SourceRegionName string `json:"SourceRegionName,omitempty"`
	SourceSnapshotId string `json:"SourceSnapshotId,omitempty"`
}

// implements the service definition of CopySnapshotResponse
type CopySnapshotResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
	Snapshot        Snapshot        `json:"Snapshot,omitempty"`
}

// implements the service definition of CopySnapshot_PermissionToCreateVolume
type CopySnapshot_PermissionToCreateVolume struct {
	AccountIds       []string `json:"AccountIds,omitempty"`
	GlobalPermission bool     `json:"GlobalPermission,omitempty"`
}

// implements the service definition of CreateFirewallRuleInboundRequest
type CreateFirewallRuleInboundRequest struct {
	DryRun                          bool           `json:"DryRun,omitempty"`
	FirewallRulesSetId              string         `json:"FirewallRulesSetId,omitempty"`
	FromPortRange                   int64          `json:"FromPortRange,omitempty"`
	InboundRules                    []InboundRules `json:"InboundRules,omitempty"`
	IpProtocol                      string         `json:"IpProtocol,omitempty"`
	IpRange                         string         `json:"IpRange,omitempty"`
	SourceFirewallRulesSetAccountId string         `json:"SourceFirewallRulesSetAccountId,omitempty"`
	SourceFirewallRulesSetName      string         `json:"SourceFirewallRulesSetName,omitempty"`
	ToPortRange                     int64          `json:"ToPortRange,omitempty"`
}

// implements the service definition of CreateFirewallRuleInboundResponse
type CreateFirewallRuleInboundResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of CreateFirewallRuleOutboundRequest
type CreateFirewallRuleOutboundRequest struct {
	DestinationFirewallRulesSetAccountId string          `json:"DestinationFirewallRulesSetAccountId,omitempty"`
	DestinationFirewallRulesSetName      string          `json:"DestinationFirewallRulesSetName,omitempty"`
	DryRun                               bool            `json:"DryRun,omitempty"`
	FirewallRulesSetId                   string          `json:"FirewallRulesSetId,omitempty"`
	FromPortRange                        int64           `json:"FromPortRange,omitempty"`
	IpProtocol                           string          `json:"IpProtocol,omitempty"`
	IpRange                              string          `json:"IpRange,omitempty"`
	OutboundRules                        []OutboundRules `json:"OutboundRules,omitempty"`
	ToPortRange                          int64           `json:"ToPortRange,omitempty"`
}

// implements the service definition of CreateFirewallRuleOutboundResponse
type CreateFirewallRuleOutboundResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of CreateFirewallRulesSetRequest
type CreateFirewallRulesSetRequest struct {
	Description string `json:"Description,omitempty"`
	DryRun      bool   `json:"DryRun,omitempty"`
	Name        string `json:"Name,omitempty"`
	NetId       string `json:"NetId,omitempty"`
}

// implements the service definition of CreateFirewallRulesSetResponse
type CreateFirewallRulesSetResponse struct {
	FirewallRulesSetId string          `json:"FirewallRulesSetId,omitempty"`
	ResponseContext    ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of CreateImageExportTaskRequest
type CreateImageExportTaskRequest struct {
	DryRun    bool      `json:"DryRun,omitempty"`
	ImageId   string    `json:"ImageId,omitempty"`
	OsuExport OsuExport `json:"OsuExport,omitempty"`
}

// implements the service definition of CreateImageExportTaskResponse
type CreateImageExportTaskResponse struct {
	ImageExportTask ImageExportTask `json:"ImageExportTask,omitempty"`
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of CreateImageRequest
type CreateImageRequest struct {
	Description string `json:"Description,omitempty"`
	DryRun      bool   `json:"DryRun,omitempty"`
	Name        string `json:"Name,omitempty"`
	NoReboot    bool   `json:"NoReboot,omitempty"`
	VmId        string `json:"VmId,omitempty"`
}

// implements the service definition of CreateImageResponse
type CreateImageResponse struct {
	Image           Image           `json:"Image,omitempty"`
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of CreateKeypairRequest
type CreateKeypairRequest struct {
	DryRun      bool   `json:"DryRun,omitempty"`
	KeypairName string `json:"KeypairName,omitempty"`
}

// implements the service definition of CreateKeypairResponse
type CreateKeypairResponse struct {
	Keypair         Keypair         `json:"Keypair,omitempty"`
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of CreateNatServiceRequest
type CreateNatServiceRequest struct {
	ClientToken string `json:"ClientToken,omitempty"`
	DryRun      bool   `json:"DryRun,omitempty"`
	LinkId      string `json:"LinkId,omitempty"`
	SubnetId    string `json:"SubnetId,omitempty"`
}

// implements the service definition of CreateNatServiceResponse
type CreateNatServiceResponse struct {
	ClientToken     string          `json:"ClientToken,omitempty"`
	NatService      NatService      `json:"NatService,omitempty"`
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of CreateNatService_PublicIps
type CreateNatService_PublicIps struct {
	LinkId   string `json:"LinkId,omitempty"`
	PublicIp string `json:"PublicIp,omitempty"`
}

// implements the service definition of CreateNetInternetGatewayRequest
type CreateNetInternetGatewayRequest struct {
	DryRun bool `json:"DryRun,omitempty"`
}

// implements the service definition of CreateNetInternetGatewayResponse
type CreateNetInternetGatewayResponse struct {
	NetInternetGateway NetInternetGateway `json:"NetInternetGateway,omitempty"`
	ResponseContext    ResponseContext    `json:"ResponseContext,omitempty"`
}

// implements the service definition of CreateNetPeeringRequest
type CreateNetPeeringRequest struct {
	AccepterNetId string `json:"AccepterNetId,omitempty"`
	DryRun        bool   `json:"DryRun,omitempty"`
	SourceNetId   string `json:"SourceNetId,omitempty"`
}

// implements the service definition of CreateNetPeeringResponse
type CreateNetPeeringResponse struct {
	NetPeering      NetPeering      `json:"NetPeering,omitempty"`
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of CreateNetRequest
type CreateNetRequest struct {
	DryRun  bool   `json:"DryRun,omitempty"`
	IpRange string `json:"IpRange,omitempty"`
	Tenancy string `json:"Tenancy,omitempty"`
}

// implements the service definition of CreateNetResponse
type CreateNetResponse struct {
	Net             Net             `json:"Net,omitempty"`
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of CreateNicRequest
type CreateNicRequest struct {
	Description         string                         `json:"Description,omitempty"`
	DryRun              bool                           `json:"DryRun,omitempty"`
	FirewallRulesSetIds []string                       `json:"FirewallRulesSetIds,omitempty"`
	PrivateIps          []CreateNic_Request_PrivateIps `json:"PrivateIps,omitempty"`
	SubnetId            string                         `json:"SubnetId,omitempty"`
}

// implements the service definition of CreateNicResponse
type CreateNicResponse struct {
	Nic             Nic             `json:"Nic,omitempty"`
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of CreateNic_FirewallRulesSets
type CreateNic_FirewallRulesSets struct {
	FirewallRulesSetId   string `json:"FirewallRulesSetId,omitempty"`
	FirewallRulesSetName string `json:"FirewallRulesSetName,omitempty"`
}

// implements the service definition of CreateNic_NicLink
type CreateNic_NicLink struct {
	DeleteOnVmDeletion bool   `json:"DeleteOnVmDeletion,omitempty"`
	DeviceNumber       int64  `json:"DeviceNumber,omitempty"`
	NicLinkId          string `json:"NicLinkId,omitempty"`
	State              string `json:"State,omitempty"`
	VmAccountId        string `json:"VmAccountId,omitempty"`
	VmId               string `json:"VmId,omitempty"`
}

// implements the service definition of CreateNic_PublicIpToNicLink
type CreateNic_PublicIpToNicLink struct {
	LinkId            string `json:"LinkId,omitempty"`
	PublicDnsName     string `json:"PublicDnsName,omitempty"`
	PublicIp          string `json:"PublicIp,omitempty"`
	PublicIpAccountId string `json:"PublicIpAccountId,omitempty"`
	ReservationId     string `json:"ReservationId,omitempty"`
}

// implements the service definition of CreateNic_Request_PrivateIps
type CreateNic_Request_PrivateIps struct {
	IsPrimary bool   `json:"IsPrimary,omitempty"`
	PrivateIp string `json:"PrivateIp,omitempty"`
}

// implements the service definition of CreateNic_ResponseContext
type ResponseContext struct {
	RequestId   string `json:"RequestId,omitempty"`
	RequesterId string `json:"RequesterId,omitempty"`
}

// implements the service definition of CreateNic_Response_PrivateIps
type CreateNic_Response_PrivateIps struct {
	IsPrimary         bool                        `json:"IsPrimary,omitempty"`
	PrivateDnsName    string                      `json:"PrivateDnsName,omitempty"`
	PrivateIp         string                      `json:"PrivateIp,omitempty"`
	PublicIpToNicLink CreateNic_PublicIpToNicLink `json:"PublicIpToNicLink,omitempty"`
}

// implements the service definition of CreatePublicIpRequest
type CreatePublicIpRequest struct {
	DryRun    bool   `json:"DryRun,omitempty"`
	Placement string `json:"Placement,omitempty"`
}

// implements the service definition of CreatePublicIpResponse
type CreatePublicIpResponse struct {
	Placement       string          `json:"Placement,omitempty"`
	PublicIp        string          `json:"PublicIp,omitempty"`
	ReservationId   string          `json:"ReservationId,omitempty"`
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of CreateRouteRequest
type CreateRouteRequest struct {
	DestinationIpRange string `json:"DestinationIpRange,omitempty"`
	DryRun             bool   `json:"DryRun,omitempty"`
	GatewayId          string `json:"GatewayId,omitempty"`
	NatServiceId       string `json:"NatServiceId,omitempty"`
	NetPeeringId       string `json:"NetPeeringId,omitempty"`
	NicId              string `json:"NicId,omitempty"`
	RouteTableId       string `json:"RouteTableId,omitempty"`
	VmId               string `json:"VmId,omitempty"`
}

// implements the service definition of CreateRouteResponse
type CreateRouteResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of CreateRouteTableRequest
type CreateRouteTableRequest struct {
	DryRun bool   `json:"DryRun,omitempty"`
	NetId  string `json:"NetId,omitempty"`
}

// implements the service definition of CreateRouteTableResponse
type CreateRouteTableResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
	RouteTable      RouteTable      `json:"RouteTable,omitempty"`
}

// implements the service definition of CreateSnapshotExportTaskRequest
type CreateSnapshotExportTaskRequest struct {
	DryRun     bool      `json:"DryRun,omitempty"`
	OsuExport  OsuExport `json:"OsuExport,omitempty"`
	SnapshotId string    `json:"SnapshotId,omitempty"`
}

// implements the service definition of CreateSnapshotExportTaskResponse
type CreateSnapshotExportTaskResponse struct {
	ResponseContext    ResponseContext    `json:"ResponseContext,omitempty"`
	SnapshotExportTask SnapshotExportTask `json:"SnapshotExportTask,omitempty"`
}

// implements the service definition of CreateSnapshotRequest
type CreateSnapshotRequest struct {
	Description string `json:"Description,omitempty"`
	DryRun      bool   `json:"DryRun,omitempty"`
	VolumeId    string `json:"VolumeId,omitempty"`
}

// implements the service definition of CreateSnapshotResponse
type CreateSnapshotResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
	Snapshot        Snapshot        `json:"Snapshot,omitempty"`
}

// implements the service definition of CreateSubnetRequest
type CreateSubnetRequest struct {
	DryRun        bool   `json:"DryRun,omitempty"`
	IpRange       string `json:"IpRange,omitempty"`
	NetId         string `json:"NetId,omitempty"`
	SubRegionName string `json:"SubRegionName,omitempty"`
}

// implements the service definition of CreateSubnetResponse
type CreateSubnetResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
	Subnet          Subnet          `json:"Subnet,omitempty"`
}

// implements the service definition of CreateTagsRequest
type CreateTagsRequest struct {
	DryRun      bool     `json:"DryRun,omitempty"`
	ResourceIds []string `json:"ResourceIds,omitempty"`
	Tags        []Tags   `json:"Tags,omitempty"`
}

// implements the service definition of CreateTagsResponse
type CreateTagsResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of CreateVmsRequest
type CreateVmsRequest struct {
	BlockDeviceMappings         []CreateVms_Request_BlockDeviceMappings `json:"BlockDeviceMappings,omitempty"`
	BsuOptimized                bool                                    `json:"BsuOptimized,omitempty"`
	ClientToken                 string                                  `json:"ClientToken,omitempty"`
	DeletionProtection          bool                                    `json:"DeletionProtection,omitempty"`
	DryRun                      bool                                    `json:"DryRun,omitempty"`
	FirewallRulesSetIds         []string                                `json:"FirewallRulesSetIds,omitempty"`
	FirewallRulesSets           []string                                `json:"FirewallRulesSets,omitempty"`
	ImageId                     string                                  `json:"ImageId,omitempty"`
	KeypairName                 string                                  `json:"KeypairName,omitempty"`
	MaxVmsCount                 int64                                   `json:"MaxVmsCount,omitempty"`
	MinVmsCount                 int64                                   `json:"MinVmsCount,omitempty"`
	Nics                        []CreateVms_Request_Nics                `json:"Nics,omitempty"`
	Placement                   CreateVms_Request_Placement             `json:"Placement,omitempty"`
	PrivateIps                  []string                                `json:"PrivateIps,omitempty"`
	SubnetId                    string                                  `json:"SubnetId,omitempty"`
	Type                        string                                  `json:"Type,omitempty"`
	UserData                    string                                  `json:"UserData,omitempty"`
	VmInitiatedShutdownBehavior string                                  `json:"VmInitiatedShutdownBehavior,omitempty"`
}

// implements the service definition of CreateVmsResponse
type CreateVmsResponse struct {
	AccountId         string                        `json:"AccountId,omitempty"`
	FirewallRulesSets []CreateVms_FirewallRulesSets `json:"FirewallRulesSets,omitempty"`
	ReservationId     string                        `json:"ReservationId,omitempty"`
	ResponseContext   ResponseContext               `json:"ResponseContext,omitempty"`
	Vms               []CreateVms_Vms               `json:"Vms,omitempty"`
}

// implements the service definition of CreateVms_Bsu
type CreateVms_Bsu struct {
	DeleteOnVmDeletion bool   `json:"DeleteOnVmDeletion,omitempty"`
	State              string `json:"State,omitempty"`
	VolumeId           string `json:"VolumeId,omitempty"`
}

// implements the service definition of CreateVms_FirewallRulesSets
type CreateVms_FirewallRulesSets struct {
	FirewallRulesSetId   string `json:"FirewallRulesSetId,omitempty"`
	FirewallRulesSetName string `json:"FirewallRulesSetName,omitempty"`
}

// implements the service definition of CreateVms_NicLink
type CreateVms_NicLink struct {
	DeleteOnVmDeletion bool   `json:"DeleteOnVmDeletion,omitempty"`
	DeviceNumber       int64  `json:"DeviceNumber,omitempty"`
	NicLinkId          string `json:"NicLinkId,omitempty"`
	State              string `json:"State,omitempty"`
}

// implements the service definition of CreateVms_PublicIpToNicLink
type CreateVms_PublicIpToNicLink struct {
	PublicDnsName     string `json:"PublicDnsName,omitempty"`
	PublicIp          string `json:"PublicIp,omitempty"`
	PublicIpAccountId string `json:"PublicIpAccountId,omitempty"`
}

// implements the service definition of CreateVms_Request_BlockDeviceMappings
type CreateVms_Request_BlockDeviceMappings struct {
	DeviceName string `json:"DeviceName,omitempty"`
}

// implements the service definition of CreateVms_Request_Nics
type CreateVms_Request_Nics struct {
	DeleteOnVmDeletion      bool                           `json:"DeleteOnVmDeletion,omitempty"`
	Description             string                         `json:"Description,omitempty"`
	DeviceNumber            int64                          `json:"DeviceNumber,omitempty"`
	FirewallRulesSetIds     []string                       `json:"FirewallRulesSetIds,omitempty"`
	NicId                   string                         `json:"NicId,omitempty"`
	PrivateIps              []CreateVms_Request_PrivateIps `json:"PrivateIps,omitempty"`
	SecondaryPrivateIpCount int64                          `json:"SecondaryPrivateIpCount,omitempty"`
	SubnetId                string                         `json:"SubnetId,omitempty"`
}

// implements the service definition of CreateVms_Request_Placement
type CreateVms_Request_Placement struct {
	SubRegionName string `json:"SubRegionName,omitempty"`
	Tenancy       string `json:"Tenancy,omitempty"`
}

// implements the service definition of CreateVms_Request_PrivateIps
type CreateVms_Request_PrivateIps struct {
	IsPrimary bool   `json:"IsPrimary,omitempty"`
	PrivateIp string `json:"PrivateIp,omitempty"`
}

// implements the service definition of CreateVms_Response_BlockDeviceMappings
type CreateVms_Response_BlockDeviceMappings struct {
	Bsu        CreateVms_Bsu `json:"Bsu,omitempty"`
	DeviceName string        `json:"DeviceName,omitempty"`
}

// implements the service definition of CreateVms_Response_Nics
type CreateVms_Response_Nics struct {
	AccountId           string                          `json:"AccountId,omitempty"`
	Description         string                          `json:"Description,omitempty"`
	FirewallRulesSets   []CreateVms_FirewallRulesSets   `json:"FirewallRulesSets,omitempty"`
	IsSourceDestChecked bool                            `json:"IsSourceDestChecked,omitempty"`
	MacAddress          string                          `json:"MacAddress,omitempty"`
	NetId               string                          `json:"NetId,omitempty"`
	NicId               string                          `json:"NicId,omitempty"`
	NicLink             CreateVms_NicLink               `json:"NicLink,omitempty"`
	PrivateDnsName      string                          `json:"PrivateDnsName,omitempty"`
	PrivateIps          []CreateVms_Response_PrivateIps `json:"PrivateIps,omitempty"`
	PublicIpToNicLink   CreateVms_PublicIpToNicLink     `json:"PublicIpToNicLink,omitempty"`
	State               string                          `json:"State,omitempty"`
	SubnetId            string                          `json:"SubnetId,omitempty"`
}

// implements the service definition of CreateVms_Response_Placement
type CreateVms_Response_Placement struct {
	PlacementName string `json:"PlacementName,omitempty"`
	SubRegionName string `json:"SubRegionName,omitempty"`
	Tenancy       string `json:"Tenancy,omitempty"`
}

// implements the service definition of CreateVms_Response_PrivateIps
type CreateVms_Response_PrivateIps struct {
	IsPrimary         bool                        `json:"IsPrimary,omitempty"`
	PrivateDnsName    string                      `json:"PrivateDnsName,omitempty"`
	PrivateIp         string                      `json:"PrivateIp,omitempty"`
	PublicIpToNicLink CreateVms_PublicIpToNicLink `json:"PublicIpToNicLink,omitempty"`
}

// implements the service definition of CreateVms_Vms
type CreateVms_Vms struct {
	Architecture        string                                   `json:"Architecture,omitempty"`
	BlockDeviceMappings []CreateVms_Response_BlockDeviceMappings `json:"BlockDeviceMappings,omitempty"`
	BsuOptimized        bool                                     `json:"BsuOptimized,omitempty"`
	ClientToken         string                                   `json:"ClientToken,omitempty"`
	FirewallRulesSets   []CreateVms_FirewallRulesSets            `json:"FirewallRulesSets,omitempty"`
	Hypervisor          string                                   `json:"Hypervisor,omitempty"`
	ImageId             string                                   `json:"ImageId,omitempty"`
	IsSourceDestChecked bool                                     `json:"IsSourceDestChecked,omitempty"`
	KeypairName         string                                   `json:"KeypairName,omitempty"`
	LaunchNumber        int64                                    `json:"LaunchNumber,omitempty"`
	NetId               string                                   `json:"NetId,omitempty"`
	Nics                []CreateVms_Response_Nics                `json:"Nics,omitempty"`
	Placement           CreateVms_Response_Placement             `json:"Placement,omitempty"`
	PrivateDnsName      string                                   `json:"PrivateDnsName,omitempty"`
	PrivateIp           string                                   `json:"PrivateIp,omitempty"`
	ProductCodes        []string                                 `json:"ProductCodes,omitempty"`
	PublicDnsName       string                                   `json:"PublicDnsName,omitempty"`
	PublicIp            string                                   `json:"PublicIp,omitempty"`
	RootDeviceName      string                                   `json:"RootDeviceName,omitempty"`
	RootDeviceType      string                                   `json:"RootDeviceType,omitempty"`
	State               string                                   `json:"State,omitempty"`
	SubnetId            string                                   `json:"SubnetId,omitempty"`
	Tags                []Tags                                   `json:"Tags,omitempty"`
	Transition          string                                   `json:"Transition,omitempty"`
	Type                string                                   `json:"Type,omitempty"`
	VmId                string                                   `json:"VmId,omitempty"`
}

// implements the service definition of CreateVolumeRequest
type CreateVolumeRequest struct {
	DryRun        bool   `json:"DryRun,omitempty"`
	Iops          int64  `json:"Iops,omitempty"`
	Size          int64  `json:"Size,omitempty"`
	SnapshotId    string `json:"SnapshotId,omitempty"`
	SubRegionName string `json:"SubRegionName,omitempty"`
	Type          string `json:"Type,omitempty"`
}

// implements the service definition of CreateVolumeResponse
type CreateVolumeResponse struct {
	Iops            int64           `json:"Iops,omitempty"`
	LinkedVolumes   []LinkedVolumes `json:"LinkedVolumes,omitempty"`
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
	Size            int64           `json:"Size,omitempty"`
	SnapshotId      string          `json:"SnapshotId,omitempty"`
	State           string          `json:"State,omitempty"`
	SubRegionName   string          `json:"SubRegionName,omitempty"`
	Tags            []Tags          `json:"Tags,omitempty"`
	Type            string          `json:"Type,omitempty"`
	VolumeId        string          `json:"VolumeId,omitempty"`
}

// implements the service definition of CreateVpnGatewayRequest
type CreateVpnGatewayRequest struct {
	DryRun bool   `json:"DryRun,omitempty"`
	Type   string `json:"Type,omitempty"`
}

// implements the service definition of CreateVpnGatewayResponse
type CreateVpnGatewayResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
	VpnGateway      VpnGateway      `json:"VpnGateway,omitempty"`
}

// implements the service definition of DeleteFirewallRuleInboundRequest
type DeleteFirewallRuleInboundRequest struct {
	DryRun                          bool           `json:"DryRun,omitempty"`
	FirewallRulesSetId              string         `json:"FirewallRulesSetId,omitempty"`
	FromPortRange                   int64          `json:"FromPortRange,omitempty"`
	InboundRules                    []InboundRules `json:"InboundRules,omitempty"`
	IpProtocol                      string         `json:"IpProtocol,omitempty"`
	IpRange                         string         `json:"IpRange,omitempty"`
	SourceFirewallRulesSetAccountId string         `json:"SourceFirewallRulesSetAccountId,omitempty"`
	SourceFirewallRulesSetName      string         `json:"SourceFirewallRulesSetName,omitempty"`
	ToPortRange                     int64          `json:"ToPortRange,omitempty"`
}

// implements the service definition of DeleteFirewallRuleInboundResponse
type DeleteFirewallRuleInboundResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of DeleteFirewallRuleOutboundRequest
type DeleteFirewallRuleOutboundRequest struct {
	DestinationFirewallRulesSetAccountId string          `json:"DestinationFirewallRulesSetAccountId,omitempty"`
	DestinationFirewallRulesSetName      string          `json:"DestinationFirewallRulesSetName,omitempty"`
	DryRun                               bool            `json:"DryRun,omitempty"`
	FirewallRulesSetId                   string          `json:"FirewallRulesSetId,omitempty"`
	FromPortRange                        int64           `json:"FromPortRange,omitempty"`
	IpProtocol                           string          `json:"IpProtocol,omitempty"`
	IpRange                              string          `json:"IpRange,omitempty"`
	OutboundRules                        []OutboundRules `json:"OutboundRules,omitempty"`
	ToPortRange                          int64           `json:"ToPortRange,omitempty"`
}

// implements the service definition of DeleteFirewallRuleOutboundResponse
type DeleteFirewallRuleOutboundResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of DeleteFirewallRulesSetRequest
type DeleteFirewallRulesSetRequest struct {
	DryRun             bool   `json:"DryRun,omitempty"`
	FirewallRulesSetId string `json:"FirewallRulesSetId,omitempty"`
	Name               string `json:"Name,omitempty"`
}

// implements the service definition of DeleteFirewallRulesSetResponse
type DeleteFirewallRulesSetResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of DeleteKeypairRequest
type DeleteKeypairRequest struct {
	DryRun      bool   `json:"DryRun,omitempty"`
	KeypairName string `json:"KeypairName,omitempty"`
}

// implements the service definition of DeleteKeypairResponse
type DeleteKeypairResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of DeleteNatServiceRequest
type DeleteNatServiceRequest struct {
	DryRun       bool   `json:"DryRun,omitempty"`
	NatServiceId string `json:"NatServiceId,omitempty"`
}

// implements the service definition of DeleteNatServiceResponse
type DeleteNatServiceResponse struct {
	NatServiceId    string          `json:"NatServiceId,omitempty"`
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of DeleteNetInternetGatewayRequest
type DeleteNetInternetGatewayRequest struct {
	DryRun               bool   `json:"DryRun,omitempty"`
	NetInternetGatewayId string `json:"NetInternetGatewayId,omitempty"`
}

// implements the service definition of DeleteNetInternetGatewayResponse
type DeleteNetInternetGatewayResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of DeleteNetPeeringRequest
type DeleteNetPeeringRequest struct {
	DryRun       bool   `json:"DryRun,omitempty"`
	NetPeeringId string `json:"NetPeeringId,omitempty"`
}

// implements the service definition of DeleteNetPeeringResponse
type DeleteNetPeeringResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of DeleteNetRequest
type DeleteNetRequest struct {
	DryRun bool   `json:"DryRun,omitempty"`
	NetId  string `json:"NetId,omitempty"`
}

// implements the service definition of DeleteNetResponse
type DeleteNetResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of DeleteNicRequest
type DeleteNicRequest struct {
	DryRun bool   `json:"DryRun,omitempty"`
	NicId  string `json:"NicId,omitempty"`
}

// implements the service definition of DeleteNicResponse
type DeleteNicResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of DeletePublicIpRequest
type DeletePublicIpRequest struct {
	DryRun        bool   `json:"DryRun,omitempty"`
	PublicIp      string `json:"PublicIp,omitempty"`
	ReservationId string `json:"ReservationId,omitempty"`
}

// implements the service definition of DeletePublicIpResponse
type DeletePublicIpResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of DeleteRouteRequest
type DeleteRouteRequest struct {
	DestinationIpRange string `json:"DestinationIpRange,omitempty"`
	DryRun             bool   `json:"DryRun,omitempty"`
	RouteTableId       string `json:"RouteTableId,omitempty"`
}

// implements the service definition of DeleteRouteResponse
type DeleteRouteResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of DeleteRouteTableRequest
type DeleteRouteTableRequest struct {
	DryRun       bool   `json:"DryRun,omitempty"`
	RouteTableId string `json:"RouteTableId,omitempty"`
}

// implements the service definition of DeleteRouteTableResponse
type DeleteRouteTableResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of DeleteSnapshotRequest
type DeleteSnapshotRequest struct {
	DryRun     bool   `json:"DryRun,omitempty"`
	SnapshotId string `json:"SnapshotId,omitempty"`
}

// implements the service definition of DeleteSnapshotResponse
type DeleteSnapshotResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of DeleteSubnetRequest
type DeleteSubnetRequest struct {
	DryRun   bool   `json:"DryRun,omitempty"`
	SubnetId string `json:"SubnetId,omitempty"`
}

// implements the service definition of DeleteSubnetResponse
type DeleteSubnetResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of DeleteTagsRequest
type DeleteTagsRequest struct {
	DryRun      bool     `json:"DryRun,omitempty"`
	ResourceIds []string `json:"ResourceIds,omitempty"`
	Tags        []Tags   `json:"Tags,omitempty"`
}

// implements the service definition of DeleteTagsResponse
type DeleteTagsResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of DeleteVmsRequest
type DeleteVmsRequest struct {
	DryRun bool     `json:"DryRun,omitempty"`
	VmIds  []string `json:"VmIds,omitempty"`
}

// implements the service definition of DeleteVmsResponse
type DeleteVmsResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
	Vms             []DeleteVms_Vms `json:"Vms,omitempty"`
}

// implements the service definition of DeleteVms_Vms
type DeleteVms_Vms struct {
	CurrentState  string `json:"CurrentState,omitempty"`
	PreviousState string `json:"PreviousState,omitempty"`
	VmId          string `json:"VmId,omitempty"`
}

// implements the service definition of DeleteVolumeRequest
type DeleteVolumeRequest struct {
	DryRun   bool   `json:"DryRun,omitempty"`
	VolumeId string `json:"VolumeId,omitempty"`
}

// implements the service definition of DeleteVolumeResponse
type DeleteVolumeResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of DeleteVpnGatewayRequest
type DeleteVpnGatewayRequest struct {
	DryRun       bool   `json:"DryRun,omitempty"`
	VpnGatewayId string `json:"VpnGatewayId,omitempty"`
}

// implements the service definition of DeleteVpnGatewayResponse
type DeleteVpnGatewayResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of DeregisterImageRequest
type DeregisterImageRequest struct {
	DryRun  bool   `json:"DryRun,omitempty"`
	ImageId string `json:"ImageId,omitempty"`
}

// implements the service definition of DeregisterImageResponse
type DeregisterImageResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of FirewallRulesSetsMembers
type FirewallRulesSetsMembers struct {
	AccountId          string `json:"AccountId,omitempty"`
	FirewallRulesSetId string `json:"FirewallRulesSetId,omitempty"`
	Name               string `json:"Name,omitempty"`
}

// implements the service definition of Image
type Image struct {
	AccountAlias        string                          `json:"AccountAlias,omitempty"`
	AccountId           string                          `json:"AccountId,omitempty"`
	Architecture        string                          `json:"Architecture,omitempty"`
	BlockDeviceMappings []CopyImage_BlockDeviceMappings `json:"BlockDeviceMappings,omitempty"`
	CreationDate        string                          `json:"CreationDate,omitempty"`
	Description         string                          `json:"Description,omitempty"`
	ImageId             string                          `json:"ImageId,omitempty"`
	Name                string                          `json:"Name,omitempty"`
	OsuLocation         string                          `json:"OsuLocation,omitempty"`
	PermissionToLaunch  CopyImage_PermissionToLaunch    `json:"PermissionToLaunch,omitempty"`
	ProductCodes        []string                        `json:"ProductCodes,omitempty"`
	RootDeviceName      string                          `json:"RootDeviceName,omitempty"`
	RootDeviceType      string                          `json:"RootDeviceType,omitempty"`
	State               string                          `json:"State,omitempty"`
	StateComment        StateComment                    `json:"StateComment,omitempty"`
	Tags                []Tags                          `json:"Tags,omitempty"`
	Type                string                          `json:"Type,omitempty"`
}

// implements the service definition of ImageExportTask
type ImageExportTask struct {
	Comment   string    `json:"Comment,omitempty"`
	ImageId   string    `json:"ImageId,omitempty"`
	OsuExport OsuExport `json:"OsuExport,omitempty"`
	Progress  int64     `json:"Progress,omitempty"`
	State     string    `json:"State,omitempty"`
	TaskId    string    `json:"TaskId,omitempty"`
}

// implements the service definition of ImageExportTasks
type ImageExportTasks struct {
	Comment   string    `json:"Comment,omitempty"`
	ImageId   string    `json:"ImageId,omitempty"`
	OsuExport OsuExport `json:"OsuExport,omitempty"`
	Progress  int64     `json:"Progress,omitempty"`
	State     string    `json:"State,omitempty"`
	TaskId    string    `json:"TaskId,omitempty"`
}

// implements the service definition of Images
type Images struct {
	AccountAlias        string                           `json:"AccountAlias,omitempty"`
	AccountId           string                           `json:"AccountId,omitempty"`
	Architecture        string                           `json:"Architecture,omitempty"`
	BlockDeviceMappings []ReadImages_BlockDeviceMappings `json:"BlockDeviceMappings,omitempty"`
	CreationDate        string                           `json:"CreationDate,omitempty"`
	Description         string                           `json:"Description,omitempty"`
	ImageId             string                           `json:"ImageId,omitempty"`
	Name                string                           `json:"Name,omitempty"`
	OsuLocation         string                           `json:"OsuLocation,omitempty"`
	PermissionToLaunch  ReadImages_PermissionToLaunch    `json:"PermissionToLaunch,omitempty"`
	ProductCodes        []string                         `json:"ProductCodes,omitempty"`
	RootDeviceName      string                           `json:"RootDeviceName,omitempty"`
	RootDeviceType      string                           `json:"RootDeviceType,omitempty"`
	State               string                           `json:"State,omitempty"`
	StateComment        StateComment                     `json:"StateComment,omitempty"`
	Tags                []Tags                           `json:"Tags,omitempty"`
	Type                string                           `json:"Type,omitempty"`
}

// implements the service definition of ImportSnapshotRequest
type ImportSnapshotRequest struct {
	Description  string `json:"Description,omitempty"`
	DryRun       bool   `json:"DryRun,omitempty"`
	OsuLocation  string `json:"OsuLocation,omitempty"`
	SnapshotSize int64  `json:"SnapshotSize,omitempty"`
}

// implements the service definition of ImportSnapshotResponse
type ImportSnapshotResponse struct {
	Snapshot Snapshot `json:"Snapshot,omitempty"`
}

// implements the service definition of InboundRules
type InboundRules struct {
	FirewallRulesSetsMembers []FirewallRulesSetsMembers `json:"FirewallRulesSetsMembers,omitempty"`
	FromPortRange            int64                      `json:"FromPortRange,omitempty"`
	IpProtocol               string                     `json:"IpProtocol,omitempty"`
	IpRanges                 []string                   `json:"IpRanges,omitempty"`
	PrefixListIds            []string                   `json:"PrefixListIds,omitempty"`
	ToPortRange              int64                      `json:"ToPortRange,omitempty"`
}

// implements the service definition of Keypair
type Keypair struct {
	KeypairFingerprint string `json:"KeypairFingerprint,omitempty"`
	KeypairName        string `json:"KeypairName,omitempty"`
	PrivateKey         string `json:"PrivateKey,omitempty"`
}

// implements the service definition of Keypairs
type Keypairs struct {
	KeypairFingerprint string `json:"KeypairFingerprint,omitempty"`
	KeypairName        string `json:"KeypairName,omitempty"`
}

// implements the service definition of LinkNetInternetGatewayRequest
type LinkNetInternetGatewayRequest struct {
	DryRun               bool   `json:"DryRun,omitempty"`
	NetId                string `json:"NetId,omitempty"`
	NetInternetGatewayId string `json:"NetInternetGatewayId,omitempty"`
}

// implements the service definition of LinkNetInternetGatewayResponse
type LinkNetInternetGatewayResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of LinkNicRequest
type LinkNicRequest struct {
	DeviceNumber int64  `json:"DeviceNumber,omitempty"`
	DryRun       bool   `json:"DryRun,omitempty"`
	NicId        string `json:"NicId,omitempty"`
	VmId         string `json:"VmId,omitempty"`
}

// implements the service definition of LinkNicResponse
type LinkNicResponse struct {
	NicLinkId       string          `json:"NicLinkId,omitempty"`
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of LinkPublicIpRequest
type LinkPublicIpRequest struct {
	AllowRelink   bool   `json:"AllowRelink,omitempty"`
	DryRun        bool   `json:"DryRun,omitempty"`
	NicId         string `json:"NicId,omitempty"`
	PrivateIp     string `json:"PrivateIp,omitempty"`
	PublicIp      string `json:"PublicIp,omitempty"`
	ReservationId string `json:"ReservationId,omitempty"`
	VmId          string `json:"VmId,omitempty"`
}

// implements the service definition of LinkPublicIpResponse
type LinkPublicIpResponse struct {
	LinkId          string          `json:"LinkId,omitempty"`
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of LinkRouteTableRequest
type LinkRouteTableRequest struct {
	DryRun       bool   `json:"DryRun,omitempty"`
	RouteTableId string `json:"RouteTableId,omitempty"`
	SubnetId     string `json:"SubnetId,omitempty"`
}

// implements the service definition of LinkRouteTableResponse
type LinkRouteTableResponse struct {
	LinkId          string          `json:"LinkId,omitempty"`
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of LinkVolumeRequest
type LinkVolumeRequest struct {
	DeviceName string `json:"DeviceName,omitempty"`
	DryRun     bool   `json:"DryRun,omitempty"`
	VmId       string `json:"VmId,omitempty"`
	VolumeId   string `json:"VolumeId,omitempty"`
}

// implements the service definition of LinkVolumeResponse
type LinkVolumeResponse struct {
	DeleteOnVmDeletion bool            `json:"DeleteOnVmDeletion,omitempty"`
	DeviceName         string          `json:"DeviceName,omitempty"`
	ResponseContext    ResponseContext `json:"ResponseContext,omitempty"`
	State              string          `json:"State,omitempty"`
	VmId               string          `json:"VmId,omitempty"`
	VolumeId           string          `json:"VolumeId,omitempty"`
}

// implements the service definition of LinkVpnGatewayRequest
type LinkVpnGatewayRequest struct {
	DryRun       bool   `json:"DryRun,omitempty"`
	NetId        string `json:"NetId,omitempty"`
	VpnGatewayId string `json:"VpnGatewayId,omitempty"`
}

// implements the service definition of LinkVpnGatewayResponse
type LinkVpnGatewayResponse struct {
	NetToVpnGatewayLink NetToVpnGatewayLink `json:"NetToVpnGatewayLink,omitempty"`
	ResponseContext     ResponseContext     `json:"ResponseContext,omitempty"`
}

// implements the service definition of LinkedVolumes
type LinkedVolumes struct {
	DeleteOnVmDeletion bool   `json:"DeleteOnVmDeletion,omitempty"`
	DeviceName         string `json:"DeviceName,omitempty"`
	State              string `json:"State,omitempty"`
	VmId               string `json:"VmId,omitempty"`
	VolumeId           string `json:"VolumeId,omitempty"`
}

// implements the service definition of Links
type Links struct {
	Main                     bool   `json:"Main,omitempty"`
	RouteTableId             string `json:"RouteTableId,omitempty"`
	RouteTableToSubnetLinkId string `json:"RouteTableToSubnetLinkId,omitempty"`
	SubnetId                 string `json:"SubnetId,omitempty"`
}

// implements the service definition of MaintenanceEvents
type MaintenanceEvents struct {
	Description string `json:"Description,omitempty"`
	NotAfter    string `json:"NotAfter,omitempty"`
	NotBefore   string `json:"NotBefore,omitempty"`
	Reason      string `json:"Reason,omitempty"`
}

// implements the service definition of NatService
type NatService struct {
	NatServiceId string                       `json:"NatServiceId,omitempty"`
	NetId        string                       `json:"NetId,omitempty"`
	PublicIps    []CreateNatService_PublicIps `json:"PublicIps,omitempty"`
	State        string                       `json:"State,omitempty"`
	SubnetId     string                       `json:"SubnetId,omitempty"`
}

// implements the service definition of NatServices
type NatServices struct {
	NatServiceId string                      `json:"NatServiceId,omitempty"`
	NetId        string                      `json:"NetId,omitempty"`
	PublicIps    []ReadNatServices_PublicIps `json:"PublicIps,omitempty"`
	State        string                      `json:"State,omitempty"`
	SubnetId     string                      `json:"SubnetId,omitempty"`
}

// implements the service definition of Net
type Net struct {
	DhcpOptionsSetId string `json:"DhcpOptionsSetId,omitempty"`
	IpRange          string `json:"IpRange,omitempty"`
	NetId            string `json:"NetId,omitempty"`
	State            string `json:"State,omitempty"`
	Tags             []Tags `json:"Tags,omitempty"`
	Tenancy          string `json:"Tenancy,omitempty"`
}

// implements the service definition of NetInternetGateway
type NetInternetGateway struct {
	NetInternetGatewayId         string                         `json:"NetInternetGatewayId,omitempty"`
	NetToNetInternetGatewayLinks []NetToNetInternetGatewayLinks `json:"NetToNetInternetGatewayLinks,omitempty"`
	Tags                         []Tags                         `json:"Tags,omitempty"`
}

// implements the service definition of NetInternetGateways
type NetInternetGateways struct {
	NetInternetGatewayId         string                         `json:"NetInternetGatewayId,omitempty"`
	NetToNetInternetGatewayLinks []NetToNetInternetGatewayLinks `json:"NetToNetInternetGatewayLinks,omitempty"`
	Tags                         []Tags                         `json:"Tags,omitempty"`
}

// implements the service definition of NetPeering
type NetPeering struct {
	AccepterNet  AccepterNet `json:"AccepterNet,omitempty"`
	NetPeeringId string      `json:"NetPeeringId,omitempty"`
	SourceNet    SourceNet   `json:"SourceNet,omitempty"`
	State        State       `json:"State,omitempty"`
	Tags         []Tags      `json:"Tags,omitempty"`
}

// implements the service definition of NetPeerings
type NetPeerings struct {
	AccepterNet  AccepterNet `json:"AccepterNet,omitempty"`
	NetPeeringId string      `json:"NetPeeringId,omitempty"`
	SourceNet    SourceNet   `json:"SourceNet,omitempty"`
	State        State       `json:"State,omitempty"`
	Tags         []Tags      `json:"Tags,omitempty"`
}

// implements the service definition of NetToNetInternetGatewayLinks
type NetToNetInternetGatewayLinks struct {
	NetId string `json:"NetId,omitempty"`
	State string `json:"State,omitempty"`
}

// implements the service definition of NetToVpnGatewayLink
type NetToVpnGatewayLink struct {
	NetId string `json:"NetId,omitempty"`
	State string `json:"State,omitempty"`
}

// implements the service definition of NetToVpnGatewayLinks
type NetToVpnGatewayLinks struct {
	NetId string `json:"NetId,omitempty"`
	State string `json:"State,omitempty"`
}

// implements the service definition of Nets
type Nets struct {
	DhcpOptionsSetId string `json:"DhcpOptionsSetId,omitempty"`
	IpRange          string `json:"IpRange,omitempty"`
	NetId            string `json:"NetId,omitempty"`
	State            string `json:"State,omitempty"`
	Tags             []Tags `json:"Tags,omitempty"`
	Tenancy          string `json:"Tenancy,omitempty"`
}

// implements the service definition of Nic
type Nic struct {
	AccountId           string                          `json:"AccountId,omitempty"`
	Description         string                          `json:"Description,omitempty"`
	FirewallRulesSets   []CreateNic_FirewallRulesSets   `json:"FirewallRulesSets,omitempty"`
	IsSourceDestChecked bool                            `json:"IsSourceDestChecked,omitempty"`
	MacAddress          string                          `json:"MacAddress,omitempty"`
	NetId               string                          `json:"NetId,omitempty"`
	NicId               string                          `json:"NicId,omitempty"`
	NicLink             CreateNic_NicLink               `json:"NicLink,omitempty"`
	PrivateDnsName      string                          `json:"PrivateDnsName,omitempty"`
	PrivateIps          []CreateNic_Response_PrivateIps `json:"PrivateIps,omitempty"`
	PublicIpToNicLink   CreateNic_PublicIpToNicLink     `json:"PublicIpToNicLink,omitempty"`
	State               string                          `json:"State,omitempty"`
	SubnetId            string                          `json:"SubnetId,omitempty"`
	SubregionName       string                          `json:"SubregionName,omitempty"`
	Tags                []Tags                          `json:"Tags,omitempty"`
}

// implements the service definition of OsuApiKey
type OsuApiKey struct {
	ApiKeyId  string `json:"ApiKeyId,omitempty"`
	SecretKey string `json:"SecretKey,omitempty"`
}

// implements the service definition of OsuExport
type OsuExport struct {
	DiskImageFormat string    `json:"DiskImageFormat,omitempty"`
	OsuApiKey       OsuApiKey `json:"OsuApiKey,omitempty"`
	OsuBucket       string    `json:"OsuBucket,omitempty"`
	OsuManifestUrl  string    `json:"OsuManifestUrl,omitempty"`
	OsuPrefix       string    `json:"OsuPrefix,omitempty"`
}

// implements the service definition of OutboundRules
type OutboundRules struct {
	FirewallRulesSetsMembers []FirewallRulesSetsMembers `json:"FirewallRulesSetsMembers,omitempty"`
	FromPortRange            int64                      `json:"FromPortRange,omitempty"`
	IpProtocol               string                     `json:"IpProtocol,omitempty"`
	IpRanges                 []string                   `json:"IpRanges,omitempty"`
	PrefixListIds            []string                   `json:"PrefixListIds,omitempty"`
	ToPortRange              int64                      `json:"ToPortRange,omitempty"`
}

// implements the service definition of ProductCodes
type ProductCodes struct {
	ProductCode string `json:"ProductCode,omitempty"`
	ProductType string `json:"ProductType,omitempty"`
}

// implements the service definition of ReadFirewallRulesSetsRequest
type ReadFirewallRulesSetsRequest struct {
	DryRun              bool                            `json:"DryRun,omitempty"`
	Filters             []ReadFirewallRulesSets_Filters `json:"Filters,omitempty"`
	FirewallRulesSetIds []string                        `json:"FirewallRulesSetIds,omitempty"`
	Names               []string                        `json:"Names,omitempty"`
}

// implements the service definition of ReadFirewallRulesSetsResponse
type ReadFirewallRulesSetsResponse struct {
	FirewallRulesSets []ReadFirewallRulesSets_FirewallRulesSets `json:"FirewallRulesSets,omitempty"`
	ResponseContext   ResponseContext                           `json:"ResponseContext,omitempty"`
}

// implements the service definition of ReadFirewallRulesSets_Filters
type ReadFirewallRulesSets_Filters struct {
	Name   string   `json:"Name,omitempty"`
	Values []string `json:"Values,omitempty"`
}

// implements the service definition of ReadFirewallRulesSets_FirewallRulesSets
type ReadFirewallRulesSets_FirewallRulesSets struct {
	AccountId          string          `json:"AccountId,omitempty"`
	Description        string          `json:"Description,omitempty"`
	FirewallRulesSetId string          `json:"FirewallRulesSetId,omitempty"`
	InboundRules       []InboundRules  `json:"InboundRules,omitempty"`
	Name               string          `json:"Name,omitempty"`
	NetId              string          `json:"NetId,omitempty"`
	OutboundRules      []OutboundRules `json:"OutboundRules,omitempty"`
	Tags               []Tags          `json:"Tags,omitempty"`
}

// implements the service definition of ReadImageExportTasksRequest
type ReadImageExportTasksRequest struct {
	DryRun           bool                         `json:"DryRun,omitempty"`
	Filters          ReadImageExportTasks_Filters `json:"Filters,omitempty"`
	MaxResults       int64                        `json:"MaxResults,omitempty"`
	NextResultsToken string                       `json:"NextResultsToken,omitempty"`
}

// implements the service definition of ReadImageExportTasksResponse
type ReadImageExportTasksResponse struct {
	ImageExportTasks []ImageExportTasks `json:"ImageExportTasks,omitempty"`
	ResponseContext  ResponseContext    `json:"ResponseContext,omitempty"`
}

// implements the service definition of ReadImageExportTasks_Filters
type ReadImageExportTasks_Filters struct {
	TaskIds []string `json:"TaskIds,omitempty"`
}

// implements the service definition of ReadImagesRequest
type ReadImagesRequest struct {
	DryRun  bool               `json:"DryRun,omitempty"`
	Filters ReadImages_Filters `json:"Filters,omitempty"`
}

// implements the service definition of ReadImagesResponse
type ReadImagesResponse struct {
	Images          []Images        `json:"Images,omitempty"`
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of ReadImages_BlockDeviceMappings
type ReadImages_BlockDeviceMappings struct {
	Bsu               ReadImages_Bsu `json:"Bsu,omitempty"`
	DeviceName        string         `json:"DeviceName,omitempty"`
	VirtualDeviceName string         `json:"VirtualDeviceName,omitempty"`
}

// implements the service definition of ReadImages_Bsu
type ReadImages_Bsu struct {
	DeleteOnVmDeletion bool   `json:"DeleteOnVmDeletion,omitempty"`
	Iops               int64  `json:"Iops,omitempty"`
	SnapshotId         string `json:"SnapshotId,omitempty"`
	VolumeSize         int64  `json:"VolumeSize,omitempty"`
	VolumeType         string `json:"VolumeType,omitempty"`
}

// implements the service definition of ReadImages_Filters
type ReadImages_Filters struct {
	AccountAliases                           []string `json:"AccountAliases,omitempty"`
	AccountIds                               []string `json:"AccountIds,omitempty"`
	Architectures                            []string `json:"Architectures,omitempty"`
	BlockDeviceMappingDeleteOnVmTerminations []bool   `json:"BlockDeviceMappingDeleteOnVmTerminations,omitempty"`
	BlockDeviceMappingDeviceNames            []string `json:"BlockDeviceMappingDeviceNames,omitempty"`
	BlockDeviceMappingSnapshotIds            []string `json:"BlockDeviceMappingSnapshotIds,omitempty"`
	BlockDeviceMappingVolumeSize             []int64  `json:"BlockDeviceMappingVolumeSize,omitempty"`
	BlockDeviceMappingVolumeType             []string `json:"BlockDeviceMappingVolumeType,omitempty"`
	Descriptions                             []string `json:"Descriptions,omitempty"`
	Hypervisors                              []string `json:"Hypervisors,omitempty"`
	ImageIds                                 []string `json:"ImageIds,omitempty"`
	ImageNames                               []string `json:"ImageNames,omitempty"`
	ImageTypes                               []string `json:"ImageTypes,omitempty"`
	KernelIds                                []string `json:"KernelIds,omitempty"`
	ManifestLocation                         []string `json:"ManifestLocation,omitempty"`
	PermissionToLaunchAccountIds             []string `json:"PermissionToLaunchAccountIds,omitempty"`
	PermissionToLaunchGlobalPermissions      []bool   `json:"PermissionToLaunchGlobalPermissions,omitempty"`
	ProductCodes                             []string `json:"ProductCodes,omitempty"`
	RamDiskIds                               []string `json:"RamDiskIds,omitempty"`
	RootDeviceNames                          []string `json:"RootDeviceNames,omitempty"`
	RootDeviceTypes                          []string `json:"RootDeviceTypes,omitempty"`
	States                                   []string `json:"States,omitempty"`
	System                                   []string `json:"System,omitempty"`
	TagKeys                                  []string `json:"TagKeys,omitempty"`
	TagValues                                []string `json:"TagValues,omitempty"`
	Tags                                     []string `json:"Tags,omitempty"`
	VirtualizationTypes                      []string `json:"VirtualizationTypes,omitempty"`
}

// implements the service definition of ReadImages_PermissionToLaunch
type ReadImages_PermissionToLaunch struct {
	AccountIds       []string `json:"AccountIds,omitempty"`
	GlobalPermission bool     `json:"GlobalPermission,omitempty"`
}

// implements the service definition of ReadKeypairsRequest
type ReadKeypairsRequest struct {
	DryRun  bool                 `json:"DryRun,omitempty"`
	Filters ReadKeypairs_Filters `json:"Filters,omitempty"`
}

// implements the service definition of ReadKeypairsResponse
type ReadKeypairsResponse struct {
	Keypairs        []Keypairs      `json:"Keypairs,omitempty"`
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of ReadKeypairs_Filters
type ReadKeypairs_Filters struct {
	KeypairFingerprints []string `json:"KeypairFingerprints,omitempty"`
	KeypairNames        []string `json:"KeypairNames,omitempty"`
}

// implements the service definition of ReadNatServicesRequest
type ReadNatServicesRequest struct {
	DryRun           bool                    `json:"DryRun,omitempty"`
	Filters          ReadNatServices_Filters `json:"Filters,omitempty"`
	MaxResults       int64                   `json:"MaxResults,omitempty"`
	NextResultsToken string                  `json:"NextResultsToken,omitempty"`
}

// implements the service definition of ReadNatServicesResponse
type ReadNatServicesResponse struct {
	NatServices      []NatServices   `json:"NatServices,omitempty"`
	NextResultsToken string          `json:"NextResultsToken,omitempty"`
	ResponseContext  ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of ReadNatServices_Filters
type ReadNatServices_Filters struct {
	NatServiceIds []string `json:"NatServiceIds,omitempty"`
	NetIds        []string `json:"NetIds,omitempty"`
	States        []string `json:"States,omitempty"`
	SubnetIds     []string `json:"SubnetIds,omitempty"`
	TagKeys       []string `json:"TagKeys,omitempty"`
	TagValues     []string `json:"TagValues,omitempty"`
	Tags          []string `json:"Tags,omitempty"`
}

// implements the service definition of ReadNatServices_PublicIps
type ReadNatServices_PublicIps struct {
	LinkId   string `json:"LinkId,omitempty"`
	PublicIp string `json:"PublicIp,omitempty"`
}

// implements the service definition of ReadNetInternetGatewaysRequest
type ReadNetInternetGatewaysRequest struct {
	DryRun                bool                              `json:"DryRun,omitempty"`
	Filters               []ReadNetInternetGateways_Filters `json:"Filters,omitempty"`
	NetInternetGatewayIds []string                          `json:"NetInternetGatewayIds,omitempty"`
}

// implements the service definition of ReadNetInternetGatewaysResponse
type ReadNetInternetGatewaysResponse struct {
	NetInternetGateways []NetInternetGateways `json:"NetInternetGateways,omitempty"`
	ResponseContext     ResponseContext       `json:"ResponseContext,omitempty"`
}

// implements the service definition of ReadNetInternetGateways_Filters
type ReadNetInternetGateways_Filters struct {
	Name   string   `json:"Name,omitempty"`
	Values []string `json:"Values,omitempty"`
}

// implements the service definition of ReadNetPeeringsRequest
type ReadNetPeeringsRequest struct {
	DryRun  bool                    `json:"DryRun,omitempty"`
	Filters ReadNetPeerings_Filters `json:"Filters,omitempty"`
}

// implements the service definition of ReadNetPeeringsResponse
type ReadNetPeeringsResponse struct {
	NetPeerings     []NetPeerings   `json:"NetPeerings,omitempty"`
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of ReadNetPeerings_Filters
type ReadNetPeerings_Filters struct {
	AccepterNetAccountIds []string `json:"AccepterNetAccountIds,omitempty"`
	AccepterNetIpRanges   []string `json:"AccepterNetIpRanges,omitempty"`
	AccepterNetNetIds     []string `json:"AccepterNetNetIds,omitempty"`
	NetPeeringIds         []string `json:"NetPeeringIds,omitempty"`
	SourceNetAccountIds   []string `json:"SourceNetAccountIds,omitempty"`
	SourceNetIpRanges     []string `json:"SourceNetIpRanges,omitempty"`
	SourceNetNetIds       []string `json:"SourceNetNetIds,omitempty"`
	StateMessages         []string `json:"StateMessages,omitempty"`
	StateNames            []string `json:"StateNames,omitempty"`
	TagKeys               []string `json:"TagKeys,omitempty"`
	TagValues             []string `json:"TagValues,omitempty"`
	Tags                  []string `json:"Tags,omitempty"`
}

// implements the service definition of ReadNetsRequest
type ReadNetsRequest struct {
	DryRun  bool             `json:"DryRun,omitempty"`
	Filters ReadNets_Filters `json:"Filters,omitempty"`
}

// implements the service definition of ReadNetsResponse
type ReadNetsResponse struct {
	Nets            []Nets          `json:"Nets,omitempty"`
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of ReadNets_Filters
type ReadNets_Filters struct {
	DhcpOptionsSetIds []string `json:"DhcpOptionsSetIds,omitempty"`
	IpRanges          []string `json:"IpRanges,omitempty"`
	IsDefault         []string `json:"IsDefault,omitempty"`
	NetIds            []string `json:"NetIds,omitempty"`
	States            []string `json:"States,omitempty"`
	TagKeys           []string `json:"TagKeys,omitempty"`
	TagValues         []string `json:"TagValues,omitempty"`
	Tags              []string `json:"Tags,omitempty"`
}

// implements the service definition of ReadNicsRequest
type ReadNicsRequest struct {
	DryRun  bool               `json:"DryRun,omitempty"`
	Filters []ReadNics_Filters `json:"Filters,omitempty"`
	NicIds  []string           `json:"NicIds,omitempty"`
}

// implements the service definition of ReadNicsResponse
type ReadNicsResponse struct {
	Nics            []ReadNics_Nics `json:"Nics,omitempty"`
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of ReadNics_Filters
type ReadNics_Filters struct {
	Name   string   `json:"Name,omitempty"`
	Values []string `json:"Values,omitempty"`
}

// implements the service definition of ReadNics_FirewallRulesSets
type ReadNics_FirewallRulesSets struct {
	FirewallRulesSetId   string `json:"FirewallRulesSetId,omitempty"`
	FirewallRulesSetName string `json:"FirewallRulesSetName,omitempty"`
}

// implements the service definition of ReadNics_NicLink
type ReadNics_NicLink struct {
	DeleteOnVmDeletion bool   `json:"DeleteOnVmDeletion,omitempty"`
	DeviceNumber       int64  `json:"DeviceNumber,omitempty"`
	NicLinkId          string `json:"NicLinkId,omitempty"`
	State              string `json:"State,omitempty"`
	VmAccountId        string `json:"VmAccountId,omitempty"`
	VmId               string `json:"VmId,omitempty"`
}

// implements the service definition of ReadNics_Nics
type ReadNics_Nics struct {
	AccountId           string                         `json:"AccountId,omitempty"`
	Description         string                         `json:"Description,omitempty"`
	FirewallRulesSets   []ReadNics_FirewallRulesSets   `json:"FirewallRulesSets,omitempty"`
	IsSourceDestChecked bool                           `json:"IsSourceDestChecked,omitempty"`
	MacAddress          string                         `json:"MacAddress,omitempty"`
	NetId               string                         `json:"NetId,omitempty"`
	NicId               string                         `json:"NicId,omitempty"`
	NicLink             ReadNics_NicLink               `json:"NicLink,omitempty"`
	PrivateDnsName      string                         `json:"PrivateDnsName,omitempty"`
	PrivateIps          []ReadNics_Response_PrivateIps `json:"PrivateIps,omitempty"`
	PublicIpToNicLink   ReadNics_PublicIpToNicLink     `json:"PublicIpToNicLink,omitempty"`
	State               string                         `json:"State,omitempty"`
	SubnetId            string                         `json:"SubnetId,omitempty"`
	SubregionName       string                         `json:"SubregionName,omitempty"`
	Tags                []Tags                         `json:"Tags,omitempty"`
}

// implements the service definition of ReadNics_PublicIpToNicLink
type ReadNics_PublicIpToNicLink struct {
	LinkId            string `json:"LinkId,omitempty"`
	PublicDnsName     string `json:"PublicDnsName,omitempty"`
	PublicIp          string `json:"PublicIp,omitempty"`
	PublicIpAccountId string `json:"PublicIpAccountId,omitempty"`
	ReservationId     string `json:"ReservationId,omitempty"`
}

// implements the service definition of ReadNics_Response_PrivateIps
type ReadNics_Response_PrivateIps struct {
	IsPrimary         bool                       `json:"IsPrimary,omitempty"`
	PrivateDnsName    string                     `json:"PrivateDnsName,omitempty"`
	PrivateIp         string                     `json:"PrivateIp,omitempty"`
	PublicIpToNicLink ReadNics_PublicIpToNicLink `json:"PublicIpToNicLink,omitempty"`
}

// implements the service definition of ReadPublicIpsRequest
type ReadPublicIpsRequest struct {
	DryRun  bool                  `json:"DryRun,omitempty"`
	Filters ReadPublicIps_Filters `json:"Filters,omitempty"`
}

// implements the service definition of ReadPublicIpsResponse
type ReadPublicIpsResponse struct {
	PublicIps       []ReadPublicIps_PublicIps `json:"PublicIps,omitempty"`
	ResponseContext ResponseContext           `json:"ResponseContext,omitempty"`
}

// implements the service definition of ReadPublicIps_Filters
type ReadPublicIps_Filters struct {
	LinkIds        []string `json:"LinkIds,omitempty"`
	NicAccountIds  []string `json:"NicAccountIds,omitempty"`
	NicIds         []string `json:"NicIds,omitempty"`
	Placements     []string `json:"Placements,omitempty"`
	PrivateIps     []string `json:"PrivateIps,omitempty"`
	PublicIps      []string `json:"PublicIps,omitempty"`
	ReservationIds []string `json:"ReservationIds,omitempty"`
	VmIds          []string `json:"VmIds,omitempty"`
}

// implements the service definition of ReadPublicIps_PublicIps
type ReadPublicIps_PublicIps struct {
	LinkId        string `json:"LinkId,omitempty"`
	NicAccountId  string `json:"NicAccountId,omitempty"`
	NicId         string `json:"NicId,omitempty"`
	Placement     string `json:"Placement,omitempty"`
	PrivateIp     string `json:"PrivateIp,omitempty"`
	PublicIp      string `json:"PublicIp,omitempty"`
	ReservationId string `json:"ReservationId,omitempty"`
	VmId          string `json:"VmId,omitempty"`
}

// implements the service definition of ReadRouteTablesRequest
type ReadRouteTablesRequest struct {
	DryRun  bool                    `json:"DryRun,omitempty"`
	Filters ReadRouteTables_Filters `json:"Filters,omitempty"`
}

// implements the service definition of ReadRouteTablesResponse
type ReadRouteTablesResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
	RouteTables     []RouteTables   `json:"RouteTables,omitempty"`
}

// implements the service definition of ReadRouteTables_Filters
type ReadRouteTables_Filters struct {
	LinkMain                      bool     `json:"LinkMain,omitempty"`
	LinkRouteTableLinkIds         []string `json:"LinkRouteTableLinkIds,omitempty"`
	LinkSubnetIds                 []string `json:"LinkSubnetIds,omitempty"`
	NetIds                        []string `json:"NetIds,omitempty"`
	RouteCreationMethods          []string `json:"RouteCreationMethods,omitempty"`
	RouteDestinationIpRanges      []string `json:"RouteDestinationIpRanges,omitempty"`
	RouteDestinationPrefixListIds []string `json:"RouteDestinationPrefixListIds,omitempty"`
	RouteGatewayIds               []string `json:"RouteGatewayIds,omitempty"`
	RouteNatServiceIds            []string `json:"RouteNatServiceIds,omitempty"`
	RouteNetPeeringIds            []string `json:"RouteNetPeeringIds,omitempty"`
	RouteStates                   []string `json:"RouteStates,omitempty"`
	RouteTableIds                 []string `json:"RouteTableIds,omitempty"`
	RouteVmIds                    []string `json:"RouteVmIds,omitempty"`
	TagKeys                       []string `json:"TagKeys,omitempty"`
	TagValues                     []string `json:"TagValues,omitempty"`
	Tags                          []Tags   `json:"Tags,omitempty"`
}

// implements the service definition of ReadSnapshotExportTasksRequest
type ReadSnapshotExportTasksRequest struct {
	DryRun  bool     `json:"DryRun,omitempty"`
	TaskIds []string `json:"TaskIds,omitempty"`
}

// implements the service definition of ReadSnapshotExportTasksResponse
type ReadSnapshotExportTasksResponse struct {
	ResponseContext     ResponseContext       `json:"ResponseContext,omitempty"`
	SnapshotExportTasks []SnapshotExportTasks `json:"SnapshotExportTasks,omitempty"`
}

// implements the service definition of ReadSnapshotsRequest
type ReadSnapshotsRequest struct {
	DryRun           bool                  `json:"DryRun,omitempty"`
	Filters          ReadSnapshots_Filters `json:"Filters,omitempty"`
	MaxResults       int64                 `json:"MaxResults,omitempty"`
	NextResultsToken string                `json:"NextResultsToken,omitempty"`
}

// implements the service definition of ReadSnapshotsResponse
type ReadSnapshotsResponse struct {
	NextResultsToken string          `json:"NextResultsToken,omitempty"`
	ResponseContext  ResponseContext `json:"ResponseContext,omitempty"`
	Snapshots        []Snapshots     `json:"Snapshots,omitempty"`
}

// implements the service definition of ReadSnapshots_Filters
type ReadSnapshots_Filters struct {
	AccountAliases                            []string `json:"AccountAliases,omitempty"`
	AccountIds                                []string `json:"AccountIds,omitempty"`
	Descriptions                              []string `json:"Descriptions,omitempty"`
	PermissionToCreateVolumeAccountIds        []string `json:"PermissionToCreateVolumeAccountIds,omitempty"`
	PermissionToCreateVolumeGlobalPermissions []bool   `json:"PermissionToCreateVolumeGlobalPermissions,omitempty"`
	Progresses                                []string `json:"Progresses,omitempty"`
	SnapshotIds                               []string `json:"SnapshotIds,omitempty"`
	States                                    []string `json:"States,omitempty"`
	TagKeys                                   []string `json:"TagKeys,omitempty"`
	TagValues                                 []string `json:"TagValues,omitempty"`
	Tags                                      []string `json:"Tags,omitempty"`
	VolumeIds                                 []string `json:"VolumeIds,omitempty"`
	VolumeSizes                               []string `json:"VolumeSizes,omitempty"`
}

// implements the service definition of ReadSnapshots_PermissionToCreateVolume
type ReadSnapshots_PermissionToCreateVolume struct {
	AccountIds       []string `json:"AccountIds,omitempty"`
	GlobalPermission bool     `json:"GlobalPermission,omitempty"`
}

// implements the service definition of ReadSubnetsRequest
type ReadSubnetsRequest struct {
	DryRun    bool                  `json:"DryRun,omitempty"`
	Filters   []ReadSubnets_Filters `json:"Filters,omitempty"`
	SubnetIds []string              `json:"SubnetIds,omitempty"`
}

// implements the service definition of ReadSubnetsResponse
type ReadSubnetsResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
	Subnets         []Subnets       `json:"Subnets,omitempty"`
}

// implements the service definition of ReadSubnets_Filters
type ReadSubnets_Filters struct {
	Name   string   `json:"Name,omitempty"`
	Values []string `json:"Values,omitempty"`
}

// implements the service definition of ReadTagsRequest
type ReadTagsRequest struct {
	DryRun           bool             `json:"DryRun,omitempty"`
	Filters          ReadTags_Filters `json:"Filters,omitempty"`
	MaxResults       int64            `json:"MaxResults,omitempty"`
	NextResultsToken string           `json:"NextResultsToken,omitempty"`
}

// implements the service definition of ReadTagsResponse
type ReadTagsResponse struct {
	NextResultsToken string          `json:"NextResultsToken,omitempty"`
	ResponseContext  ResponseContext `json:"ResponseContext,omitempty"`
	Tags             []Tags          `json:"Tags,omitempty"`
}

// implements the service definition of ReadTags_Filters
type ReadTags_Filters struct {
	Keys          []string `json:"Keys,omitempty"`
	ResourceIds   []string `json:"ResourceIds,omitempty"`
	ResourceTypes []string `json:"ResourceTypes,omitempty"`
	Values        []string `json:"Values,omitempty"`
}

// implements the service definition of ReadTags_Tags
type Tags struct {
	Key          string `json:"Key,omitempty"`
	ResourceId   string `json:"ResourceId,omitempty"`
	ResourceType string `json:"ResourceType,omitempty"`
	Value        string `json:"Value,omitempty"`
}

// implements the service definition of ReadVmAttributeRequest
type ReadVmAttributeRequest struct {
	Attribute string `json:"Attribute,omitempty"`
	DryRun    bool   `json:"DryRun,omitempty"`
	VmId      string `json:"VmId,omitempty"`
}

// implements the service definition of ReadVmAttributeResponse
type ReadVmAttributeResponse struct {
	BlockDeviceMappings         []ReadVmAttribute_BlockDeviceMappings `json:"BlockDeviceMappings,omitempty"`
	BsuOptimized                bool                                  `json:"BsuOptimized,omitempty"`
	DeletionProtection          bool                                  `json:"DeletionProtection,omitempty"`
	FirewallRulesSets           []ReadVmAttribute_FirewallRulesSets   `json:"FirewallRulesSets,omitempty"`
	IsSourceDestChecked         bool                                  `json:"IsSourceDestChecked,omitempty"`
	KeypairName                 string                                `json:"KeypairName,omitempty"`
	ProductCodes                []ProductCodes                        `json:"ProductCodes,omitempty"`
	ResponseContext             ResponseContext                       `json:"ResponseContext,omitempty"`
	RootDeviceName              string                                `json:"RootDeviceName,omitempty"`
	Type                        string                                `json:"Type,omitempty"`
	UserData                    string                                `json:"UserData,omitempty"`
	VmId                        string                                `json:"VmId,omitempty"`
	VmInitiatedShutdownBehavior string                                `json:"VmInitiatedShutdownBehavior,omitempty"`
}

// implements the service definition of ReadVmAttribute_BlockDeviceMappings
type ReadVmAttribute_BlockDeviceMappings struct {
	Bsu        ReadVmAttribute_Bsu `json:"Bsu,omitempty"`
	DeviceName string              `json:"DeviceName,omitempty"`
}

// implements the service definition of ReadVmAttribute_Bsu
type ReadVmAttribute_Bsu struct {
	DeleteOnVmDeletion bool   `json:"DeleteOnVmDeletion,omitempty"`
	LinkDate           string `json:"LinkDate,omitempty"`
	State              string `json:"State,omitempty"`
	VolumeId           string `json:"VolumeId,omitempty"`
}

// implements the service definition of ReadVmAttribute_FirewallRulesSets
type ReadVmAttribute_FirewallRulesSets struct {
	FirewallRulesSetId   string `json:"FirewallRulesSetId,omitempty"`
	FirewallRulesSetName string `json:"FirewallRulesSetName,omitempty"`
}

// implements the service definition of ReadVmsRequest
type ReadVmsRequest struct {
	DryRun           bool              `json:"DryRun,omitempty"`
	Filters          []ReadVms_Filters `json:"Filters,omitempty"`
	MaxResults       int64             `json:"MaxResults,omitempty"`
	NextResultsToken string            `json:"NextResultsToken,omitempty"`
}

// implements the service definition of ReadVmsResponse
type ReadVmsResponse struct {
	NextResultsToken string          `json:"NextResultsToken,omitempty"`
	ResponseContext  ResponseContext `json:"ResponseContext,omitempty"`
	Vms              []ReadVms_Vms   `json:"Vms,omitempty"`
}

// implements the service definition of ReadVmsStateRequest
type ReadVmsStateRequest struct {
	AllVms           bool                   `json:"AllVms,omitempty"`
	DryRun           bool                   `json:"DryRun,omitempty"`
	Filters          []ReadVmsState_Filters `json:"Filters,omitempty"`
	MaxResults       int64                  `json:"MaxResults,omitempty"`
	NextResultsToken string                 `json:"NextResultsToken,omitempty"`
}

// implements the service definition of ReadVmsStateResponse
type ReadVmsStateResponse struct {
	NextResultsToken string          `json:"NextResultsToken,omitempty"`
	ResponseContext  ResponseContext `json:"ResponseContext,omitempty"`
	VmStates         []VmStates      `json:"VmStates,omitempty"`
}

// implements the service definition of ReadVmsState_Filters
type ReadVmsState_Filters struct {
	MaintenanceEventDescriptions []string `json:"MaintenanceEventDescriptions,omitempty"`
	MaintenanceEventReasons      []string `json:"MaintenanceEventReasons,omitempty"`
	MaintenanceEventsNotAfter    []string `json:"MaintenanceEventsNotAfter,omitempty"`
	MaintenanceEventsNotBefore   []string `json:"MaintenanceEventsNotBefore,omitempty"`
	SubRegionNames               []string `json:"SubRegionNames,omitempty"`
	VmIds                        []string `json:"VmIds,omitempty"`
	VmStates                     []string `json:"VmStates,omitempty"`
}

// implements the service definition of ReadVms_Bsu
type ReadVms_Bsu struct {
	DeleteOnVmDeletion bool   `json:"DeleteOnVmDeletion,omitempty"`
	State              string `json:"State,omitempty"`
	VolumeId           string `json:"VolumeId,omitempty"`
}

// implements the service definition of ReadVms_Filters
type ReadVms_Filters struct {
	AccountIds                            []string `json:"AccountIds,omitempty"`
	ActivatedChecks                       []bool   `json:"ActivatedChecks,omitempty"`
	Architectures                         []string `json:"Architectures,omitempty"`
	BlockDeviceMappingDeleteOnVmDeletions []bool   `json:"BlockDeviceMappingDeleteOnVmDeletions,omitempty"`
	BlockDeviceMappingDeviceNames         []string `json:"BlockDeviceMappingDeviceNames,omitempty"`
	BlockDeviceMappingLinkDates           []string `json:"BlockDeviceMappingLinkDates,omitempty"`
	BlockDeviceMappingStates              []string `json:"BlockDeviceMappingStates,omitempty"`
	BlockDeviceMappingVolumeIds           []string `json:"BlockDeviceMappingVolumeIds,omitempty"`
	Comments                              []string `json:"Comments,omitempty"`
	CreationDates                         []string `json:"CreationDates,omitempty"`
	DnsNames                              []string `json:"DnsNames,omitempty"`
	FirewallRulesSetIds                   []string `json:"FirewallRulesSetIds,omitempty"`
	FirewallRulesSetNames                 []string `json:"FirewallRulesSetNames,omitempty"`
	Hypervisors                           []string `json:"Hypervisors,omitempty"`
	ImageIds                              []string `json:"ImageIds,omitempty"`
	KernelIds                             []string `json:"KernelIds,omitempty"`
	KeypairNames                          []string `json:"KeypairNames,omitempty"`
	LaunchSortNumbers                     []int64  `json:"LaunchSortNumbers,omitempty"`
	MonitoringStates                      []string `json:"MonitoringStates,omitempty"`
	NetIds                                []string `json:"NetIds,omitempty"`
	NicAccountIds                         []string `json:"NicAccountIds,omitempty"`
	NicActivatedChecks                    []bool   `json:"NicActivatedChecks,omitempty"`
	NicDescriptions                       []string `json:"NicDescriptions,omitempty"`
	NicFirewallRulesSetIds                []string `json:"NicFirewallRulesSetIds,omitempty"`
	NicFirewallRulesSetNames              []string `json:"NicFirewallRulesSetNames,omitempty"`
	NicIpLinkPrivateIpAccountIds          []string `json:"NicIpLinkPrivateIpAccountIds,omitempty"`
	NicIpLinkPublicIps                    []string `json:"NicIpLinkPublicIps,omitempty"`
	NicIpPrimaryIps                       []string `json:"NicIpPrimaryIps,omitempty"`
	NicIpPrivateIps                       []string `json:"NicIpPrivateIps,omitempty"`
	NicLinkDeleteOnVmDeletions            []bool   `json:"NicLinkDeleteOnVmDeletions,omitempty"`
	NicLinkLinkDates                      []string `json:"NicLinkLinkDates,omitempty"`
	NicLinkLinkIds                        []string `json:"NicLinkLinkIds,omitempty"`
	NicLinkNicLinkIds                     []string `json:"NicLinkNicLinkIds,omitempty"`
	NicLinkNicSortNumbers                 []int64  `json:"NicLinkNicSortNumbers,omitempty"`
	NicLinkPublicIpAccountIds             []string `json:"NicLinkPublicIpAccountIds,omitempty"`
	NicLinkPublicIps                      []string `json:"NicLinkPublicIps,omitempty"`
	NicLinkReservationIds                 []string `json:"NicLinkReservationIds,omitempty"`
	NicLinkStates                         []string `json:"NicLinkStates,omitempty"`
	NicLinkVmAccountIds                   []string `json:"NicLinkVmAccountIds,omitempty"`
	NicLinkVmIds                          []string `json:"NicLinkVmIds,omitempty"`
	NicMacAddresses                       []string `json:"NicMacAddresses,omitempty"`
	NicNetIds                             []string `json:"NicNetIds,omitempty"`
	NicNicIds                             []string `json:"NicNicIds,omitempty"`
	NicPrivateDnsNames                    []string `json:"NicPrivateDnsNames,omitempty"`
	NicRequesterIds                       []string `json:"NicRequesterIds,omitempty"`
	NicRequesterManaged                   []string `json:"NicRequesterManaged,omitempty"`
	NicStates                             []string `json:"NicStates,omitempty"`
	NicSubRegionNames                     []string `json:"NicSubRegionNames,omitempty"`
	NicSubnetIds                          []string `json:"NicSubnetIds,omitempty"`
	PlacementGroups                       []string `json:"PlacementGroups,omitempty"`
	PrivateDnsNames                       []string `json:"PrivateDnsNames,omitempty"`
	PrivateIps                            []string `json:"PrivateIps,omitempty"`
	ProductCodes                          []string `json:"ProductCodes,omitempty"`
	PublicIps                             []string `json:"PublicIps,omitempty"`
	RamDiskIds                            []string `json:"RamDiskIds,omitempty"`
	RequesterIds                          []string `json:"RequesterIds,omitempty"`
	ReservationIds                        []string `json:"ReservationIds,omitempty"`
	RootDeviceNames                       []string `json:"RootDeviceNames,omitempty"`
	RootDeviceTypes                       []string `json:"RootDeviceTypes,omitempty"`
	SpotVmRequestIds                      []string `json:"SpotVmRequestIds,omitempty"`
	SpotVms                               []string `json:"SpotVms,omitempty"`
	StateComments                         []string `json:"StateComments,omitempty"`
	SubRegionNames                        []string `json:"SubRegionNames,omitempty"`
	SubnetIds                             []string `json:"SubnetIds,omitempty"`
	Systems                               []string `json:"Systems,omitempty"`
	TagKeys                               []string `json:"TagKeys,omitempty"`
	TagValues                             []string `json:"TagValues,omitempty"`
	Tags                                  []string `json:"Tags,omitempty"`
	Tenancies                             []string `json:"Tenancies,omitempty"`
	Tokens                                []string `json:"Tokens,omitempty"`
	VirtualizationTypes                   []string `json:"VirtualizationTypes,omitempty"`
	VmIds                                 []string `json:"VmIds,omitempty"`
	VmStates                              []string `json:"VmStates,omitempty"`
	VmTypes                               []string `json:"VmTypes,omitempty"`
	VmsFirewallRulesSetIds                []string `json:"VmsFirewallRulesSetIds,omitempty"`
	VmsFirewallRulesSetNames              []string `json:"VmsFirewallRulesSetNames,omitempty"`
}

// implements the service definition of ReadVms_FirewallRulesSets
type ReadVms_FirewallRulesSets struct {
	FirewallRulesSetId   string `json:"FirewallRulesSetId,omitempty"`
	FirewallRulesSetName string `json:"FirewallRulesSetName,omitempty"`
}

// implements the service definition of ReadVms_NicLink
type ReadVms_NicLink struct {
	DeleteOnVmDeletion bool   `json:"DeleteOnVmDeletion,omitempty"`
	DeviceNumber       int64  `json:"DeviceNumber,omitempty"`
	NicLinkId          string `json:"NicLinkId,omitempty"`
	State              string `json:"State,omitempty"`
}

// implements the service definition of ReadVms_Placement
type ReadVms_Placement struct {
	Affinity        string `json:"Affinity,omitempty"`
	DedicatedHostId string `json:"DedicatedHostId,omitempty"`
	PlacementName   string `json:"PlacementName,omitempty"`
	SubRegionName   string `json:"SubRegionName,omitempty"`
	Tenancy         string `json:"Tenancy,omitempty"`
}

// implements the service definition of ReadVms_PublicIpToNicLink
type ReadVms_PublicIpToNicLink struct {
	PublicDnsName     string `json:"PublicDnsName,omitempty"`
	PublicIp          string `json:"PublicIp,omitempty"`
	PublicIpAccountId string `json:"PublicIpAccountId,omitempty"`
}

// implements the service definition of ReadVms_Response_BlockDeviceMappings
type ReadVms_Response_BlockDeviceMappings struct {
	Bsu        ReadVms_Bsu `json:"Bsu,omitempty"`
	DeviceName string      `json:"DeviceName,omitempty"`
}

// implements the service definition of ReadVms_Response_Nics
type ReadVms_Response_Nics struct {
	AccountId           string                        `json:"AccountId,omitempty"`
	Description         string                        `json:"Description,omitempty"`
	FirewallRulesSets   []ReadVms_FirewallRulesSets   `json:"FirewallRulesSets,omitempty"`
	IsSourceDestChecked bool                          `json:"IsSourceDestChecked,omitempty"`
	MacAddress          string                        `json:"MacAddress,omitempty"`
	NetId               string                        `json:"NetId,omitempty"`
	NicId               string                        `json:"NicId,omitempty"`
	NicLink             ReadVms_NicLink               `json:"NicLink,omitempty"`
	PrivateDnsName      string                        `json:"PrivateDnsName,omitempty"`
	PrivateIps          []ReadVms_Response_PrivateIps `json:"PrivateIps,omitempty"`
	PublicIpToNicLink   ReadVms_PublicIpToNicLink     `json:"PublicIpToNicLink,omitempty"`
	State               string                        `json:"State,omitempty"`
	SubnetId            string                        `json:"SubnetId,omitempty"`
}

// implements the service definition of ReadVms_Response_PrivateIps
type ReadVms_Response_PrivateIps struct {
	IsPrimary         bool                      `json:"IsPrimary,omitempty"`
	PrivateDnsName    string                    `json:"PrivateDnsName,omitempty"`
	PrivateIp         string                    `json:"PrivateIp,omitempty"`
	PublicIpToNicLink ReadVms_PublicIpToNicLink `json:"PublicIpToNicLink,omitempty"`
}

// implements the service definition of ReadVms_Vms
type ReadVms_Vms struct {
	Architecture        string                                 `json:"Architecture,omitempty"`
	BlockDeviceMappings []ReadVms_Response_BlockDeviceMappings `json:"BlockDeviceMappings,omitempty"`
	BsuOptimized        bool                                   `json:"BsuOptimized,omitempty"`
	ClientToken         string                                 `json:"ClientToken,omitempty"`
	Comment             string                                 `json:"Comment,omitempty"`
	FirewallRulesSets   []ReadVms_FirewallRulesSets            `json:"FirewallRulesSets,omitempty"`
	ImageId             string                                 `json:"ImageId,omitempty"`
	IsSourceDestChecked bool                                   `json:"IsSourceDestChecked,omitempty"`
	KeypairName         string                                 `json:"KeypairName,omitempty"`
	LaunchNumber        int64                                  `json:"LaunchNumber,omitempty"`
	NetId               string                                 `json:"NetId,omitempty"`
	Nics                []ReadVms_Response_Nics                `json:"Nics,omitempty"`
	OsFamily            string                                 `json:"OsFamily,omitempty"`
	Placement           ReadVms_Placement                      `json:"Placement,omitempty"`
	PrivateDnsName      string                                 `json:"PrivateDnsName,omitempty"`
	PrivateIp           string                                 `json:"PrivateIp,omitempty"`
	ProductCodes        []ProductCodes                         `json:"ProductCodes,omitempty"`
	PublicDnsName       string                                 `json:"PublicDnsName,omitempty"`
	PublicIp            string                                 `json:"PublicIp,omitempty"`
	ReservationId       string                                 `json:"ReservationId,omitempty"`
	RootDeviceName      string                                 `json:"RootDeviceName,omitempty"`
	RootDeviceType      string                                 `json:"RootDeviceType,omitempty"`
	State               string                                 `json:"State,omitempty"`
	SubnetId            string                                 `json:"SubnetId,omitempty"`
	Tags                []Tags                                 `json:"Tags,omitempty"`
	Transition          Transition                             `json:"Transition,omitempty"`
	Type                string                                 `json:"Type,omitempty"`
	VmId                string                                 `json:"VmId,omitempty"`
}

// implements the service definition of ReadVolumesRequest
type ReadVolumesRequest struct {
	DryRun           bool                `json:"DryRun,omitempty"`
	Filters          ReadVolumes_Filters `json:"Filters,omitempty"`
	MaxResults       int64               `json:"MaxResults,omitempty"`
	NextResultsToken string              `json:"NextResultsToken,omitempty"`
}

// implements the service definition of ReadVolumesResponse
type ReadVolumesResponse struct {
	NextResultsToken string          `json:"NextResultsToken,omitempty"`
	ResponseContext  ResponseContext `json:"ResponseContext,omitempty"`
	Volumes          []Volumes       `json:"Volumes,omitempty"`
}

// implements the service definition of ReadVolumes_Filters
type ReadVolumes_Filters struct {
	CreationDates  []string `json:"CreationDates,omitempty"`
	SnapshotIds    []string `json:"SnapshotIds,omitempty"`
	SubRegionNames []string `json:"SubRegionNames,omitempty"`
	TagKeys        []string `json:"TagKeys,omitempty"`
	TagValues      []string `json:"TagValues,omitempty"`
	Tags           []Tags   `json:"Tags,omitempty"`
	VolumeIds      []string `json:"VolumeIds,omitempty"`
	VolumeSizes    []int64  `json:"VolumeSizes,omitempty"`
	VolumeTypes    []string `json:"VolumeTypes,omitempty"`
}

// implements the service definition of ReadVpnGatewaysRequest
type ReadVpnGatewaysRequest struct {
	DryRun  bool                      `json:"DryRun,omitempty"`
	Filters []ReadVpnGateways_Filters `json:"Filters,omitempty"`
}

// implements the service definition of ReadVpnGatewaysResponse
type ReadVpnGatewaysResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
	VpnGateways     []VpnGateways   `json:"VpnGateways,omitempty"`
}

// implements the service definition of ReadVpnGateways_Filters
type ReadVpnGateways_Filters struct {
	Name   string   `json:"Name,omitempty"`
	Values []string `json:"Values,omitempty"`
}

// implements the service definition of RebootVmsRequest
type RebootVmsRequest struct {
	DryRun bool     `json:"DryRun,omitempty"`
	VmIds  []string `json:"VmIds,omitempty"`
}

// implements the service definition of RebootVmsResponse
type RebootVmsResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of RegisterImageRequest
type RegisterImageRequest struct {
	Architecture        string                              `json:"Architecture,omitempty"`
	BlockDeviceMappings []RegisterImage_BlockDeviceMappings `json:"BlockDeviceMappings,omitempty"`
	Description         string                              `json:"Description,omitempty"`
	DryRun              bool                                `json:"DryRun,omitempty"`
	Name                string                              `json:"Name,omitempty"`
	OsuLocation         string                              `json:"OsuLocation,omitempty"`
	RootDeviceName      string                              `json:"RootDeviceName,omitempty"`
}

// implements the service definition of RegisterImageResponse
type RegisterImageResponse struct {
	Image           Image           `json:"Image,omitempty"`
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of RegisterImage_BlockDeviceMappings
type RegisterImage_BlockDeviceMappings struct {
	Bsu               RegisterImage_Bsu `json:"Bsu,omitempty"`
	DeviceName        string            `json:"DeviceName,omitempty"`
	VirtualDeviceName string            `json:"VirtualDeviceName,omitempty"`
}

// implements the service definition of RegisterImage_Bsu
type RegisterImage_Bsu struct {
	DeleteOnVmDeletion bool   `json:"DeleteOnVmDeletion,omitempty"`
	Iops               int64  `json:"Iops,omitempty"`
	SnapshotId         string `json:"SnapshotId,omitempty"`
	VolumeSize         int64  `json:"VolumeSize,omitempty"`
	VolumeType         string `json:"VolumeType,omitempty"`
}

// implements the service definition of RejectNetPeeringRequest
type RejectNetPeeringRequest struct {
	DryRun       bool   `json:"DryRun,omitempty"`
	NetPeeringId string `json:"NetPeeringId,omitempty"`
}

// implements the service definition of RejectNetPeeringResponse
type RejectNetPeeringResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of Removal
type Removal struct {
	AccountIds       []string `json:"AccountIds,omitempty"`
	GlobalPermission bool     `json:"GlobalPermission,omitempty"`
}

// implements the service definition of RoutePropagatingVpnGateways
type RoutePropagatingVpnGateways struct {
	VpnGatewayId string `json:"VpnGatewayId,omitempty"`
}

// implements the service definition of RouteTable
type RouteTable struct {
	Links                       []Links                       `json:"Links,omitempty"`
	NetId                       string                        `json:"NetId,omitempty"`
	RoutePropagatingVpnGateways []RoutePropagatingVpnGateways `json:"RoutePropagatingVpnGateways,omitempty"`
	RouteTableId                string                        `json:"RouteTableId,omitempty"`
	Routes                      []Routes                      `json:"Routes,omitempty"`
	Tags                        []Tags                        `json:"Tags,omitempty"`
}

// implements the service definition of RouteTables
type RouteTables struct {
	Links                       []Links                       `json:"Links,omitempty"`
	NetId                       string                        `json:"NetId,omitempty"`
	RoutePropagatingVpnGateways []RoutePropagatingVpnGateways `json:"RoutePropagatingVpnGateways,omitempty"`
	RouteTableId                string                        `json:"RouteTableId,omitempty"`
	Routes                      []Routes                      `json:"Routes,omitempty"`
	Tags                        []Tags                        `json:"Tags,omitempty"`
}

// implements the service definition of Routes
type Routes struct {
	CreationMethod          string `json:"CreationMethod,omitempty"`
	DestinationIpRange      string `json:"DestinationIpRange,omitempty"`
	DestinationPrefixListId string `json:"DestinationPrefixListId,omitempty"`
	GatewayId               string `json:"GatewayId,omitempty"`
	NatServiceId            string `json:"NatServiceId,omitempty"`
	NetPeeringId            string `json:"NetPeeringId,omitempty"`
	NicId                   string `json:"NicId,omitempty"`
	State                   string `json:"State,omitempty"`
	VmAccountId             string `json:"VmAccountId,omitempty"`
	VmId                    string `json:"VmId,omitempty"`
}

// implements the service definition of Snapshot
type Snapshot struct {
	AccountAlias             string                                `json:"AccountAlias,omitempty"`
	AccountId                string                                `json:"AccountId,omitempty"`
	Description              string                                `json:"Description,omitempty"`
	PermissionToCreateVolume CopySnapshot_PermissionToCreateVolume `json:"PermissionToCreateVolume,omitempty"`
	Progress                 int64                                 `json:"Progress,omitempty"`
	SnapshotId               string                                `json:"SnapshotId,omitempty"`
	State                    string                                `json:"State,omitempty"`
	Tags                     []Tags                                `json:"Tags,omitempty"`
	VolumeId                 string                                `json:"VolumeId,omitempty"`
	VolumeSize               int64                                 `json:"VolumeSize,omitempty"`
}

// implements the service definition of SnapshotExportTask
type SnapshotExportTask struct {
	Comment    string    `json:"Comment,omitempty"`
	OsuExport  OsuExport `json:"OsuExport,omitempty"`
	Progress   int64     `json:"Progress,omitempty"`
	SnapshotId string    `json:"SnapshotId,omitempty"`
	State      string    `json:"State,omitempty"`
	TaskId     string    `json:"TaskId,omitempty"`
}

// implements the service definition of SnapshotExportTasks
type SnapshotExportTasks struct {
	Comment    string    `json:"Comment,omitempty"`
	OsuExport  OsuExport `json:"OsuExport,omitempty"`
	Progress   int64     `json:"Progress,omitempty"`
	SnapshotId string    `json:"SnapshotId,omitempty"`
	State      string    `json:"State,omitempty"`
	TaskId     string    `json:"TaskId,omitempty"`
}

// implements the service definition of Snapshots
type Snapshots struct {
	AccountAlias             string                                 `json:"AccountAlias,omitempty"`
	AccountId                string                                 `json:"AccountId,omitempty"`
	Description              string                                 `json:"Description,omitempty"`
	PermissionToCreateVolume ReadSnapshots_PermissionToCreateVolume `json:"PermissionToCreateVolume,omitempty"`
	Progress                 int64                                  `json:"Progress,omitempty"`
	SnapshotId               string                                 `json:"SnapshotId,omitempty"`
	State                    string                                 `json:"State,omitempty"`
	Tags                     []Tags                                 `json:"Tags,omitempty"`
	VolumeId                 string                                 `json:"VolumeId,omitempty"`
	VolumeSize               int64                                  `json:"VolumeSize,omitempty"`
}

// implements the service definition of SourceNet
type SourceNet struct {
	AccountId string   `json:"AccountId,omitempty"`
	IpRanges  []string `json:"IpRanges,omitempty"`
	NetId     string   `json:"NetId,omitempty"`
}

// implements the service definition of StartVmsRequest
type StartVmsRequest struct {
	DryRun bool     `json:"DryRun,omitempty"`
	VmIds  []string `json:"VmIds,omitempty"`
}

// implements the service definition of StartVmsResponse
type StartVmsResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
	Vms             []StartVms_Vms  `json:"Vms,omitempty"`
}

// implements the service definition of StartVms_Vms
type StartVms_Vms struct {
	CurrentState  string `json:"CurrentState,omitempty"`
	PreviousState string `json:"PreviousState,omitempty"`
	VmId          string `json:"VmId,omitempty"`
}

// implements the service definition of State
type State struct {
	Message string `json:"Message,omitempty"`
	Name    string `json:"Name,omitempty"`
}

// implements the service definition of StateComment
type StateComment struct {
	StateCode    string `json:"StateCode,omitempty"`
	StateMessage string `json:"StateMessage,omitempty"`
}

// implements the service definition of StopVmsRequest
type StopVmsRequest struct {
	DryRun    bool     `json:"DryRun,omitempty"`
	ForceStop bool     `json:"ForceStop,omitempty"`
	VmIds     []string `json:"VmIds,omitempty"`
}

// implements the service definition of StopVmsResponse
type StopVmsResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
	Vms             []StopVms_Vms   `json:"Vms,omitempty"`
}

// implements the service definition of StopVms_Vms
type StopVms_Vms struct {
	CurrentState  string `json:"CurrentState,omitempty"`
	PreviousState string `json:"PreviousState,omitempty"`
	VmId          string `json:"VmId,omitempty"`
}

// implements the service definition of Subnet
type Subnet struct {
	AvailableIpsCount int64  `json:"AvailableIpsCount,omitempty"`
	IpRange           string `json:"IpRange,omitempty"`
	NetId             string `json:"NetId,omitempty"`
	State             string `json:"State,omitempty"`
	SubRegionName     string `json:"SubRegionName,omitempty"`
	SubnetId          string `json:"SubnetId,omitempty"`
	Tags              []Tags `json:"Tags,omitempty"`
}

// implements the service definition of Subnets
type Subnets struct {
	AvailableIpsCount int64  `json:"AvailableIpsCount,omitempty"`
	IpRange           string `json:"IpRange,omitempty"`
	NetId             string `json:"NetId,omitempty"`
	State             string `json:"State,omitempty"`
	SubRegionName     string `json:"SubRegionName,omitempty"`
	SubnetId          string `json:"SubnetId,omitempty"`
	Tags              []Tags `json:"Tags,omitempty"`
}

// implements the service definition of Transition
type Transition struct {
	Code    string `json:"Code,omitempty"`
	Message string `json:"Message,omitempty"`
}

// implements the service definition of UnlinkNetInternetGatewayRequest
type UnlinkNetInternetGatewayRequest struct {
	DryRun               bool   `json:"DryRun,omitempty"`
	NetId                string `json:"NetId,omitempty"`
	NetInternetGatewayId string `json:"NetInternetGatewayId,omitempty"`
}

// implements the service definition of UnlinkNetInternetGatewayResponse
type UnlinkNetInternetGatewayResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of UnlinkNicRequest
type UnlinkNicRequest struct {
	DryRun    bool   `json:"DryRun,omitempty"`
	NicLinkId string `json:"NicLinkId,omitempty"`
}

// implements the service definition of UnlinkNicResponse
type UnlinkNicResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of UnlinkPrivateIpsRequest
type UnlinkPrivateIpsRequest struct {
	DryRun     bool     `json:"DryRun,omitempty"`
	NicId      string   `json:"NicId,omitempty"`
	PrivateIps []string `json:"PrivateIps,omitempty"`
}

// implements the service definition of UnlinkPrivateIpsResponse
type UnlinkPrivateIpsResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of UnlinkPublicIpRequest
type UnlinkPublicIpRequest struct {
	DryRun   bool   `json:"DryRun,omitempty"`
	LinkId   string `json:"LinkId,omitempty"`
	PublicIp string `json:"PublicIp,omitempty"`
}

// implements the service definition of UnlinkPublicIpResponse
type UnlinkPublicIpResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of UnlinkRouteTableRequest
type UnlinkRouteTableRequest struct {
	DryRun bool   `json:"DryRun,omitempty"`
	LinkId string `json:"LinkId,omitempty"`
}

// implements the service definition of UnlinkRouteTableResponse
type UnlinkRouteTableResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of UnlinkVolumeRequest
type UnlinkVolumeRequest struct {
	DeviceName  string `json:"DeviceName,omitempty"`
	DryRun      bool   `json:"DryRun,omitempty"`
	ForceUnlink bool   `json:"ForceUnlink,omitempty"`
	VmId        string `json:"VmId,omitempty"`
	VolumeId    string `json:"VolumeId,omitempty"`
}

// implements the service definition of UnlinkVolumeResponse
type UnlinkVolumeResponse struct {
	DeleteOnVmDeletion bool            `json:"DeleteOnVmDeletion,omitempty"`
	DeviceName         string          `json:"DeviceName,omitempty"`
	ResponseContext    ResponseContext `json:"ResponseContext,omitempty"`
	State              string          `json:"State,omitempty"`
	VmId               string          `json:"VmId,omitempty"`
	VolumeId           string          `json:"VolumeId,omitempty"`
}

// implements the service definition of UnlinkVpnGatewayRequest
type UnlinkVpnGatewayRequest struct {
	DryRun       bool   `json:"DryRun,omitempty"`
	NetId        string `json:"NetId,omitempty"`
	VpnGatewayId string `json:"VpnGatewayId,omitempty"`
}

// implements the service definition of UnlinkVpnGatewayResponse
type UnlinkVpnGatewayResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of UpdateImageRequest
type UpdateImageRequest struct {
	DryRun             bool                           `json:"DryRun,omitempty"`
	ImageId            string                         `json:"ImageId,omitempty"`
	PermissionToLaunch UpdateImage_PermissionToLaunch `json:"PermissionToLaunch,omitempty"`
}

// implements the service definition of UpdateImageResponse
type UpdateImageResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of UpdateImage_PermissionToLaunch
type UpdateImage_PermissionToLaunch struct {
	Addition Addition `json:"Addition,omitempty"`
	Removal  Removal  `json:"Removal,omitempty"`
}

// implements the service definition of UpdateRoutePropagationRequest
type UpdateRoutePropagationRequest struct {
	DryRun       bool   `json:"DryRun,omitempty"`
	Enable       bool   `json:"Enable,omitempty"`
	RouteTableId string `json:"RouteTableId,omitempty"`
	VpnGatewayId string `json:"VpnGatewayId,omitempty"`
}

// implements the service definition of UpdateRoutePropagationResponse
type UpdateRoutePropagationResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of UpdateRouteRequest
type UpdateRouteRequest struct {
	DestinationIpRange string `json:"DestinationIpRange,omitempty"`
	DryRun             bool   `json:"DryRun,omitempty"`
	GatewayId          string `json:"GatewayId,omitempty"`
	NatServiceId       string `json:"NatServiceId,omitempty"`
	NetPeeringId       string `json:"NetPeeringId,omitempty"`
	NicId              string `json:"NicId,omitempty"`
	RouteTableId       string `json:"RouteTableId,omitempty"`
	VmId               string `json:"VmId,omitempty"`
}

// implements the service definition of UpdateRouteResponse
type UpdateRouteResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of UpdateRouteTableLinkRequest
type UpdateRouteTableLinkRequest struct {
	DryRun       bool   `json:"DryRun,omitempty"`
	LinkId       string `json:"LinkId,omitempty"`
	RouteTableId string `json:"RouteTableId,omitempty"`
}

// implements the service definition of UpdateRouteTableLinkResponse
type UpdateRouteTableLinkResponse struct {
	LinkId          string          `json:"LinkId,omitempty"`
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of UpdateSnapshotRequest
type UpdateSnapshotRequest struct {
	DryRun                   bool                                    `json:"DryRun,omitempty"`
	PermissionToCreateVolume UpdateSnapshot_PermissionToCreateVolume `json:"PermissionToCreateVolume,omitempty"`
	SnapshotId               string                                  `json:"SnapshotId,omitempty"`
}

// implements the service definition of UpdateSnapshotResponse
type UpdateSnapshotResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of UpdateSnapshot_PermissionToCreateVolume
type UpdateSnapshot_PermissionToCreateVolume struct {
	Addition Addition `json:"Addition,omitempty"`
	Removal  Removal  `json:"Removal,omitempty"`
}

// implements the service definition of UpdateVmAttributeRequest
type UpdateVmAttributeRequest struct {
	BlockDeviceMappings         []UpdateVmAttribute_BlockDeviceMappings `json:"BlockDeviceMappings,omitempty"`
	BsuOptimized                bool                                    `json:"BsuOptimized,omitempty"`
	DeletionProtection          bool                                    `json:"DeletionProtection,omitempty"`
	DryRun                      bool                                    `json:"DryRun,omitempty"`
	FirewallRulesSetIds         []string                                `json:"FirewallRulesSetIds,omitempty"`
	IsSourceDestChecked         bool                                    `json:"IsSourceDestChecked,omitempty"`
	KeypairName                 string                                  `json:"KeypairName,omitempty"`
	Type                        string                                  `json:"Type,omitempty"`
	UserData                    string                                  `json:"UserData,omitempty"`
	VmId                        string                                  `json:"VmId,omitempty"`
	VmInitiatedShutdownBehavior string                                  `json:"VmInitiatedShutdownBehavior,omitempty"`
}

// implements the service definition of UpdateVmAttributeResponse
type UpdateVmAttributeResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of UpdateVmAttribute_BlockDeviceMappings
type UpdateVmAttribute_BlockDeviceMappings struct {
	Bsu               UpdateVmAttribute_Bsu `json:"Bsu,omitempty"`
	DeviceName        string                `json:"DeviceName,omitempty"`
	NoDevice          string                `json:"NoDevice,omitempty"`
	VirtualDeviceName string                `json:"VirtualDeviceName,omitempty"`
}

// implements the service definition of UpdateVmAttribute_Bsu
type UpdateVmAttribute_Bsu struct {
	DeleteOnVmDeletion bool   `json:"DeleteOnVmDeletion,omitempty"`
	VolumeId           string `json:"VolumeId,omitempty"`
}

// implements the service definition of VmStates
type VmStates struct {
	MaintenanceEvents []MaintenanceEvents `json:"MaintenanceEvents,omitempty"`
	SubRegionName     string              `json:"SubRegionName,omitempty"`
	VmId              string              `json:"VmId,omitempty"`
	VmState           string              `json:"VmState,omitempty"`
}

// implements the service definition of Volumes
type Volumes struct {
	Iops          int64           `json:"Iops,omitempty"`
	LinkedVolumes []LinkedVolumes `json:"LinkedVolumes,omitempty"`
	Size          int64           `json:"Size,omitempty"`
	SnapshotId    string          `json:"SnapshotId,omitempty"`
	State         string          `json:"State,omitempty"`
	SubRegionName string          `json:"SubRegionName,omitempty"`
	Tags          []Tags          `json:"Tags,omitempty"`
	Type          string          `json:"Type,omitempty"`
	VolumeId      string          `json:"VolumeId,omitempty"`
}

// implements the service definition of VpnGateway
type VpnGateway struct {
	NetToVpnGatewayLinks []NetToVpnGatewayLinks `json:"NetToVpnGatewayLinks,omitempty"`
	State                string                 `json:"State,omitempty"`
	Tags                 []Tags                 `json:"Tags,omitempty"`
	Type                 string                 `json:"Type,omitempty"`
	VpnGatewayId         string                 `json:"VpnGatewayId,omitempty"`
}

// implements the service definition of VpnGateways
type VpnGateways struct {
	NetToVpnGatewayLinks []NetToVpnGatewayLinks `json:"NetToVpnGatewayLinks,omitempty"`
	State                string                 `json:"State,omitempty"`
	Tags                 []Tags                 `json:"Tags,omitempty"`
	Type                 string                 `json:"Type,omitempty"`
	VpnGatewayId         string                 `json:"VpnGatewayId,omitempty"`
}

// POST_AcceptNetPeeringParameters holds parameters to POST_AcceptNetPeering
type POST_AcceptNetPeeringParameters struct {
	Acceptnetpeeringrequest AcceptNetPeeringRequest `json:"acceptnetpeeringrequest,omitempty"`
}

// POST_AcceptNetPeeringResponses holds responses of POST_AcceptNetPeering
type POST_AcceptNetPeeringResponses struct {
	OK *AcceptNetPeeringResponse
}

// POST_CancelExportTaskParameters holds parameters to POST_CancelExportTask
type POST_CancelExportTaskParameters struct {
	Cancelexporttaskrequest CancelExportTaskRequest `json:"cancelexporttaskrequest,omitempty"`
}

// POST_CancelExportTaskResponses holds responses of POST_CancelExportTask
type POST_CancelExportTaskResponses struct {
	OK *CancelExportTaskResponse
}

// POST_CopyImageParameters holds parameters to POST_CopyImage
type POST_CopyImageParameters struct {
	Copyimagerequest CopyImageRequest `json:"copyimagerequest,omitempty"`
}

// POST_CopyImageResponses holds responses of POST_CopyImage
type POST_CopyImageResponses struct {
	OK *CopyImageResponse
}

// POST_CopySnapshotParameters holds parameters to POST_CopySnapshot
type POST_CopySnapshotParameters struct {
	Copysnapshotrequest CopySnapshotRequest `json:"copysnapshotrequest,omitempty"`
}

// POST_CopySnapshotResponses holds responses of POST_CopySnapshot
type POST_CopySnapshotResponses struct {
	OK *CopySnapshotResponse
}

// POST_CreateFirewallRuleInboundParameters holds parameters to POST_CreateFirewallRuleInbound
type POST_CreateFirewallRuleInboundParameters struct {
	Createfirewallruleinboundrequest CreateFirewallRuleInboundRequest `json:"createfirewallruleinboundrequest,omitempty"`
}

// POST_CreateFirewallRuleInboundResponses holds responses of POST_CreateFirewallRuleInbound
type POST_CreateFirewallRuleInboundResponses struct {
	OK *CreateFirewallRuleInboundResponse
}

// POST_CreateFirewallRuleOutboundParameters holds parameters to POST_CreateFirewallRuleOutbound
type POST_CreateFirewallRuleOutboundParameters struct {
	Createfirewallruleoutboundrequest CreateFirewallRuleOutboundRequest `json:"createfirewallruleoutboundrequest,omitempty"`
}

// POST_CreateFirewallRuleOutboundResponses holds responses of POST_CreateFirewallRuleOutbound
type POST_CreateFirewallRuleOutboundResponses struct {
	OK *CreateFirewallRuleOutboundResponse
}

// POST_CreateFirewallRulesSetParameters holds parameters to POST_CreateFirewallRulesSet
type POST_CreateFirewallRulesSetParameters struct {
	Createfirewallrulessetrequest CreateFirewallRulesSetRequest `json:"createfirewallrulessetrequest,omitempty"`
}

// POST_CreateFirewallRulesSetResponses holds responses of POST_CreateFirewallRulesSet
type POST_CreateFirewallRulesSetResponses struct {
	OK *CreateFirewallRulesSetResponse
}

// POST_CreateImageParameters holds parameters to POST_CreateImage
type POST_CreateImageParameters struct {
	Createimagerequest CreateImageRequest `json:"createimagerequest,omitempty"`
}

// POST_CreateImageResponses holds responses of POST_CreateImage
type POST_CreateImageResponses struct {
	OK *CreateImageResponse
}

// POST_CreateImageExportTaskParameters holds parameters to POST_CreateImageExportTask
type POST_CreateImageExportTaskParameters struct {
	Createimageexporttaskrequest CreateImageExportTaskRequest `json:"createimageexporttaskrequest,omitempty"`
}

// POST_CreateImageExportTaskResponses holds responses of POST_CreateImageExportTask
type POST_CreateImageExportTaskResponses struct {
	OK *CreateImageExportTaskResponse
}

// POST_CreateKeypairParameters holds parameters to POST_CreateKeypair
type POST_CreateKeypairParameters struct {
	Createkeypairrequest CreateKeypairRequest `json:"createkeypairrequest,omitempty"`
}

// POST_CreateKeypairResponses holds responses of POST_CreateKeypair
type POST_CreateKeypairResponses struct {
	OK *CreateKeypairResponse
}

// POST_CreateNatServiceParameters holds parameters to POST_CreateNatService
type POST_CreateNatServiceParameters struct {
	Createnatservicerequest CreateNatServiceRequest `json:"createnatservicerequest,omitempty"`
}

// POST_CreateNatServiceResponses holds responses of POST_CreateNatService
type POST_CreateNatServiceResponses struct {
	OK *CreateNatServiceResponse
}

// POST_CreateNetParameters holds parameters to POST_CreateNet
type POST_CreateNetParameters struct {
	Createnetrequest CreateNetRequest `json:"createnetrequest,omitempty"`
}

// POST_CreateNetResponses holds responses of POST_CreateNet
type POST_CreateNetResponses struct {
	OK *CreateNetResponse
}

// POST_CreateNetInternetGatewayParameters holds parameters to POST_CreateNetInternetGateway
type POST_CreateNetInternetGatewayParameters struct {
	Createnetinternetgatewayrequest CreateNetInternetGatewayRequest `json:"createnetinternetgatewayrequest,omitempty"`
}

// POST_CreateNetInternetGatewayResponses holds responses of POST_CreateNetInternetGateway
type POST_CreateNetInternetGatewayResponses struct {
	OK *CreateNetInternetGatewayResponse
}

// POST_CreateNetPeeringParameters holds parameters to POST_CreateNetPeering
type POST_CreateNetPeeringParameters struct {
	Createnetpeeringrequest CreateNetPeeringRequest `json:"createnetpeeringrequest,omitempty"`
}

// POST_CreateNetPeeringResponses holds responses of POST_CreateNetPeering
type POST_CreateNetPeeringResponses struct {
	OK *CreateNetPeeringResponse
}

// POST_CreateNicParameters holds parameters to POST_CreateNic
type POST_CreateNicParameters struct {
	Createnicrequest CreateNicRequest `json:"createnicrequest,omitempty"`
}

// POST_CreateNicResponses holds responses of POST_CreateNic
type POST_CreateNicResponses struct {
	OK *CreateNicResponse
}

// POST_CreatePublicIpParameters holds parameters to POST_CreatePublicIp
type POST_CreatePublicIpParameters struct {
	Createpubliciprequest CreatePublicIpRequest `json:"createpubliciprequest,omitempty"`
}

// POST_CreatePublicIpResponses holds responses of POST_CreatePublicIp
type POST_CreatePublicIpResponses struct {
	OK *CreatePublicIpResponse
}

// POST_CreateRouteParameters holds parameters to POST_CreateRoute
type POST_CreateRouteParameters struct {
	Createrouterequest CreateRouteRequest `json:"createrouterequest,omitempty"`
}

// POST_CreateRouteResponses holds responses of POST_CreateRoute
type POST_CreateRouteResponses struct {
	OK *CreateRouteResponse
}

// POST_CreateRouteTableParameters holds parameters to POST_CreateRouteTable
type POST_CreateRouteTableParameters struct {
	Createroutetablerequest CreateRouteTableRequest `json:"createroutetablerequest,omitempty"`
}

// POST_CreateRouteTableResponses holds responses of POST_CreateRouteTable
type POST_CreateRouteTableResponses struct {
	OK *CreateRouteTableResponse
}

// POST_CreateSnapshotParameters holds parameters to POST_CreateSnapshot
type POST_CreateSnapshotParameters struct {
	Createsnapshotrequest CreateSnapshotRequest `json:"createsnapshotrequest,omitempty"`
}

// POST_CreateSnapshotResponses holds responses of POST_CreateSnapshot
type POST_CreateSnapshotResponses struct {
	OK *CreateSnapshotResponse
}

// POST_CreateSnapshotExportTaskParameters holds parameters to POST_CreateSnapshotExportTask
type POST_CreateSnapshotExportTaskParameters struct {
	Createsnapshotexporttaskrequest CreateSnapshotExportTaskRequest `json:"createsnapshotexporttaskrequest,omitempty"`
}

// POST_CreateSnapshotExportTaskResponses holds responses of POST_CreateSnapshotExportTask
type POST_CreateSnapshotExportTaskResponses struct {
	OK *CreateSnapshotExportTaskResponse
}

// POST_CreateSubnetParameters holds parameters to POST_CreateSubnet
type POST_CreateSubnetParameters struct {
	Createsubnetrequest CreateSubnetRequest `json:"createsubnetrequest,omitempty"`
}

// POST_CreateSubnetResponses holds responses of POST_CreateSubnet
type POST_CreateSubnetResponses struct {
	OK *CreateSubnetResponse
}

// POST_CreateTagsParameters holds parameters to POST_CreateTags
type POST_CreateTagsParameters struct {
	Createtagsrequest CreateTagsRequest `json:"createtagsrequest,omitempty"`
}

// POST_CreateTagsResponses holds responses of POST_CreateTags
type POST_CreateTagsResponses struct {
	OK *CreateTagsResponse
}

// POST_CreateVmsParameters holds parameters to POST_CreateVms
type POST_CreateVmsParameters struct {
	Createvmsrequest CreateVmsRequest `json:"createvmsrequest,omitempty"`
}

// POST_CreateVmsResponses holds responses of POST_CreateVms
type POST_CreateVmsResponses struct {
	OK *CreateVmsResponse
}

// POST_CreateVolumeParameters holds parameters to POST_CreateVolume
type POST_CreateVolumeParameters struct {
	Createvolumerequest CreateVolumeRequest `json:"createvolumerequest,omitempty"`
}

// POST_CreateVolumeResponses holds responses of POST_CreateVolume
type POST_CreateVolumeResponses struct {
	OK *CreateVolumeResponse
}

// POST_CreateVpnGatewayParameters holds parameters to POST_CreateVpnGateway
type POST_CreateVpnGatewayParameters struct {
	Createvpngatewayrequest CreateVpnGatewayRequest `json:"createvpngatewayrequest,omitempty"`
}

// POST_CreateVpnGatewayResponses holds responses of POST_CreateVpnGateway
type POST_CreateVpnGatewayResponses struct {
	OK *CreateVpnGatewayResponse
}

// POST_DeleteFirewallRuleInboundParameters holds parameters to POST_DeleteFirewallRuleInbound
type POST_DeleteFirewallRuleInboundParameters struct {
	Deletefirewallruleinboundrequest DeleteFirewallRuleInboundRequest `json:"deletefirewallruleinboundrequest,omitempty"`
}

// POST_DeleteFirewallRuleInboundResponses holds responses of POST_DeleteFirewallRuleInbound
type POST_DeleteFirewallRuleInboundResponses struct {
	OK *DeleteFirewallRuleInboundResponse
}

// POST_DeleteFirewallRuleOutboundParameters holds parameters to POST_DeleteFirewallRuleOutbound
type POST_DeleteFirewallRuleOutboundParameters struct {
	Deletefirewallruleoutboundrequest DeleteFirewallRuleOutboundRequest `json:"deletefirewallruleoutboundrequest,omitempty"`
}

// POST_DeleteFirewallRuleOutboundResponses holds responses of POST_DeleteFirewallRuleOutbound
type POST_DeleteFirewallRuleOutboundResponses struct {
	OK *DeleteFirewallRuleOutboundResponse
}

// POST_DeleteFirewallRulesSetParameters holds parameters to POST_DeleteFirewallRulesSet
type POST_DeleteFirewallRulesSetParameters struct {
	Deletefirewallrulessetrequest DeleteFirewallRulesSetRequest `json:"deletefirewallrulessetrequest,omitempty"`
}

// POST_DeleteFirewallRulesSetResponses holds responses of POST_DeleteFirewallRulesSet
type POST_DeleteFirewallRulesSetResponses struct {
	OK *DeleteFirewallRulesSetResponse
}

// POST_DeleteKeypairParameters holds parameters to POST_DeleteKeypair
type POST_DeleteKeypairParameters struct {
	Deletekeypairrequest DeleteKeypairRequest `json:"deletekeypairrequest,omitempty"`
}

// POST_DeleteKeypairResponses holds responses of POST_DeleteKeypair
type POST_DeleteKeypairResponses struct {
	OK *DeleteKeypairResponse
}

// POST_DeleteNatServiceParameters holds parameters to POST_DeleteNatService
type POST_DeleteNatServiceParameters struct {
	Deletenatservicerequest DeleteNatServiceRequest `json:"deletenatservicerequest,omitempty"`
}

// POST_DeleteNatServiceResponses holds responses of POST_DeleteNatService
type POST_DeleteNatServiceResponses struct {
	OK *DeleteNatServiceResponse
}

// POST_DeleteNetParameters holds parameters to POST_DeleteNet
type POST_DeleteNetParameters struct {
	Deletenetrequest DeleteNetRequest `json:"deletenetrequest,omitempty"`
}

// POST_DeleteNetResponses holds responses of POST_DeleteNet
type POST_DeleteNetResponses struct {
	OK *DeleteNetResponse
}

// POST_DeleteNetInternetGatewayParameters holds parameters to POST_DeleteNetInternetGateway
type POST_DeleteNetInternetGatewayParameters struct {
	Deletenetinternetgatewayrequest DeleteNetInternetGatewayRequest `json:"deletenetinternetgatewayrequest,omitempty"`
}

// POST_DeleteNetInternetGatewayResponses holds responses of POST_DeleteNetInternetGateway
type POST_DeleteNetInternetGatewayResponses struct {
	OK *DeleteNetInternetGatewayResponse
}

// POST_DeleteNetPeeringParameters holds parameters to POST_DeleteNetPeering
type POST_DeleteNetPeeringParameters struct {
	Deletenetpeeringrequest DeleteNetPeeringRequest `json:"deletenetpeeringrequest,omitempty"`
}

// POST_DeleteNetPeeringResponses holds responses of POST_DeleteNetPeering
type POST_DeleteNetPeeringResponses struct {
	OK *DeleteNetPeeringResponse
}

// POST_DeleteNicParameters holds parameters to POST_DeleteNic
type POST_DeleteNicParameters struct {
	Deletenicrequest DeleteNicRequest `json:"deletenicrequest,omitempty"`
}

// POST_DeleteNicResponses holds responses of POST_DeleteNic
type POST_DeleteNicResponses struct {
	OK *DeleteNicResponse
}

// POST_DeletePublicIpParameters holds parameters to POST_DeletePublicIp
type POST_DeletePublicIpParameters struct {
	Deletepubliciprequest DeletePublicIpRequest `json:"deletepubliciprequest,omitempty"`
}

// POST_DeletePublicIpResponses holds responses of POST_DeletePublicIp
type POST_DeletePublicIpResponses struct {
	OK *DeletePublicIpResponse
}

// POST_DeleteRouteParameters holds parameters to POST_DeleteRoute
type POST_DeleteRouteParameters struct {
	Deleterouterequest DeleteRouteRequest `json:"deleterouterequest,omitempty"`
}

// POST_DeleteRouteResponses holds responses of POST_DeleteRoute
type POST_DeleteRouteResponses struct {
	OK *DeleteRouteResponse
}

// POST_DeleteRouteTableParameters holds parameters to POST_DeleteRouteTable
type POST_DeleteRouteTableParameters struct {
	Deleteroutetablerequest DeleteRouteTableRequest `json:"deleteroutetablerequest,omitempty"`
}

// POST_DeleteRouteTableResponses holds responses of POST_DeleteRouteTable
type POST_DeleteRouteTableResponses struct {
	OK *DeleteRouteTableResponse
}

// POST_DeleteSnapshotParameters holds parameters to POST_DeleteSnapshot
type POST_DeleteSnapshotParameters struct {
	Deletesnapshotrequest DeleteSnapshotRequest `json:"deletesnapshotrequest,omitempty"`
}

// POST_DeleteSnapshotResponses holds responses of POST_DeleteSnapshot
type POST_DeleteSnapshotResponses struct {
	OK *DeleteSnapshotResponse
}

// POST_DeleteSubnetParameters holds parameters to POST_DeleteSubnet
type POST_DeleteSubnetParameters struct {
	Deletesubnetrequest DeleteSubnetRequest `json:"deletesubnetrequest,omitempty"`
}

// POST_DeleteSubnetResponses holds responses of POST_DeleteSubnet
type POST_DeleteSubnetResponses struct {
	OK *DeleteSubnetResponse
}

// POST_DeleteTagsParameters holds parameters to POST_DeleteTags
type POST_DeleteTagsParameters struct {
	Deletetagsrequest DeleteTagsRequest `json:"deletetagsrequest,omitempty"`
}

// POST_DeleteTagsResponses holds responses of POST_DeleteTags
type POST_DeleteTagsResponses struct {
	OK *DeleteTagsResponse
}

// POST_DeleteVmsParameters holds parameters to POST_DeleteVms
type POST_DeleteVmsParameters struct {
	Deletevmsrequest DeleteVmsRequest `json:"deletevmsrequest,omitempty"`
}

// POST_DeleteVmsResponses holds responses of POST_DeleteVms
type POST_DeleteVmsResponses struct {
	OK *DeleteVmsResponse
}

// POST_DeleteVolumeParameters holds parameters to POST_DeleteVolume
type POST_DeleteVolumeParameters struct {
	Deletevolumerequest DeleteVolumeRequest `json:"deletevolumerequest,omitempty"`
}

// POST_DeleteVolumeResponses holds responses of POST_DeleteVolume
type POST_DeleteVolumeResponses struct {
	OK *DeleteVolumeResponse
}

// POST_DeleteVpnGatewayParameters holds parameters to POST_DeleteVpnGateway
type POST_DeleteVpnGatewayParameters struct {
	Deletevpngatewayrequest DeleteVpnGatewayRequest `json:"deletevpngatewayrequest,omitempty"`
}

// POST_DeleteVpnGatewayResponses holds responses of POST_DeleteVpnGateway
type POST_DeleteVpnGatewayResponses struct {
	OK *DeleteVpnGatewayResponse
}

// POST_DeregisterImageParameters holds parameters to POST_DeregisterImage
type POST_DeregisterImageParameters struct {
	Deregisterimagerequest DeregisterImageRequest `json:"deregisterimagerequest,omitempty"`
}

// POST_DeregisterImageResponses holds responses of POST_DeregisterImage
type POST_DeregisterImageResponses struct {
	OK *DeregisterImageResponse
}

// POST_ImportSnapshotParameters holds parameters to POST_ImportSnapshot
type POST_ImportSnapshotParameters struct {
	Importsnapshotrequest ImportSnapshotRequest `json:"importsnapshotrequest,omitempty"`
}

// POST_ImportSnapshotResponses holds responses of POST_ImportSnapshot
type POST_ImportSnapshotResponses struct {
	OK *ImportSnapshotResponse
}

// POST_LinkNetInternetGatewayParameters holds parameters to POST_LinkNetInternetGateway
type POST_LinkNetInternetGatewayParameters struct {
	Linknetinternetgatewayrequest LinkNetInternetGatewayRequest `json:"linknetinternetgatewayrequest,omitempty"`
}

// POST_LinkNetInternetGatewayResponses holds responses of POST_LinkNetInternetGateway
type POST_LinkNetInternetGatewayResponses struct {
	OK *LinkNetInternetGatewayResponse
}

// POST_LinkNicParameters holds parameters to POST_LinkNic
type POST_LinkNicParameters struct {
	Linknicrequest LinkNicRequest `json:"linknicrequest,omitempty"`
}

// POST_LinkNicResponses holds responses of POST_LinkNic
type POST_LinkNicResponses struct {
	OK *LinkNicResponse
}

// POST_LinkPublicIpParameters holds parameters to POST_LinkPublicIp
type POST_LinkPublicIpParameters struct {
	Linkpubliciprequest LinkPublicIpRequest `json:"linkpubliciprequest,omitempty"`
}

// POST_LinkPublicIpResponses holds responses of POST_LinkPublicIp
type POST_LinkPublicIpResponses struct {
	OK *LinkPublicIpResponse
}

// POST_LinkRouteTableParameters holds parameters to POST_LinkRouteTable
type POST_LinkRouteTableParameters struct {
	Linkroutetablerequest LinkRouteTableRequest `json:"linkroutetablerequest,omitempty"`
}

// POST_LinkRouteTableResponses holds responses of POST_LinkRouteTable
type POST_LinkRouteTableResponses struct {
	OK *LinkRouteTableResponse
}

// POST_LinkVolumeParameters holds parameters to POST_LinkVolume
type POST_LinkVolumeParameters struct {
	Linkvolumerequest LinkVolumeRequest `json:"linkvolumerequest,omitempty"`
}

// POST_LinkVolumeResponses holds responses of POST_LinkVolume
type POST_LinkVolumeResponses struct {
	OK *LinkVolumeResponse
}

// POST_LinkVpnGatewayParameters holds parameters to POST_LinkVpnGateway
type POST_LinkVpnGatewayParameters struct {
	Linkvpngatewayrequest LinkVpnGatewayRequest `json:"linkvpngatewayrequest,omitempty"`
}

// POST_LinkVpnGatewayResponses holds responses of POST_LinkVpnGateway
type POST_LinkVpnGatewayResponses struct {
	OK *LinkVpnGatewayResponse
}

// POST_ReadFirewallRulesSetsParameters holds parameters to POST_ReadFirewallRulesSets
type POST_ReadFirewallRulesSetsParameters struct {
	Readfirewallrulessetsrequest ReadFirewallRulesSetsRequest `json:"readfirewallrulessetsrequest,omitempty"`
}

// POST_ReadFirewallRulesSetsResponses holds responses of POST_ReadFirewallRulesSets
type POST_ReadFirewallRulesSetsResponses struct {
	OK *ReadFirewallRulesSetsResponse
}

// POST_ReadImageExportTasksParameters holds parameters to POST_ReadImageExportTasks
type POST_ReadImageExportTasksParameters struct {
	Readimageexporttasksrequest ReadImageExportTasksRequest `json:"readimageexporttasksrequest,omitempty"`
}

// POST_ReadImageExportTasksResponses holds responses of POST_ReadImageExportTasks
type POST_ReadImageExportTasksResponses struct {
	OK *ReadImageExportTasksResponse
}

// POST_ReadImagesParameters holds parameters to POST_ReadImages
type POST_ReadImagesParameters struct {
	Readimagesrequest ReadImagesRequest `json:"readimagesrequest,omitempty"`
}

// POST_ReadImagesResponses holds responses of POST_ReadImages
type POST_ReadImagesResponses struct {
	OK *ReadImagesResponse
}

// POST_ReadKeypairsParameters holds parameters to POST_ReadKeypairs
type POST_ReadKeypairsParameters struct {
	Readkeypairsrequest ReadKeypairsRequest `json:"readkeypairsrequest,omitempty"`
}

// POST_ReadKeypairsResponses holds responses of POST_ReadKeypairs
type POST_ReadKeypairsResponses struct {
	OK *ReadKeypairsResponse
}

// POST_ReadNatServicesParameters holds parameters to POST_ReadNatServices
type POST_ReadNatServicesParameters struct {
	Readnatservicesrequest ReadNatServicesRequest `json:"readnatservicesrequest,omitempty"`
}

// POST_ReadNatServicesResponses holds responses of POST_ReadNatServices
type POST_ReadNatServicesResponses struct {
	OK *ReadNatServicesResponse
}

// POST_ReadNetInternetGatewaysParameters holds parameters to POST_ReadNetInternetGateways
type POST_ReadNetInternetGatewaysParameters struct {
	Readnetinternetgatewaysrequest ReadNetInternetGatewaysRequest `json:"readnetinternetgatewaysrequest,omitempty"`
}

// POST_ReadNetInternetGatewaysResponses holds responses of POST_ReadNetInternetGateways
type POST_ReadNetInternetGatewaysResponses struct {
	OK *ReadNetInternetGatewaysResponse
}

// POST_ReadNetPeeringsParameters holds parameters to POST_ReadNetPeerings
type POST_ReadNetPeeringsParameters struct {
	Readnetpeeringsrequest ReadNetPeeringsRequest `json:"readnetpeeringsrequest,omitempty"`
}

// POST_ReadNetPeeringsResponses holds responses of POST_ReadNetPeerings
type POST_ReadNetPeeringsResponses struct {
	OK *ReadNetPeeringsResponse
}

// POST_ReadNetsParameters holds parameters to POST_ReadNets
type POST_ReadNetsParameters struct {
	Readnetsrequest ReadNetsRequest `json:"readnetsrequest,omitempty"`
}

// POST_ReadNetsResponses holds responses of POST_ReadNets
type POST_ReadNetsResponses struct {
	OK *ReadNetsResponse
}

// POST_ReadNicsParameters holds parameters to POST_ReadNics
type POST_ReadNicsParameters struct {
	Readnicsrequest ReadNicsRequest `json:"readnicsrequest,omitempty"`
}

// POST_ReadNicsResponses holds responses of POST_ReadNics
type POST_ReadNicsResponses struct {
	OK *ReadNicsResponse
}

// POST_ReadPublicIpsParameters holds parameters to POST_ReadPublicIps
type POST_ReadPublicIpsParameters struct {
	Readpublicipsrequest ReadPublicIpsRequest `json:"readpublicipsrequest,omitempty"`
}

// POST_ReadPublicIpsResponses holds responses of POST_ReadPublicIps
type POST_ReadPublicIpsResponses struct {
	OK *ReadPublicIpsResponse
}

// POST_ReadRouteTablesParameters holds parameters to POST_ReadRouteTables
type POST_ReadRouteTablesParameters struct {
	Readroutetablesrequest ReadRouteTablesRequest `json:"readroutetablesrequest,omitempty"`
}

// POST_ReadRouteTablesResponses holds responses of POST_ReadRouteTables
type POST_ReadRouteTablesResponses struct {
	OK *ReadRouteTablesResponse
}

// POST_ReadSnapshotExportTasksParameters holds parameters to POST_ReadSnapshotExportTasks
type POST_ReadSnapshotExportTasksParameters struct {
	Readsnapshotexporttasksrequest ReadSnapshotExportTasksRequest `json:"readsnapshotexporttasksrequest,omitempty"`
}

// POST_ReadSnapshotExportTasksResponses holds responses of POST_ReadSnapshotExportTasks
type POST_ReadSnapshotExportTasksResponses struct {
	OK *ReadSnapshotExportTasksResponse
}

// POST_ReadSnapshotsParameters holds parameters to POST_ReadSnapshots
type POST_ReadSnapshotsParameters struct {
	Readsnapshotsrequest ReadSnapshotsRequest `json:"readsnapshotsrequest,omitempty"`
}

// POST_ReadSnapshotsResponses holds responses of POST_ReadSnapshots
type POST_ReadSnapshotsResponses struct {
	OK *ReadSnapshotsResponse
}

// POST_ReadSubnetsParameters holds parameters to POST_ReadSubnets
type POST_ReadSubnetsParameters struct {
	Readsubnetsrequest ReadSubnetsRequest `json:"readsubnetsrequest,omitempty"`
}

// POST_ReadSubnetsResponses holds responses of POST_ReadSubnets
type POST_ReadSubnetsResponses struct {
	OK *ReadSubnetsResponse
}

// POST_ReadTagsParameters holds parameters to POST_ReadTags
type POST_ReadTagsParameters struct {
	Readtagsrequest ReadTagsRequest `json:"readtagsrequest,omitempty"`
}

// POST_ReadTagsResponses holds responses of POST_ReadTags
type POST_ReadTagsResponses struct {
	OK *ReadTagsResponse
}

// POST_ReadVmAttributeParameters holds parameters to POST_ReadVmAttribute
type POST_ReadVmAttributeParameters struct {
	Readvmattributerequest ReadVmAttributeRequest `json:"readvmattributerequest,omitempty"`
}

// POST_ReadVmAttributeResponses holds responses of POST_ReadVmAttribute
type POST_ReadVmAttributeResponses struct {
	OK *ReadVmAttributeResponse
}

// POST_ReadVmsParameters holds parameters to POST_ReadVms
type POST_ReadVmsParameters struct {
	Readvmsrequest ReadVmsRequest `json:"readvmsrequest,omitempty"`
}

// POST_ReadVmsResponses holds responses of POST_ReadVms
type POST_ReadVmsResponses struct {
	OK *ReadVmsResponse
}

// POST_ReadVmsStateParameters holds parameters to POST_ReadVmsState
type POST_ReadVmsStateParameters struct {
	Readvmsstaterequest ReadVmsStateRequest `json:"readvmsstaterequest,omitempty"`
}

// POST_ReadVmsStateResponses holds responses of POST_ReadVmsState
type POST_ReadVmsStateResponses struct {
	OK *ReadVmsStateResponse
}

// POST_ReadVolumesParameters holds parameters to POST_ReadVolumes
type POST_ReadVolumesParameters struct {
	Readvolumesrequest ReadVolumesRequest `json:"readvolumesrequest,omitempty"`
}

// POST_ReadVolumesResponses holds responses of POST_ReadVolumes
type POST_ReadVolumesResponses struct {
	OK *ReadVolumesResponse
}

// POST_ReadVpnGatewaysParameters holds parameters to POST_ReadVpnGateways
type POST_ReadVpnGatewaysParameters struct {
	Readvpngatewaysrequest ReadVpnGatewaysRequest `json:"readvpngatewaysrequest,omitempty"`
}

// POST_ReadVpnGatewaysResponses holds responses of POST_ReadVpnGateways
type POST_ReadVpnGatewaysResponses struct {
	OK *ReadVpnGatewaysResponse
}

// POST_RebootVmsParameters holds parameters to POST_RebootVms
type POST_RebootVmsParameters struct {
	Rebootvmsrequest RebootVmsRequest `json:"rebootvmsrequest,omitempty"`
}

// POST_RebootVmsResponses holds responses of POST_RebootVms
type POST_RebootVmsResponses struct {
	OK *RebootVmsResponse
}

// POST_RegisterImageParameters holds parameters to POST_RegisterImage
type POST_RegisterImageParameters struct {
	Registerimagerequest RegisterImageRequest `json:"registerimagerequest,omitempty"`
}

// POST_RegisterImageResponses holds responses of POST_RegisterImage
type POST_RegisterImageResponses struct {
	OK *RegisterImageResponse
}

// POST_RejectNetPeeringParameters holds parameters to POST_RejectNetPeering
type POST_RejectNetPeeringParameters struct {
	Rejectnetpeeringrequest RejectNetPeeringRequest `json:"rejectnetpeeringrequest,omitempty"`
}

// POST_RejectNetPeeringResponses holds responses of POST_RejectNetPeering
type POST_RejectNetPeeringResponses struct {
	OK *RejectNetPeeringResponse
}

// POST_StartVmsParameters holds parameters to POST_StartVms
type POST_StartVmsParameters struct {
	Startvmsrequest StartVmsRequest `json:"startvmsrequest,omitempty"`
}

// POST_StartVmsResponses holds responses of POST_StartVms
type POST_StartVmsResponses struct {
	OK *StartVmsResponse
}

// POST_StopVmsParameters holds parameters to POST_StopVms
type POST_StopVmsParameters struct {
	Stopvmsrequest StopVmsRequest `json:"stopvmsrequest,omitempty"`
}

// POST_StopVmsResponses holds responses of POST_StopVms
type POST_StopVmsResponses struct {
	OK *StopVmsResponse
}

// POST_UnlinkNetInternetGatewayParameters holds parameters to POST_UnlinkNetInternetGateway
type POST_UnlinkNetInternetGatewayParameters struct {
	Unlinknetinternetgatewayrequest UnlinkNetInternetGatewayRequest `json:"unlinknetinternetgatewayrequest,omitempty"`
}

// POST_UnlinkNetInternetGatewayResponses holds responses of POST_UnlinkNetInternetGateway
type POST_UnlinkNetInternetGatewayResponses struct {
	OK *UnlinkNetInternetGatewayResponse
}

// POST_UnlinkNicParameters holds parameters to POST_UnlinkNic
type POST_UnlinkNicParameters struct {
	Unlinknicrequest UnlinkNicRequest `json:"unlinknicrequest,omitempty"`
}

// POST_UnlinkNicResponses holds responses of POST_UnlinkNic
type POST_UnlinkNicResponses struct {
	OK *UnlinkNicResponse
}

// POST_UnlinkPrivateIpsParameters holds parameters to POST_UnlinkPrivateIps
type POST_UnlinkPrivateIpsParameters struct {
	Unlinkprivateipsrequest UnlinkPrivateIpsRequest `json:"unlinkprivateipsrequest,omitempty"`
}

// POST_UnlinkPrivateIpsResponses holds responses of POST_UnlinkPrivateIps
type POST_UnlinkPrivateIpsResponses struct {
	OK *UnlinkPrivateIpsResponse
}

// POST_UnlinkPublicIpParameters holds parameters to POST_UnlinkPublicIp
type POST_UnlinkPublicIpParameters struct {
	Unlinkpubliciprequest UnlinkPublicIpRequest `json:"unlinkpubliciprequest,omitempty"`
}

// POST_UnlinkPublicIpResponses holds responses of POST_UnlinkPublicIp
type POST_UnlinkPublicIpResponses struct {
	OK *UnlinkPublicIpResponse
}

// POST_UnlinkRouteTableParameters holds parameters to POST_UnlinkRouteTable
type POST_UnlinkRouteTableParameters struct {
	Unlinkroutetablerequest UnlinkRouteTableRequest `json:"unlinkroutetablerequest,omitempty"`
}

// POST_UnlinkRouteTableResponses holds responses of POST_UnlinkRouteTable
type POST_UnlinkRouteTableResponses struct {
	OK *UnlinkRouteTableResponse
}

// POST_UnlinkVolumeParameters holds parameters to POST_UnlinkVolume
type POST_UnlinkVolumeParameters struct {
	Unlinkvolumerequest UnlinkVolumeRequest `json:"unlinkvolumerequest,omitempty"`
}

// POST_UnlinkVolumeResponses holds responses of POST_UnlinkVolume
type POST_UnlinkVolumeResponses struct {
	OK *UnlinkVolumeResponse
}

// POST_UnlinkVpnGatewayParameters holds parameters to POST_UnlinkVpnGateway
type POST_UnlinkVpnGatewayParameters struct {
	Unlinkvpngatewayrequest UnlinkVpnGatewayRequest `json:"unlinkvpngatewayrequest,omitempty"`
}

// POST_UnlinkVpnGatewayResponses holds responses of POST_UnlinkVpnGateway
type POST_UnlinkVpnGatewayResponses struct {
	OK *UnlinkVpnGatewayResponse
}

// POST_UpdateImageParameters holds parameters to POST_UpdateImage
type POST_UpdateImageParameters struct {
	Updateimagerequest UpdateImageRequest `json:"updateimagerequest,omitempty"`
}

// POST_UpdateImageResponses holds responses of POST_UpdateImage
type POST_UpdateImageResponses struct {
	OK *UpdateImageResponse
}

// POST_UpdateRouteParameters holds parameters to POST_UpdateRoute
type POST_UpdateRouteParameters struct {
	Updaterouterequest UpdateRouteRequest `json:"updaterouterequest,omitempty"`
}

// POST_UpdateRouteResponses holds responses of POST_UpdateRoute
type POST_UpdateRouteResponses struct {
	OK *UpdateRouteResponse
}

// POST_UpdateRoutePropagationParameters holds parameters to POST_UpdateRoutePropagation
type POST_UpdateRoutePropagationParameters struct {
	Updateroutepropagationrequest UpdateRoutePropagationRequest `json:"updateroutepropagationrequest,omitempty"`
}

// POST_UpdateRoutePropagationResponses holds responses of POST_UpdateRoutePropagation
type POST_UpdateRoutePropagationResponses struct {
	OK *UpdateRoutePropagationResponse
}

// POST_UpdateRouteTableLinkParameters holds parameters to POST_UpdateRouteTableLink
type POST_UpdateRouteTableLinkParameters struct {
	Updateroutetablelinkrequest UpdateRouteTableLinkRequest `json:"updateroutetablelinkrequest,omitempty"`
}

// POST_UpdateRouteTableLinkResponses holds responses of POST_UpdateRouteTableLink
type POST_UpdateRouteTableLinkResponses struct {
	OK *UpdateRouteTableLinkResponse
}

// POST_UpdateSnapshotParameters holds parameters to POST_UpdateSnapshot
type POST_UpdateSnapshotParameters struct {
	Updatesnapshotrequest UpdateSnapshotRequest `json:"updatesnapshotrequest,omitempty"`
}

// POST_UpdateSnapshotResponses holds responses of POST_UpdateSnapshot
type POST_UpdateSnapshotResponses struct {
	OK *UpdateSnapshotResponse
}

// POST_UpdateVmAttributeParameters holds parameters to POST_UpdateVmAttribute
type POST_UpdateVmAttributeParameters struct {
	Updatevmattributerequest UpdateVmAttributeRequest `json:"updatevmattributerequest,omitempty"`
}

// POST_UpdateVmAttributeResponses holds responses of POST_UpdateVmAttribute
type POST_UpdateVmAttributeResponses struct {
	OK *UpdateVmAttributeResponse
}
