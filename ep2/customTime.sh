#!/bin/bash

# Pega o tempo real, usuário, system e percentual de CPU usado no formato CSV
TIMEFORMAT='%E, %U, %S, %P'

function roda30 {
	#$1 -> numero de threads
	#$2 -> modo
	#$3 -> precisao
	#$4 -> x
	#$5 -> ultimo argumento
	TIMEFORMAT=$1','$3','$4',%E,%U,%S,%P'
	#TIMEFORMAT=$1$TIMEFORMAT
	sep='_'
	csv='.csv'
	for j in `seq 1 30`;
	do
		{ time bin/cos $1 $2 $3 $4 $5 1> /dev/null; } 2>> $2$sep$5$csv
	done
}

function teste {
	#$1 -> modo
	#$2 -> precisao
	#$3 -> x
	#$4 -> max threads
	echo "precisao 10^-$2"
	echo "S"
	roda30 1 $1 $2 $3 s

	# for t in `seq 10 $4`;
	# do
	# 	echo "$t threads"
	# 	roda30 $t $1 $2 $3 ''
	# done

	# for t in `seq 1 $4`;
	# do
	# 	echo "$t threads, D"
		# roda30 $t $1 $2 $3 d
	# done
}

pi=3.14159265359

echo "Inicio"
echo "x=PI"
echo "modo f"
# teste f 10 $pi 10
# teste f 100 $pi 10
# teste f 1000 $pi 10
echo "modo m"
# teste m 10 $pi 10
# teste m 100 $pi 10
# teste m 1000 $pi 10

echo "x=1"
echo "modo f"
# teste f 10 1 20
# teste f 100 1 20
# teste f 1000 1 20
teste f 10000 1 16
echo "modo m"
# teste m 10 1 20
# teste m 100 1 20
# teste m 1000 1 20
# teste m 10000 1 16

echo "x=0"
echo "modo f"
# teste f 10 0 50
# teste f 100 0 50
# teste f 1000 0 50
# teste f 10000 0 50
# teste f 100000 0 50
echo "modo m"
# teste m 10 0 50
# teste m 100 0 50
# teste m 1000 0 50
# teste m 10000 0 50
# teste m 100000 0 50

echo "x=0.1"
echo "modo f"
# teste f 10 0.1 20
# teste f 100 0.1 20
# teste f 1000 0.1 20
teste f 10000 0.1 10
echo "modo m"
# teste m 10 0.1 16
# teste m 100 0.1 16
# teste m 1000 0.1 16
# teste m 10000 0.1 10
