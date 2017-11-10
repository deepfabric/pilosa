// Copyright 2017 Pilosa Corp.
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

package pilosa_test

import (
	"reflect"
	"testing"

	"github.com/deepfabric/pilosa/test"
)

// Ensure database can set and retrieve column attributes.
func TestAttrStore_Attrs(t *testing.T) {
	s := test.MustOpenAttrStore()
	defer s.Close()

	// Set attributes.
	if err := s.SetAttrs(1, map[string]interface{}{"A": 100, "C": -27}); err != nil {
		t.Fatal(err)
	} else if err := s.SetAttrs(2, map[string]interface{}{"A": uint64(200)}); err != nil {
		t.Fatal(err)
	} else if err := s.SetAttrs(1, map[string]interface{}{"B": "VALUE"}); err != nil {
		t.Fatal(err)
	}

	// Retrieve attributes for column #1.
	if m, err := s.Attrs(1); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(m, map[string]interface{}{"A": int64(100), "B": "VALUE", "C": int64(-27)}) {
		t.Fatalf("unexpected attrs(1): %#v", m)
	}

	// Retrieve attributes for column #2.
	if m, err := s.Attrs(2); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(m, map[string]interface{}{"A": int64(200)}) {
		t.Fatalf("unexpected attrs(2): %#v", m)
	}
}

// Ensure database returns a non-nil empty map if unset.
func TestAttrStore_Attrs_Empty(t *testing.T) {
	s := test.MustOpenAttrStore()
	defer s.Close()

	if m, err := s.Attrs(100); err != nil {
		t.Fatal(err)
	} else if m == nil || len(m) > 0 {
		t.Fatalf("unexpected attrs: %#v", m)
	}
}

// Ensure database can unset attributes if explicitly set to nil.
func TestAttrStore_Attrs_Unset(t *testing.T) {
	s := test.MustOpenAttrStore()
	defer s.Close()

	// Set attributes.
	if err := s.SetAttrs(1, map[string]interface{}{"A": "X", "B": "Y"}); err != nil {
		t.Fatal(err)
	} else if err := s.SetAttrs(1, map[string]interface{}{"B": nil}); err != nil {
		t.Fatal(err)
	}

	// Verify attributes.
	if m, err := s.Attrs(1); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(m, map[string]interface{}{"A": "X"}) {
		t.Fatalf("unexpected attrs: %#v", m)
	}
}

// Ensure attribute block checksums can be returned.
func TestAttrStore_Blocks(t *testing.T) {
	s := test.MustOpenAttrStore()
	defer s.Close()

	// Set attributes.
	if err := s.SetAttrs(1, map[string]interface{}{"A": uint64(100)}); err != nil {
		t.Fatal(err)
	} else if err := s.SetAttrs(2, map[string]interface{}{"A": uint64(200)}); err != nil {
		t.Fatal(err)
	} else if err := s.SetAttrs(100, map[string]interface{}{"B": "VALUE"}); err != nil {
		t.Fatal(err)
	} else if err := s.SetAttrs(350, map[string]interface{}{"C": "FOO"}); err != nil {
		t.Fatal(err)
	}

	// Retrieve blocks.
	blks0, err := s.Blocks()
	if err != nil {
		t.Fatal(err)
	} else if len(blks0) != 3 || blks0[0].ID != 0 || blks0[1].ID != 1 || blks0[2].ID != 3 {
		t.Fatalf("unexpected blocks: %#v", blks0)
	}

	// Change second block.
	if err := s.SetAttrs(100, map[string]interface{}{"X": 12}); err != nil {
		t.Fatal(err)
	}

	// Ensure second block changed.
	blks1, err := s.Blocks()
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(blks0[0], blks1[0]) {
		t.Fatalf("block 0 mismatch: %#v != %#v", blks0[0], blks1[0])
	} else if reflect.DeepEqual(blks0[1], blks1[1]) {
		t.Fatalf("block 1 match: %#v ", blks0[0])
	} else if !reflect.DeepEqual(blks0[2], blks1[2]) {
		t.Fatalf("block 2 mismatch: %#v != %#v", blks0[2], blks1[2])
	}
}
