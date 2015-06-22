package filosofos

import (
	"github.com/jaodsilv/ep3/monitor"

	"math/rand"
	"sync"
	"time"
)

const (
	// Existe um erro na espera do processo, esse erro é da ordem de 0.1 milisegundos
	thinkingConstant = 100 * time.Microsecond
)

// Filosofo representa um filosofo.
type Filosofo struct {
	Index, Comidas, Peso int32
}

func newFilosofo(index, peso int32) *Filosofo {
	return &Filosofo{
		Index:   index,
		Peso:    peso,
		Comidas: int32(0),
	}
}

func pensa() {
	time.Sleep(thinkingConstant * time.Duration(rand.Int31n(100)+1))
}

func (f *Filosofo) janta(m *monitor.FilosofoMonitor, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		// A alternância entre pares e impares garante que não haverá deadlock nem starvation
		switch m.Come(f.Comidas, f.Peso, f.Index) {
		case monitor.Comeu:
			f.Comidas++
		case monitor.Cheio:
		case monitor.AcabouComida:
			return
		}
		pensa()
	}
}
