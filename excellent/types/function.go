package types

import (
	"github.com/nyaruka/goflow/utils"
)

// XFunction is a callable function
type XFunction func(env utils.Environment, args ...XValue) XValue

// Describe returns a representation of this type for error messages
func (x XFunction) Describe() string { return x.String() }

// Reduce returns the primitive version of this type (i.e. itself)
func (x XFunction) Reduce(env utils.Environment) XPrimitive { return x }

// ToXText converts this type to text
func (x XFunction) ToXText(env utils.Environment) XText {
	return NewXText(x.String())
}

// ToXBoolean converts this type to a bool
func (x XFunction) ToXBoolean(env utils.Environment) XBoolean { return XBooleanTrue }

// ToXJSON is called when this type is passed to @(json(...))
func (x XFunction) ToXJSON(env utils.Environment) XText { return MustMarshalToXText(x.String()) }

// String returns the native string representation of this type
func (x XFunction) String() string { return "function" }

var _ XPrimitive = XFunction(nil)