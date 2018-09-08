package aux

import (
	"context"
)

// A Corporation represents an EVE Online corporation.
type Corporation struct {
	// ID
	id int64

	// Name
	name string

	members lazyList
}

func (corp *Corporation) Name() string {
	return corp.name
}

func (corp *Corporation) Members() ([]*Character, error) {
	if corp.members.loaded() {
		return corp.members.load([]*Character), nil
	}

	lst, err := esi.Corporations.Members(context.Background(), corp.ID)
	if err != nil {
		return nil, err
	}

}
