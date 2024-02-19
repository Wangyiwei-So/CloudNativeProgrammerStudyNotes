# color终端颜色库

简介: 有时需要对终端输出进行上色，可以使用github.com/fatih/color

## 使用
直接上示例

```go
import "github.com/fatih/color"

var BlueBold = color.New(color.FgBlue).Add(color.Bold)
var GreenBold = color.New(color.FgGreen).Add(color.Bold)
var RedBold = color.New(color.FgRed).Add(color.Bold)
var YellowBold = color.New(color.FgYellow).Add(color.Bold)
var WhiteBold = color.New(color.FgWhite).Add(color.Bold)
var CyanBold = color.New(color.FgHiCyan).Add(color.Bold)

var Blue = color.New(color.FgBlue)
var Green = color.New(color.FgGreen)
var Red = color.New(color.FgRed)
var Yellow = color.New(color.FgYellow)
var White = color.New(color.FgWhite)
var Cyan = color.New(color.FgHiCyan)

func HeaderLn(a ...interface{}) {
	b := append([]interface{}{"[INFO]"}, a...)
	_, _ = CyanBold.Println(b...)
}

func SubHeaderLn(a ...interface{}) {
	b := append([]interface{}{"[INFO]"}, a...)
	_, _ = Cyan.Println(b...)
}

func InfoLn(a ...interface{}) {
	b := append([]interface{}{"[OK]"}, a...)
	_, _ = BlueBold.Println(b...)
}

func WarnLn(a ...interface{}) {
	b := append([]interface{}{"[WARN]"}, a...)
	_, _ = YellowBold.Println(b...)
}

func ErrLn(a ...interface{}) {
	b := append([]interface{}{"[ERROR]"}, a...)
	_, _ = RedBold.Println(b...)
	panic("")
}

func HeaderF(format string, a ...interface{}) {
	format = "[INFO] " + format
	_, _ = CyanBold.Printf(format, a...)
}

func SubHeaderF(format string, a ...interface{}) {
	format = "[INFO] " + format
	_, _ = Cyan.Printf(format, a...)
}

func InfoF(format string, a ...interface{}) {
	format = "[OK] " + format
	_, _ = BlueBold.Printf(format, a...)
}

func WarnF(format string, a ...interface{}) {
	format = "[WARN] " + format
	_, _ = YellowBold.Printf(format, a...)
}

func ErrF(format string, a ...interface{}) {
	format = "[ERROR] " + format
	_, _ = RedBold.Printf(format, a...)
	panic("")
}

func printJson(obj any) {
	j, _ := json.Marshal(obj)
	_, _ = White.Println("[data]", string(j))
}
```
