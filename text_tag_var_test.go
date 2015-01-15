// Copyright 2014, Kevin Ko <kevin@faveset.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package htmlgen

import (
	"bytes"
	"fmt"
	"testing"
)

type ParseVarsTest struct {
	s        string
	expected []Var
}

type Vars []Var

func (vars Vars) Equals(o Vars) bool {
	for ii, v := range vars {
		if v != o[ii] {
			return false
		}
	}

	return true
}

func TestParseVars(t *testing.T) {
	tests := []ParseVarsTest{
		{"", []Var{}},
		{"hello", []Var{}},
		{"01234 $worl", []Var{
			{6, 5},
		}},
		{"01234 $worl ", []Var{
			{6, 5},
		}},
		{"01234 $worl foo", []Var{
			{6, 5},
		}},
		{"01234 $worl $ $$ $", []Var{
			{6, 5},
			{12, 1},
			{14, 2},
			{17, 1},
		}},
	}

	for _, test := range tests {
		if vars := parseVars(test.s); !Vars(vars).Equals(test.expected) {
			t.Error(fmt.Sprintf("%v => mismatch %v != %v expected", test.s, vars, test.expected))
		}
	}
}

type ExpandTest struct {
	s        string
	env      Environment
	expected string
	isUnsafe bool
}

func TestExpand(t *testing.T) {
	tests := []ExpandTest{
		{s: "", env: Environment{}, expected: ""},
		{s: "",
			env: Environment{
				"name": StringValue("foo"),
			},
			expected: ""},
		{s: "hello, $name",
			env: Environment{
				"name": StringValue("foo"),
			},
			expected: "hello, foo"},
		{s: "hello, $name.",
			env: Environment{
				"name": StringValue("foo"),
			},
			expected: "hello, foo."},
		{s: "hello, $",
			env: Environment{
				"name": StringValue("foo"),
			},
			expected: "hello, undefined"},
		{s: "hello, $.",
			env: Environment{
				"name": StringValue("foo"),
			},
			expected: "hello, undefined."},
		{s: "hello, $name.  How much is $$5 for $bar?",
			env: Environment{
				"name": StringValue("foo"),
				"bar":  StringValue("apple"),
			},
			expected: "hello, foo.  How much is $5 for apple?"},
		{s: "$$$bar?",
			env: Environment{
				"name": StringValue("foo"),
				"bar":  StringValue("apple"),
			},
			expected: "$apple?"},
		{s: "$0bar",
			env: Environment{
				"name": StringValue("foo"),
				"bar":  StringValue("apple"),
				"0bar": StringValue("not possible"),
			},
			expected: "undefined0bar"},
		{s: "$bar0",
			env: Environment{
				"name": StringValue("foo"),
				"bar":  StringValue("apple"),
				"bar0": StringValue("orange"),
			},
			expected: "orange"},
		{s: "hello, $name & $name.  Is 3 < 5?",
			env: Environment{
				"name": StringValue("foo"),
			},
			expected: "hello, foo & foo.  Is 3 < 5?",
			isUnsafe: true},
		{s: "hello, $name & $name.  Is 3 < 5?",
			env: Environment{
				"name": StringValue("foo&bar"),
			},
			expected: "hello, foo&amp;bar &amp; foo&amp;bar.  Is 3 &lt; 5?"},
	}

	for _, test := range tests {
		tag := NewTextTagVar(test.s)

		if test.isUnsafe {
			tag.SetUnsafe(true)
		}

		buf := bytes.Buffer{}
		n, err := tag.write(&buf, test.env)
		if err != nil {
			t.Error(err)
		}

		if n != buf.Len() {
			t.Error(fmt.Sprintf("%v (%v) => buf.Len() = %d != %d expected", test.s, test.env, buf.Len(), n))
		}

		cmp := buf.String()
		if cmp != test.expected {
			t.Error(fmt.Sprintf("%v (%v) => mismatch %v != %v expected", test.s, test.env, cmp, test.expected))
		}
	}
}
