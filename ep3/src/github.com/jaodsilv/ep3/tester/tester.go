// Este package é só para ter resultados em um formato agradável para fazer análises.
package main

import (
	"fmt"
	printer "github.com/jaodsilv/ep3/filosofoPrinter"
	"github.com/jaodsilv/ep3/filosofos"
	"os"
)

const (
	rs  = 10000
	rb  = 1000000
	rb2 = 1000000
	// Se for rodar este teste troque pelo endereço correto
	smallPath      = "/home/joao/workspace/go/src/github.com/jaodsilv/ep3/sampleSmall"
	bigPath        = "/home/joao/workspace/go/src/github.com/jaodsilv/ep3/sampleBig"
	bigPath2       = "/home/joao/workspace/go/src/github.com/jaodsilv/ep3/sampleBig2"
	u              = "U"
	p              = "P"
	numeroDeTestes = 50
)

var rsString = fmt.Sprintf("%d", rs)
var rbString = fmt.Sprintf("%d", rb)
var rb2String = fmt.Sprintf("%d", rb2)

func main() {

	f1, f2, f3, f4, f5, f6 := make([][]*filosofos.Filosofo, numeroDeTestes), make([][]*filosofos.Filosofo, numeroDeTestes), make([][]*filosofos.Filosofo, numeroDeTestes), make([][]*filosofos.Filosofo, numeroDeTestes), make([][]*filosofos.Filosofo, numeroDeTestes), make([][]*filosofos.Filosofo, numeroDeTestes)

	log, err := os.Create("/tmp/logEP3")
	if err != nil {
		fmt.Printf("Error creating file.\n")
		panic(err)
	}

	log.WriteString("Criando arquivos\n")
	fsu, err := os.Create("/tmp/Small-U.csv")
	if err != nil {
		fmt.Printf("Error creating file.\n")
		panic(err)
	}
	fsp, err := os.Create("/tmp/Small-P.csv")
	if err != nil {
		fmt.Printf("Error creating file.\n")
		panic(err)
	}
	fbu, err := os.Create("/tmp/Big-U.csv")
	if err != nil {
		fmt.Printf("Error creating file.\n")
		panic(err)
	}
	fbp, err := os.Create("/tmp/Big-P.csv")
	if err != nil {
		fmt.Printf("Error creating file.\n")
		panic(err)
	}
	fb2u, err := os.Create("/tmp/Big2-U.csv")
	if err != nil {
		fmt.Printf("Error creating file.\n")
		panic(err)
	}
	fb2p, err := os.Create("/tmp/Big2-P.csv")
	if err != nil {
		fmt.Printf("Error creating file.\n")
		panic(err)
	}
	log.WriteString("Arquivos criados\n")

	log.WriteString("Rodando Testes\n")
	for i := 0; i < numeroDeTestes; i++ {
		log.WriteString(fmt.Sprintf("Run %d\n", i))
		log.WriteString("Small U\n")
		f1[i] = filosofos.Run(smallPath, rsString, u)
		log.WriteString("Small P\n")
		f2[i] = filosofos.Run(smallPath, rsString, p)
		log.WriteString("Medium U\n")
		f3[i] = filosofos.Run(bigPath, rbString, u)
		log.WriteString("Medium P\n")
		f4[i] = filosofos.Run(bigPath, rbString, p)
		log.WriteString("Big U\n")
		f5[i] = filosofos.Run(bigPath2, rb2String, u)
		log.WriteString("Big P\n")
		f6[i] = filosofos.Run(bigPath2, rb2String, p)
	}

	log.WriteString("Gravando em Arquivos\n")
	// Invertendo os índices para ficar mais fácil fazer análises
	log.WriteString("Small\n")
	for fi := range f1[0] {
		log.WriteString(fmt.Sprintf("Filósofo %d\n", fi))
		log.WriteString("U\n")
		for _, fj := range f1 {
			f := fj[fi]
			fsu.WriteString(printer.FilosofoCSV(f.Index, f.Peso, f.Comidas))
		}
		log.WriteString("P\n")
		for _, fj := range f2 {
			f := fj[fi]
			fsp.WriteString(printer.FilosofoCSV(f.Index, f.Peso, f.Comidas))
		}
	}
	log.WriteString("Medium\n")
	for fi := range f3[0] {
		log.WriteString(fmt.Sprintf("Filósofo %d\n", fi))
		log.WriteString("U\n")
		for _, fj := range f3 {
			f := fj[fi]
			fbu.WriteString(printer.FilosofoCSV(f.Index, f.Peso, f.Comidas))
		}
		log.WriteString("P\n")
		for _, fj := range f4 {
			f := fj[fi]
			fbp.WriteString(printer.FilosofoCSV(f.Index, f.Peso, f.Comidas))
		}
	}
	log.WriteString("Big\n")
	for fi := range f5[0] {
		log.WriteString(fmt.Sprintf("Filósofo %d\n", fi))
		log.WriteString("U\n")
		for _, fj := range f5 {
			f := fj[fi]
			fb2u.WriteString(printer.FilosofoCSV(f.Index, f.Peso, f.Comidas))
		}
		log.WriteString("P\n")
		for _, fj := range f6 {
			f := fj[fi]
			fb2p.WriteString(printer.FilosofoCSV(f.Index, f.Peso, f.Comidas))
		}
	}

	log.WriteString("Fechando arquivos\n")
	fsu.Close()
	fsp.Close()
	fbu.Close()
	fbp.Close()
	fb2u.Close()
	fb2p.Close()
	log.Close()
}
