package unitsetter_test

import (
	"errors"
	"testing"

	"github.com/nickwells/param.mod/v5/param"
	"github.com/nickwells/param.mod/v5/paramtest"
	"github.com/nickwells/testhelper.mod/testhelper"
	"github.com/nickwells/units.mod/units"
	"github.com/nickwells/unitsetter.mod/v3/unitsetter"
)

const (
	updFlagNameUnitsetter     = "upd-gf-unitsetter"
	keepBadFlagNameUnitsetter = "keep-bad-unitsetter"
)

var commonUnitsetterGFC = testhelper.GoldenFileCfg{
	DirNames:               []string{"testdata", "UnitSetter"},
	Pfx:                    "gf",
	Sfx:                    "txt",
	UpdFlagName:            updFlagNameUnitsetter,
	KeepBadResultsFlagName: keepBadFlagNameUnitsetter,
}

func init() {
	commonUnitsetterGFC.AddUpdateFlag()
	commonUnitsetterGFC.AddKeepBadResultsFlag()
}

func TestUnitSetter(t *testing.T) {
	const dfltParamName = "set-unit"
	goodFam := units.Family{
		BaseUnitName: "test",
		Description:  "unit of test",
	}
	dfltUnit := units.Unit{
		ConvPreAdd:  1,
		ConvPostAdd: 2,
		ConvFactor:  3,
		Fam:         goodFam,
		Abbrev:      "dflt",
		Name:        "dflt",
		NamePlural:  "dflts",
		Notes:       "desc of default Unit value",
	}
	var u units.Unit
	goodUnitDetails := units.UnitDetails{
		Fam: goodFam,
		AltU: map[string]units.Unit{
			"test": {
				ConvPreAdd:  0,
				ConvPostAdd: 0,
				ConvFactor:  1,
				Fam:         goodFam,
				Abbrev:      "t",
				Name:        "test",
				NamePlural:  "tests",
				Notes:       "desc",
			},
			"other-unit": {
				ConvPreAdd:  0,
				ConvPostAdd: 0,
				ConvFactor:  2,
				Fam:         goodFam,
				Abbrev:      "o-u-t",
				Name:        "o-u-test",
				NamePlural:  "o-u-tests",
				Notes:       "desc of other-unit",
			},
			"otherunit": {
				ConvPreAdd:  0,
				ConvPostAdd: 0,
				ConvFactor:  3,
				Fam:         goodFam,
				Abbrev:      "ou-t",
				Name:        "ou-test",
				NamePlural:  "ou-tests",
				Notes:       "desc of otherunit",
			},
		},
		Aliases: map[string]units.Alias{
			"other-units": {
				UnitName: "other-unit",
				Notes:    "plural",
			},
		},
	}
	testCases := []paramtest.Setter{
		{
			ID: testhelper.MkID("bad-value"),
			PSetter: unitsetter.UnitSetter{
				UD: goodUnitDetails,
			},
			ExpPanic: testhelper.MkExpPanic(
				dfltParamName + ": unitsetter.UnitSetter Check failed: " +
					"the Value to be set is nil"),
		},
		{
			ID: testhelper.MkID("bad-unit-details-nil-AltU"),
			PSetter: unitsetter.UnitSetter{
				Value: &u,
			},
			ExpPanic: testhelper.MkExpPanic(
				dfltParamName + ": unitsetter.UnitSetter Check failed: " +
					"there are no valid alternative units"),
		},
		{
			ID: testhelper.MkID("bad-unit-details-empty-AltU"),
			PSetter: unitsetter.UnitSetter{
				Value: &u,
				UD: units.UnitDetails{
					AltU: map[string]units.Unit{},
				},
			},
			ExpPanic: testhelper.MkExpPanic(
				dfltParamName + ": unitsetter.UnitSetter Check failed: " +
					"there are no valid alternative units"),
		},
		{
			ID: testhelper.MkID("bad-unit-details-nil-checkfunc"),
			PSetter: unitsetter.UnitSetter{
				Value:  &u,
				UD:     goodUnitDetails,
				Checks: []unitsetter.UnitCheckFunc{nil},
			},
			ExpPanic: testhelper.MkExpPanic(
				dfltParamName + ": unitsetter.UnitSetter Check failed: " +
					"one of the check functions is nil"),
		},
		{
			ID: testhelper.MkID("good-unit-details-bad-val"),
			PSetter: unitsetter.UnitSetter{
				Value: &u,
				UD:    goodUnitDetails,
			},
			ParamVal: "nonesuch",
			SetWithValErr: testhelper.MkExpErr(
				`there is no unit of test called "nonesuch".`),
		},
		{
			ID: testhelper.MkID("good-unit-details-bad-val-close-match"),
			PSetter: unitsetter.UnitSetter{
				Value: &u,
				UD:    goodUnitDetails,
			},
			ParamVal: "other-unnit",
			SetWithValErr: testhelper.MkExpErr(
				`there is no unit of test called "other-unnit".`,
				"Did you mean: other-unit or other-units or otherunit"),
		},
		{
			ID: testhelper.MkID("good-unit-details-good-val"),
			PSetter: unitsetter.UnitSetter{
				Value: &u,
				UD:    goodUnitDetails,
			},
			ParamVal:     "other-unit",
			ValDescriber: true,
		},
		{
			ID: testhelper.MkID("good-unit-details-good-val-described-val"),
			PSetter: unitsetter.UnitSetter{
				Value:   &u,
				UD:      goodUnitDetails,
				ValDesc: "val-name",
			},
			ParamVal:     "other-unit",
			ValDescriber: true,
		},
		{
			ID: testhelper.MkID("good-unit-details-good-val-failing-check"),
			PSetter: unitsetter.UnitSetter{
				Value: &u,
				UD:    goodUnitDetails,
				Checks: []unitsetter.UnitCheckFunc{
					func(u units.Unit) error {
						return errors.New("always fails")
					},
				},
			},
			SetWithValErr: testhelper.MkExpErr("always fails"),
			ParamVal:      "other-unit",
		},
	}

	for _, tc := range testCases {
		tc.GFC = commonUnitsetterGFC
		if tc.ParamName == "" {
			tc.ParamName = dfltParamName
		}
		tc.SetVR(param.Mandatory)

		u = dfltUnit

		tc.Test(t)
	}
}
