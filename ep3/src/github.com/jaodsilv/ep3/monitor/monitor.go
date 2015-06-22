package monitor

// Este arquivo é a parte que contém a parte básica do monitor, isto é, as funções wait(cv),
// signal(cv) e signal_all(cv)
// O código desta parte é reutilizável, e não contém nada específico do problema dos filósofos.

import (
	"sync"
)

// Em go são esportados métodos e strucs que começam com letra maiúsculas.
// Atributos de structs que começam com letra minúscula não são exportados.
// Variáveis permanentes, uma por condição.
type cond struct {
	ocupado bool
	mutex   *sync.Mutex
	fila    []chan bool
}

// Variáveis permanentes, uma por monitor
// monitor genérico
type monitor struct {
	conds map[string]*cond
}

// Comandos de inicialização
func newMonitor(variaveis ...string) *monitor {
	conds := make(map[string]*cond)
	for _, cv := range variaveis {
		conds[cv] = &cond{
			ocupado: false,
			mutex:   &sync.Mutex{},
			fila:    make([]chan bool, 0),
		}
	}
	return &monitor{
		conds: conds,
	}
}

// Demais Procedimentos

// Não sendo usado fora daqui, então está como método privado
func (m *monitor) empty(cv string) bool {
	return len(m.conds[cv].fila) == 0
}

// Havia feito uma versão com rank usando uma fila de prioridade ao invés de um slice(array).
// Mas percebi que para este problema ele não faz diferença alguma.
func (m *monitor) wait(cv string) {
	condit := m.conds[cv]
	condit.mutex.Lock()
	defer condit.mutex.Unlock()

	if condit.ocupado {
		p := make(chan bool)
		// Ele já entrar na fila dentro do lock garante que a mensagem cehga nele
		// mesmo antes dele chegar em "<-p".
		condit.fila = append(condit.fila, p)
		condit.mutex.Unlock()
		// Espera sinal
		<-p
		condit.mutex.Lock()
	}
	condit.ocupado = true
}

func (m *monitor) signal(cv string) {
	condit := m.conds[cv]
	condit.mutex.Lock()
	defer condit.mutex.Unlock()

	condit.ocupado = false
	if !m.empty(cv) {
		c := condit.fila[0]
		condit.fila = condit.fila[1:]
		c <- true
	}
}
