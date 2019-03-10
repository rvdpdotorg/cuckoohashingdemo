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
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
)

func main() {
	input := bufio.NewScanner(os.Stdin)

	hashBits := 1
	done := false
	for !done {
		fmt.Print("Number of hash bits (hash table size 2^bits): ")
		input.Scan()
		hashBits, _ = strconv.Atoi(input.Text())
		if hashBits <= 0 || hashBits > 32 {
			fmt.Fprintf(os.Stderr,
				"Please enter integer in interval (0, 32]\n")
		} else {
			done = true
		}
	}

	var cuckoo bool
	done = false
	for !done {
		fmt.Print("Use Cuckoo hashing (y/n)? ")
		input.Scan()
		s := input.Text()
		switch s[:1] {
		case "y", "Y":
			cuckoo = true
			done = true
		case "n", "N":
			cuckoo = false
			done = true
		default:
			fmt.Fprintf(os.Stderr,
				"Please enter 'y' or 'n'\n")
			continue
		}
	}

	// Run Monte Carlo simulation.
	data := simulate(hashBits, cuckoo)

	// Print results on stdout.
	var points []int
	for point := range data {
		points = append(points, point)
	}
	sort.Ints(points)
	for _, point := range points {
		fmt.Printf("%d %f\n", point, data[point])
	}
}
