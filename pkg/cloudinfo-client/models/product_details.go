// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"strconv"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// ProductDetails ProductDetails extended view of the virtual machine details
// swagger:model ProductDetails
type ProductDetails struct {

	// attributes
	Attributes map[string]string `json:"attributes,omitempty"`

	// Burst this is derived for now
	Burst bool `json:"burst,omitempty"`

	// category
	Category string `json:"category,omitempty"`

	// cpus
	Cpus float64 `json:"cpusPerVm,omitempty"`

	// CurrentGen signals whether the instance type generation is the current one. Only applies for amazon
	CurrentGen bool `json:"currentGen,omitempty"`

	// gpus
	Gpus float64 `json:"gpusPerVm,omitempty"`

	// mem
	Mem float64 `json:"memPerVm,omitempty"`

	// ntw perf
	NtwPerf string `json:"ntwPerf,omitempty"`

	// ntw perf cat
	NtwPerfCat string `json:"ntwPerfCategory,omitempty"`

	// on demand price
	OnDemandPrice float64 `json:"onDemandPrice,omitempty"`

	// spot price
	SpotPrice []*ZonePrice `json:"spotPrice"`

	// type
	Type string `json:"type,omitempty"`

	// zones
	Zones []string `json:"zones"`
}

// Validate validates this product details
func (m *ProductDetails) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateSpotPrice(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ProductDetails) validateSpotPrice(formats strfmt.Registry) error {

	if swag.IsZero(m.SpotPrice) { // not required
		return nil
	}

	for i := 0; i < len(m.SpotPrice); i++ {
		if swag.IsZero(m.SpotPrice[i]) { // not required
			continue
		}

		if m.SpotPrice[i] != nil {
			if err := m.SpotPrice[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("spotPrice" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *ProductDetails) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ProductDetails) UnmarshalBinary(b []byte) error {
	var res ProductDetails
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
