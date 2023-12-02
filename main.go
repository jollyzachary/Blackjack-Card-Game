package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type Card struct {
	value int
	suit  int // 0 - spades, 1 - hearts, 2 - diamonds, 3 - clubs
}

func (card Card) getString() string {
	var suit string
	var value string

	switch card.suit {
	case 0:
		suit = "♠"
	case 1:
		suit = "♥"
	case 2:
		suit = "♦"
	case 3:
		suit = "♣"
	}

	switch card.value {
	case 11:
		value = "J"
	case 12:
		value = "Q"
	case 13:
		value = "K"
	case 1:
		value = "A"
	default:
		value = fmt.Sprintf("%d", card.value)
	}

	return value + suit
}

type Deck struct {
	cards []Card
}

func (d *Deck) deal(num uint) []Card {
	dealtCards := []Card{}

	for i := uint(0); i < num; i++ {
		card := d.cards[0]
		dealtCards = append(dealtCards, card)
		d.cards = d.cards[1:] // remove first element from slice
	}

	return dealtCards
}

func (d *Deck) create() {
	for i := 1; i <= 13; i++ {
		for j := 0; j < 4; j++ {
			card := Card{value: i, suit: j}
			d.cards = append(d.cards, card)
		}
	}
}

func (d *Deck) shuffle() {
	rand.Shuffle(len(d.cards), func(i, j int) { d.cards[i], d.cards[j] = d.cards[j], d.cards[i] })
}

type Game struct {
	deck        Deck
	playerCards []Card
	dealerCards []Card
}

func (game *Game) dealStartingCards() {
	game.deck.create()
	game.deck.shuffle()
	game.dealerCards = game.deck.deal(2)
	game.playerCards = game.deck.deal(2)

	fmt.Printf("Player has been dealt: ")
	displayCards(game.playerCards)
	fmt.Printf("\nDealer shows: %s\n", game.dealerCards[0].getString())
}

func (game *Game) play(bet float64) float64 {
	defer fmt.Printf("\n--------------------------\n\n")

	fmt.Printf("\n--------------------------\n\n")

	game.dealStartingCards()

	// check for blackjack
	playerHasBlackjack := isBlackjack(game.playerCards)
	dealerHasBlackjack := isBlackjack(game.dealerCards)
	if playerHasBlackjack && !dealerHasBlackjack {
		fmt.Println("Blackjack!")
		return float64(bet) * 1.5
	} else if playerHasBlackjack && dealerHasBlackjack {
		return 0
	} else if dealerHasBlackjack {
		fmt.Println("The dealers other card is:", game.dealerCards[1].getString())
		fmt.Println("The dealer has blackjack, you lost :/")
		return -bet
	}

	// playerTurn() returns true if the player lost during their turn
	if game.playerTurn() {
		return -bet
	}

	fmt.Println()
	game.dealerTurn()

	// Determine the winner
	isWinner := isPlayerWinner(game.playerCards, game.dealerCards) // 0 = push, 1 = win, -1 = lost
	if isWinner == 1 {
		fmt.Printf("You won $%v\n", bet)
		return bet
	} else if isWinner == 0 {
		fmt.Println("Push")
		return 0
	}

	fmt.Println("You lost...")
	return -bet
}

func (game *Game) playerTurn() bool {
	for true {
		fmt.Printf("\nWould you like to hit or stay (H/S)? ")
		hitOrStay := enterString()

		if hitOrStay == "H" || hitOrStay == "h" {
			card := game.deck.deal(1)[0]
			game.playerCards = append(game.playerCards, card)
			fmt.Printf("You are dealt: %v\n", card.getString())
		} else {
			return false // didn't lose
		}
		fmt.Printf("You now have: ")
		displayCards(game.playerCards)

		value := getCardValues(game.playerCards)
		if value > 21 {
			fmt.Println("Oops, you busted :/")
			return true
		} else if value == 21 {
			fmt.Println("21! Nice.")
			return false
		}
	}
	return false
}

func (game *Game) dealerTurn() {
	fmt.Printf("Dealer reveals: ")
	displayCards(game.dealerCards)

	for true {
		time.Sleep(1 * time.Second)
		if shouldDealerHit(game.dealerCards) {
			card := game.deck.deal(1)[0]
			game.dealerCards = append(game.dealerCards, card)
			fmt.Printf("Dealer hits and receives: %v\n", card.getString())
		} else {
			break
		}
		fmt.Printf("Dealer now has: ")
		displayCards(game.dealerCards)
		fmt.Println()
	}
}

func enterString() string {
	reader := bufio.NewReader(os.Stdin)
	// ReadString will block until the delimiter is entered
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("An error occurred while reading input. Please try again", err)
		return ""
	}

	// remove the delimiter from the string
	input = strings.TrimSuffix(input, "\n")
	return input
}

func getCardValues(cards []Card) int {
	aces := 0
	value := 0

	for _, card := range cards {
		if card.value == 1 {
			aces++
		} else {
			value += int(math.Min(10, float64(card.value)))
		}
	}

	if aces == 0 {
		return value
	}

	if value < (11 + aces - 1) {
		return value + 11 + aces - 1
	} else {
		return value + aces
	}
}

func isPlayerWinner(playerHand []Card, dealerHand []Card) int {
	playerValue := getCardValues(playerHand)
	dealerValue := getCardValues(dealerHand)

	if playerValue == dealerValue {
		return 0
	} else if playerValue > dealerValue || dealerValue > 21 {
		return 1
	} else {
		return -1
	}
}

func shouldDealerHit(dealerCards []Card) bool {
	value := getCardValues(dealerCards)
	return value <= 16
}

func isBlackjack(cards []Card) bool {
	value := getCardValues(cards)
	return value == 21 && len(cards) == 2
}

func displayCards(cards []Card) {
	value := getCardValues(cards)
	displayStr := ""

	for i, card := range cards {
		displayStr += card.getString()
		if i < len(cards)-1 {
			displayStr += " "
		}
	}

	fmt.Printf("%v = %v\n", displayStr, value)
}

func main() {
	balance := float64(100)

	for balance > 0 {
		fmt.Printf("Your balance is: $%.2f\n", balance)
		fmt.Printf("Enter your bet (q to quit): ")
		bet, err := strconv.ParseFloat(enterString(), 64)
		if err != nil {
			break
		}
		if bet > balance || bet <= 0 {
			fmt.Println("Invalid bet.")
			continue
		}

		game := Game{}
		balance += game.play(bet)
	}

	fmt.Printf("You left with: $%2.f\n", balance)
}
