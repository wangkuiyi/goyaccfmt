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
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGoYaccFmt(t *testing.T) {
	a := assert.New(t)

	in := strings.NewReader(`
  %{
	package parser
	import 	"fmt"
func Print() { fmt.Println("Hello")
}
  %}
%type  <eslt> sqlflow_select_stmt
%type  <slct> standard_select_stmt

%%

sqlflow_select_stmt
: standard_select_stmt end_of_stmt {
	parseResult = &SQLFlowSelectStmt{
		Extended: false,
		StandardSelect: $1}
  };
  %%
func main() { Print() }
`)
	var out bytes.Buffer
	a.NoError(goyaccfmt(in, &out))
	a.Equal(`package parser

import "fmt"

func Print() {
	fmt.Println("Hello")
}
func main() { Print() }
`, out.String())
}
