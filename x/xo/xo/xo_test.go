package xo_test

import (
	"testing"

	. "github.com/ayzatziko/stuff/x/xo/xo"
)

func TestFlowTestPlayWithWaitingOpponent(t *testing.T) {
	t.Cleanup(CleanDatabase)

	first, second := "user1", "user2"

	err := RegisterUser(first, "")
	failIfError(t, err)
	tokenFirst, err := Login(first, "")
	failIfError(t, err)

	err = RegisterSelfAsParticipant(tokenFirst, SignO)
	failIfError(t, err)

	err = RegisterUser(second, "")
	failIfError(t, err)
	tokenSecond, err := Login(second, "")
	failIfError(t, err)

	opponents := SearchOpponents()
	secSign := oppositeSign[opponents[0].Sign()]
	err = StartPlayingWithWaitingOpponent(tokenSecond, secSign, opponents[0].User())
	failIfError(t, err)

	currentMoveSession, nextMoveSession := tokenFirst, tokenSecond

	move := func(x, y int, winner bool) *TypeBoard {
		t.Helper()

		cell, err := NewCell(x, y)
		failIfError(t, err)

		b, msg, err := MakeAMove(currentMoveSession, cell)
		failIfError(t, err)

		if !winner && (b != nil || msg != "") {
			t.Fatalf("unexpected winner: %v", msg)
		} else if winner && (b == nil || msg == "") {
			t.Fatalf("expected winner but it is absent")
		}

		currentMoveSession, nextMoveSession = nextMoveSession, currentMoveSession
		return b
	}

	move(0, 0, false)
	move(1, 0, false)
	move(0, 1, false)
	move(1, 1, false)
	b := move(0, 2, true)
	ok, end := b.Winner(TypeUser(first))

	failIfFalseFmt(t, end, "expected play is finished")
	failIfFalseFmt(t, ok, "expected first player has won")
}

var oppositeSign = map[TypeSign]TypeSign{
	SignO: SignX,
	SignX: SignO,
}

func failIfFalseFmt(t *testing.T, ok bool, msg string, args ...any) {
	t.Helper()

	if ok {
		return
	}

	t.Fatalf(msg, args...)
}

func failIfError(t *testing.T, err error) {
	t.Helper()

	if err == nil {
		return
	}

	t.Fatal(err)
}
