// cos(x) = Sum_0^∞((−1)^n * x^(2n) / (2n)!)
package main

import (
	"runtime"
	"fmt"
	"flag"
	"strconv"
	"math/big"
	"sync"
)

// Meu tipo de barreira
type Barreira func(i int64) bool

var printPrecision = 100000

var precision, x, cos *big.Rat
var xString string
var fat []*big.Int //fat[i] = 2i!
var pot []*big.Rat //pot[i] = x^(2i)
var pexp int
var f, d, s bool

// variável global usada pelas barreiras
var stop = false

// Canais são a forma padrão de sincronização em Go
var sumPart chan *big.Rat

func potx2(turno, q int64) {
	for i := q*(turno-1); i < q*turno; i++ {
		pot = append(pot, new(big.Rat).Mul(pot[i], x))
		pot[i+1].Mul(pot[i+1], x)
	}
}

func fat2(turno, q int64) {
	for i := q*(turno-1); i < q*turno; i++ {
		fat = append(fat, new(big.Int).Mul(fat[i], big.NewInt((2*i+1)*(2*i+2))))
	}
}

func potx2C(turno, q int64, c chan bool) {
	potx2(turno, q)
	c <- true
}

func fat2C(turno, q int64, c chan bool) {
	fat2(turno, q)
	c <- true
}

// Barreiras iguais as do EP1, agora usando closures
func newBarreira(q int64) Barreira {
	mutex := &sync.Mutex{}
	cond  := sync.NewCond(mutex)
	count := int64(0)
	s := false

	return func(i int64) bool {
		mutex.Lock()
		
  		count++
  		if count == q+1 {
	    	count = 0

	    	// A closure garante que todos que entram na barreira saem com o mesmo valor de s
	    	s = !stop

	  		if d {
	  			fmt.Printf("%d\n", i)
	  			fmt.Println("Cos(", xString, ") = ", cos.FloatString(printPrecision))
	  		}
	    	cond.Broadcast()
	  	} else {
	    	if d {
	    		fmt.Printf("%d, ", i)
	    	}
	    	cond.Wait()
	  	}

	  	mutex.Unlock()
	  	return s
	}
}

func calcula(index, q int64, barreira Barreira) {
	var sum *big.Rat
	n := index

	//(−1)^n * x^(2n) / (2n)!
	// O começo do for serve como barreira e sinal se deve continuar no loop
	for barreira(index) {
		sum = big.NewRat(int64(0), int64(1))
		for i := int64(0); i < q; i++ {
			elem := new(big.Rat)
			elem.SetInt(fat[n])
			elem.Quo(pot[n], elem)
			if !f && elem.Cmp(precision) == -1 {
				stop = true
			}
			neg := (n%2 == 1)
			if neg {
				elem.Neg(elem)
			}
			sum.Add(sum, elem)
			n += q
		}
		// espera todos terminarem, espera sinal da thread principal
		sumPart <- sum
	}
}

func main() {
	// Lê argumentos do terminal
	flag.Parse()

	// Número de threads
	q, err := strconv.Atoi(flag.Arg(0))
	if err != nil {
		fmt.Printf("Error on first argument: %s\n", err)
		return
	}
	if q == 0 {
		q = runtime.NumCPU()
	}
	q64 := int64(q)
	q2 := q64*q64

	// f ou m
	f = (flag.Arg(1) == "f")
	
	// valor de f ou m
	pexp, err = strconv.Atoi(flag.Arg(2))
	if err != nil {
		fmt.Printf("Error on third argument: %s\n", err)
		return
	}
	precision = new(big.Rat)
	precision.SetString(fmt.Sprintf("1e-%d", pexp))

	// Undefined precision float
	x = new(big.Rat)
	xString = flag.Arg(3)
	_, success := x.SetString(xString)
	if !success {
		fmt.Printf("Error parsing fourth argument\n")
		return
	}

	// d ou s
	d = false
	s = false
	if flag.NArg() == 5 {
	 	d = (flag.Arg(4) == "d")
	 	s = (flag.Arg(4) == "s")
	}


	barreira := newBarreira(q64)
	
	sumPart = make(chan *big.Rat)
	fatpot := make(chan bool)

	turno := int64(0)
	cos = big.NewRat(int64(0), int64(1))
	fat = []*big.Int{big.NewInt(int64(1))}
	pot = []*big.Rat{big.NewRat(int64(1), int64(1))}

	if !s {
		cosOld := big.NewRat(int64(0), int64(1))

		go fat2C(1, q2, fatpot)
		go potx2C(1, q2, fatpot)


		for i := int64(0); i < q64; i++ {
			go calcula(i, q64, barreira)
		}

		for !stop {
			turno++
			
			// espera o cálculo de fatorial e de potencias
			<-fatpot
			<-fatpot

			// thread principal também entra na barreira e também entra na conta
			if !barreira(q64) {
				turno--
				break
			}

			go fat2C(turno+1, q2, fatpot)
			go potx2C(turno+1, q2, fatpot)
		
			//espera todos terminarem
			for i := int64(0); i < q64; i++ {
				cos.Add(cos, <-sumPart)
			}

			// espera o cálculo de potência e 

			if f && cosOld.Sub(cos, cosOld).Abs(cosOld).Cmp(precision) == -1 {
				stop = true
			} else {
				cosOld.Set(cos)
			}
		}
	} else {
		elem := new(big.Rat)
		for !stop {
			elem = new(big.Rat)
			fat2(turno+1, int64(1))
			potx2(turno+1, int64(1))
			elem.SetInt(fat[turno])
			elem.Quo(pot[turno], elem)
			if elem.Cmp(precision) == -1 {
				stop = true
			}
			neg := (turno%2 == 1)
			if neg {
				elem.Neg(elem)
			}
			cos.Add(cos, elem)
			turno++
			fmt.Println("Cos(", xString, ") = ", cos.FloatString(printPrecision))
		}
	}

	if !s {
		fmt.Println("Número de rodadas: ", turno)
	} else {
		fmt.Println("Número de Termos: ", turno)
	}
	fmt.Println("Cos(", xString, ") = ", cos.FloatString(printPrecision))
}