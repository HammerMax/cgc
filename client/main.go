package main

import "fmt"

type demo struct {
	str string
}

func (d *demo) changeStr() {
	d.str = "changeStr le"
}

func (d demo) cc() {
	d.str = "ccc"
}

type mapDemo map[string]string

func (m mapDemo) mapchange() {
	m["11"] = "11"
}

func main() {
	d := demo{str: "f"}
	fmt.Println(d)

	d.changeStr()
	fmt.Println(d)

	d.cc()
	fmt.Println(d)

	mm := mapDemo{
		"11" : "demo",
	}
	fmt.Println(mm)

	mm.mapchange()
	fmt.Println(mm)
}

