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

// implements the service definition of AccessLog
type AccessLog struct {
	IsEnabled           bool   `json:"IsEnabled,omitempty"`
	OsuBucketName       string `json:"OsuBucketName,omitempty"`
	OsuBucketPrefix     string `json:"OsuBucketPrefix,omitempty"`
	PublicationInterval int64  `json:"PublicationInterval,omitempty"`
}

// implements the service definition of Account
type Account struct {
	AccountId     string `json:"AccountId,omitempty"`
	City          string `json:"City,omitempty"`
	CompanyName   string `json:"CompanyName,omitempty"`
	Country       string `json:"Country,omitempty"`
	CustomerId    string `json:"CustomerId,omitempty"`
	Email         string `json:"Email,omitempty"`
	FirstName     string `json:"FirstName,omitempty"`
	JobTitle      string `json:"JobTitle,omitempty"`
	LastName      string `json:"LastName,omitempty"`
	Mobile        string `json:"Mobile,omitempty"`
	Phone         string `json:"Phone,omitempty"`
	StateProvince string `json:"StateProvince,omitempty"`
	VatNumber     string `json:"VatNumber,omitempty"`
	ZipCode       string `json:"ZipCode,omitempty"`
}

// implements the service definition of Additions
type Additions struct {
	AccountId        string `json:"AccountId,omitempty"`
	GlobalPermission string `json:"GlobalPermission,omitempty"`
}

// implements the service definition of ApiKey
type ApiKey struct {
	AccountId string `json:"AccountId,omitempty"`
	ApiKeyId  string `json:"ApiKeyId,omitempty"`
	SecretKey string `json:"SecretKey,omitempty"`
	State     string `json:"State,omitempty"`
	Tags      []Tags `json:"Tags,omitempty"`
	UserName  string `json:"UserName,omitempty"`
}

// implements the service definition of ApiKeys
type ApiKeys struct {
	ApiKeyId  string `json:"ApiKeyId,omitempty"`
	SecretKey string `json:"SecretKey,omitempty"`
	Tags      []Tags `json:"Tags,omitempty"`
}

// implements the service definition of ApplicationStickyCookiePolicies
type ApplicationStickyCookiePolicies struct {
	CookieName string `json:"CookieName,omitempty"`
	PolicyName string `json:"PolicyName,omitempty"`
}

// implements the service definition of AuthenticateAccountRequest
type AuthenticateAccountRequest struct {
	Login    string `json:"Login,omitempty"`
	Password string `json:"Password,omitempty"`
}

// implements the service definition of AuthenticateAccountResponse
type AuthenticateAccountResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of BackendVmsHealth
type BackendVmsHealth struct {
	Description string `json:"Description,omitempty"`
	State       string `json:"State,omitempty"`
	StateReason string `json:"StateReason,omitempty"`
	VmId        string `json:"VmId,omitempty"`
}

// implements the service definition of BlockDeviceMappings
type BlockDeviceMappings struct {
	Bsu        Bsu    `json:"Bsu,omitempty"`
	DeviceName string `json:"DeviceName,omitempty"`
}

// implements the service definition of Bsu
type Bsu struct {
	DeleteOnVmDeletion bool   `json:"DeleteOnVmDeletion,omitempty"`
	State              string `json:"State,omitempty"`
	VolumeId           string `json:"VolumeId,omitempty"`
}

// implements the service definition of CancelExportTaskRequest
type CancelExportTaskRequest struct {
	ExportTaskId string `json:"ExportTaskId,omitempty"`
}

// implements the service definition of CancelExportTaskResponse
type CancelExportTaskResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of Catalog
type Catalog struct {
	Domain           string `json:"Domain,omitempty"`
	Instance         string `json:"Instance,omitempty"`
	SourceRegionName string `json:"SourceRegionName,omitempty"`
	TargetRegionName string `json:"TargetRegionName,omitempty"`
	Version          string `json:"Version,omitempty"`
}

// implements the service definition of CheckSignatureRequest
type CheckSignatureRequest struct {
	ApiKeyId      string `json:"ApiKeyId,omitempty"`
	RegionName    string `json:"RegionName,omitempty"`
	RequestDate   string `json:"RequestDate,omitempty"`
	Service       string `json:"Service,omitempty"`
	Signature     string `json:"Signature,omitempty"`
	SignedContent string `json:"SignedContent,omitempty"`
}

// implements the service definition of CheckSignatureResponse
type CheckSignatureResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of ClientEndpoint
type ClientEndpoint struct {
	BgpAsn   int64  `json:"BgpAsn,omitempty"`
	ClientId string `json:"ClientId,omitempty"`
	PublicIp string `json:"PublicIp,omitempty"`
	State    string `json:"State,omitempty"`
	Tags     []Tags `json:"Tags,omitempty"`
	Type     string `json:"Type,omitempty"`
}

// implements the service definition of ClientEndpoints
type ClientEndpoints struct {
	BgpAsn           string `json:"BgpAsn,omitempty"`
	ClientEndpointId string `json:"ClientEndpointId,omitempty"`
	PublicIp         string `json:"PublicIp,omitempty"`
	State            string `json:"State,omitempty"`
	Tags             []Tags `json:"Tags,omitempty"`
	Type             string `json:"Type,omitempty"`
}

// implements the service definition of ConsumptionEntries
type ConsumptionEntries struct {
	Category         string `json:"Category,omitempty"`
	ConsumptionValue string `json:"ConsumptionValue,omitempty"`
	Entry            string `json:"Entry,omitempty"`
	Service          string `json:"Service,omitempty"`
	ShortDescription string `json:"ShortDescription,omitempty"`
	Type             string `json:"Type,omitempty"`
}

// implements the service definition of CopyAccountRequest
type CopyAccountRequest struct {
	DestinationRegionName string `json:"DestinationRegionName,omitempty"`
	QuotaProfile          string `json:"QuotaProfile,omitempty"`
}

// implements the service definition of CopyAccountResponse
type CopyAccountResponse struct {
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
	ImageId         string          `json:"ImageId,omitempty"`
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of CopySnapshotRequest
type CopySnapshotRequest struct {
	Description           string `json:"Description,omitempty"`
	DestinationRegionName string `json:"DestinationRegionName,omitempty"`
	DryRun                bool   `json:"DryRun,omitempty"`
	SourceRegionName      string `json:"SourceRegionName,omitempty"`
	SourceSnapshotId      string `json:"SourceSnapshotId,omitempty"`
}

// implements the service definition of CopySnapshotResponse
type CopySnapshotResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
	SnapshotId      string          `json:"SnapshotId,omitempty"`
}

// implements the service definition of CreateAccountRequest
type CreateAccountRequest struct {
	AccountId     string    `json:"AccountId,omitempty"`
	ApiKeys       []ApiKeys `json:"ApiKeys,omitempty"`
	City          string    `json:"City,omitempty"`
	CompanyName   string    `json:"CompanyName,omitempty"`
	Country       string    `json:"Country,omitempty"`
	CustomerId    string    `json:"CustomerId,omitempty"`
	Email         string    `json:"Email,omitempty"`
	FirstName     string    `json:"FirstName,omitempty"`
	JobTitle      string    `json:"JobTitle,omitempty"`
	LastName      string    `json:"LastName,omitempty"`
	Mobile        string    `json:"Mobile,omitempty"`
	Password      string    `json:"Password,omitempty"`
	Phone         string    `json:"Phone,omitempty"`
	QuotaProfile  string    `json:"QuotaProfile,omitempty"`
	StateProvince string    `json:"StateProvince,omitempty"`
	VatNumber     string    `json:"VatNumber,omitempty"`
	ZipCode       string    `json:"ZipCode,omitempty"`
}

// implements the service definition of CreateAccountResponse
type CreateAccountResponse struct {
	Account         Account         `json:"Account,omitempty"`
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of CreateApiKeyRequest
type CreateApiKeyRequest struct {
	ApiKeyId  string `json:"ApiKeyId,omitempty"`
	SecretKey string `json:"SecretKey,omitempty"`
	Tags      []Tags `json:"Tags,omitempty"`
	UserName  string `json:"UserName,omitempty"`
}

// implements the service definition of CreateApiKeyResponse
type CreateApiKeyResponse struct {
	ApiKey          ApiKey          `json:"ApiKey,omitempty"`
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of CreateClientEndpointRequest
type CreateClientEndpointRequest struct {
	BgpAsn   int64  `json:"BgpAsn,omitempty"`
	DryRun   bool   `json:"DryRun,omitempty"`
	PublicIp string `json:"PublicIp,omitempty"`
	Type     string `json:"Type,omitempty"`
}

// implements the service definition of CreateClientEndpointResponse
type CreateClientEndpointResponse struct {
	ClientEndpoint  ClientEndpoint  `json:"ClientEndpoint,omitempty"`
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of CreateDhcpOptionsRequest
type CreateDhcpOptionsRequest struct {
	DhcpConfigurations []DhcpConfigurations `json:"DhcpConfigurations,omitempty"`
	DryRun             bool                 `json:"DryRun,omitempty"`
}

// implements the service definition of CreateDhcpOptionsResponse
type CreateDhcpOptionsResponse struct {
	DhcpOptionsSet  DhcpOptionsSet  `json:"DhcpOptionsSet,omitempty"`
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of CreateDirectLinkInterfaceRequest
type CreateDirectLinkInterfaceRequest struct {
	DirectLinkId        string              `json:"DirectLinkId,omitempty"`
	DirectLinkInterface DirectLinkInterface `json:"DirectLinkInterface,omitempty"`
}

// implements the service definition of CreateDirectLinkInterfaceResponse
type CreateDirectLinkInterfaceResponse struct {
	AccountId               string          `json:"AccountId,omitempty"`
	BgpAsn                  int64           `json:"BgpAsn,omitempty"`
	BgpKey                  string          `json:"BgpKey,omitempty"`
	ClientPrivateIp         string          `json:"ClientPrivateIp,omitempty"`
	DirectLinkId            string          `json:"DirectLinkId,omitempty"`
	DirectLinkInterfaceId   string          `json:"DirectLinkInterfaceId,omitempty"`
	DirectLinkInterfaceName string          `json:"DirectLinkInterfaceName,omitempty"`
	OutscalePrivateIp       string          `json:"OutscalePrivateIp,omitempty"`
	ResponseContext         ResponseContext `json:"ResponseContext,omitempty"`
	Site                    string          `json:"Site,omitempty"`
	State                   string          `json:"State,omitempty"`
	Type                    string          `json:"Type,omitempty"`
	Vlan                    int64           `json:"Vlan,omitempty"`
	VpnGatewayId            string          `json:"VpnGatewayId,omitempty"`
}

// implements the service definition of CreateDirectLinkRequest
type CreateDirectLinkRequest struct {
	Bandwidth      string `json:"Bandwidth,omitempty"`
	DirectLinkName string `json:"DirectLinkName,omitempty"`
	Site           string `json:"Site,omitempty"`
}

// implements the service definition of CreateDirectLinkResponse
type CreateDirectLinkResponse struct {
	AccountId       string          `json:"AccountId,omitempty"`
	Bandwidth       string          `json:"Bandwidth,omitempty"`
	DirectLinkId    string          `json:"DirectLinkId,omitempty"`
	DirectLinkName  string          `json:"DirectLinkName,omitempty"`
	RegionName      string          `json:"RegionName,omitempty"`
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
	Site            string          `json:"Site,omitempty"`
	State           string          `json:"State,omitempty"`
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

// implements the service definition of CreateGroupRequest
type CreateGroupRequest struct {
	GroupName string `json:"GroupName,omitempty"`
	Path      string `json:"Path,omitempty"`
}

// implements the service definition of CreateGroupResponse
type CreateGroupResponse struct {
	Group           Group           `json:"Group,omitempty"`
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of CreateImageExportTaskRequest
type CreateImageExportTaskRequest struct {
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
	ImageId         string          `json:"ImageId,omitempty"`
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of CreateKeypairRequest
type CreateKeypairRequest struct {
	DryRun      bool   `json:"DryRun,omitempty"`
	KeypairName string `json:"KeypairName,omitempty"`
}

// implements the service definition of CreateKeypairResponse
type CreateKeypairResponse struct {
	KeypairFingerprint string          `json:"KeypairFingerprint,omitempty"`
	KeypairName        string          `json:"KeypairName,omitempty"`
	PrivateKey         string          `json:"PrivateKey,omitempty"`
	ResponseContext    ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of CreateListenerRuleRequest
type CreateListenerRuleRequest struct {
	Listener     Listener     `json:"Listener,omitempty"`
	ListenerRule ListenerRule `json:"ListenerRule,omitempty"`
	VmIds        []string     `json:"VmIds,omitempty"`
}

// implements the service definition of CreateListenerRuleResponse
type CreateListenerRuleResponse struct {
	ListenerId      string          `json:"ListenerId,omitempty"`
	ListenerRule    ListenerRule    `json:"ListenerRule,omitempty"`
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
	VmIds           []string        `json:"VmIds,omitempty"`
}

// implements the service definition of CreateLoadBalancerListenersRequest
type CreateLoadBalancerListenersRequest struct {
	Listeners        []Listeners `json:"Listeners,omitempty"`
	LoadBalancerName string      `json:"LoadBalancerName,omitempty"`
}

// implements the service definition of CreateLoadBalancerListenersResponse
type CreateLoadBalancerListenersResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of CreateLoadBalancerRequest
type CreateLoadBalancerRequest struct {
	FirewallRulesSets []string    `json:"FirewallRulesSets,omitempty"`
	Listeners         []Listeners `json:"Listeners,omitempty"`
	LoadBalancerName  string      `json:"LoadBalancerName,omitempty"`
	LoadBalancerType  string      `json:"LoadBalancerType,omitempty"`
	SubRegionNames    []string    `json:"SubRegionNames,omitempty"`
	Subnets           []string    `json:"Subnets,omitempty"`
	Tags              []Tags      `json:"Tags,omitempty"`
}

// implements the service definition of CreateLoadBalancerResponse
type CreateLoadBalancerResponse struct {
	DnsName         string          `json:"DnsName,omitempty"`
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of CreateNatServiceRequest
type CreateNatServiceRequest struct {
	ClientToken string `json:"ClientToken,omitempty"`
	LinkId      string `json:"LinkId,omitempty"`
	SubnetId    string `json:"SubnetId,omitempty"`
}

// implements the service definition of CreateNatServiceResponse
type CreateNatServiceResponse struct {
	ClientToken     string          `json:"ClientToken,omitempty"`
	NatService      NatService      `json:"NatService,omitempty"`
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of CreateNetAccessRequest
type CreateNetAccessRequest struct {
	DryRun         bool     `json:"DryRun,omitempty"`
	NetId          string   `json:"NetId,omitempty"`
	PrefixListName string   `json:"PrefixListName,omitempty"`
	RouteTableIds  []string `json:"RouteTableIds,omitempty"`
}

// implements the service definition of CreateNetAccessResponse
type CreateNetAccessResponse struct {
	NetAccess       NetAccess       `json:"NetAccess,omitempty"`
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
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
	Description         string       `json:"Description,omitempty"`
	DryRun              bool         `json:"DryRun,omitempty"`
	FirewallRulesSetIds []string     `json:"FirewallRulesSetIds,omitempty"`
	PrivateIps          []PrivateIps `json:"PrivateIps,omitempty"`
	SubnetId            string       `json:"SubnetId,omitempty"`
}

// implements the service definition of CreateNicResponse
type CreateNicResponse struct {
	Nic             Nic             `json:"Nic,omitempty"`
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of CreatePolicyRequest
type CreatePolicyRequest struct {
	Description string `json:"Description,omitempty"`
	Document    string `json:"Document,omitempty"`
	Path        string `json:"Path,omitempty"`
	PolicyName  string `json:"PolicyName,omitempty"`
}

// implements the service definition of CreatePolicyResponse
type CreatePolicyResponse struct {
	Policy          Policy          `json:"Policy,omitempty"`
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
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
	AccountId       string          `json:"AccountId,omitempty"`
	Description     string          `json:"Description,omitempty"`
	Progress        string          `json:"Progress,omitempty"`
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
	SnapshotId      string          `json:"SnapshotId,omitempty"`
	State           string          `json:"State,omitempty"`
	Tags            []Tags          `json:"Tags,omitempty"`
	VolumeId        string          `json:"VolumeId,omitempty"`
	VolumeSize      int64           `json:"VolumeSize,omitempty"`
}

// implements the service definition of CreateStickyCookiePolicyRequest
type CreateStickyCookiePolicyRequest struct {
	CookieName       string `json:"CookieName,omitempty"`
	LoadBalancerName string `json:"LoadBalancerName,omitempty"`
	PolicyName       string `json:"PolicyName,omitempty"`
	Type             string `json:"Type,omitempty"`
}

// implements the service definition of CreateStickyCookiePolicyResponse
type CreateStickyCookiePolicyResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
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

// implements the service definition of CreateUserRequest
type CreateUserRequest struct {
	Path     string `json:"Path,omitempty"`
	UserName string `json:"UserName,omitempty"`
}

// implements the service definition of CreateUserResponse
type CreateUserResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
	User            User            `json:"User,omitempty"`
}

// implements the service definition of CreateVmsRequest
type CreateVmsRequest struct {
	BlockDeviceMappings         []BlockDeviceMappings `json:"BlockDeviceMappings,omitempty"`
	BsuOptimized                bool                  `json:"BsuOptimized,omitempty"`
	ClientToken                 string                `json:"ClientToken,omitempty"`
	DeletionProtection          bool                  `json:"DeletionProtection,omitempty"`
	DryRun                      bool                  `json:"DryRun,omitempty"`
	FirewallRulesSetIds         []string              `json:"FirewallRulesSetIds,omitempty"`
	FirewallRulesSets           []string              `json:"FirewallRulesSets,omitempty"`
	ImageId                     string                `json:"ImageId,omitempty"`
	KeypairName                 string                `json:"KeypairName,omitempty"`
	MaxVmsCount                 int64                 `json:"MaxVmsCount,omitempty"`
	MinVmsCount                 int64                 `json:"MinVmsCount,omitempty"`
	Nics                        []Nics                `json:"Nics,omitempty"`
	Placement                   Placement             `json:"Placement,omitempty"`
	PrivateIps                  []string              `json:"PrivateIps,omitempty"`
	SubnetId                    string                `json:"SubnetId,omitempty"`
	Type                        string                `json:"Type,omitempty"`
	UserData                    string                `json:"UserData,omitempty"`
	VmInitiatedShutdownBehavior string                `json:"VmInitiatedShutdownBehavior,omitempty"`
}

// implements the service definition of CreateVmsResponse
type CreateVmsResponse struct {
	AccountId         string              `json:"AccountId,omitempty"`
	FirewallRulesSets []FirewallRulesSets `json:"FirewallRulesSets,omitempty"`
	ReservationId     string              `json:"ReservationId,omitempty"`
	ResponseContext   ResponseContext     `json:"ResponseContext,omitempty"`
	Vms               []Vms               `json:"Vms,omitempty"`
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

// implements the service definition of CreateVpnConnectionRequest
type CreateVpnConnectionRequest struct {
	ClientEndpointId string `json:"ClientEndpointId,omitempty"`
	DryRun           bool   `json:"DryRun,omitempty"`
	StaticRoutesOnly bool   `json:"StaticRoutesOnly,omitempty"`
	Type             string `json:"Type,omitempty"`
	VpnGatewayId     string `json:"VpnGatewayId,omitempty"`
}

// implements the service definition of CreateVpnConnectionResponse
type CreateVpnConnectionResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
	VpnConnection   VpnConnection   `json:"VpnConnection,omitempty"`
}

// implements the service definition of CreateVpnConnectionRouteRequest
type CreateVpnConnectionRouteRequest struct {
	DestinationIpRange string `json:"DestinationIpRange,omitempty"`
	VpnConnectionId    string `json:"VpnConnectionId,omitempty"`
}

// implements the service definition of CreateVpnConnectionRouteResponse
type CreateVpnConnectionRouteResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
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

// implements the service definition of DeleteApiKeyRequest
type DeleteApiKeyRequest struct {
	ApiKeyId string `json:"ApiKeyId,omitempty"`
}

// implements the service definition of DeleteApiKeyResponse
type DeleteApiKeyResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of DeleteClientEndpointRequest
type DeleteClientEndpointRequest struct {
	ClientEndpointId string `json:"ClientEndpointId,omitempty"`
	DryRun           bool   `json:"DryRun,omitempty"`
}

// implements the service definition of DeleteClientEndpointResponse
type DeleteClientEndpointResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of DeleteDhcpOptionsRequest
type DeleteDhcpOptionsRequest struct {
	DhcpOptionsSetId string `json:"DhcpOptionsSetId,omitempty"`
	DryRun           bool   `json:"DryRun,omitempty"`
}

// implements the service definition of DeleteDhcpOptionsResponse
type DeleteDhcpOptionsResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of DeleteDirectLinkInterfaceRequest
type DeleteDirectLinkInterfaceRequest struct {
	DirectLinkInterfaceId string `json:"DirectLinkInterfaceId,omitempty"`
}

// implements the service definition of DeleteDirectLinkInterfaceResponse
type DeleteDirectLinkInterfaceResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of DeleteDirectLinkRequest
type DeleteDirectLinkRequest struct {
	DirectLinkId string `json:"DirectLinkId,omitempty"`
}

// implements the service definition of DeleteDirectLinkResponse
type DeleteDirectLinkResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
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

// implements the service definition of DeleteGroupRequest
type DeleteGroupRequest struct {
	GroupName string `json:"GroupName,omitempty"`
}

// implements the service definition of DeleteGroupResponse
type DeleteGroupResponse struct {
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

// implements the service definition of DeleteListenerRuleRequest
type DeleteListenerRuleRequest struct {
	ListenerRuleName string `json:"ListenerRuleName,omitempty"`
}

// implements the service definition of DeleteListenerRuleResponse
type DeleteListenerRuleResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of DeleteLoadBalancerListenersRequest
type DeleteLoadBalancerListenersRequest struct {
	LoadBalancerName  string `json:"LoadBalancerName,omitempty"`
	LoadBalancerPorts []int  `json:"LoadBalancerPorts,omitempty"`
}

// implements the service definition of DeleteLoadBalancerListenersResponse
type DeleteLoadBalancerListenersResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of DeleteLoadBalancerPolicyRequest
type DeleteLoadBalancerPolicyRequest struct {
	LoadBalancerName string `json:"LoadBalancerName,omitempty"`
	PolicyName       string `json:"PolicyName,omitempty"`
}

// implements the service definition of DeleteLoadBalancerPolicyResponse
type DeleteLoadBalancerPolicyResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of DeleteLoadBalancerRequest
type DeleteLoadBalancerRequest struct {
	LoadBalancerName string `json:"LoadBalancerName,omitempty"`
}

// implements the service definition of DeleteLoadBalancerResponse
type DeleteLoadBalancerResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of DeleteNatServiceRequest
type DeleteNatServiceRequest struct {
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

// implements the service definition of DeletePolicyRequest
type DeletePolicyRequest struct {
	PolicyId string `json:"PolicyId,omitempty"`
}

// implements the service definition of DeletePolicyResponse
type DeletePolicyResponse struct {
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

// implements the service definition of DeleteServerCertificateRequest
type DeleteServerCertificateRequest struct {
	ServerCertificateName string `json:"ServerCertificateName,omitempty"`
}

// implements the service definition of DeleteServerCertificateResponse
type DeleteServerCertificateResponse struct {
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

// implements the service definition of DeleteUserRequest
type DeleteUserRequest struct {
	UserName string `json:"UserName,omitempty"`
}

// implements the service definition of DeleteUserResponse
type DeleteUserResponse struct {
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
	Vms             []Vms           `json:"Vms,omitempty"`
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

// implements the service definition of DeleteVpcEndpointsRequest
type DeleteVpcEndpointsRequest struct {
	DryRun       bool     `json:"DryRun,omitempty"`
	NetAccessIds []string `json:"NetAccessIds,omitempty"`
}

// implements the service definition of DeleteVpcEndpointsResponse
type DeleteVpcEndpointsResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of DeleteVpnConnectionRequest
type DeleteVpnConnectionRequest struct {
	DryRun          bool   `json:"DryRun,omitempty"`
	VpnConnectionId string `json:"VpnConnectionId,omitempty"`
}

// implements the service definition of DeleteVpnConnectionResponse
type DeleteVpnConnectionResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of DeleteVpnConnectionRouteRequest
type DeleteVpnConnectionRouteRequest struct {
	DestinationIpRange string `json:"DestinationIpRange,omitempty"`
	DryRun             bool   `json:"DryRun,omitempty"`
	VpnConnectionId    string `json:"VpnConnectionId,omitempty"`
}

// implements the service definition of DeleteVpnConnectionRouteResponse
type DeleteVpnConnectionRouteResponse struct {
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

// implements the service definition of DeregisterUserInGroupRequest
type DeregisterUserInGroupRequest struct {
	GroupName string `json:"GroupName,omitempty"`
	UserName  string `json:"UserName,omitempty"`
}

// implements the service definition of DeregisterUserInGroupResponse
type DeregisterUserInGroupResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of DeregisterVmsInListenerRuleRequest
type DeregisterVmsInListenerRuleRequest struct {
	ListenerRuleName string   `json:"ListenerRuleName,omitempty"`
	VmIds            []string `json:"VmIds,omitempty"`
}

// implements the service definition of DeregisterVmsInListenerRuleResponse
type DeregisterVmsInListenerRuleResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
	VmIds           []string        `json:"VmIds,omitempty"`
}

// implements the service definition of DeregisterVmsInLoadBalancerRequest
type DeregisterVmsInLoadBalancerRequest struct {
	BackendVmsIds    []string `json:"BackendVmsIds,omitempty"`
	LoadBalancerName string   `json:"LoadBalancerName,omitempty"`
}

// implements the service definition of DeregisterVmsInLoadBalancerResponse
type DeregisterVmsInLoadBalancerResponse struct {
	BackendVmsIds   []string        `json:"BackendVmsIds,omitempty"`
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of DhcpConfigurations
type DhcpConfigurations struct {
	Key    string   `json:"Key,omitempty"`
	Values []string `json:"Values,omitempty"`
}

// implements the service definition of DhcpOptionsSet
type DhcpOptionsSet struct {
	DhcpConfigurations []DhcpConfigurations `json:"DhcpConfigurations,omitempty"`
	DhcpOptionsSetId   string               `json:"DhcpOptionsSetId,omitempty"`
	Tags               []Tags               `json:"Tags,omitempty"`
}

// implements the service definition of DhcpOptionsSets
type DhcpOptionsSets struct {
	DhcpConfigurations []DhcpConfigurations `json:"DhcpConfigurations,omitempty"`
	DhcpOptionsSetId   string               `json:"DhcpOptionsSetId,omitempty"`
	Tags               []Tags               `json:"Tags,omitempty"`
}

// implements the service definition of DirectLinkInterface
type DirectLinkInterface struct {
	BgpAsn                  int64  `json:"BgpAsn,omitempty"`
	BgpKey                  string `json:"BgpKey,omitempty"`
	ClientPrivateIp         string `json:"ClientPrivateIp,omitempty"`
	DirectLinkInterfaceName string `json:"DirectLinkInterfaceName,omitempty"`
	OutscalePrivateIp       string `json:"OutscalePrivateIp,omitempty"`
	Vlan                    int64  `json:"Vlan,omitempty"`
	VpnGatewayId            string `json:"VpnGatewayId,omitempty"`
}

// implements the service definition of DirectLinkInterfaces
type DirectLinkInterfaces struct {
	AccountId               string `json:"AccountId,omitempty"`
	BgpAsn                  int64  `json:"BgpAsn,omitempty"`
	BgpKey                  string `json:"BgpKey,omitempty"`
	ClientPrivateIp         string `json:"ClientPrivateIp,omitempty"`
	DirectLinkId            string `json:"DirectLinkId,omitempty"`
	DirectLinkInterfaceId   string `json:"DirectLinkInterfaceId,omitempty"`
	DirectLinkInterfaceName string `json:"DirectLinkInterfaceName,omitempty"`
	OutscalePrivateIp       string `json:"OutscalePrivateIp,omitempty"`
	Site                    string `json:"Site,omitempty"`
	State                   string `json:"State,omitempty"`
	Type                    string `json:"Type,omitempty"`
	Vlan                    int64  `json:"Vlan,omitempty"`
	VpnGatewayId            string `json:"VpnGatewayId,omitempty"`
}

// implements the service definition of DirectLinks
type DirectLinks struct {
	AccountId      string `json:"AccountId,omitempty"`
	Bandwidth      string `json:"Bandwidth,omitempty"`
	DirectLinkId   string `json:"DirectLinkId,omitempty"`
	DirectLinkName string `json:"DirectLinkName,omitempty"`
	RegionName     string `json:"RegionName,omitempty"`
	Site           string `json:"Site,omitempty"`
	State          string `json:"State,omitempty"`
}

// implements the service definition of Filters
type Filters struct {
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

// implements the service definition of FirewallRulesSetLogging
type FirewallRulesSetLogging struct {
	IsEnabled    bool   `json:"IsEnabled,omitempty"`
	RateLimit    string `json:"RateLimit,omitempty"`
	SyslogServer string `json:"SyslogServer,omitempty"`
}

// implements the service definition of FirewallRulesSets
type FirewallRulesSets struct {
	FirewallRulesSetId   string `json:"FirewallRulesSetId,omitempty"`
	FirewallRulesSetName string `json:"FirewallRulesSetName,omitempty"`
}

// implements the service definition of FirewallRulesSetsMembers
type FirewallRulesSetsMembers struct {
	AccountId          string `json:"AccountId,omitempty"`
	FirewallRulesSetId string `json:"FirewallRulesSetId,omitempty"`
	Name               string `json:"Name,omitempty"`
}

// implements the service definition of GetBillableDigestRequest
type GetBillableDigestRequest struct {
	AccountId      string `json:"AccountId,omitempty"`
	FromDate       string `json:"FromDate,omitempty"`
	InvoiceState   string `json:"InvoiceState,omitempty"`
	IsConsolidated bool   `json:"IsConsolidated,omitempty"`
	ToDate         string `json:"ToDate,omitempty"`
}

// implements the service definition of GetBillableDigestResponse
type GetBillableDigestResponse struct {
	Items           []Items         `json:"Items,omitempty"`
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of GetRegionConfigRequest
type GetRegionConfigRequest struct {
	FromDate string `json:"FromDate,omitempty"`
}

// implements the service definition of GetRegionConfigResponse
type GetRegionConfigResponse struct {
	RegionConfig    RegionConfig    `json:"RegionConfig,omitempty"`
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of Group
type Group struct {
	GroupId   string `json:"GroupId,omitempty"`
	GroupName string `json:"GroupName,omitempty"`
	Path      string `json:"Path,omitempty"`
}

// implements the service definition of Groups
type Groups struct {
	GroupId   string `json:"GroupId,omitempty"`
	GroupName string `json:"GroupName,omitempty"`
	Path      string `json:"Path,omitempty"`
}

// implements the service definition of HealthCheck
type HealthCheck struct {
	CheckInterval      int64  `json:"CheckInterval,omitempty"`
	CheckedVm          string `json:"CheckedVm,omitempty"`
	HealthyThreshold   int64  `json:"HealthyThreshold,omitempty"`
	Timeout            int64  `json:"Timeout,omitempty"`
	UnhealthyThreshold int64  `json:"UnhealthyThreshold,omitempty"`
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
	AccountAlias        string                `json:"AccountAlias,omitempty"`
	AccountId           string                `json:"AccountId,omitempty"`
	Architecture        string                `json:"Architecture,omitempty"`
	BlockDeviceMappings []BlockDeviceMappings `json:"BlockDeviceMappings,omitempty"`
	CreationDate        string                `json:"CreationDate,omitempty"`
	Description         string                `json:"Description,omitempty"`
	ImageId             string                `json:"ImageId,omitempty"`
	IsPublic            bool                  `json:"IsPublic,omitempty"`
	Name                string                `json:"Name,omitempty"`
	OsuLocation         string                `json:"OsuLocation,omitempty"`
	ProductCodes        []ProductCodes        `json:"ProductCodes,omitempty"`
	RootDeviceName      string                `json:"RootDeviceName,omitempty"`
	RootDeviceType      string                `json:"RootDeviceType,omitempty"`
	State               string                `json:"State,omitempty"`
	StateComment        StateComment          `json:"StateComment,omitempty"`
	Tags                []Tags                `json:"Tags,omitempty"`
	Type                string                `json:"Type,omitempty"`
}

// implements the service definition of ImportKeyPairRequest
type ImportKeyPairRequest struct {
	DryRun      bool   `json:"DryRun,omitempty"`
	KeypairName string `json:"KeypairName,omitempty"`
	PublicKey   string `json:"PublicKey,omitempty"`
}

// implements the service definition of ImportKeyPairResponse
type ImportKeyPairResponse struct {
	KeypairFingerprint string          `json:"KeypairFingerprint,omitempty"`
	ResponseContext    ResponseContext `json:"ResponseContext,omitempty"`
	KeypairName        string          `json:"keypairName,omitempty"`
}

// implements the service definition of ImportServerCertificateRequest
type ImportServerCertificateRequest struct {
	PrivateKey             string `json:"PrivateKey,omitempty"`
	ServerCertificateBody  string `json:"ServerCertificateBody,omitempty"`
	ServerCertificateChain string `json:"ServerCertificateChain,omitempty"`
	ServerCertificateName  string `json:"ServerCertificateName,omitempty"`
	ServerCertificatePath  string `json:"ServerCertificatePath,omitempty"`
}

// implements the service definition of ImportServerCertificateResponse
type ImportServerCertificateResponse struct {
	ResponseContext   ResponseContext   `json:"ResponseContext,omitempty"`
	ServerCertificate ServerCertificate `json:"ServerCertificate,omitempty"`
}

// implements the service definition of ImportSnaptShotRequest
type ImportSnaptShotRequest struct {
	Description  string `json:"Description,omitempty"`
	OsuLocation  string `json:"OsuLocation,omitempty"`
	SnapshotSize int64  `json:"SnapshotSize,omitempty"`
}

// implements the service definition of ImportSnaptShotResponse
type ImportSnaptShotResponse struct {
	AccountAlias    string          `json:"AccountAlias,omitempty"`
	Description     string          `json:"Description,omitempty"`
	IsEncrypted     bool            `json:"IsEncrypted,omitempty"`
	Progress        int64           `json:"Progress,omitempty"`
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
	SnapshotId      string          `json:"SnapshotId,omitempty"`
	SnapshotState   string          `json:"SnapshotState,omitempty"`
	StartDate       string          `json:"StartDate,omitempty"`
	VolumeSize      int64           `json:"VolumeSize,omitempty"`
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

// implements the service definition of Items
type Items struct {
	AccountId       string    `json:"AccountId,omitempty"`
	Catalog         []Catalog `json:"Catalog,omitempty"`
	ComsuptionValue int       `json:"ComsuptionValue,omitempty"`
	Entry           string    `json:"Entry,omitempty"`
	FromDate        string    `json:"FromDate,omitempty"`
	PayingAccountId string    `json:"PayingAccountId,omitempty"`
	Service         string    `json:"Service,omitempty"`
	SubRegionName   string    `json:"SubRegionName,omitempty"`
	ToDate          string    `json:"ToDate,omitempty"`
	Type            string    `json:"Type,omitempty"`
}

// implements the service definition of Keypairs
type Keypairs struct {
	KeypairFingerprint string `json:"KeypairFingerprint,omitempty"`
	KeypairName        string `json:"KeypairName,omitempty"`
}

// implements the service definition of LinkDhcpOptionsRequest
type LinkDhcpOptionsRequest struct {
	DhcpOptionsSetId string `json:"DhcpOptionsSetId,omitempty"`
	DryRun           bool   `json:"DryRun,omitempty"`
	NetId            string `json:"NetId,omitempty"`
}

// implements the service definition of LinkDhcpOptionsResponse
type LinkDhcpOptionsResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of LinkLoadBalancerServerCertificateRequest
type LinkLoadBalancerServerCertificateRequest struct {
	LoadBalancerName    string `json:"LoadBalancerName,omitempty"`
	LoadBalancerPort    int64  `json:"LoadBalancerPort,omitempty"`
	ServerCertificateId string `json:"ServerCertificateId,omitempty"`
}

// implements the service definition of LinkLoadBalancerServerCertificateResponse
type LinkLoadBalancerServerCertificateResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
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

// implements the service definition of LinkPolicyRequest
type LinkPolicyRequest struct {
	GroupName string `json:"GroupName,omitempty"`
	PolicyId  string `json:"PolicyId,omitempty"`
	UserName  string `json:"UserName,omitempty"`
}

// implements the service definition of LinkPolicyResponse
type LinkPolicyResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of LinkPrivateIpRequest
type LinkPrivateIpRequest struct {
	AllowRelink             bool     `json:"AllowRelink,omitempty"`
	NicId                   string   `json:"NicId,omitempty"`
	PrivateIps              []string `json:"PrivateIps,omitempty"`
	SecondaryPrivateIpCount int64    `json:"SecondaryPrivateIpCount,omitempty"`
}

// implements the service definition of LinkPrivateIpResponse
type LinkPrivateIpResponse struct {
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

// implements the service definition of ListGroupsForUserRequest
type ListGroupsForUserRequest struct {
	MaxResults       int64  `json:"MaxResults,omitempty"`
	NextResultsToken string `json:"NextResultsToken,omitempty"`
	UserName         string `json:"UserName,omitempty"`
}

// implements the service definition of ListGroupsForUserResponse
type ListGroupsForUserResponse struct {
	Groups           []Groups        `json:"Groups,omitempty"`
	NextResultsToken string          `json:"NextResultsToken,omitempty"`
	ResponseContext  ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of Listener
type Listener struct {
	LoadBalancerName string `json:"LoadBalancerName,omitempty"`
	LoadBalancerPort int64  `json:"LoadBalancerPort,omitempty"`
}

// implements the service definition of ListenerRule
type ListenerRule struct {
	Action           string `json:"Action,omitempty"`
	HostNamePattern  string `json:"HostNamePattern,omitempty"`
	ListenerRuleId   string `json:"ListenerRuleId,omitempty"`
	ListenerRuleName string `json:"ListenerRuleName,omitempty"`
	PathPattern      string `json:"PathPattern,omitempty"`
	Priority         int64  `json:"Priority,omitempty"`
}

// implements the service definition of ListenerRules
type ListenerRules struct {
	ListenerId   string       `json:"ListenerId,omitempty"`
	ListenerRule ListenerRule `json:"ListenerRule,omitempty"`
	VmIds        []string     `json:"VmIds,omitempty"`
}

// implements the service definition of Listeners
type Listeners struct {
	BackendPort          int64  `json:"BackendPort,omitempty"`
	BackendProtocol      string `json:"BackendProtocol,omitempty"`
	LoadBalancerPort     int64  `json:"LoadBalancerPort,omitempty"`
	LoadBalancerProtocol string `json:"LoadBalancerProtocol,omitempty"`
	ServerCertificateId  string `json:"ServerCertificateId,omitempty"`
}

// implements the service definition of LoadBalancerStickyCookiePolicies
type LoadBalancerStickyCookiePolicies struct {
	PolicyName string `json:"PolicyName,omitempty"`
}

// implements the service definition of LoadBalancers
type LoadBalancers struct {
	ApplicationStickyCookiePolicies  []ApplicationStickyCookiePolicies  `json:"ApplicationStickyCookiePolicies,omitempty"`
	BackendVmsIds                    []string                           `json:"BackendVmsIds,omitempty"`
	DnsName                          string                             `json:"DnsName,omitempty"`
	FirewallRulesSets                []string                           `json:"FirewallRulesSets,omitempty"`
	HealthCheck                      HealthCheck                        `json:"HealthCheck,omitempty"`
	Listeners                        []Listeners                        `json:"Listeners,omitempty"`
	LoadBalancerName                 string                             `json:"LoadBalancerName,omitempty"`
	LoadBalancerStickyCookiePolicies []LoadBalancerStickyCookiePolicies `json:"LoadBalancerStickyCookiePolicies,omitempty"`
	LoadBalancerType                 string                             `json:"LoadBalancerType,omitempty"`
	NetId                            string                             `json:"NetId,omitempty"`
	SourceFirewallRulesSet           SourceFirewallRulesSet             `json:"SourceFirewallRulesSet,omitempty"`
	SubRegionNames                   []string                           `json:"SubRegionNames,omitempty"`
	Subnets                          []string                           `json:"Subnets,omitempty"`
}

// implements the service definition of Logs
type Logs struct {
	CallDuration       int64  `json:"CallDuration,omitempty"`
	QueryAccessKey     string `json:"QueryAccessKey,omitempty"`
	QueryApiName       string `json:"QueryApiName,omitempty"`
	QueryApiVersion    string `json:"QueryApiVersion,omitempty"`
	QueryCallName      string `json:"QueryCallName,omitempty"`
	QueryDate          string `json:"QueryDate,omitempty"`
	QueryIpAddress     string `json:"QueryIpAddress,omitempty"`
	QueryRaw           string `json:"QueryRaw,omitempty"`
	QuerySize          int64  `json:"QuerySize,omitempty"`
	QueryUserAgent     string `json:"QueryUserAgent,omitempty"`
	ResponseId         string `json:"ResponseId,omitempty"`
	ResponseSize       int64  `json:"ResponseSize,omitempty"`
	ResponseStatusCode int64  `json:"ResponseStatusCode,omitempty"`
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
	NatServiceId string      `json:"NatServiceId,omitempty"`
	NetId        string      `json:"NetId,omitempty"`
	PublicIps    []PublicIps `json:"PublicIps,omitempty"`
	State        string      `json:"State,omitempty"`
	SubnetId     string      `json:"SubnetId,omitempty"`
}

// implements the service definition of NatServices
type NatServices struct {
	NatServiceId string      `json:"NatServiceId,omitempty"`
	NetId        string      `json:"NetId,omitempty"`
	PublicIps    []PublicIps `json:"PublicIps,omitempty"`
	State        string      `json:"State,omitempty"`
	SubnetId     string      `json:"SubnetId,omitempty"`
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

// implements the service definition of NetAccess
type NetAccess struct {
	NetAccessId    string   `json:"NetAccessId,omitempty"`
	NetId          string   `json:"NetId,omitempty"`
	PrefixListName string   `json:"PrefixListName,omitempty"`
	RouteTableIds  []string `json:"RouteTableIds,omitempty"`
	State          string   `json:"State,omitempty"`
}

// implements the service definition of NetAccesses
type NetAccesses struct {
	NetAccessId    string   `json:"NetAccessId,omitempty"`
	NetId          string   `json:"NetId,omitempty"`
	PrefixListName string   `json:"PrefixListName,omitempty"`
	RouteTableIds  []string `json:"RouteTableIds,omitempty"`
	State          string   `json:"State,omitempty"`
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
	AccountId           string              `json:"AccountId,omitempty"`
	Description         string              `json:"Description,omitempty"`
	FirewallRulesSets   []FirewallRulesSets `json:"FirewallRulesSets,omitempty"`
	IsSourceDestChecked bool                `json:"IsSourceDestChecked,omitempty"`
	MacAddress          string              `json:"MacAddress,omitempty"`
	NetId               string              `json:"NetId,omitempty"`
	NicId               string              `json:"NicId,omitempty"`
	NicLink             NicLink             `json:"NicLink,omitempty"`
	PrivateDnsName      string              `json:"PrivateDnsName,omitempty"`
	PrivateIps          []PrivateIps        `json:"PrivateIps,omitempty"`
	PublicIpToNicLink   PublicIpToNicLink   `json:"PublicIpToNicLink,omitempty"`
	State               string              `json:"State,omitempty"`
	SubnetId            string              `json:"SubnetId,omitempty"`
	SubregionName       string              `json:"SubregionName,omitempty"`
	Tags                []Tags              `json:"Tags,omitempty"`
}

// implements the service definition of NicLink
type NicLink struct {
	DeleteOnVmDeletion bool   `json:"DeleteOnVmDeletion,omitempty"`
	NicLinkId          string `json:"NicLinkId,omitempty"`
}

// implements the service definition of Nics
type Nics struct {
	AccountId           string              `json:"AccountId,omitempty"`
	Description         string              `json:"Description,omitempty"`
	FirewallRulesSets   []FirewallRulesSets `json:"FirewallRulesSets,omitempty"`
	IsSourceDestChecked bool                `json:"IsSourceDestChecked,omitempty"`
	MacAddress          string              `json:"MacAddress,omitempty"`
	NetId               string              `json:"NetId,omitempty"`
	NicId               string              `json:"NicId,omitempty"`
	NicLink             NicLink             `json:"NicLink,omitempty"`
	PrivateDnsName      string              `json:"PrivateDnsName,omitempty"`
	PrivateIps          []PrivateIps        `json:"PrivateIps,omitempty"`
	PublicIpToNicLink   PublicIpToNicLink   `json:"PublicIpToNicLink,omitempty"`
	State               string              `json:"State,omitempty"`
	SubnetId            string              `json:"SubnetId,omitempty"`
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

// implements the service definition of Permission
type Permission struct {
	Additions []Additions `json:"Additions,omitempty"`
	Removals  []Removals  `json:"Removals,omitempty"`
}

// implements the service definition of Permissions
type Permissions struct {
	AccountId        string `json:"AccountId,omitempty"`
	GlobalPermission string `json:"GlobalPermission,omitempty"`
}

// implements the service definition of PermissionsToCreateVolumes
type PermissionsToCreateVolumes struct {
	AccountId        string `json:"AccountId,omitempty"`
	GlobalPermission string `json:"GlobalPermission,omitempty"`
}

// implements the service definition of Placement
type Placement struct {
	Affinity        string `json:"Affinity,omitempty"`
	DedicatedHostId string `json:"DedicatedHostId,omitempty"`
	PlacementName   string `json:"PlacementName,omitempty"`
	SubRegionName   string `json:"SubRegionName,omitempty"`
	Tenancy         string `json:"Tenancy,omitempty"`
}

// implements the service definition of Policies
type Policies struct {
	Description            string `json:"Description,omitempty"`
	IsLinkable             bool   `json:"IsLinkable,omitempty"`
	Path                   string `json:"Path,omitempty"`
	PolicyDefaultVersionId string `json:"PolicyDefaultVersionId,omitempty"`
	PolicyId               string `json:"PolicyId,omitempty"`
	PolicyName             string `json:"PolicyName,omitempty"`
	ResourcesCount         int64  `json:"ResourcesCount,omitempty"`
}

// implements the service definition of Policy
type Policy struct {
	Description            string `json:"Description,omitempty"`
	IsLinkable             bool   `json:"IsLinkable,omitempty"`
	Path                   string `json:"Path,omitempty"`
	PolicyDefaultVersionId string `json:"PolicyDefaultVersionId,omitempty"`
	PolicyId               string `json:"PolicyId,omitempty"`
	PolicyName             string `json:"PolicyName,omitempty"`
	ResourcesCount         int64  `json:"ResourcesCount,omitempty"`
}

// implements the service definition of PrefixLists
type PrefixLists struct {
	IpRanges       []string `json:"IpRanges,omitempty"`
	PrefixListId   string   `json:"PrefixListId,omitempty"`
	PrefixListName string   `json:"PrefixListName,omitempty"`
}

// implements the service definition of PricingDetails
type PricingDetails struct {
	Count int64 `json:"Count,omitempty"`
}

// implements the service definition of PrivateIps
type PrivateIps struct {
	IsPrimary         bool              `json:"IsPrimary,omitempty"`
	PrivateDnsName    string            `json:"PrivateDnsName,omitempty"`
	PrivateIp         string            `json:"PrivateIp,omitempty"`
	PublicIpToNicLink PublicIpToNicLink `json:"PublicIpToNicLink,omitempty"`
}

// implements the service definition of ProductCodes
type ProductCodes struct {
	ProductCode string `json:"ProductCode,omitempty"`
	ProductType string `json:"ProductType,omitempty"`
}

// implements the service definition of ProductTypes
type ProductTypes struct {
	Description   string `json:"Description,omitempty"`
	ProductTypeId string `json:"ProductTypeId,omitempty"`
	Vendor        string `json:"Vendor,omitempty"`
}

// implements the service definition of PublicIpToNicLink
type PublicIpToNicLink struct {
	PublicDnsName     string `json:"PublicDnsName,omitempty"`
	PublicIp          string `json:"PublicIp,omitempty"`
	PublicIpAccountId string `json:"PublicIpAccountId,omitempty"`
}

// implements the service definition of PublicIps
type PublicIps struct {
	LinkId   string `json:"LinkId,omitempty"`
	PublicIp string `json:"PublicIp,omitempty"`
}

// implements the service definition of PurchaseReservedVmsOfferRequest
type PurchaseReservedVmsOfferRequest struct {
	DryRun             bool   `json:"DryRun,omitempty"`
	ReservedVmsOfferId string `json:"ReservedVmsOfferId,omitempty"`
	VmCount            int64  `json:"VmCount,omitempty"`
}

// implements the service definition of PurchaseReservedVmsOfferResponse
type PurchaseReservedVmsOfferResponse struct {
	ReservedVmsId   string          `json:"ReservedVmsId,omitempty"`
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of QuotaTypes
type QuotaTypes struct {
	QuotaType string   `json:"QuotaType,omitempty"`
	Quotas    []Quotas `json:"Quotas,omitempty"`
}

// implements the service definition of Quotas
type Quotas struct {
	AccountId        string `json:"AccountId,omitempty"`
	Description      string `json:"Description,omitempty"`
	MaxValue         int64  `json:"MaxValue,omitempty"`
	Name             string `json:"Name,omitempty"`
	QuotaCollection  string `json:"QuotaCollection,omitempty"`
	ShortDescription string `json:"ShortDescription,omitempty"`
	UsedValue        int64  `json:"UsedValue,omitempty"`
}

// implements the service definition of ReadAccountConsumptionRequest
type ReadAccountConsumptionRequest struct {
	FromDate string `json:"FromDate,omitempty"`
	ToDate   string `json:"ToDate,omitempty"`
}

// implements the service definition of ReadAccountConsumptionResponse
type ReadAccountConsumptionResponse struct {
	ConsumptionEntries ConsumptionEntries `json:"ConsumptionEntries,omitempty"`
	ResponseContext    ResponseContext    `json:"ResponseContext,omitempty"`
}

// implements the service definition of ReadAccountRequest
type ReadAccountRequest struct {
}

// implements the service definition of ReadAccountResponse
type ReadAccountResponse struct {
	Account         Account         `json:"Account,omitempty"`
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of ReadAdminPasswordRequest
type ReadAdminPasswordRequest struct {
	DryRun bool   `json:"DryRun,omitempty"`
	VmId   string `json:"VmId,omitempty"`
}

// implements the service definition of ReadAdminPasswordResponse
type ReadAdminPasswordResponse struct {
	AdminPassword   string          `json:"AdminPassword,omitempty"`
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
	VmId            string          `json:"VmId,omitempty"`
}

// implements the service definition of ReadApiKeysRequest
type ReadApiKeysRequest struct {
	MaxResults       int64  `json:"MaxResults,omitempty"`
	NextResultsToken string `json:"NextResultsToken,omitempty"`
	Tags             []Tags `json:"Tags,omitempty"`
	UserName         string `json:"UserName,omitempty"`
}

// implements the service definition of ReadApiKeysResponse
type ReadApiKeysResponse struct {
	ApiKeys          []ApiKeys       `json:"ApiKeys,omitempty"`
	NextResultsToken string          `json:"NextResultsToken,omitempty"`
	ResponseContext  ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of ReadApiLogsRequest
type ReadApiLogsRequest struct {
	DryRun           bool    `json:"DryRun,omitempty"`
	Filters          Filters `json:"Filters,omitempty"`
	MaxResults       int64   `json:"MaxResults,omitempty"`
	NextResultsToken string  `json:"NextResultsToken,omitempty"`
	With             With    `json:"With,omitempty"`
}

// implements the service definition of ReadApiLogsResponse
type ReadApiLogsResponse struct {
	Logs             []Logs          `json:"Logs,omitempty"`
	NextResultsToken string          `json:"NextResultsToken,omitempty"`
	ResponseContext  ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of ReadCatalogRequest
type ReadCatalogRequest struct {
}

// implements the service definition of ReadCatalogResponse
type ReadCatalogResponse struct {
	Catalog         Catalog         `json:"Catalog,omitempty"`
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of ReadClientEndpointsRequest
type ReadClientEndpointsRequest struct {
	ClientEndpointIds []string  `json:"ClientEndpointIds,omitempty"`
	DryRun            bool      `json:"DryRun,omitempty"`
	Filters           []Filters `json:"Filters,omitempty"`
}

// implements the service definition of ReadClientEndpointsResponse
type ReadClientEndpointsResponse struct {
	ClientEndpoints []ClientEndpoints `json:"ClientEndpoints,omitempty"`
	ResponseContext ResponseContext   `json:"ResponseContext,omitempty"`
}

// implements the service definition of ReadConsoleOutputRequest
type ReadConsoleOutputRequest struct {
	DryRun bool   `json:"DryRun,omitempty"`
	VmId   string `json:"VmId,omitempty"`
}

// implements the service definition of ReadConsoleOutputResponse
type ReadConsoleOutputResponse struct {
	ConsoleOutput   string          `json:"ConsoleOutput,omitempty"`
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
	VmId            string          `json:"VmId,omitempty"`
}

// implements the service definition of ReadDhcpOptionsRequest
type ReadDhcpOptionsRequest struct {
	DhcpOptionsSetIds []string  `json:"DhcpOptionsSetIds,omitempty"`
	DryRun            bool      `json:"DryRun,omitempty"`
	Filters           []Filters `json:"Filters,omitempty"`
}

// implements the service definition of ReadDhcpOptionsResponse
type ReadDhcpOptionsResponse struct {
	DhcpOptionsSets []DhcpOptionsSets `json:"DhcpOptionsSets,omitempty"`
	ResponseContext ResponseContext   `json:"ResponseContext,omitempty"`
}

// implements the service definition of ReadDirectLinkInterfacesRequest
type ReadDirectLinkInterfacesRequest struct {
	DirectLinkId          string `json:"DirectLinkId,omitempty"`
	DirectLinkInterfaceId string `json:"DirectLinkInterfaceId,omitempty"`
}

// implements the service definition of ReadDirectLinkInterfacesResponse
type ReadDirectLinkInterfacesResponse struct {
	DirectLinkInterfaces []DirectLinkInterfaces `json:"DirectLinkInterfaces,omitempty"`
	ResponseContext      ResponseContext        `json:"ResponseContext,omitempty"`
}

// implements the service definition of ReadDirectLinksRequest
type ReadDirectLinksRequest struct {
	DirectLinkId string `json:"DirectLinkId,omitempty"`
}

// implements the service definition of ReadDirectLinksResponse
type ReadDirectLinksResponse struct {
	DirectLinks     []DirectLinks   `json:"DirectLinks,omitempty"`
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of ReadFirewallRulesSetsRequest
type ReadFirewallRulesSetsRequest struct {
	DryRun              bool      `json:"DryRun,omitempty"`
	Filters             []Filters `json:"Filters,omitempty"`
	FirewallRulesSetIds []string  `json:"FirewallRulesSetIds,omitempty"`
	Names               []string  `json:"Names,omitempty"`
}

// implements the service definition of ReadFirewallRulesSetsResponse
type ReadFirewallRulesSetsResponse struct {
	FirewallRulesSets []FirewallRulesSets `json:"FirewallRulesSets,omitempty"`
	ResponseContext   ResponseContext     `json:"ResponseContext,omitempty"`
}

// implements the service definition of ReadGroupsRequest
type ReadGroupsRequest struct {
	MaxResults       int64  `json:"MaxResults,omitempty"`
	NextResultsToken string `json:"NextResultsToken,omitempty"`
	Path             string `json:"Path,omitempty"`
}

// implements the service definition of ReadGroupsResponse
type ReadGroupsResponse struct {
	Groups           []Groups        `json:"Groups,omitempty"`
	NextResultsToken string          `json:"NextResultsToken,omitempty"`
	ResponseContext  ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of ReadImageAttributeRequest
type ReadImageAttributeRequest struct {
	Attribute string `json:"Attribute,omitempty"`
	DryRun    bool   `json:"DryRun,omitempty"`
	ImageId   string `json:"ImageId,omitempty"`
}

// implements the service definition of ReadImageAttributeResponse
type ReadImageAttributeResponse struct {
	Description     string          `json:"Description,omitempty"`
	ImageId         string          `json:"ImageId,omitempty"`
	Permissions     []Permissions   `json:"Permissions,omitempty"`
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of ReadImageExportTasksRequest
type ReadImageExportTasksRequest struct {
	TaskIds []string `json:"TaskIds,omitempty"`
}

// implements the service definition of ReadImageExportTasksResponse
type ReadImageExportTasksResponse struct {
	ImageExportTasks []ImageExportTasks `json:"ImageExportTasks,omitempty"`
	ResponseContext  ResponseContext    `json:"ResponseContext,omitempty"`
}

// implements the service definition of ReadImagesRequest
type ReadImagesRequest struct {
	AccountIds  []string  `json:"AccountIds,omitempty"`
	DryRun      bool      `json:"DryRun,omitempty"`
	Filters     []Filters `json:"Filters,omitempty"`
	ImageIds    []string  `json:"ImageIds,omitempty"`
	Permissions []string  `json:"Permissions,omitempty"`
}

// implements the service definition of ReadImagesResponse
type ReadImagesResponse struct {
	Images          []Images        `json:"Images,omitempty"`
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of ReadKeypairsRequest
type ReadKeypairsRequest struct {
	DryRun       bool      `json:"DryRun,omitempty"`
	Filters      []Filters `json:"Filters,omitempty"`
	KeypairNames []string  `json:"KeypairNames,omitempty"`
}

// implements the service definition of ReadKeypairsResponse
type ReadKeypairsResponse struct {
	Keypairs        []Keypairs      `json:"Keypairs,omitempty"`
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of ReadListenerRulesRequest
type ReadListenerRulesRequest struct {
	ListenerRuleNames []string `json:"ListenerRuleNames,omitempty"`
}

// implements the service definition of ReadListenerRulesResponse
type ReadListenerRulesResponse struct {
	ListenerRules   []ListenerRules `json:"ListenerRules,omitempty"`
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of ReadLoadBalancerAttributesRequest
type ReadLoadBalancerAttributesRequest struct {
	LoadBalancerName string `json:"LoadBalancerName,omitempty"`
}

// implements the service definition of ReadLoadBalancerAttributesResponse
type ReadLoadBalancerAttributesResponse struct {
	AccessLog       AccessLog       `json:"AccessLog,omitempty"`
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of ReadLoadBalancersRequest
type ReadLoadBalancersRequest struct {
	LoadBalancerNames []string `json:"LoadBalancerNames,omitempty"`
	MaxResults        int64    `json:"MaxResults,omitempty"`
	NextResultsToken  string   `json:"NextResultsToken,omitempty"`
}

// implements the service definition of ReadLoadBalancersResponse
type ReadLoadBalancersResponse struct {
	LoadBalancers    []LoadBalancers `json:"LoadBalancers,omitempty"`
	NextResultsToken string          `json:"NextResultsToken,omitempty"`
	ResponseContext  ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of ReadNatServicesRequest
type ReadNatServicesRequest struct {
	Filters          []Filters `json:"Filters,omitempty"`
	MaxResults       int64     `json:"MaxResults,omitempty"`
	NatServiceIds    []string  `json:"NatServiceIds,omitempty"`
	NextResultsToken string    `json:"NextResultsToken,omitempty"`
}

// implements the service definition of ReadNatServicesResponse
type ReadNatServicesResponse struct {
	NatServices      []NatServices   `json:"NatServices,omitempty"`
	NextResultsToken string          `json:"NextResultsToken,omitempty"`
	ResponseContext  ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of ReadNetAccessesRequest
type ReadNetAccessesRequest struct {
	DryRun       bool      `json:"DryRun,omitempty"`
	Filters      []Filters `json:"Filters,omitempty"`
	NetAccessIds []string  `json:"NetAccessIds,omitempty"`
}

// implements the service definition of ReadNetAccessesResponse
type ReadNetAccessesResponse struct {
	NetAccesses     []NetAccesses   `json:"NetAccesses,omitempty"`
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of ReadNetInternetGatewaysRequest
type ReadNetInternetGatewaysRequest struct {
	DryRun                bool      `json:"DryRun,omitempty"`
	Filters               []Filters `json:"Filters,omitempty"`
	NetInternetGatewayIds []string  `json:"NetInternetGatewayIds,omitempty"`
}

// implements the service definition of ReadNetInternetGatewaysResponse
type ReadNetInternetGatewaysResponse struct {
	NetInternetGateways []NetInternetGateways `json:"NetInternetGateways,omitempty"`
	ResponseContext     ResponseContext       `json:"ResponseContext,omitempty"`
}

// implements the service definition of ReadNetOptionsRequest
type ReadNetOptionsRequest struct {
	NetId string `json:"NetId,omitempty"`
}

// implements the service definition of ReadNetOptionsResponse
type ReadNetOptionsResponse struct {
	FirewallRulesSetLogging FirewallRulesSetLogging `json:"FirewallRulesSetLogging,omitempty"`
	NetId                   string                  `json:"NetId,omitempty"`
	ResponseContext         ResponseContext         `json:"ResponseContext,omitempty"`
}

// implements the service definition of ReadNetPeeringsRequest
type ReadNetPeeringsRequest struct {
	DryRun      bool      `json:"DryRun,omitempty"`
	Filters     []Filters `json:"Filters,omitempty"`
	NetPeerings []string  `json:"NetPeerings,omitempty"`
}

// implements the service definition of ReadNetPeeringsResponse
type ReadNetPeeringsResponse struct {
	NetPeerings     []NetPeerings   `json:"NetPeerings,omitempty"`
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of ReadNetServicesRequest
type ReadNetServicesRequest struct {
	DryRun           bool   `json:"DryRun,omitempty"`
	MaxResults       int64  `json:"MaxResults,omitempty"`
	NextResultsToken string `json:"NextResultsToken,omitempty"`
}

// implements the service definition of ReadNetServicesResponse
type ReadNetServicesResponse struct {
	NextResultsToken string          `json:"NextResultsToken,omitempty"`
	ResponseContext  ResponseContext `json:"ResponseContext,omitempty"`
	ServiceNames     []string        `json:"ServiceNames,omitempty"`
}

// implements the service definition of ReadNetsRequest
type ReadNetsRequest struct {
	DryRun  bool    `json:"DryRun,omitempty"`
	Filters Filters `json:"Filters,omitempty"`
}

// implements the service definition of ReadNetsResponse
type ReadNetsResponse struct {
	Nets            []Nets          `json:"Nets,omitempty"`
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of ReadNicsRequest
type ReadNicsRequest struct {
	DryRun  bool      `json:"DryRun,omitempty"`
	Filters []Filters `json:"Filters,omitempty"`
	NicIds  []string  `json:"NicIds,omitempty"`
}

// implements the service definition of ReadNicsResponse
type ReadNicsResponse struct {
	Nics            []Nics          `json:"Nics,omitempty"`
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of ReadPoliciesRequest
type ReadPoliciesRequest struct {
	GroupName        string `json:"GroupName,omitempty"`
	IsLinked         bool   `json:"IsLinked,omitempty"`
	MaxResults       int64  `json:"MaxResults,omitempty"`
	NextResultsToken string `json:"NextResultsToken,omitempty"`
	Path             string `json:"Path,omitempty"`
	UserName         string `json:"UserName,omitempty"`
}

// implements the service definition of ReadPoliciesResponse
type ReadPoliciesResponse struct {
	NextResultsToken string          `json:"NextResultsToken,omitempty"`
	Policies         []Policies      `json:"Policies,omitempty"`
	ResponseContext  ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of ReadPrefixListsRequest
type ReadPrefixListsRequest struct {
	DryRun           bool      `json:"DryRun,omitempty"`
	Filters          []Filters `json:"Filters,omitempty"`
	MaxResults       int64     `json:"MaxResults,omitempty"`
	NextResultsToken string    `json:"NextResultsToken,omitempty"`
	PrefixListIds    []string  `json:"PrefixListIds,omitempty"`
}

// implements the service definition of ReadPrefixListsResponse
type ReadPrefixListsResponse struct {
	NextResultsToken string          `json:"NextResultsToken,omitempty"`
	PrefixLists      []PrefixLists   `json:"PrefixLists,omitempty"`
	ResponseContext  ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of ReadProductTypesRequest
type ReadProductTypesRequest struct {
	Filters []Filters `json:"Filters,omitempty"`
}

// implements the service definition of ReadProductTypesResponse
type ReadProductTypesResponse struct {
	ProductTypes    []ProductTypes  `json:"ProductTypes,omitempty"`
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of ReadPublicCatalogRequest
type ReadPublicCatalogRequest struct {
}

// implements the service definition of ReadPublicCatalogResponse
type ReadPublicCatalogResponse struct {
	Catalog         Catalog         `json:"Catalog,omitempty"`
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of ReadPublicIpRangesRequest
type ReadPublicIpRangesRequest struct {
}

// implements the service definition of ReadPublicIpRangesResponse
type ReadPublicIpRangesResponse struct {
	PublicIps       []string        `json:"PublicIps,omitempty"`
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

type ReadPublicIpsRequest struct {
	DryRun         bool     `json:"DryRun,omitempty"`
	Filters        Filters  `json:"Filters,omitempty"`
	PublicIps      []string `json:"PublicIps,omitempty"`
	ReservationIds []string `json:"ReservationIds,omitempty"`
}

// implements the service definition of ReadPublicIpsResponse
type ReadPublicIpsResponse struct {
	PublicIps       []PublicIps     `json:"PublicIps,omitempty"`
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of ReadQuotasRequest
type ReadQuotasRequest struct {
	DryRun           bool      `json:"DryRun,omitempty"`
	Filters          []Filters `json:"Filters,omitempty"`
	MaxResults       int64     `json:"MaxResults,omitempty"`
	NextResultsToken string    `json:"NextResultsToken,omitempty"`
	QuotaNames       []string  `json:"QuotaNames,omitempty"`
}

// implements the service definition of ReadQuotasResponse
type ReadQuotasResponse struct {
	NextResultsToken string          `json:"NextResultsToken,omitempty"`
	QuotaTypes       []QuotaTypes    `json:"QuotaTypes,omitempty"`
	ResponseContext  ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of ReadRegionsRequest
type ReadRegionsRequest struct {
	DryRun      bool      `json:"DryRun,omitempty"`
	Filters     []Filters `json:"Filters,omitempty"`
	RegionNames []string  `json:"RegionNames,omitempty"`
}

// implements the service definition of ReadRegionsResponse
type ReadRegionsResponse struct {
	Regions         []Regions       `json:"Regions,omitempty"`
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of ReadReservedVmOffersRequest
type ReadReservedVmOffersRequest struct {
	DryRun              bool      `json:"DryRun,omitempty"`
	Filters             []Filters `json:"Filters,omitempty"`
	MaxResults          int64     `json:"MaxResults,omitempty"`
	NextResultsToken    string    `json:"NextResultsToken,omitempty"`
	OfferingType        string    `json:"OfferingType,omitempty"`
	ProductType         string    `json:"ProductType,omitempty"`
	ReservedVmsOfferIds []string  `json:"ReservedVmsOfferIds,omitempty"`
	SubRegionName       string    `json:"SubRegionName,omitempty"`
	Tenancy             string    `json:"Tenancy,omitempty"`
	Type                string    `json:"Type,omitempty"`
}

// implements the service definition of ReadReservedVmOffersResponse
type ReadReservedVmOffersResponse struct {
	NextResultsToken  string              `json:"NextResultsToken,omitempty"`
	ReservedVmsOffers []ReservedVmsOffers `json:"ReservedVmsOffers,omitempty"`
	ResponseContext   ResponseContext     `json:"ResponseContext,omitempty"`
}

// implements the service definition of ReadReservedVmsRequest
type ReadReservedVmsRequest struct {
	DryRun         bool      `json:"DryRun,omitempty"`
	Filters        []Filters `json:"Filters,omitempty"`
	OfferingType   string    `json:"OfferingType,omitempty"`
	ReservedVmsIds []string  `json:"ReservedVmsIds,omitempty"`
	SubRegionName  string    `json:"SubRegionName,omitempty"`
}

// implements the service definition of ReadReservedVmsResponse
type ReadReservedVmsResponse struct {
	ReservedVms     []ReservedVms   `json:"ReservedVms,omitempty"`
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of ReadRouteTablesRequest
type ReadRouteTablesRequest struct {
	DryRun  bool    `json:"DryRun,omitempty"`
	Filters Filters `json:"Filters,omitempty"`
}

// implements the service definition of ReadRouteTablesResponse
type ReadRouteTablesResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
	RouteTables     []RouteTables   `json:"RouteTables,omitempty"`
}

// implements the service definition of ReadServerCertificatesRequest
type ReadServerCertificatesRequest struct {
	MaxResults       int64  `json:"MaxResults,omitempty"`
	NextResultsToken string `json:"NextResultsToken,omitempty"`
	Path             string `json:"Path,omitempty"`
}

// implements the service definition of ReadServerCertificatesResponse
type ReadServerCertificatesResponse struct {
	NextResultsToken   string               `json:"NextResultsToken,omitempty"`
	ResponseContext    ResponseContext      `json:"ResponseContext,omitempty"`
	ServerCertificates []ServerCertificates `json:"ServerCertificates,omitempty"`
}

// implements the service definition of ReadSitesRequest
type ReadSitesRequest struct {
}

// implements the service definition of ReadSitesResponse
type ReadSitesResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
	Sites           []Sites         `json:"Sites,omitempty"`
}

// implements the service definition of ReadSnapshotAttributeRequest
type ReadSnapshotAttributeRequest struct {
	Attribute  string `json:"Attribute,omitempty"`
	DryRun     bool   `json:"DryRun,omitempty"`
	SnapshotId string `json:"SnapshotId,omitempty"`
}

// implements the service definition of ReadSnapshotAttributeResponse
type ReadSnapshotAttributeResponse struct {
	PermissionsToCreateVolumes []PermissionsToCreateVolumes `json:"PermissionsToCreateVolumes,omitempty"`
	ResponseContext            ResponseContext              `json:"ResponseContext,omitempty"`
	SnapshotId                 string                       `json:"SnapshotId,omitempty"`
}

// implements the service definition of ReadSnapshotExportTasksRequest
type ReadSnapshotExportTasksRequest struct {
	TaskIds []string `json:"TaskIds,omitempty"`
}

// implements the service definition of ReadSnapshotExportTasksResponse
type ReadSnapshotExportTasksResponse struct {
	ResponseContext     ResponseContext       `json:"ResponseContext,omitempty"`
	SnapshotExportTasks []SnapshotExportTasks `json:"SnapshotExportTasks,omitempty"`
}

// implements the service definition of ReadSnapshotsRequest
type ReadSnapshotsRequest struct {
	AccountIds                 []string  `json:"AccountIds,omitempty"`
	DryRun                     bool      `json:"DryRun,omitempty"`
	Filters                    []Filters `json:"Filters,omitempty"`
	MaxResults                 int64     `json:"MaxResults,omitempty"`
	NextResultsToken           string    `json:"NextResultsToken,omitempty"`
	PermissionsToCreateVolumes []string  `json:"PermissionsToCreateVolumes,omitempty"`
	SnapshotIds                []string  `json:"SnapshotIds,omitempty"`
}

// implements the service definition of ReadSnapshotsResponse
type ReadSnapshotsResponse struct {
	NextResultsToken string          `json:"NextResultsToken,omitempty"`
	ResponseContext  ResponseContext `json:"ResponseContext,omitempty"`
	Snapshots        []Snapshots     `json:"Snapshots,omitempty"`
}

// implements the service definition of ReadSubRegionsRequest
type ReadSubRegionsRequest struct {
	DryRun         bool      `json:"DryRun,omitempty"`
	Filters        []Filters `json:"Filters,omitempty"`
	SubRegionNames []string  `json:"SubRegionNames,omitempty"`
}

// implements the service definition of ReadSubRegionsResponse
type ReadSubRegionsResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
	SubRegions      []SubRegions    `json:"SubRegions,omitempty"`
}

// implements the service definition of ReadSubnetsRequest
type ReadSubnetsRequest struct {
	DryRun    bool      `json:"DryRun,omitempty"`
	Filters   []Filters `json:"Filters,omitempty"`
	SubnetIds []string  `json:"SubnetIds,omitempty"`
}

// implements the service definition of ReadSubnetsResponse
type ReadSubnetsResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
	Subnets         []Subnets       `json:"Subnets,omitempty"`
}

// implements the service definition of ReadTagsRequest
type ReadTagsRequest struct {
	DryRun           bool      `json:"DryRun,omitempty"`
	Filters          []Filters `json:"Filters,omitempty"`
	MaxResults       int64     `json:"MaxResults,omitempty"`
	NextResultsToken string    `json:"NextResultsToken,omitempty"`
}

// implements the service definition of ReadTagsResponse
type ReadTagsResponse struct {
	NextResultsToken string          `json:"NextResultsToken,omitempty"`
	ResponseContext  ResponseContext `json:"ResponseContext,omitempty"`
	Tags             []Tags          `json:"Tags,omitempty"`
}

// implements the service definition of ReadUsersRequest
type ReadUsersRequest struct {
	MaxResults       int64  `json:"MaxResults,omitempty"`
	NextResultsToken string `json:"NextResultsToken,omitempty"`
	Path             string `json:"Path,omitempty"`
}

// implements the service definition of ReadUsersResponse
type ReadUsersResponse struct {
	NextResultsToken string          `json:"NextResultsToken,omitempty"`
	ResponseContext  ResponseContext `json:"ResponseContext,omitempty"`
	Users            []Users         `json:"Users,omitempty"`
}

// implements the service definition of ReadVmAttributeRequest
type ReadVmAttributeRequest struct {
	Attribute string `json:"Attribute,omitempty"`
	DryRun    bool   `json:"DryRun,omitempty"`
	VmId      string `json:"VmId,omitempty"`
}

// implements the service definition of ReadVmAttributeResponse
type ReadVmAttributeResponse struct {
	BlockDeviceMappings         []BlockDeviceMappings `json:"BlockDeviceMappings,omitempty"`
	BsuOptimized                bool                  `json:"BsuOptimized,omitempty"`
	DeletionProtection          bool                  `json:"DeletionProtection,omitempty"`
	FirewallRulesSets           []FirewallRulesSets   `json:"FirewallRulesSets,omitempty"`
	IsSourceDestChecked         bool                  `json:"IsSourceDestChecked,omitempty"`
	KeypairName                 string                `json:"KeypairName,omitempty"`
	ProductCodes                []ProductCodes        `json:"ProductCodes,omitempty"`
	ResponseContext             ResponseContext       `json:"ResponseContext,omitempty"`
	RootDeviceName              string                `json:"RootDeviceName,omitempty"`
	Type                        string                `json:"Type,omitempty"`
	UserData                    string                `json:"UserData,omitempty"`
	VmId                        string                `json:"VmId,omitempty"`
	VmInitiatedShutdownBehavior string                `json:"VmInitiatedShutdownBehavior,omitempty"`
}

// implements the service definition of ReadVmTypesRequest
type ReadVmTypesRequest struct {
	Filters []Filters `json:"Filters,omitempty"`
}

// implements the service definition of ReadVmTypesResponse
type ReadVmTypesResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
	Types           []Types         `json:"Types,omitempty"`
}

// implements the service definition of ReadVmsHealthRequest
type ReadVmsHealthRequest struct {
	BackendVmsIds    []string `json:"BackendVmsIds,omitempty"`
	LoadBalancerName string   `json:"LoadBalancerName,omitempty"`
}

// implements the service definition of ReadVmsHealthResponse
type ReadVmsHealthResponse struct {
	BackendVmsHealth []BackendVmsHealth `json:"BackendVmsHealth,omitempty"`
	ResponseContext  ResponseContext    `json:"ResponseContext,omitempty"`
}

// implements the service definition of ReadVmsRequest
type ReadVmsRequest struct {
	DryRun           bool      `json:"DryRun,omitempty"`
	Filters          []Filters `json:"Filters,omitempty"`
	MaxResults       int64     `json:"MaxResults,omitempty"`
	NextResultsToken string    `json:"NextResultsToken,omitempty"`
}

// implements the service definition of ReadVmsResponse
type ReadVmsResponse struct {
	NextResultsToken string          `json:"NextResultsToken,omitempty"`
	ResponseContext  ResponseContext `json:"ResponseContext,omitempty"`
	Vms              []Vms           `json:"Vms,omitempty"`
}

// implements the service definition of ReadVmsStateRequest
type ReadVmsStateRequest struct {
	AllVms           bool      `json:"AllVms,omitempty"`
	DryRun           bool      `json:"DryRun,omitempty"`
	Filters          []Filters `json:"Filters,omitempty"`
	MaxResults       int64     `json:"MaxResults,omitempty"`
	NextResultsToken string    `json:"NextResultsToken,omitempty"`
}

// implements the service definition of ReadVmsStateResponse
type ReadVmsStateResponse struct {
	NextResultsToken string          `json:"NextResultsToken,omitempty"`
	ResponseContext  ResponseContext `json:"ResponseContext,omitempty"`
	VmStates         []VmStates      `json:"VmStates,omitempty"`
}

// implements the service definition of ReadVolumesRequest
type ReadVolumesRequest struct {
	DryRun           bool    `json:"DryRun,omitempty"`
	Filters          Filters `json:"Filters,omitempty"`
	MaxResults       int64   `json:"MaxResults,omitempty"`
	NextResultsToken string  `json:"NextResultsToken,omitempty"`
}

// implements the service definition of ReadVolumesResponse
type ReadVolumesResponse struct {
	NextResultsToken string          `json:"NextResultsToken,omitempty"`
	ResponseContext  ResponseContext `json:"ResponseContext,omitempty"`
	Volumes          []Volumes       `json:"Volumes,omitempty"`
}

// implements the service definition of ReadVpnConnectionsRequest
type ReadVpnConnectionsRequest struct {
	DryRun           bool      `json:"DryRun,omitempty"`
	Filters          []Filters `json:"Filters,omitempty"`
	VpnConnectionIds []string  `json:"VpnConnectionIds,omitempty"`
}

// implements the service definition of ReadVpnConnectionsResponse
type ReadVpnConnectionsResponse struct {
	ResponseContext ResponseContext  `json:"ResponseContext,omitempty"`
	VpnConnections  []VpnConnections `json:"vpnConnections,omitempty"`
}

// implements the service definition of ReadVpnGatewaysRequest
type ReadVpnGatewaysRequest struct {
	DryRun        bool      `json:"DryRun,omitempty"`
	Filters       []Filters `json:"Filters,omitempty"`
	VpnGatewayIds []string  `json:"VpnGatewayIds,omitempty"`
}

// implements the service definition of ReadVpnGatewaysResponse
type ReadVpnGatewaysResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
	VpnGateways     []VpnGateways   `json:"VpnGateways,omitempty"`
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

// implements the service definition of RecurringCharges
type RecurringCharges struct {
	Frequency string `json:"Frequency,omitempty"`
}

// implements the service definition of RegionConfig
type RegionConfig struct {
	FromDate     string       `json:"FromDate,omitempty"`
	Regions      []Regions    `json:"Regions,omitempty"`
	TargetRegion TargetRegion `json:"TargetRegion,omitempty"`
}

// implements the service definition of Regions
type Regions struct {
	RegionEndpoint string `json:"RegionEndpoint,omitempty"`
	RegionName     string `json:"RegionName,omitempty"`
}

// implements the service definition of RegisterImageRequest
type RegisterImageRequest struct {
	Architecture        string                `json:"Architecture,omitempty"`
	BlockDeviceMappings []BlockDeviceMappings `json:"BlockDeviceMappings,omitempty"`
	Description         string                `json:"Description,omitempty"`
	DryRun              bool                  `json:"DryRun,omitempty"`
	Name                string                `json:"Name,omitempty"`
	OsuLocation         string                `json:"OsuLocation,omitempty"`
	RootDeviceName      string                `json:"RootDeviceName,omitempty"`
}

// implements the service definition of RegisterImageResponse
type RegisterImageResponse struct {
	ImageId         string          `json:"ImageId,omitempty"`
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of RegisterUserInGroupRequest
type RegisterUserInGroupRequest struct {
	GroupName string `json:"GroupName,omitempty"`
	UserName  string `json:"UserName,omitempty"`
}

// implements the service definition of RegisterUserInGroupResponse
type RegisterUserInGroupResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of RegisterVmsInListenerRuleRequest
type RegisterVmsInListenerRuleRequest struct {
	ListenerRuleName string   `json:"ListenerRuleName,omitempty"`
	VmIds            []string `json:"VmIds,omitempty"`
}

// implements the service definition of RegisterVmsInListenerRuleResponse
type RegisterVmsInListenerRuleResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
	VmIds           []string        `json:"VmIds,omitempty"`
}

// implements the service definition of RegisterVmsInLoadBalancerRequest
type RegisterVmsInLoadBalancerRequest struct {
	BackendVmsIds    []string `json:"BackendVmsIds,omitempty"`
	LoadBalancerName string   `json:"LoadBalancerName,omitempty"`
}

// implements the service definition of RegisterVmsInLoadBalancerResponse
type RegisterVmsInLoadBalancerResponse struct {
	BackendVmsIds   []string        `json:"BackendVmsIds,omitempty"`
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
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

// implements the service definition of Removals
type Removals struct {
	AccountId        string `json:"AccountId,omitempty"`
	GlobalPermission string `json:"GlobalPermission,omitempty"`
}

// implements the service definition of ReservedVms
type ReservedVms struct {
	CurrencyCode     string             `json:"CurrencyCode,omitempty"`
	OfferingType     string             `json:"OfferingType,omitempty"`
	ProductType      string             `json:"ProductType,omitempty"`
	RecurringCharges []RecurringCharges `json:"RecurringCharges,omitempty"`
	ReservedVmsId    string             `json:"ReservedVmsId,omitempty"`
	State            string             `json:"State,omitempty"`
	SubRegionName    string             `json:"SubRegionName,omitempty"`
	Tenancy          string             `json:"Tenancy,omitempty"`
	Type             string             `json:"Type,omitempty"`
	VmCount          int64              `json:"VmCount,omitempty"`
}

// implements the service definition of ReservedVmsOffers
type ReservedVmsOffers struct {
	CurrencyCode       string             `json:"CurrencyCode,omitempty"`
	Duration           int64              `json:"Duration,omitempty"`
	FixedPrice         int                `json:"FixedPrice,omitempty"`
	OfferingType       string             `json:"OfferingType,omitempty"`
	PricingDetails     []PricingDetails   `json:"PricingDetails,omitempty"`
	ProductType        string             `json:"ProductType,omitempty"`
	RecurringCharges   []RecurringCharges `json:"RecurringCharges,omitempty"`
	ReservedVmsOfferId string             `json:"ReservedVmsOfferId,omitempty"`
	SubRegionName      string             `json:"SubRegionName,omitempty"`
	Tenancy            string             `json:"Tenancy,omitempty"`
	Type               string             `json:"Type,omitempty"`
	UsagePrice         int                `json:"UsagePrice,omitempty"`
}

// implements the service definition of ResetAccountPasswordRequest
type ResetAccountPasswordRequest struct {
	Password      string `json:"Password,omitempty"`
	PasswordToken string `json:"PasswordToken,omitempty"`
}

// implements the service definition of ResetAccountPasswordResponse
type ResetAccountPasswordResponse struct {
	Email           string          `json:"Email,omitempty"`
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of ResponseContext
type ResponseContext struct {
	RequestId string `json:"RequestId,omitempty"`
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

// implements the service definition of SendResetPasswordEmailRequest
type SendResetPasswordEmailRequest struct {
	Email string `json:"Email,omitempty"`
}

// implements the service definition of SendResetPasswordEmailResponse
type SendResetPasswordEmailResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of ServerCertificate
type ServerCertificate struct {
	Path                  string `json:"Path,omitempty"`
	ServerCertificateId   string `json:"ServerCertificateId,omitempty"`
	ServerCertificateName string `json:"ServerCertificateName,omitempty"`
}

// implements the service definition of ServerCertificates
type ServerCertificates struct {
	Path                  string `json:"Path,omitempty"`
	ServerCertificateId   string `json:"ServerCertificateId,omitempty"`
	ServerCertificateName string `json:"ServerCertificateName,omitempty"`
}

// implements the service definition of Sites
type Sites struct {
	Code string `json:"Code,omitempty"`
	Name string `json:"Name,omitempty"`
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
	AccountId   string `json:"AccountId,omitempty"`
	Description string `json:"Description,omitempty"`
	Progress    string `json:"Progress,omitempty"`
	SnapshotId  string `json:"SnapshotId,omitempty"`
	State       string `json:"State,omitempty"`
	Tags        []Tags `json:"Tags,omitempty"`
	VolumeId    string `json:"VolumeId,omitempty"`
	VolumeSize  int64  `json:"VolumeSize,omitempty"`
}

// implements the service definition of SourceFirewallRulesSet
type SourceFirewallRulesSet struct {
	FirewallRulesSetAccountId string `json:"FirewallRulesSetAccountId,omitempty"`
	FirewallRulesSetName      string `json:"FirewallRulesSetName,omitempty"`
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
	Vms             []Vms           `json:"Vms,omitempty"`
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
	Vms             []Vms           `json:"Vms,omitempty"`
}

// implements the service definition of SubRegions
type SubRegions struct {
	RegionName    string `json:"RegionName,omitempty"`
	State         string `json:"State,omitempty"`
	SubRegionName string `json:"SubRegionName,omitempty"`
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

// implements the service definition of Tags
type Tags struct {
	Key   string `json:"Key,omitempty"`
	Value string `json:"Value,omitempty"`
}

// implements the service definition of TargetRegion
type TargetRegion struct {
	RegionDomain string `json:"RegionDomain,omitempty"`
	RegionId     string `json:"RegionId,omitempty"`
	RegionName   string `json:"RegionName,omitempty"`
}

// implements the service definition of Transition
type Transition struct {
	Code    string `json:"Code,omitempty"`
	Message string `json:"Message,omitempty"`
}

// implements the service definition of Types
type Types struct {
	IsBsuOptimized bool   `json:"IsBsuOptimized,omitempty"`
	MaxPrivateIps  int64  `json:"MaxPrivateIps,omitempty"`
	MemorySize     int64  `json:"MemorySize,omitempty"`
	Name           string `json:"Name,omitempty"`
	StorageCount   int64  `json:"StorageCount,omitempty"`
	StorageSize    int64  `json:"StorageSize,omitempty"`
	VcoreCount     int64  `json:"VcoreCount,omitempty"`
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

// implements the service definition of UnlinkPolicyRequest
type UnlinkPolicyRequest struct {
	GroupName string `json:"GroupName,omitempty"`
	PolicyId  string `json:"PolicyId,omitempty"`
	UserName  string `json:"UserName,omitempty"`
}

// implements the service definition of UnlinkPolicyResponse
type UnlinkPolicyResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of UnlinkPrivateIpsRequest
type UnlinkPrivateIpsRequest struct {
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

// implements the service definition of UpdateAccountRequest
type UpdateAccountRequest struct {
	City          string `json:"City,omitempty"`
	CompanyName   string `json:"CompanyName,omitempty"`
	Country       string `json:"Country,omitempty"`
	Email         string `json:"Email,omitempty"`
	FirstName     string `json:"FirstName,omitempty"`
	JobTitle      string `json:"JobTitle,omitempty"`
	LastName      string `json:"LastName,omitempty"`
	Mobile        string `json:"Mobile,omitempty"`
	Password      string `json:"Password,omitempty"`
	Phone         string `json:"Phone,omitempty"`
	StateProvince string `json:"StateProvince,omitempty"`
	VatNumber     string `json:"VatNumber,omitempty"`
	ZipCode       string `json:"ZipCode,omitempty"`
}

// implements the service definition of UpdateAccountResponse
type UpdateAccountResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of UpdateApiKeyRequest
type UpdateApiKeyRequest struct {
	ApiKeyId string `json:"ApiKeyId,omitempty"`
	State    string `json:"State,omitempty"`
}

// implements the service definition of UpdateApiKeyResponse
type UpdateApiKeyResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of UpdateGroupRequest
type UpdateGroupRequest struct {
	GroupName    string `json:"GroupName,omitempty"`
	NewGroupName string `json:"NewGroupName,omitempty"`
	NewPath      string `json:"NewPath,omitempty"`
}

// implements the service definition of UpdateGroupResponse
type UpdateGroupResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of UpdateHealthCheckRequest
type UpdateHealthCheckRequest struct {
	HealthCheck      HealthCheck `json:"HealthCheck,omitempty"`
	LoadBalancerName string      `json:"LoadBalancerName,omitempty"`
}

// implements the service definition of UpdateHealthCheckResponse
type UpdateHealthCheckResponse struct {
	HealthCheck     HealthCheck     `json:"HealthCheck,omitempty"`
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of UpdateImageAttributeRequest
type UpdateImageAttributeRequest struct {
	DryRun     bool       `json:"DryRun,omitempty"`
	ImageId    string     `json:"ImageId,omitempty"`
	Permission Permission `json:"Permission,omitempty"`
}

// implements the service definition of UpdateImageAttributeResponse
type UpdateImageAttributeResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of UpdateKeypairRequest
type UpdateKeypairRequest struct {
	KeypairName string `json:"KeypairName,omitempty"`
	PublicKey   string `json:"PublicKey,omitempty"`
}

// implements the service definition of UpdateKeypairResponse
type UpdateKeypairResponse struct {
	KeypairFingerprint string          `json:"KeypairFingerprint,omitempty"`
	KeypairName        string          `json:"KeypairName,omitempty"`
	ResponseContext    ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of UpdateListenerRuleRequest
type UpdateListenerRuleRequest struct {
	Attribute        string `json:"Attribute,omitempty"`
	ListenerRuleName string `json:"ListenerRuleName,omitempty"`
	Value            string `json:"Value,omitempty"`
}

// implements the service definition of UpdateListenerRuleResponse
type UpdateListenerRuleResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of UpdateLoadBalancerAttributesRequest
type UpdateLoadBalancerAttributesRequest struct {
	AccessLog        AccessLog `json:"AccessLog,omitempty"`
	LoadBalancerName string    `json:"LoadBalancerName,omitempty"`
}

// implements the service definition of UpdateLoadBalancerAttributesResponse
type UpdateLoadBalancerAttributesResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of UpdateLoadBalancerPoliciesRequest
type UpdateLoadBalancerPoliciesRequest struct {
	LoadBalancerName string   `json:"LoadBalancerName,omitempty"`
	LoadBalancerPort int64    `json:"LoadBalancerPort,omitempty"`
	PolicyNames      []string `json:"PolicyNames,omitempty"`
}

// implements the service definition of UpdateLoadBalancerPoliciesResponse
type UpdateLoadBalancerPoliciesResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of UpdateNetAccessRequest
type UpdateNetAccessRequest struct {
	AddRouteTableIds    []string `json:"AddRouteTableIds,omitempty"`
	DryRun              bool     `json:"DryRun,omitempty"`
	NetAccessId         string   `json:"NetAccessId,omitempty"`
	RemoveRouteTableIds []string `json:"RemoveRouteTableIds,omitempty"`
}

// implements the service definition of UpdateNetAccessResponse
type UpdateNetAccessResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of UpdateNetOptionsRequest
type UpdateNetOptionsRequest struct {
	FirewallRulesSetLogging FirewallRulesSetLogging `json:"FirewallRulesSetLogging,omitempty"`
	NetId                   string                  `json:"NetId,omitempty"`
}

// implements the service definition of UpdateNetOptionsResponse
type UpdateNetOptionsResponse struct {
	FirewallRulesSetLogging FirewallRulesSetLogging `json:"FirewallRulesSetLogging,omitempty"`
	NetId                   string                  `json:"NetId,omitempty"`
	ResponseContext         ResponseContext         `json:"ResponseContext,omitempty"`
}

// implements the service definition of UpdateNicAttributeRequest
type UpdateNicAttributeRequest struct {
	Description         string   `json:"Description,omitempty"`
	DryRun              bool     `json:"DryRun,omitempty"`
	FirewallRulesSetIds []string `json:"FirewallRulesSetIds,omitempty"`
	NicId               string   `json:"NicId,omitempty"`
	NicLink             NicLink  `json:"NicLink,omitempty"`
}

// implements the service definition of UpdateNicAttributeResponse
type UpdateNicAttributeResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of UpdateRoutePropagationRequest
type UpdateRoutePropagationRequest struct {
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

// implements the service definition of UpdateServerCertificateRequest
type UpdateServerCertificateRequest struct {
	NewPath                  string `json:"NewPath,omitempty"`
	NewServerCertificateName string `json:"NewServerCertificateName,omitempty"`
	ServerCertificateName    string `json:"ServerCertificateName,omitempty"`
}

// implements the service definition of UpdateServerCertificateResponse
type UpdateServerCertificateResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of UpdateSnapshotAttributeRequest
type UpdateSnapshotAttributeRequest struct {
	DryRun                     bool                       `json:"DryRun,omitempty"`
	PermissionsToCreateVolumes PermissionsToCreateVolumes `json:"PermissionsToCreateVolumes,omitempty"`
	SnapshotId                 string                     `json:"SnapshotId,omitempty"`
}

// implements the service definition of UpdateSnapshotAttributeResponse
type UpdateSnapshotAttributeResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of UpdateUserRequest
type UpdateUserRequest struct {
	NewPath     string `json:"NewPath,omitempty"`
	NewUserName string `json:"NewUserName,omitempty"`
	UserName    string `json:"UserName,omitempty"`
}

// implements the service definition of UpdateUserResponse
type UpdateUserResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of UpdateVmAttributeRequest
type UpdateVmAttributeRequest struct {
	BlockDeviceMappings         []BlockDeviceMappings `json:"BlockDeviceMappings,omitempty"`
	BsuOptimized                bool                  `json:"BsuOptimized,omitempty"`
	DeletionProtection          bool                  `json:"DeletionProtection,omitempty"`
	DryRun                      bool                  `json:"DryRun,omitempty"`
	FirewallRulesSetIds         []string              `json:"FirewallRulesSetIds,omitempty"`
	IsSourceDestChecked         bool                  `json:"IsSourceDestChecked,omitempty"`
	KeypairName                 string                `json:"KeypairName,omitempty"`
	Type                        string                `json:"Type,omitempty"`
	UserData                    string                `json:"UserData,omitempty"`
	VmId                        string                `json:"VmId,omitempty"`
	VmInitiatedShutdownBehavior string                `json:"VmInitiatedShutdownBehavior,omitempty"`
}

// implements the service definition of UpdateVmAttributeResponse
type UpdateVmAttributeResponse struct {
	ResponseContext ResponseContext `json:"ResponseContext,omitempty"`
}

// implements the service definition of User
type User struct {
	Path     string `json:"Path,omitempty"`
	UserId   string `json:"UserId,omitempty"`
	UserName string `json:"UserName,omitempty"`
}

// implements the service definition of Users
type Users struct {
	Path     string `json:"Path,omitempty"`
	UserId   string `json:"UserId,omitempty"`
	UserName string `json:"UserName,omitempty"`
}

// implements the service definition of VmStates
type VmStates struct {
	MaintenanceEvents []MaintenanceEvents `json:"MaintenanceEvents,omitempty"`
	SubRegionName     string              `json:"SubRegionName,omitempty"`
	VmId              string              `json:"VmId,omitempty"`
	VmState           string              `json:"VmState,omitempty"`
}

// implements the service definition of Vms
type Vms struct {
	Architecture        string                `json:"Architecture,omitempty"`
	BlockDeviceMappings []BlockDeviceMappings `json:"BlockDeviceMappings,omitempty"`
	BsuOptimized        bool                  `json:"BsuOptimized,omitempty"`
	ClientToken         string                `json:"ClientToken,omitempty"`
	Comment             string                `json:"Comment,omitempty"`
	FirewallRulesSets   []FirewallRulesSets   `json:"FirewallRulesSets,omitempty"`
	ImageId             string                `json:"ImageId,omitempty"`
	IsSourceDestChecked bool                  `json:"IsSourceDestChecked,omitempty"`
	KeypairName         string                `json:"KeypairName,omitempty"`
	LaunchNumber        int64                 `json:"LaunchNumber,omitempty"`
	NetId               string                `json:"NetId,omitempty"`
	Nics                []Nics                `json:"Nics,omitempty"`
	OsFamily            string                `json:"OsFamily,omitempty"`
	Placement           Placement             `json:"Placement,omitempty"`
	PrivateDnsName      string                `json:"PrivateDnsName,omitempty"`
	PrivateIp           string                `json:"PrivateIp,omitempty"`
	ProductCodes        []ProductCodes        `json:"ProductCodes,omitempty"`
	PublicDnsName       string                `json:"PublicDnsName,omitempty"`
	PublicIp            string                `json:"PublicIp,omitempty"`
	ReservationId       string                `json:"ReservationId,omitempty"`
	RootDeviceName      string                `json:"RootDeviceName,omitempty"`
	RootDeviceType      string                `json:"RootDeviceType,omitempty"`
	State               string                `json:"State,omitempty"`
	SubnetId            string                `json:"SubnetId,omitempty"`
	Tags                []Tags                `json:"Tags,omitempty"`
	Transition          Transition            `json:"Transition,omitempty"`
	Type                string                `json:"Type,omitempty"`
	VmId                string                `json:"VmId,omitempty"`
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

// implements the service definition of VpnConnection
type VpnConnection struct {
	ClientEndpointConfiguration string   `json:"ClientEndpointConfiguration,omitempty"`
	ClientEndpointId            string   `json:"ClientEndpointId,omitempty"`
	Routes                      []Routes `json:"Routes,omitempty"`
	State                       string   `json:"State,omitempty"`
	StaticRoutesOnly            bool     `json:"StaticRoutesOnly,omitempty"`
	Tags                        []Tags   `json:"Tags,omitempty"`
	Type                        string   `json:"Type,omitempty"`
	VpnConnectionId             string   `json:"VpnConnectionId,omitempty"`
	VpnGatewayId                string   `json:"VpnGatewayId,omitempty"`
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

// implements the service definition of With
type With struct {
	CallDuration       bool `json:"CallDuration,omitempty"`
	QueryAccessKey     bool `json:"QueryAccessKey,omitempty"`
	QueryApiName       bool `json:"QueryApiName,omitempty"`
	QueryApiVersion    bool `json:"QueryApiVersion,omitempty"`
	QueryCallName      bool `json:"QueryCallName,omitempty"`
	QueryDate          bool `json:"QueryDate,omitempty"`
	QueryIpAddress     bool `json:"QueryIpAddress,omitempty"`
	QueryRaw           bool `json:"QueryRaw,omitempty"`
	QuerySize          bool `json:"QuerySize,omitempty"`
	QueryUserAgent     bool `json:"QueryUserAgent,omitempty"`
	ResponseId         bool `json:"ResponseId,omitempty"`
	ResponseSize       bool `json:"ResponseSize,omitempty"`
	ResponseStatusCode bool `json:"ResponseStatusCode,omitempty"`
}

// implements the service definition of vpnConnections
type VpnConnections struct {
	ClientEndpointConfiguration string   `json:"ClientEndpointConfiguration,omitempty"`
	ClientEndpointId            string   `json:"ClientEndpointId,omitempty"`
	Routes                      []Routes `json:"Routes,omitempty"`
	State                       string   `json:"State,omitempty"`
	StaticRoutesOnly            bool     `json:"StaticRoutesOnly,omitempty"`
	Tags                        []Tags   `json:"Tags,omitempty"`
	Type                        string   `json:"Type,omitempty"`
	VpnConnectionId             string   `json:"VpnConnectionId,omitempty"`
	VpnGatewayId                string   `json:"VpnGatewayId,omitempty"`
}

// POST_AcceptNetPeeringParameters holds parameters to POST_AcceptNetPeering
type POST_AcceptNetPeeringParameters struct {
	Acceptnetpeeringrequest AcceptNetPeeringRequest `json:"acceptnetpeeringrequest,omitempty"`
}

// POST_AcceptNetPeeringResponses holds responses of POST_AcceptNetPeering
type POST_AcceptNetPeeringResponses struct {
	OK *AcceptNetPeeringResponse
}

// POST_AuthenticateAccountParameters holds parameters to POST_AuthenticateAccount
type POST_AuthenticateAccountParameters struct {
	Authenticateaccountrequest AuthenticateAccountRequest `json:"authenticateaccountrequest,omitempty"`
}

// POST_AuthenticateAccountResponses holds responses of POST_AuthenticateAccount
type POST_AuthenticateAccountResponses struct {
	OK *AuthenticateAccountResponse
}

// POST_CancelExportTaskParameters holds parameters to POST_CancelExportTask
type POST_CancelExportTaskParameters struct {
	Cancelexporttaskrequest CancelExportTaskRequest `json:"cancelexporttaskrequest,omitempty"`
}

// POST_CancelExportTaskResponses holds responses of POST_CancelExportTask
type POST_CancelExportTaskResponses struct {
	OK *CancelExportTaskResponse
}

// POST_CheckSignatureParameters holds parameters to POST_CheckSignature
type POST_CheckSignatureParameters struct {
	Checksignaturerequest CheckSignatureRequest `json:"checksignaturerequest,omitempty"`
}

// POST_CheckSignatureResponses holds responses of POST_CheckSignature
type POST_CheckSignatureResponses struct {
	OK *CheckSignatureResponse
}

// POST_CopyAccountParameters holds parameters to POST_CopyAccount
type POST_CopyAccountParameters struct {
	Copyaccountrequest CopyAccountRequest `json:"copyaccountrequest,omitempty"`
}

// POST_CopyAccountResponses holds responses of POST_CopyAccount
type POST_CopyAccountResponses struct {
	OK *CopyAccountResponse
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

// POST_CreateAccountParameters holds parameters to POST_CreateAccount
type POST_CreateAccountParameters struct {
	Createaccountrequest CreateAccountRequest `json:"createaccountrequest,omitempty"`
}

// POST_CreateAccountResponses holds responses of POST_CreateAccount
type POST_CreateAccountResponses struct {
	OK *CreateAccountResponse
}

// POST_CreateApiKeyParameters holds parameters to POST_CreateApiKey
type POST_CreateApiKeyParameters struct {
	Createapikeyrequest CreateApiKeyRequest `json:"createapikeyrequest,omitempty"`
}

// POST_CreateApiKeyResponses holds responses of POST_CreateApiKey
type POST_CreateApiKeyResponses struct {
	OK *CreateApiKeyResponse
}

// POST_CreateClientEndpointParameters holds parameters to POST_CreateClientEndpoint
type POST_CreateClientEndpointParameters struct {
	Createclientendpointrequest CreateClientEndpointRequest `json:"createclientendpointrequest,omitempty"`
}

// POST_CreateClientEndpointResponses holds responses of POST_CreateClientEndpoint
type POST_CreateClientEndpointResponses struct {
	OK *CreateClientEndpointResponse
}

// POST_CreateDhcpOptionsParameters holds parameters to POST_CreateDhcpOptions
type POST_CreateDhcpOptionsParameters struct {
	Createdhcpoptionsrequest CreateDhcpOptionsRequest `json:"createdhcpoptionsrequest,omitempty"`
}

// POST_CreateDhcpOptionsResponses holds responses of POST_CreateDhcpOptions
type POST_CreateDhcpOptionsResponses struct {
	OK *CreateDhcpOptionsResponse
}

// POST_CreateDirectLinkParameters holds parameters to POST_CreateDirectLink
type POST_CreateDirectLinkParameters struct {
	Createdirectlinkrequest CreateDirectLinkRequest `json:"createdirectlinkrequest,omitempty"`
}

// POST_CreateDirectLinkResponses holds responses of POST_CreateDirectLink
type POST_CreateDirectLinkResponses struct {
	OK *CreateDirectLinkResponse
}

// POST_CreateDirectLinkInterfaceParameters holds parameters to POST_CreateDirectLinkInterface
type POST_CreateDirectLinkInterfaceParameters struct {
	Createdirectlinkinterfacerequest CreateDirectLinkInterfaceRequest `json:"createdirectlinkinterfacerequest,omitempty"`
}

// POST_CreateDirectLinkInterfaceResponses holds responses of POST_CreateDirectLinkInterface
type POST_CreateDirectLinkInterfaceResponses struct {
	OK *CreateDirectLinkInterfaceResponse
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

// POST_CreateGroupParameters holds parameters to POST_CreateGroup
type POST_CreateGroupParameters struct {
	Creategrouprequest CreateGroupRequest `json:"creategrouprequest,omitempty"`
}

// POST_CreateGroupResponses holds responses of POST_CreateGroup
type POST_CreateGroupResponses struct {
	OK *CreateGroupResponse
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

// POST_CreateListenerRuleParameters holds parameters to POST_CreateListenerRule
type POST_CreateListenerRuleParameters struct {
	Createlistenerrulerequest CreateListenerRuleRequest `json:"createlistenerrulerequest,omitempty"`
}

// POST_CreateListenerRuleResponses holds responses of POST_CreateListenerRule
type POST_CreateListenerRuleResponses struct {
	OK *CreateListenerRuleResponse
}

// POST_CreateLoadBalancerParameters holds parameters to POST_CreateLoadBalancer
type POST_CreateLoadBalancerParameters struct {
	Createloadbalancerrequest CreateLoadBalancerRequest `json:"createloadbalancerrequest,omitempty"`
}

// POST_CreateLoadBalancerResponses holds responses of POST_CreateLoadBalancer
type POST_CreateLoadBalancerResponses struct {
	OK *CreateLoadBalancerResponse
}

// POST_CreateLoadBalancerListenersParameters holds parameters to POST_CreateLoadBalancerListeners
type POST_CreateLoadBalancerListenersParameters struct {
	Createloadbalancerlistenersrequest CreateLoadBalancerListenersRequest `json:"createloadbalancerlistenersrequest,omitempty"`
}

// POST_CreateLoadBalancerListenersResponses holds responses of POST_CreateLoadBalancerListeners
type POST_CreateLoadBalancerListenersResponses struct {
	OK *CreateLoadBalancerListenersResponse
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

// POST_CreateNetAccessParameters holds parameters to POST_CreateNetAccess
type POST_CreateNetAccessParameters struct {
	Createnetaccessrequest CreateNetAccessRequest `json:"createnetaccessrequest,omitempty"`
}

// POST_CreateNetAccessResponses holds responses of POST_CreateNetAccess
type POST_CreateNetAccessResponses struct {
	OK *CreateNetAccessResponse
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

// POST_CreatePolicyParameters holds parameters to POST_CreatePolicy
type POST_CreatePolicyParameters struct {
	Createpolicyrequest CreatePolicyRequest `json:"createpolicyrequest,omitempty"`
}

// POST_CreatePolicyResponses holds responses of POST_CreatePolicy
type POST_CreatePolicyResponses struct {
	OK *CreatePolicyResponse
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

// POST_CreateStickyCookiePolicyParameters holds parameters to POST_CreateStickyCookiePolicy
type POST_CreateStickyCookiePolicyParameters struct {
	Createstickycookiepolicyrequest CreateStickyCookiePolicyRequest `json:"createstickycookiepolicyrequest,omitempty"`
}

// POST_CreateStickyCookiePolicyResponses holds responses of POST_CreateStickyCookiePolicy
type POST_CreateStickyCookiePolicyResponses struct {
	OK *CreateStickyCookiePolicyResponse
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

// POST_CreateUserParameters holds parameters to POST_CreateUser
type POST_CreateUserParameters struct {
	Createuserrequest CreateUserRequest `json:"createuserrequest,omitempty"`
}

// POST_CreateUserResponses holds responses of POST_CreateUser
type POST_CreateUserResponses struct {
	OK *CreateUserResponse
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

// POST_CreateVpnConnectionParameters holds parameters to POST_CreateVpnConnection
type POST_CreateVpnConnectionParameters struct {
	Createvpnconnectionrequest CreateVpnConnectionRequest `json:"createvpnconnectionrequest,omitempty"`
}

// POST_CreateVpnConnectionResponses holds responses of POST_CreateVpnConnection
type POST_CreateVpnConnectionResponses struct {
	OK *CreateVpnConnectionResponse
}

// POST_CreateVpnConnectionRouteParameters holds parameters to POST_CreateVpnConnectionRoute
type POST_CreateVpnConnectionRouteParameters struct {
	Createvpnconnectionrouterequest CreateVpnConnectionRouteRequest `json:"createvpnconnectionrouterequest,omitempty"`
}

// POST_CreateVpnConnectionRouteResponses holds responses of POST_CreateVpnConnectionRoute
type POST_CreateVpnConnectionRouteResponses struct {
	OK *CreateVpnConnectionRouteResponse
}

// POST_CreateVpnGatewayParameters holds parameters to POST_CreateVpnGateway
type POST_CreateVpnGatewayParameters struct {
	Createvpngatewayrequest CreateVpnGatewayRequest `json:"createvpngatewayrequest,omitempty"`
}

// POST_CreateVpnGatewayResponses holds responses of POST_CreateVpnGateway
type POST_CreateVpnGatewayResponses struct {
	OK *CreateVpnGatewayResponse
}

// POST_DeleteApiKeyParameters holds parameters to POST_DeleteApiKey
type POST_DeleteApiKeyParameters struct {
	Deleteapikeyrequest DeleteApiKeyRequest `json:"deleteapikeyrequest,omitempty"`
}

// POST_DeleteApiKeyResponses holds responses of POST_DeleteApiKey
type POST_DeleteApiKeyResponses struct {
	OK *DeleteApiKeyResponse
}

// POST_DeleteClientEndpointParameters holds parameters to POST_DeleteClientEndpoint
type POST_DeleteClientEndpointParameters struct {
	Deleteclientendpointrequest DeleteClientEndpointRequest `json:"deleteclientendpointrequest,omitempty"`
}

// POST_DeleteClientEndpointResponses holds responses of POST_DeleteClientEndpoint
type POST_DeleteClientEndpointResponses struct {
	OK *DeleteClientEndpointResponse
}

// POST_DeleteDhcpOptionsParameters holds parameters to POST_DeleteDhcpOptions
type POST_DeleteDhcpOptionsParameters struct {
	Deletedhcpoptionsrequest DeleteDhcpOptionsRequest `json:"deletedhcpoptionsrequest,omitempty"`
}

// POST_DeleteDhcpOptionsResponses holds responses of POST_DeleteDhcpOptions
type POST_DeleteDhcpOptionsResponses struct {
	OK *DeleteDhcpOptionsResponse
}

// POST_DeleteDirectLinkParameters holds parameters to POST_DeleteDirectLink
type POST_DeleteDirectLinkParameters struct {
	Deletedirectlinkrequest DeleteDirectLinkRequest `json:"deletedirectlinkrequest,omitempty"`
}

// POST_DeleteDirectLinkResponses holds responses of POST_DeleteDirectLink
type POST_DeleteDirectLinkResponses struct {
	OK *DeleteDirectLinkResponse
}

// POST_DeleteDirectLinkInterfaceParameters holds parameters to POST_DeleteDirectLinkInterface
type POST_DeleteDirectLinkInterfaceParameters struct {
	Deletedirectlinkinterfacerequest DeleteDirectLinkInterfaceRequest `json:"deletedirectlinkinterfacerequest,omitempty"`
}

// POST_DeleteDirectLinkInterfaceResponses holds responses of POST_DeleteDirectLinkInterface
type POST_DeleteDirectLinkInterfaceResponses struct {
	OK *DeleteDirectLinkInterfaceResponse
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

// POST_DeleteGroupParameters holds parameters to POST_DeleteGroup
type POST_DeleteGroupParameters struct {
	Deletegrouprequest DeleteGroupRequest `json:"deletegrouprequest,omitempty"`
}

// POST_DeleteGroupResponses holds responses of POST_DeleteGroup
type POST_DeleteGroupResponses struct {
	OK *DeleteGroupResponse
}

// POST_DeleteKeypairParameters holds parameters to POST_DeleteKeypair
type POST_DeleteKeypairParameters struct {
	Deletekeypairrequest DeleteKeypairRequest `json:"deletekeypairrequest,omitempty"`
}

// POST_DeleteKeypairResponses holds responses of POST_DeleteKeypair
type POST_DeleteKeypairResponses struct {
	OK *DeleteKeypairResponse
}

// POST_DeleteListenerRuleParameters holds parameters to POST_DeleteListenerRule
type POST_DeleteListenerRuleParameters struct {
	Deletelistenerrulerequest DeleteListenerRuleRequest `json:"deletelistenerrulerequest,omitempty"`
}

// POST_DeleteListenerRuleResponses holds responses of POST_DeleteListenerRule
type POST_DeleteListenerRuleResponses struct {
	OK *DeleteListenerRuleResponse
}

// POST_DeleteLoadBalancerParameters holds parameters to POST_DeleteLoadBalancer
type POST_DeleteLoadBalancerParameters struct {
	Deleteloadbalancerrequest DeleteLoadBalancerRequest `json:"deleteloadbalancerrequest,omitempty"`
}

// POST_DeleteLoadBalancerResponses holds responses of POST_DeleteLoadBalancer
type POST_DeleteLoadBalancerResponses struct {
	OK *DeleteLoadBalancerResponse
}

// POST_DeleteLoadBalancerListenersParameters holds parameters to POST_DeleteLoadBalancerListeners
type POST_DeleteLoadBalancerListenersParameters struct {
	Deleteloadbalancerlistenersrequest DeleteLoadBalancerListenersRequest `json:"deleteloadbalancerlistenersrequest,omitempty"`
}

// POST_DeleteLoadBalancerListenersResponses holds responses of POST_DeleteLoadBalancerListeners
type POST_DeleteLoadBalancerListenersResponses struct {
	OK *DeleteLoadBalancerListenersResponse
}

// POST_DeleteLoadBalancerPolicyParameters holds parameters to POST_DeleteLoadBalancerPolicy
type POST_DeleteLoadBalancerPolicyParameters struct {
	Deleteloadbalancerpolicyrequest DeleteLoadBalancerPolicyRequest `json:"deleteloadbalancerpolicyrequest,omitempty"`
}

// POST_DeleteLoadBalancerPolicyResponses holds responses of POST_DeleteLoadBalancerPolicy
type POST_DeleteLoadBalancerPolicyResponses struct {
	OK *DeleteLoadBalancerPolicyResponse
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

// POST_DeletePolicyParameters holds parameters to POST_DeletePolicy
type POST_DeletePolicyParameters struct {
	Deletepolicyrequest DeletePolicyRequest `json:"deletepolicyrequest,omitempty"`
}

// POST_DeletePolicyResponses holds responses of POST_DeletePolicy
type POST_DeletePolicyResponses struct {
	OK *DeletePolicyResponse
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

// POST_DeleteServerCertificateParameters holds parameters to POST_DeleteServerCertificate
type POST_DeleteServerCertificateParameters struct {
	Deleteservercertificaterequest DeleteServerCertificateRequest `json:"deleteservercertificaterequest,omitempty"`
}

// POST_DeleteServerCertificateResponses holds responses of POST_DeleteServerCertificate
type POST_DeleteServerCertificateResponses struct {
	OK *DeleteServerCertificateResponse
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

// POST_DeleteUserParameters holds parameters to POST_DeleteUser
type POST_DeleteUserParameters struct {
	Deleteuserrequest DeleteUserRequest `json:"deleteuserrequest,omitempty"`
}

// POST_DeleteUserResponses holds responses of POST_DeleteUser
type POST_DeleteUserResponses struct {
	OK *DeleteUserResponse
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

// POST_DeleteVpcEndpointsParameters holds parameters to POST_DeleteVpcEndpoints
type POST_DeleteVpcEndpointsParameters struct {
	Deletevpcendpointsrequest DeleteVpcEndpointsRequest `json:"deletevpcendpointsrequest,omitempty"`
}

// POST_DeleteVpcEndpointsResponses holds responses of POST_DeleteVpcEndpoints
type POST_DeleteVpcEndpointsResponses struct {
	OK *DeleteVpcEndpointsResponse
}

// POST_DeleteVpnConnectionParameters holds parameters to POST_DeleteVpnConnection
type POST_DeleteVpnConnectionParameters struct {
	Deletevpnconnectionrequest DeleteVpnConnectionRequest `json:"deletevpnconnectionrequest,omitempty"`
}

// POST_DeleteVpnConnectionResponses holds responses of POST_DeleteVpnConnection
type POST_DeleteVpnConnectionResponses struct {
	OK *DeleteVpnConnectionResponse
}

// POST_DeleteVpnConnectionRouteParameters holds parameters to POST_DeleteVpnConnectionRoute
type POST_DeleteVpnConnectionRouteParameters struct {
	Deletevpnconnectionrouterequest DeleteVpnConnectionRouteRequest `json:"deletevpnconnectionrouterequest,omitempty"`
}

// POST_DeleteVpnConnectionRouteResponses holds responses of POST_DeleteVpnConnectionRoute
type POST_DeleteVpnConnectionRouteResponses struct {
	OK *DeleteVpnConnectionRouteResponse
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

// POST_DeregisterUserInGroupParameters holds parameters to POST_DeregisterUserInGroup
type POST_DeregisterUserInGroupParameters struct {
	Deregisteruseringrouprequest DeregisterUserInGroupRequest `json:"deregisteruseringrouprequest,omitempty"`
}

// POST_DeregisterUserInGroupResponses holds responses of POST_DeregisterUserInGroup
type POST_DeregisterUserInGroupResponses struct {
	OK *DeregisterUserInGroupResponse
}

// POST_DeregisterVmsInListenerRuleParameters holds parameters to POST_DeregisterVmsInListenerRule
type POST_DeregisterVmsInListenerRuleParameters struct {
	Deregistervmsinlistenerrulerequest DeregisterVmsInListenerRuleRequest `json:"deregistervmsinlistenerrulerequest,omitempty"`
}

// POST_DeregisterVmsInListenerRuleResponses holds responses of POST_DeregisterVmsInListenerRule
type POST_DeregisterVmsInListenerRuleResponses struct {
	OK *DeregisterVmsInListenerRuleResponse
}

// POST_DeregisterVmsInLoadBalancerParameters holds parameters to POST_DeregisterVmsInLoadBalancer
type POST_DeregisterVmsInLoadBalancerParameters struct {
	Deregistervmsinloadbalancerrequest DeregisterVmsInLoadBalancerRequest `json:"deregistervmsinloadbalancerrequest,omitempty"`
}

// POST_DeregisterVmsInLoadBalancerResponses holds responses of POST_DeregisterVmsInLoadBalancer
type POST_DeregisterVmsInLoadBalancerResponses struct {
	OK *DeregisterVmsInLoadBalancerResponse
}

// POST_GetBillableDigestParameters holds parameters to POST_GetBillableDigest
type POST_GetBillableDigestParameters struct {
	Getbillabledigestrequest GetBillableDigestRequest `json:"getbillabledigestrequest,omitempty"`
}

// POST_GetBillableDigestResponses holds responses of POST_GetBillableDigest
type POST_GetBillableDigestResponses struct {
	OK *GetBillableDigestResponse
}

// POST_GetRegionConfigParameters holds parameters to POST_GetRegionConfig
type POST_GetRegionConfigParameters struct {
	Getregionconfigrequest GetRegionConfigRequest `json:"getregionconfigrequest,omitempty"`
}

// POST_GetRegionConfigResponses holds responses of POST_GetRegionConfig
type POST_GetRegionConfigResponses struct {
	OK *GetRegionConfigResponse
}

// POST_ImportKeyPairParameters holds parameters to POST_ImportKeyPair
type POST_ImportKeyPairParameters struct {
	Importkeypairrequest ImportKeyPairRequest `json:"importkeypairrequest,omitempty"`
}

// POST_ImportKeyPairResponses holds responses of POST_ImportKeyPair
type POST_ImportKeyPairResponses struct {
	OK *ImportKeyPairResponse
}

// POST_ImportServerCertificateParameters holds parameters to POST_ImportServerCertificate
type POST_ImportServerCertificateParameters struct {
	Importservercertificaterequest ImportServerCertificateRequest `json:"importservercertificaterequest,omitempty"`
}

// POST_ImportServerCertificateResponses holds responses of POST_ImportServerCertificate
type POST_ImportServerCertificateResponses struct {
	OK *ImportServerCertificateResponse
}

// POST_ImportSnaptShotParameters holds parameters to POST_ImportSnaptShot
type POST_ImportSnaptShotParameters struct {
	Importsnaptshotrequest ImportSnaptShotRequest `json:"importsnaptshotrequest,omitempty"`
}

// POST_ImportSnaptShotResponses holds responses of POST_ImportSnaptShot
type POST_ImportSnaptShotResponses struct {
	OK *ImportSnaptShotResponse
}

// POST_LinkDhcpOptionsParameters holds parameters to POST_LinkDhcpOptions
type POST_LinkDhcpOptionsParameters struct {
	Linkdhcpoptionsrequest LinkDhcpOptionsRequest `json:"linkdhcpoptionsrequest,omitempty"`
}

// POST_LinkDhcpOptionsResponses holds responses of POST_LinkDhcpOptions
type POST_LinkDhcpOptionsResponses struct {
	OK *LinkDhcpOptionsResponse
}

// POST_LinkLoadBalancerServerCertificateParameters holds parameters to POST_LinkLoadBalancerServerCertificate
type POST_LinkLoadBalancerServerCertificateParameters struct {
	Linkloadbalancerservercertificaterequest LinkLoadBalancerServerCertificateRequest `json:"linkloadbalancerservercertificaterequest,omitempty"`
}

// POST_LinkLoadBalancerServerCertificateResponses holds responses of POST_LinkLoadBalancerServerCertificate
type POST_LinkLoadBalancerServerCertificateResponses struct {
	OK *LinkLoadBalancerServerCertificateResponse
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

// POST_LinkPolicyParameters holds parameters to POST_LinkPolicy
type POST_LinkPolicyParameters struct {
	Linkpolicyrequest LinkPolicyRequest `json:"linkpolicyrequest,omitempty"`
}

// POST_LinkPolicyResponses holds responses of POST_LinkPolicy
type POST_LinkPolicyResponses struct {
	OK *LinkPolicyResponse
}

// POST_LinkPrivateIpParameters holds parameters to POST_LinkPrivateIp
type POST_LinkPrivateIpParameters struct {
	Linkprivateiprequest LinkPrivateIpRequest `json:"linkprivateiprequest,omitempty"`
}

// POST_LinkPrivateIpResponses holds responses of POST_LinkPrivateIp
type POST_LinkPrivateIpResponses struct {
	OK *LinkPrivateIpResponse
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

// POST_ListGroupsForUserParameters holds parameters to POST_ListGroupsForUser
type POST_ListGroupsForUserParameters struct {
	Listgroupsforuserrequest ListGroupsForUserRequest `json:"listgroupsforuserrequest,omitempty"`
}

// POST_ListGroupsForUserResponses holds responses of POST_ListGroupsForUser
type POST_ListGroupsForUserResponses struct {
	OK *ListGroupsForUserResponse
}

// POST_PurchaseReservedVmsOfferParameters holds parameters to POST_PurchaseReservedVmsOffer
type POST_PurchaseReservedVmsOfferParameters struct {
	Purchasereservedvmsofferrequest PurchaseReservedVmsOfferRequest `json:"purchasereservedvmsofferrequest,omitempty"`
}

// POST_PurchaseReservedVmsOfferResponses holds responses of POST_PurchaseReservedVmsOffer
type POST_PurchaseReservedVmsOfferResponses struct {
	OK *PurchaseReservedVmsOfferResponse
}

// POST_ReadAccountParameters holds parameters to POST_ReadAccount
type POST_ReadAccountParameters struct {
	Readaccountrequest ReadAccountRequest `json:"readaccountrequest,omitempty"`
}

// POST_ReadAccountResponses holds responses of POST_ReadAccount
type POST_ReadAccountResponses struct {
	OK *ReadAccountResponse
}

// POST_ReadAccountConsumptionParameters holds parameters to POST_ReadAccountConsumption
type POST_ReadAccountConsumptionParameters struct {
	Readaccountconsumptionrequest ReadAccountConsumptionRequest `json:"readaccountconsumptionrequest,omitempty"`
}

// POST_ReadAccountConsumptionResponses holds responses of POST_ReadAccountConsumption
type POST_ReadAccountConsumptionResponses struct {
	OK *ReadAccountConsumptionResponse
}

// POST_ReadAdminPasswordParameters holds parameters to POST_ReadAdminPassword
type POST_ReadAdminPasswordParameters struct {
	Readadminpasswordrequest ReadAdminPasswordRequest `json:"readadminpasswordrequest,omitempty"`
}

// POST_ReadAdminPasswordResponses holds responses of POST_ReadAdminPassword
type POST_ReadAdminPasswordResponses struct {
	OK *ReadAdminPasswordResponse
}

// POST_ReadApiKeysParameters holds parameters to POST_ReadApiKeys
type POST_ReadApiKeysParameters struct {
	Readapikeysrequest ReadApiKeysRequest `json:"readapikeysrequest,omitempty"`
}

// POST_ReadApiKeysResponses holds responses of POST_ReadApiKeys
type POST_ReadApiKeysResponses struct {
	OK *ReadApiKeysResponse
}

// POST_ReadApiLogsParameters holds parameters to POST_ReadApiLogs
type POST_ReadApiLogsParameters struct {
	Readapilogsrequest ReadApiLogsRequest `json:"readapilogsrequest,omitempty"`
}

// POST_ReadApiLogsResponses holds responses of POST_ReadApiLogs
type POST_ReadApiLogsResponses struct {
	OK *ReadApiLogsResponse
}

// POST_ReadCatalogParameters holds parameters to POST_ReadCatalog
type POST_ReadCatalogParameters struct {
	Readcatalogrequest ReadCatalogRequest `json:"readcatalogrequest,omitempty"`
}

// POST_ReadCatalogResponses holds responses of POST_ReadCatalog
type POST_ReadCatalogResponses struct {
	OK *ReadCatalogResponse
}

// POST_ReadClientEndpointsParameters holds parameters to POST_ReadClientEndpoints
type POST_ReadClientEndpointsParameters struct {
	Readclientendpointsrequest ReadClientEndpointsRequest `json:"readclientendpointsrequest,omitempty"`
}

// POST_ReadClientEndpointsResponses holds responses of POST_ReadClientEndpoints
type POST_ReadClientEndpointsResponses struct {
	OK *ReadClientEndpointsResponse
}

// POST_ReadConsoleOutputParameters holds parameters to POST_ReadConsoleOutput
type POST_ReadConsoleOutputParameters struct {
	Readconsoleoutputrequest ReadConsoleOutputRequest `json:"readconsoleoutputrequest,omitempty"`
}

// POST_ReadConsoleOutputResponses holds responses of POST_ReadConsoleOutput
type POST_ReadConsoleOutputResponses struct {
	OK *ReadConsoleOutputResponse
}

// POST_ReadDhcpOptionsParameters holds parameters to POST_ReadDhcpOptions
type POST_ReadDhcpOptionsParameters struct {
	Readdhcpoptionsrequest ReadDhcpOptionsRequest `json:"readdhcpoptionsrequest,omitempty"`
}

// POST_ReadDhcpOptionsResponses holds responses of POST_ReadDhcpOptions
type POST_ReadDhcpOptionsResponses struct {
	OK *ReadDhcpOptionsResponse
}

// POST_ReadDirectLinkInterfacesParameters holds parameters to POST_ReadDirectLinkInterfaces
type POST_ReadDirectLinkInterfacesParameters struct {
	Readdirectlinkinterfacesrequest ReadDirectLinkInterfacesRequest `json:"readdirectlinkinterfacesrequest,omitempty"`
}

// POST_ReadDirectLinkInterfacesResponses holds responses of POST_ReadDirectLinkInterfaces
type POST_ReadDirectLinkInterfacesResponses struct {
	OK *ReadDirectLinkInterfacesResponse
}

// POST_ReadDirectLinksParameters holds parameters to POST_ReadDirectLinks
type POST_ReadDirectLinksParameters struct {
	Readdirectlinksrequest ReadDirectLinksRequest `json:"readdirectlinksrequest,omitempty"`
}

// POST_ReadDirectLinksResponses holds responses of POST_ReadDirectLinks
type POST_ReadDirectLinksResponses struct {
	OK *ReadDirectLinksResponse
}

// POST_ReadFirewallRulesSetsParameters holds parameters to POST_ReadFirewallRulesSets
type POST_ReadFirewallRulesSetsParameters struct {
	Readfirewallrulessetsrequest ReadFirewallRulesSetsRequest `json:"readfirewallrulessetsrequest,omitempty"`
}

// POST_ReadFirewallRulesSetsResponses holds responses of POST_ReadFirewallRulesSets
type POST_ReadFirewallRulesSetsResponses struct {
	OK *ReadFirewallRulesSetsResponse
}

// POST_ReadGroupsParameters holds parameters to POST_ReadGroups
type POST_ReadGroupsParameters struct {
	Readgroupsrequest ReadGroupsRequest `json:"readgroupsrequest,omitempty"`
}

// POST_ReadGroupsResponses holds responses of POST_ReadGroups
type POST_ReadGroupsResponses struct {
	OK *ReadGroupsResponse
}

// POST_ReadImageAttributeParameters holds parameters to POST_ReadImageAttribute
type POST_ReadImageAttributeParameters struct {
	Readimageattributerequest ReadImageAttributeRequest `json:"readimageattributerequest,omitempty"`
}

// POST_ReadImageAttributeResponses holds responses of POST_ReadImageAttribute
type POST_ReadImageAttributeResponses struct {
	OK *ReadImageAttributeResponse
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

// POST_ReadListenerRulesParameters holds parameters to POST_ReadListenerRules
type POST_ReadListenerRulesParameters struct {
	Readlistenerrulesrequest ReadListenerRulesRequest `json:"readlistenerrulesrequest,omitempty"`
}

// POST_ReadListenerRulesResponses holds responses of POST_ReadListenerRules
type POST_ReadListenerRulesResponses struct {
	OK *ReadListenerRulesResponse
}

// POST_ReadLoadBalancerAttributesParameters holds parameters to POST_ReadLoadBalancerAttributes
type POST_ReadLoadBalancerAttributesParameters struct {
	Readloadbalancerattributesrequest ReadLoadBalancerAttributesRequest `json:"readloadbalancerattributesrequest,omitempty"`
}

// POST_ReadLoadBalancerAttributesResponses holds responses of POST_ReadLoadBalancerAttributes
type POST_ReadLoadBalancerAttributesResponses struct {
	OK *ReadLoadBalancerAttributesResponse
}

// POST_ReadLoadBalancersParameters holds parameters to POST_ReadLoadBalancers
type POST_ReadLoadBalancersParameters struct {
	Readloadbalancersrequest ReadLoadBalancersRequest `json:"readloadbalancersrequest,omitempty"`
}

// POST_ReadLoadBalancersResponses holds responses of POST_ReadLoadBalancers
type POST_ReadLoadBalancersResponses struct {
	OK *ReadLoadBalancersResponse
}

// POST_ReadNatServicesParameters holds parameters to POST_ReadNatServices
type POST_ReadNatServicesParameters struct {
	Readnatservicesrequest ReadNatServicesRequest `json:"readnatservicesrequest,omitempty"`
}

// POST_ReadNatServicesResponses holds responses of POST_ReadNatServices
type POST_ReadNatServicesResponses struct {
	OK *ReadNatServicesResponse
}

// POST_ReadNetAccessesParameters holds parameters to POST_ReadNetAccesses
type POST_ReadNetAccessesParameters struct {
	Readnetaccessesrequest ReadNetAccessesRequest `json:"readnetaccessesrequest,omitempty"`
}

// POST_ReadNetAccessesResponses holds responses of POST_ReadNetAccesses
type POST_ReadNetAccessesResponses struct {
	OK *ReadNetAccessesResponse
}

// POST_ReadNetInternetGatewaysParameters holds parameters to POST_ReadNetInternetGateways
type POST_ReadNetInternetGatewaysParameters struct {
	Readnetinternetgatewaysrequest ReadNetInternetGatewaysRequest `json:"readnetinternetgatewaysrequest,omitempty"`
}

// POST_ReadNetInternetGatewaysResponses holds responses of POST_ReadNetInternetGateways
type POST_ReadNetInternetGatewaysResponses struct {
	OK *ReadNetInternetGatewaysResponse
}

// POST_ReadNetOptionsParameters holds parameters to POST_ReadNetOptions
type POST_ReadNetOptionsParameters struct {
	Readnetoptionsrequest ReadNetOptionsRequest `json:"readnetoptionsrequest,omitempty"`
}

// POST_ReadNetOptionsResponses holds responses of POST_ReadNetOptions
type POST_ReadNetOptionsResponses struct {
	OK *ReadNetOptionsResponse
}

// POST_ReadNetPeeringsParameters holds parameters to POST_ReadNetPeerings
type POST_ReadNetPeeringsParameters struct {
	Readnetpeeringsrequest ReadNetPeeringsRequest `json:"readnetpeeringsrequest,omitempty"`
}

// POST_ReadNetPeeringsResponses holds responses of POST_ReadNetPeerings
type POST_ReadNetPeeringsResponses struct {
	OK *ReadNetPeeringsResponse
}

// POST_ReadNetServicesParameters holds parameters to POST_ReadNetServices
type POST_ReadNetServicesParameters struct {
	Readnetservicesrequest ReadNetServicesRequest `json:"readnetservicesrequest,omitempty"`
}

// POST_ReadNetServicesResponses holds responses of POST_ReadNetServices
type POST_ReadNetServicesResponses struct {
	OK *ReadNetServicesResponse
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

// POST_ReadPoliciesParameters holds parameters to POST_ReadPolicies
type POST_ReadPoliciesParameters struct {
	Readpoliciesrequest ReadPoliciesRequest `json:"readpoliciesrequest,omitempty"`
}

// POST_ReadPoliciesResponses holds responses of POST_ReadPolicies
type POST_ReadPoliciesResponses struct {
	OK *ReadPoliciesResponse
}

// POST_ReadPrefixListsParameters holds parameters to POST_ReadPrefixLists
type POST_ReadPrefixListsParameters struct {
	Readprefixlistsrequest ReadPrefixListsRequest `json:"readprefixlistsrequest,omitempty"`
}

// POST_ReadPrefixListsResponses holds responses of POST_ReadPrefixLists
type POST_ReadPrefixListsResponses struct {
	OK *ReadPrefixListsResponse
}

// POST_ReadProductTypesParameters holds parameters to POST_ReadProductTypes
type POST_ReadProductTypesParameters struct {
	Readproducttypesrequest ReadProductTypesRequest `json:"readproducttypesrequest,omitempty"`
}

// POST_ReadProductTypesResponses holds responses of POST_ReadProductTypes
type POST_ReadProductTypesResponses struct {
	OK *ReadProductTypesResponse
}

// POST_ReadPublicCatalogParameters holds parameters to POST_ReadPublicCatalog
type POST_ReadPublicCatalogParameters struct {
	Readpubliccatalogrequest ReadPublicCatalogRequest `json:"readpubliccatalogrequest,omitempty"`
}

// POST_ReadPublicCatalogResponses holds responses of POST_ReadPublicCatalog
type POST_ReadPublicCatalogResponses struct {
	OK *ReadPublicCatalogResponse
}

// POST_ReadPublicIpRangesParameters holds parameters to POST_ReadPublicIpRanges
type POST_ReadPublicIpRangesParameters struct {
	Readpubliciprangesrequest ReadPublicIpRangesRequest `json:"readpubliciprangesrequest,omitempty"`
}

// POST_ReadPublicIpRangesResponses holds responses of POST_ReadPublicIpRanges
type POST_ReadPublicIpRangesResponses struct {
	OK *ReadPublicIpRangesResponse
}

// POST_ReadPublicIpsParameters holds parameters to POST_ReadPublicIps
type POST_ReadPublicIpsParameters struct {
	Readpublicipsrequest ReadPublicIpsRequest `json:"readpublicipsrequest,omitempty"`
}

// POST_ReadPublicIpsResponses holds responses of POST_ReadPublicIps
type POST_ReadPublicIpsResponses struct {
	OK *ReadPublicIpsResponse
}

// POST_ReadQuotasParameters holds parameters to POST_ReadQuotas
type POST_ReadQuotasParameters struct {
	Readquotasrequest ReadQuotasRequest `json:"readquotasrequest,omitempty"`
}

// POST_ReadQuotasResponses holds responses of POST_ReadQuotas
type POST_ReadQuotasResponses struct {
	OK *ReadQuotasResponse
}

// POST_ReadRegionsParameters holds parameters to POST_ReadRegions
type POST_ReadRegionsParameters struct {
	Readregionsrequest ReadRegionsRequest `json:"readregionsrequest,omitempty"`
}

// POST_ReadRegionsResponses holds responses of POST_ReadRegions
type POST_ReadRegionsResponses struct {
	OK *ReadRegionsResponse
}

// POST_ReadReservedVmOffersParameters holds parameters to POST_ReadReservedVmOffers
type POST_ReadReservedVmOffersParameters struct {
	Readreservedvmoffersrequest ReadReservedVmOffersRequest `json:"readreservedvmoffersrequest,omitempty"`
}

// POST_ReadReservedVmOffersResponses holds responses of POST_ReadReservedVmOffers
type POST_ReadReservedVmOffersResponses struct {
	OK *ReadReservedVmOffersResponse
}

// POST_ReadReservedVmsParameters holds parameters to POST_ReadReservedVms
type POST_ReadReservedVmsParameters struct {
	Readreservedvmsrequest ReadReservedVmsRequest `json:"readreservedvmsrequest,omitempty"`
}

// POST_ReadReservedVmsResponses holds responses of POST_ReadReservedVms
type POST_ReadReservedVmsResponses struct {
	OK *ReadReservedVmsResponse
}

// POST_ReadRouteTablesParameters holds parameters to POST_ReadRouteTables
type POST_ReadRouteTablesParameters struct {
	Readroutetablesrequest ReadRouteTablesRequest `json:"readroutetablesrequest,omitempty"`
}

// POST_ReadRouteTablesResponses holds responses of POST_ReadRouteTables
type POST_ReadRouteTablesResponses struct {
	OK *ReadRouteTablesResponse
}

// POST_ReadServerCertificatesParameters holds parameters to POST_ReadServerCertificates
type POST_ReadServerCertificatesParameters struct {
	Readservercertificatesrequest ReadServerCertificatesRequest `json:"readservercertificatesrequest,omitempty"`
}

// POST_ReadServerCertificatesResponses holds responses of POST_ReadServerCertificates
type POST_ReadServerCertificatesResponses struct {
	OK *ReadServerCertificatesResponse
}

// POST_ReadSitesParameters holds parameters to POST_ReadSites
type POST_ReadSitesParameters struct {
	Readsitesrequest ReadSitesRequest `json:"readsitesrequest,omitempty"`
}

// POST_ReadSitesResponses holds responses of POST_ReadSites
type POST_ReadSitesResponses struct {
	OK *ReadSitesResponse
}

// POST_ReadSnapshotAttributeParameters holds parameters to POST_ReadSnapshotAttribute
type POST_ReadSnapshotAttributeParameters struct {
	Readsnapshotattributerequest ReadSnapshotAttributeRequest `json:"readsnapshotattributerequest,omitempty"`
}

// POST_ReadSnapshotAttributeResponses holds responses of POST_ReadSnapshotAttribute
type POST_ReadSnapshotAttributeResponses struct {
	OK *ReadSnapshotAttributeResponse
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

// POST_ReadSubRegionsParameters holds parameters to POST_ReadSubRegions
type POST_ReadSubRegionsParameters struct {
	Readsubregionsrequest ReadSubRegionsRequest `json:"readsubregionsrequest,omitempty"`
}

// POST_ReadSubRegionsResponses holds responses of POST_ReadSubRegions
type POST_ReadSubRegionsResponses struct {
	OK *ReadSubRegionsResponse
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

// POST_ReadUsersParameters holds parameters to POST_ReadUsers
type POST_ReadUsersParameters struct {
	Readusersrequest ReadUsersRequest `json:"readusersrequest,omitempty"`
}

// POST_ReadUsersResponses holds responses of POST_ReadUsers
type POST_ReadUsersResponses struct {
	OK *ReadUsersResponse
}

// POST_ReadVmAttributeParameters holds parameters to POST_ReadVmAttribute
type POST_ReadVmAttributeParameters struct {
	Readvmattributerequest ReadVmAttributeRequest `json:"readvmattributerequest,omitempty"`
}

// POST_ReadVmAttributeResponses holds responses of POST_ReadVmAttribute
type POST_ReadVmAttributeResponses struct {
	OK *ReadVmAttributeResponse
}

// POST_ReadVmTypesParameters holds parameters to POST_ReadVmTypes
type POST_ReadVmTypesParameters struct {
	Readvmtypesrequest ReadVmTypesRequest `json:"readvmtypesrequest,omitempty"`
}

// POST_ReadVmTypesResponses holds responses of POST_ReadVmTypes
type POST_ReadVmTypesResponses struct {
	OK *ReadVmTypesResponse
}

// POST_ReadVmsParameters holds parameters to POST_ReadVms
type POST_ReadVmsParameters struct {
	Readvmsrequest ReadVmsRequest `json:"readvmsrequest,omitempty"`
}

// POST_ReadVmsResponses holds responses of POST_ReadVms
type POST_ReadVmsResponses struct {
	OK *ReadVmsResponse
}

// POST_ReadVmsHealthParameters holds parameters to POST_ReadVmsHealth
type POST_ReadVmsHealthParameters struct {
	Readvmshealthrequest ReadVmsHealthRequest `json:"readvmshealthrequest,omitempty"`
}

// POST_ReadVmsHealthResponses holds responses of POST_ReadVmsHealth
type POST_ReadVmsHealthResponses struct {
	OK *ReadVmsHealthResponse
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

// POST_ReadVpnConnectionsParameters holds parameters to POST_ReadVpnConnections
type POST_ReadVpnConnectionsParameters struct {
	Readvpnconnectionsrequest ReadVpnConnectionsRequest `json:"readvpnconnectionsrequest,omitempty"`
}

// POST_ReadVpnConnectionsResponses holds responses of POST_ReadVpnConnections
type POST_ReadVpnConnectionsResponses struct {
	OK *ReadVpnConnectionsResponse
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

// POST_RegisterUserInGroupParameters holds parameters to POST_RegisterUserInGroup
type POST_RegisterUserInGroupParameters struct {
	Registeruseringrouprequest RegisterUserInGroupRequest `json:"registeruseringrouprequest,omitempty"`
}

// POST_RegisterUserInGroupResponses holds responses of POST_RegisterUserInGroup
type POST_RegisterUserInGroupResponses struct {
	OK *RegisterUserInGroupResponse
}

// POST_RegisterVmsInListenerRuleParameters holds parameters to POST_RegisterVmsInListenerRule
type POST_RegisterVmsInListenerRuleParameters struct {
	Registervmsinlistenerrulerequest RegisterVmsInListenerRuleRequest `json:"registervmsinlistenerrulerequest,omitempty"`
}

// POST_RegisterVmsInListenerRuleResponses holds responses of POST_RegisterVmsInListenerRule
type POST_RegisterVmsInListenerRuleResponses struct {
	OK *RegisterVmsInListenerRuleResponse
}

// POST_RegisterVmsInLoadBalancerParameters holds parameters to POST_RegisterVmsInLoadBalancer
type POST_RegisterVmsInLoadBalancerParameters struct {
	Registervmsinloadbalancerrequest RegisterVmsInLoadBalancerRequest `json:"registervmsinloadbalancerrequest,omitempty"`
}

// POST_RegisterVmsInLoadBalancerResponses holds responses of POST_RegisterVmsInLoadBalancer
type POST_RegisterVmsInLoadBalancerResponses struct {
	OK *RegisterVmsInLoadBalancerResponse
}

// POST_RejectNetPeeringParameters holds parameters to POST_RejectNetPeering
type POST_RejectNetPeeringParameters struct {
	Rejectnetpeeringrequest RejectNetPeeringRequest `json:"rejectnetpeeringrequest,omitempty"`
}

// POST_RejectNetPeeringResponses holds responses of POST_RejectNetPeering
type POST_RejectNetPeeringResponses struct {
	OK *RejectNetPeeringResponse
}

// POST_ResetAccountPasswordParameters holds parameters to POST_ResetAccountPassword
type POST_ResetAccountPasswordParameters struct {
	Resetaccountpasswordrequest ResetAccountPasswordRequest `json:"resetaccountpasswordrequest,omitempty"`
}

// POST_ResetAccountPasswordResponses holds responses of POST_ResetAccountPassword
type POST_ResetAccountPasswordResponses struct {
	OK *ResetAccountPasswordResponse
}

// POST_SendResetPasswordEmailParameters holds parameters to POST_SendResetPasswordEmail
type POST_SendResetPasswordEmailParameters struct {
	Sendresetpasswordemailrequest SendResetPasswordEmailRequest `json:"sendresetpasswordemailrequest,omitempty"`
}

// POST_SendResetPasswordEmailResponses holds responses of POST_SendResetPasswordEmail
type POST_SendResetPasswordEmailResponses struct {
	OK *SendResetPasswordEmailResponse
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

// POST_UnlinkPolicyParameters holds parameters to POST_UnlinkPolicy
type POST_UnlinkPolicyParameters struct {
	Unlinkpolicyrequest UnlinkPolicyRequest `json:"unlinkpolicyrequest,omitempty"`
}

// POST_UnlinkPolicyResponses holds responses of POST_UnlinkPolicy
type POST_UnlinkPolicyResponses struct {
	OK *UnlinkPolicyResponse
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

// POST_UpdateAccountParameters holds parameters to POST_UpdateAccount
type POST_UpdateAccountParameters struct {
	Updateaccountrequest UpdateAccountRequest `json:"updateaccountrequest,omitempty"`
}

// POST_UpdateAccountResponses holds responses of POST_UpdateAccount
type POST_UpdateAccountResponses struct {
	OK *UpdateAccountResponse
}

// POST_UpdateApiKeyParameters holds parameters to POST_UpdateApiKey
type POST_UpdateApiKeyParameters struct {
	Updateapikeyrequest UpdateApiKeyRequest `json:"updateapikeyrequest,omitempty"`
}

// POST_UpdateApiKeyResponses holds responses of POST_UpdateApiKey
type POST_UpdateApiKeyResponses struct {
	OK *UpdateApiKeyResponse
}

// POST_UpdateGroupParameters holds parameters to POST_UpdateGroup
type POST_UpdateGroupParameters struct {
	Updategrouprequest UpdateGroupRequest `json:"updategrouprequest,omitempty"`
}

// POST_UpdateGroupResponses holds responses of POST_UpdateGroup
type POST_UpdateGroupResponses struct {
	OK *UpdateGroupResponse
}

// POST_UpdateHealthCheckParameters holds parameters to POST_UpdateHealthCheck
type POST_UpdateHealthCheckParameters struct {
	Updatehealthcheckrequest UpdateHealthCheckRequest `json:"updatehealthcheckrequest,omitempty"`
}

// POST_UpdateHealthCheckResponses holds responses of POST_UpdateHealthCheck
type POST_UpdateHealthCheckResponses struct {
	OK *UpdateHealthCheckResponse
}

// POST_UpdateImageAttributeParameters holds parameters to POST_UpdateImageAttribute
type POST_UpdateImageAttributeParameters struct {
	Updateimageattributerequest UpdateImageAttributeRequest `json:"updateimageattributerequest,omitempty"`
}

// POST_UpdateImageAttributeResponses holds responses of POST_UpdateImageAttribute
type POST_UpdateImageAttributeResponses struct {
	OK *UpdateImageAttributeResponse
}

// POST_UpdateKeypairParameters holds parameters to POST_UpdateKeypair
type POST_UpdateKeypairParameters struct {
	Updatekeypairrequest UpdateKeypairRequest `json:"updatekeypairrequest,omitempty"`
}

// POST_UpdateKeypairResponses holds responses of POST_UpdateKeypair
type POST_UpdateKeypairResponses struct {
	OK *UpdateKeypairResponse
}

// POST_UpdateListenerRuleParameters holds parameters to POST_UpdateListenerRule
type POST_UpdateListenerRuleParameters struct {
	Updatelistenerrulerequest UpdateListenerRuleRequest `json:"updatelistenerrulerequest,omitempty"`
}

// POST_UpdateListenerRuleResponses holds responses of POST_UpdateListenerRule
type POST_UpdateListenerRuleResponses struct {
	OK *UpdateListenerRuleResponse
}

// POST_UpdateLoadBalancerAttributesParameters holds parameters to POST_UpdateLoadBalancerAttributes
type POST_UpdateLoadBalancerAttributesParameters struct {
	Updateloadbalancerattributesrequest UpdateLoadBalancerAttributesRequest `json:"updateloadbalancerattributesrequest,omitempty"`
}

// POST_UpdateLoadBalancerAttributesResponses holds responses of POST_UpdateLoadBalancerAttributes
type POST_UpdateLoadBalancerAttributesResponses struct {
	OK *UpdateLoadBalancerAttributesResponse
}

// POST_UpdateLoadBalancerPoliciesParameters holds parameters to POST_UpdateLoadBalancerPolicies
type POST_UpdateLoadBalancerPoliciesParameters struct {
	Updateloadbalancerpoliciesrequest UpdateLoadBalancerPoliciesRequest `json:"updateloadbalancerpoliciesrequest,omitempty"`
}

// POST_UpdateLoadBalancerPoliciesResponses holds responses of POST_UpdateLoadBalancerPolicies
type POST_UpdateLoadBalancerPoliciesResponses struct {
	OK *UpdateLoadBalancerPoliciesResponse
}

// POST_UpdateNetAccessParameters holds parameters to POST_UpdateNetAccess
type POST_UpdateNetAccessParameters struct {
	Updatenetaccessrequest UpdateNetAccessRequest `json:"updatenetaccessrequest,omitempty"`
}

// POST_UpdateNetAccessResponses holds responses of POST_UpdateNetAccess
type POST_UpdateNetAccessResponses struct {
	OK *UpdateNetAccessResponse
}

// POST_UpdateNetOptionsParameters holds parameters to POST_UpdateNetOptions
type POST_UpdateNetOptionsParameters struct {
	Updatenetoptionsrequest UpdateNetOptionsRequest `json:"updatenetoptionsrequest,omitempty"`
}

// POST_UpdateNetOptionsResponses holds responses of POST_UpdateNetOptions
type POST_UpdateNetOptionsResponses struct {
	OK *UpdateNetOptionsResponse
}

// POST_UpdateNicAttributeParameters holds parameters to POST_UpdateNicAttribute
type POST_UpdateNicAttributeParameters struct {
	Updatenicattributerequest UpdateNicAttributeRequest `json:"updatenicattributerequest,omitempty"`
}

// POST_UpdateNicAttributeResponses holds responses of POST_UpdateNicAttribute
type POST_UpdateNicAttributeResponses struct {
	OK *UpdateNicAttributeResponse
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

// POST_UpdateServerCertificateParameters holds parameters to POST_UpdateServerCertificate
type POST_UpdateServerCertificateParameters struct {
	Updateservercertificaterequest UpdateServerCertificateRequest `json:"updateservercertificaterequest,omitempty"`
}

// POST_UpdateServerCertificateResponses holds responses of POST_UpdateServerCertificate
type POST_UpdateServerCertificateResponses struct {
	OK *UpdateServerCertificateResponse
}

// POST_UpdateSnapshotAttributeParameters holds parameters to POST_UpdateSnapshotAttribute
type POST_UpdateSnapshotAttributeParameters struct {
	Updatesnapshotattributerequest UpdateSnapshotAttributeRequest `json:"updatesnapshotattributerequest,omitempty"`
}

// POST_UpdateSnapshotAttributeResponses holds responses of POST_UpdateSnapshotAttribute
type POST_UpdateSnapshotAttributeResponses struct {
	OK *UpdateSnapshotAttributeResponse
}

// POST_UpdateUserParameters holds parameters to POST_UpdateUser
type POST_UpdateUserParameters struct {
	Updateuserrequest UpdateUserRequest `json:"updateuserrequest,omitempty"`
}

// POST_UpdateUserResponses holds responses of POST_UpdateUser
type POST_UpdateUserResponses struct {
	OK *UpdateUserResponse
}

// POST_UpdateVmAttributeParameters holds parameters to POST_UpdateVmAttribute
type POST_UpdateVmAttributeParameters struct {
	Updatevmattributerequest UpdateVmAttributeRequest `json:"updatevmattributerequest,omitempty"`
}

// POST_UpdateVmAttributeResponses holds responses of POST_UpdateVmAttribute
type POST_UpdateVmAttributeResponses struct {
	OK *UpdateVmAttributeResponse
}
