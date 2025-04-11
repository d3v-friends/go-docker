package dkEngine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/d3v-friends/go-tools/fnError"
	"github.com/d3v-friends/go-tools/fnPanic"
	"io"
	"net/http"
	"time"
)

func QueryNetworks(
	host string,
) (ls Networks, err error) {
	var request *http.Request
	if request, err = http.NewRequest(
		"GET",
		fmt.Sprintf("%s/networks", host),
		nil,
	); err != nil {
		return
	}

	var resp *http.Response
	if resp, err = (&http.Client{
		Timeout: time.Second * 10,
	}).Do(request); err != nil {
		return
	}

	switch resp.StatusCode {
	case 200:
		ls = make(Networks, 0)
		if err = json.NewDecoder(resp.Body).Decode(&ls); err != nil {
			return
		}
		return
	default:
		err = fnError.NewF("%s", fnPanic.Value(io.ReadAll(resp.Body)))
		return
	}
}

/* ------------------------------------------------------------------------------------------------------------ */

type CreateNetworkRequest struct {
	Name     string `json:"name"`
	Driver   string `json:"driver"`
	Internal bool   `json:"internal"`
}

type CreateNetworkResponse struct {
	Id      string `json:"Tag"`
	Warning string `json:"Warning"`
}

func CreateNetwork(
	host string,
	name string,
	driver string,
	internal bool,
) (res *CreateNetworkResponse, err error) {
	var body []byte
	if body, err = json.Marshal(&CreateNetworkRequest{
		Name:     name,
		Driver:   driver,
		Internal: internal,
	}); err != nil {
		return
	}

	var request *http.Request
	if request, err = http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/networks/create", host),
		bytes.NewReader(body),
	); err != nil {
		return
	}

	request.Header.Set(httpHeaderKeyContentType, httpHeaderValueApplicationJson)

	var resp *http.Response
	if resp, err = (&http.Client{
		Timeout: time.Second * 10,
	}).Do(request); err != nil {
		return
	}

	switch resp.StatusCode {
	case 201:
		res = &CreateNetworkResponse{}
		if err = json.NewDecoder(resp.Body).Decode(res); err != nil {
			return
		}
	default:
		err = fnError.NewF("%s", fnPanic.Value(io.ReadAll(resp.Body)))
		return
	}

	return
}

/* ------------------------------------------------------------------------------------------------------------ */

func DeleteNetwork(
	host string,
	networkName string,
) (err error) {
	var request *http.Request
	if request, err = http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("%s/networks/%s", host, networkName),
		nil,
	); err != nil {
		return
	}

	var resp *http.Response
	if resp, err = (&http.Client{
		Timeout: time.Second * 10,
	}).Do(request); err != nil {
		return
	}

	switch resp.StatusCode {
	case 204:
	default:
		err = fnError.NewF("%s", fnPanic.Value(io.ReadAll(resp.Body)))
		return
	}
	return
}
