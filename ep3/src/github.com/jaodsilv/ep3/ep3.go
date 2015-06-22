package main

import (
	"flag"
	"github.com/jaodsilv/ep3/filosofos"
)

func main() {
	flag.Parse()
	filosofos.Run(flag.Arg(0), flag.Arg(1), flag.Arg(2))
}
