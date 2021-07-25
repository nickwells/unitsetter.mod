package unitsetter

import (
	"fmt"
	"sort"
	"strings"

	"github.com/nickwells/param.mod/v5/param/psetter"
	"github.com/nickwells/strdist.mod/strdist"
	"github.com/nickwells/units.mod/v2/units"
)

// TagSetter is a parameter setter used to populate units.Tag values.
type TagSetter struct {
	psetter.ValueReqMandatory

	Value *units.Tag
}

// suggestAltVal will suggest a possible alternative value for the parameter
// value. It will find those strings in the set of possible values that are
// closest to the given value
func (s TagSetter) suggestAltVal(val string) string {
	suggestedNames := ""
	matches :=
		strdist.CaseBlindCosineFinder.FindNStrLike(
			3, val, units.GetTagNames()...)

	if len(matches) > 0 {
		sort.Strings(matches)
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
func (s TagSetter) SetWithVal(_ string, paramVal string) error {
	tag := units.Tag(paramVal)
	if !tag.IsValid() {
		return fmt.Errorf("There is no unit tag called %q.%s",
			tag, s.suggestAltVal(paramVal))
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
