package unitsetter

import (
	"fmt"
	"sort"
	"strings"

	"github.com/nickwells/param.mod/v6/psetter"
	"github.com/nickwells/units.mod/v2/units"
)

// TagSetter is a parameter setter used to populate units.Tag values.
type TagSetter struct {
	psetter.ValueReqMandatory

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
			tag, psetter.SuggestionString(
				psetter.SuggestedVals(
					paramVal,
					units.GetTagNames(),
				)))
	}

	*s.Value = tag

	return nil
}

// AllowedValues returns a string describing the allowed values
func (s TagSetter) AllowedValues() string {
	names := units.GetTagNames()
	sort.Strings(names)
	rval := strings.Join(names, ", ")

	return rval
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
