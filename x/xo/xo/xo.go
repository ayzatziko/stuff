package xo

import (
	"fmt"
	"sync"

	"github.com/ayzatziko/stuff/xerrors"
)

type TypeSign string

const (
	signNull TypeSign = ""
	SignX    TypeSign = "x"
	SignO    TypeSign = "o"
)

const (
	constBoardSizeMax = 3
	constUsersNum     = 2
)

type TypeUser string

func winnerString(u TypeUser) string {
	if u == "" {
		return "None"
	}

	return string(u)
}

func newUserSign(user TypeUser, sign TypeSign) (typeUserSign, error) {
	userSign := typeUserSign{user, sign}
	return userSign, validateSignOfUserSign(userSign)
}

type typeUserSign struct {
	user TypeUser
	sign TypeSign
}

func (userSign typeUserSign) User() TypeUser { return userSign.user }
func (userSign typeUserSign) Sign() TypeSign { return userSign.sign }

func (userSign typeUserSign) String() string {
	return fmt.Sprintf("{user: %q, sign: %q}", userSign.user, userSign.sign)
}

func validateSignOfUserSign(userSign typeUserSign) error {
	if userSign.sign != SignO && userSign.sign != SignX {
		return fmt.Errorf(
			"invalid user sign %s, valid signs %q and %q",
			userSign, SignO, SignX,
		)
	}

	return nil
}

func NewCell(x, y int) (TypeCell, error) {
	c := TypeCell{y, x}
	return c, validateCell(c)
}

type TypeCell struct{ x, y int }

func (c TypeCell) String() string { return fmt.Sprintf("{y: %v, x: %v}", c.y, c.x) }

func validateCell(cell TypeCell) error {
	if constBoardSizeMax > cell.y && cell.y >= 0 && cell.x >= 0 && constBoardSizeMax > cell.x {
		return nil
	}

	return fmt.Errorf("invalid cell %s", cell)
}

type TypeBoard struct {
	rows             [constBoardSizeMax][constBoardSizeMax]TypeSign
	participants     [constUsersNum]typeUserSign
	lastMoveIsDoneBy TypeUser

	winnerSet bool
	winner    TypeUser
}

func (board *TypeBoard) Winner(user TypeUser) (bool, bool) {
	return board.winner == user, board.winnerSet
}

func (board *TypeBoard) Draw(user TypeUser) (bool, bool) {
	return board.winner == "", board.winnerSet
}

func newBoard(user1, user2 typeUserSign, first typeUserSign) (_ *TypeBoard, err error) {
	defer xerrors.Wrap(&err, "NewBoard(user1: %s, user2: %s, first: %s)", user1, user2, first)

	if err := validateSignOfUserSign(user1); err != nil {
		return nil, fmt.Errorf("invalid first user: %v", err)
	} else if err := validateSignOfUserSign(user2); err != nil {
		return nil, fmt.Errorf("invalid second user: %v", err)
	} else if user1.user == user2.user {
		return nil, fmt.Errorf("cannot start game with yourself")
	} else if user1.sign == user2.sign {
		return nil, fmt.Errorf("cannot start game with equal signs")
	} else if user1 != first && user2 != first {
		return nil, fmt.Errorf("passed first user %s is not in partisipants list(%s, %s)", first, user1, user2)
	}

	last := user1
	if first == user1 {
		last = user2
	}

	board := TypeBoard{
		participants:     [2]typeUserSign{user1, user2},
		lastMoveIsDoneBy: last.user,
	}

	return &board, nil
}

func move(board *TypeBoard, cell TypeCell, user TypeUser) (err error) {
	if err := validateBoard(board); err != nil {
		return err
	} else if err := validateCell(cell); err != nil {
		return err
	} else if curSign := board.rows[cell.y][cell.x]; curSign != signNull {
		return fmt.Errorf("cell %s has already value %v, cannot overwrite it", cell, curSign)
	} else if board.participants[0].user != user && board.participants[1].user != user {
		return fmt.Errorf("user %q is not a participant of current game", user)
	} else if board.lastMoveIsDoneBy == user {
		return fmt.Errorf("it is not allowed to make second move in a row, user %q", user)
	} else if board.winnerSet {
		return fmt.Errorf("game is finished, %q is the winner", winnerString(board.winner))
	}

	var sign TypeSign
	if board.participants[0].user == user {
		sign = board.participants[0].sign
	} else {
		sign = board.participants[1].sign
	}

	board.rows[cell.y][cell.x] = sign
	board.lastMoveIsDoneBy = user

	// check is winner
	for _, comb := range winnerCombinations {
		ok := true
		for _, cell := range comb {
			if board.rows[cell.y][cell.x] != sign {
				ok = false
			}
		}

		if ok {
			board.winnerSet = true
			board.winner = user
			break
		}
	}

	if board.winnerSet {
		return nil
	}

	end := true
	for _, row := range board.rows {
		for _, boardSign := range row {
			if boardSign == signNull {
				end = false
			}
		}
	}

	if end {
		board.winnerSet = true
	}

	return nil
}

func validateBoard(b *TypeBoard) error {
	if b == nil {
		return fmt.Errorf("nil board")
	} else if rowsn := len(b.rows); rowsn != constBoardSizeMax {
		return fmt.Errorf("rows number is %d, required to be %d", rowsn, constBoardSizeMax)
	} else if colsn := len(b.rows[0]); colsn != constBoardSizeMax {
		return fmt.Errorf("columns number is %d, required to be %d", colsn, constBoardSizeMax)
	}

	return nil
}

// can be optimized by using a tree or graph data structure.
var winnerCombinations = [...][constBoardSizeMax]TypeCell{
	/*
		xxx
		___
		___
	*/
	{{0, 0}, {0, 1}, {0, 2}},

	/*
		___
		xxx
		___
	*/
	{{1, 0}, {1, 1}, {1, 2}},

	/*
		___
		___
		xxx
	*/
	{{2, 0}, {2, 1}, {2, 2}},

	/*
		x__
		x__
		x__
	*/
	{{0, 0}, {1, 0}, {2, 0}},

	/*
		_x_
		_x_
		_x_
	*/
	{{0, 1}, {1, 1}, {2, 1}},

	/*
		__x
		__x
		__x
	*/
	{{0, 2}, {1, 2}, {2, 2}},

	/*
		x__
		_x_
		__x
	*/
	{{0, 0}, {1, 1}, {2, 2}},

	/*
		__x
		_x_
		x__
	*/
	{{0, 2}, {1, 1}, {2, 0}},
}

// storages
var (
	mu sync.Mutex

	waitingOpponents = map[string]typeUserSign{}
	userBoard        = map[string]*TypeBoard{}
	activeUserToken  = map[string]string{}
	activeTokenUser  = map[string]TypeUser{}

	playsHistory = []typeHistoryRecord{}

	registeredUser = map[string]typeLoginPass{}
)

func cleanDatabase() {
	mu.Lock()
	defer mu.Unlock()

	dropLocked(waitingOpponents)
	dropLocked(userBoard)
	dropLocked(activeTokenUser)
	dropLocked(activeUserToken)
	playsHistory = playsHistory[:0]
	dropLocked(registeredUser)
}

func dropLocked[K comparable, V any](m map[K]V) {
	for k := range m {
		delete(m, k)
	}
}

type typeHistoryRecord struct {
	mayBeWinner, user2 string
	result             typeResult
}

type typeResult bool

const (
	typeResultFirstWon = true
	typeResultDraw     = false
)

func RegisterSelfAsParticipant(sessionToken string, sign TypeSign) (err error) {
	mu.Lock()
	defer mu.Unlock()

	user, ok := activeTokenUser[sessionToken]
	if !ok {
		return fmt.Errorf("session is not found")
	}

	defer xerrors.Wrap(&err, "RegisterSelfAsParticipant(%s, %s)", user, sign)

	userSign, err := newUserSign(user, sign)
	if err != nil {
		return err
	}

	if err := validateSignOfUserSign(userSign); err != nil {
		return err
	} else if board, ok := userBoard[string(userSign.user)]; ok {
		opponent := board.participants[0]
		if board.participants[1].user != userSign.user {
			opponent = board.participants[1]
		}

		return fmt.Errorf("already playing with %s", opponent)
	}

	waitingOpponents[string(userSign.user)] = userSign
	return nil
}

func SearchOpponents() []typeUserSign {
	mu.Lock()
	defer mu.Unlock()

	return valuesOfMap(waitingOpponents)
}

func valuesOfMap[K comparable, V any](m map[K]V) []V {
	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

func StartPlayingWithWaitingOpponent(sessionToken string, sign TypeSign, opponentUser TypeUser) (err error) {
	mu.Lock()
	defer mu.Unlock()

	user, ok := activeTokenUser[sessionToken]
	if !ok {
		return fmt.Errorf("session is not found")
	}

	defer xerrors.Wrap(&err, "StartPlayingWithWaitingOpponent(%s, %s, %s)", user, sign, opponentUser)

	opponentUserSign, ok := waitingOpponents[string(opponentUser)]
	if !ok {
		return fmt.Errorf("opponent %s not found", opponentUser)
	}

	firstUserSign, err := newUserSign(user, sign)
	if err != nil {
		return err
	}

	board, err := newBoard(firstUserSign, opponentUserSign, opponentUserSign)
	if err != nil {
		return err
	}

	delete(waitingOpponents, string(firstUserSign.user))
	delete(waitingOpponents, string(opponentUser))

	userBoard[string(firstUserSign.user)] = board
	userBoard[string(opponentUser)] = board

	return nil
}

func MakeAMove(sessionToken string, cell TypeCell) (_ *TypeBoard, _ string, err error) {
	mu.Lock()
	defer mu.Unlock()

	user, ok := activeTokenUser[sessionToken]
	if !ok {
		return nil, "", fmt.Errorf("session not found")
	}

	defer xerrors.Wrap(&err, "MakeAMove(%s, %s)", user, cell)

	board, ok := userBoard[string(user)]
	if !ok {
		return nil, "", fmt.Errorf("user %s does not participate in any play", user)
	}

	if err = move(board, cell, user); err != nil {
		return nil, "", err
	}

	if board.winnerSet && board.winner != "" {
		winner, loser := board.participants[0].user, board.participants[1].user
		if loser == user {
			winner, loser = board.participants[1].user, board.participants[0].user
		}

		delete(userBoard, string(user))
		delete(userBoard, string(loser))

		playsHistory = append(playsHistory, typeHistoryRecord{mayBeWinner: string(winner), user2: string(loser), result: typeResultFirstWon})

		return board, fmt.Sprintf("%s wins %s", winner, loser), nil
	} else if board.winnerSet {
		user1, user2 := string(board.participants[0].user), string(board.participants[1].user)
		delete(userBoard, string(user1))
		delete(userBoard, string(user2))

		playsHistory = append(playsHistory, typeHistoryRecord{mayBeWinner: user1, user2: user2, result: typeResultDraw})

		return board, "draw", nil
	}

	return nil, "", nil
}

type typeLoginPass struct {
	username string
	password string // use bcrypt
}

func RegisterUser(username, password string) error {
	mu.Lock()
	defer mu.Unlock()

	_, ok := registeredUser[username]
	if ok {
		return fmt.Errorf("user %q already exists", username)
	}

	registeredUser[username] = typeLoginPass{username: username, password: password}

	return nil
}

func Login(username, password string) (_ string, err error) {
	defer xerrors.Wrap(&err, "Login(%s, *****)", username)

	mu.Lock()
	defer mu.Unlock()

	user, ok := registeredUser[username]
	if !ok {
		return "", fmt.Errorf("user %q not found", username)
	}

	if user.password != password {
		return "", fmt.Errorf("password does not match")
	}

	sessionToken := randomStringLocked()

	// deregister from other session
	existingToken, ok := activeUserToken[username]
	if ok {
		delete(activeUserToken, username)
		delete(activeTokenUser, existingToken)

		// immediately make an opponent a winner
		board, boardExists := userBoard[username]
		_, waitingAPlay := waitingOpponents[username]
		if boardExists {
			// the same logic as in make a move, abstarct
			winner, loser := board.participants[0].user, board.participants[1].user
			if string(winner) == username {
				winner, loser = loser, winner
			}

			delete(userBoard, string(winner))
			delete(userBoard, string(loser))

			playsHistory = append(playsHistory, typeHistoryRecord{mayBeWinner: string(winner), user2: string(loser), result: typeResultFirstWon})
		} else if waitingAPlay {
			delete(waitingOpponents, username)
		}
	}

	activeTokenUser[sessionToken] = TypeUser(username)
	activeUserToken[username] = sessionToken

	return sessionToken, nil
}

func Logout(sessionToken string) error {
	mu.Lock()
	defer mu.Unlock()

	user, ok := activeTokenUser[sessionToken]
	if !ok {
		return fmt.Errorf("session not found")
	}
	username := string(user)

	delete(activeUserToken, username)
	delete(activeTokenUser, sessionToken)

	// immediately make an opponent a winner
	board, boardExists := userBoard[username]
	_, waitingAPlay := waitingOpponents[username]
	if boardExists { // same as login, abstract
		// the same logic as in make a move, abstarct
		winner, loser := board.participants[0].user, board.participants[1].user
		if winner == user {
			winner, loser = loser, winner
		}

		delete(userBoard, string(winner))
		delete(userBoard, string(loser))

		playsHistory = append(playsHistory, typeHistoryRecord{mayBeWinner: string(winner), user2: string(loser), result: typeResultFirstWon})
	} else if waitingAPlay {
		delete(waitingOpponents, username)
	}

	return nil
}

var incToken int

func randomStringLocked() string {
	incToken++
	return fmt.Sprintf("%d", incToken)
}
