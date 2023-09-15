package unitsetter_test

import (
	"errors"
	"testing"

	"github.com/nickwells/param.mod/v6/param"
	"github.com/nickwells/param.mod/v6/paramtest"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
	"github.com/nickwells/units.mod/v2/units"
	"github.com/nickwells/unitsetter.mod/v4/unitsetter"
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
	var u units.Unit
	testCases := []paramtest.Setter{
		{
			ID: testhelper.MkID("bad-value"),
			PSetter: unitsetter.UnitSetter{
				F: units.SampleFamily,
			},
			ExpPanic: testhelper.MkExpPanic(
				dfltParamName + ": unitsetter.UnitSetter Check failed: " +
					"the Value to be set is nil"),
		},
		{
			ID: testhelper.MkID("bad-setter-no-Family"),
			PSetter: unitsetter.UnitSetter{
				Value: &u,
			},
			ExpPanic: testhelper.MkExpPanic(
				dfltParamName + ": unitsetter.UnitSetter Check failed: " +
					"the Family (F) has not been set"),
		},
		{
			ID: testhelper.MkID("bad-setter-no-units"),
			PSetter: unitsetter.UnitSetter{
				Value: &u,
				F:     units.BadSampleFamily,
			},
			ExpPanic: testhelper.MkExpPanic(
				dfltParamName + ": unitsetter.UnitSetter Check failed: " +
					`the Family ("` +
					units.BadSampleFamily.Name() +
					`") has no units`),
		},
		{
			ID: testhelper.MkID("bad-setter-nil-checkfunc"),
			PSetter: unitsetter.UnitSetter{
				Value:  &u,
				F:      units.SampleFamily,
				Checks: []unitsetter.UnitCheckFunc{nil},
			},
			ExpPanic: testhelper.MkExpPanic(
				dfltParamName + ": unitsetter.UnitSetter Check failed: " +
					"one of the check functions is nil"),
		},
		{
			ID: testhelper.MkID("good-setter-bad-val"),
			PSetter: unitsetter.UnitSetter{
				Value: &u,
				F:     units.SampleFamily,
			},
			ParamVal: "nonesuch",
			SetWithValErr: testhelper.MkExpErr(
				`there is no unit of test called "nonesuch".`),
		},
		{
			ID: testhelper.MkID("good-setter-bad-val-close-match"),
			PSetter: unitsetter.UnitSetter{
				Value: &u,
				F:     units.SampleFamily,
			},
			ParamVal: "sampl",
			SetWithValErr: testhelper.MkExpErr(
				`there is no unit of test called "sampl".`,
				`Did you mean: "sample"`),
		},
		{
			ID: testhelper.MkID("good-setter-good-val"),
			PSetter: unitsetter.UnitSetter{
				Value: &u,
				F:     units.SampleFamily,
			},
			ParamVal:     units.SampleUnitA,
			ValDescriber: true,
		},
		{
			ID: testhelper.MkID("good-setter-good-val-described-val"),
			PSetter: unitsetter.UnitSetter{
				Value:   &u,
				F:       units.SampleFamily,
				ValDesc: "val-name",
			},
			ParamVal:     units.SampleUnitA,
			ValDescriber: true,
		},
		{
			ID: testhelper.MkID("good-setter-good-val-failing-check"),
			PSetter: unitsetter.UnitSetter{
				Value: &u,
				F:     units.SampleFamily,
				Checks: []unitsetter.UnitCheckFunc{
					func(u units.Unit) error {
						return errors.New("always fails")
					},
				},
			},
			SetWithValErr: testhelper.MkExpErr("always fails"),
			ParamVal:      units.SampleUnitA,
		},
	}

	for _, tc := range testCases {
		tc.GFC = commonUnitsetterGFC
		if tc.ParamName == "" {
			tc.ParamName = dfltParamName
		}
		tc.SetVR(param.Mandatory)

		u = units.SampleFamily.GetUnitOrPanic(units.SampleFamily.BaseUnitName())

		tc.Test(t)
	}
}
