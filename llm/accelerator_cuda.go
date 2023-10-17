//go:build cuda

package llm

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"

	"github.com/jmorganca/ollama/format"
)

const (
	acceleratorName = "cuda"
)

var errNoGPU = errors.New("nvidia-smi command failed")

// CheckVRAM returns the free VRAM in bytes on Linux machines with NVIDIA GPUs
func CheckVRAM() (int64, error) {
	cmd := exec.Command("nvidia-smi", "--query-gpu=memory.free", "--format=csv,noheader,nounits")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err := cmd.Run()
	if err != nil {
		return 0, errNoGPU
	}

	var freeMiB int64
	scanner := bufio.NewScanner(&stdout)
	for scanner.Scan() {
		line := scanner.Text()
		vram, err := strconv.ParseInt(strings.TrimSpace(line), 10, 64)
		if err != nil {
			return 0, fmt.Errorf("failed to parse available VRAM: %v", err)
		}

		freeMiB += vram
	}

	freeBytes := freeMiB * 1024 * 1024
	if freeBytes < 2*format.GigaByte {
		log.Printf("less than 2 GB VRAM available, falling back to CPU only")
		freeMiB = 0
	}

	return freeBytes, nil
}
