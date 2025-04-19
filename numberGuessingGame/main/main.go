package main

import (
	"fmt"
	"math/rand"
)

func difficultDescriber() {
	fmt.Print(`Please Select the difficulty level:
                1. Easy (10 chances)
                2. Medium (5 chances)
                3. Hard (3 chances) :
`)
}



type Difficulty struct {
	Name   string
	Chances int
}

func main() {
	var difficultyLevel int
	var usersGuess int
	var chancesLeft int
	computerNum := rand.Intn(100) + 1 
 println(computerNum)
	difficultyLevels := map[int]Difficulty{
		1: {"Easy", 10},
		2: {"Medium", 7},
		3: {"Hard", 3},
	}
	fmt.Println("Welcome to the Number Guessing Game!")
	fmt.Println("I'm thinking of a number between 1 and 100.")
	for {
		difficultDescriber()
		num, err := fmt.Scan(&difficultyLevel)
		if err != nil || num != 1 || difficultyLevel < 1 || difficultyLevel > 3 {
			fmt.Println("Invalid value entered. Please enter a number between 1 and 3.")
			fmt.Scanln()
			continue
		}
		break
	}
selectedValue := difficultyLevels[difficultyLevel]
	fmt.Printf("Great! You have selected the %v difficulty level \n", selectedValue.Name)
	fmt.Println("Let's start the game!")
	chancesLeft = selectedValue.Chances
	
	for i :=range chancesLeft {
		
		fmt.Print("Enter your guess: ")
		_, err := fmt.Scan(&usersGuess)
		 if err != nil {
			fmt.Print("There is an error somewhere")
		 }

		 if usersGuess < 1 || usersGuess > 100 {
			fmt.Println("You have entered an invalid number")
		 }

		 if usersGuess != computerNum {
			
				if usersGuess > computerNum {
					fmt.Printf("Incorrect! The number is less than %v.\n",usersGuess)
				}else{
					fmt.Printf("Incorrect! The number is greater than %v.\n",usersGuess)
				}
	
		 }else{

			 fmt.Printf("Congratulations! You guessed the correct number in %v attempts.\n",i+1)
			 break
		 }

	
	if i+1 == chancesLeft {
		fmt.Println("You have Exhausted number of tries , Game Over")
	}


	}
	
}
