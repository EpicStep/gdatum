// Code generated by ogen, DO NOT EDIT.

package api

import (
	"time"

	"github.com/go-faster/errors"
)

type GetMultiplayersSummaryOKItem struct {
	// Name of a multiplayer.
	Name   string `json:"name"`
	Online int64  `json:"online"`
}

// GetName returns the value of Name.
func (s *GetMultiplayersSummaryOKItem) GetName() string {
	return s.Name
}

// GetOnline returns the value of Online.
func (s *GetMultiplayersSummaryOKItem) GetOnline() int64 {
	return s.Online
}

// SetName sets the value of Name.
func (s *GetMultiplayersSummaryOKItem) SetName(val string) {
	s.Name = val
}

// SetOnline sets the value of Online.
func (s *GetMultiplayersSummaryOKItem) SetOnline(val int64) {
	s.Online = val
}

type GetMultiplayersSummaryOrder string

const (
	GetMultiplayersSummaryOrderAsc  GetMultiplayersSummaryOrder = "asc"
	GetMultiplayersSummaryOrderDesc GetMultiplayersSummaryOrder = "desc"
)

// AllValues returns all GetMultiplayersSummaryOrder values.
func (GetMultiplayersSummaryOrder) AllValues() []GetMultiplayersSummaryOrder {
	return []GetMultiplayersSummaryOrder{
		GetMultiplayersSummaryOrderAsc,
		GetMultiplayersSummaryOrderDesc,
	}
}

// MarshalText implements encoding.TextMarshaler.
func (s GetMultiplayersSummaryOrder) MarshalText() ([]byte, error) {
	switch s {
	case GetMultiplayersSummaryOrderAsc:
		return []byte(s), nil
	case GetMultiplayersSummaryOrderDesc:
		return []byte(s), nil
	default:
		return nil, errors.Errorf("invalid value: %q", s)
	}
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (s *GetMultiplayersSummaryOrder) UnmarshalText(data []byte) error {
	switch GetMultiplayersSummaryOrder(data) {
	case GetMultiplayersSummaryOrderAsc:
		*s = GetMultiplayersSummaryOrderAsc
		return nil
	case GetMultiplayersSummaryOrderDesc:
		*s = GetMultiplayersSummaryOrderDesc
		return nil
	default:
		return errors.Errorf("invalid value: %q", data)
	}
}

// GetServerByIDNotFound is response for GetServerByID operation.
type GetServerByIDNotFound struct{}

func (*GetServerByIDNotFound) getServerByIDRes() {}

type GetServerByIDOK struct {
	Name     string    `json:"name"`
	URL      OptString `json:"url"`
	Gamemode OptString `json:"gamemode"`
	Lang     OptString `json:"lang"`
}

// GetName returns the value of Name.
func (s *GetServerByIDOK) GetName() string {
	return s.Name
}

// GetURL returns the value of URL.
func (s *GetServerByIDOK) GetURL() OptString {
	return s.URL
}

// GetGamemode returns the value of Gamemode.
func (s *GetServerByIDOK) GetGamemode() OptString {
	return s.Gamemode
}

// GetLang returns the value of Lang.
func (s *GetServerByIDOK) GetLang() OptString {
	return s.Lang
}

// SetName sets the value of Name.
func (s *GetServerByIDOK) SetName(val string) {
	s.Name = val
}

// SetURL sets the value of URL.
func (s *GetServerByIDOK) SetURL(val OptString) {
	s.URL = val
}

// SetGamemode sets the value of Gamemode.
func (s *GetServerByIDOK) SetGamemode(val OptString) {
	s.Gamemode = val
}

// SetLang sets the value of Lang.
func (s *GetServerByIDOK) SetLang(val OptString) {
	s.Lang = val
}

func (*GetServerByIDOK) getServerByIDRes() {}

// GetServerStatsByIDNotFound is response for GetServerStatsByID operation.
type GetServerStatsByIDNotFound struct{}

func (*GetServerStatsByIDNotFound) getServerStatsByIDRes() {}

type GetServerStatsByIDOKApplicationJSON []GetServerStatsByIDOKItem

func (*GetServerStatsByIDOKApplicationJSON) getServerStatsByIDRes() {}

type GetServerStatsByIDOKItem struct {
	Timestamp time.Time `json:"timestamp"`
	Online    int32     `json:"online"`
}

// GetTimestamp returns the value of Timestamp.
func (s *GetServerStatsByIDOKItem) GetTimestamp() time.Time {
	return s.Timestamp
}

// GetOnline returns the value of Online.
func (s *GetServerStatsByIDOKItem) GetOnline() int32 {
	return s.Online
}

// SetTimestamp sets the value of Timestamp.
func (s *GetServerStatsByIDOKItem) SetTimestamp(val time.Time) {
	s.Timestamp = val
}

// SetOnline sets the value of Online.
func (s *GetServerStatsByIDOKItem) SetOnline(val int32) {
	s.Online = val
}

type GetServerStatsByIDOrder string

const (
	GetServerStatsByIDOrderAsc  GetServerStatsByIDOrder = "asc"
	GetServerStatsByIDOrderDesc GetServerStatsByIDOrder = "desc"
)

// AllValues returns all GetServerStatsByIDOrder values.
func (GetServerStatsByIDOrder) AllValues() []GetServerStatsByIDOrder {
	return []GetServerStatsByIDOrder{
		GetServerStatsByIDOrderAsc,
		GetServerStatsByIDOrderDesc,
	}
}

// MarshalText implements encoding.TextMarshaler.
func (s GetServerStatsByIDOrder) MarshalText() ([]byte, error) {
	switch s {
	case GetServerStatsByIDOrderAsc:
		return []byte(s), nil
	case GetServerStatsByIDOrderDesc:
		return []byte(s), nil
	default:
		return nil, errors.Errorf("invalid value: %q", s)
	}
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (s *GetServerStatsByIDOrder) UnmarshalText(data []byte) error {
	switch GetServerStatsByIDOrder(data) {
	case GetServerStatsByIDOrderAsc:
		*s = GetServerStatsByIDOrderAsc
		return nil
	case GetServerStatsByIDOrderDesc:
		*s = GetServerStatsByIDOrderDesc
		return nil
	default:
		return errors.Errorf("invalid value: %q", data)
	}
}

// GetServersByMultiplayerNotFound is response for GetServersByMultiplayer operation.
type GetServersByMultiplayerNotFound struct{}

func (*GetServersByMultiplayerNotFound) getServersByMultiplayerRes() {}

type GetServersByMultiplayerOKApplicationJSON []GetServersByMultiplayerOKItem

func (*GetServersByMultiplayerOKApplicationJSON) getServersByMultiplayerRes() {}

type GetServersByMultiplayerOKItem struct {
	Name   string `json:"name"`
	Online int32  `json:"online"`
}

// GetName returns the value of Name.
func (s *GetServersByMultiplayerOKItem) GetName() string {
	return s.Name
}

// GetOnline returns the value of Online.
func (s *GetServersByMultiplayerOKItem) GetOnline() int32 {
	return s.Online
}

// SetName sets the value of Name.
func (s *GetServersByMultiplayerOKItem) SetName(val string) {
	s.Name = val
}

// SetOnline sets the value of Online.
func (s *GetServersByMultiplayerOKItem) SetOnline(val int32) {
	s.Online = val
}

// NewOptDateTime returns new OptDateTime with value set to v.
func NewOptDateTime(v time.Time) OptDateTime {
	return OptDateTime{
		Value: v,
		Set:   true,
	}
}

// OptDateTime is optional time.Time.
type OptDateTime struct {
	Value time.Time
	Set   bool
}

// IsSet returns true if OptDateTime was set.
func (o OptDateTime) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptDateTime) Reset() {
	var v time.Time
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptDateTime) SetTo(v time.Time) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptDateTime) Get() (v time.Time, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptDateTime) Or(d time.Time) time.Time {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// NewOptGetMultiplayersSummaryOrder returns new OptGetMultiplayersSummaryOrder with value set to v.
func NewOptGetMultiplayersSummaryOrder(v GetMultiplayersSummaryOrder) OptGetMultiplayersSummaryOrder {
	return OptGetMultiplayersSummaryOrder{
		Value: v,
		Set:   true,
	}
}

// OptGetMultiplayersSummaryOrder is optional GetMultiplayersSummaryOrder.
type OptGetMultiplayersSummaryOrder struct {
	Value GetMultiplayersSummaryOrder
	Set   bool
}

// IsSet returns true if OptGetMultiplayersSummaryOrder was set.
func (o OptGetMultiplayersSummaryOrder) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptGetMultiplayersSummaryOrder) Reset() {
	var v GetMultiplayersSummaryOrder
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptGetMultiplayersSummaryOrder) SetTo(v GetMultiplayersSummaryOrder) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptGetMultiplayersSummaryOrder) Get() (v GetMultiplayersSummaryOrder, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptGetMultiplayersSummaryOrder) Or(d GetMultiplayersSummaryOrder) GetMultiplayersSummaryOrder {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// NewOptGetServerStatsByIDOrder returns new OptGetServerStatsByIDOrder with value set to v.
func NewOptGetServerStatsByIDOrder(v GetServerStatsByIDOrder) OptGetServerStatsByIDOrder {
	return OptGetServerStatsByIDOrder{
		Value: v,
		Set:   true,
	}
}

// OptGetServerStatsByIDOrder is optional GetServerStatsByIDOrder.
type OptGetServerStatsByIDOrder struct {
	Value GetServerStatsByIDOrder
	Set   bool
}

// IsSet returns true if OptGetServerStatsByIDOrder was set.
func (o OptGetServerStatsByIDOrder) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptGetServerStatsByIDOrder) Reset() {
	var v GetServerStatsByIDOrder
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptGetServerStatsByIDOrder) SetTo(v GetServerStatsByIDOrder) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptGetServerStatsByIDOrder) Get() (v GetServerStatsByIDOrder, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptGetServerStatsByIDOrder) Or(d GetServerStatsByIDOrder) GetServerStatsByIDOrder {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// NewOptInt32 returns new OptInt32 with value set to v.
func NewOptInt32(v int32) OptInt32 {
	return OptInt32{
		Value: v,
		Set:   true,
	}
}

// OptInt32 is optional int32.
type OptInt32 struct {
	Value int32
	Set   bool
}

// IsSet returns true if OptInt32 was set.
func (o OptInt32) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptInt32) Reset() {
	var v int32
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptInt32) SetTo(v int32) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptInt32) Get() (v int32, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptInt32) Or(d int32) int32 {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// NewOptString returns new OptString with value set to v.
func NewOptString(v string) OptString {
	return OptString{
		Value: v,
		Set:   true,
	}
}

// OptString is optional string.
type OptString struct {
	Value string
	Set   bool
}

// IsSet returns true if OptString was set.
func (o OptString) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptString) Reset() {
	var v string
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptString) SetTo(v string) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptString) Get() (v string, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptString) Or(d string) string {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}
