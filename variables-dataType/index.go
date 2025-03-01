package main

import (
	"fmt"
	"math"
)

func main() {
	var investmentAmount float64
	const inflationRate = 6.5
	const expectedReturnRate = 5.5
	var years float64 
	fmt.Print("Enter investment Amount: ")
	fmt.Scan(&investmentAmount)
	fmt.Print("Enter years: ")
	fmt.Scan(&years)

	var futureValue = investmentAmount +   math.Pow(1 + expectedReturnRate/100,years)
	futureRealValue := futureValue/ math.Pow(1 + inflationRate/100,years)
	
	fmt.Println("Future Value is: ",futureValue)
	fmt.Println("Future real value is: ",futureRealValue)
	
}
