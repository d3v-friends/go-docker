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

type ExecRequest struct {
	AttachStdin  *bool    `json:"AttachStdin,omitempty"`
	AttachStdout *bool    `json:"AttachStdout,omitempty"`
	AttachStderr *bool    `json:"AttachStderr,omitempty"`
	DetachKeys   *string  `json:"DetachKeys,omitempty"`
	Tty          *bool    `json:"Tty,omitempty"`
	Cmd          []string `json:"Cmd,omitempty"`
	Env          []string `json:"Env,omitempty"`
	Privileged   *bool    `json:"Privileged,omitempty"`
	User         *string  `json:"User,omitempty"`
	WorkingDir   *string  `json:"WorkingDir,omitempty"`
}

type ExecResponse struct {
	Id      string `json:"Id,omitempty"`
	Message string `json:"Message,omitempty"`
}

func Exec(
	host string,
	id string,
	args *ExecRequest,
) (res *ExecResponse, err error) {
	var body []byte
	if body, err = json.Marshal(args); err != nil {
		return
	}

	var request *http.Request
	if request, err = http.NewRequest(
		http.MethodPost,
		fmt.Sprintf(
			"%s/containers/%s/exec",
			host,
			id,
		),
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
	case 200, 201:
		var result = &ExecResponse{}
		if err = json.NewDecoder(resp.Body).Decode(result); err != nil {
			return
		}
		res = result
		return
	default:
		err = fnError.NewF("%s", fnPanic.Value(io.ReadAll(resp.Body)))
		return
	}
}
