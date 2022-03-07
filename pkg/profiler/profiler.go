package profiler

import (
	"fmt"
	"time"
)

type runtime struct {
	duration   time.Duration
	sectorName string
}

type Profiler struct {
	runtimes []runtime
}

func NewProfiler() *Profiler {
	return &Profiler{}
}

func (p *Profiler) AddRuntime(startTime time.Time, sectorName string) {
	sectorDuration := time.Since(startTime)
	p.runtimes = append(p.runtimes,
		runtime{
			duration:   sectorDuration,
			sectorName: sectorName})
}

func (p *Profiler) PrintRuntimes() {
	fmt.Println("Profiler results:")
	for _, runtime := range p.runtimes {
		fmt.Println("Sector", runtime.sectorName, "finished in", runtime.duration)
	}
}
