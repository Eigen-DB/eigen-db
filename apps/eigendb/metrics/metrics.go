package metrics

import (
	"os"
	"time"

	"github.com/shirou/gopsutil/v4/process"
)

var programStartTime time.Time

// Initializes the metrics
func Init() {
	programStartTime = time.Now()
}

// Returns the program's uptime
func GetUptime() time.Duration {
	return time.Since(programStartTime)
}

// Returns memory usage in Megabytes
func GetMemUsage() (float32, error) {
	p, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		return -1.0, err
	}
	usage, err := p.MemoryPercent()
	if err != nil {
		return -1.0, err
	}
	return usage, nil
}

func GetCpuUsage() (float64, error) {
	p, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		return -1.0, err
	}
	usage, err := p.CPUPercent()
	if err != nil {
		return -1.0, err
	}
	return usage, nil
}
