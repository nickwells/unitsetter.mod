package unitsetter

import (
	"fmt"
	"slices"
	"sort"
	"strings"

	"github.com/nickwells/param.mod/v7/psetter"
	"github.com/nickwells/param.mod/v7/ptypes"
	"github.com/nickwells/strdist.mod/v2/strdist"
	"github.com/nickwells/units.mod/v2/units"
)

// TagListAppender is a parameter setter used to add to a slice of units.Tag
// values.
type TagListAppender struct {
	psetter.ValueReqMandatory

	// The StrListSeparator allows you to override the default separator
	// between list elements.
	psetter.StrListSeparator

	// Value must be set, the program will panic if not. This is the slice of
	// Tag values that this setter is setting.
	Value *[]units.Tag
}

// SetWithVal (called when a value follows the parameter) checks that the
// value can be found in the map of Units, if it cannot it returns an
// error. It then checks that the slice does not already contain the tag and
// returns an error if it does.  Only if all these checks pass is the Value
// set.
func (s TagListAppender) SetWithVal(_ string, paramVal string) error {
	sep := s.GetSeparator()

	vals := strings.SplitSeq(paramVal, sep)
	for v := range vals {
		tag := units.Tag(v)
		if !tag.IsValid() {
			return fmt.Errorf("there is no unit tag called %q%s",
				tag, strdist.SuggestionString(
					strdist.SuggestedVals(
						v,
						units.GetTagNames(),
					)))
		}

		if slices.Contains(*s.Value, tag) {
			return fmt.Errorf("tag  %q is already in the list of tags", tag)
		}

		*s.Value = append(*s.Value, tag)
	}

	return nil
}

// AllowedValues returns a string describing the allowed values
func (s TagListAppender) AllowedValues() string {
	return s.ListValDesc("tag names")
}

// AllowedValuesMap returns a map of allowed tags. This will be used by the
// standard help package to generate a list of allowed values.
func (s TagListAppender) AllowedValuesMap() ptypes.AllowedVals[string] {
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
func (s TagListAppender) ValDescribe() string {
	return "unit-tag"
}

// CurrentValue returns the current setting of the parameter value
func (s TagListAppender) CurrentValue() string {
	if s.Value == nil {
		return ""
	}

	var cv strings.Builder

	sep := ""

	for _, tag := range *s.Value {
		cv.WriteString(sep)
		cv.WriteString(string(tag))

		sep = ", "
	}

	return cv.String()
}

// CheckSetter panics if the setter has not been properly created - if the
// Value is nil.
func (s TagListAppender) CheckSetter(name string) {
	intro := name + ": unitsetter.TagListAppender Check failed: "

	if s.Value == nil {
		panic(intro + "the Value to be set is nil")
	}
}
