package handler

import (
	"fmt"
	"net/http"
)

const mediaTypeURLEncoded = "application/x-www-form-urlencoded"
const mediaType = "application/x-amz-json-1.1"

// SetHeaders sets the headers for the request
func SetHeaders(agent string, req *http.Request, operation string) {
	req.Header.Add("X-Amz-Target", fmt.Sprintf("%s.%s", agent, operation))
	commonHeadres(agent, mediaTypeURLEncoded, req)
}

// SetHeadersICU sets the headers for the request
func SetHeadersICU(agent string, req *http.Request, operation string) {
	req.Header.Add("X-Amz-Target", fmt.Sprintf("TinaIcuService.%s", operation))
	commonHeadres(agent, mediaType, req)
}

// SetHeadersDL sets the headers for the request
func SetHeadersDL(agent string, req *http.Request, operation string) {
	req.Header.Add("X-Amz-Target", fmt.Sprintf("OvertureService.%s", operation))
	commonHeadres(agent, mediaType, req)
}

func commonHeadres(agent, media string, req *http.Request) {
	req.Header.Add("User-Agent", agent)
	req.Header.Add("Content-Type", media)
}
