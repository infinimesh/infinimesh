package log

import (
	"io/ioutil"
	"strings"

	"go.uber.org/zap"
)

// NewProdOrDev returns a Prod logger if the application is running in Docker,
// otherwise it returns a Dev logger.
func NewProdOrDev() (log *zap.Logger, err error) {
	if runningInDocker() {
		return zap.NewProduction()
	} else {
		return zap.NewDevelopment()
	}
}

// TODO Will this work with CRI-O, rkt and other container runtimes beside docker?
func runningInDocker() bool {
	f, err := ioutil.ReadFile("/proc/1/cgroup")
	if err != nil {
		return false
	}

	if strings.Contains(string(f), "docker") {
		return true
	}
	return false
}
