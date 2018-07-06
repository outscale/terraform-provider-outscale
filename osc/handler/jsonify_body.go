package handler

import (
	"encoding/json"
)

// Bind ...
func Bind(operation string, body interface{}) {}

// BindICU ...
func BindICU(operation string, body interface{}) string {
	v := struct {
		Action               string `json:"Action"`
		Version              string `json:"Version"`
		AuthenticationMethod string `json:"AuthenticationMethod"`
	}{operation, "2018-05-14", "accesskey"}

	return commonStuctre(v, body)
}

// BindDL ...
func BindDL(operation string, body interface{}) string {
	return commonStuctre(nil, body)
}

func commonStuctre(v, body interface{}) string {
	var m map[string]interface{}

	jb, _ := json.Marshal(body)
	json.Unmarshal(jb, &m)

	ja, _ := json.Marshal(v)
	json.Unmarshal(ja, &m)

	jm, _ := json.Marshal(m)

	return string(jm)
}
