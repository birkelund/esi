package esi

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestFleetsEndpoint_GetCharacterFleetInfo(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v1/characters/42/fleet/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `
			{
				"fleet_id": 1234567890,
				"role": "fleet_commander",
				"squad_id": -1,
				"wing_id": -1
			}
		`)
	})

	fleetInfo, _, err := client.Fleets.GetCharacterFleet(context.Background(), 42)
	if err != nil {
		t.Errorf("Fleets.GetCharacterFleetInfo returned error: %v", err)
	}

	want := &CharacterFleetResponse{
		FleetID: Int(1234567890),
		Role:    String("fleet_commander"),
		SquadID: Int(-1),
		WingID:  Int(-1),
	}
	if !reflect.DeepEqual(fleetInfo, want) {
		t.Errorf("Fleets.GetCharacterFleetInfo returned %+v, want %+v", fleetInfo, want)
	}
}

func TestFleetsEndpoint_GetFleetInfo(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v1/fleets/42/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `
			{
				"is_free_move": false,
				"is_registered": false,
				"is_voice_enabled": false,
				"motd": "This is an <b>awesome</b> fleet!"
		  	}
		`)
	})

	fleetInfo, _, err := client.Fleets.Get(context.Background(), 42)
	if err != nil {
		t.Errorf("Fleets.GetFleetInfo returned error: %v", err)
	}

	want := &FleetResponse{
		IsFreeMove:     Bool(false),
		IsRegistered:   Bool(false),
		IsVoiceEnabled: Bool(false),
		MOTD:           String("This is an <b>awesome</b> fleet!"),
	}
	if !reflect.DeepEqual(fleetInfo, want) {
		t.Errorf("Fleets.GetFleetInfo returned %+v, want %+v", fleetInfo, want)
	}
}

func TestFleetsEndpoint_Update(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v1/fleets/42/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testBody(t, r, `{"is_free_move":true,"motd":"some motd"}`+"\n")
	})

	_, err := client.Fleets.Update(context.Background(), 42, &FleetSettings{
		IsFreeMove: true,
		MOTD:       "some motd",
	})
	if err != nil {
		t.Errorf("Fleets.Update returned error: %v", err)
	}
}
