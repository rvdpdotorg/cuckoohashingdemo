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
	"testing"
)

func TestHashing(t *testing.T) {
	// IEEE hash(0) = 2144df1c
	// IEEE hash(1) = 5643ef8a
	h, _ := newHashtable(8, false)

	var value Value
	var found bool

	// search key=0 in empty table
	_, found = h.search(0)
	if found {
		t.Error("Search(0) should have returned fasle")
	}

	// insert and verify key=0, value=100
	if h.insert(0, 100) != nil {
		t.Error("Insert(0, 100) should have succeeded")
	}
	value, found = h.search(0)
	if !found {
		t.Error("Search(0) should have succeeded")
	}
	if value != 100 {
		t.Error("Expected 100, got ", value)
	}

	// insert and verify key=1, value=200
	if h.insert(1, 200) != nil {
		t.Error("Insert(1, 200) should have succeeded")
	}
	value, found = h.search(1)
	if !found {
		t.Error("Search(1) should have succeeded")
	}
	if value != 200 {
		t.Error("Expected 200, got ", value)
	}

	// overwrite key=0 and verify
	if h.insert(0, 300) != nil {
		t.Error("Insert(0, 300) should have succeeded")
	}
	value, found = h.search(0)
	if !found {
		t.Error("Search(0) should have succeeded")
	}
	if value != 300 {
		t.Error("Expected 300, got ", value)
	}

	// IEEE hash(120) = 7f9a2612
	// IEEE hash(284) = 2c5eb212
	// collision for hash table size 2^8
	if h.insert(120, 120) != nil {
		t.Error("Insert(120, 120) should have succeeded")
	}
	if h.insert(284, 284) == nil {
		t.Error("Insert(284, 284) should have failed (collision)")
	}
}

func TestCuckoo(t *testing.T) {
	h, _ := newHashtable(14, true)

	if h.insert(1, 1) != nil {
		t.Error("Insert(1, 1) should have succeeded")
	}
	if h.insert(2, 2) != nil {
		t.Error("Insert(1, 1) should have succeeded")
	}
	if h.insert(3, 3) != nil {
		t.Error("Insert(3, 3) should have succeeded")
	}

	// IEEE hash(120) = 7f9a2612
	// IEEE hash(284) = 2c5eb212
	// should not collide for cuckoo hashing
	if h.insert(120, 120) != nil {
		t.Error("Insert(120, 120) should have succeeded")
	}
	if h.insert(284, 284) != nil {
		t.Error("Insert(284, 284) should have succeeded")
	}

	h.clear()
	data := make(map[Key]Value)
	for i := 0; i < 1000; i++ {
		k := rand.Uint32()
		data[Key(k)] = Value(k)
		if err := h.insert(Key(k), Value(k)); err != nil {
			fmt.Printf("%v\n", err)
			fmt.Printf("failed for %x\n", k)
			fmt.Printf("IEEE %x\n",
				computeHash(Key(k), IEEETable, h.mask))
			fmt.Printf("CastagnoliTable %x\n",
				computeHash(Key(k), CastagnoliTable, h.mask))
			fmt.Printf("KoopmanTable %x\n",
				computeHash(Key(k), KoopmanTable, h.mask))
			fmt.Printf("Crc32qTable %x\n",
				computeHash(Key(k), Crc32qTable, h.mask))
			break
		}
	}
	for k, v := range data {
		value, found := h.search(k)
		if !found {
			t.Errorf("search(%d) should have succeeded",
				uint32(k))
		}
		if value != v {
			t.Errorf("search(%d) returns %d, %d expected",
				uint32(k), uint32(value), uint32(v))
		}
	}
}
