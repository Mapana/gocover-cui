# gocover-cui

Reference go tool cover -html=c.out -o coverage.html display it to the terminal

for example when using in docker, I can't display html, I can try to display the terminal UI.

    gocover-cui -cui c.out

## Install
##### go get
``` bash
go get github.com/Mapana/gocover-cui
```

##### git
``` bash
git clone https://github.com/Mapana/gocover-cui.git
cd gocover-cui
go install
```

## Example
``` bash
cd $GOPATH/src/github.com/Mapana/gocover-cui

gocover-cui -cui=example/example_ls.out -log=example/example_ls.log # Can run -cui or -log separately

gocover-cui -cui=example/example_hs.out -log=example/example_hs.log
```

#### Focus in `Cover Files`
![image](https://github.com/Mapana/gocover-cui/blob/master/gocover-cui-1.png)

#### toggle option for `Cover Files`
![image](https://github.com/Mapana/gocover-cui/blob/master/gocover-cui-2.png)

#### Focus in `Data View`
![image](https://github.com/Mapana/gocover-cui/blob/master/gocover-cui-3.png)

## Plan
- [x] support log display
- [ ] add keyboard prompt help