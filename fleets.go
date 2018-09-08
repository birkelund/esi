package esi

import (
	"context"
	"fmt"
)

// FleetsEndpoint handles communication with the fleets related methods of
// the ESI API.
type FleetsEndpoint endpoint

// CharacterFleetResponse holds details about the character's fleet.
type CharacterFleetResponse struct {
	FleetID *int    `json:"fleet_id,omitempty"`
	Role    *string `json:"role,omitempty"`
	SquadID *int    `json:"squad_id,omitempty"`
	WingID  *int    `json:"wing_id,omitempty"`
}

// GetCharacterFleet returns the fleet ID the is in, if any.
func (e *FleetsEndpoint) GetCharacterFleet(ctx context.Context, cid int) (*CharacterFleetResponse, *Response, error) {
	u := fmt.Sprintf("v1/characters/%d/fleet/", cid)

	req, err := e.api.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	characterFleetResponse := new(CharacterFleetResponse)
	resp, err := e.api.Do(ctx, req, characterFleetResponse)
	if err != nil {
		return nil, resp, err
	}

	return characterFleetResponse, resp, nil
}

// FleetResponse holds details about a fleet.
type FleetResponse struct {
	IsFreeMove     *bool   `json:"is_free_move,omitempty"`
	IsRegistered   *bool   `json:"is_registered,omitempty"`
	IsVoiceEnabled *bool   `json:"is_voice_enabled,omitempty"`
	MOTD           *string `json:"motd,omitempty"`
}

// Get returns details about a fleet.
func (e *FleetsEndpoint) Get(ctx context.Context, fid int) (*FleetResponse, *Response, error) {
	u := fmt.Sprintf("v1/fleets/%d/", fid)

	req, err := e.api.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	fleetResponse := new(FleetResponse)
	resp, err := e.api.Do(ctx, req, fleetResponse)
	if err != nil {
		return nil, resp, err
	}

	return fleetResponse, resp, nil
}

// FleetSettings holds what to update for a fleet.
type FleetSettings struct {
	IsFreeMove bool   `json:"is_free_move,omitempty"`
	MOTD       string `json:"motd,omitempty"`
}

// Update updates settings about a fleet.
func (e *FleetsEndpoint) Update(ctx context.Context, fid int, settings *FleetSettings) (*Response, error) {
	u := fmt.Sprintf("v1/fleets/%d/", fid)

	req, err := e.api.NewRequest("PUT", u, settings)
	if err != nil {
		return nil, err
	}

	return e.api.Do(ctx, req, nil)
}

// FleetMember holds details of a fleet member.
type FleetMember struct {
	CharacterID    *int       `json:"character_id,omitempty"`
	JoinTime       *Timestamp `json:"join_time,omitempty"`
	Role           *string    `json:"role,omitempty"`
	RoleName       *string    `json:"role_name,omitempty"`
	ShipTypeID     *string    `json:"ship_type_id,omitempty"`
	SolarSystemID  *int       `json:"solar_system_id,omitempty"`
	SquadID        *int       `json:"squad_id,omitempty"`
	StationID      *int       `json:"station_id,omitempty"`
	TakesFleetWarp *bool      `json:"takes_fleet_warp,omitempty"`
	WingID         *int       `json:"wing_id,omitempty"`
}

// FleetMembersResponse holds information about fleet members.
type FleetMembersResponse []*FleetMember

// GetMembers returns information about fleet members.
func (e *FleetsEndpoint) GetMembers(ctx context.Context, fid int, opt *I18NOptions) (FleetMembersResponse, *Response, error) {
	u := fmt.Sprintf("v1/fleets/%d/members/", fid)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := e.api.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var fleetMembersResponse FleetMembersResponse
	resp, err := e.api.Do(ctx, req, &fleetMembersResponse)
	if err != nil {
		return nil, resp, err
	}

	return fleetMembersResponse, resp, nil
}

// FleetInvitation holds details of a fleet invitation.
type FleetInvitation struct {
	CharacterID int    `json:"character_id,omitempty"`
	Role        string `json:"role,omitempty"`
	SquadID     int    `json:"squad_id,omitempty"`
	WingID      int    `json:"wing_id,omitempty"`
}

// Invite invites a character into the fleet. If a character has a CSPA charge
// set it is not possible to invite them to the fleet using ESI.
func (e *FleetsEndpoint) Invite(ctx context.Context, fid int, invitation *FleetInvitation) (*Response, error) {
	u := fmt.Sprintf("v1/fleets/%d/members/", fid)

	req, err := e.api.NewRequest("POST", u, invitation)
	if err != nil {
		return nil, err
	}

	return e.api.Do(ctx, req, nil)
}

// Kick kicks a fleet member.
func (e *FleetsEndpoint) Kick(ctx context.Context, fid int, cid int) (*Response, error) {
	u := fmt.Sprintf("v1/fleets/%d/members/%d/", fid, cid)

	req, err := e.api.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return e.api.Do(ctx, req, nil)
}

// FleetMemberMovement holds details of the fleet member movement.
type FleetMemberMovement struct {
	Role    string `json:"role,omitempty"`
	SquadID int    `json:"squad_id,omitempty"`
	WingID  int    `json:"wing_id,omitempty"`
}

// Move moves a fleet member between squads and wings.
func (e *FleetsEndpoint) Move(ctx context.Context, fid int, cid int, m *FleetMemberMovement) (*Response, error) {
	u := fmt.Sprintf("v1/fleets/%d/members/%d/", fid, cid)

	req, err := e.api.NewRequest("PUT", u, m)
	if err != nil {
		return nil, err
	}

	return e.api.Do(ctx, req, nil)
}

// DeleteSquad deletes a fleet squad. Only empty squads can be deleted.
func (e *FleetsEndpoint) DeleteSquad(ctx context.Context, fid int, sid int) (*Response, error) {
	u := fmt.Sprintf("v1/fleets/%d/squads/%d/", fid, sid)

	req, err := e.api.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return e.api.Do(ctx, req, nil)
}

// RenameSquad renames a fleet squad.
func (e *FleetsEndpoint) RenameSquad(ctx context.Context, fid int, sid int, name string) (*Response, error) {
	u := fmt.Sprintf("v1/fleets/%d/squads/%d/", fid, sid)

	req, err := e.api.NewRequest("PUT", u, struct {
		Name string `json:"name"`
	}{name})
	if err != nil {
		return nil, err
	}

	return e.api.Do(ctx, req, nil)
}

// FleetSquad contains info about a fleet squad.
type FleetSquad struct {
	ID   *int    `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
}

// FleetWing contains info about a fleet wing, including a list of squads.
type FleetWing struct {
	ID     *int          `json:"id,omitempty"`
	Name   *string       `json:"name,omitempty"`
	Squads []*FleetSquad `json:"squads,omitempty"`
}

// FleetWingsResponse holds information about wings in a fleet.
type FleetWingsResponse []*FleetWing

// GetWings returns information about wings in a fleet.
func (e *FleetsEndpoint) GetWings(ctx context.Context, fid int, opt *I18NOptions) (FleetWingsResponse, *Response, error) {
	u := fmt.Sprintf("v1/fleets/%d/wings/", fid)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := e.api.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var fleetWingsResponse FleetWingsResponse
	resp, err := e.api.Do(ctx, req, &fleetWingsResponse)
	if err != nil {
		return nil, resp, err
	}

	return fleetWingsResponse, resp, nil
}

// CreateWing creates a new wing in a fleet.
func (e *FleetsEndpoint) CreateWing(ctx context.Context, fid int) (int, *Response, error) {
	u := fmt.Sprintf("v1/fleets/%d/wings/", fid)

	req, err := e.api.NewRequest("PUT", u, nil)
	if err != nil {
		return -1, nil, err
	}

	var v struct {
		WingID int `json:"wing_id"`
	}

	resp, err := e.api.Do(ctx, req, &v)
	if err != nil {
		return -1, resp, err
	}

	return v.WingID, resp, nil
}

// DeleteWing deletes a wing in a fleet. Only empty wings may be deleted. The
// wing may contain squads but the squads must be empty.
func (e *FleetsEndpoint) DeleteWing(ctx context.Context, fid int, wid int) (*Response, error) {
	u := fmt.Sprintf("v1/fleets/%d/wings/%d/", fid, wid)

	req, err := e.api.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return e.api.Do(ctx, req, nil)
}

// RenameWing renames a fleet wing.
func (e *FleetsEndpoint) RenameWing(ctx context.Context, fid int, wid int, name string) (*Response, error) {
	u := fmt.Sprintf("v1/fleets/%d/wings/%d/", wid, wid)

	req, err := e.api.NewRequest("PUT", u, struct {
		Name string `json:"name"`
	}{name})
	if err != nil {
		return nil, err
	}

	return e.api.Do(ctx, req, nil)
}

// CreateSquad creates a new squad in a fleet.
func (e *FleetsEndpoint) CreateSquad(ctx context.Context, fid int) (int, *Response, error) {
	u := fmt.Sprintf("v1/fleets/%d/wings/squads/", fid)

	req, err := e.api.NewRequest("PUT", u, nil)
	if err != nil {
		return -1, nil, err
	}

	var v struct {
		SquadID int `json:"squad_id"`
	}

	resp, err := e.api.Do(ctx, req, &v)
	if err != nil {
		return -1, resp, err
	}

	return v.SquadID, resp, nil
}
