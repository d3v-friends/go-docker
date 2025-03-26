package dkRegistry

const (
	httpHeaderKeyContentType       = "Content-Type"
	httpHeaderValueApplicationJson = "application/json"
)

type Registry interface {
	GetServerAddress() string
	GetUsername() string
	GetPassword() string
}
