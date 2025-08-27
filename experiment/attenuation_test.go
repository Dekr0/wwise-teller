package experiment_test

import (
	"testing"

	"github.com/Dekr0/wwise-teller/integration/helldivers"
	"github.com/Dekr0/wwise-teller/parser"
	"github.com/Dekr0/wwise-teller/utils"
	"github.com/Dekr0/wwise-teller/wwise"
)

func TestAttenuation(t *testing.T) {
	bnk, _ := parser.ParseBank(
		"./content_audio_stratagems_sentry_machine_gun_base.st_bnk",
		t.Context(),
		false,
	)

	hirc := bnk.HIRC()
	v, _ := hirc.Actions.Load(uint32(997520562))
	refStopFireAction := v.(*wwise.Action)
	newStopFireActionId, _ := utils.ShortID()
	newStopFireAction := refStopFireAction.Clone(newStopFireActionId, 769717269)
	hirc.AppendNewActionToEvent(&newStopFireAction, 3043568243)

	v, _ = hirc.Attenuations.Load(uint32(705537317))
	closeAttenuation := v.(*wwise.Attenuation)
	closeAttenuation.AttenuationConversionTables[0] = wwise.AttenuationConversionTable{
		EnumScaling: 2,
		RTPCGraphPointsX: []float32{0, 16, 100, 320},
		RTPCGraphPointsY: []float32{0, 0, -6.4, -6.4},
		RTPCGraphPointsInterp: []uint32{
			uint32(wwise.InterpCurveTypeConst), 
			uint32(wwise.InterpCurveTypeLog3),
			uint32(wwise.InterpCurveTypeConst),
			uint32(wwise.InterpCurveTypeLinear),
		},
	}

	midAttenuationId, _ := utils.ShortID()
	midAttenuation := wwise.Attenuation{}
	closeAttenuation.Clone(midAttenuationId, &midAttenuation)
	midAttenuation.AttenuationConversionTables[0] = wwise.AttenuationConversionTable{
		EnumScaling: 2,
		RTPCGraphPointsX: []float32{0, 72, 320},
		RTPCGraphPointsY: []float32{-200, -0.4, -0.4},
		RTPCGraphPointsInterp: []uint32{
			uint32(wwise.InterpCurveTypeExp3),
			uint32(wwise.InterpCurveTypeConst),
			uint32(wwise.InterpCurveTypeLinear),
		},
	}
	hirc.AppendNewAttenuation(&midAttenuation)

	v, _ = hirc.ActorMixerHirc.Load(uint32(769717269))
	midRand := v.(*wwise.RanSeqCntr)
	midRand.BaseParam.PropBundle.Remove(wwise.TPriorityDistanceOffset, 154)
	idx, _ := midRand.BaseParam.PropBundle.HasPid(wwise.TAttenuationID, 154)
	midRand.BaseParam.PropBundle.SetPropByIdxU32(idx, midAttenuationId)

	data, err := bnk.Encode(t.Context(), true, false)
	if err != nil {
		t.Fatal(err)
	}
	helldivers.GenHelldiversPatchStable(data, bnk.META().B, ".")
}
