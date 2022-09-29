package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/DecentralCardGame/go-faucet/cardchain"
	"github.com/DecentralCardGame/go-faucet/token"
	"github.com/joho/godotenv"
)

var (
	CHAIN_USER string = "alice"
)

type Payload struct {
	Token   string
	Address string
}

func handleClaimTokens(w http.ResponseWriter, r *http.Request) {
	log.Print("Endpoint Hit: ClaimTokens")
	w.Header().Set("Content-Type", "application/json")
	enableCors(&w)

	payload := Payload{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		handleInternalServerError(w, err)
		return
	}

	err = json.Unmarshal(body, &payload)
	if err != nil {
		handleInternalServerError(w, err)
		return
	}

	isValid, err := token.ValidateToken(payload.Token)
	if err != nil {
		handleInternalServerError(w, err)
		return
	}

	if !isValid {
		http.Error(
			w,
			"User failed captcha",
			http.StatusForbidden,
		)
		return
	}

	cResp, err := cardchain.CreateUser(CHAIN_USER, "newUser", payload.Address)
	if err != nil {
		handleInternalServerError(w, err)
		return
	}

	if cResp.Code != 0 {
		http.Error(
			w,
			fmt.Sprintf("Cardchain responded with code %d: %s", cResp.Code, cResp.RawLog),
			http.StatusForbidden,
		)
	}
}

func handleInternalServerError(w http.ResponseWriter, err error) {
	http.Error(w, "Internal server error", http.StatusInternalServerError)
	log.Printf("Error: %s", err.Error())
}

func handleRequests() {
	http.HandleFunc("/claimTokens", handleClaimTokens)
	log.Print("Server running at port 4500")
	log.Fatal(http.ListenAndServe(":4500", nil))
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
	handleRequests()
}
