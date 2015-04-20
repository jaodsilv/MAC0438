/****************************************/
/** Aluno: João Marco Maciel da Silva  **/
/** Matéria: MAC0438                   **/
/** Prof: Daniel Macedo Batista        **/
/** EP1: Corrida por eliminação        **/
/****************************************/



#include <pthread.h>
#include <semaphore.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>


/* 
 * 1 passo = 72ms
 * 50 km/h = 1m/passo
 * 25 km/h = 1m/2passos
 *   14,4s = 200 passos
 */

#define FALSE 0
#define TRUE 1

 typedef char bool; 

 typedef struct Rider {
  /* Classif, Pos, Index ajudam a de qualquer
    uma das tabelas chegar no corredor em outra tabela. */
  int volta, classif, pos, index, track;
  bool vivo, meia_pos, slow;
} rider;

typedef struct Metro {
  int ciclistas[4];
  int trilhas_ocupadas;
} metro;


bool acaba_volta;
rider *corredores;
int *classificacao, n, ni, passo = 0, d;
char* nomes[82] = {"Dick Dastardly and Muttley", "The Slag Brothers", 
"The Gruesome Twosome", "Professor Pat Pending", "The Red Max", 
"Penelope Pitstop", "Sergeant Blast and Private Meekly", "The Ant Hill Mob", 
"Lazy Luke and Blubber Bear", "Peter Perfect", "Rufus Ruffcut and Sawtooth",
"Takumi Fujiwara", "Bunta Fujiwara", "Ryosuke Takahashi", "Keisuke Takahashi",
"Go Mifune(Speed Racer)", "Racer X", "Ken'ichi Mifune(Rex Racer)", 
"Captain Terror", "Snake Oiler", "Dominic Toretto", "Leticia \"Letty\" Ortiz",
"Johnny Tran", "Brian O'Conner", "Mia Toretto", "Roman Pearce", "Tej Parker",
"Han Seoul-Oh", "Sean Boswell", "Gisele Yashar", "Carter Verone", "D.K.", 
"Arturo Braga", "Hernan Reyes", "Owen Shaw", "Riley Hicks", "Deckard Shaw", 
"Vince", "Leon", "Chan Foh To", "Warner \"Cougar\" Kaugman", 
"Lightning McQueen", "Strip \"The King\" Weathers", "Francesco Bernoulli", 
"Bill Whipple", "MacDonald" , "Steve Grayson", "Steve McQueen", "Jim Douglas",
"Peter Thorndyke", "Cole Trickle", "Russ Wheeler", "Guy Martin", 
"Ian Hutchinson", "Ayrton Senna", "Frankenstein", 
"\"Machine Gun\" Joe Viterbo", "Jensen Ames", "Nero the Hero", 
"Matilda the Hun", "Calamity Jane", "Natasha Martin", 
"Sonoshee \"Cherry Boy Hunter\" McLaren", "Jean-Pierre Sarti", "Pete Aron",
"Scott Stoddard", "Nino Barlini", "Nino Barlini", "Jeff", "Flat Top", 
"Michel Vaillant", "Henri Vaillant", "Steve Warson", "The \"Leader\"",
"Yves Douleac", "Gabriele Spangenberg", "Bob Cramer", "Jimmy Bly",
"Joe Tanto", "Memo Moreno", "Beau Brandenburg", "Aubrey James"};
pthread_mutex_t mutex_pista, mutex_count;
int count;
pthread_barrier_t barreira;


int anda_passo(rider *corr) {
  pthread_mutex_lock(&mutex_pista);
  pos = corr->pos;
  classif = corr->classif;
  track = corr->track;
  if(pos==d-1) next = 0;
  else next = pos+1;
  
  if(pista[next].trilhas_ocupadas<4) {
    pista[next].ciclistas[pista[next].trilhas_ocupadas] = corr->index;
    pista[next].trilhas_ocupadas++;
    for(i=track+1; i<pista[pos].trilhas_ocupadas; i++) {
      corredores[pista[pos].ciclistas[i]].track--;
      if(corredores[pista[pos].ciclistas[i]].volta==corr->volta &&
        corredores[pista[pos].classificacao[i]].classif>classif) {
        corredores[pista[pos].classificacao[i]].classif--;
      corr->classif++;
    }
    pista[pos].ciclistas[i-1] = pista[pos].ciclistas[i];
  }

  pista[pos].trilhas_ocupadas--;
  pthread_mutex_unlock(&mutex_pista);
  if(next==0) {
    corr->volta++;
    if(corr->classif==0) acaba_volta = TRUE;
    corr->slow = rand()%2;
  }
} else {
  pthread_mutex_unlock(&mutex_pista);
}
}

void * ThreadPasso (void * argument) {
  long index = (long) argument;
  int pos, next, i;
  rider *corr = &corredores[index];

  while(corr->vivo) {
    pthread_barrier_wait(&barrier);
    if(corr->slow && !corr->meia-pos){
      corr->meia-pos=TRUE;
    } else if((corr->slow && corr->meia-pos) || 
      !corr->slow) {
      anda_passo(corr);
    }

    pthread_mutex_lock(&mutex_count);
    count++;
    pthread_mutex_unlock(&mutex_count);
  }
  pthread_exit(NULL);
}

void imprime_classificacao() {
  int i;
  char* ativo;
  printf("Tempo: %dms\n", passo*72);
  for(i=0; i<ni; i++) {
    if (corredores[classificacao[i]].vivo) {
      ativo = "Ativo";
    } else {
      ativo = "Eliminado";
    }
    if (ni>82) {
      printf("%d: %d - %s\n", i+1, classificacao[i], ativo);
    } else {
      printf("%d: %s - %s\n", i+1, nomes[classificacao[i]], ativo);
    }
  }
}

/* Esse método é chamado apenas quando todas as threads estão na barreira */
/* Então não há necessidade de passar pelo semáforo para alterar os vetores */
void swap_classif(int class1, int class2) {
  int aux = classificacao[class1];
  classificacao[class1] = classificacao[class2];
  classificacao[class2] = aux;

  corredores[aux].classif = class2;
  corredores[classificacao[class1]].classif = class1;
}

void move_to_last(int i, int n) {
  int classif;
  int aux;
  for(classif = corredores[i].classif; i<n-1; classif++) {
    aux = classificacao[classif];
    classificacao[classif] = classificacao[classif+1];
    classificacao[classif+1] = aux;

    corredores[aux].classif = classif+1;
    corredores[classificacao[classif]].classif = classif;
  }
}

int main(int argc, char *argv[]) {
  int voltas=0;
  metro *pista;
  long i;
  int j;
  char mode;
  bool uniform, debug;
  pthread_t *ciclista;
  char* medalha[3] = {"ouro", "prata", "bronze"};
  
    /* Tratamento de opções, no caso, apenas debug ou nada é aceito */
  if(argc == 1) {
    debug = FALSE;
  } else if(argc == 2) {
    if(!strcmp("debug", argv[1])) {
      debug = TRUE;
    } else {
      printf("Invalid option %s\n", argv[1]);
      exit(EXIT_FAILURE);
    }
  } else {
    printf("Invalid number of options: %d\n", argc-1);
    exit(EXIT_FAILURE);
  }

  /* Ler entrada */
  printf("Entre o tamanho da pista(m), número de ciclistas e o modo(u/v)\n");
  scanf("%d %d %c", &d, &ni, &mode);
  n = ni;
  uniform = (mode == 'u');

  pthread_mutex_init(&mutex_pista,NULL);
  pthread_mutex_init(&mutex_count, NULL);
  pthread_barrierattr_t attr;
  pthread_barrierattr_init(&attr);
  pthread_barrier_init(&barreira, &attr, ni+1);
  pthread_cond_t cond;
  pthread_cond_init(&cond, NULL);
  count = 0;

  ciclista = malloc(n*sizeof(pthread_t));
  if(ciclista == NULL) {
    printf("Socorro! malloc devolveu NULL!\n");
    exit(EXIT_FAILURE);
  }
  for(i=0; i<n; i++) {
    if (pthread_create(&ciclista[i], NULL, ThreadPasso, (void *) i)) {
      printf("\n ERROR creating thread 1");
      exit(1);
    }
  }

  pista = malloc(d*sizeof(metro));
  if(pista == NULL) {
    printf("Socorro! malloc devolveu NULL!\n");
    exit(EXIT_FAILURE);
  }

  classificacao = malloc(n*sizeof(int));
  if(classificacao == NULL) {
    printf("Socorro! malloc devolveu NULL!\n");
    exit(EXIT_FAILURE);
  }

  corredores = malloc(n*sizeof(rider));
  if(corredores == NULL) {
    printf("Socorro! malloc devolveu NULL!\n");
    exit(EXIT_FAILURE);
  }

  if(uniform) {
      /* coloca todas as velocidades para 50km/h */
    for (j=0; j<n; j++) {
      corredores[j].slow=FALSE;
    }
  }

  srand(time(NULL));

  while(n>1) {
    passo++;
    /* Executa passo */
    if(debug && !(passo % 200)) {
      imprime_classificacao();
    }
      pthread_barrierattr_setpshared(&attr, ni+1);
      /*<corram!!!>
      <await(terminados==n);>
      terminados = 0*/
      if (acaba_volta) {
        if(corredores[classificacao[0]].volta%2 == 0) {
          n--;
          corredores[classificacao[n]].vivo=FALSE;
        }
        if(corredores[classificacao[0]].volta%4 == 0 && n>3 && !(rand() % 100)) {
          /* colocar um corredor aleatorio em ultimo e eliminá-lo*/
          move_to_last(rand()%n, n);
          n--;
          corredores[classificacao[n]].vivo=FALSE;
        }
        if(n < 3) {
          /* Imprime eliminado */
          if (ni>82) {
            printf("medalha de %s: ciclista%i", medalha[n], classificacao[n]);
          } else {
            printf("medalha de %s: %s", medalha[n], nomes[classificacao[n]]);
          }
        }
        acaba_volta = FALSE;
      }
      pthread_barrierattr_setpshared(&attr, ni+1);
    }
    if (ni>82) {
      printf("medalha de %s: ciclista%i", medalha[0], classificacao[0]);
    } else {
      printf("medalha de %s: %s", medalha[0], nomes[classificacao[0]]);
    }
    
    pthread_mutex_destroy(&mutex_pista);
    free(ciclista);
    free(pista);
    exit(EXIT_SUCCESS);
  }
