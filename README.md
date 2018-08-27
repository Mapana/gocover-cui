# gocover-cui

模拟 go tool cover -html=c.out -o coverage.html 的显示效果显示到Terminal
比如docker中使用时，我无法显示html，可以尝试显示Terminal UI

    gocover-cui -cui c.out

## Install
##### go get
    go get github.com/Mapana/gocover-cui

##### git
    git https://github.com/Mapana/gocover-cui.git
    cd gocover-cui
    go install

## Plan
1. support log display [x]
2. add keyboard prompt help