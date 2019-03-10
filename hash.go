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
	"encoding/binary"
	"fmt"
	"hash/crc32"
)

type Key int32
type Value int32

// Useful hash collision probability calculator:
// http://davidjohnstone.net/pages/hash-collision-probability
//
// Checksum calculator: https://crccalc.com/
// hash/crc32 uses polynomial function bits in reverse order.
// E.g. for CRC-32Q, 0x814141AB in reverse bit order is 0xD5828281.

type Hashtable struct {
	mask     uint32 // used to limit hash result
	hashBits uint32 // hash table size is 2 ^ hashBits
	cuckoo   bool   // true if Cuckoo hashing is used
	occupied []bool // true is this slot contains a key
	key      []Key
	value    []Value
}

var IEEETable = crc32.MakeTable(crc32.IEEE)
var CastagnoliTable = crc32.MakeTable(crc32.Castagnoli)
var KoopmanTable = crc32.MakeTable(crc32.Koopman)
var Crc32qTable = crc32.MakeTable(0xD5828281)
var PosixTable = crc32.MakeTable(0x76dc419)

var polynomialTables = []*crc32.Table{
	IEEETable,
	CastagnoliTable,
	KoopmanTable,
	Crc32qTable,
	PosixTable,
}

func getPolynomialTables() []*crc32.Table {
	return polynomialTables
}

// Hash 'key' using polynomial 'table' to 'mask' bits.
func computeHash(key Key, table *crc32.Table, mask uint32) uint32 {
	bs := make([]byte, 4)
	binary.BigEndian.PutUint32(bs, uint32(key))
	hash := crc32.Checksum(bs, table)
	hash &= mask
	return hash
}

// Create an empty hash table of size 2 ^ hashBits.
// 'cuckoo' is true when Cuckoo hashing is used.
func newHashtable(hashBits int, cuckoo bool) (*Hashtable, error) {
	if hashBits <= 0 || hashBits > 32 {
		return nil, fmt.Errorf("hash bits must be 0 < bits <= 32")
	}
	size := 1 << uint(hashBits)
	h := &Hashtable{
		mask:     (1 << uint(hashBits)) - 1,
		hashBits: uint32(hashBits),
		cuckoo:   cuckoo,
		occupied: make([]bool, size),
		key:      make([]Key, size),
		value:    make([]Value, size),
	}
	return h, nil
}

// Make the hash table empty.
func (h *Hashtable) clear() {
	size := 1 << uint(h.hashBits)
	for i := 0; i < size; i++ {
		h.occupied[i] = false
	}
}

// Print the contents of the hash table on stdout.
func (h *Hashtable) dump() {
	size := 1 << uint(h.hashBits)
	var elements int
	for i := 0; i < size; i++ {
		if h.occupied[i] {
			fmt.Printf("(%x)%x->%x\n", i,
				uint32(h.key[i]), uint32(h.value[i]))
			elements++
		}
	}
	fmt.Printf("%d elements\n", elements)
}

// Return true and the value corresponding to 'key' when the key is
// present in the hash table. Return false otherwise.
func (h *Hashtable) search(key Key) (Value, bool) {
	if h.cuckoo {
		_, value, found := h.cuckooSearch(key)
		return value, found
	} else {
		return h.simpleSearch(key)
	}
}

// Insert or update the 'key' in the hash table. Return an
// error when there is a hash collision.
func (h *Hashtable) insert(key Key, value Value) error {
	if h.cuckoo {
		return h.cuckooInsert(key, value)
	} else {
		return h.simpleInsert(key, value)
	}
}
