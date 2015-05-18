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
    uma das listas chegar no corredor em outra tabela. */
  int volta, classif, pos, index, track;
  bool vivo, meia_pos, slow;
} rider;

typedef struct Metro {
  int ciclistas[4];
  int trilhas_ocupadas;
} metro;


bool acaba_volta = FALSE, uniform = FALSE, debug = FALSE;
rider *corredores;
metro *pista;
int *classificacao, n, ni, passo=0, d, eliminados;
char mode;
pthread_t *ciclista;

/* Sincronização */
int count1=0, count2=0;
pthread_mutex_t mutex1, mutex2, mutex_pista;
pthread_cond_t cond1, cond2;
sem_t sem;
/* -------------- */

void anda_passo(rider *corr);
void barreira1();
void barreira2();
void elimina_corredor(int i);
void finaliza_pista();
void finaliza_pthread();
char * get_nome(int i);
void imprime_3_ultimos();
void imprime_classificacao(bool fim);
void inicializa_pista();
void inicializa_pthread();
int * knuth_shuffle(int n);
void le_entrada();
void move_ciclistas();
void move_para_ultimo_vivo(int i, int n);
void set_debug(int argc, char *argv[]);
void swap_classif(int class1, int class2);
void * ThreadPasso (void * argument);

int main(int argc, char *argv[]) {
  set_debug(argc, argv);

  le_entrada();/* Ler entrada */

  srand(time(NULL));

  inicializa_pista(); 

  inicializa_pthread();

  while(n > 1) {
    move_ciclistas();
  }
  imprime_classificacao(TRUE);
   
  finaliza_pthread();
  finaliza_pista();
    
  exit(EXIT_SUCCESS);
}
/* Anda 72ms de um ciclista */
void anda_passo(rider *corr) {
  int pos, track, next, i;
  pthread_mutex_lock(&mutex_pista);
  pos = corr->pos;
  track = corr->track;
  if(pos == d-1) next = 0;
  else next = pos+1;
 
  if(pista[next].trilhas_ocupadas<4) {
    /* Coloca o ciclista da última trilha ocupada na trilha do ciclista atual */
    pista[pos].trilhas_ocupadas--;
    pista[pos].ciclistas[track] = pista[pos].ciclistas[pista[pos].trilhas_ocupadas];
    corredores[pista[pos].ciclistas[track]].track = track;

    /* Atualiza nova posição do ciclista */
    corr->pos = next;
    pista[next].ciclistas[pista[next].trilhas_ocupadas] = corr->index;
    corr->track = pista[next].trilhas_ocupadas;
    pista[next].trilhas_ocupadas++;

    /* Atualiza classificação */
    for(i=0; i<pista[pos].trilhas_ocupadas; i++) {
      if(corredores[pista[pos].ciclistas[i]].volta == corr->volta &&
        corredores[pista[pos].ciclistas[i]].classif < corr->classif) {
        corredores[pista[pos].ciclistas[i]].classif++;
        corr->classif--;
        classificacao[corredores[pista[pos].ciclistas[i]].classif] = pista[pos].ciclistas[i];
      }
    }
    classificacao[corr->classif]=corr->index;

    if(next==0) {
      corr->volta++;
      if(corr->classif == 0) acaba_volta = TRUE;
      if(uniform == FALSE) corr->slow = rand()%2;
    }
  }
  pthread_mutex_unlock(&mutex_pista);
}

void barreira1() {
  pthread_mutex_lock(&mutex1);
  count1++;
  if(count1 == n+1+eliminados) {
    count1=0;
    pthread_cond_broadcast(&cond1);
  } else {
    pthread_cond_wait(&cond1, &mutex1);
  }
  pthread_mutex_unlock(&mutex1);
}

void barreira2() {
  pthread_mutex_lock(&mutex2);
  count2++;
  if(count2 == n+1) {
    count2=0;
    pthread_cond_broadcast(&cond2);
  } else {
    pthread_cond_wait(&cond2, &mutex2);
  }
  pthread_mutex_unlock(&mutex2);
}

/* Elimina corredor de qualquer posicao */
void elimina_corredor(int i) {
  int index = classificacao[i];
  int pos = corredores[index].pos;
  int track = corredores[index].track;

  /* Marca como eliminado */
  corredores[index].vivo = FALSE;

  move_para_ultimo_vivo(i, n);
  n--;

  /* Remove cilicsta da pista */  
  pista[pos].trilhas_ocupadas--;
  pista[pos].ciclistas[track] = pista[pos].ciclistas[pista[pos].trilhas_ocupadas];
  corredores[pista[pos].ciclistas[track]].track = track;
}

void finaliza_pista() {
  free(ciclista);
  free(pista);
  free(corredores);
  free(classificacao);
}

void finaliza_pthread() {
    pthread_mutex_destroy(&mutex_pista);
    pthread_mutex_destroy(&mutex1);
    pthread_mutex_destroy(&mutex2);
    pthread_cond_destroy(&cond1);
    pthread_cond_destroy(&cond2);
}

/* Firula de nomes dos ciclistas. São nomes de personagens de filmes/jogos/desenhos de corrida. */
char * get_nome(int i) {
  char * nome[81] = {"Dick Dastardly and Muttley", "The Slag Brothers", "The Gruesome Twosome", 
      "Professor Pat Pending", "The Red Max", "Penelope Pitstop", "Sergeant Blast and Private Meekly", 
      "The Ant Hill Mob", "Lazy Luke and Blubber Bear", "Peter Perfect", "Rufus Ruffcut and Sawtooth", 
      "Takumi Fujiwara", "Bunta Fujiwara", "Ryosuke Takahashi", "Keisuke Takahashi", "Go Mifune(Speed Racer)", 
      "Racer X", "Ken'ichi Mifune(Rex Racer)", "Captain Terror", "Snake Oiler", "Dominic Toretto", 
      "Leticia \"Letty\" Ortiz", "Johnny Tran", "Brian O'Conner", "Mia Toretto", "Roman Pearce", "Tej Parker", 
      "Han Seoul-Oh", "Sean Boswell", "Gisele Yashar", "Carter Verone", "D.K.", "Arturo Braga", "Hernan Reyes", 
      "Owen Shaw", "Riley Hicks", "Deckard Shaw", "Vince", "Leon", "Chan Foh To", "Warner \"Cougar\" Kaugman", 
      "Lightning McQueen", "Strip \"The King\" Weathers", "Francesco Bernoulli", "Bill Whipple", "MacDonald" , 
      "Steve Grayson", "Steve McQueen", "Jim Douglas", "Peter Thorndyke", "Cole Trickle", "Russ Wheeler", "Guy Martin",
      "Ian Hutchinson", "Ayrton Senna", "Frankenstein", "\"Machine Gun\" Joe Viterbo", "Jensen Ames", "Nero the Hero", 
      "Matilda the Hun", "Calamity Jane", "Natasha Martin", "Sonoshee \"Cherry Boy Hunter\" McLaren", 
      "Jean-Pierre Sarti", "Pete Aron", "Scott Stoddard", "Nino Barlini", "Jeff", "Flat Top", "Michel Vaillant", 
      "Henri Vaillant", "Steve Warson", "The \"Leader\"", "Yves Douleac", "Gabriele Spangenberg", "Bob Cramer", 
      "Jimmy Bly", "Joe Tanto", "Memo Moreno", "Beau Brandenburg", "Aubrey James"};
  if (ni > 81) {
    sprintf(nome[0], "ciclista%d", i);
    return nome[0];
  } else {
    return nome[i];
  }
}

void imprime_3_ultimos() {
  int m = n+eliminados;
  printf("\n------------------------------------------\n");
  printf("\n3 últimos da volta %d\n\n", corredores[classificacao[0]].volta+1);
  if(m>=3) {
    printf("%dº: %s\n", m-2, get_nome(classificacao[m-3]));
  } else {
    printf("%dº: %s\n", m-1, get_nome(classificacao[m-2]));
  }
  
  if(m>=3) {
    if(eliminados == 2) {
      printf("Eliminado: %s\n", get_nome(classificacao[m-2]));
    } else {
      printf("%dº: %s\n", m-1, get_nome(classificacao[m-2]));
    }
  } else {
    if(m == 2 && eliminados == 1) {
      printf("Eliminado: %s\n", get_nome(classificacao[m-1]));
    } else {
      printf("%dº: %s\n", m, get_nome(classificacao[m-1]));
    }
  }

  if(m>=3) {
    if(eliminados >= 1) {
      printf("Eliminado: %s\n", get_nome(classificacao[m-1]));
    } else {
      printf("%dº: %s\n", m, get_nome(classificacao[m-1]));
    }
  }
}

void imprime_classificacao(bool fim) {
  int i;
  int v;
  char* ativo;
  printf("\n------------------------------------------\n");
  printf("Tempo: %dms\n\n", passo*72);
  if(fim) {
    printf("MEDALHA DE OURO\n");
    printf("1º: %s - %dª volta\n\n", get_nome(classificacao[0]), corredores[classificacao[0]].volta+1);
    printf("MEDALHA DE PRATA\n");
    printf("2º: %s - %dª volta\n\n", get_nome(classificacao[1]), corredores[classificacao[1]].volta+1);
    printf("MEDALHA DE BRONZE\n");
    printf("3º: %s - %dª volta\n\n", get_nome(classificacao[2]), corredores[classificacao[2]].volta+1);
  } else {
    if(corredores[classificacao[0]].slow) v = 25;
    else v = 50;
    printf("1º: %s - Ativo - %dª volta - v: %dkm/h\n", get_nome(classificacao[0]), corredores[classificacao[0]].volta+1, v);
    if(corredores[classificacao[1]].slow) v = 25;
    else v = 50;
    printf("2º: %s - Ativo - %dª volta - v: %dkm/h\n", get_nome(classificacao[1]), corredores[classificacao[1]].volta+1, v);
    if (corredores[classificacao[2]].vivo) {
      ativo = "Ativo";
    } else {
      ativo = "Eliminado";
    }
    if(corredores[classificacao[2]].slow) v = 25;
    else v = 50;
    printf("3º: %s - %s - %dª volta - v: %dkm/h\n", get_nome(classificacao[2]), ativo, corredores[classificacao[2]].volta+1, v);
  }
  for(i=3; i<ni; i++) {
    if (corredores[classificacao[i]].vivo) {
      ativo = "Ativo";
    } else {
      ativo = "Eliminado";
    }
    if(fim) {
      printf("%dº: %s - %dª volta\n", i+1, get_nome(classificacao[i]), corredores[classificacao[i]].volta+1);
    } else {
      if(corredores[classificacao[i]].slow) v = 25;
      else v = 50;
      printf("%dº: %s - %s - %dª volta - v: %dkm/h\n", i+1, get_nome(classificacao[i]), ativo, corredores[classificacao[i]].volta+1, v);
    }
  }
}

void inicializa_pista() {
  int i;

  pista = malloc(d*sizeof(metro));
  if(pista == NULL) {
    printf("Socorro! malloc devolveu NULL!\n");
    exit(EXIT_FAILURE);
  }
  classificacao = knuth_shuffle(ni);
  corredores = malloc(n*sizeof(rider));
  if(corredores == NULL) {
    printf("Socorro! malloc devolveu NULL!\n");
    exit(EXIT_FAILURE);
  }
  for(i=0; i<n; i++) {
    pista[i].ciclistas[0]=classificacao[n-i-1];
    pista[i].trilhas_ocupadas = 1;
    corredores[classificacao[i]].volta = 0;
    corredores[classificacao[i]].classif = i;
    corredores[classificacao[i]].pos = n-i-1;
    corredores[classificacao[i]].index = classificacao[i];
    corredores[classificacao[i]].track = 0;
    corredores[classificacao[i]].vivo = TRUE;
    corredores[classificacao[i]].meia_pos = FALSE;
    if(uniform) corredores[classificacao[i]].slow = FALSE;
    else corredores[classificacao[i]].slow = TRUE;
  }
  for(i=n; i<d; i++) {
    pista[i].trilhas_ocupadas = 0;
  }
}

void inicializa_pthread() {
  long i;
  pthread_mutex_init(&mutex1, NULL);
  pthread_mutex_init(&mutex2, NULL);
  pthread_mutex_init(&mutex_pista, NULL);
  pthread_cond_init(&cond1, NULL);
  pthread_cond_init(&cond2, NULL);

  ciclista = malloc(n*sizeof(pthread_t));
  if(ciclista == NULL) {
    printf("Socorro! malloc devolveu NULL!\n");
    exit(EXIT_FAILURE);
  }
  for(i=0; i<n; i++) {
    if (pthread_create(&ciclista[i], NULL, ThreadPasso, (void *) i)) {
      printf("\n ERROR creating thread 1");
      exit(EXIT_FAILURE);
    }
  }
}

int * knuth_shuffle(int n) {
  int i, *inteiros, aux1, aux2;
  inteiros = malloc(n*sizeof(int));
  if(inteiros == NULL) {
    printf("Socorro! malloc devolveu NULL!\n");
    exit(EXIT_FAILURE);
  }
  for(i=0; i<n; i++) {
    inteiros[i]=i;
  }
  for(i=n-1; i>0; i--) {
    aux1 = rand()%(i+1);
    aux2 = inteiros[i];
    inteiros[i] = inteiros[aux1];
    inteiros[aux1] = aux2;
  }
  return inteiros;
}

void le_entrada() {
  printf("Entre o tamanho da pista(m), número de ciclistas e o modo(u/v)\n");
  scanf("%d %d %c", &d, &ni, &mode);
  n = ni;
  uniform = (mode == 'u');
}

void move_ciclistas() {
  /* Executa passo */
  if(debug && !(passo % 200)) {
    imprime_classificacao(FALSE);
  }
  passo++;
  
  /* Autoriza começar e espera terminar o turno dos ciclistas. */
  barreira1();
  eliminados = 0;
  barreira2();

  if (acaba_volta) {
    eliminados = 0;
    if(corredores[classificacao[0]].volta%2 == 0) {
      elimina_corredor(n-1);
      eliminados++;
    }
    if(corredores[classificacao[0]].volta%4 == 0 && n>3 && !(rand() % 100)) {
      elimina_corredor(rand()%n);
      eliminados++;
    }
    imprime_3_ultimos(eliminados);
    acaba_volta = FALSE;
  }
}

void move_para_ultimo_vivo(int i, int n) {
  int classif;
  rider* corr = &corredores[classificacao[i]];
  for(classif = i; classif<n-1; classif++) {
    classificacao[classif] = classificacao[classif+1];
    corredores[classificacao[classif]].classif = classif;
  }
  classificacao[n-1] = corr->index;
  corr->classif = n-1;
}

void set_debug(int argc, char *argv[]) {
 /* Tratamento de opções, no caso, apenas debug ou nada é aceito */
  if(argc == 2) {
    if(!strcmp("debug", argv[1])) {
      printf("Debug set to TRUE\n");
      debug = TRUE;
    } else {
      printf("Invalid option %s\n", argv[1]);
      exit(EXIT_FAILURE);
    }
  } else if(argc>2) {
    printf("Invalid number of options: %d\n", argc-1);
    exit(EXIT_FAILURE);
  }
}

void swap_classif(int class1, int class2) {
  int aux = classificacao[class1];
  classificacao[class1] = classificacao[class2];
  classificacao[class2] = aux;

  corredores[aux].classif = class2;
  corredores[classificacao[class1]].classif = class1;
}

void * ThreadPasso (void * argument) {
  long index = (long) argument;
  
  rider *corr = &corredores[index];
  while(TRUE) {
    barreira1();
    if(corr->vivo) {
      if(corr->slow && !corr->meia_pos){
        corr->meia_pos=TRUE;
      } else {
        corr->meia_pos=FALSE;
        anda_passo(corr);
      }
      barreira2();
    } else {
      pthread_exit(NULL);
    }
  }
  pthread_exit(NULL);
}