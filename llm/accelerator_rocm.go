//go:build rocm

package llm

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	acceleratorName = "rocm"
)

var errNoGPU = errors.New("rocm-smi command failed")

// CheckVRAM returns the available VRAM in MiB on Linux machines with AMD GPUs
func CheckVRAM() (int64, error) {
	rocmHome := os.Getenv("ROCM_PATH")
	if rocmHome == "" {
		rocmHome = os.Getenv("ROCM_HOME")
	}
	if rocmHome == "" {
		log.Println("warning: ROCM_PATH is not set. Trying a likely fallback path, but it is recommended to set this variable in the environment.")
		rocmHome = "/opt/rocm"
	}
	cmd := exec.Command(filepath.Join(rocmHome, "bin/rocm-smi"), "--showmeminfo", "VRAM", "--csv")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err := cmd.Run()
	if err != nil {
		return 0, errNoGPU
	}
	csvData := csv.NewReader(&stdout)
	// llama.cpp or ROCm don't seem to understand splitting the VRAM allocations across them properly, so try to find the biggest card instead :(. FIXME.
	totalBiggestCard := int64(0)
	bigCardName := ""
	for {
		record, err := csvData.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, fmt.Errorf("failed to parse available VRAM: %v", err)
		}
		if !strings.HasPrefix(record[0], "card") {
			continue
		}
		cardTotal, err := strconv.ParseInt(record[1], 10, 64)
		if err != nil {
			return 0, err
		}
		cardUsed, err := strconv.ParseInt(record[2], 10, 64)
		if err != nil {
			return 0, err
		}
		possible := (cardTotal - cardUsed)
		if possible > totalBiggestCard {
			totalBiggestCard = possible
			bigCardName = record[0]
		}
	}
	log.Printf("ROCm presenting %d bytes of available VRAM on device %q", totalBiggestCard, bigCardName)
	return totalBiggestCard, nil
}
