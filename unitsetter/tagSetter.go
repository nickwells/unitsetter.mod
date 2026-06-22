package unitsetter

import (
	"fmt"
	"sort"

	"github.com/nickwells/param.mod/v7/psetter"
	"github.com/nickwells/param.mod/v7/ptypes"
	"github.com/nickwells/strdist.mod/v2/strdist"
	"github.com/nickwells/units.mod/v2/units"
)

// TagSetter is a parameter setter used to populate a units.Tag value.
type TagSetter struct {
	psetter.ValueReqMandatory

	// Value must be set, the program will panic if not. This is the Tag
	// value that this setter is setting.
	Value *units.Tag
}

// SetWithVal (called when a value follows the parameter) checks that the
// value can be found in the map of Units, if it cannot it returns an
// error. If there are checks and any check is violated it returns an
// error. Only if the value is parsed successfully and no checks are violated
// is the Value set.
func (s TagSetter) SetWithVal(_ string, paramVal string) error {
	tag := units.Tag(paramVal)
	if !tag.IsValid() {
		return fmt.Errorf("there is no unit tag called %q%s",
			tag, strdist.SuggestionString(
				strdist.SuggestedVals(
					paramVal,
					units.GetTagNames(),
				)))
	}

	*s.Value = tag

	return nil
}

// AllowedValues returns a string describing the allowed values
func (s TagSetter) AllowedValues() string {
	return "a tag name"
}

// AllowedValuesMap returns a map of allowed tags. This will be used by the
// standard help package to generate a list of allowed values.
func (s TagSetter) AllowedValuesMap() ptypes.AllowedVals[string] {
	names := units.GetTagNames()
	sort.Strings(names)

	avm := ptypes.AllowedVals[string]{}

	for _, t := range names {
		avm[t] = units.Tag(t).Notes()
	}

	return avm
}

// ValDescribe returns a string describing the value that can follow the
// parameter
func (s TagSetter) ValDescribe() string {
	return "unit-tag"
}

// CurrentValue returns the current setting of the parameter value
func (s TagSetter) CurrentValue() string {
	if s.Value == nil {
		return ""
	}

	return string(*s.Value)
}

// CheckSetter panics if the setter has not been properly created - if the
// Value is nil.
func (s TagSetter) CheckSetter(name string) {
	intro := name + ": unitsetter.TagSetter Check failed: "

	if s.Value == nil {
		panic(intro + "the Value to be set is nil")
	}
}
