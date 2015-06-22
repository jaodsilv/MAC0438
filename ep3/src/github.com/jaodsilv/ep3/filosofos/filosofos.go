// Package filosofos é um simulador do problema do jantar dos filósofos
package filosofos

import (
	"bufio"
	"fmt"
	printer "github.com/jaodsilv/ep3/filosofoPrinter"
	"github.com/jaodsilv/ep3/monitor"
	"io"
	"os"
	"strconv"
	"sync"
)

func readInt(reader *bufio.Reader, delim byte) int32 {
	s, err := reader.ReadString(delim)
	if err != nil && err != io.EOF {
		fmt.Printf("Error reading file.\n")
		panic(err)
	}
	if err == nil {
		s = s[:len(s)-1]
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		fmt.Printf("Error converting string %s to int.\n", s)
		panic(err)
	}
	return int32(i)
}

func parseOptions(filepath, rString, uString string) (*os.File, int32, bool) {
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Printf("Error reading file %s: %s\n", filepath, err)
		panic(err)
	}

	r, err := strconv.Atoi(rString)
	if err != nil {
		fmt.Printf("Error on second argument %s: %s\n", rString, err)
		panic(err)
	}

	// U ou P
	u := (uString == "U")

	return file, int32(r), u
}

// Run é a execução em si
func Run(filepath, rString, uString string) []*Filosofo {
	// Arquivo a ser lido
	file, r, u := parseOptions(filepath, rString, uString)
	reader := bufio.NewReader(file)
	// número de filosofos
	n := readInt(reader, '\n')

	wg := &sync.WaitGroup{}
	filosofos := make([]*Filosofo, n)
	pesoTotal := int32(0)

	p := &printer.Printer{}
	p.InitNomes(n)

	for i := int32(0); i < n; i++ {
		peso := int32(1)
		if i != n-1 {
			peso = readInt(reader, ' ')
		} else {
			peso = readInt(reader, '\n')
		}
		pesoTotal += peso
		filosofos[i] = newFilosofo(i, peso)
	}

	m := monitor.Init(n, r, pesoTotal, u, p)

	// Deixei separado do loop anterior para que a diferença de tempo do inicio do jantar de cada filósofo seja igual
	for _, f := range filosofos {
		wg.Add(1)
		go f.janta(m, wg)
	}

	file.Close()

	// espera todos terminarem de jantar
	wg.Wait()

	// Imprimir tudo serialmente pode demorar muito se forem muitos filósofos (Em especial
	// para o teste feito com mais de 900 filósofos demora muuuuuuuito)
	for _, f := range filosofos {
		p.PrintFilosofo(f.Index, f.Peso, f.Comidas)
	}

	// Este retorno é para ser usado nos meus testes.
	return filosofos
}
