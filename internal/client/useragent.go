package client

import (
	"fmt"
	"runtime"
)

var (
	Version = "dev"
	Commit  = "none"
)

func UserAgent() string {
	return fmt.Sprintf("ploicloud-cli/%s (%s; %s)", Version, runtime.GOOS, runtime.GOARCH)
}
