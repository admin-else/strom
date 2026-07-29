package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// BlockedServers fetches the list of SHA-1 hashes of blocked server IPs from
// Mojang's session server. The returned strings are hex-encoded hashes.
func BlockedServers() (s []string, err error) {
	resp, err := http.Get("https://sessionserver.mojang.com/blockedservers")
	if err != nil {
		return
	}
	defer resp.Body.Close()
	d, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	s = strings.Split(string(d), "\n")
	return
}

type ProfileProperty struct {
	Name      string `json:"name"`
	Signature string `json:"signature,omitempty"`
	Value     string `json:"value"`
}

type ProfileResp struct {
	Id         string            `json:"id"`
	Name       string            `json:"name"`
	Legacy     bool              `json:"legacy"`
	Properties []ProfileProperty `json:"properties"`
}

// HasJoined checks with Mojang's session server whether a player has joined a
// given server. name is the player name, server is the server hash, and ip is
// the optional player IP address. Returns the player's profile on success.
func HasJoined(name, server, ip string) (p ProfileResp, err error) {
	url := "https://sessionserver.mojang.com/session/minecraft/hasJoined?username=" + name + "&serverId=" + server
	if ip != "" {
		url += "&ip=" + ip
	}
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		err = fmt.Errorf("bad status code %v", resp.StatusCode)
		return
	}
	err = json.NewDecoder(resp.Body).Decode(&p)
	return
}

type MojangPublicKeysResp struct {
	ProfilePropertyKeys   []string `json:"profilePropertyKeys"`
	PlayerCertificateKeys []string `json:"playerCertificateKeys"`

	// AuthenticationKeys are undocumented by wiki.vg but they exist as of the time of writing so ill include them
	AuthenticationKeys []string `json:"authenticationKeys"`
}

// MojangPublicKeys fetches Mojang's public keys for profile property signing and
// player certificate verification from the session server.
func MojangPublicKeys() (keys MojangPublicKeysResp, err error) {
	resp, err := http.Get("https://sessionserver.mojang.com/session/minecraft/profile/public_key_by_token")
	if err != nil {
		return
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&keys)
	return
}

// Profile fetches a player's profile by UUID from Mojang's session server.
// If unsigned is false, signed profile properties are included in the response.
func Profile(id string, unsinged bool) (p ProfileResp, err error) {
	url := "https://sessionserver.mojang.com/session/minecraft/profile/" + id
	if !unsinged {
		url += "?unsigned=false"
	}
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&p)
	return
}

type IdName struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// NameToId resolves a Minecraft player name to a UUID and canonical name by
// making an HTTP GET to Mojang's API.
func NameToId(name string) (idName IdName, err error) {
	// may also use these:
	//https://api.mojang.com/users/profiles/minecraft/<player name>
	//https://api.minecraftservices.com/minecraft/profile/lookup/name/<player name>
	//https://api.mojang.com/minecraft/profile/lookup/name/<player name>
	resp, err := http.Get("https://api.minecraftservices.com/minecraft/profile/lookup/name/" + name)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&idName)
	return
}

// NameToIdBulk resolves multiple Minecraft player names to UUIDs and canonical
// names in a single HTTP POST to Mojang's API.
func NameToIdBulk(names []string) (idNames []IdName, err error) {
	b := bytes.NewBuffer(nil)
	err = json.NewEncoder(b).Encode(names)
	if err != nil {
		return
	}
	// may also use these:
	//https://api.mojang.com/profiles/minecraft
	//https://api.minecraftservices.com/minecraft/profile/lookup/bulk/byname
	//https://api.mojang.com/minecraft/profile/lookup/bulk/byname
	r, err := http.NewRequest("POST", "https://api.minecraftservices.com/minecraft/profile/lookup/bulk/byname", b)
	if err != nil {
		return
	}
	r.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&idNames)
	return
}
