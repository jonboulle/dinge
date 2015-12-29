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
	"bufio"
	"bytes"
	"io"
)

type skippingReader struct {
	r     *bufio.Reader
	pat   []byte
	found bool
}

// NewSkippingReader wraps the given io.Reader and returns one that will
// skip everything until the pattern pat is reached.
func NewSkippingReader(r io.Reader, pat []byte) *skippingReader {
	return &skippingReader{
		r:     bufio.NewReader(r),
		pat:   pat,
		found: false,
	}
}

func (sr *skippingReader) Read(p []byte) (n int, err error) {
	if sr.found {
		return sr.r.Read(p)
	}
outer:
	for {
		for _, w := range sr.pat {
			var c byte
			c, err = sr.r.ReadByte()
			if err != nil {
				n = 0
				return
			}
			if w != c {
				continue outer
			}
		}
		sr.found = true
		// re-add the pat we just skipped
		sr.r = bufio.NewReader(io.MultiReader(bytes.NewReader(sr.pat), sr.r))
		return sr.Read(p)
	}
}

type readUntilReader struct {
	r   *bufio.Reader
	pat []byte
	fi  int
}

// NewReadUntilReader wraps the given io.Reader and returns one that, on being
// read, will not return anything until the supplied pattern pat is found.
func NewReadUntilReader(r io.Reader, pat []byte) *readUntilReader {
	return &readUntilReader{
		r:   bufio.NewReader(r),
		pat: pat,
		fi:  0,
	}
}

func (rtr *readUntilReader) Read(p []byte) (n int, err error) {
	for n < len(p) && rtr.fi < len(rtr.pat) {
		var c byte
		c, err = rtr.r.ReadByte()
		if err != nil && n > 0 {
			err = nil
			return
		} else if err != nil {
			return
		}
		if rtr.pat[rtr.fi] != c {
			rtr.fi = 0
		} else {
			rtr.fi++
		}
		if n >= len(p) {
			return
		}
		p[n] = c
		n++
		if rtr.fi == 0 {
			continue
		}
	}
	if n == 0 {
		err = io.EOF
	}
	return
}
