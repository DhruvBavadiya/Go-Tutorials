package main

import (
	"fmt"
	"math/rand"
)

func main() {
	nguess()
}

// Score struct to store each game's score details
type Score struct {
	isFind     bool
	attempt    int
	Difficulty // Embedding Difficulty struct
}

// Difficulty struct to hold difficulty details
type Difficulty struct {
	tries                   int
	difficulty_level_string string
	points                  int
}

// Slice of different difficulty levels
var tries = []Difficulty{
	{
		10,
		"easy",
		50,
	},
	{
		5,
		"medium",
		150,
	},
	{
		3,
		"hard",
		250,
	},
}

// Slice to store user scores for each game
var user_score = []Score{}

// nguess function for the number guessing game
func nguess() {
	fmt.Println("Welcome to the game")

	var want_to_play bool = true
	for want_to_play {
		// Difficulty selection
		fmt.Println("Select the difficulty \n 1. Easy (10 chances) \n 2. Medium (5 chances) \n 3. Hard (3 chances)")

		var difficulty int
		_, err := fmt.Scan(&difficulty)
		if err != nil {
			fmt.Println("Please select a valid number.")
			continue
		}

		if difficulty < 1 || difficulty > 3 {
			fmt.Println("Please select a valid difficulty (1, 2, or 3).")
			continue
		}

		var usertries Difficulty = tries[difficulty-1]
		fmt.Printf("You have selected %s difficulty with %d tries.\n", usertries.difficulty_level_string, usertries.tries)

		fmt.Println("Guess a number between 1 and 100")
		num := rand.Intn(100) + 1 // Generate random number between 1 and 100
		var attempts int = 0      // Initialize attempt counter

		for usertries.tries > 0 {
			fmt.Printf("You have %d tries left. Enter your guess: ", usertries.tries)
			var guess int
			_, err := fmt.Scan(&guess)
			if err != nil {
				fmt.Println("Please enter a valid number.")
				continue
			}
			attempts++ // Increment attempts count

			if guess == num {
				fmt.Printf("Congratulations! You guessed the number in %d attempts!\n", attempts)
				var score = Score{
					isFind:     true,
					attempt:    attempts,
					Difficulty: usertries,
				}
				user_score = append(user_score, score)
				break
			} else if guess > num {
				fmt.Println("Your guess is greater than the number.")
			} else {
				fmt.Println("Your guess is less than the number.")
			}

			usertries.tries-- // Reduce remaining tries
		}

		// If the player runs out of tries, show the correct number
		if usertries.tries == 0 {
			fmt.Printf("You have run out of tries! The correct number was %d.\n", num)
			var score = Score{
				isFind:     false,
				attempt:    attempts,
				Difficulty: usertries,
			}
			user_score = append(user_score, score)
		}

		// Ask the user if they want to play again
		fmt.Println("Do you want to play again? Select 1 for 'Yes', 0 for 'No'")
		var user_answer int
		_, err = fmt.Scan(&user_answer)
		if err != nil {
			fmt.Println("Please select a valid option (1 or 0).")
			return
		}

		// Determine if the user wants to play again
		if user_answer == 1 {
			want_to_play = true
		} else {
			want_to_play = false
		}
	}

	// Display all scores at the end
	fmt.Println("\nGame Over. Here's your game summary:")
	var finalscore int = 0
	for index, score := range user_score {
		if score.isFind {
			finalscore += (score.tries) * score.points
		}
		fmt.Printf("Game %d: Difficulty = %s, Guessed = %t, Attempts = %d\n",
			index+1, score.difficulty_level_string, score.isFind, score.attempt)
	}

	fmt.Println("Your final score is %d ", finalscore)

}
