package main

import "fmt"

func main() {
	var revenue float64
	var taxRate float64
	var expenses float64
	fmt.Print("Enter Revenue: ")
	fmt.Scan(&revenue)

	fmt.Print("Enter Expenses: ")
	fmt.Scan(&expenses)

	fmt.Print("Enter tax rate: ")
	fmt.Scan(&taxRate)



	

	var ebt = revenue - expenses
	var profit = ebt * (1-taxRate/100)
	var ratio = ebt/profit
	fmt.Println("Ebt",ebt)
	fmt.Println("Profit",profit)
	fmt.Println("Ratio",ratio)
}
