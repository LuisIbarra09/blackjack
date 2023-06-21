package blackjack

import (
	"fmt"

	"github.com/LuisIbarra09/deck"
)

type state int8

const (
	statePlayerTurn state = iota
	stateDealerTurn
	stateHandOver
)

func New() GameState {
	return GameState{
		state:   statePlayerTurn,
		balance: 0,
	}
}

type GameState struct {
	// unexported fields
	deck    []deck.Card
	state   state
	player  []deck.Card
	dealer  []deck.Card
	balance int
}

func (gs *GameState) currentHand() *[]deck.Card {
	switch gs.state {
	case statePlayerTurn:
		return &gs.player
	case stateDealerTurn:
		return &gs.dealer
	default:
		panic("It isn't currently any player's turn")
	}
}

func deal(gs *GameState) {
	gs.player = make([]deck.Card, 0, 5)
	gs.dealer = make([]deck.Card, 0, 5)
	var card deck.Card
	for i := 0; i < 2; i++ {
		card, gs.deck = draw(gs.deck)
		gs.player = append(gs.player, card)
		card, gs.deck = draw(gs.deck)
		gs.dealer = append(gs.dealer, card)
	}
	gs.state = statePlayerTurn
}

func (gs *GameState) Play(ai AI) int {
	gs.deck = deck.New(deck.Deck(3), deck.Shuffle)

	for i := 0; i < 10; i++ {
		deal(gs)

		for gs.state == statePlayerTurn {
			hand := make([]deck.Card, len(gs.player))
			copy(hand, gs.player)
			move := ai.Play(hand, gs.dealer[0])
			move(gs)
		}

		// Dealer Tunr == gs.State = StateDealerTurn
		for gs.state == stateDealerTurn {
			dScore := Score(gs.dealer...)
			if dScore <= 16 || (dScore == 17 && Soft(gs.dealer...)) {
				MoveHit(gs)
			} else {
				MoveStand(gs)
			}
		}

		// Score calculation
		endHand(gs, ai)
	}
	return 0
}

type Move func(*GameState)

func MoveHit(gs *GameState) {
	hand := gs.currentHand()
	var card deck.Card
	card, gs.deck = draw(gs.deck)
	*hand = append(*hand, card)
	if Score(*hand...) > 21 {
		MoveStand(gs)
	}
}

func MoveStand(gs *GameState) {
	gs.state++
}

func draw(cards []deck.Card) (deck.Card, []deck.Card) {
	return cards[0], cards[1:]
}

func endHand(gs *GameState, ai AI) {
	pScore, dScore := Score(gs.player...), Score(gs.dealer...)
	// TODO:  Figure out winnings and add/subtract them
	switch {
	case pScore > 21:
		fmt.Println("You busted")
		gs.balance--
	case dScore > 21:
		fmt.Println("Dealer busted")
		gs.balance++
	case pScore > dScore:
		fmt.Println("You win!")
		gs.balance++
	case dScore > pScore:
		fmt.Println("You lose")
		gs.balance--
	case pScore == dScore:
		fmt.Println("Draw")
	}
	// Marcar una nueva linea
	fmt.Println()
	ai.Results([][]deck.Card{gs.player}, gs.dealer)
	// Reset hands
	gs.player = nil
	gs.dealer = nil
}

// Score will take in a hand of cards and return the best blackjack score
// possible with the hand.
func Score(hand ...deck.Card) int {
	minScore := minScore(hand...)
	if minScore > 11 {
		return minScore
	}
	for _, c := range hand {
		if c.Rank == deck.Ace {
			// ace is currently worth 1, and we are changing it to be worth 11
			minScore += 10
		}
	}
	return minScore
}

// Soft returns true if the score of a hand is a soft score - that is if an ace
// is being countend as 11 points.
func Soft(hand ...deck.Card) bool {
	minScore := minScore(hand...)
	score := Score(hand...)
	return minScore != score
}

// Func para asignar un score a la mano de los jugadores (puntuajes minimos)
func minScore(hand ...deck.Card) int {
	score := 0
	for _, c := range hand {
		score += min(int(c.Rank), 10)
	}
	return score
}

// Sirve para asignar que las cartas J, Q y K tienen un valor de 10
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
