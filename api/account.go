package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/admin-else/strom/crypto"
	"github.com/google/uuid"
)

type Account struct {
	Name string
	Uuid uuid.UUID
	Ygg  string
}

const (
	BaseUrlApi = "https://api.minecraftservices.com"
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
	r.Header.Set("Authorization", "Bearer "+a.Ygg)

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
	err = a.doMcApiRequest("POST", BaseUrlApi+"/player/certificates", nil, &keys)
	return
}

func (a *Account) JoinServer(serverId string) (err error) {
	body, err := json.Marshal(struct {
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
	resp, err := http.Post("https://sessionserver.mojang.com/session/minecraft/join", "application/json", bytes.NewReader(body))
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

type NameChangResp struct {
	ChangedAt         string `json:"changedAt"`
	CreatedAt         string `json:"createdAt"`
	NameChangeAllowed bool   `json:"nameChangeAllowed"`
}

func (a *Account) NameChangeInfo() (r NameChangResp, err error) {
	err = a.doMcApiRequest("GET", "https://api.minecraftservices.com/minecraft/profile/namechange", nil, &r)
	return
}

func NewOfflineAccount(name string) *Account {
	return &Account{
		Name: name,
		Uuid: crypto.FromOfflinePlayer(name),
	}
}

// PlayerProfile represents the player's profile information
type PlayerProfile struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Skins []struct {
		ID      string `json:"id"`
		State   string `json:"state"`
		URL     string `json:"url"`
		Variant string `json:"variant"`
	} `json:"skins"`
	Capes []struct {
		ID    string `json:"id"`
		State string `json:"state"`
		URL   string `json:"url"`
		Alias string `json:"alias"`
	} `json:"capes"`
}

// QueryProfile queries the player's profile
func (a *Account) QueryProfile() (profile PlayerProfile, err error) {
	err = a.doMcApiRequest("GET", BaseUrlApi+"/minecraft/profile", nil, &profile)
	return
}

// PlayerAttributes represents the player's attributes
type PlayerAttributes struct {
	Privileges struct {
		OnlineChat struct {
			Enabled bool `json:"enabled"`
		} `json:"onlineChat"`
		MultiplayerServer struct {
			Enabled bool `json:"enabled"`
		} `json:"multiplayerServer"`
		MultiplayerRealms struct {
			Enabled bool `json:"enabled"`
		} `json:"multiplayerRealms"`
		Telemetry struct {
			Enabled bool `json:"enabled"`
		} `json:"telemetry"`
		OptionalTelemetry struct {
			Enabled bool `json:"enabled"`
		} `json:"optionalTelemetry"`
	} `json:"privileges"`
	ProfanityFilterPreferences struct {
		ProfanityFilterOn bool `json:"profanityFilterOn"`
	} `json:"profanityFilterPreferences"`
	BanStatus struct {
		BannedScopes map[string]struct {
			BanID         string `json:"banId"`
			Expires       string `json:"expires,omitempty"`
			Reason        string `json:"reason"`
			ReasonMessage string `json:"reasonMessage"`
		} `json:"bannedScopes"`
	} `json:"banStatus"`
}

// QueryAttributes queries the player's attributes
func (a *Account) QueryAttributes() (attrs PlayerAttributes, err error) {
	err = a.doMcApiRequest("GET", BaseUrlApi+"/player/attributes", nil, &attrs)
	return
}

// ModifyAttributesRequest represents the request to modify player attributes
type ModifyAttributesRequest struct {
	ProfanityFilterPreferences struct {
		ProfanityFilterOn bool `json:"profanityFilterOn"`
	} `json:"profanityFilterPreferences"`
}

// ModifyAttributes modifies the player's attributes
func (a *Account) ModifyAttributes(req ModifyAttributesRequest) (attrs PlayerAttributes, err error) {
	err = a.doMcApiRequest("POST", BaseUrlApi+"/player/attributes", req, &attrs)
	return
}

// BlockedUsers represents the list of blocked users
type BlockedUsers struct {
	BlockedProfiles []string `json:"blockedProfiles"`
}

// GetBlocklist gets the list of blocked users
func (a *Account) GetBlocklist() (blocklist BlockedUsers, err error) {
	err = a.doMcApiRequest("GET", BaseUrlApi+"/privacy/blocklist", nil, &blocklist)
	return
}

// CheckNameAvailability checks if a name is available
type NameAvailability struct {
	Status string `json:"status"` // DUPLICATE, AVAILABLE, NOT_ALLOWED
}

// CheckNameAvailability checks if a name is available
func (a *Account) CheckNameAvailability(name string) (avail NameAvailability, err error) {
	err = a.doMcApiRequest("GET", BaseUrlApi+"/minecraft/profile/name/"+name+"/available", nil, &avail)
	return
}

// ChangeName changes the player's name
func (a *Account) ChangeName(newName string) (profile PlayerProfile, err error) {
	err = a.doMcApiRequest("PUT", BaseUrlApi+"/minecraft/profile/name/"+newName, nil, &profile)
	return
}

// ChangeSkinRequest represents the request to change skin
type ChangeSkinRequest struct {
	Variant string `json:"variant"` // "classic" or "slim"
	URL     string `json:"url"`
}

// ChangeSkin changes the player's skin
func (a *Account) ChangeSkin(req ChangeSkinRequest) (profile PlayerProfile, err error) {
	err = a.doMcApiRequest("POST", BaseUrlApi+"/minecraft/profile/skins", req, &profile)
	return
}

// ResetSkin resets the player's skin to default
func (a *Account) ResetSkin() (profile PlayerProfile, err error) {
	err = a.doMcApiRequest("DELETE", BaseUrlApi+"/minecraft/profile/skins/active", nil, &profile)
	return
}

// HideCape hides the player's active cape
func (a *Account) HideCape() (profile PlayerProfile, err error) {
	err = a.doMcApiRequest("DELETE", BaseUrlApi+"/minecraft/profile/capes/active", nil, &profile)
	return
}

// ShowCapeRequest represents the request to show a cape
type ShowCapeRequest struct {
	CapeID string `json:"capeId"`
}

// ShowCape shows a specific cape
func (a *Account) ShowCape(capeID string) (profile PlayerProfile, err error) {
	req := ShowCapeRequest{CapeID: capeID}
	err = a.doMcApiRequest("PUT", BaseUrlApi+"/minecraft/profile/capes/active", req, &profile)
	return
}

// CheckGiftCodeValidity checks if a gift code is valid
func (a *Account) CheckGiftCodeValidity(giftCode string) (valid bool, err error) {
	url := BaseUrlApi + "/productvoucher/giftcode"
	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}

	r.Header.Set("Authorization", "Bearer "+a.Ygg)

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// Returns HTTP 200 or 204 if valid, 404 if invalid
	if resp.StatusCode == 200 || resp.StatusCode == 204 {
		return true, nil
	}
	return false, nil
}
