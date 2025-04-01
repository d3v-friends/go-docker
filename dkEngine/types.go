package dkEngine

import (
	"encoding/base64"
	"encoding/json"
	"github.com/d3v-friends/go-tools/fnPointer"
	"strings"
	"time"
)

const (
	ErrAlreadyHasSameContainerName = "already_has_same_container_name"
)

const (
	xRegistryAuthHeader            = "X-Registry-Auth"
	httpHeaderKeyContentType       = "Content-Type"
	httpHeaderValueApplicationJson = "application/json"
)

type Platform string

const (
	PlatformLinuxAmd64 Platform = "linux/amd64"
)

func (x Platform) String() string {
	return string(x)
}

type Network struct {
	Name       string         `json:"ContainerName"`
	Id         string         `json:"Tag"`
	Created    time.Time      `json:"Created"`
	Scope      string         `json:"Scope"`
	Driver     string         `json:"Driver"`
	EnableIPv6 bool           `json:"EnableIPv6"`
	IPAM       map[string]any `json:"IPAM"`
	Internal   bool           `json:"Internal"`
	Attachable bool           `json:"Attachable"`
	ConfigForm map[string]any `json:"ConfigForm"`
	ConfigOnly bool           `json:"ConfigOnly"`
	Containers map[string]any `json:"Containers"`
	Options    map[string]any `json:"Options"`
	Labels     map[string]any `json:"Labels"`
}

type Networks []*Network

type Container struct {
	Id         string            `json:"Id"`
	Names      Names             `json:"Names"`
	Image      string            `json:"Image"`
	ImageID    string            `json:"ImageID"`
	Command    string            `json:"Command"`
	Created    int64             `json:"Created"`
	Ports      []*ContainerPort  `json:"Ports"`
	Labels     map[string]string `json:"Labels"`
	State      string            `json:"State"`
	HostConfig map[string]string `json:"HostConfig"`
}

type Names []string

func (x Names) Has(name string) bool {
	for _, str := range x {
		if strings.TrimPrefix(str, "/") == name {
			return true
		}
	}
	return false
}

type ContainerPort struct {
	IP          string `json:"IP"`
	PrivatePort int64  `json:"PrivatePort"`
	PublicPort  int64  `json:"PublicPort"`
	Type        string `json:"Type"`
}

type Containers []*Container

// CreateContainerRequest
// https://docs.docker.com/reference/api/engine/version/v1.47/#tag/Container/operation/ContainerCreate
type CreateContainerRequest struct {
	Cmd              []string          `json:"Cmd,omitempty"`
	Hostname         *string           `json:"ApiHost,omitempty"`
	User             *string           `json:"User,omitempty"`
	Env              []string          `json:"Env,omitempty"`
	Image            *string           `json:"Image,omitempty"`
	Volumes          map[string]string `json:"Volumes,omitempty"`
	HostConfig       *HostConfig       `json:"HostConfig,omitempty"`
	ExposedPorts     ExposedPorts      `json:"ExposedPorts,omitempty"`
	NetworkingConfig *NetworkingConfig `json:"NetworkingConfig,omitempty"`
}

type CreateContainerResponse struct {
	Id string `json:"Id"`
}

type NetworkingConfig struct {
	EndpointsConfig EndpointsConfig
}

type EndpointsConfig map[string]*EndpointSettings

type EndpointSettings struct {
	IPAMConfig          *EndpointIPAMConfig
	Links               []string
	Aliases             []string
	MacAddress          string
	DriverOpts          map[string]string
	NetworkID           string
	EndpointID          string
	Gateway             string
	IPAddress           string
	IPPrefixLen         int
	IPv6Gateway         string
	GlobalIPv6Address   string
	GlobalIPv6PrefixLen int
	DNSNames            []string
}

type EndpointIPAMConfig struct {
	IPv4Address  string   `json:",omitempty"`
	IPv6Address  string   `json:",omitempty"`
	LinkLocalIPs []string `json:",omitempty"`
}

type ExposedPorts map[string]map[string]string

type HostConfig struct {
	LogConfig    *LogConfig   `json:"LogConfig,omitempty"`
	PortBindings PortBindings `json:"PortBindings,omitempty"`
	NetworkMode  *string      `json:"NetworkMode,omitempty"`
	Binds        []string     `json:"Binds,omitempty"`
	Privileged   *bool        `json:"Privileged,omitempty"`
}

type LogConfig struct {
	Type   *string           `json:"Type,omitempty"`
	Config map[string]string `json:"Config,omitempty"`
}

type PortBinding struct {
	HostIp   *string `json:"HostIp,omitempty"`
	HostPort *string `json:"HostPort,omitempty"`
}

type PortBindings map[string][]*PortBinding

// ContainerInspection
// 필요한 정보만 추출함
// 추후 다른 내용이 필요시 다음 문서에서 찾아보기
// https://docker-docs.uclv.cu/engine/api/v1.40/#operation/ContainerInspect
type ContainerInspection struct {
	Id             string                     `json:"Id"`
	Created        *string                    `json:"Created"`
	Path           string                     `json:"Path"`
	Args           []string                   `json:"Args"`
	State          *ContainerInspectionState  `json:"State"`
	Image          string                     `json:"Image"`
	ResolvConfPath string                     `json:"ResolvConfPath"`
	HostnamePath   string                     `json:"HostnamePath"`
	HostsPath      string                     `json:"HostsPath"`
	LogPath        string                     `json:"LogPath"`
	Name           string                     `json:"Name"`
	RestartCount   int                        `json:"RestartCount"`
	Driver         string                     `json:"Driver"`
	Platform       string                     `json:"Platform"`
	MountLabel     string                     `json:"MountLabel"`
	ProcessLabel   string                     `json:"ProcessLabel"`
	Config         *ContainerInspectionConfig `json:"Config"`
}

type ContainerInspectionConfig struct {
	Hostname     string   `json:"Hostname"`
	Domainname   string   `json:"Domainname"`
	User         string   `json:"User"`
	AttachStdin  bool     `json:"AttachStdin"`
	AttachStdout bool     `json:"AttachStdout"`
	AttachStderr bool     `json:"AttachStderr"`
	Tty          bool     `json:"Tty"`
	OpenStdin    bool     `json:"OpenStdin"`
	StdinOnce    bool     `json:"StdinOnce"`
	Env          []string `json:"Env"`
	Cmd          []string `json:"Cmd"`
}

type ContainerInspectionState struct {
	Status     string                          `json:"Status"`
	Running    bool                            `json:"Running"`
	Paused     bool                            `json:"Paused"`
	Restarting bool                            `json:"Restarting"`
	OOMKilled  bool                            `json:"OOMKilled"`
	Dead       bool                            `json:"Dead"`
	Pid        int                             `json:"Pid"`
	ExitCode   int                             `json:"ExitCode"`
	Error      string                          `json:"Error"`
	StartedAt  string                          `json:"StartedAt"`
	FinishedAt string                          `json:"FinishedAt"`
	Health     *ContainerInspectionStateHealth `json:"Health"`
}

type ContainerInspectionStateHealth struct {
	Status        string `json:"Status"`
	FailingStreak int    `json:"FailingStreak"`
	Log           []struct {
		Start    string `json:"Start"`
		End      string `json:"End"`
		Output   string `json:"Output"`
		ExitCode int    `json:"ExitCode"`
	} `json:"Log"`
}

func (x *ContainerInspection) Env() (m map[string]string) {
	m = make(map[string]string)
	if fnPointer.IsNil(x.Config) {
		return
	}

	for _, s := range x.Config.Env {
		var ls = strings.Split(s, "=")
		if len(ls) != 2 {
			continue
		}
		m[ls[0]] = ls[1]
	}

	return
}

/* ------------------------------------------------------------------------------------------------------------ */

type VolumeOption string

func (x VolumeOption) String() string {
	return strings.ToLower(string(x))
}

func (x VolumeOption) IsValid() bool {
	for _, option := range VolumeOptionAll {
		if option == x {
			return true
		}
	}
	return false
}

const (
	VolumeOptionReadonly  VolumeOption = "RO"
	VolumeOptionReadWrite VolumeOption = "RW"
	VolumeOptionNocopy    VolumeOption = "NOCOPY"
	VolumeOptionZ         VolumeOption = "Z"
	VolumeOptionShared    VolumeOption = "SHARED"
	VolumeOptionSlave     VolumeOption = "SLAVE"
	VolumeOptionPrivate   VolumeOption = "PRIVATE"
)

var VolumeOptionAll = []VolumeOption{
	VolumeOptionReadonly,
	VolumeOptionReadWrite,
	VolumeOptionNocopy,
	VolumeOptionZ,
	VolumeOptionShared,
	VolumeOptionSlave,
	VolumeOptionPrivate,
}

/* ------------------------------------------------------------------------------------------------------------ */

type Registry interface {
	GetServerAddress() string
	GetUsername() string
	GetPassword() string
	GetEmail() string
}

func createRegistryToken(args Registry) (token string, err error) {
	var body []byte
	if body, err = json.Marshal(map[string]string{
		"serveraddress": args.GetServerAddress(),
		"username":      args.GetUsername(),
		"password":      args.GetPassword(),
		"email":         args.GetEmail(),
	}); err != nil {
		return
	}
	token = base64.StdEncoding.EncodeToString(body)
	return
}
