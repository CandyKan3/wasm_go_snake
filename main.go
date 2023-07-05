package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"

	"github.com/candykan31/wasm_go_snake/models"
	spinhttp "github.com/fermyon/spin/sdk/go/http"
	"github.com/julienschmidt/httprouter"
)

const (
	contentType           = "Content-Type"
	applicationJson       = "application/json"
	internalServiceErrMsg = "internal service error"
)

func init() {
	spinhttp.Handle(func(w http.ResponseWriter, r *http.Request) {
		router := spinhttp.NewRouter()
		router.GET("/", info())
		router.POST("/start", noOp)
		router.POST("/move", HandleMove)
		router.POST("/end", noOp)
		router.ServeHTTP(w, r)
	})
}

/**
* Handlers
**/
func info() spinhttp.RouterHandle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var info = models.BattlesnakeInfoResponse{
			APIVersion: "1",
			Author:     "",        // TODO: Your Battlesnake username
			Color:      "#888888", // TODO: Choose color
			Head:       "default", // TODO: Choose head
			Tail:       "default", // TODO: Choose tail
		}
		if err := json.NewEncoder(w).Encode(info); err != nil {
			http.Error(w, internalServiceErrMsg, http.StatusInternalServerError)
			log.Printf("failed to encode: %s", err)
			return
		}
		w.Header().Set(contentType, applicationJson)
	}
}

func HandleMove(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	state := models.GameState{}
	err := json.NewDecoder(r.Body).Decode(&state)
	if err != nil {
		log.Printf("ERROR: Failed to decode move json, %s", err)
		return
	}

	response := move(state)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("ERROR: Failed to encode move response, %s", err)
		return
	}
}

/**
* Core logic
**/
func move(state models.GameState) models.BattlesnakeMoveResponse {

	isMoveSafe := map[string]bool{
		"up":    true,
		"down":  true,
		"left":  true,
		"right": true,
	}

	// We've included code to prevent your Battlesnake from moving backwards
	myHead := state.You.Body[0] // Coordinates of your head
	myNeck := state.You.Body[1] // Coordinates of your "neck"

	if myNeck.X < myHead.X { // Neck is left of head, don't move left
		isMoveSafe["left"] = false

	} else if myNeck.X > myHead.X { // Neck is right of head, don't move right
		isMoveSafe["right"] = false

	} else if myNeck.Y < myHead.Y { // Neck is below head, don't move down
		isMoveSafe["down"] = false

	} else if myNeck.Y > myHead.Y { // Neck is above head, don't move up
		isMoveSafe["up"] = false
	}

	// TODO: Step 1 - Prevent your Battlesnake from moving out of bounds
	// boardWidth := state.Board.Width
	// boardHeight := state.Board.Height

	// TODO: Step 2 - Prevent your Battlesnake from colliding with itself
	// mybody := state.You.Body

	// TODO: Step 3 - Prevent your Battlesnake from colliding with other Battlesnakes
	// opponents := state.Board.Snakes

	// Are there any safe moves left?
	safeMoves := []string{}
	for move, isSafe := range isMoveSafe {
		if isSafe {
			safeMoves = append(safeMoves, move)
		}
	}

	if len(safeMoves) == 0 {
		log.Printf("MOVE %d: No safe moves detected! Moving down\n", state.Turn)
		return models.BattlesnakeMoveResponse{Move: "down"}
	}

	// Choose a random move from the safe ones
	nextMove := safeMoves[rand.Intn(len(safeMoves))]

	// TODO: Step 4 - Move towards food instead of random, to regain health and survive longer
	// food := state.Board.Food

	log.Printf("MOVE %d: %s\n", state.Turn, nextMove)
	return models.BattlesnakeMoveResponse{Move: nextMove}
}

func noOp(w http.ResponseWriter, r *http.Request, p httprouter.Params) {}

func main() {}
