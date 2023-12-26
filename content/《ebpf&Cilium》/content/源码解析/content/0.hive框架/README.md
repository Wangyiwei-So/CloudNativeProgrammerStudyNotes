# 0.hive框架

简介: cilium中很多组件都是通过hive框架来组织的。

hive是基于uber/dig的依赖注入框架

# 1. uber/dig库

## 1.1 安装

dig是uber开源的依赖注入库

```bash
go get go.uber.org/dig
```

## 1.2 hello world

```go
package main

import (
	"fmt"

	"go.uber.org/dig"
)

type Option struct { //依赖
	Name string
}

func main() {
	container := dig.New() //new一个，没啥好说的
	container.Provide(func() (*Option, error) { //单例式的初始化依赖
		return &Option{
			Name: "wangyiwei",
		}, nil
	})
	container.Invoke( //使用依赖
		func(opt *Option) { 
			fmt.Println(opt)
		},
	)
}
```

dig主要提供了两个能力，声明依赖和使用依赖

## 1.3 声明依赖

`container.Provider`函数，需要传入一个函数。

- 这个函数的返回值就是要声明的依赖项
- 可以返回多个对象
- 可以返回错误，如果错误不为nil，则不会声明依赖
- 一种类型的依赖只能声明一次，是单例的，多次声明只会生效第一次
- 声明interface的话也只能声明一次，dig是按照反射的类型来区分的，并不能区分interface背后的实际类型（但dig还可以指定名字，往下看）
- 返回多个对象有两种使用方式，一种是直接函数返回多个返回值，另一种是直接返回一个潜入dig.Out的对象，下面是一个示例

```go
package main

import (
	"flag"
	"fmt"

	"go.uber.org/dig"
)

type Option struct {
	Name string
}

type Option2 struct {
	Id string
}

type Options struct { //内嵌dig.Out就给dig说明了这个结构体包含了多个需要声明的依赖
	dig.Out
	Option1 *Option
	Option2 *Option2
}

func main() {
	useDigOut := false
	flag.BoolVar(&useDigOut, "b", false, "")
	flag.Parse()
	container := dig.New()
	if useDigOut {
		container.Provide(func() (Options, error) { //嫌返回多个参数不好看的话，可以使用这种方式
			return Options{
				Option1: &Option{
					Name: "wangyiwei",
				},
				Option2: &Option2{
					Id: "12345",
				},
			}, nil
		})
	} else {
		container.Provide(func() (*Option, *Option2, error) { //直接返回多个返回值
			return &Option{Name: "xuqinqin"}, &Option2{Id: "67890"}, nil
		})
	}

	container.Invoke(
		func(opt *Option) {
			fmt.Println(opt)
		},
	)
}
```

这样运行，查看结果

```bash
❯ go run main.go
&{xuqinqin}

❯ go run main.go -b
&{wangyiwei}
```

## 1.4 使用依赖

container.Invoke函数，需要传入一个函数，函数的参数就是要使用的依赖的类型，如果这个类型的依赖已经被声明，则dig会给这个函数传入参数

- 类型必须匹配
- 可以传入多个参数
- 传入多个参数同样有两种方式，一种是入参直接写多个，另一种是嵌套dig.In，下面是一个例子

```bash
package main

import (
	"flag"
	"fmt"

	"go.uber.org/dig"
)

type Option struct { //依赖1
	Name string
}

type Option2 struct { //依赖2
	Id string
}

type Options struct {
	dig.In
	Opt1 *Option
	Opt2 *Option2
}

func main() {
	useDigIn := false
	flag.BoolVar(&useDigIn, "b", false, "")
	flag.Parse()

	container := dig.New()
	container.Provide(func() (*Option, error) {
		return &Option{
			Name: "wangyiwei",
		}, nil
	})
	container.Provide(func() (*Option2, error) { //
		return &Option2{
			Id: "12345",
		}, nil
	})
	if useDigIn {
		container.Invoke(func(opt Options) {
			fmt.Println(opt.Opt1, opt.Opt2)
		})
	} else {
		container.Invoke( //使用依赖
			func(opt *Option, opt2 *Option2) {
				fmt.Println(opt, opt2)
			},
		)
	}
}
```

测试运行

```bash
❯ go run main.go -b
&{wangyiwei} &{12345}

❯ go run main.go   
&{wangyiwei} &{12345}
```

在container.Provider的func如果传入参数，也会被注入，并且dig会自动处理依赖的顺序。

```go
package main

import (
	"fmt"
	"go.uber.org/dig"
)

type Option struct { //依赖
	Name string
}

type Option2 struct {
	Opt *Option
}

func main() {
	container := dig.New()
	container.Provide(func(opt *Option) *Option2 { //Provider的func里也可以传入依赖
		fmt.Println("===", opt)
		return &Option2{
			Opt: opt,
		}
	})
	container.Provide(func() (*Option, error) { //dig会自动处理依赖的顺序
		return &Option{
			Name: "wangyiwei",
		}, nil
	})

	container.Invoke( //使用依赖
		func(opt *Option2) {
			fmt.Println(opt.Opt)
		},
	)
}
```

测试运行

```go
❯ go run main.go
=== &{wangyiwei}
&{wangyiwei}
```

## 1.5 指定依赖的名字

container.Provider以及container.Invoke时都可以指定依赖的名字

指定名字主要是为了能够多次声明同一种类型的依赖：

- 同一种类型的结构体，但初始化的参数不同
- 要注入interface

要注入interface，考虑下面这种情况

```bash
package main

import (
	"fmt"

	"go.uber.org/dig"
)

type Option1 struct { //依赖
	Name string
}

func (o *Option1) GetName() string {
	return o.Name
}

type Option2 struct { //依赖
	Name string
}

func (o *Option2) GetName() string {
	return o.Name
}

type Option interface {
	GetName() string
}

func main() {
	container := dig.New()
	container.Provide(func() (Option, error) {
		return &Option1{
			Name: "wangyiwei",
		}, nil
	})
	container.Provide(func() (Option, error) {
		return &Option2{
			Name: "xuqinqin",
		}, nil
	})
	container.Invoke( //使用依赖
		func(opt Option) {
			fmt.Println(opt.GetName())
		},
	)
}
```

一个interface可以被多个类型实现，但我无法区分，因为是单例注入，本例中只会生效wangyiwei这段

指定名字可以解决这个问题

```go
package main

import (
	"fmt"

	"go.uber.org/dig"
)

type Option1 struct { //依赖
	Name string
}

func (o *Option1) GetName() string {
	return o.Name
}

type Option2 struct { //依赖
	Name string
}

func (o *Option2) GetName() string {
	return o.Name
}

type Option interface {
	GetName() string
}

type Options struct {
	dig.In
	Opt1 Option `name:"option1"` //这里指定依赖的名字
	Opt2 Option `name:"option2"`
}

func main() {
	container := dig.New()
	container.Provide(func() (Option, error) {
		return &Option1{
			Name: "wangyiwei",
		}, nil
	}, dig.Name("option1")) //指定声明的依赖的名字
	container.Provide(func() (Option, error) {
		return &Option2{
			Name: "xuqinqin",
		}, nil
	}, dig.Name("option2"))
	container.Invoke( //使用依赖
		func(opts Options) {
			fmt.Println(opts.Opt1.GetName(), opts.Opt2.GetName())
		},
	)
}
```

测试运行

```go
❯ go run main.go 
wangyiwei xuqinqin
```

## 1.6 组

dig还支持注入组，组就是一些相同类型的对象的切片

```go
package main

import (
	"fmt"

	"go.uber.org/dig"
)

type Option struct { //依赖
	Name string
}

type Options struct {
	dig.In
	Options []*Option `group:"options"` //指定组
}

func main() {
	container := dig.New()
	container.Provide(func() (*Option, error) {
		return &Option{
			Name: "wangyiwei",
		}, nil
	}, dig.Group("options")) //往一个组里面插入俩元素
	container.Provide(func() (*Option, error) {
		return &Option{
			Name: "xuqinqin",
		}, nil
	}, dig.Group("options"))

	container.Invoke( //使用依赖
		func(opts Options) {
			for _, opt := range opts.Options {
				fmt.Println(opt.Name)
			}
		},
	)
}
```

测试运行

```go
❯ go run main.go 
wangyiwei
xuqinqin
```

# 2. hive
