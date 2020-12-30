package util

import (
	"github.com/pkg/errors"
	"os"
	"strings"
)

func GetHostName(hostnameOverride string) (string, error) {
	hostName := hostnameOverride
	if len(hostName) == 0 {
		nodeName, err := os.Hostname()
		if err != nil {
			return "", errors.Wrapf(err, "couldn't determine hostname")
		}
		hostName = nodeName
	}

	hostName = strings.TrimSpace(hostName)
	if len(hostName) == 0 {
		return "", errors.New("empty hostname is invalid")
	}

	return strings.ToLower(hostName), nil
}
