// Copyright 2019 The SQLFlow Authors. All rights reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: goyaccfmt [-w] path\n")
	flag.PrintDefaults()
}

func main() {
	overwrite := flag.Bool("w", false, "overwrite source file instead of stdout")
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() != 1 {
		usage()
		os.Exit(-1)
	}

	if e := goyaccfmtMain(flag.Arg(0), *overwrite); e != nil {
		fmt.Fprintf(os.Stderr, "%s", e)
		os.Exit(-1)
	}
}

func goyaccfmtMain(path string, overwrite bool) error {
	in, e := os.Open(path)
	if e != nil {
		return fmt.Errorf("Cannot open input %s: %v", path, e)
	}

	out, e := ioutil.TempFile("", "")
	if e != nil {
		return fmt.Errorf("Cannot create temp file: %v", e)
	}

	e = goyaccfmt(in, out)
	if e != nil {
		return fmt.Errorf("goyaccfmt: %v", e)
	}

	if e := in.Close(); e != nil {
		return fmt.Errorf("Failed closing source file: %v", e)
	}

	if e := out.Close(); e != nil {
		return fmt.Errorf("Cannot close temp output file: %v", e)
	}

	if overwrite {
		return os.Rename(out.Name(), path)
	}
	return cat(out.Name(), os.Stdout)
}

func cat(filename string, w io.Writer) error {
	f, e := os.Open(filename)
	if e != nil {
		return e
	}

	_, e = io.Copy(w, f)
	return e
}

func goyaccfmt(in io.Reader, out io.Writer) error {
	const (
		HEAD     = iota // content before %{
		PREEMBLE        // between %{ and %}
		TYPES           // between %} and %%
		RULES           // bewteen the first and the second %%
		APPENDIX        // after the second %%
	)

	var fmtr *gofmt
	var e error
	current := HEAD

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		switch l := strings.TrimSpace(scanner.Text()); l {
		case "%{", "%}", "%%":
			current++
			switch current {
			case PREEMBLE, APPENDIX:
				fmtr, e = newGofmt(out)
				if e != nil {
					return fmt.Errorf("newGofmt: %v", e)
				}
			case TYPES:
				fmtr.Close()
			}
			fmt.Fprintf(out, "%s\n", l)

		default:
			var w io.Writer = out
			if current == PREEMBLE || current == APPENDIX {
				w = fmtr
			}
			if _, e := w.Write(scanner.Bytes()); e != nil {
				return fmt.Errorf("Copying content error: %v", e)
			}
			w.Write([]byte("\n"))
		}
	}

	if e := scanner.Err(); e != nil {
		return fmt.Errorf("Scanner error: %v", e)
	}

	return fmtr.Close()
}

type gofmt struct {
	pr     *io.PipeReader
	pw     *io.PipeWriter
	cmd    *exec.Cmd
	stderr bytes.Buffer
}

func newGofmt(out io.Writer) (*gofmt, error) {
	f := &gofmt{}
	f.pr, f.pw = io.Pipe()

	f.cmd = exec.Command("gofmt")
	f.cmd.Stdin = f.pr
	f.cmd.Stdout = out
	f.cmd.Stderr = &f.stderr

	if e := f.cmd.Start(); e != nil {
		return nil, fmt.Errorf("Cannot start gofmt: %v", e)
	}
	return f, nil
}

func (fmtr *gofmt) Write(b []byte) (int, error) {
	return fmtr.pw.Write(b)
}

func (fmtr *gofmt) Close() error {
	if e := fmtr.pw.Close(); e != nil { // Signal the end of content.
		return fmt.Errorf("Close pipe writer: %v", e)
	}

	if e := fmtr.cmd.Wait(); e != nil {
		return fmt.Errorf("Waiting for gofmt: %v. %s", e, fmtr.stderr.String())
	}

	return nil
}
