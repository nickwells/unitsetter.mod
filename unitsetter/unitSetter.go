package unitsetter

import (
	"fmt"
	"sort"
	"strings"

	"github.com/nickwells/check.mod/v2/check"
	"github.com/nickwells/param.mod/v7/psetter"
	"github.com/nickwells/strdist.mod/v2/strdist"
	"github.com/nickwells/units.mod/v2/units"
)

// UnitCheckFunc is an alias for the UnitSetter value checks
type UnitCheckFunc = check.ValCk[units.Unit]

// UnitSetter allows you to specify a parameter that can be used to set a
// Unit value. You can also supply check functions that will validate the
// Value.
type UnitSetter struct {
	psetter.ValueReqMandatory
	psetter.ValueChecker[units.Unit]

	// Value must be set, the program will panic if not. This is the value
	// being set
	Value *units.Unit
	// F must be set, the program will panic if not. This is the units.Family
	// to which the units being set must belong.
	F *units.Family
	// ValDesc can be used to describe the value in a help message describing
	// the parameter. If it is not set the unit.Family description will be
	// used (with any spaces replaced by dashes).
	ValDesc string
}

// suggestAltVal will suggest a possible alternative value for the parameter
// value. It will find those strings in the set of possible values that are
// closest to the given value
func (s UnitSetter) suggestAltVal(val string) string {
	names := s.F.GetUnitNames()
	names = append(names, s.F.GetUnitAliases()...)

	return strdist.SuggestionString(strdist.SuggestedVals(val, names))
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

	if err := s.ApplyChecks(v); err != nil {
		return err
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
	if s.Value == nil {
		panic(psetter.NilValueMessage(name, fmt.Sprintf("%T", s)))
	}

	if s.F == nil {
		panic(psetter.BadSetterMessage(name, fmt.Sprintf("%T", s),
			"the Family (F) has not been set"))
	}

	if len(s.F.GetUnitNames()) == 0 {
		panic(psetter.BadSetterMessage(name, fmt.Sprintf("%T", s),
			fmt.Sprintf("the Family (%q) has no units", s.F.Name())))
	}

	s.VerifyChecks(name, fmt.Sprintf("%T", s))
}
