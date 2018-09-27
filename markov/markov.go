package markov

import (
	"strings"

	"github.com/abrisene/gocausal/distribution"
	"github.com/abrisene/gocausal/random"
	// . "github.com/logrusorgru/aurora"
)

type Model struct {
	config    *Config
	sequences [][]string
	states    map[string]int
	grams     map[string]*gram
	generator *random.Random
}

type Config struct {
	maxOrder       int
	delimiter      string
	startDelimiter string
	endDelimiter   string
}

type gram struct {
	ID        string
	Last      distribution.Distribution
	Next      distribution.Distribution
	Order     int
	Frequency int
	DegreeIn  int
	DegreeOut int
}

func New(sequences [][]string, order int, generator *random.Random) *Model {
	config := &Config{
		maxOrder:       order,
		delimiter:      "⏐",
		startDelimiter: "○",
		endDelimiter:   "◍",
	}

	model := &Model{
		config:    config,
		sequences: sequences,
		states:    make(map[string]int),
		grams:     make(map[string]*gram),
		generator: generator,
	}
	model.Add(sequences)
	return model
}

func (m *Model) getGramID(array []string) string {
	id := strings.Join(array, m.config.delimiter)
	return id
}

func (m *Model) getGram(array []string) *gram {
	id := m.getGramID(array)
	return m.grams[id]
}

func (m *Model) GetGrams() map[string]*gram {
	return m.grams
}

func (m *Model) Add(sequences [][]string) *Model {
	for _, sequence := range sequences {
		m.AddSequence(sequence)
	}
	return m
}

func (m *Model) AddSequence(sequence []string) *Model {
	m.sequences = append(m.sequences, sequence)

	// Add delimiters
	seq := append([]string{m.config.startDelimiter}, sequence...)
	seq = append(seq, m.config.endDelimiter)

	for order := 1; order <= m.config.maxOrder; order++ {
		for pos := 0; pos < len(seq); pos++ {
			nextPos := pos + order
			lastPos := pos - 1

			if nextPos > len(seq)-1 {
				break
			}

			// Find the previous and next states
			var lastState, nextState string
			if lastPos >= 0 {
				lastState = seq[lastPos]
			}
			if nextPos < len(seq) {
				nextState = seq[nextPos]
			}

			// fmt.Println(order, seq, pos, nextPos, len(seq))
			// Get the gram sequence and id
			gramSeq := seq[pos:nextPos]
			gramID := m.getGramID(gramSeq)

			// fmt.Println(Red(lastState), Green(gramID), Cyan(nextState))

			// Add the gram & edge
			m.addEdge(gramID, lastState, nextState, order)

			// Break if we've hit the end delimiter
			if nextState == m.config.endDelimiter {
				break
			}
		}
	}

	return m
}

func (m *Model) addEdge(id string, last string, next string, order int) *gram {
	// Add the gram if it doesn't exist
	g := m.grams[id]
	if g == nil {
		g = &gram{
			ID:    id,
			Last:  *distribution.New(make(map[string]float64), m.generator),
			Next:  *distribution.New(make(map[string]float64), m.generator),
			Order: order,
		}
		m.grams[id] = g
	}

	// Add the edges to the matrix
	if last != "" {
		g.DegreeIn++
		g.Last.Add(last, 1)
		g.Frequency++
	}
	if next != "" {
		g.DegreeOut++
		g.Next.Add(next, 1)
		g.Frequency++
	}
	return g
}

func (m *Model) Next(gramSeq []string, mask []string) string {
	gram := m.getGram(gramSeq)
	var pick string
	if gram != nil {
		pick = gram.Next.Pick(mask)
	}
	return pick
}

func (m *Model) Last(gramSeq []string, mask []string) string {
	gram := m.getGram(gramSeq)
	var pick string
	if gram != nil {
		pick = gram.Last.Pick(mask)
	}
	return pick
}

func (m *Model) Generate(start []string, order int, min int, max int, strict bool) []string {
	sequence := append([]string{m.config.startDelimiter}, start...)
	currentOrder := len(sequence)
	// if strict {
	// 	currentOrder = order
	// }

	var nextState string

	for i := 0; i < max; i++ {
		if i < max-1 {
			var gram *gram
			// Set the mask if applicable
			mask := []string{}
			if i < min {
				mask = append(mask, m.config.endDelimiter)
			}

			// Find a gram of the highest possible order
			for o := currentOrder; o > 0; o-- {
				gramSeq := sequence[len(sequence)-o : len(sequence)]
				gram = m.getGram(gramSeq)
				// fmt.Println(gram != nil, i < min, gram.ID, gram.Order, currentOrder, order, gram.Next.GetProbability(m.config.endDelimiter))
				if gram != nil {
					// fmt.Println(gram.ID, gram.Order, "/", order)
					// fmt.Println(gram != nil, i < min, gram.ID, gram.Order, currentOrder, order, gram.Next.GetProbability(m.config.endDelimiter))
					if i > min || gram.Next.GetProbability(m.config.endDelimiter) < 1 {
						break
					}
				}
			}

			// fmt.Println(gram.ID, gram.Order)

			// Add our next state to the sequence
			nextState = gram.Next.Pick(mask)
			sequence = append(sequence, nextState)

			// Break if our next state is the end delimiter
			if nextState == m.config.endDelimiter {
				break
			}

			// Adjust the order if dynamic
			if currentOrder < order {
				currentOrder++
			} else if currentOrder == order && order > 1 && !strict {
				currentOrder--
			}
		} else {
			// Set final state to end delimiter and break
			nextState = m.config.endDelimiter
			sequence = append(sequence, nextState)
			break
		}
	}
	return sequence[1 : len(sequence)-1]
}
