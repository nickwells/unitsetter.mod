package unitsetter

import (
	"fmt"
	"sort"
	"strings"

	"github.com/nickwells/param.mod/v3/param"
	"github.com/nickwells/param.mod/v3/param/psetter"
	"github.com/nickwells/strdist.mod/strdist"
	"github.com/nickwells/units.mod/units"
)

// UnitCheckFunc is the type of the check function for this setter. It takes
// a Unit parameter and returns an error (or nil)
type UnitCheckFunc func(units.Unit) error

// UnitSetter allows you to specify a parameter that can be used to set a
// Unit value. You can also supply a check function that will validate the
// Value.
type UnitSetter struct {
	param.ValueReqMandatory
	param.NilAVM

	Value  *units.Unit
	UD     units.UnitDetails
	Checks []UnitCheckFunc
}

// CountChecks returns the number of check functions
func (s UnitSetter) CountChecks() int {
	return len(s.Checks)
}

// suggestAltVal will suggest a possible alternative value for the parameter
// value. It will find those strings in the set of possible values that are
// closest to the given value
func (s UnitSetter) suggestAltVal(val string) string {
	suggestedNames := ""
	matches :=
		strdist.CaseBlindCosineFinder.FindNStrLike(3, val, s.validNames()...)

	if len(matches) > 0 {
		suggestedNames = " Did you mean: " +
			strings.Join(matches, " or ") +
			"?"
	}
	return suggestedNames
}

// SetWithVal (called when a value follows the parameter) checks that the
// value can be found in the map of Units, if it cannot it returns an
// error. If there are checks and any check is violated it returns an
// error. Only if the value is parsed successfully and no checks are violated
// is the Value set.
func (s UnitSetter) SetWithVal(_ string, paramVal string) error {
	v, ok := s.UD.AltU[paramVal]
	if !ok {

		return fmt.Errorf("'%s' is not a recognised %s.%s",
			paramVal, s.UD.Fam.Description, s.suggestAltVal(paramVal))
	}

	if len(s.Checks) != 0 {
		for _, check := range s.Checks {
			if check == nil {
				continue
			}

			err := check(v)
			if err != nil {
				return err
			}
		}
	}

	*s.Value = v
	return nil
}

// validNames returns a slice containing the names of allowed values
func (s UnitSetter) validNames() []string {
	var names []string

	for k := range s.UD.AltU {
		names = append(names, k)
	}

	return names
}

// AllowedValues returns a string describing the allowed values
func (s UnitSetter) AllowedValues() string {
	var names = s.validNames()
	if len(names) == 0 {
		return "there are no valid conversions for this unit type: " +
			s.UD.Fam.BaseUnitName
	}

	sort.Strings(names)
	rval := strings.Join(names, ", ")

	rval += psetter.HasChecks(s)

	return rval
}

// CurrentValue returns the current setting of the parameter value
func (s UnitSetter) CurrentValue() string {
	return s.Value.Name
}

// CheckSetter panics if the setter has not been properly created - if the
// Value is nil, if the base unit is invalid or if one of the check functions
// is nil.
func (s UnitSetter) CheckSetter(name string) {
	if s.Value == nil {
		panic(name +
			": UnitSetter Check failed: the Value to be set is nil")
	}
	if s.UD.AltU == nil {
		panic(name +
			": UnitSetter Check failed: there are no valid alternative units")
	}
	for _, check := range s.Checks {
		if check == nil {
			panic(name +
				": UnitSetter Check failed: one of the check functions is nil")
		}
	}
}
