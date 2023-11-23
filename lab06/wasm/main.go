package main

import (
	"fmt"
	"math/big"
	"strconv"
	"syscall/js"
)

func CheckPrime(this js.Value, args []js.Value) interface{} {
	// TODO: Check if the number is prime
	str := js.Global().Get("document").Call("getElementById", "value").Get("value").String()
	num, _ := strconv.ParseInt(str, 10, 64)
	x := big.NewInt(num)
	if x.ProbablyPrime(0) {
		js.Global().Get("document").Call("getElementById", "answer").Set("innerText", "It's prime")
	} else {
		js.Global().Get("document").Call("getElementById", "answer").Set("innerText", "It's not prime")
	}
	return nil
}

func registerCallbacks() {
	// TODO: Register the function CheckPrime
	js.Global().Set("CheckPrime", js.FuncOf(CheckPrime))
}

func main() {
	fmt.Println("Golang main function executed")
	registerCallbacks()

	//need block the main thread forever
	select {}
}
