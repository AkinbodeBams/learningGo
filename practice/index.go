package main

import "fmt"

func main() {
	var revenue float64
	var taxRate float64
	var expenses float64
	fmt.Print("Enter Revenue: ")
	fmt.Scan(&revenue)
	fmt.Print("Enter tax rate: ")
	fmt.Scan(taxRate)
	fmt.Print("Enter Expenses: ")
	fmt.Scan(expenses)

	var ebt = revenue - expenses
	var eat = (revenue -expenses ) -taxRate
	var ratio = ebt/eat
	fmt.Print("Ebt",ebt)
	fmt.Print("Profit",eat)
	fmt.Print("Ratio",ratio)
}
