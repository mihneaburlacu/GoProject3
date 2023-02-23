package main

import (
	"GoProject3/handlers"
	"fmt"
	"net/http"
)

func main() {
	err := handlerMain(http.ListenAndServe)
	if err != nil {
		fmt.Printf("Error while calling in main function: %v\n", err)
		return
	}
}

func handlerMain(serveFunc func(addr string, handler http.Handler) error) error {
	http.HandleFunc("/till-salary/how-much", handlers.TillSalaryHandler)
	http.HandleFunc("/till-salary/pay-day/", handlers.ListPayDayDatesHandler)

	err := serveFunc(":8080", nil)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}

	return nil
}
