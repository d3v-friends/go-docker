package dkEngine

import (
	"fmt"
	"github.com/d3v-friends/go-tools/fnError"
	"github.com/d3v-friends/go-tools/fnPanic"
	"io"
	"net/http"
	"time"
)

func Ping(
	host string,
) (err error) {
	var request *http.Request
	if request, err = http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/_ping", host),
		nil,
	); err != nil {
		return
	}

	var c = &http.Client{
		Timeout: time.Second * 5,
	}

	var resp *http.Response
	if resp, err = c.Do(request); err != nil {
		return
	}

	switch resp.StatusCode {
	case 200:
		return
	default:
		err = fnError.NewF("%s", fnPanic.Value(io.ReadAll(resp.Body)))
		return
	}
}
