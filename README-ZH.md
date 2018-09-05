# gocover-cui

参考`go tool cover -html=c.out -o coverage.html `将结果显示在终端
```
gocover-cui -cui c.out
```

## 须知
**终端主题或bash颜色将影响最终的`gocover-cui`显示颜色**

## 安装
#### go get
``` bash
go get github.com/Mapana/gocover-cui
```

#### git
``` bash
git clone https://github.com/Mapana/gocover-cui.git
cd gocover-cui
go install
```

## 按键帮助
这里提供了分支和手册
- [branch](https://github.com/Mapana/gocover-cui/tree/key-help)
- [manuals](https://github.com/Mapana/gocover-cui/wiki)

## 示例
``` bash
cd $GOPATH/src/github.com/Mapana/gocover-cui
gocover-cui -cui=example/example_ls.out -log=example/example_ls.log # 可以单独运行 -cui 或 -log
gocover-cui -cui=example/example_hs.out -log=example/example_hs.log
```

#### 当焦点处于`Cover Files`
![image](https://github.com/Mapana/public/blob/master/gocover-cui-1.png)

#### 展开`Cover Files`选项
![image](https://github.com/Mapana/public/blob/master/gocover-cui-2.png)

#### 当焦点处于`Data View`
![image](https://github.com/Mapana/public/blob/master/gocover-cui-3.png)

## 计划
- [x] 支持log显示
- [x] 增加按键帮助