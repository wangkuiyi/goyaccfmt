# `goyaccfmt`

`goyaccfmt` auto reformats [`goyacc`](https://godoc.org/golang.org/x/tools/cmd/goyacc) source code by calling [`gofmt`](https://golang.org/cmd/gofmt/).

The following command reformats a source file `a.y` and outputs to stdout.

```bash
goyaccfmt a.y
```

For inline reformat, please use option `-w`.

```bash
goyaccfmt -w a.y
```
