# file=$(mktemp)
file="/tmp/a"
cat <<EOF > $file
%{
package parser
import ("fmt"
)
func f() {fmt.Printf(
"a")}
%}
%type  <eslt> sqlflow_select_stmt
%type  <tran> train_clause
%%
end_of_stmt
: ';'         {
   do_something($1)
}
;
%%
func g() {fmt.Printf("g")}
EOF

# gold=$(mktemp)
gold="/tmp/b"
cat <<EOF > $gold
%{
package parser

import (
	"fmt"
)

func f() {
	fmt.Printf(
		"a")
}
%}
%type  <eslt> sqlflow_select_stmt
%type  <tran> train_clause
%%
end_of_stmt
: ';'         {
   do_something($1)
}
;
%%
func g() { fmt.Printf("g") }
EOF

go install
test $? -eq 0 || { echo "Failed to compile"; exit 1; }

go run goyaccfmt.go -w $file
test $? -eq 0 || { echo "Failed to run with replace mode"; exit 1; }

cmp $file $gold || { echo "Unexpected output"; diff $file $gold; exit 1; }

# out=$(mktemp)
out="/tmp/c"
cat $file | go run goyaccfmt.go > $out
test $? -eq 0 || { echo "Failed to run reading stdin"; exit 1; }

cmp $out $gold || { echo "Unexpected output"; diff $out $gold; exit 1; }

file1=$(mktemp)
cp $file $file1
go run goyaccfmt.go -w $file $file1 || { echo "Failed to replace multiple files"; exit 1; }

cmp $file1 $gold || { echo "Unexpected output from replacing multiple files"; diff $file1 $gold; exit 1; }
