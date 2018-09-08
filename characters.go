package esi

import (
	"context"
	"fmt"
)

// CharactersEndpoint handles communication with the characters related methods
// of the ESI API.
type CharactersEndpoint endpoint

// CharacterPublicInfo holds public information about a character.
type CharacterPublicInfo struct {
	AllianceID     *int       `json:"alliance_id,omitempty"`
	AncestryID     *int       `json:"ancestry_id,omitempty"`
	Birthday       *Timestamp `json:"birthday,omitempty"`
	BloodlineID    *int       `json:"bloodline_id,omitempty"`
	CorporationID  *int       `json:"corporation_id,omitempty"`
	Description    *string    `json:"description,omitempty"`
	FactionID      *int       `json:"faction_id,omitempty"`
	Gender         *string    `json:"gender,omitempty"`
	Name           *string    `json:"name,omitempty"`
	RaceID         *int       `json:"race_id,omitempty"`
	SecurityStatus *float64   `json:"security_status,omitempty"`
}

func (s CharacterPublicInfo) String() string {
	return Stringify(s)
}

// GetCharacter returns public information about a character.
func (e *CharactersEndpoint) GetCharacter(ctx context.Context, cid int) (*CharacterPublicInfo, *Response, error) {
	u := fmt.Sprintf("v1/characters/%d/", cid)

	req, err := e.api.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	characterPublicInfo := new(CharacterPublicInfo)
	resp, err := e.api.Do(ctx, req, characterPublicInfo)
	if err != nil {
		return nil, resp, err
	}

	return characterPublicInfo, resp, nil
}
