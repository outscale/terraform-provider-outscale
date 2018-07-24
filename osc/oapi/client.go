// GENERATED FILE: DO NOT EDIT!

package oapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/signer/v4"
)

type Client struct {
	service string

	signer *v4.Signer

	client *http.Client

	config *Config
}

type Config struct {
	AccessKey string
	SecretKey string
	Region    string
	URL       string

	//Only Used for OAPI
	Service string

	// User agent for client
	UserAgent string
}

func (c Config) ServiceURL() string {
	s := fmt.Sprintf("https://%s.%s.%s/", c.Service, c.Region, c.URL)

	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}

	return u.String()
}

// NewClient creates an API client.
func NewClient(config *Config, c *http.Client) *Client {
	client := &Client{}
	client.service = config.ServiceURL()
	if c != nil {
		client.client = c
	} else {
		client.client = http.DefaultClient
	}

	s := &v4.Signer{
		Credentials: credentials.NewStaticCredentials(config.AccessKey,
			config.SecretKey, ""),
	}

	client.signer = s
	client.config = config

	return client
}

func (c *Client) Sign(req *http.Request, body []byte) error {
	reader := strings.NewReader(string(body))
	timestamp := time.Now()
	_, err := c.signer.Sign(req, reader, "oapi", c.config.Region, timestamp)
	return err

}

//
func (client *Client) POST_AcceptNetPeering(
	acceptnetpeeringrequest AcceptNetPeeringRequest,
) (
	response *POST_AcceptNetPeeringResponses,
	err error,
) {
	path := client.service + "/AcceptNetPeering"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(acceptnetpeeringrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_AcceptNetPeeringResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &AcceptNetPeeringResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_AuthenticateAccount(
	authenticateaccountrequest AuthenticateAccountRequest,
) (
	response *POST_AuthenticateAccountResponses,
	err error,
) {
	path := client.service + "/AuthenticateAccount"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(authenticateaccountrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_AuthenticateAccountResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &AuthenticateAccountResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_CancelExportTask(
	cancelexporttaskrequest CancelExportTaskRequest,
) (
	response *POST_CancelExportTaskResponses,
	err error,
) {
	path := client.service + "/CancelExportTask"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(cancelexporttaskrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_CancelExportTaskResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &CancelExportTaskResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_CheckSignature(
	checksignaturerequest CheckSignatureRequest,
) (
	response *POST_CheckSignatureResponses,
	err error,
) {
	path := client.service + "/CheckSignature"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(checksignaturerequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_CheckSignatureResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &CheckSignatureResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_CopyAccount(
	copyaccountrequest CopyAccountRequest,
) (
	response *POST_CopyAccountResponses,
	err error,
) {
	path := client.service + "/CopyAccount"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(copyaccountrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_CopyAccountResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &CopyAccountResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_CopyImage(
	copyimagerequest CopyImageRequest,
) (
	response *POST_CopyImageResponses,
	err error,
) {
	path := client.service + "/CopyImage"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(copyimagerequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_CopyImageResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &CopyImageResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_CopySnapshot(
	copysnapshotrequest CopySnapshotRequest,
) (
	response *POST_CopySnapshotResponses,
	err error,
) {
	path := client.service + "/CopySnapshot"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(copysnapshotrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_CopySnapshotResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &CopySnapshotResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_CreateAccount(
	createaccountrequest CreateAccountRequest,
) (
	response *POST_CreateAccountResponses,
	err error,
) {
	path := client.service + "/CreateAccount"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(createaccountrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_CreateAccountResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &CreateAccountResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_CreateApiKey(
	createapikeyrequest CreateApiKeyRequest,
) (
	response *POST_CreateApiKeyResponses,
	err error,
) {
	path := client.service + "/CreateApiKey"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(createapikeyrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_CreateApiKeyResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &CreateApiKeyResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_CreateClientEndpoint(
	createclientendpointrequest CreateClientEndpointRequest,
) (
	response *POST_CreateClientEndpointResponses,
	err error,
) {
	path := client.service + "/CreateClientEndpoint"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(createclientendpointrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_CreateClientEndpointResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &CreateClientEndpointResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_CreateDhcpOptions(
	createdhcpoptionsrequest CreateDhcpOptionsRequest,
) (
	response *POST_CreateDhcpOptionsResponses,
	err error,
) {
	path := client.service + "/CreateDhcpOptions"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(createdhcpoptionsrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_CreateDhcpOptionsResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &CreateDhcpOptionsResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_CreateDirectLink(
	createdirectlinkrequest CreateDirectLinkRequest,
) (
	response *POST_CreateDirectLinkResponses,
	err error,
) {
	path := client.service + "/CreateDirectLink"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(createdirectlinkrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_CreateDirectLinkResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &CreateDirectLinkResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_CreateDirectLinkInterface(
	createdirectlinkinterfacerequest CreateDirectLinkInterfaceRequest,
) (
	response *POST_CreateDirectLinkInterfaceResponses,
	err error,
) {
	path := client.service + "/CreateDirectLinkInterface"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(createdirectlinkinterfacerequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_CreateDirectLinkInterfaceResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &CreateDirectLinkInterfaceResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_CreateFirewallRuleInbound(
	createfirewallruleinboundrequest CreateFirewallRuleInboundRequest,
) (
	response *POST_CreateFirewallRuleInboundResponses,
	err error,
) {
	path := client.service + "/CreateFirewallRuleInbound"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(createfirewallruleinboundrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_CreateFirewallRuleInboundResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &CreateFirewallRuleInboundResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_CreateFirewallRuleOutbound(
	createfirewallruleoutboundrequest CreateFirewallRuleOutboundRequest,
) (
	response *POST_CreateFirewallRuleOutboundResponses,
	err error,
) {
	path := client.service + "/CreateFirewallRuleOutbound"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(createfirewallruleoutboundrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_CreateFirewallRuleOutboundResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &CreateFirewallRuleOutboundResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_CreateFirewallRulesSet(
	createfirewallrulessetrequest CreateFirewallRulesSetRequest,
) (
	response *POST_CreateFirewallRulesSetResponses,
	err error,
) {
	path := client.service + "/CreateFirewallRulesSet"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(createfirewallrulessetrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_CreateFirewallRulesSetResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &CreateFirewallRulesSetResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_CreateGroup(
	creategrouprequest CreateGroupRequest,
) (
	response *POST_CreateGroupResponses,
	err error,
) {
	path := client.service + "/CreateGroup"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(creategrouprequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_CreateGroupResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &CreateGroupResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_CreateImage(
	createimagerequest CreateImageRequest,
) (
	response *POST_CreateImageResponses,
	err error,
) {
	path := client.service + "/CreateImage"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(createimagerequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_CreateImageResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &CreateImageResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_CreateImageExportTask(
	createimageexporttaskrequest CreateImageExportTaskRequest,
) (
	response *POST_CreateImageExportTaskResponses,
	err error,
) {
	path := client.service + "/CreateImageExportTask"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(createimageexporttaskrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_CreateImageExportTaskResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &CreateImageExportTaskResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_CreateKeypair(
	createkeypairrequest CreateKeypairRequest,
) (
	response *POST_CreateKeypairResponses,
	err error,
) {
	path := client.service + "/CreateKeypair"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(createkeypairrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_CreateKeypairResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &CreateKeypairResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_CreateListenerRule(
	createlistenerrulerequest CreateListenerRuleRequest,
) (
	response *POST_CreateListenerRuleResponses,
	err error,
) {
	path := client.service + "/CreateListenerRule"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(createlistenerrulerequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_CreateListenerRuleResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &CreateListenerRuleResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_CreateLoadBalancer(
	createloadbalancerrequest CreateLoadBalancerRequest,
) (
	response *POST_CreateLoadBalancerResponses,
	err error,
) {
	path := client.service + "/CreateLoadBalancer"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(createloadbalancerrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_CreateLoadBalancerResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &CreateLoadBalancerResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_CreateLoadBalancerListeners(
	createloadbalancerlistenersrequest CreateLoadBalancerListenersRequest,
) (
	response *POST_CreateLoadBalancerListenersResponses,
	err error,
) {
	path := client.service + "/CreateLoadBalancerListeners"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(createloadbalancerlistenersrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_CreateLoadBalancerListenersResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &CreateLoadBalancerListenersResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_CreateNatService(
	createnatservicerequest CreateNatServiceRequest,
) (
	response *POST_CreateNatServiceResponses,
	err error,
) {
	path := client.service + "/CreateNatService"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(createnatservicerequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_CreateNatServiceResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &CreateNatServiceResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_CreateNet(
	createnetrequest CreateNetRequest,
) (
	response *POST_CreateNetResponses,
	err error,
) {
	path := client.service + "/CreateNet"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(createnetrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_CreateNetResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &CreateNetResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_CreateNetAccess(
	createnetaccessrequest CreateNetAccessRequest,
) (
	response *POST_CreateNetAccessResponses,
	err error,
) {
	path := client.service + "/CreateNetAccess"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(createnetaccessrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_CreateNetAccessResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &CreateNetAccessResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_CreateNetInternetGateway(
	createnetinternetgatewayrequest CreateNetInternetGatewayRequest,
) (
	response *POST_CreateNetInternetGatewayResponses,
	err error,
) {
	path := client.service + "/CreateNetInternetGateway"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(createnetinternetgatewayrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_CreateNetInternetGatewayResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &CreateNetInternetGatewayResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_CreateNetPeering(
	createnetpeeringrequest CreateNetPeeringRequest,
) (
	response *POST_CreateNetPeeringResponses,
	err error,
) {
	path := client.service + "/CreateNetPeering"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(createnetpeeringrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_CreateNetPeeringResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &CreateNetPeeringResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_CreateNic(
	createnicrequest CreateNicRequest,
) (
	response *POST_CreateNicResponses,
	err error,
) {
	path := client.service + "/CreateNic"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(createnicrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_CreateNicResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &CreateNicResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_CreatePolicy(
	createpolicyrequest CreatePolicyRequest,
) (
	response *POST_CreatePolicyResponses,
	err error,
) {
	path := client.service + "/CreatePolicy"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(createpolicyrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_CreatePolicyResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &CreatePolicyResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_CreatePublicIp(
	createpubliciprequest CreatePublicIpRequest,
) (
	response *POST_CreatePublicIpResponses,
	err error,
) {
	path := client.service + "/CreatePublicIp"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(createpubliciprequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_CreatePublicIpResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &CreatePublicIpResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_CreateRoute(
	createrouterequest CreateRouteRequest,
) (
	response *POST_CreateRouteResponses,
	err error,
) {
	path := client.service + "/CreateRoute"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(createrouterequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_CreateRouteResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &CreateRouteResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_CreateRouteTable(
	createroutetablerequest CreateRouteTableRequest,
) (
	response *POST_CreateRouteTableResponses,
	err error,
) {
	path := client.service + "/CreateRouteTable"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(createroutetablerequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_CreateRouteTableResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &CreateRouteTableResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_CreateSnapshot(
	createsnapshotrequest CreateSnapshotRequest,
) (
	response *POST_CreateSnapshotResponses,
	err error,
) {
	path := client.service + "/CreateSnapshot"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(createsnapshotrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_CreateSnapshotResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &CreateSnapshotResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_CreateSnapshotExportTask(
	createsnapshotexporttaskrequest CreateSnapshotExportTaskRequest,
) (
	response *POST_CreateSnapshotExportTaskResponses,
	err error,
) {
	path := client.service + "/CreateSnapshotExportTask"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(createsnapshotexporttaskrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_CreateSnapshotExportTaskResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &CreateSnapshotExportTaskResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_CreateStickyCookiePolicy(
	createstickycookiepolicyrequest CreateStickyCookiePolicyRequest,
) (
	response *POST_CreateStickyCookiePolicyResponses,
	err error,
) {
	path := client.service + "/CreateStickyCookiePolicy"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(createstickycookiepolicyrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_CreateStickyCookiePolicyResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &CreateStickyCookiePolicyResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_CreateSubnet(
	createsubnetrequest CreateSubnetRequest,
) (
	response *POST_CreateSubnetResponses,
	err error,
) {
	path := client.service + "/CreateSubnet"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(createsubnetrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_CreateSubnetResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &CreateSubnetResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_CreateTags(
	createtagsrequest CreateTagsRequest,
) (
	response *POST_CreateTagsResponses,
	err error,
) {
	path := client.service + "/CreateTags"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(createtagsrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_CreateTagsResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &CreateTagsResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_CreateUser(
	createuserrequest CreateUserRequest,
) (
	response *POST_CreateUserResponses,
	err error,
) {
	path := client.service + "/CreateUser"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(createuserrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_CreateUserResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &CreateUserResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_CreateVms(
	createvmsrequest CreateVmsRequest,
) (
	response *POST_CreateVmsResponses,
	err error,
) {
	path := client.service + "/CreateVms"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(createvmsrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_CreateVmsResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &CreateVmsResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_CreateVolume(
	createvolumerequest CreateVolumeRequest,
) (
	response *POST_CreateVolumeResponses,
	err error,
) {
	path := client.service + "/CreateVolume"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(createvolumerequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_CreateVolumeResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &CreateVolumeResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_CreateVpnConnection(
	createvpnconnectionrequest CreateVpnConnectionRequest,
) (
	response *POST_CreateVpnConnectionResponses,
	err error,
) {
	path := client.service + "/CreateVpnConnection"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(createvpnconnectionrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_CreateVpnConnectionResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &CreateVpnConnectionResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_CreateVpnConnectionRoute(
	createvpnconnectionrouterequest CreateVpnConnectionRouteRequest,
) (
	response *POST_CreateVpnConnectionRouteResponses,
	err error,
) {
	path := client.service + "/CreateVpnConnectionRoute"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(createvpnconnectionrouterequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_CreateVpnConnectionRouteResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &CreateVpnConnectionRouteResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_CreateVpnGateway(
	createvpngatewayrequest CreateVpnGatewayRequest,
) (
	response *POST_CreateVpnGatewayResponses,
	err error,
) {
	path := client.service + "/CreateVpnGateway"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(createvpngatewayrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_CreateVpnGatewayResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &CreateVpnGatewayResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_DeleteApiKey(
	deleteapikeyrequest DeleteApiKeyRequest,
) (
	response *POST_DeleteApiKeyResponses,
	err error,
) {
	path := client.service + "/DeleteApiKey"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(deleteapikeyrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_DeleteApiKeyResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &DeleteApiKeyResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_DeleteClientEndpoint(
	deleteclientendpointrequest DeleteClientEndpointRequest,
) (
	response *POST_DeleteClientEndpointResponses,
	err error,
) {
	path := client.service + "/DeleteClientEndpoint"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(deleteclientendpointrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_DeleteClientEndpointResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &DeleteClientEndpointResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_DeleteDhcpOptions(
	deletedhcpoptionsrequest DeleteDhcpOptionsRequest,
) (
	response *POST_DeleteDhcpOptionsResponses,
	err error,
) {
	path := client.service + "/DeleteDhcpOptions"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(deletedhcpoptionsrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_DeleteDhcpOptionsResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &DeleteDhcpOptionsResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_DeleteDirectLink(
	deletedirectlinkrequest DeleteDirectLinkRequest,
) (
	response *POST_DeleteDirectLinkResponses,
	err error,
) {
	path := client.service + "/DeleteDirectLink"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(deletedirectlinkrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_DeleteDirectLinkResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &DeleteDirectLinkResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_DeleteDirectLinkInterface(
	deletedirectlinkinterfacerequest DeleteDirectLinkInterfaceRequest,
) (
	response *POST_DeleteDirectLinkInterfaceResponses,
	err error,
) {
	path := client.service + "/DeleteDirectLinkInterface"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(deletedirectlinkinterfacerequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_DeleteDirectLinkInterfaceResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &DeleteDirectLinkInterfaceResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_DeleteFirewallRuleInbound(
	deletefirewallruleinboundrequest DeleteFirewallRuleInboundRequest,
) (
	response *POST_DeleteFirewallRuleInboundResponses,
	err error,
) {
	path := client.service + "/DeleteFirewallRuleInbound"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(deletefirewallruleinboundrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_DeleteFirewallRuleInboundResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &DeleteFirewallRuleInboundResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_DeleteFirewallRuleOutbound(
	deletefirewallruleoutboundrequest DeleteFirewallRuleOutboundRequest,
) (
	response *POST_DeleteFirewallRuleOutboundResponses,
	err error,
) {
	path := client.service + "/DeleteFirewallRuleOutbound"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(deletefirewallruleoutboundrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_DeleteFirewallRuleOutboundResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &DeleteFirewallRuleOutboundResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_DeleteFirewallRulesSet(
	deletefirewallrulessetrequest DeleteFirewallRulesSetRequest,
) (
	response *POST_DeleteFirewallRulesSetResponses,
	err error,
) {
	path := client.service + "/DeleteFirewallRulesSet"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(deletefirewallrulessetrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_DeleteFirewallRulesSetResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &DeleteFirewallRulesSetResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_DeleteGroup(
	deletegrouprequest DeleteGroupRequest,
) (
	response *POST_DeleteGroupResponses,
	err error,
) {
	path := client.service + "/DeleteGroup"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(deletegrouprequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_DeleteGroupResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &DeleteGroupResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_DeleteKeypair(
	deletekeypairrequest DeleteKeypairRequest,
) (
	response *POST_DeleteKeypairResponses,
	err error,
) {
	path := client.service + "/DeleteKeypair"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(deletekeypairrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_DeleteKeypairResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &DeleteKeypairResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_DeleteListenerRule(
	deletelistenerrulerequest DeleteListenerRuleRequest,
) (
	response *POST_DeleteListenerRuleResponses,
	err error,
) {
	path := client.service + "/DeleteListenerRule"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(deletelistenerrulerequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_DeleteListenerRuleResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &DeleteListenerRuleResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_DeleteLoadBalancer(
	deleteloadbalancerrequest DeleteLoadBalancerRequest,
) (
	response *POST_DeleteLoadBalancerResponses,
	err error,
) {
	path := client.service + "/DeleteLoadBalancer"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(deleteloadbalancerrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_DeleteLoadBalancerResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &DeleteLoadBalancerResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_DeleteLoadBalancerListeners(
	deleteloadbalancerlistenersrequest DeleteLoadBalancerListenersRequest,
) (
	response *POST_DeleteLoadBalancerListenersResponses,
	err error,
) {
	path := client.service + "/DeleteLoadBalancerListeners"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(deleteloadbalancerlistenersrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_DeleteLoadBalancerListenersResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &DeleteLoadBalancerListenersResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_DeleteLoadBalancerPolicy(
	deleteloadbalancerpolicyrequest DeleteLoadBalancerPolicyRequest,
) (
	response *POST_DeleteLoadBalancerPolicyResponses,
	err error,
) {
	path := client.service + "/DeleteLoadBalancerPolicy"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(deleteloadbalancerpolicyrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_DeleteLoadBalancerPolicyResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &DeleteLoadBalancerPolicyResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_DeleteNatService(
	deletenatservicerequest DeleteNatServiceRequest,
) (
	response *POST_DeleteNatServiceResponses,
	err error,
) {
	path := client.service + "/DeleteNatService"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(deletenatservicerequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_DeleteNatServiceResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &DeleteNatServiceResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_DeleteNet(
	deletenetrequest DeleteNetRequest,
) (
	response *POST_DeleteNetResponses,
	err error,
) {
	path := client.service + "/DeleteNet"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(deletenetrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_DeleteNetResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &DeleteNetResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_DeleteNetInternetGateway(
	deletenetinternetgatewayrequest DeleteNetInternetGatewayRequest,
) (
	response *POST_DeleteNetInternetGatewayResponses,
	err error,
) {
	path := client.service + "/DeleteNetInternetGateway"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(deletenetinternetgatewayrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_DeleteNetInternetGatewayResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &DeleteNetInternetGatewayResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_DeleteNetPeering(
	deletenetpeeringrequest DeleteNetPeeringRequest,
) (
	response *POST_DeleteNetPeeringResponses,
	err error,
) {
	path := client.service + "/DeleteNetPeering"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(deletenetpeeringrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_DeleteNetPeeringResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &DeleteNetPeeringResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_DeleteNic(
	deletenicrequest DeleteNicRequest,
) (
	response *POST_DeleteNicResponses,
	err error,
) {
	path := client.service + "/DeleteNic"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(deletenicrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_DeleteNicResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &DeleteNicResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_DeletePolicy(
	deletepolicyrequest DeletePolicyRequest,
) (
	response *POST_DeletePolicyResponses,
	err error,
) {
	path := client.service + "/DeletePolicy"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(deletepolicyrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_DeletePolicyResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &DeletePolicyResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_DeletePublicIp(
	deletepubliciprequest DeletePublicIpRequest,
) (
	response *POST_DeletePublicIpResponses,
	err error,
) {
	path := client.service + "/DeletePublicIp"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(deletepubliciprequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_DeletePublicIpResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &DeletePublicIpResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_DeleteRoute(
	deleterouterequest DeleteRouteRequest,
) (
	response *POST_DeleteRouteResponses,
	err error,
) {
	path := client.service + "/DeleteRoute"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(deleterouterequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_DeleteRouteResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &DeleteRouteResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_DeleteRouteTable(
	deleteroutetablerequest DeleteRouteTableRequest,
) (
	response *POST_DeleteRouteTableResponses,
	err error,
) {
	path := client.service + "/DeleteRouteTable"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(deleteroutetablerequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_DeleteRouteTableResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &DeleteRouteTableResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_DeleteServerCertificate(
	deleteservercertificaterequest DeleteServerCertificateRequest,
) (
	response *POST_DeleteServerCertificateResponses,
	err error,
) {
	path := client.service + "/DeleteServerCertificate"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(deleteservercertificaterequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_DeleteServerCertificateResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &DeleteServerCertificateResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_DeleteSnapshot(
	deletesnapshotrequest DeleteSnapshotRequest,
) (
	response *POST_DeleteSnapshotResponses,
	err error,
) {
	path := client.service + "/DeleteSnapshot"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(deletesnapshotrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_DeleteSnapshotResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &DeleteSnapshotResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_DeleteSubnet(
	deletesubnetrequest DeleteSubnetRequest,
) (
	response *POST_DeleteSubnetResponses,
	err error,
) {
	path := client.service + "/DeleteSubnet"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(deletesubnetrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_DeleteSubnetResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &DeleteSubnetResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_DeleteTags(
	deletetagsrequest DeleteTagsRequest,
) (
	response *POST_DeleteTagsResponses,
	err error,
) {
	path := client.service + "/DeleteTags"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(deletetagsrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_DeleteTagsResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &DeleteTagsResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_DeleteUser(
	deleteuserrequest DeleteUserRequest,
) (
	response *POST_DeleteUserResponses,
	err error,
) {
	path := client.service + "/DeleteUser"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(deleteuserrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_DeleteUserResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &DeleteUserResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_DeleteVms(
	deletevmsrequest DeleteVmsRequest,
) (
	response *POST_DeleteVmsResponses,
	err error,
) {
	path := client.service + "/DeleteVms"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(deletevmsrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_DeleteVmsResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &DeleteVmsResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_DeleteVolume(
	deletevolumerequest DeleteVolumeRequest,
) (
	response *POST_DeleteVolumeResponses,
	err error,
) {
	path := client.service + "/DeleteVolume"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(deletevolumerequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_DeleteVolumeResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &DeleteVolumeResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_DeleteVpcEndpoints(
	deletevpcendpointsrequest DeleteVpcEndpointsRequest,
) (
	response *POST_DeleteVpcEndpointsResponses,
	err error,
) {
	path := client.service + "/DeleteVpcEndpoints"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(deletevpcendpointsrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_DeleteVpcEndpointsResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &DeleteVpcEndpointsResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_DeleteVpnConnection(
	deletevpnconnectionrequest DeleteVpnConnectionRequest,
) (
	response *POST_DeleteVpnConnectionResponses,
	err error,
) {
	path := client.service + "/DeleteVpnConnection"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(deletevpnconnectionrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_DeleteVpnConnectionResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &DeleteVpnConnectionResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_DeleteVpnConnectionRoute(
	deletevpnconnectionrouterequest DeleteVpnConnectionRouteRequest,
) (
	response *POST_DeleteVpnConnectionRouteResponses,
	err error,
) {
	path := client.service + "/DeleteVpnConnectionRoute"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(deletevpnconnectionrouterequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_DeleteVpnConnectionRouteResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &DeleteVpnConnectionRouteResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_DeleteVpnGateway(
	deletevpngatewayrequest DeleteVpnGatewayRequest,
) (
	response *POST_DeleteVpnGatewayResponses,
	err error,
) {
	path := client.service + "/DeleteVpnGateway"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(deletevpngatewayrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_DeleteVpnGatewayResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &DeleteVpnGatewayResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_DeregisterImage(
	deregisterimagerequest DeregisterImageRequest,
) (
	response *POST_DeregisterImageResponses,
	err error,
) {
	path := client.service + "/DeregisterImage"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(deregisterimagerequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_DeregisterImageResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &DeregisterImageResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_DeregisterUserInGroup(
	deregisteruseringrouprequest DeregisterUserInGroupRequest,
) (
	response *POST_DeregisterUserInGroupResponses,
	err error,
) {
	path := client.service + "/DeregisterUserInGroup"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(deregisteruseringrouprequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_DeregisterUserInGroupResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &DeregisterUserInGroupResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_DeregisterVmsInListenerRule(
	deregistervmsinlistenerrulerequest DeregisterVmsInListenerRuleRequest,
) (
	response *POST_DeregisterVmsInListenerRuleResponses,
	err error,
) {
	path := client.service + "/DeregisterVmsInListenerRule"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(deregistervmsinlistenerrulerequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_DeregisterVmsInListenerRuleResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &DeregisterVmsInListenerRuleResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_DeregisterVmsInLoadBalancer(
	deregistervmsinloadbalancerrequest DeregisterVmsInLoadBalancerRequest,
) (
	response *POST_DeregisterVmsInLoadBalancerResponses,
	err error,
) {
	path := client.service + "/DeregisterVmsInLoadBalancer"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(deregistervmsinloadbalancerrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_DeregisterVmsInLoadBalancerResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &DeregisterVmsInLoadBalancerResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_GetBillableDigest(
	getbillabledigestrequest GetBillableDigestRequest,
) (
	response *POST_GetBillableDigestResponses,
	err error,
) {
	path := client.service + "/GetBillableDigest"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(getbillabledigestrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_GetBillableDigestResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &GetBillableDigestResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_GetRegionConfig(
	getregionconfigrequest GetRegionConfigRequest,
) (
	response *POST_GetRegionConfigResponses,
	err error,
) {
	path := client.service + "/GetRegionConfig"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(getregionconfigrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_GetRegionConfigResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &GetRegionConfigResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_ImportKeyPair(
	importkeypairrequest ImportKeyPairRequest,
) (
	response *POST_ImportKeyPairResponses,
	err error,
) {
	path := client.service + "/ImportKeyPair"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(importkeypairrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_ImportKeyPairResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &ImportKeyPairResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_ImportServerCertificate(
	importservercertificaterequest ImportServerCertificateRequest,
) (
	response *POST_ImportServerCertificateResponses,
	err error,
) {
	path := client.service + "/ImportServerCertificate"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(importservercertificaterequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_ImportServerCertificateResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &ImportServerCertificateResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_ImportSnaptShot(
	importsnaptshotrequest ImportSnaptShotRequest,
) (
	response *POST_ImportSnaptShotResponses,
	err error,
) {
	path := client.service + "/ImportSnaptShot"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(importsnaptshotrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_ImportSnaptShotResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &ImportSnaptShotResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_LinkDhcpOptions(
	linkdhcpoptionsrequest LinkDhcpOptionsRequest,
) (
	response *POST_LinkDhcpOptionsResponses,
	err error,
) {
	path := client.service + "/LinkDhcpOptions"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(linkdhcpoptionsrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_LinkDhcpOptionsResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &LinkDhcpOptionsResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_LinkLoadBalancerServerCertificate(
	linkloadbalancerservercertificaterequest LinkLoadBalancerServerCertificateRequest,
) (
	response *POST_LinkLoadBalancerServerCertificateResponses,
	err error,
) {
	path := client.service + "/LinkLoadBalancerServerCertificate"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(linkloadbalancerservercertificaterequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_LinkLoadBalancerServerCertificateResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &LinkLoadBalancerServerCertificateResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_LinkNetInternetGateway(
	linknetinternetgatewayrequest LinkNetInternetGatewayRequest,
) (
	response *POST_LinkNetInternetGatewayResponses,
	err error,
) {
	path := client.service + "/LinkNetInternetGateway"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(linknetinternetgatewayrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_LinkNetInternetGatewayResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &LinkNetInternetGatewayResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_LinkNic(
	linknicrequest LinkNicRequest,
) (
	response *POST_LinkNicResponses,
	err error,
) {
	path := client.service + "/LinkNic"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(linknicrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_LinkNicResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &LinkNicResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_LinkPolicy(
	linkpolicyrequest LinkPolicyRequest,
) (
	response *POST_LinkPolicyResponses,
	err error,
) {
	path := client.service + "/LinkPolicy"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(linkpolicyrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_LinkPolicyResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &LinkPolicyResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_LinkPrivateIp(
	linkprivateiprequest LinkPrivateIpRequest,
) (
	response *POST_LinkPrivateIpResponses,
	err error,
) {
	path := client.service + "/LinkPrivateIp"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(linkprivateiprequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_LinkPrivateIpResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &LinkPrivateIpResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_LinkPublicIp(
	linkpubliciprequest LinkPublicIpRequest,
) (
	response *POST_LinkPublicIpResponses,
	err error,
) {
	path := client.service + "/LinkPublicIp"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(linkpubliciprequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_LinkPublicIpResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &LinkPublicIpResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_LinkRouteTable(
	linkroutetablerequest LinkRouteTableRequest,
) (
	response *POST_LinkRouteTableResponses,
	err error,
) {
	path := client.service + "/LinkRouteTable"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(linkroutetablerequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_LinkRouteTableResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &LinkRouteTableResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_LinkVolume(
	linkvolumerequest LinkVolumeRequest,
) (
	response *POST_LinkVolumeResponses,
	err error,
) {
	path := client.service + "/LinkVolume"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(linkvolumerequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_LinkVolumeResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &LinkVolumeResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_LinkVpnGateway(
	linkvpngatewayrequest LinkVpnGatewayRequest,
) (
	response *POST_LinkVpnGatewayResponses,
	err error,
) {
	path := client.service + "/LinkVpnGateway"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(linkvpngatewayrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_LinkVpnGatewayResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &LinkVpnGatewayResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_ListGroupsForUser(
	listgroupsforuserrequest ListGroupsForUserRequest,
) (
	response *POST_ListGroupsForUserResponses,
	err error,
) {
	path := client.service + "/ListGroupsForUser"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(listgroupsforuserrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_ListGroupsForUserResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &ListGroupsForUserResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_PurchaseReservedVmsOffer(
	purchasereservedvmsofferrequest PurchaseReservedVmsOfferRequest,
) (
	response *POST_PurchaseReservedVmsOfferResponses,
	err error,
) {
	path := client.service + "/PurchaseReservedVmsOffer"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(purchasereservedvmsofferrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_PurchaseReservedVmsOfferResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &PurchaseReservedVmsOfferResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_ReadAccount(
	readaccountrequest ReadAccountRequest,
) (
	response *POST_ReadAccountResponses,
	err error,
) {
	path := client.service + "/ReadAccount"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(readaccountrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_ReadAccountResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &ReadAccountResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_ReadAccountConsumption(
	readaccountconsumptionrequest ReadAccountConsumptionRequest,
) (
	response *POST_ReadAccountConsumptionResponses,
	err error,
) {
	path := client.service + "/ReadAccountConsumption"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(readaccountconsumptionrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_ReadAccountConsumptionResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &ReadAccountConsumptionResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_ReadAdminPassword(
	readadminpasswordrequest ReadAdminPasswordRequest,
) (
	response *POST_ReadAdminPasswordResponses,
	err error,
) {
	path := client.service + "/ReadAdminPassword"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(readadminpasswordrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_ReadAdminPasswordResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &ReadAdminPasswordResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_ReadApiKeys(
	readapikeysrequest ReadApiKeysRequest,
) (
	response *POST_ReadApiKeysResponses,
	err error,
) {
	path := client.service + "/ReadApiKeys"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(readapikeysrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_ReadApiKeysResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &ReadApiKeysResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_ReadApiLogs(
	readapilogsrequest ReadApiLogsRequest,
) (
	response *POST_ReadApiLogsResponses,
	err error,
) {
	path := client.service + "/ReadApiLogs"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(readapilogsrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_ReadApiLogsResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &ReadApiLogsResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_ReadCatalog(
	readcatalogrequest ReadCatalogRequest,
) (
	response *POST_ReadCatalogResponses,
	err error,
) {
	path := client.service + "/ReadCatalog"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(readcatalogrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_ReadCatalogResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &ReadCatalogResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_ReadClientEndpoints(
	readclientendpointsrequest ReadClientEndpointsRequest,
) (
	response *POST_ReadClientEndpointsResponses,
	err error,
) {
	path := client.service + "/ReadClientEndpoints"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(readclientendpointsrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_ReadClientEndpointsResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &ReadClientEndpointsResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_ReadConsoleOutput(
	readconsoleoutputrequest ReadConsoleOutputRequest,
) (
	response *POST_ReadConsoleOutputResponses,
	err error,
) {
	path := client.service + "/ReadConsoleOutput"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(readconsoleoutputrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_ReadConsoleOutputResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &ReadConsoleOutputResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_ReadDhcpOptions(
	readdhcpoptionsrequest ReadDhcpOptionsRequest,
) (
	response *POST_ReadDhcpOptionsResponses,
	err error,
) {
	path := client.service + "/ReadDhcpOptions"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(readdhcpoptionsrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_ReadDhcpOptionsResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &ReadDhcpOptionsResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_ReadDirectLinkInterfaces(
	readdirectlinkinterfacesrequest ReadDirectLinkInterfacesRequest,
) (
	response *POST_ReadDirectLinkInterfacesResponses,
	err error,
) {
	path := client.service + "/ReadDirectLinkInterfaces"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(readdirectlinkinterfacesrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_ReadDirectLinkInterfacesResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &ReadDirectLinkInterfacesResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_ReadDirectLinks(
	readdirectlinksrequest ReadDirectLinksRequest,
) (
	response *POST_ReadDirectLinksResponses,
	err error,
) {
	path := client.service + "/ReadDirectLinks"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(readdirectlinksrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_ReadDirectLinksResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &ReadDirectLinksResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_ReadFirewallRulesSets(
	readfirewallrulessetsrequest ReadFirewallRulesSetsRequest,
) (
	response *POST_ReadFirewallRulesSetsResponses,
	err error,
) {
	path := client.service + "/ReadFirewallRulesSets"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(readfirewallrulessetsrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_ReadFirewallRulesSetsResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &ReadFirewallRulesSetsResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_ReadGroups(
	readgroupsrequest ReadGroupsRequest,
) (
	response *POST_ReadGroupsResponses,
	err error,
) {
	path := client.service + "/ReadGroups"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(readgroupsrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_ReadGroupsResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &ReadGroupsResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_ReadImageAttribute(
	readimageattributerequest ReadImageAttributeRequest,
) (
	response *POST_ReadImageAttributeResponses,
	err error,
) {
	path := client.service + "/ReadImageAttribute"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(readimageattributerequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_ReadImageAttributeResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &ReadImageAttributeResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_ReadImageExportTasks(
	readimageexporttasksrequest ReadImageExportTasksRequest,
) (
	response *POST_ReadImageExportTasksResponses,
	err error,
) {
	path := client.service + "/ReadImageExportTasks"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(readimageexporttasksrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_ReadImageExportTasksResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &ReadImageExportTasksResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_ReadImages(
	readimagesrequest ReadImagesRequest,
) (
	response *POST_ReadImagesResponses,
	err error,
) {
	path := client.service + "/ReadImages"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(readimagesrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_ReadImagesResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &ReadImagesResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_ReadKeypairs(
	readkeypairsrequest ReadKeypairsRequest,
) (
	response *POST_ReadKeypairsResponses,
	err error,
) {
	path := client.service + "/ReadKeypairs"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(readkeypairsrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_ReadKeypairsResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &ReadKeypairsResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_ReadListenerRules(
	readlistenerrulesrequest ReadListenerRulesRequest,
) (
	response *POST_ReadListenerRulesResponses,
	err error,
) {
	path := client.service + "/ReadListenerRules"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(readlistenerrulesrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_ReadListenerRulesResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &ReadListenerRulesResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_ReadLoadBalancerAttributes(
	readloadbalancerattributesrequest ReadLoadBalancerAttributesRequest,
) (
	response *POST_ReadLoadBalancerAttributesResponses,
	err error,
) {
	path := client.service + "/ReadLoadBalancerAttributes"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(readloadbalancerattributesrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_ReadLoadBalancerAttributesResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &ReadLoadBalancerAttributesResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_ReadLoadBalancers(
	readloadbalancersrequest ReadLoadBalancersRequest,
) (
	response *POST_ReadLoadBalancersResponses,
	err error,
) {
	path := client.service + "/ReadLoadBalancers"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(readloadbalancersrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_ReadLoadBalancersResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &ReadLoadBalancersResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_ReadNatServices(
	readnatservicesrequest ReadNatServicesRequest,
) (
	response *POST_ReadNatServicesResponses,
	err error,
) {
	path := client.service + "/ReadNatServices"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(readnatservicesrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_ReadNatServicesResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &ReadNatServicesResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_ReadNetAccesses(
	readnetaccessesrequest ReadNetAccessesRequest,
) (
	response *POST_ReadNetAccessesResponses,
	err error,
) {
	path := client.service + "/ReadNetAccesses"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(readnetaccessesrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_ReadNetAccessesResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &ReadNetAccessesResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_ReadNetInternetGateways(
	readnetinternetgatewaysrequest ReadNetInternetGatewaysRequest,
) (
	response *POST_ReadNetInternetGatewaysResponses,
	err error,
) {
	path := client.service + "/ReadNetInternetGateways"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(readnetinternetgatewaysrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_ReadNetInternetGatewaysResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &ReadNetInternetGatewaysResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_ReadNetOptions(
	readnetoptionsrequest ReadNetOptionsRequest,
) (
	response *POST_ReadNetOptionsResponses,
	err error,
) {
	path := client.service + "/ReadNetOptions"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(readnetoptionsrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_ReadNetOptionsResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &ReadNetOptionsResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_ReadNetPeerings(
	readnetpeeringsrequest ReadNetPeeringsRequest,
) (
	response *POST_ReadNetPeeringsResponses,
	err error,
) {
	path := client.service + "/ReadNetPeerings"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(readnetpeeringsrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_ReadNetPeeringsResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &ReadNetPeeringsResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_ReadNetServices(
	readnetservicesrequest ReadNetServicesRequest,
) (
	response *POST_ReadNetServicesResponses,
	err error,
) {
	path := client.service + "/ReadNetServices"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(readnetservicesrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_ReadNetServicesResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &ReadNetServicesResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_ReadNets(
	readnetsrequest ReadNetsRequest,
) (
	response *POST_ReadNetsResponses,
	err error,
) {
	path := client.service + "/ReadNets"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(readnetsrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_ReadNetsResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &ReadNetsResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_ReadNics(
	readnicsrequest ReadNicsRequest,
) (
	response *POST_ReadNicsResponses,
	err error,
) {
	path := client.service + "/ReadNics"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(readnicsrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_ReadNicsResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &ReadNicsResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_ReadPolicies(
	readpoliciesrequest ReadPoliciesRequest,
) (
	response *POST_ReadPoliciesResponses,
	err error,
) {
	path := client.service + "/ReadPolicies"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(readpoliciesrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_ReadPoliciesResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &ReadPoliciesResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_ReadPrefixLists(
	readprefixlistsrequest ReadPrefixListsRequest,
) (
	response *POST_ReadPrefixListsResponses,
	err error,
) {
	path := client.service + "/ReadPrefixLists"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(readprefixlistsrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_ReadPrefixListsResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &ReadPrefixListsResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_ReadProductTypes(
	readproducttypesrequest ReadProductTypesRequest,
) (
	response *POST_ReadProductTypesResponses,
	err error,
) {
	path := client.service + "/ReadProductTypes"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(readproducttypesrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_ReadProductTypesResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &ReadProductTypesResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_ReadPublicCatalog(
	readpubliccatalogrequest ReadPublicCatalogRequest,
) (
	response *POST_ReadPublicCatalogResponses,
	err error,
) {
	path := client.service + "/ReadPublicCatalog"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(readpubliccatalogrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_ReadPublicCatalogResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &ReadPublicCatalogResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_ReadPublicIpRanges(
	readpubliciprangesrequest ReadPublicIpRangesRequest,
) (
	response *POST_ReadPublicIpRangesResponses,
	err error,
) {
	path := client.service + "/ReadPublicIpRanges"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(readpubliciprangesrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_ReadPublicIpRangesResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &ReadPublicIpRangesResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_ReadPublicIps(
	readpublicipsrequest ReadPublicIpsRequest,
) (
	response *POST_ReadPublicIpsResponses,
	err error,
) {
	path := client.service + "/ReadPublicIps"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(readpublicipsrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_ReadPublicIpsResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &ReadPublicIpsResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_ReadQuotas(
	readquotasrequest ReadQuotasRequest,
) (
	response *POST_ReadQuotasResponses,
	err error,
) {
	path := client.service + "/ReadQuotas"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(readquotasrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_ReadQuotasResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &ReadQuotasResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_ReadRegions(
	readregionsrequest ReadRegionsRequest,
) (
	response *POST_ReadRegionsResponses,
	err error,
) {
	path := client.service + "/ReadRegions"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(readregionsrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_ReadRegionsResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &ReadRegionsResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_ReadReservedVmOffers(
	readreservedvmoffersrequest ReadReservedVmOffersRequest,
) (
	response *POST_ReadReservedVmOffersResponses,
	err error,
) {
	path := client.service + "/ReadReservedVmOffers"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(readreservedvmoffersrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_ReadReservedVmOffersResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &ReadReservedVmOffersResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_ReadReservedVms(
	readreservedvmsrequest ReadReservedVmsRequest,
) (
	response *POST_ReadReservedVmsResponses,
	err error,
) {
	path := client.service + "/ReadReservedVms"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(readreservedvmsrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_ReadReservedVmsResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &ReadReservedVmsResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_ReadRouteTables(
	readroutetablesrequest ReadRouteTablesRequest,
) (
	response *POST_ReadRouteTablesResponses,
	err error,
) {
	path := client.service + "/ReadRouteTables"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(readroutetablesrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_ReadRouteTablesResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &ReadRouteTablesResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_ReadServerCertificates(
	readservercertificatesrequest ReadServerCertificatesRequest,
) (
	response *POST_ReadServerCertificatesResponses,
	err error,
) {
	path := client.service + "/ReadServerCertificates"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(readservercertificatesrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_ReadServerCertificatesResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &ReadServerCertificatesResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_ReadSites(
	readsitesrequest ReadSitesRequest,
) (
	response *POST_ReadSitesResponses,
	err error,
) {
	path := client.service + "/ReadSites"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(readsitesrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_ReadSitesResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &ReadSitesResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_ReadSnapshotAttribute(
	readsnapshotattributerequest ReadSnapshotAttributeRequest,
) (
	response *POST_ReadSnapshotAttributeResponses,
	err error,
) {
	path := client.service + "/ReadSnapshotAttribute"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(readsnapshotattributerequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_ReadSnapshotAttributeResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &ReadSnapshotAttributeResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_ReadSnapshotExportTasks(
	readsnapshotexporttasksrequest ReadSnapshotExportTasksRequest,
) (
	response *POST_ReadSnapshotExportTasksResponses,
	err error,
) {
	path := client.service + "/ReadSnapshotExportTasks"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(readsnapshotexporttasksrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_ReadSnapshotExportTasksResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &ReadSnapshotExportTasksResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_ReadSnapshots(
	readsnapshotsrequest ReadSnapshotsRequest,
) (
	response *POST_ReadSnapshotsResponses,
	err error,
) {
	path := client.service + "/ReadSnapshots"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(readsnapshotsrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_ReadSnapshotsResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &ReadSnapshotsResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_ReadSubRegions(
	readsubregionsrequest ReadSubRegionsRequest,
) (
	response *POST_ReadSubRegionsResponses,
	err error,
) {
	path := client.service + "/ReadSubRegions"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(readsubregionsrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_ReadSubRegionsResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &ReadSubRegionsResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_ReadSubnets(
	readsubnetsrequest ReadSubnetsRequest,
) (
	response *POST_ReadSubnetsResponses,
	err error,
) {
	path := client.service + "/ReadSubnets"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(readsubnetsrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_ReadSubnetsResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &ReadSubnetsResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_ReadTags(
	readtagsrequest ReadTagsRequest,
) (
	response *POST_ReadTagsResponses,
	err error,
) {
	path := client.service + "/ReadTags"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(readtagsrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_ReadTagsResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &ReadTagsResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_ReadUsers(
	readusersrequest ReadUsersRequest,
) (
	response *POST_ReadUsersResponses,
	err error,
) {
	path := client.service + "/ReadUsers"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(readusersrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_ReadUsersResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &ReadUsersResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_ReadVmAttribute(
	readvmattributerequest ReadVmAttributeRequest,
) (
	response *POST_ReadVmAttributeResponses,
	err error,
) {
	path := client.service + "/ReadVmAttribute"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(readvmattributerequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_ReadVmAttributeResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &ReadVmAttributeResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_ReadVmTypes(
	readvmtypesrequest ReadVmTypesRequest,
) (
	response *POST_ReadVmTypesResponses,
	err error,
) {
	path := client.service + "/ReadVmTypes"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(readvmtypesrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_ReadVmTypesResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &ReadVmTypesResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_ReadVms(
	readvmsrequest ReadVmsRequest,
) (
	response *POST_ReadVmsResponses,
	err error,
) {
	path := client.service + "/ReadVms"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(readvmsrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_ReadVmsResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &ReadVmsResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_ReadVmsHealth(
	readvmshealthrequest ReadVmsHealthRequest,
) (
	response *POST_ReadVmsHealthResponses,
	err error,
) {
	path := client.service + "/ReadVmsHealth"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(readvmshealthrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_ReadVmsHealthResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &ReadVmsHealthResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_ReadVmsState(
	readvmsstaterequest ReadVmsStateRequest,
) (
	response *POST_ReadVmsStateResponses,
	err error,
) {
	path := client.service + "/ReadVmsState"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(readvmsstaterequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_ReadVmsStateResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &ReadVmsStateResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_ReadVolumes(
	readvolumesrequest ReadVolumesRequest,
) (
	response *POST_ReadVolumesResponses,
	err error,
) {
	path := client.service + "/ReadVolumes"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(readvolumesrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_ReadVolumesResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &ReadVolumesResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_ReadVpnConnections(
	readvpnconnectionsrequest ReadVpnConnectionsRequest,
) (
	response *POST_ReadVpnConnectionsResponses,
	err error,
) {
	path := client.service + "/ReadVpnConnections"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(readvpnconnectionsrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_ReadVpnConnectionsResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &ReadVpnConnectionsResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_ReadVpnGateways(
	readvpngatewaysrequest ReadVpnGatewaysRequest,
) (
	response *POST_ReadVpnGatewaysResponses,
	err error,
) {
	path := client.service + "/ReadVpnGateways"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(readvpngatewaysrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_ReadVpnGatewaysResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &ReadVpnGatewaysResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_RebootVms(
	rebootvmsrequest RebootVmsRequest,
) (
	response *POST_RebootVmsResponses,
	err error,
) {
	path := client.service + "/RebootVms"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(rebootvmsrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_RebootVmsResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &RebootVmsResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_RegisterImage(
	registerimagerequest RegisterImageRequest,
) (
	response *POST_RegisterImageResponses,
	err error,
) {
	path := client.service + "/RegisterImage"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(registerimagerequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_RegisterImageResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &RegisterImageResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_RegisterUserInGroup(
	registeruseringrouprequest RegisterUserInGroupRequest,
) (
	response *POST_RegisterUserInGroupResponses,
	err error,
) {
	path := client.service + "/RegisterUserInGroup"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(registeruseringrouprequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_RegisterUserInGroupResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &RegisterUserInGroupResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_RegisterVmsInListenerRule(
	registervmsinlistenerrulerequest RegisterVmsInListenerRuleRequest,
) (
	response *POST_RegisterVmsInListenerRuleResponses,
	err error,
) {
	path := client.service + "/RegisterVmsInListenerRule"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(registervmsinlistenerrulerequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_RegisterVmsInListenerRuleResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &RegisterVmsInListenerRuleResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_RegisterVmsInLoadBalancer(
	registervmsinloadbalancerrequest RegisterVmsInLoadBalancerRequest,
) (
	response *POST_RegisterVmsInLoadBalancerResponses,
	err error,
) {
	path := client.service + "/RegisterVmsInLoadBalancer"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(registervmsinloadbalancerrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_RegisterVmsInLoadBalancerResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &RegisterVmsInLoadBalancerResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_RejectNetPeering(
	rejectnetpeeringrequest RejectNetPeeringRequest,
) (
	response *POST_RejectNetPeeringResponses,
	err error,
) {
	path := client.service + "/RejectNetPeering"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(rejectnetpeeringrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_RejectNetPeeringResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &RejectNetPeeringResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_ResetAccountPassword(
	resetaccountpasswordrequest ResetAccountPasswordRequest,
) (
	response *POST_ResetAccountPasswordResponses,
	err error,
) {
	path := client.service + "/ResetAccountPassword"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(resetaccountpasswordrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_ResetAccountPasswordResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &ResetAccountPasswordResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_SendResetPasswordEmail(
	sendresetpasswordemailrequest SendResetPasswordEmailRequest,
) (
	response *POST_SendResetPasswordEmailResponses,
	err error,
) {
	path := client.service + "/SendResetPasswordEmail"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(sendresetpasswordemailrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_SendResetPasswordEmailResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &SendResetPasswordEmailResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_StartVms(
	startvmsrequest StartVmsRequest,
) (
	response *POST_StartVmsResponses,
	err error,
) {
	path := client.service + "/StartVms"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(startvmsrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_StartVmsResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &StartVmsResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_StopVms(
	stopvmsrequest StopVmsRequest,
) (
	response *POST_StopVmsResponses,
	err error,
) {
	path := client.service + "/StopVms"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(stopvmsrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_StopVmsResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &StopVmsResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_UnlinkNetInternetGateway(
	unlinknetinternetgatewayrequest UnlinkNetInternetGatewayRequest,
) (
	response *POST_UnlinkNetInternetGatewayResponses,
	err error,
) {
	path := client.service + "/UnlinkNetInternetGateway"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(unlinknetinternetgatewayrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_UnlinkNetInternetGatewayResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &UnlinkNetInternetGatewayResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_UnlinkNic(
	unlinknicrequest UnlinkNicRequest,
) (
	response *POST_UnlinkNicResponses,
	err error,
) {
	path := client.service + "/UnlinkNic"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(unlinknicrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_UnlinkNicResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &UnlinkNicResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_UnlinkPolicy(
	unlinkpolicyrequest UnlinkPolicyRequest,
) (
	response *POST_UnlinkPolicyResponses,
	err error,
) {
	path := client.service + "/UnlinkPolicy"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(unlinkpolicyrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_UnlinkPolicyResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &UnlinkPolicyResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_UnlinkPrivateIps(
	unlinkprivateipsrequest UnlinkPrivateIpsRequest,
) (
	response *POST_UnlinkPrivateIpsResponses,
	err error,
) {
	path := client.service + "/UnlinkPrivateIps"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(unlinkprivateipsrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_UnlinkPrivateIpsResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &UnlinkPrivateIpsResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_UnlinkPublicIp(
	unlinkpubliciprequest UnlinkPublicIpRequest,
) (
	response *POST_UnlinkPublicIpResponses,
	err error,
) {
	path := client.service + "/UnlinkPublicIp"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(unlinkpubliciprequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_UnlinkPublicIpResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &UnlinkPublicIpResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_UnlinkRouteTable(
	unlinkroutetablerequest UnlinkRouteTableRequest,
) (
	response *POST_UnlinkRouteTableResponses,
	err error,
) {
	path := client.service + "/UnlinkRouteTable"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(unlinkroutetablerequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_UnlinkRouteTableResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &UnlinkRouteTableResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_UnlinkVolume(
	unlinkvolumerequest UnlinkVolumeRequest,
) (
	response *POST_UnlinkVolumeResponses,
	err error,
) {
	path := client.service + "/UnlinkVolume"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(unlinkvolumerequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_UnlinkVolumeResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &UnlinkVolumeResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_UnlinkVpnGateway(
	unlinkvpngatewayrequest UnlinkVpnGatewayRequest,
) (
	response *POST_UnlinkVpnGatewayResponses,
	err error,
) {
	path := client.service + "/UnlinkVpnGateway"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(unlinkvpngatewayrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_UnlinkVpnGatewayResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &UnlinkVpnGatewayResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_UpdateAccount(
	updateaccountrequest UpdateAccountRequest,
) (
	response *POST_UpdateAccountResponses,
	err error,
) {
	path := client.service + "/UpdateAccount"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(updateaccountrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_UpdateAccountResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &UpdateAccountResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_UpdateApiKey(
	updateapikeyrequest UpdateApiKeyRequest,
) (
	response *POST_UpdateApiKeyResponses,
	err error,
) {
	path := client.service + "/UpdateApiKey"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(updateapikeyrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_UpdateApiKeyResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &UpdateApiKeyResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_UpdateGroup(
	updategrouprequest UpdateGroupRequest,
) (
	response *POST_UpdateGroupResponses,
	err error,
) {
	path := client.service + "/UpdateGroup"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(updategrouprequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_UpdateGroupResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &UpdateGroupResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_UpdateHealthCheck(
	updatehealthcheckrequest UpdateHealthCheckRequest,
) (
	response *POST_UpdateHealthCheckResponses,
	err error,
) {
	path := client.service + "/UpdateHealthCheck"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(updatehealthcheckrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_UpdateHealthCheckResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &UpdateHealthCheckResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_UpdateImageAttribute(
	updateimageattributerequest UpdateImageAttributeRequest,
) (
	response *POST_UpdateImageAttributeResponses,
	err error,
) {
	path := client.service + "/UpdateImageAttribute"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(updateimageattributerequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_UpdateImageAttributeResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &UpdateImageAttributeResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_UpdateKeypair(
	updatekeypairrequest UpdateKeypairRequest,
) (
	response *POST_UpdateKeypairResponses,
	err error,
) {
	path := client.service + "/UpdateKeypair"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(updatekeypairrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_UpdateKeypairResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &UpdateKeypairResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_UpdateListenerRule(
	updatelistenerrulerequest UpdateListenerRuleRequest,
) (
	response *POST_UpdateListenerRuleResponses,
	err error,
) {
	path := client.service + "/UpdateListenerRule"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(updatelistenerrulerequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_UpdateListenerRuleResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &UpdateListenerRuleResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_UpdateLoadBalancerAttributes(
	updateloadbalancerattributesrequest UpdateLoadBalancerAttributesRequest,
) (
	response *POST_UpdateLoadBalancerAttributesResponses,
	err error,
) {
	path := client.service + "/UpdateLoadBalancerAttributes"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(updateloadbalancerattributesrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_UpdateLoadBalancerAttributesResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &UpdateLoadBalancerAttributesResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_UpdateLoadBalancerPolicies(
	updateloadbalancerpoliciesrequest UpdateLoadBalancerPoliciesRequest,
) (
	response *POST_UpdateLoadBalancerPoliciesResponses,
	err error,
) {
	path := client.service + "/UpdateLoadBalancerPolicies"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(updateloadbalancerpoliciesrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_UpdateLoadBalancerPoliciesResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &UpdateLoadBalancerPoliciesResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_UpdateNetAccess(
	updatenetaccessrequest UpdateNetAccessRequest,
) (
	response *POST_UpdateNetAccessResponses,
	err error,
) {
	path := client.service + "/UpdateNetAccess"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(updatenetaccessrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_UpdateNetAccessResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &UpdateNetAccessResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_UpdateNetOptions(
	updatenetoptionsrequest UpdateNetOptionsRequest,
) (
	response *POST_UpdateNetOptionsResponses,
	err error,
) {
	path := client.service + "/UpdateNetOptions"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(updatenetoptionsrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_UpdateNetOptionsResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &UpdateNetOptionsResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_UpdateNicAttribute(
	updatenicattributerequest UpdateNicAttributeRequest,
) (
	response *POST_UpdateNicAttributeResponses,
	err error,
) {
	path := client.service + "/UpdateNicAttribute"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(updatenicattributerequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_UpdateNicAttributeResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &UpdateNicAttributeResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_UpdateRoute(
	updaterouterequest UpdateRouteRequest,
) (
	response *POST_UpdateRouteResponses,
	err error,
) {
	path := client.service + "/UpdateRoute"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(updaterouterequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_UpdateRouteResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &UpdateRouteResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_UpdateRoutePropagation(
	updateroutepropagationrequest UpdateRoutePropagationRequest,
) (
	response *POST_UpdateRoutePropagationResponses,
	err error,
) {
	path := client.service + "/UpdateRoutePropagation"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(updateroutepropagationrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_UpdateRoutePropagationResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &UpdateRoutePropagationResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_UpdateRouteTableLink(
	updateroutetablelinkrequest UpdateRouteTableLinkRequest,
) (
	response *POST_UpdateRouteTableLinkResponses,
	err error,
) {
	path := client.service + "/UpdateRouteTableLink"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(updateroutetablelinkrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_UpdateRouteTableLinkResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &UpdateRouteTableLinkResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_UpdateServerCertificate(
	updateservercertificaterequest UpdateServerCertificateRequest,
) (
	response *POST_UpdateServerCertificateResponses,
	err error,
) {
	path := client.service + "/UpdateServerCertificate"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(updateservercertificaterequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_UpdateServerCertificateResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &UpdateServerCertificateResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_UpdateSnapshotAttribute(
	updatesnapshotattributerequest UpdateSnapshotAttributeRequest,
) (
	response *POST_UpdateSnapshotAttributeResponses,
	err error,
) {
	path := client.service + "/UpdateSnapshotAttribute"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(updatesnapshotattributerequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_UpdateSnapshotAttributeResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &UpdateSnapshotAttributeResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_UpdateUser(
	updateuserrequest UpdateUserRequest,
) (
	response *POST_UpdateUserResponses,
	err error,
) {
	path := client.service + "/UpdateUser"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(updateuserrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_UpdateUserResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &UpdateUserResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}

//
func (client *Client) POST_UpdateVmAttribute(
	updatevmattributerequest UpdateVmAttributeRequest,
) (
	response *POST_UpdateVmAttributeResponses,
	err error,
) {
	path := client.service + "/UpdateVmAttribute"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(updatevmattributerequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_UpdateVmAttributeResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &UpdateVmAttributeResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		response.OK = result
	default:
		break
	}
	return
}
