package log

import (
	"io/ioutil"
	"strings"

	"go.uber.org/zap"
)

func NewProdOrDev() (log *zap.Logger, err error) {
	if runningInDocker() {
		return zap.NewProduction()
	} else {
		return zap.NewDevelopment()
	}
}

// Will this work with CRI-O or other container runtimes?
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
