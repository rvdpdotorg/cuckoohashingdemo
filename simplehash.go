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

func (hashTable *Hashtable) simpleSearch(key Key) (Value, bool) {
	bs := make([]byte, 4)
	binary.BigEndian.PutUint32(bs, uint32(key))
	hash := crc32.Checksum(bs, IEEETable)
	hash &= hashTable.mask
	if !hashTable.occupied[hash] {
		return 0, false
	}
	if hashTable.key[hash] == key {
		return hashTable.value[hash], true
	} else {
		return 0, false
	}
}

func (hashTable *Hashtable) simpleInsert(key Key, value Value) error {
	bs := make([]byte, 4)
	binary.BigEndian.PutUint32(bs, uint32(key))
	hash := crc32.Checksum(bs, IEEETable)
	hash &= hashTable.mask
	/*
	 * There are three options for the hash slot:
	 * - hash slot is empty
	 * - hash slot contains 'key'
	 * - hash slot contains different key
	 */
	if !hashTable.occupied[hash] {
		hashTable.key[hash] = key
		hashTable.value[hash] = value
		hashTable.occupied[hash] = true
		return nil
	}
	_, found := hashTable.simpleSearch(key)
	if found {
		// key is present in hash slot, update its value
		hashTable.value[hash] = value
		return nil
	}
	return fmt.Errorf("Insert: hash collision for %d", key)
}
