package profiler

import (
	"time"

	units "github.com/docker/go-units"
	log "github.com/sirupsen/logrus"
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
	if len(p.runtimes) > 0 {
		log.Info("Profiler results:")
		for _, runtime := range p.runtimes {
			log.Infof("Sector %s: %s elapsed", runtime.sectorName, units.HumanDuration(runtime.duration))
		}
	}
}
