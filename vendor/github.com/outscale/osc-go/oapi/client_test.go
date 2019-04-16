package oapi

import (
	"log"
	"testing"
)

func TestCheckErrorResponse(t *testing.T) {
	jsonError := `
	{
		"Errors":
		[{
			"Code":"6003",
			"Details":"",
			"Type":"InvalidState"
		}],
		"ResponseContext":{
			"RequestId":"b52ceadc-c8c4-4815-b037-446c82e1e074"
		}
	}
	`

	fmtErr, err := fmtErrorResponse([]byte(jsonError))

	if err != nil {
		t.Errorf("got error %s", err)
	}

	log.Printf("[Debug] %s", fmtErr)
}
