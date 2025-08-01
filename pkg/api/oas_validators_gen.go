// Code generated by ogen, DO NOT EDIT.

package api

import (
	"github.com/go-faster/errors"
)

func (s GetMultiplayersSummaryOrder) Validate() error {
	switch s {
	case "asc":
		return nil
	case "desc":
		return nil
	default:
		return errors.Errorf("invalid value: %v", s)
	}
}

func (s GetServerStatsByIDOKApplicationJSON) Validate() error {
	alias := ([]GetServerStatsByIDOKItem)(s)
	if alias == nil {
		return errors.New("nil is invalid value")
	}
	return nil
}

func (s GetServerStatsByIDOrder) Validate() error {
	switch s {
	case "asc":
		return nil
	case "desc":
		return nil
	default:
		return errors.Errorf("invalid value: %v", s)
	}
}

func (s GetServersByMultiplayerOKApplicationJSON) Validate() error {
	alias := ([]GetServersByMultiplayerOKItem)(s)
	if alias == nil {
		return errors.New("nil is invalid value")
	}
	return nil
}
