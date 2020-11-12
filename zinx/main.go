package main

import "fmt"

func main() {
	fmt.Println("dsdaasd")

	is := false
	SetTrue(&is)
	fmt.Println(is)
}

func SetTrue(is* bool){
	*is = true
}