package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/DecentralCardGame/Cardchain/x/cardchain/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/ignite-hq/cli/ignite/pkg/cosmosclient"
	"github.com/joho/godotenv"
)

var (
	CHAIN_USER string = "alice"
)

type Payload struct {
	Token   string
	Address string
}

type CaptchaResponse struct {
	Success     bool
	Credit      bool
	Hostname    string
	ChallengeTs string `json:"challenge_ts"`
}

func getClient() (cosmosclient.Client, error) {
	config := sdktypes.GetConfig()
	config.SetBech32PrefixForAccount("cc", "ccpub")

	return cosmosclient.New(context.Background(), cosmosclient.WithAddressPrefix("cc"), cosmosclient.WithNodeAddress(os.Getenv("RPC_NODE")))
}

func getAddr(cosmos cosmosclient.Client, user string) (sdktypes.AccAddress, error) {
	address, err := cosmos.Address(user)
	if err != nil {
		return nil, err
	}
	return address, nil
}

func broadcastMsg(cosmos cosmosclient.Client, creator string, msg sdktypes.Msg) error {
	resp, err := cosmos.BroadcastTx(creator, msg)
	if err != nil {
		return err
	}

	log.Printf("%d", int(resp.Code))

	return nil
}

func createUser(creator string, alias string, userAddressString string) error {
	cosmos, err := getClient()
	if err != nil {
		return err
	}

	address, err := getAddr(cosmos, creator)
	if err != nil {
		return err
	}

	userAddr, err := sdktypes.AccAddressFromBech32(userAddressString)
	if err != nil {
		return err
	}

	msg := types.NewMsgCreateuser(
		address.String(),
		userAddr.String(),
		alias,
	)

	err = broadcastMsg(cosmos, creator, msg)

	return err
}

func isValidToken(token string) (bool, error) {
	data := url.Values{
		"secret":   {os.Getenv("SECRET_KEY")},
		"response": {token},
	}

	resp, err := http.PostForm(
		"https://hcaptcha.com/siteverify",
		data,
	)
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	captchaResponse := CaptchaResponse{}
	err = json.Unmarshal(body, &captchaResponse)
	if err != nil {
		return false, err
	}

	log.Printf("%#v", captchaResponse)

	return captchaResponse.Success, nil
}

func handleClaimTokens(w http.ResponseWriter, r *http.Request) {
	log.Print("Endpoint Hit: ClaimTokens")
	w.Header().Set("Content-Type", "application/json")
	enableCors(&w)

	payload := Payload{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	log.Print(r.Header.Get("Content-Type"))
	log.Print(string(body))
	err = json.Unmarshal(body, &payload)
	if err != nil {
		panic(err)
	}

	isValid, err := isValidToken(payload.Token)
	if err != nil {
		panic(err)
	}

	if !isValid {
		return
	}

	err = createUser(CHAIN_USER, "newUser", payload.Address)
	if err != nil {
		log.Fatal(err)
	}
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
