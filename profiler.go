// Copyright © 2019 Alexandre Kovac <contact@kovacou.fr>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package database

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/mozillazg/go-slugify"
)

// newProfiler create a new profiler.
func newProfiler(output string) *profiler {
	p := &profiler{
		DirectoryOutput: output,
	}
	p.start()

	return p
}

type profiler struct {
	DirectoryOutput string
	store           chan profile
	i               uint
}

func (p *profiler) start() {
	if p.store == nil {
		p.store = make(chan profile)
	}

	go func() {
		for m := range p.store {
			flag := slugify.Slugify(m.flag)
			if flag == "" {
				flag = "default"
			}

			filename := fmt.Sprintf(
				"%s/%s/%s/%s/%d____%s.sql",
				p.DirectoryOutput,
				m.qs.Start().Format("2006_01_02"),
				flag,
				m.qs.ContextID(),
				p.i,
				m.qs.Runtime().String(),
			)

			p.write(filename, m.qs.Bytes())
			p.i++
		}
	}()
}

func (p *profiler) close() error {
	if p.store != nil {
		close(p.store)
	}
	return nil
}

func (p *profiler) write(filename string, body []byte) error {
	filename = path.Clean(filename)
	os.MkdirAll(filepath.Dir(filename), os.ModePerm)
	return ioutil.WriteFile(filename, body, os.ModePerm)
}

// Push a new profile into the profiler.
func (p *profiler) Push(qs QueryState) {
	if p.store != nil {
		p.store <- profile{
			qs:   qs,
			flag: strings.Join(qs.ContextFlag(), " "),
		}
	}
}

type profile struct {
	qs   QueryState
	flag string
}

// -------------------------------------------------

// HasProfiler says if the connection has a profiler.
func (conn *db) HasProfiler() bool {
	return conn.profiler != nil
}
