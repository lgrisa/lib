package utils

import (
	"fmt"
	"os"
	"os/user"
)

var value string

func init() {
	value = fmt.Sprintf("%s@%s", getUsername(), getHostname())
}

func getUsername() string {
	username := os.Getenv("USER")

	if u, err := user.Current(); err == nil {
		if username != u.Username {
			username = fmt.Sprintf("%s(%s)", username, u.Username)
			if u.Username != u.Name {
				username = fmt.Sprintf("%s(%s)", username, u.Name)
			}
		} else {
			if username != u.Name {
				username = fmt.Sprintf("%s(%s)", username, u.Name)
			}
		}
	}

	return username
}

func getHostname() string {
	hostname := os.Getenv("HOSTNAME")
	if hostname == "" {
		hostname = os.Getenv("HOST")
	}
	if hostname == "" {
		hostname, _ = os.Hostname()
	}
	if hostname == "" {
		hostname = "localhost"
	}
	return hostname
}

func GetValue() string {
	return value
}
