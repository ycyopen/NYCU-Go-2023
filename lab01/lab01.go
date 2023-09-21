package main

import "fmt"

func main() {
	fmt.Println("Welcome to Simple Calculator")

	var a, b int64
	fmt.Print("Enter first number: ")
	fmt.Scan(&a)

	fmt.Print("Enter second number: ")
	fmt.Scan(&b)

	fmt.Println("Add:", Add(a, b))
	fmt.Println("Subtract:", Sub(a, b))
	fmt.Println("Multiply:", Mul(a, b))
	fmt.Println("Divide:", Div(a, b))
}

func Add(a int64,b int64)int64{
	return a+b
}

func Sub(a int64,b int64)int64{
	return a-b
}

func Mul(a int64,b int64)int64{
	return a*b
}

func Div(a int64,b int64)int64{
	return a/b
}
// TODO: Create `Add`, `Sub`, `Mul`, `Div` function here
