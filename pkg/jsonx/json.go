package jsonx

import (
	jsoniter "github.com/json-iterator/go"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
	// Marshal Marshal by jsoniter
	Marshal = json.Marshal
	// Unmarshal Unmarshal by jsoniter
	Unmarshal = json.Unmarshal
	// MarshalIndent MarshalIndent by jsoniter
	MarshalIndent = json.MarshalIndent
	// NewDecoder NewDecoder by jsoniter
	NewDecoder = json.NewDecoder
	// NewEncoder NewEncoder by jsoniter
	NewEncoder = json.NewEncoder
)

// MarshalToString add func of jsoniter.MarshalToString
func MarshalToString(v interface{}) string {
	s, err := jsoniter.MarshalToString(v)
	if err != nil {
		return ""
	}
	return s
}
