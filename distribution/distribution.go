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

func (d *Distribution) GetProbability(key string) float64 {
	return d.distributionNrm[key]
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
