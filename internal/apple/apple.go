package apple

import (
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const authTokenUrl string = "https://appleid.apple.com/auth/token"

type SignInWithApple struct {
	secret      *ecdsa.PrivateKey
	keyId       string
	teamId      string
	clientId    string
	servicesId  string
	redirectUri string
}

// New
func New(filename, keyId, teamId, clientId, servicesId, redirectUri string) (*SignInWithApple, error) {
	s := &SignInWithApple{
		keyId:       keyId,
		teamId:      teamId,
		clientId:    clientId,
		servicesId:  servicesId,
		redirectUri: redirectUri,
	}

	// load apple .p8 auth key file
	if err := s.loadAuthKey(filename); err != nil {
		return nil, err
	}

	return s, nil
}

type AuthTokenResponse struct {
	Error        string `json:"error"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	IDToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

func (s *SignInWithApple) do(form url.Values) (*AuthTokenResponse, error) {
	var request *http.Request
	var err error
	if request, err = http.NewRequest("POST", authTokenUrl, strings.NewReader(form.Encode())); err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var response *http.Response
	if response, err = http.DefaultClient.Do(request); nil != err {
		return nil, err
	}
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	res := &AuthTokenResponse{}
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}

	if res.Error != "" {
		return res, errors.New(res.Error)
	}
	return res, nil
}

func (s *SignInWithApple) AuthTokenWithApp(code string) (*AuthTokenResponse, error) {
	clientSecret, err := s.clientSecret()
	if err != nil {
		return nil, err
	}

	form := url.Values{}
	form.Set("client_id", s.clientId)
	form.Set("client_secret", clientSecret)
	form.Set("code", code)
	form.Set("grant_type", "authorization_code")
	form.Set("redirect_uri", s.redirectUri)

	return s.do(form)
}

func (s *SignInWithApple) AuthTokenWithWeb(code string) (*AuthTokenResponse, error) {
	clientSecret, err := s.clientSecret()
	if err != nil {
		return nil, err
	}

	form := url.Values{}
	form.Set("client_id", s.servicesId)
	form.Set("client_secret", clientSecret)
	form.Set("code", code)
	form.Set("grant_type", "authorization_code")
	form.Set("redirect_uri", s.redirectUri)

	return s.do(form)
}

func (s *SignInWithApple) RefreshToken(refreshToken string) (*AuthTokenResponse, error) {
	clientSecret, err := s.clientSecret()
	if err != nil {
		return nil, err
	}

	form := url.Values{}
	form.Set("client_id", s.clientId)
	form.Set("client_secret", clientSecret)
	form.Set("refresh_token", refreshToken)
	form.Set("grant_type", "refresh_token")

	return s.do(form)
}
