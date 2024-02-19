package main

import "fmt"

type Add struct {
	Name string
}

func a(aa interface{}) {
	switch real := aa.(type) {
	case *Add:
		fmt.Println(real.Name)
	default:
		fmt.Println("没有匹配到")
	}
}

func main() {
	a(&Add{Name: "wyw"})
}
