#!/bin/bash

# Pega o tempo real, usuário, system e percentual de CPU usado no formato CSV
TIMEFORMAT='%E, %U, %S, %P'

function roda50 {
	#$1 -> arquivo
	#$2 -> pratos servidos
	#$3 -> modo
	#$4 -> amostra
	TIMEFORMAT=$4',%E,%U,%S,%P'
	sep='_'
	csv='.csv'
	for j in `seq 1 50`;
	do
		{ time bin/ep3 $1 $2 $3 1> /dev/null; } 2>> $4$sep$3$csv
	done
}

small="/home/joao/workspace/go/src/github.com/jaodsilv/ep3/sampleSmall"
big="/home/joao/workspace/go/src/github.com/jaodsilv/ep3/sampleBig"
big2="/home/joao/workspace/go/src/github.com/jaodsilv/ep3/sampleBig2"
rs=10000
rb=1000000

echo "Inicio"
echo "Teste S U"
roda50 $small $rs U "S"
echo "Teste S P"
roda50 $small $rs P "S"
echo "Teste B U"
roda50 $big $rb U "B"
echo "Teste B P"
roda50 $big $rb P "B"
echo "Teste B2 U"
roda50 $big2 $rb U "B2"
echo "Teste B2 P"
roda50 $big2 $rb P "B2"