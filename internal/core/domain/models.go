package domain

import (
	"encoding/json"
	"strconv"
)

// FlexInt handles both int and string JSON inputs
type FlexInt int

// UnmarshalJSON implements custom unmarshaling
func (fi *FlexInt) UnmarshalJSON(b []byte) error {
	if b[0] == '"' {
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		val, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		*fi = FlexInt(val)
		return nil
	}
	var i int
	if err := json.Unmarshal(b, &i); err != nil {
		return err
	}
	*fi = FlexInt(i)
	return nil
}

// SearchRequest represents the request payload for the search endpoint
type SearchRequest struct {
	MaxLength     FlexInt    `json:"maxLength"`
	MinLength     FlexInt    `json:"minLength"`
	MaxWords      FlexInt    `json:"maxWords"`
	LettersMatrix [][]string `json:"lettersMatrix"`
}

// UpdateWordsRequest represents the request payload for updating words
type UpdateWordsRequest struct {
	Words   []string `json:"words"`
	Include bool     `json:"include"`
}

// MergeResponse represents the response for the merge operation
type MergeResponse struct {
	AddedCount   int `json:"addedCount"`
	RemovedCount int `json:"removedCount"`
}

// LookupResultResponseItem represents an item in the lookup result
type LookupResultResponseItem struct {
	Word      string `json:"word"`
	Found     bool   `json:"found"`
	Source    string `json:"source"`
	Line      int    `json:"line"`
	Timestamp string `json:"timestamp"`
}
