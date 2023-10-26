package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// TODO: Create a struct to hold the data sent to the template
type Data struct {
	op         string
	num1       int64
	num2       int64
	err1       error
	err2       error
	Result     int64
	Expression string
}

func gcd(a, b int64) int64 {
	if b == 0 {
		return a
	}
	return gcd(b, a%b)
}

func Calculator(w http.ResponseWriter, r *http.Request) {
	// TODO: Finish this function
	t, err_t := template.ParseFiles("index.html")
	errpage, err_e := template.ParseFiles("error.html")
	cal := Data{
		op: r.URL.Query().Get("op"),
	}
	num1, err1 := strconv.ParseInt(r.URL.Query().Get("num1"), 10, 64)
	num2, err2 := strconv.ParseInt(r.URL.Query().Get("num2"), 10, 64)
	cal.num1 = num1
	cal.num2 = num2
	cal.err1 = err1
	cal.err2 = err2
	if err_t != nil || err_e != nil || cal.op == "" || cal.err1 != nil || cal.err2 != nil {
		errpage.Execute(w, cal)
		return
	} else if cal.op == "add" {
		cal.Result = cal.num1 + cal.num2
		cal.Expression = fmt.Sprintf("%d + %d", cal.num1, cal.num2)
	} else if cal.op == "sub" {
		cal.Result = cal.num1 - cal.num2
		cal.Expression = fmt.Sprintf("%d - %d", cal.num1, cal.num2)
	} else if cal.op == "mul" {
		cal.Result = cal.num1 * cal.num2
		cal.Expression = fmt.Sprintf("%d * %d", cal.num1, cal.num2)
	} else if cal.op == "div" {
		if cal.num2 == 0 {
			errpage.Execute(w, cal)
			return
		}
		cal.Result = cal.num1 / cal.num2
		cal.Expression = fmt.Sprintf("%d / %d", cal.num1, cal.num2)
	} else if cal.op == "gcd" {
		cal.Expression = fmt.Sprintf("GCD(%d, %d)", cal.num1, cal.num2)
		cal.Result = gcd(cal.num1, cal.num2)
	} else if cal.op == "lcm" {
		cal.Result = cal.num1 * cal.num2 / gcd(cal.num1, cal.num2)
		cal.Expression = fmt.Sprintf("LCM(%d, %d)", cal.num1, cal.num2)
	} else {
		errpage.Execute(w, cal)
		return
	}
	t.Execute(w, cal)
}

func main() {
	http.HandleFunc("/", Calculator)
	log.Fatal(http.ListenAndServe(":8084", nil))
}
