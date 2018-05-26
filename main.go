package main // import "gophercise/blackjack"

import (
	"fmt"
	"strings"

	"github.com/Athulus/deck"
)

type hand []deck.Card

type state uint8

const (
	newGame state = iota
	playersTurn
	dealersTurn
	endHand
)

type game struct {
	deck      []deck.Card
	gameState state
	player    hand
	dealer    hand
}

func (h hand) String() string {
	s := make([]string, len(h))
	for i := range h {
		s[i] = h[i].String()
	}
	return strings.Join(s, ", ")
}
func (h hand) dealerString() string {
	return h[0].String() + ", ***SECRET***"
}

func (h hand) minScore() int {
	sum := 0
	for _, card := range h {
		if (card.Rank) > 10 {
			sum += 10
		} else {
			sum += int(card.Rank)
		}
	}
	return sum
}

func (h hand) score() int {
	finalScore := h.minScore()
	if finalScore > 11 {
		return finalScore
	}
	for _, c := range h {
		if c.Rank == deck.Ace {
			return finalScore + 10
		}
	}
	return finalScore

}

func cloneGame(g game) game {
	ret := game{
		deck:      make([]deck.Card, len(g.deck)),
		gameState: g.gameState,
		player:    make(hand, len(g.player)),
		dealer:    make(hand, len(g.dealer)),
	}

	copy(ret.deck, g.deck)
	copy(ret.player, g.player)
	copy(ret.dealer, g.dealer)
	return ret
}

func (g *game) currentPlayer() *hand {
	switch g.gameState {
	case playersTurn:
		return &g.player
	case dealersTurn:
		return &g.dealer
	default:
		panic("it is not currently anyones turn")
	}
}

func shuffle(g game) game {
	new := cloneGame(g)
	new.deck = deck.New(deck.Deck(3), deck.Shuffle)
	return new
}

func deal(g game) game {
	new := cloneGame(g)
	new.player = make(hand, 0, 5)
	new.dealer = make(hand, 0, 5)
	for i := 0; i < 2; i++ {
		new.player = append(new.player, draw(&new.deck))
		new.dealer = append(new.dealer, draw(&new.deck))
	}
	new.gameState = playersTurn
	return new
}

func hit(g game) game {
	new := cloneGame(g)
	hand := new.currentPlayer()
	*hand = append(*hand, draw(&new.deck))
	if hand.score() > 21 {
		return stand(new)
	}
	return new
}

func stand(g game) game {
	new := cloneGame(g)
	new.gameState++
	return new
}

func cleanUpHand(g game) game {
	new := cloneGame(g)
	pScore, dScore := new.player.score(), new.dealer.score()
	fmt.Println("--Hand is over, who won?--")
	fmt.Println("dealer hand:", new.dealer)
	fmt.Println("player hand:", new.player)
	switch {
	case pScore > 21:
		fmt.Println("Player busted")
	case dScore > 21:
		fmt.Println("dealer busted, player wins")
	case pScore > dScore:
		fmt.Println("Player wins")
	case pScore < dScore:
		fmt.Println("dealer wins")
	case pScore == dScore:
		fmt.Println("draw")
	}
	fmt.Println()

	new.gameState = newGame
	new.player = nil
	new.dealer = nil

	return new
}

func main() {

	var g game
	g = shuffle(g)

	for i := 0; i < 10; i++ {

		// //deal initial hand
		fmt.Printf("NEW HAND! %d/10 \n", i+1)
		g = deal(g)
		fmt.Println("dealer hand:", g.dealer.dealerString())
		fmt.Println("player hand:", g.player)

		//allow player to hit
		var input string
		fmt.Println("PLAYERS TURN!")
		for g.gameState == playersTurn {
			fmt.Println("Current score: ", g.player.score())
			fmt.Println("what do you want to do? hit(h) or stand(s)")
			fmt.Scanf("%s \n", &input)

			switch input {
			case "h":
				g = hit(g)
				fmt.Println("player hand:", g.player)
			case "s":
				g = stand(g)
			default:
				fmt.Println("that was not an accepted input, type `h` or `s`")
			}

		}

		//dealer logic
		fmt.Println("DEALERS TURN!")
		for g.gameState == dealersTurn {
			fmt.Println("dealer hand:", g.dealer)
			if g.dealer.score() <= 16 || (g.dealer.score() == 17 && g.dealer.minScore() > 17) {
				fmt.Println("dealer hits")
				g = hit(g)
			} else {
				fmt.Println("dealer stands")
				g = stand(g)
			}

		}

		//display score and winner as well as reset the player and dealers hands
		g = cleanUpHand(g)
		if i != 9 {
			fmt.Println("press Enter to play another hand, or press Ctl+c to quit")
			var x string
			fmt.Scanln(&x)
		}
	}
}

func draw(d *[]deck.Card) deck.Card {
	if len(*d) < 1 {
		return deck.Card{}
	}
	var card deck.Card
	card, *d = (*d)[0], (*d)[1:]
	return card
}

func dealerTurn() {

}
