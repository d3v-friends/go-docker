package dkRegistry

import (
	"encoding/json"
	"fmt"
	"github.com/d3v-friends/go-tools/fnError"
	"github.com/d3v-friends/go-tools/fnPanic"
	"io"
	"net/http"
	"time"
)

func Ping(
	args Registry,
) (err error) {
	var request *http.Request
	if request, err = http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("https://%s/v2", args.GetServerAddress()),
		nil,
	); err != nil {
		return
	}

	request.SetBasicAuth(args.GetUsername(), args.GetPassword())

	var response *http.Response
	if response, err = (&http.Client{
		Timeout: time.Second * 10,
	}).Do(request); err != nil {
		return
	}

	switch response.StatusCode {
	case 200:
		return
	default:
		err = fnError.NewF("%s", fnPanic.Value(io.ReadAll(response.Body)))
		return
	}
}

/* ------------------------------------------------------------------------------------------------------------ */

type RepositoryResponse struct {
	Repositories []string `json:"repositories"`
}

func QueryRepositories(
	args Registry,
) (res []string, err error) {
	var request *http.Request
	if request, err = http.NewRequest(
		http.MethodGet,
		fmt.Sprintf(
			"https://%s/v2/_catalog",
			args.GetServerAddress(),
		),
		nil,
	); err != nil {
		return
	}

	request.SetBasicAuth(args.GetUsername(), args.GetPassword())
	request.Header.Set(httpHeaderKeyContentType, httpHeaderValueApplicationJson)

	var response *http.Response
	if response, err = (&http.Client{
		Timeout: time.Second * 10,
	}).Do(request); err != nil {
		return
	}

	switch response.StatusCode {
	case 200:
		var result = &RepositoryResponse{}
		if err = json.
			NewDecoder(response.Body).
			Decode(result); err != nil {
			return
		}
		res = result.Repositories
		return
	default:
		err = fnError.NewF("%s", fnPanic.Value(io.ReadAll(response.Body)))
		return
	}
}

/* ------------------------------------------------------------------------------------------------------------ */

type QueryTagsResponse struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

func QueryTags(
	args Registry,
	repository string,
) (res []string, err error) {
	var request *http.Request
	if request, err = http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("https://%s/v2/%s/tags/list", args.GetServerAddress(), repository),
		nil,
	); err != nil {
		return
	}

	request.SetBasicAuth(args.GetUsername(), args.GetPassword())
	request.Header.Set(httpHeaderKeyContentType, httpHeaderValueApplicationJson)

	var response *http.Response
	if response, err = (&http.Client{
		Timeout: time.Second * 10,
	}).Do(request); err != nil {
		return
	}

	switch response.StatusCode {
	case 200:
		var result = &QueryTagsResponse{}
		if err = json.NewDecoder(response.Body).Decode(result); err != nil {
			return
		}

		res = result.Tags
		return
	default:
		err = fnError.NewF("%s", fnPanic.Value(io.ReadAll(response.Body)))
		return
	}
}
