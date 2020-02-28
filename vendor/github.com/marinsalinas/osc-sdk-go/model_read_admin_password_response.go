/*
 * 3DS OUTSCALE API
 *
 * Welcome to the 3DS OUTSCALE's API documentation.<br /><br />  The 3DS OUTSCALE API enables you to manage your resources in the 3DS OUTSCALE Cloud. This documentation describes the different actions available along with code examples.<br /><br />  Note that the 3DS OUTSCALE Cloud is compatible with Amazon Web Services (AWS) APIs, but some resources have different names in AWS than in the 3DS OUTSCALE API. You can find a list of the differences [here](https://wiki.outscale.net/display/EN/3DS+OUTSCALE+APIs+Reference).<br /><br />  You can also manage your resources using the [Cockpit](https://wiki.outscale.net/display/EN/About+Cockpit) web interface.
 *
 * API version: 0.15
 * Contact: support@outscale.com
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package oscgo

import (
	"bytes"
	"encoding/json"
)

// ReadAdminPasswordResponse struct for ReadAdminPasswordResponse
type ReadAdminPasswordResponse struct {
	// The password of the VM. After the first boot, returns an empty string.
	AdminPassword   *string          `json:"AdminPassword,omitempty"`
	ResponseContext *ResponseContext `json:"ResponseContext,omitempty"`
	// The ID of the VM.
	VmId *string `json:"VmId,omitempty"`
}

// GetAdminPassword returns the AdminPassword field value if set, zero value otherwise.
func (o *ReadAdminPasswordResponse) GetAdminPassword() string {
	if o == nil || o.AdminPassword == nil {
		var ret string
		return ret
	}
	return *o.AdminPassword
}

// GetAdminPasswordOk returns a tuple with the AdminPassword field value if set, zero value otherwise
// and a boolean to check if the value has been set.
func (o *ReadAdminPasswordResponse) GetAdminPasswordOk() (string, bool) {
	if o == nil || o.AdminPassword == nil {
		var ret string
		return ret, false
	}
	return *o.AdminPassword, true
}

// HasAdminPassword returns a boolean if a field has been set.
func (o *ReadAdminPasswordResponse) HasAdminPassword() bool {
	if o != nil && o.AdminPassword != nil {
		return true
	}

	return false
}

// SetAdminPassword gets a reference to the given string and assigns it to the AdminPassword field.
func (o *ReadAdminPasswordResponse) SetAdminPassword(v string) {
	o.AdminPassword = &v
}

// GetResponseContext returns the ResponseContext field value if set, zero value otherwise.
func (o *ReadAdminPasswordResponse) GetResponseContext() ResponseContext {
	if o == nil || o.ResponseContext == nil {
		var ret ResponseContext
		return ret
	}
	return *o.ResponseContext
}

// GetResponseContextOk returns a tuple with the ResponseContext field value if set, zero value otherwise
// and a boolean to check if the value has been set.
func (o *ReadAdminPasswordResponse) GetResponseContextOk() (ResponseContext, bool) {
	if o == nil || o.ResponseContext == nil {
		var ret ResponseContext
		return ret, false
	}
	return *o.ResponseContext, true
}

// HasResponseContext returns a boolean if a field has been set.
func (o *ReadAdminPasswordResponse) HasResponseContext() bool {
	if o != nil && o.ResponseContext != nil {
		return true
	}

	return false
}

// SetResponseContext gets a reference to the given ResponseContext and assigns it to the ResponseContext field.
func (o *ReadAdminPasswordResponse) SetResponseContext(v ResponseContext) {
	o.ResponseContext = &v
}

// GetVmId returns the VmId field value if set, zero value otherwise.
func (o *ReadAdminPasswordResponse) GetVmId() string {
	if o == nil || o.VmId == nil {
		var ret string
		return ret
	}
	return *o.VmId
}

// GetVmIdOk returns a tuple with the VmId field value if set, zero value otherwise
// and a boolean to check if the value has been set.
func (o *ReadAdminPasswordResponse) GetVmIdOk() (string, bool) {
	if o == nil || o.VmId == nil {
		var ret string
		return ret, false
	}
	return *o.VmId, true
}

// HasVmId returns a boolean if a field has been set.
func (o *ReadAdminPasswordResponse) HasVmId() bool {
	if o != nil && o.VmId != nil {
		return true
	}

	return false
}

// SetVmId gets a reference to the given string and assigns it to the VmId field.
func (o *ReadAdminPasswordResponse) SetVmId(v string) {
	o.VmId = &v
}

type NullableReadAdminPasswordResponse struct {
	Value        ReadAdminPasswordResponse
	ExplicitNull bool
}

func (v NullableReadAdminPasswordResponse) MarshalJSON() ([]byte, error) {
	switch {
	case v.ExplicitNull:
		return []byte("null"), nil
	default:
		return json.Marshal(v.Value)
	}
}

func (v *NullableReadAdminPasswordResponse) UnmarshalJSON(src []byte) error {
	if bytes.Equal(src, []byte("null")) {
		v.ExplicitNull = true
		return nil
	}

	return json.Unmarshal(src, &v.Value)
}