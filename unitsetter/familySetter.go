package unitsetter

import (
	"fmt"
	"sort"
	"strings"

	"github.com/nickwells/param.mod/v6/psetter"
	"github.com/nickwells/units.mod/v2/units"
)

// FamilySetter is a parameter setter used to populate units.Family values.
type FamilySetter struct {
	psetter.ValueReqMandatory

	Value **units.Family
}

// SetWithVal (called when a value follows the parameter) checks that the
// value can be found in the map of Units, if it cannot it returns an
// error. If there are checks and any check is violated it returns an
// error. Only if the value is parsed successfully and no checks are violated
// is the Value set.
func (s FamilySetter) SetWithVal(_ string, paramVal string) error {
	v, err := units.GetFamily(paramVal)
	if err != nil {
		return fmt.Errorf("%v%s",
			err,
			psetter.SuggestionString(
				psetter.SuggestedVals(
					paramVal,
					units.GetFamilyNames(),
				)))
	}

	*s.Value = v

	return nil
}

// AllowedValues returns a string describing the allowed values
func (s FamilySetter) AllowedValues() string {
	names := units.GetFamilyNames()
	names = append(names, units.GetFamilyAliases()...)

	sort.Strings(names)
	rval := strings.Join(names, ", ")

	return rval
}

// ValDescribe returns a string describing the value that can follow the
// parameter
func (s FamilySetter) ValDescribe() string {
	return "unit-family"
}

// CurrentValue returns the current setting of the parameter value
func (s FamilySetter) CurrentValue() string {
	if *s.Value == nil {
		return ""
	}

	return (*s.Value).Name()
}

// CheckSetter panics if the setter has not been properly created - if the
// Value is nil.
func (s FamilySetter) CheckSetter(name string) {
	intro := name + ": unitsetter.FamilySetter Check failed: "
	if s.Value == nil {
		panic(intro + "the Value to be set is nil")
	}
}
