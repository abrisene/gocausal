package random

import (
	"math/rand"
	"time"
)

type Random struct {
	generator *rand.Rand
	seed      int64
	uses      uint64
}

func New(seed, uses int64) *Random {
	r := &Random{}
	// If seed is 0, auto seed with current time
	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	source := rand.NewSource(seed)
	r.generator = rand.New(source)
	r.seed = seed

	if uses > 0 {
		r.Discard(uses)
	}

	return r
}

func (r *Random) GetSeed() int64 {
	return r.seed
}

func (r *Random) Discard(count int64) *Random {
	for i := int64(0); i < count; i++ {
		r.uses++
		r.generator.Intn(1)
	}
	return r
}

func (r *Random) Float32() float32 {
	r.uses++
	return r.generator.Float32()
}

func (r *Random) Float64() float64 {
	r.uses++
	return r.generator.Float64()
}

func (r *Random) ExpFloat64() float64 {
	r.uses++
	return r.generator.ExpFloat64()
}

func (r *Random) NormFloat64(params ...float64) float64 {
	var mean, stDev float64 = 0, 1
	if params[0] != 0 {
		mean = params[0]
	}
	if params[1] != 0 {
		stDev = params[1]
	}
	r.uses++
	return r.generator.NormFloat64()*stDev + mean
}

func (r *Random) Int() int {
	r.uses++
	return r.generator.Int()
}

func (r *Random) Intn(n int) int {
	r.uses++
	return r.generator.Intn(n)
}

func (r *Random) Intzn(n int) int {
	r.uses++
	return r.generator.Intn(n + 1)
}

func (r *Random) IntRange(min, max int) int {
	return r.Intzn(max-min) + min
}

func (r *Random) Int31() int32 {
	r.uses++
	return r.generator.Int31()
}

func (r *Random) Int31n(n int32) int32 {
	r.uses++
	return r.generator.Int31n(n)
}

func (r *Random) Int31zn(n int32) int32 {
	r.uses++
	return r.generator.Int31n(n + 1)
}

func (r *Random) Int31Range(min, max int32) int32 {
	return r.Int31zn(max-min) + min
}

func (r *Random) Int63() int64 {
	r.uses++
	return r.generator.Int63()
}

func (r *Random) Int63n(n int64) int64 {
	r.uses++
	return r.generator.Int63n(n)
}

func (r *Random) Int63zn(n int64) int64 {
	r.uses++
	return r.generator.Int63n(n + 1)
}

func (r *Random) Int63Range(min, max int64) int64 {
	return r.Int63zn(max-min) + min
}

func (r *Random) Uint32() uint32 {
	r.uses++
	return r.generator.Uint32()
}

func (r *Random) Uint64() uint64 {
	r.uses++
	return r.generator.Uint64()
}

func (r *Random) PickWeighted(set map[string]float64, order []string, mask []string) string {
	var lastValid, result string
	var sum float64
	value := r.Float64()

	// fmt.Println("__", value)

	for _, key := range order {
		v := set[key]
		sum += v
		// fmt.Println("==>", key, v, "  [%f]", sum)
		// If the mask doesn't contain the current key
		if !containsString(key, mask) {
			if sum >= value {
				result = key
				return result
			} else {
				lastValid = key
			}
		} else if sum >= value {
			result = lastValid
			return result
		}
	}
	return result
}

func containsString(value string, array []string) bool {
	for _, v := range array {
		if v == value {
			return true
		}
	}
	return false
}
