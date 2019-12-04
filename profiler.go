// Copyright © 2019 Alexandre Kovac <contact@kovacou.com>.
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
	return &profiler{
		DirectoryOutput: output,
	}
}

// profiler is a feature that export computed string into files.
type profiler struct {
	DirectoryOutput string
	i               uint
}

// write content into filename and create directories recursively.
func (p *profiler) write(filename string, body []byte) error {
	filename = path.Clean(filename)
	os.MkdirAll(filepath.Dir(filename), os.ModePerm)
	return ioutil.WriteFile(filename, body, os.ModePerm)
}

// Push a new profile into the profiler.
func (p *profiler) Push(qs QueryState) {
	flag := slugify.Slugify(strings.Join(qs.ContextFlag(), " "))

	if flag == "" {
		flag = "default"
	}

	filename := fmt.Sprintf(
		"%s/%s/%s/%s/%d____%s.sql",
		p.DirectoryOutput,
		qs.Start().Format("2006_01_02"),
		flag,
		qs.ContextID(),
		p.i,
		qs.Runtime().String(),
	)

	p.write(filename, qs.Bytes())
	p.i++
}

// HasProfiler says if the connection has a profiler.
func (conn *db) HasProfiler() bool {
	return conn.profiler != nil
}
