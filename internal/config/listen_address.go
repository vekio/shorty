package config

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

func ValidateListenAddress(address string) error {
	if address == "" || strings.TrimSpace(address) != address {
		return fmt.Errorf("address must be a non-empty TCP host:port")
	}
	_, port, err := net.SplitHostPort(address)
	if err != nil {
		return fmt.Errorf("address must be a TCP host:port: %w", err)
	}
	portNumber, err := strconv.Atoi(port)
	if err != nil || portNumber < 1 || portNumber > 65535 {
		return fmt.Errorf("address port must be between 1 and 65535")
	}
	return nil
}
