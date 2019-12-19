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
	"os"
	"os/exec"
	"strings"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: goyaccfmt path\n")
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

	var out io.Writer
	var buf bytes.Buffer
	if overwrite {
		out = &buf
	} else {
		out = os.Stdout
	}

	e = goyaccfmt(in, out)
	in.Close()

	if e != nil {
		return e
	}

	if overwrite {
		f, e := os.Create(path)
		if e != nil {
			return e
		}
		defer f.Close()
	}

	return nil
}

const (
	HEAD = iota
	PREEMBLE
	TYPES
	RULES
	APPENDIX
)

func goyaccfmt(in io.Reader, out io.Writer) error {
	pr, pw := io.Pipe()

	var stderr bytes.Buffer
	cmd := exec.Command("gofmt")
	cmd.Stdin = pr
	cmd.Stdout = out
	cmd.Stderr = &stderr

	if e := cmd.Start(); e != nil {
		return fmt.Errorf("Cannot start gofmt: %v", e)
	}

	content := HEAD

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		switch l := strings.TrimSpace(scanner.Text()); l {
		case "%{", "%}", "%%":
			content++
		default:
			if content == PREEMBLE || content == APPENDIX {
				if _, e := pw.Write(scanner.Bytes()); e != nil {
					return fmt.Errorf("Copying content error: %v", e)
				}
				pw.Write([]byte("\n"))
			}
		}
	}
	pw.Close() // Signal the end of content.

	if e := scanner.Err(); e != nil {
		return fmt.Errorf("Scanner error: %v", e)
	}

	if e := cmd.Wait(); e != nil {
		return fmt.Errorf("Waiting for gofmt: %v. %s", e, stderr.String())
	}

	return nil
}
