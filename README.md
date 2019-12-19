# `goyaccfmt`

`goyaccfmt` auto reformats [`goyacc`](https://godoc.org/golang.org/x/tools/cmd/goyacc) source code by calling [`gofmt`](https://golang.org/cmd/gofmt/). 

<table border=0>
<tr><td>

The following command reformats a source file `a.y` and outputs to stdout.

```bash
goyaccfmt a.y
```

For inline reformat, please use option `-w`.

```bash
goyaccfmt -w a.y
```

To the right is the difference before and after auto reformatting the grammar rule file of [SQLFlow](https://sqlflow.org/sqlflow).

</td><td>

![](opendiff-goyaccfmt.png)

</td></tr>
</table>
