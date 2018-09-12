// GENERATED FILE: DO NOT EDIT!

package oapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
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
	s := fmt.Sprintf("https://%s.%s.%s", c.Service, c.Region, c.URL)

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

// Sign ...
func (c *Client) Sign(req *http.Request, body []byte) error {
	reader := strings.NewReader(string(body))
	timestamp := time.Now()
	_, err := c.signer.Sign(req, reader, "oapi", c.config.Region, timestamp)
	utils.DebugRequest(req)
	return err

}

// Do ...
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		log.Printf("[Debug] Error in Do Request %s", err)
	}
	utils.DebugResponse(resp)
	return resp, err
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
func (client *Client) POST_ImportSnapshot(
	importsnapshotrequest ImportSnapshotRequest,
) (
	response *POST_ImportSnapshotResponses,
	err error,
) {
	path := client.service + "/ImportSnapshot"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(importsnapshotrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_ImportSnapshotResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &ImportSnapshotResponse{}
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
func (client *Client) POST_LinkPrivateIps(
	linkprivateipsrequest LinkPrivateIpsRequest,
) (
	response *POST_LinkPrivateIpsResponses,
	err error,
) {
	path := client.service + "/LinkPrivateIps"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(linkprivateipsrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_LinkPrivateIpsResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &LinkPrivateIpsResponse{}
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
func (client *Client) POST_UpdateImage(
	updateimagerequest UpdateImageRequest,
) (
	response *POST_UpdateImageResponses,
	err error,
) {
	path := client.service + "/UpdateImage"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(updateimagerequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_UpdateImageResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &UpdateImageResponse{}
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
	resp, err := client.Do(req)
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
func (client *Client) POST_UpdateNic(
	updatenicrequest UpdateNicRequest,
) (
	response *POST_UpdateNicResponses,
	err error,
) {
	path := client.service + "/UpdateNic"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(updatenicrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_UpdateNicResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &UpdateNicResponse{}
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
	resp, err := client.Do(req)
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
func (client *Client) POST_UpdateSnapshot(
	updatesnapshotrequest UpdateSnapshotRequest,
) (
	response *POST_UpdateSnapshotResponses,
	err error,
) {
	path := client.service + "/UpdateSnapshot"
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(updatesnapshotrequest)
	req, err := http.NewRequest("POST", path, body)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	req.Header = reqHeaders

	client.Sign(req, body.Bytes())

	if err != nil {
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	response = &POST_UpdateSnapshotResponses{}
	switch {
	case resp.StatusCode == 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := &UpdateSnapshotResponse{}
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
