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
)

const (
	maxEvictions = 60 // maximum number of Cuckoo evictions
)

// Search the 'key' using all hash functions used in Cuckoo hashing.
// Each hash function hash an associated polynomial table.
func (hashTable *Hashtable) cuckooSearch(key Key) (uint32, Value, bool) {
	polynomialTables := getPolynomialTables()
	for _, table := range polynomialTables {
		hash := computeHash(key, table, hashTable.mask)
		if hashTable.occupied[hash] && hashTable.key[hash] == key {
			return hash, hashTable.value[hash], true
		}
	}
	return 0, 0, false
}

func (hashTable *Hashtable) cuckooInsert(key Key, value Value) error {
	polynomialTables := getPolynomialTables()
	for _, table := range polynomialTables {
		hash := computeHash(key, table, hashTable.mask)
		if !hashTable.occupied[hash] {
			hashTable.key[hash] = key
			hashTable.value[hash] = value
			hashTable.occupied[hash] = true
			return nil
		} else if hashTable.key[hash] == key {
			hashTable.value[hash] = value
			return nil
		}
	}

	// at this point there is a collision for all hash functions

	var hash uint32
	var oldKey Key
	var oldValue Value
	var oldTable int32
	var n int32
	round := maxEvictions
	for round > 0 {
		// Use a random hash function and make sure it is
		// different from the one used in the previous round.
		for n == oldTable {
			n = rand.Int31n(int32(len(polynomialTables)))
		}
		table := polynomialTables[n]
		oldTable = n
		hash = computeHash(key, table, hashTable.mask)
		if !hashTable.occupied[hash] {
			// Hash slot is empty, store 'key' and return
			hashTable.key[hash] = key
			hashTable.value[hash] = value
			hashTable.occupied[hash] = true
			return nil
		}
		// Store 'key' and remember evicted 'key'.
		oldKey = hashTable.key[hash]
		oldValue = hashTable.value[hash]
		hashTable.key[hash] = key
		hashTable.value[hash] = value
		key = oldKey
		value = oldValue
		round--
	}
	return fmt.Errorf("Insert: hash collision for %x, hash=%x",
		uint32(key), hash)
}
