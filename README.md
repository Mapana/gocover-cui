# gocover-cui

模拟 go tool cover -html=c.out -o coverage.html 的显示效果显示到Terminal

比如docker中使用时，我无法显示html，可以尝试显示Terminal UI

    gocover-cui -cui c.out

## Install
##### go get
``` bash
go get github.com/Mapana/gocover-cui
```

##### git
``` bash
git https://github.com/Mapana/gocover-cui.git
cd gocover-cui
go install
```

## Example
``` bash
cd example

go test -run TestSUNLs -coverprofile=example_ls.out
go test -v -run TestSUNLs > example_ls.log
gocover-cui -cui=example/example_ls.out -log=example/example_ls.log # Can run -cui or -log separately

go test -run TestSUHs -coverprofile=example_hs.out
go test -v -run TestSUNHs > example_hs.log
gocover-cui -cui=example/example_hs.out -log=example/example_hs.log
```

## Plan
- [x] support log display
- [ ] add keyboard prompt help