package database

import (
	"database/sql"
	"encoding/json"
)

// NullString is an alias for sql.NullString data type
type NullString struct {
	sql.NullString
}

// MarshalJSON for NullString
func (ns NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}

// NullInt32 is an alias for sql.NullInt32 data type
type NullInt32 struct {
	sql.NullInt32
}

// MarshalJSON for NullInt32
func (ni NullInt32) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ni.Int32)
}
