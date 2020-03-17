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
	"flag"
	"fmt"
	"go/format"
	"io"
	"io/ioutil"
	"log"
	"os"
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

	if flag.NArg() > 0 {
		for _, f := range flag.Args() {
			if e := formatFile(f, *overwrite); e != nil {
				log.Fatal(e)
			}
		}
	} else {
		goyaccfmt(os.Stdin, os.Stdout)
	}
}

func formatFile(path string, overwrite bool) error {
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
	f.Close()
	return e
}

func goyaccfmt(in io.Reader, out io.Writer) error {
	const (
		HEAD     = iota // content before %{
		PREEMBLE        // between %{ and %}, need gofmt
		TYPES           // between %} and %%
		RULES           // bewteen the first and the second %%
		APPENDIX        // after the second %%, need gofmt
	)
	current := HEAD
	var code string

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		switch l := strings.TrimSpace(scanner.Text()); l {
		case "%{", "%}", "%%": // control lines
			current++
			switch current {
			case PREEMBLE, APPENDIX:
				code = "" // clear out for accumulation
			case TYPES:
				if e := gofmt(code, out); e != nil {
					return e
				}
			}
			fmt.Fprintf(out, "%s\n", l)

		default: // normal lines
			switch current {
			case PREEMBLE, APPENDIX:
				code += l + "\n"
			default:
				fmt.Fprintf(out, "%s\n", l)
			}
		}
	}

	if e := scanner.Err(); e != nil {
		return fmt.Errorf("Scanner error: %v", e)
	}

	return gofmt(code, out) // formatted appendix.
}

func gofmt(code string, out io.Writer) error {
	src, e := format.Source([]byte(code))
	if e != nil {
		return e
	}
	_, e = out.Write(src)
	return e
}
