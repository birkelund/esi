package aux

import "time"

// A Character represents an EVE Online character.
type Character struct {
	lazyResource

	// The character ID
	ID int64

	// Character name
	Name string

	// The corporation this character is a member of
	Corporation *Corporation

	// The alliance this character's corporation is a member of, if any
	Alliance *Alliance

	// The faction this character belongs to
	Faction *Faction

	// The character's race
	Race string

	// The ancestry of the character
	Ancestry string

	// The bloodline of the character
	Bloodline string

	// Birthday of the character
	Birthday time.Time

	// The biography of the character
	Bio string

	// The gender of the character
	Gender string

	// The security status of the character
	SecurityStatus float64
}
