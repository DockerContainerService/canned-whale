package tcp

import (
	"fmt"
	"net"
)

func IsPortAvailable(port int) bool {
	address := fmt.Sprintf("%s:%d", "0.0.0.0", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return false
	}

	defer listener.Close()
	return true
}
