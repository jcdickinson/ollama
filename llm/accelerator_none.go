//go:build !rocm && !cuda

package llm

import (
	"errors"
)

const (
	acceleratorName = "none"
)

var errNoGPU = errors.New("no accelerator found")

// CheckVRAM is a stub with no accelerator.
func CheckVRAM() (int64, error) {
	return 0, errNoGPU
}
