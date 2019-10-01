package events

import (
	"github.com/greatnonprofits-nfp/goflow/assets"
	"github.com/greatnonprofits-nfp/goflow/flows"
)

func init() {
	RegisterType(TypeContactFieldChanged, func() flows.Event { return &ContactFieldChangedEvent{} })
}

// TypeContactFieldChanged is the type of our save to contact event
const TypeContactFieldChanged string = "contact_field_changed"

// ContactFieldChangedEvent events are created when a custom field value of the contact has been changed.
// A null values indicates that the field value has been cleared.
//
//   {
//     "type": "contact_field_changed",
//     "created_on": "2006-01-02T15:04:05Z",
//     "field": {"key": "gender", "name": "Gender"},
//     "value": {"text": "Male"}
//   }
//
// @event contact_field_changed
type ContactFieldChangedEvent struct {
	BaseEvent

	Field *assets.FieldReference `json:"field" validate:"required"`
	Value *flows.Value           `json:"value"`
}

// NewContactFieldChangedEvent returns a new save to contact event
func NewContactFieldChangedEvent(field *flows.Field, value *flows.Value) *ContactFieldChangedEvent {
	return &ContactFieldChangedEvent{
		BaseEvent: NewBaseEvent(TypeContactFieldChanged),
		Field:     field.Reference(),
		Value:     value,
	}
}
