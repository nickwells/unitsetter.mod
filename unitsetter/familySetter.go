package unitsetter

import (
	"fmt"

	"github.com/nickwells/param.mod/v7/psetter"
	"github.com/nickwells/param.mod/v7/ptypes"
	"github.com/nickwells/strdist.mod/v2/strdist"
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
			strdist.SuggestionString(
				strdist.SuggestedVals(
					paramVal,
					units.GetFamilyNames(),
				)))
	}

	*s.Value = v

	return nil
}

// AllowedValues returns a string describing the allowed values
func (s FamilySetter) AllowedValues() string {
	return "a family name"
}

// AllowedValuesMap returns a map of allowed family names. This will be used
// by the standard help package to generate a list of allowed values.
func (s FamilySetter) AllowedValuesMap() ptypes.AllowedVals[string] {
	families := units.GetFamilies()

	avm := ptypes.AllowedVals[string]{}

	for _, f := range families {
		avm[f.Name()] = f.Description()
	}

	return avm
}

// AllowedValuesAliasMap returns a map of allowed family name aliasess. This
// will be used by the standard help package to generate a list of aliases
// for the allowed values.
func (s FamilySetter) AllowedValuesAliasMap() ptypes.Aliases[string] {
	families := units.GetFamilies()

	avam := ptypes.Aliases[string]{}

	for _, f := range families {
		for _, fa := range f.FamilyAliases() {
			avam[fa] = append(avam[fa], f.Name())
		}
	}

	return avam
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
