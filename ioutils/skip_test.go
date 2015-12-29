// Copyright 2015 Jonathan Boulle
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

package ioutils

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestSkippingReader(t *testing.T) {
	tests := []struct {
		in  string
		pat string

		wout string
	}{
		{
			`abcdefghijklmnopqrstuvwxyz`,
			`jkl`,
			`jklmnopqrstuvwxyz`,
		},
		{
			`<html>
<body>
<table>
btw, this is a stupid way to parse html
</table>
</body>
</html>`,
			`<table>`,

			`<table>
btw, this is a stupid way to parse html
</table>
</body>
</html>`,
		},
		{
			`<html><foo></foo><bar><bla><html><foo>`,
			`<foo>`,

			`<foo></foo><bar><bla><html><foo>`,
		},
	}
	for i, tt := range tests {
		r := bytes.NewBufferString(tt.in)
		sr := NewSkippingReader(r, []byte(tt.pat))
		b, err := ioutil.ReadAll(sr)
		if err != nil {
			t.Errorf("#%d: unexpected error: %v", i, err)
		}
		if g := string(b); g != tt.wout {
			t.Errorf("#%d: got:\n%q\nwant:\n%q", i, g, tt.wout)
		}
	}
}

func TestReadUntilReader(t *testing.T) {
	tests := []struct {
		in  string
		pat string

		wout string
	}{
		{
			`abcdefghijklmnopqrstuvwxyz`,
			`jkl`,
			`abcdefghijkl`,
		},
		{
			`<html>
<body>
<table>
btw, this is a stupid way to parse html
</table>
</body>
</html>`,
			`<table>`,

			`<html>
<body>
<table>`,
		},
		{
			`<html><bla><foo></foo><bar></bla><html><foo>`,
			`</bla>`,

			`<html><bla><foo></foo><bar></bla>`,
		},
	}
	for i, tt := range tests {
		r := bytes.NewBufferString(tt.in)
		sr := NewReadUntilReader(r, []byte(tt.pat))
		b, err := ioutil.ReadAll(sr)
		if err != nil {
			t.Errorf("#%d: unexpected error: %v", i, err)
		}
		if g := string(b); g != tt.wout {
			t.Errorf("#%d: got:\n%q\nwant:\n%q", i, g, tt.wout)
		}
	}
}
