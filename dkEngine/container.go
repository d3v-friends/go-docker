package dkEngine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/d3v-friends/go-tools/fnError"
	"github.com/d3v-friends/go-tools/fnPanic"
	"github.com/d3v-friends/go-tools/fnPointer"
	"io"
	"net/http"
	"time"
)

type CreateContainerArgs struct {
	args          *CreateContainerRequest
	platform      Platform
	containerName string
	networkName   string
}

func NewCreateContainerArgs(
	containerName string,
	networkName string,
	image string,
	platform Platform,
) *CreateContainerArgs {
	return &CreateContainerArgs{
		args: &CreateContainerRequest{
			Cmd:      nil,
			Hostname: nil,
			User:     nil,
			Env:      make([]string, 0),
			Image:    fnPointer.Make(image),
			Volumes:  map[string]string{},
			HostConfig: &HostConfig{
				LogConfig: &LogConfig{
					Type: fnPointer.Make("json-file"),
					Config: map[string]string{
						"max-size": "10m",
						"max-file": "100",
					},
				},
				PortBindings: PortBindings{},
				NetworkMode:  fnPointer.Make("bridge"),
				Binds:        []string{},
				Privileged:   fnPointer.Make(false),
			},
			ExposedPorts: ExposedPorts{},
			NetworkingConfig: &NetworkingConfig{
				EndpointsConfig: EndpointsConfig{
					networkName: {
						Aliases: []string{
							containerName,
						},
						DNSNames: []string{
							containerName,
						},
					},
				},
			},
		},
		platform:      platform,
		containerName: containerName,
		networkName:   networkName,
	}
}

func (x *CreateContainerArgs) Body() ([]byte, error) {
	return json.Marshal(x.args)
}

func (x *CreateContainerArgs) AppendVolumeBinds(
	host string,
	container string,
	options ...VolumeOption,
) {
	var args = fmt.Sprintf("%s:%s", host, container)
	if len(options) == 1 {
		args = fmt.Sprintf("%s:%s", args, options[0].String())
	}
	x.args.HostConfig.Binds = append(x.args.HostConfig.Binds, args)
}

func (x *CreateContainerArgs) AppendEnv(key, value string) {
	x.args.Env = append(x.args.Env, fmt.Sprintf("%s=%s", key, value))
}

func (x *CreateContainerArgs) AppendPortBind(
	host uint64,
	container uint64,
) {
	var hostPort = fnPointer.Make(fmt.Sprintf("%d", host))
	x.args.HostConfig.PortBindings[fmt.Sprintf("%d/tcp", container)] = []*PortBinding{
		{HostIp: fnPointer.Make("0.0.0.0"), HostPort: hostPort},
		{HostIp: fnPointer.Make("::"), HostPort: hostPort},
	}
	x.args.ExposedPorts[fmt.Sprintf("%d/tcp", host)] = map[string]string{}
}

func (x *CreateContainerArgs) SetCmd(cmd []string) {
	x.args.Cmd = cmd
}

func (x *CreateContainerArgs) SetPrivileged(v bool) {
	x.args.HostConfig.Privileged = fnPointer.Make(v)
}

func (x *CreateContainerArgs) SetNetworkMode(mode string) {
	x.args.HostConfig.NetworkMode = fnPointer.Make(mode)
}

func (x *CreateContainerArgs) SetLogConfig(config *LogConfig) {
	x.args.HostConfig.LogConfig = config
}

/* ------------------------------------------------------------------------------------------------------------ */

func CreateContainer(
	host string,
	args *CreateContainerArgs,
	registries ...Registry,
) (id string, err error) {
	var body []byte
	if body, err = args.Body(); err != nil {
		return
	}

	var request *http.Request
	if request, err = http.NewRequest(
		http.MethodPost,
		fmt.Sprintf(
			"%s/containers/create?name=%s&platform=%s",
			host,
			args.containerName,
			args.platform.String(),
		),
		bytes.NewReader(body),
	); err != nil {
		return
	}

	request.Header.Set(httpHeaderKeyContentType, httpHeaderValueApplicationJson)

	if len(registries) == 1 {
		var token string
		if token, err = createRegistryToken(registries[0]); err != nil {
			return
		}
		request.Header.Set(xRegistryAuthHeader, token)
	}

	var resp *http.Response
	if resp, err = (&http.Client{
		Timeout: time.Second * 60,
	}).Do(request); err != nil {
		return
	}

	switch resp.StatusCode {
	case 200, 201:
		var result = &CreateContainerResponse{}
		if err = json.NewDecoder(resp.Body).Decode(result); err != nil {
			return
		}

		id = result.Id
		return
	default:
		err = fnError.NewF("%s", fnPanic.Value(io.ReadAll(resp.Body)))
		return
	}
}

/* ------------------------------------------------------------------------------------------------------------ */

func QueryContainers(
	host string,
) (ls Containers, err error) {
	var request *http.Request
	if request, err = http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/containers/json?all=true", host),
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
	default:
		err = fnError.NewF("%s", fnPanic.Value(io.ReadAll(resp.Body)))
		return
	}

	ls = make(Containers, 0)
	if err = json.NewDecoder(resp.Body).Decode(&ls); err != nil {
		return
	}

	return
}

/* ------------------------------------------------------------------------------------------------------------ */

func Start(
	host string,
	id string,
) (err error) {
	var request *http.Request
	if request, err = http.NewRequest(
		http.MethodPost,
		fmt.Sprintf(
			"%s/containers/%s/start",
			host,
			id,
		),
		nil,
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
	case 200, 201, 204:
		return
	default:
		err = fnError.NewF("%s", fnPanic.Value(io.ReadAll(resp.Body)))
		return
	}
}

/* ------------------------------------------------------------------------------------------------------------ */

func Pull(
	host string,
	image string,
	registries ...Registry,
) (err error) {
	var request *http.Request
	if request, err = http.NewRequest(
		http.MethodPost,
		fmt.Sprintf(
			"%s/images/create?fromImage=%s",
			host,
			image,
		),
		nil,
	); err != nil {
		return
	}

	request.Header.Set(httpHeaderKeyContentType, httpHeaderValueApplicationJson)

	if len(registries) == 1 {
		var token string
		if token, err = createRegistryToken(registries[0]); err != nil {
			return
		}
		request.Header.Set(xRegistryAuthHeader, token)
	}

	var resp *http.Response
	if resp, err = (&http.Client{
		Timeout: time.Second * 60,
	}).Do(request); err != nil {
		return
	}

	switch resp.StatusCode {
	case 200, 201, 204:
		// 결과가 스트림으로 보내지므로 연결을 지속해야 작업이 이뤄진다
		if _, err = io.ReadAll(resp.Body); err != nil {
			return
		}
		return
	default:
		err = fnError.NewF("%s", fnPanic.Value(io.ReadAll(resp.Body)))
		return
	}
}

/* ------------------------------------------------------------------------------------------------------------ */

func Stop(
	host string,
	id string,
) (err error) {
	var request *http.Request
	if request, err = http.NewRequest(
		http.MethodPost,
		fmt.Sprintf(
			"%s/containers/%s/stop",
			host,
			id,
		),
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
	case 204, 304:
		return
	default:
		err = fnError.NewF("%s", fnPanic.Value(io.ReadAll(resp.Body)))
		return
	}
}

/* ------------------------------------------------------------------------------------------------------------ */

func Kill(
	host string,
	id string,
) (err error) {
	var request *http.Request
	if request, err = http.NewRequest(
		"POST",
		fmt.Sprintf(
			"%s/containers/%s/kill",
			host,
			id,
		),
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
	case 200, 201, 204:
		return
	default:
		err = fnError.NewF("%s", fnPanic.Value(io.ReadAll(resp.Body)))
		return
	}
}

/* ------------------------------------------------------------------------------------------------------------ */

func Remove(
	host string,
	id string,
) (err error) {
	var request *http.Request
	if request, err = http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf(
			"%s/containers/%s",
			host,
			id,
		),
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
	case 200, 201, 204:
	default:
		err = fnError.NewF("%s", fnPanic.Value(io.ReadAll(resp.Body)))
		return
	}

	return
}

/* ------------------------------------------------------------------------------------------------------------ */

func Inspect(
	host string,
	id string,
) (res *ContainerInspection, err error) {
	var request *http.Request
	if request, err = http.NewRequest(
		http.MethodGet,
		fmt.Sprintf(
			"%s/containers/%s/json",
			host,
			id,
		),
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
		res = &ContainerInspection{}
		if err = json.NewDecoder(resp.Body).Decode(res); err != nil {
			return
		}
		return
	default:
		err = fnError.NewF("%s", fnPanic.Value(io.ReadAll(resp.Body)))
		return
	}
}
