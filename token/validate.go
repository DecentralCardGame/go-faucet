package token

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

type captchaResponse struct {
	Success     bool
	Credit      bool
	Hostname    string
	ChallengeTs string   `json:"challenge_ts"`
	ErrorCodes  []string `json:"error-codes"`
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

	if len(cResp.ErrorCodes) > 0 {
		return false, fmt.Errorf("Captcha responded with errors: %s", cResp.ErrorCodes)
	}

	return cResp.Success, nil
}
