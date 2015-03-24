/****************************************/
/** Aluno: João Marco Maciel da Silva  **/
/** Matéria: MAC0438                   **/
/** Prof: Daniel Macedo Batista        **/
/** EP1: Corrida por eliminação        **/
/****************************************/



#include <pthread.h>
#include <stdio.h>
#include <stdlib.h>
#include <time.h>


/* 
 * 1 passo = 72ms
 * 50 km/h = 1m/passo
 * 25 km/h = 1m/2passos
 *   14,4s = 200 passos
 */

typedef struct Rider {
  int classificação, volta, index;
} rider;

int main(int argc, char *argv[]) {
  int d, n, passo = 0;
  char mode;
  bool uniform, debug;
  pthread_t *ciclista, *(pista[4]);

  /* Tratamento de opções, no caso, apenas debug ou nada é aceito */
  if(argc == 1) {
    debug = false;
  } else if(argc == 2) {
    if(!strcmp("debug", argv[1])) {
      debug = true;
    } else {
      printf("Invalid option %s\n", argv[1]);
      exit(EXIT_FAILURE);
    }
  } else {
      printf("Invalid number of options: %d\n", argc-1);
      exit(EXIT_FAILURE);
  }

  /* Ler entrada */
  scanf("%d %d %c", &d, &n, &mode);
  uniform = (mode == 'u');

  ciclista = malloc(n*sizeof(pthread_t));
  if(ciclista == null) {
    printf("Socorro! malloc devolveu NULL!\n");
    exit(EXIT_FAILURE);
  }

  pista = malloc(d*sizeof(pthread_t[4]));
  if(pista == null) {
    printf("Socorro! malloc devolveu NULL!\n");
    exit(EXIT_FAILURE);
  }

  srand(time(NULL));

  while(n) {
    passo++;
    /* Executa passo */
      /*Anda ciclistas */
    if(debug && !(passo % 200)) {
      /* Imprime tudo */
    }
    /* Se volta acabou? */
      /* Se volta é par elimina o último corredor */
      n--;
      if(n < 3) {
        /* Imprime eliminado */
      }
      /* Se volta é 0 mod4 && n < 3*/
      if(n>3 && !(rand() % 100)) {
        /* Elimina corredor rand() % n*/
        n--;
      }
  }

  free(ciclista);
  free(pista)
  exit(EXIT_SUCCESS);
}
