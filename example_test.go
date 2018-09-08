package esi_test

/*
import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"golang.org/x/oauth2"

	"github.com/gregjones/httpcache"

	"corpus.space/esi"
)

const PascalCharacterID = 440659656

func ExampleConfig_simple() {
	type Character struct {
		Name string `json:"name"`
	}

	var character Character

	if err := esi.Get(context.Background(), &character, "/v4/characters/%d/", PascalCharacterID); err != nil {
		log.Fatal(err)
	}

	fmt.Println(character.Name)

	// Output:
	// Pascal d'Mier
}

func ExampleConfig_custom_transport() {
	ctx := context.Background()

	ctx = context.WithValue(ctx,
		esi.HTTPClient, httpcache.NewMemoryCacheTransport().Client(),
	)

	type Character struct {
		Name string `json:"name"`
		Sex  string `json:"gender"`
	}

	var character Character

	if err := esi.Get(ctx, &character, "/v4/characters/%d/", PascalCharacterID); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s (%s)\n", character.Name, character.Sex)

	// Output:
	// Pascal d'Mier (male)
}

func ExampleConfig_authenticated() {
	ctx := context.Background()

	oauth2conf := &oauth2.Config{
		ClientID:     "YOUR_CLIENT_ID",
		ClientSecret: "YOUR_CLIENT_SECRET",

		Scopes: []string{
			"esi-wallet.read_character_wallet.v1",
		},

		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://login.eveonline.com/oauth/authorize",
			TokenURL: "https://login.eveonline.com/oauth/token",
		},
	}

	url := oauth2conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	fmt.Printf("Visit the URL for the auth dialog: %v", url)

	fmt.Printf("Code: ")
	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatal(err)
	}

	tok, err := oauth2conf.Exchange(ctx, code)
	if err != nil {
		log.Fatal(err)
	}

	if !tok.Valid() {
		log.Fatal("invalid token received")
	}

	authenticatedClient := oauth2conf.Client(ctx, tok)

	// check and store who we authenticated as
	resp, err := authenticatedClient.Get("https://login.eveonline.com/oauth/verify")
	if err != nil {
		log.Fatal(err)
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var info struct {
		ID        int    `json:"CharacterID"`
		OwnerHash string `json:"CharacterOwnerHash"`
	}

	if err := json.Unmarshal(buf, &info); err != nil {
		log.Fatal(err)
		return
	}

	ctx = context.WithValue(ctx,
		esi.HTTPClient, oauth2conf.Client(ctx, tok),
	)

	var balance float64

	if err := esi.Get(ctx, &balance, "/v1/characters/%d/wallet/", info.ID); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("wallet balance: %f\n", balance)
}
*/
