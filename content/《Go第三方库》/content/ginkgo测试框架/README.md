# ginkgo测试框架

ginkgo是一个go语言测试框架，通常会合gomega——匹配库，来结合使用。

[ginkgo](https://github.com/onsi/ginkgo)

[gomega](https://github.com/onsi/gomega)

使用时需要
```shell
go install github.com/onsi/ginkgo/v2/ginkgo #安装命令行
go get github.com/onsi/gomega
go get github.com/onsi/ginkgo
```

直接看example

```shell
#生成*_suite_test.go 
#cd 到包文件夹
glinkgo bootstrap
#生成*_test.go
#cd 到包文件夹
glinkgo generate <.go文件名>
#然后编写测试例以及配置
glinkgo #运行，等同于go test
```

*_suite_test.go是ginkgo框架的一个包的唯一入口点，只会有一个TestX函数
somelab展示了基础用法

report展示了如何输出junit格式的，这个报告可以被gitlab或jenkins在做CI时解析
