package osc

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestError(t *testing.T) {
	if err := SendError(ErrMsg["HTTP"], errors.New("Testing Error Message")); err != nil {
		//t.Fatalf("err: %s", err)
		t.Log(err)
	}

}

func TestUnmarshallErrorHandler(t *testing.T) {
	test := &http.Response{
		Body: ioutil.NopCloser(bytes.NewBufferString("Hello World")),
	}

	buff := bytes.NewBuffer(nil)
	test.Write(buff)

	fmt.Println(buff)

	if err := UnmarshalErrorHandler(test); err != nil {
		//t.Fatalf("err: %s", err)
		t.Log(err)
	}
}
