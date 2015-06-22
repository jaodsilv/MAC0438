// Package monitor é o monitor de filosofos
package monitor

// Este arquivo é a parte que contém a parte exportada do monitor. Esta parte contém
// código específico do problema dos filósofos

import (
	"fmt"
	printer "github.com/jaodsilv/ep3/filosofoPrinter"
	"math/rand"
	"time"
)

const (
	// Existe um erro na espera do processo, esse erro é da ordem de 0.1 milisegundos
	eatTime = 5 * time.Millisecond
)

var (
	come = func() { time.Sleep(eatTime) }
)

// Estados
const (
	Cheio = iota
	Comeu
	AcabouComida
)

// FilosofoMonitor é o monitor de fato.
type FilosofoMonitor struct {
	turno, filosofos, pratos int32
	pesoTotal, servicos      int32
	monitor                  *monitor
	panela                   chan bool
	uniforme                 bool
	printer                  *printer.Printer
}

// Init inicializa o monitor
func Init(n, pratos, pesoTotal int32, uniforme bool, p *printer.Printer) *FilosofoMonitor {
	// Firula
	rand.Seed(time.Now().UnixNano())
	variaveis := make([]string, n+3)
	for i := int32(0); i < n; i++ {
		variaveis[i] = getGarfoVar(n, i)
	}
	monitor := &FilosofoMonitor{
		turno:     int32(1),
		filosofos: n,
		pratos:    pratos,
		pesoTotal: pesoTotal,
		monitor:   newMonitor(variaveis...),
		panela:    make(chan bool),
		uniforme:  uniforme,
		printer:   p,
		servicos:  int32(0),
	}
	go monitor.garcom()
	return monitor
}

// Controla o serviço de pratos para não comam mais que o máximo
// Como ele é um contador independente que não compartilha a variável de pratos para escrita ele
// não precisa do monitor
func (f *FilosofoMonitor) garcom() {
	fmt.Println(f.pesoTotal)
	for ; f.servicos < f.pratos; f.servicos++ {
		if !f.uniforme && f.servicos == f.turno*f.pesoTotal {
			f.turno++
		}
		f.panela <- true
	}

	for i := int32(0); i < f.filosofos; i++ {
		f.panela <- false
	}
}

func (f *FilosofoMonitor) temComida() bool {
	if f.servicos == f.pratos {
		return false
	}
	return true
}

// Come inclui pegar os garfos, comer e soltar os garfos
func (f *FilosofoMonitor) Come(comidas, peso, index int32) int {
	if !f.uniforme && comidas >= f.turno*peso && f.temComida() {
		// Esta cheio?
		return Cheio
	}
	if !(<-f.panela) {
		// Acabou a comida na panela?
		return AcabouComida
	}

	// Oba, tem comida!!!
	f.pegaGarfos(index)
	f.printer.PrintInicio(index)
	come()
	f.printer.PrintFim(index)
	f.soltaGarfos(index)
	return Comeu
}

func getGarfoVar(n, index int32) string {
	return fmt.Sprintf("%d", index%n)
}

func (f *FilosofoMonitor) pegaGarfos(index int32) {
	left := getGarfoVar(f.filosofos, index)
	right := getGarfoVar(f.filosofos, index+1)
	if index%2 == 0 {
		f.monitor.wait(left)
		f.monitor.wait(right)
	} else {
		f.monitor.wait(right)
		f.monitor.wait(left)
	}
}

func (f *FilosofoMonitor) soltaGarfos(index int32) {
	left := getGarfoVar(f.filosofos, index)
	right := getGarfoVar(f.filosofos, index+1)
	f.monitor.signal(left)
	f.monitor.signal(right)
}
