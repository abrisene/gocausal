package distribution

import (
	"sort"

	"github.com/abrisene/gocausal/normalize"
	"github.com/abrisene/gocausal/random"
)

type Distribution struct {
	keys            []string
	distribution    map[string]float64
	distributionNrm map[string]float64
	generator       *random.Random
}

func New(distribution map[string]float64, generator *random.Random) *Distribution {
	d := &Distribution{
		distribution: distribution,
		generator:    generator,
	}
	d.Regenerate()
	return d
}

func (d *Distribution) SetGenerator(generator *random.Random) *Distribution {
	d.generator = generator
	return d
}

func (d *Distribution) Regenerate() *Distribution {
	d.keys = sortKeys(d.distribution)
	d.distributionNrm = normalize.NormalizeMap(d.distribution)
	return d
}

func (d *Distribution) Pick(mask []string) string {
	return d.generator.PickWeighted(d.distributionNrm, d.keys, mask)
}

func (d *Distribution) Add(key string, value float64) *Distribution {
	d.distribution[key] += value
	if d.distribution[key] < 0 {
		delete(d.distribution, key)
	}
	d.Regenerate()
	return d
}

func sortKeys(collection map[string]float64) []string {
	keys := []string{}
	for k := range collection {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// MAIN =====================
/* func main() {
	fmt.Println("\n====\n")
	r := random.New(1, 0)
	d := map[string]float64{
		"a": 5,
		"b": 2,
		"c": 8,
		"d": 8,
		"e": 1,
	}
	w := New(d, r)
	w.Add("ab", float64(-12))
	w.Add("b", float64(8))
	fmt.Println(w)
	fmt.Println("\n")
	// fmt.Println(r.GetSeed())
	fmt.Println("\n")
	for i := 0; i < 5; i++ {
		fmt.Println(w.Pick([]string{}))
	}
	fmt.Println("\n====\n")
}
*/
