package token

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

type captchaResponse struct {
	Success     bool
	Credit      bool
	Hostname    string
	ChallengeTs string `json:"challenge_ts"`
}

func ValidateToken(token string) (bool, error) {
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

	cResp := captchaResponse{}
	err = json.Unmarshal(body, &cResp)
	if err != nil {
		return false, err
	}

	log.Printf("%#v", cResp)

	return cResp.Success, nil
}
