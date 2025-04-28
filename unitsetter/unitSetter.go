package unitsetter

import (
	"fmt"
	"sort"
	"strings"

	"github.com/nickwells/param.mod/v6/psetter"
	"github.com/nickwells/units.mod/v2/units"
)

// UnitCheckFunc is the type of the check function for this setter. It takes
// a Unit parameter and returns an error (or nil)
type UnitCheckFunc func(units.Unit) error

// UnitSetter allows you to specify a parameter that can be used to set a
// Unit value. You can also supply check functions that will validate the
// Value.
//
// If you give a ValDesc then that is used as the value description in the
// help message, otherwise the Unit Family description is used (with spaces
// replaced by dashes)
type UnitSetter struct {
	psetter.ValueReqMandatory

	Value   *units.Unit
	F       *units.Family
	Checks  []UnitCheckFunc
	ValDesc string
}

// CountChecks returns the number of check functions
func (s UnitSetter) CountChecks() int {
	return len(s.Checks)
}

// suggestAltVal will suggest a possible alternative value for the parameter
// value. It will find those strings in the set of possible values that are
// closest to the given value
func (s UnitSetter) suggestAltVal(val string) string {
	names := s.F.GetUnitNames()
	names = append(names, s.F.GetUnitAliases()...)

	return psetter.SuggestionString(psetter.SuggestedVals(val, names))
}

// SetWithVal (called when a value follows the parameter) checks that the
// value can be found in the map of Units, if it cannot it returns an
// error. If there are checks and any check is violated it returns an
// error. Only if the value is parsed successfully and no checks are violated
// is the Value set.
func (s UnitSetter) SetWithVal(_ string, paramVal string) error {
	v, err := s.F.GetUnit(paramVal)
	if err != nil {
		return fmt.Errorf("%v%s", err, s.suggestAltVal(paramVal))
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

// AllowedValues returns a string describing the allowed values
func (s UnitSetter) AllowedValues() string {
	names := s.F.GetUnitNames()
	if len(names) == 0 {
		return "there are no units in this unit family: " +
			s.F.Description()
	}

	names = append(names, s.F.GetUnitAliases()...)
	sort.Slice(names, func(i, j int) bool {
		// sort the family base name to the front
		if names[i] == s.F.BaseUnitName() {
			return true
		}

		if names[j] == s.F.BaseUnitName() {
			return false
		}

		// then prefer shorter names
		if len(names[i]) != len(names[j]) {
			return len(names[i]) < len(names[j])
		}

		// then alphabetically
		return names[i] < names[j]
	})

	rval := strings.Join(names, ", ")

	rval += psetter.HasChecks(s)

	return rval
}

// ValDescribe returns a string describing the value that can follow the
// parameter
func (s UnitSetter) ValDescribe() string {
	if s.ValDesc != "" {
		return s.ValDesc
	}

	return strings.ReplaceAll(s.F.Description(), " ", "-")
}

// CurrentValue returns the current setting of the parameter value
func (s UnitSetter) CurrentValue() string {
	return s.Value.Name()
}

// CheckSetter panics if the setter has not been properly created - if the
// Value is nil, if the base unit is invalid or if one of the check functions
// is nil.
func (s UnitSetter) CheckSetter(name string) {
	intro := name + ": unitsetter.UnitSetter Check failed:"

	if s.Value == nil {
		panic(intro + " the Value to be set is nil")
	}

	if s.F == nil {
		panic(intro + " the Family (F) has not been set")
	}

	if len(s.F.GetUnitNames()) == 0 {
		panic(fmt.Sprintf("%s the Family (%q) has no units", intro, s.F.Name()))
	}

	for _, check := range s.Checks {
		if check == nil {
			panic(intro + " one of the check functions is nil")
		}
	}
}
