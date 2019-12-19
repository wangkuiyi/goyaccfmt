# `goyaccfmt`

<table>

<tr>
<td>

`goyaccfmt` auto reformats [`goyacc`](https://godoc.org/golang.org/x/tools/cmd/goyacc) source code by calling [`gofmt`](https://golang.org/cmd/gofmt/). 

</td><td>

Here is an example of the difference before and after auto reformatting the grammar file of [SQLFlow](https://sqlflow.org/sqlflow).

</td>
</tr>

<tr>
<td>

The following command reformats a source file `a.y` and outputs to stdout.

```bash
goyaccfmt a.y
```

For inline reformat, please use option `-w`.

```bash
goyaccfmt -w a.y
```

</td><td>

![](opendiff-goyaccfmt.png)

</td>
</tr>

</table>
