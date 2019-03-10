// Copyright [2018-present] Ronald van der Pol <Ronald.vanderPol@rvdp.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	simulations = 10000 // .01 accuracy for 95% confidence interval
)

// Run Monte Carlo simulations. Take a key 'k' out of the key space
// [0, 2^32) and try to insert it into the hash table. Repeat this
// 'simulations' times and count the number of occurences when
// the insert fails with a hash collision.
func simulate(hashBits int, cuckoo bool) map[int]float32 {
	data := make(map[int]float32)
	rand.Seed(time.Now().UnixNano())
	hashSize := 1 << uint(hashBits)
	h, _ := newHashtable(hashBits, cuckoo)

	step := 1
	// These step parameters work well to limit nr of measurements.
	if cuckoo {
		if hashBits >= 8 {
			step = step << (uint(hashBits) - 5)
		}
	} else {
		step = (1<<uint(hashBits))/500 + 1
	}
	smallSteps := false
	fmt.Printf("Number of elements: ")
	for elements := 2; elements <= hashSize; elements += step {
		fmt.Printf("[%d] ", elements)
		collisions := 0
		for n := 0; n < simulations; n++ {
			h.clear()
			collision := false
			for i := 0; i < elements; i++ {
				k := rand.Uint32()
				if h.insert(Key(k), Value(k)) != nil {
					collision = true
				}
			}
			if collision {
				collisions++
			}
		}
		probability := float32(collisions) / simulations
		data[elements] = probability
		if probability > 0.99 {
			break
		}
		// Take smaller steps when the probability starts to increase.
		if !smallSteps && probability > 0.05 {
			step = step/5 + 1
			smallSteps = true
		}
	}
	fmt.Printf("\n")
	return data
}
