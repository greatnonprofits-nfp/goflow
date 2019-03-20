package routers

import (
	"encoding/json"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"

	"github.com/pkg/errors"
)

type category struct {
	uuid     flows.CategoryUUID
	name     string
	exitUUID flows.ExitUUID
}

// NewCategory creates a new category
func NewCategory(uuid flows.CategoryUUID, name string, exit flows.ExitUUID) flows.Category {
	return &category{uuid: uuid, name: name, exitUUID: exit}
}

func (c *category) UUID() flows.CategoryUUID { return c.uuid }
func (c *category) Name() string             { return c.name }
func (c *category) ExitUUID() flows.ExitUUID { return c.exitUUID }

// LocalizationUUID gets the UUID which identifies this object for localization
func (c *category) LocalizationUUID() utils.UUID { return utils.UUID(c.uuid) }

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type categoryEnvelope struct {
	UUID     flows.CategoryUUID `json:"uuid"                validate:"required,uuid4"`
	Name     string             `json:"name,omitempty"`
	ExitUUID flows.ExitUUID     `json:"exit_uuid,omitempty" validate:"omitempty,uuid4"`
}

// UnmarshalJSON unmarshals a node category from the given JSON
func (c *category) UnmarshalJSON(data []byte) error {
	e := &categoryEnvelope{}

	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return errors.Wrap(err, "unable to read exit")
	}

	c.uuid = e.UUID
	c.name = e.Name
	c.exitUUID = e.ExitUUID
	return nil
}

// MarshalJSON marshals this node category into JSON
func (c *category) MarshalJSON() ([]byte, error) {
	return json.Marshal(&categoryEnvelope{
		c.uuid,
		c.name,
		c.exitUUID,
	})
}
