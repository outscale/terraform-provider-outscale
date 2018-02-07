package handler

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// UnmarshalXML unmarshals a response body for the XML protocol.
func UnmarshalXML(v interface{}, r *http.Response) error {
	defer r.Body.Close()
	decoder, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("Read body: %v", err)
	}
	var data interface{}

	if err := xml.Unmarshal([]byte(decoder), &data); err != nil {
		return errors.New("SerializationError" + "failed decoding EC2 Query response" + fmt.Sprint(err))
	}
	return nil
}
