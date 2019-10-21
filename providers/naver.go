package providers

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/bitly/oauth2_proxy/api"
)

type NaverProvider struct {
	*ProviderData
}

func NewNaverProvider(p *ProviderData) *NaverProvider {
	p.ProviderName = "Naver"
	if p.LoginURL == nil || p.LoginURL.String() == "" {
		p.LoginURL, _ = url.Parse("https://nss.navercorp.com/nweauthorize")
	}
	if p.RedeemURL == nil || p.RedeemURL.String() == "" {
		p.RedeemURL, _ = url.Parse("https://nss-api.navercorp.com:5001/api/Auth/token")
	}
	if p.ValidateURL == nil || p.ValidateURL.String() == "" {
		p.ValidateURL, _ = url.Parse("https://nss-api.navercorp.com:5001/api/Auth/tokenInfo")
	}
	return &NaverProvider{ProviderData: p}
}

func (p *NaverProvider) GetLoginURL(redirectURI, finalRedirect string) string {
	params := p.LoginURL.Query()
	params.Set("redirect_uri", redirectURI)
	params.Set("client_id", p.ClientID)
	params.Set("svcId", p.ClientID)
	params.Set("response_type", "code")
	p.LoginURL.RawQuery = params.Encode()
	return p.LoginURL.String()
}

func (p *NaverProvider) Redeem(redirectURL, code string) (s *SessionState, err error) {
	if code == "" {
		err = errors.New("missing code")
		return
	}

	params := url.Values{}
	params.Add("svcId", p.ClientID)
	params.Add("code", code)

	req, err := http.NewRequest("POST", p.RedeemURL.String(), bytes.NewBufferString(params.Encode()))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var netClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := netClient.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		err = fmt.Errorf("got %d from %q %s", resp.StatusCode, p.RedeemURL.String(), body)
		return
	}

	// blindly try json and x-www-form-urlencoded
	var jsonResponse struct {
		AccessToken string `json:"access_token"`
	}
	err = json.Unmarshal(body, &jsonResponse)
	if err == nil {
		s = &SessionState{
			AccessToken: jsonResponse.AccessToken,
		}
		return
	}

	v, err := url.ParseQuery(string(body))
	if err != nil {
		return
	}
	if a := v.Get("access_token"); a != "" {
		s = &SessionState{AccessToken: a}
	} else {
		err = fmt.Errorf("no access token found %s", body)
	}
	return
}

func (p *NaverProvider) GetEmailAddress(s *SessionState) (string, error) {
	req, err := http.NewRequest("GET",
		p.ValidateURL.String()+"?access_token="+s.AccessToken+"&svcId="+p.ClientID, nil)
	if err != nil {
		log.Printf("failed building request %s", err)
		return "", err
	}
	json, err := api.RequestInsecure(req)
	if err != nil {
		log.Printf("failed making request %s", err)
		return "", err
	}
	return json.Get("user_mailaddr").String()
}
