package strom

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Account struct {
	Name string
	Uuid uuid.UUID
	Ygg  string
}

const (
	ApiBase = "https://api.minecraftservices.com"
)

func (a *Account) doMcApiRequest(method string, url string, from any, to any) (err error) {
	var b []byte
	if method != "GET" {
		b, err = json.Marshal(from)
		if err != nil {
			return err
		}
	}

	r, err := http.NewRequest(method, url, bytes.NewReader(b))
	if err != nil {
		return err
	}

	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	r.Header.Set("x-xbl-contract-version", "1")
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.Ygg))

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}

	if !(resp.StatusCode >= 200 && resp.StatusCode <= 299) {
		return fmt.Errorf("bad status code %v", resp.StatusCode)
	}
	b, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, to)
	if err != nil {
		return err
	}
	return nil
}

type PlayerKeys struct {
	ExpiresAt time.Time `json:"expiresAt"`
	KeyPair   struct {
		PublicKey  string `json:"publicKey"`
		PrivateKey string `json:"privateKey"`
	} `json:"keyPair"`
	PublicKeySignature   string    `json:"publicKeySignature"`
	PublicKeySignatureV2 string    `json:"publicKeySignatureV2"`
	RefreshedAfter       time.Time `json:"refreshedAfter"`
}

func (a *Account) FetchKeys() (keys PlayerKeys, err error) {
	err = a.doMcApiRequest("POST", ApiBase+"/player/certificates", nil, &keys)
	return
}

func (a *Account) JoinServer(serverId string) (err error) {
	var body []byte
	body, err = json.Marshal(struct {
		AccessToken     string `json:"accessToken"`
		SelectedProfile string `json:"selectedProfile"`
		ServerId        string `json:"serverId"`
	}{
		AccessToken:     a.Ygg,
		SelectedProfile: strings.ReplaceAll(a.Uuid.String(), "-", ""),
		ServerId:        serverId,
	})
	if err != nil {
		return
	}
	var resp *http.Response
	resp, err = http.Post("https://sessionserver.mojang.com/session/minecraft/join", "application/json", bytes.NewReader(body))
	if err != nil {
		return
	}
	if resp.StatusCode != 204 {
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return
		}
		err = fmt.Errorf("bad client code %v body %v", resp.StatusCode, string(body))
		return
	}
	return
}
