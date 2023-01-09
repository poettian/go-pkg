package main

type Animal struct {
	Name string
}

func (a *Animal) SetName(name string) {
	a.Name = name
}

type Cat struct {
	Animal
}

func main() {
	cat := &Cat{}
	cat.SetName("小咪")
}
