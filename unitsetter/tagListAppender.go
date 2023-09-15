package unitsetter

import (
	"fmt"
	"sort"
	"strings"

	"github.com/nickwells/param.mod/v6/psetter"
	"github.com/nickwells/strdist.mod/strdist"
	"github.com/nickwells/units.mod/v2/units"
)

// TagListAppender is a parameter setter used to add to a slice of units.Tag
// values.
type TagListAppender struct {
	psetter.ValueReqMandatory

	Value *[]units.Tag
}

// suggestAltVal will suggest a possible alternative value for the parameter
// value. It will find those strings in the set of possible values that are
// closest to the given value
func (s TagListAppender) suggestAltVal(val string) string {
	matches := strdist.CaseBlindCosineFinder.FindNStrLike(
		3, val, units.GetTagNames()...)

	return suggestionString(matches)
}

// SetWithVal (called when a value follows the parameter) checks that the
// value can be found in the map of Units, if it cannot it returns an
// error. It then checks that the slice does not already contain the tag and
// returns an error if it does.  Only if all these checks pass is the Value
// set.
func (s TagListAppender) SetWithVal(_ string, paramVal string) error {
	tag := units.Tag(paramVal)
	if !tag.IsValid() {
		return fmt.Errorf("There is no unit tag called %q.%s",
			tag, s.suggestAltVal(paramVal))
	}

	for _, existingTag := range *s.Value {
		if existingTag == tag {
			return fmt.Errorf("Tag  %q is already in the list of tags", tag)
		}
	}

	*s.Value = append(*s.Value, tag)
	return nil
}

// AllowedValues returns a string describing the allowed values
func (s TagListAppender) AllowedValues() string {
	names := units.GetTagNames()
	sort.Strings(names)
	rval := strings.Join(names, ", ")

	return rval
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
	rval := ""
	sep := ""
	for _, tag := range *s.Value {
		rval += sep + string(tag)
		sep = ", "
	}

	return rval
}

// CheckSetter panics if the setter has not been properly created - if the
// Value is nil.
func (s TagListAppender) CheckSetter(name string) {
	intro := name + ": unitsetter.TagListAppender Check failed: "
	if s.Value == nil {
		panic(intro + "the Value to be set is nil")
	}
}
