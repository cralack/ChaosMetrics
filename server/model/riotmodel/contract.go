package riotmodel

// DEFINE BY PACK JSON

// Unmarshaler is the interface implemented by types
// that can unmarshal a JSON description of themselves.
// The input can be assumed to be a valid encoding of
// a JSON value. UnmarshalJSON must copy the JSON data
// if it wishes to retain the data after returning.
//
// By convention, to approximate the behavior of Unmarshal itself,
// Unmarshalers implement UnmarshalJSON([]byte("null")) as a no-op.
// type Unmarshaler interface {
// 	UnmarshalJSON([]byte) error
// }

type DTO interface {
	UnmarshalJSON(data []byte) error
}

// TIER LEVEL
const (
	CHALLENGER = iota
	GRANDMASTER
	MASTER
	DIAMOND
	PLATINUM
	GOLD
	SILVER
	BRONZE
	IRON
)

const (
	BR1 = iota
	EUN1
	EUW1
	JP1
	KR1
	LA1
	LA2
	NA1
	OC1
	PH2
	RU
	SG2
	TH2
	TR1
	TW2
	VN2
)
