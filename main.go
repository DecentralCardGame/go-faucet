package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"

	"github.com/DecentralCardGame/go-faucet/cardchain"
	"github.com/DecentralCardGame/go-faucet/cardchain/client"
	"github.com/DecentralCardGame/go-faucet/config"
	"github.com/DecentralCardGame/go-faucet/payload"
	"github.com/DecentralCardGame/go-faucet/token"
	"github.com/joho/godotenv"
)

func handleClaimTokens(w http.ResponseWriter, r *http.Request) {
	log.Print("Endpoint Hit: ClaimTokens")
	w.Header().Set("Content-Type", "application/json")
	enableCors(&w)

	pl := payload.Payload{}

	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		handleInternalServerError(w, err)
		return
	}

	err = json.Unmarshal(body, &pl)
	if err != nil {
		http.Error(
			w,
			"Invalid json",
			http.StatusBadRequest,
		)
		return
	}

	if !pl.Verify(w) {
		return
	}

	isValid, err := token.ValidateToken(pl.Token)
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

	cResp, err := cardchain.CreateUser(
		config.Config().BlockchainUser,
		pl.Alias,
		pl.Address,
	)

	if err != nil {
		handleInternalServerError(w, err)
		return
	}

	slog.Info("%#v", cResp)

	if cResp.Code != 0 {
		http.Error(
			w,
			fmt.Sprintf(
				"Cardchain responded with code %d: %s",
				cResp.Code,
				cResp.RawLog,
			),
			http.StatusForbidden,
		)
	}
}

func handleInternalServerError(w http.ResponseWriter, err error) {
	http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
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

	err = config.FromEnv()
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	err = client.WaitForChain(config.Config().ClientConfig())
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
	log.Print("Client instanciated")

	handleRequests()
}
