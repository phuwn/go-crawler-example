package crawler

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
)

type SquadImage struct {
	Original  string `json:"original"`
	Thumbnail string `json:"thumbnail"`
}

func (si SquadImage) Value() (driver.Value, error) {
	return json.Marshal(si)
}

func (si *SquadImage) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("type assertion to []byte failed")
	}

	return json.Unmarshal(b, si)
}

type SquadAttributes map[string]string

type RawSquadAttribute struct {
	Type  string `json:"traitType"`
	Value string `json:"value"`
}

func (sa SquadAttributes) Value() (driver.Value, error) {
	return json.Marshal(sa)
}

func (sa *SquadAttributes) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("type assertion to []byte failed")
	}

	return json.Unmarshal(b, sa)
}

func (sa *SquadAttributes) UnmarshalJSON(data []byte) error {
	var rawData []RawSquadAttribute

	if err := json.Unmarshal(data, &rawData); err != nil {
		return err
	}

	*sa = make(SquadAttributes)
	for _, raw := range rawData {
		(*sa)[raw.Type] = raw.Value
	}

	return nil
}

type SquadID int

func (sid *SquadID) UnmarshalJSON(data []byte) error {
	var rawID string

	if err := json.Unmarshal(data, &rawID); err != nil {
		return err
	}

	id, err := strconv.Atoi(rawID)
	if err != nil {
		return err
	}
	*sid = SquadID(id)

	return nil
}

type PancakeSquad struct {
	ID         SquadID         `json:"tokenId"`
	Image      SquadImage      `json:"image"`
	Attributes SquadAttributes `json:"attributes"`
}

func (PancakeSquad) TableName() string {
	return "pancake_squad"
}
