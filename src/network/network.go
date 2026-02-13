package network

type Tunnel interface {
	Start() error
	Read([]byte) (int, error)
	Write([]byte) (int, error)
	Close() error
}
