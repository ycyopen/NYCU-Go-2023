package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func Calculator(w http.ResponseWriter, r *http.Request) {
	// TODO: implement a calculator
	path := strings.Split(r.URL.Path, "/")
	if len(path) != 4 {
		fmt.Fprintf(w, "Error!")
		return
	}
	a, err_a := strconv.ParseInt(path[2], 10, 64)//num1
	b, err_b := strconv.ParseInt(path[3], 10, 64)//num2
	operation := path[1]
	if err_a != nil || err_b != nil {
		fmt.Fprintf(w, "Error!")
	} else if operation == "add" {
		fmt.Fprintf(w, "%d + %d = %d", a, b, a+b)
	} else if operation == "sub" {
		fmt.Fprintf(w, "%d - %d = %d", a, b, a-b)
	} else if operation == "mul" {
		fmt.Fprintf(w, "%d * %d = %d", a, b, a*b)
	} else if operation == "div" {
		if b == 0 {
			fmt.Fprintf(w, "Error!")
			return
		}
		fmt.Fprintf(w, "%d / %d = %d, reminder = %d", a, b, a/b, a%b)
	} else {
		fmt.Fprintf(w, "Error!")
	}
}

func main() {
	http.HandleFunc("/", Calculator)
	log.Fatal(http.ListenAndServe(":8083", nil))
}
