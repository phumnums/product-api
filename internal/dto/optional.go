package dto

import "encoding/json"

type OptionalString struct {
	Set   bool
	Value *string
}

type OptionalFloat struct {
	Set   bool
	Value *float64
}

func (o *OptionalString) UnmarshalJSON(data []byte) error {
	o.Set = true

	if string(data) == "null" {
		o.Value = nil
		return nil
	}

	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	o.Value = &value

	return nil
}

func (o *OptionalFloat) UnmarshalJSON(data []byte) error {
	o.Set = true

	if string(data) == "null" {
		o.Value = nil
		return nil
	}

	var value float64
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	o.Value = &value

	return nil
}
