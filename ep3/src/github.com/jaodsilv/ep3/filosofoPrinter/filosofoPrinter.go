package printer

// Firula com nomes de filósofos de verdade que existem na wikipedia.

import (
	"fmt"
	"math/rand"
	"time"
)

// http://pt.wikipedia.org/wiki/Lista_de_fil%C3%B3sofos
var nomes = []string{"Achad Ha-am", "Uriel Acosta", "Mortimer Adler", "Theodor Adorno",
	"Giorgio Agamben", "Agostinho de Hipona", "Agripa", "Heinrich Cornelius Agrippa",
	"Alberto Magno", "Alcmeão de Crotona", "Alexandre de Afrodísias", "Alexandre de Hales",
	"Alighieri, Dante", "Anaxágoras", "Anaximandro", "Anaxímenes de Mileto", "Alan Ross Anderson",
	"John Anderson", "G.E.M. Anscombe", "Anselmo de Cantuária", "Antístenes", "Hannah Arendt",
	"Aristarco de Samos", "Aristóteles", "Aristóxenes", "Arquimedes", "Arquitas de Tarento",
	"Marco Aurélio", "John Austin", "Avicena", "Alfred Jules Ayer", "Gaston Bachelard",
	"Jean Baudrillard", "Francis Bacon", "Roger Bacon", "Mikhail Bakhtin", "Roland Barthes",
	"Zygmunt Bauman", "Simone de Beauvoir", "Cesare Beccaria", "Walter Benjamin", "Jeremy Bentham",
	"Bhaskara I", "Henri Bergson", "George Berkeley", "Max Black", "Norberto Bobbio",
	"de La Étienne Boétie", "Leonardo Boff", "Paul Boghossian", "F. H. Bradley", "Bertold Brecht",
	"Franz Brentano", "Raimundo de Farias Brito", "Giordano Bruno", "Mário Bunge", "Tyler Burge",
	"Edmund Burke", "Albert Camus", "Fritjof Capra", "Rudolf Carnap", "Cornelius Castoriadis",
	"Mario Sergio Cortella", "Roderick Chisholm", "Noam Chomsky", "Cícero", "Comte-André Sponville",
	"Auguste Comte", "Condillac", "Confúcio", "Edmundo Curvelo", "Donald Davidson",
	"Daniel Dennett", "Gilles Deleuze", "Demócrito", "Jacques Derrida", "René Descartes",
	"Léon Denis", "John Dewey", "Denis Diderot", "Diógenes de Apolônia", "Enrique Dussel",
	"Fred Dretske", "Ronald Dworkin", "Eckhart", "Umberto Eco", "Arthur Stanley Eddington",
	"Jonathan Edwards", "Albert Einstein", "George Eliot", "Thomas Stearns Eliot",
	"Friedrich Engels", "Empédocles", "Epiteto", "Epicuro", "Erasmo de Roterdão", "Eratóstenes",
	"Euclides", "Luc Ferry", "Ludwig Feuerbach", "Paul Feyerabend", "Marsílio Ficino",
	"Johann Gottlieb Fichte", "Filolau de Crotona", "Filopono de João Alexandria", "Jerry Fodor",
	"Michel Foucault", "Gottlob Frege", "Paulo Freire", "Hans-Georg Gadamer",
	"Sidarta (Buda) Gautama", "Ernest Gellner", "Edmund Gettier", "Giacóia Oswaldo Junior",
	"Giordano Bruno", "Jahann Wolfgang von Goethe", "Nelson Goodman", "Górgias de Leontini",
	"Antonio Gramsci", "Robert Grosseteste", "Félix Guattari", "Martial Gueroult", "Jürgen Habermas",
	"Georg Hegel", "Martin Heidegger", "Werner Heisenberg", "Carl G. Hempel", "Heráclito de Éfeso",
	"Heráclides do Ponto", "Johann Gottfried von Herder", "Hiérocles", "David Hilbert",
	"Hildegarda de Bingen", "Hipátia", "Hípias de Elis", "Hipócrates", "Thomas Hobbes",
	"Douglas Hofstadter", "Barão d'Holbach", "Max Horkheimer", "David Hume", "Edmund Husserl",
	"Francis Hutcheson", "Thomas Henry Huxley", "Christiaan Huygens", "Jean Hyppolite",
	"William Ralph Inge", "Isócrates", "Isaac de Stella", "Viacheslav Ivanovich Ivanov",
	"William James", "Janine Renato Ribeiro", "Karl Jaspers", "Georg Jellinek", "Joaquim de Fiore",
	"João de Damasco", "João de Salisbury", "Hans Jonas", "Carl Gustav Jung", "Ernst Jünger",
	"Justiniano I", "Immanuel Kant", "Moritz Kaposi", "Hans Kelsen", "Søren Kierkegaard",
	"Karl Korsch", "Alfred Korzybski", "Jiddu Krishnamurti", "Saul Kripke", "Piotr Kropotkin",
	"Thomas Kuhn", "Lao-Tsé", "Gottfried Wilhelm Leibniz", "Léon Denis", "Leucipo",
	"Levi-Claude Strauss", "Emmanuel Levinas", "Pierre Lévy", "David Lewis", "Gilles Lipovetsky",
	"John Locke", "Lucrécio", "George Lukács", "Jean-François Lyotard", "Alasdair MacIntyre",
	"Moisés Maimônides", "Nicholas Malebranche", "Nicolau Maquiavel", "Herbert Marcuse",
	"Julián Marías", "José Marinho", "Jacques Maritain", "Karl Marx", "Melisso de Samos",
	"Jean Meslier", "Merleau-Maurice Ponty", "James Mill", "John Stuart Mill",
	"Michel de Montaigne", "Montesquieu", "Edgar Morin", "Viviane Mosé", "G. E. Moore",
	"Emmanuel Mounier", "Arne Naess", "Nagarjuna", "Thomas Nagel", "Nemésio de Emesa",
	"John Henry Newman", "Nicolau de Cusa", "Friedrich Nietzsche", "Robert Nozick",
	"Benedito Nunes", "Nicolau de Cusa", "Nitiren Daishonin", "William de Ockham", "Lorenz Oken",
	"o Jovem Olimpiodoro", "Michel Onfray", "Hans Christian Ørsted", "Ortega y José Gasset",
	"Wilhelm Ostwald", "Rudolf Otto", "Alvin Plantinga", "Raimon Panikkar", "Parmênides",
	"George Pappas", "Blaise Pascal", "Charles Sanders Peirce", "Pitágoras", "Platão", "Plotino",
	"Luis Felipe Pondé", "Karl Popper", "Protágoras", "Samuel Pufendorf", "Hilary Putnam",
	"Ptolomeu", "António Quadros", "Willard van Orman Quine", "François Rabelais", "Ayn Rand",
	"John Rawls", "Miguel Reale", "Thomas Reid", "Paul Ricoeur", "Huberto Rohden",
	"Renato Janine Ribeiro", "Richard Rorty", "Jean Jacques Rousseau", "Bertrand Russell",
	"Michael Sandel", "Delfim Santos", "Mário Ferreira dos Santos", "Jean-Paul Sartre",
	"Fernando Savater", "Lúcio Aneu Sêneca", "Francis Schaeffer", "Arthur Schopenhauer",
	"Roger Vernon Scruton", "Sócrates", "Sófocles", "Eduardo Abranches de Soveral",
	"Baruch Spinoza", "Dugald Stewart", "Max Stirner", "Peter Strawson", "Francisco Suárez",
	"Tales de Mileto", "Alfred Tarski", "Johann Tauler", "Charles Margrave Taylor",
	"William Temple", "Temístio", "Teodoro de Cirene", "Teofrasto", "Tertuliano", "Timeu de Lócrida",
	"Tímon", "Liev Tolstói", "Tomás de Aquino", "Alexis de Tocqueville", "Trasímaco",
	"Ernst Troeltsch", "Leon Trótski", "Tucídides", "Michael Tye", "Pietro Ubaldi",
	"Miguel de Unamuno", "Peter Unger", "Hans Vaihinger", "Lorenzo Valla", "Francisco Varela",
	"Emer de Vattel", "Gianni Vattimo", "Vauvenargues", "Adolfo Sánchez Vázquez", "Giambattista Vico",
	"Eric Voegelin", "Voltaire", "Simone Weil", "Ken Wilber", "Ludwig Wittgenstein",
	"Georg wenrik von Wright", "Christian Wolff", "Wilhelm Wundt", "Xenócrates",
	"Xenófanes de Cólofon", "Xenofonte", "Yang Zhu", "Francis Parker Yockey", "Iris Marion Young",
	"Eduard Zeller", "Ernst Zermelo", "Zenão de Cítio", "Zenão de Eleia", "Zenão de Sídon",
	"Zenão de Tarso"}

// Printer imprime as coisas
type Printer struct {
	nomes []string
}

// InitNomes é uma firula para que todos os nomes tenham chance de aparecer
func (p *Printer) InitNomes(n int32) {
	rand.Seed(time.Now().UnixNano())

	numNomes := len(nomes)
	if n > int32(numNomes) {
		n = int32(numNomes)
	}

	nomes2 := make([]string, numNomes)
	r := rand.Perm(numNomes)

	for i := int32(0); i < n; i++ {
		nomes2[i] = nomes[r[i]]
	}
	// fmt.Println(nomes2)
	p.nomes = nomes2
}

func (p *Printer) nome(index int32) string {
	if index < 318 {
		return p.nomes[index]
	}
	return fmt.Sprintf("Filosofo %d", index)
}

func (p *Printer) printAction(index int32, action string) {
	fmt.Printf("%s %s a comer as %s\n", p.nome(index), action, time.Now().String())
}

// PrintInicio imprime o momento que o filosofo começou a comer
func (p *Printer) PrintInicio(index int32) {
	p.printAction(index, "começou")
}

// PrintFim imprime o momento que o filosofo terminou de comer
func (p *Printer) PrintFim(index int32) {
	p.printAction(index, "terminou")
}

// PrintFilosofo imprime a informação atual de um filosofo
func (p *Printer) PrintFilosofo(index, peso, comidas int32) {
	fmt.Printf("%s pesa %d e comeu %d porções\n", p.nome(index), peso, comidas)
}

// FilosofoCSV devolve uma string em formato CSV para facilitar o serviço de análise dos experimentos
func FilosofoCSV(index, peso, comidas int32) string {
	return fmt.Sprintf("%d,%d,%d\n", index, peso, comidas)
}
