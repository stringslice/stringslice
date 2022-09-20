// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stringslice_test

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/pfmt/pfmt"
	"github.com/pfmt/stringslice"
)

var uniqueCopyTests = []struct {
	name  string
	line  string
	src   []string
	dst   []string
	want  []string
	bench bool
	skip  bool
	keep  bool
}{
	{
		name:  "non unique",
		line:  testline(),
		src:   []string{"foo", "bar", "foo", "baz", "foo", "qux", "foo", "quux", "foo", "garply", "foo", "waldo", "foo", "fred", "foo", "plugh", "foo", "xyzzy", "foo", "thud"},
		dst:   []string{"", "", "", "", "", "", "", "", "", "", ""},
		want:  []string{"foo", "bar", "baz", "qux", "quux", "garply", "waldo", "fred", "plugh", "xyzzy", "thud"},
		bench: true,
	}, {
		name: "already unique",
		line: testline(),
		src:  []string{"foo", "bar", "baz", "qux", "quux", "garply", "waldo", "fred", "plugh", "xyzzy", "thud"},
		dst:  []string{"", "", "", "", "", "", "", "", "", "", ""},
		want: []string{"foo", "bar", "baz", "qux", "quux", "garply", "waldo", "fred", "plugh", "xyzzy", "thud"},
	}, {
		name: "non unique",
		line: testline(),
		src:  []string{"foo", "foo", "foo", "foo", "foo", "foo", "foo", "foo", "foo", "foo", "foo"},
		dst:  []string{"", "", "", "", "", "", "", "", "", "", ""},
		want: []string{"foo"},
	}, {
		name: "without destination",
		line: testline(),
		src:  []string{"foo", "bar", "baz", "qux", "quux", "garply", "waldo", "fred", "plugh", "xyzzy", "thud"},
		dst:  nil,
		want: nil,
	}, {
		name: "empty destination",
		line: testline(),
		src:  []string{"foo", "bar", "baz", "qux", "quux", "garply", "waldo", "fred", "plugh", "xyzzy", "thud"},
		dst:  []string{},
		want: []string{},
	}, {
		name: "short destination",
		line: testline(),
		src:  []string{"foo", "bar", "baz", "qux", "quux", "garply", "waldo", "fred", "plugh", "xyzzy", "thud"},
		dst:  []string{"", ""},
		want: []string{"foo", "bar"},
	}, {
		name: "very short destination",
		line: testline(),
		src:  []string{"foo", "bar", "baz", "qux", "quux", "garply", "waldo", "fred", "plugh", "xyzzy", "thud"},
		dst:  []string{""},
		want: []string{"foo"},
	},
}

func TestUniqueCopy(t *testing.T) {
	keep := uniqueCopyTests[:0]
	skip := uniqueCopyTests[:0]
	for _, tt := range uniqueCopyTests {
		if tt.keep {
			keep = append(keep, tt)
		} else {
			skip = append(skip, tt)
		}
	}

	if len(keep) == 0 {
		keep = uniqueCopyTests
	} else {
		for _, tt := range skip {
			t.Logf("%s/unkeep: %s", tt.line, tt.name)
		}
	}

	for _, tt := range keep {
		if tt.skip {
			t.Logf("%s/skip: %s", tt.line, tt.name)
			continue
		}

		tt := tt

		t.Run(tt.line+"/"+tt.name, func(t *testing.T) {
			t.Parallel()

			n := stringslice.UniqueCopy(tt.dst, tt.src)
			got := tt.dst[:n]

			if !cmp.Equal(got, tt.want) {
				t.Errorf("\nwant: %s\n got: %s\ntest: %s", pfmt.Sprint(tt.want), got, tt.line)
			}
		})
	}
}

func BenchmarkUniqueCopy(b *testing.B) {
	b.ReportAllocs()

	keep := uniqueCopyTests[:0]
	skip := uniqueCopyTests[:0]
	for _, tt := range uniqueCopyTests {
		if tt.keep {
			keep = append(keep, tt)
		} else {
			skip = append(skip, tt)
		}
	}

	if len(keep) == 0 {
		keep = uniqueCopyTests
	} else {
		for _, tt := range skip {
			b.Logf("%s/unkeep: %s", tt.line, tt.name)
		}
	}

	for _, tt := range keep {
		if tt.skip {
			b.Logf("%s/skip: %s", tt.line, tt.name)
			continue
		}

		if !tt.bench {
			continue
		}

		b.Run(tt.line, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = stringslice.UniqueCopy(tt.dst, tt.src)
			}
		})
	}
}

func TestUniqueCopyToSelf(t *testing.T) {
	src := []string{"foo", "bar", "foo", "baz", "foo", "qux", "foo", "quux", "foo", "garply", "foo", "waldo", "foo", "fred", "foo", "plugh", "foo", "xyzzy", "foo", "thud"}
	want := []string{"foo", "bar", "baz", "qux", "quux", "garply", "waldo", "fred", "plugh", "xyzzy", "thud"}

	n := stringslice.UniqueCopy(src, src)
	got := src[:n]

	if !cmp.Equal(got, want) {
		t.Errorf("\nwant: %s\n got: %s", pfmt.Sprint(want), got)
	}
}

func testline() string {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		return fmt.Sprintf("%s:%d", filepath.Base(file), line)
	}
	return "it was not possible to recover file and line number information about function invocations"
}
