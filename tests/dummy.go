package main

import "fmt"

func main() {
	var arr []string = nil
	fmt.Println(arr, arr == nil)

	array := make([]string, 0)
	arr = array
	fmt.Println(arr, arr == nil)

	arr = append(arr, "one")
	fmt.Println(arr, arr == nil)

	arr = nil
	fmt.Println(arr, arr == nil)
}
