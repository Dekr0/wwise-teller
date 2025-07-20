// TODO: batch create Random Sequence Container
package automation

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/Dekr0/wwise-teller/db"
	"github.com/Dekr0/wwise-teller/waapi"
	"github.com/Dekr0/wwise-teller/wio"
	"github.com/Dekr0/wwise-teller/wwise"
)

type ImportAsRanSeqCntrScript struct {
	Workspace      string                     `json:"workspace"`
	Conversion     string                     `json:"conversion"`
	Inputs       []string                     `json:"inputs"`
	Seq            bool                       `json:"seq"`
	Format         waapi.ConversionFormatType `json:"format"`
	PlaybackLimit  uint16                     `json:"playbackLimit"`
	MakeUpGain     float32                    `json:"makeUpGain"`
	InitialDelay   float32                    `json:"initalDelay"`
	HDRActiveRange float32                    `json:"HDRActiveRange"` // Set this to any negative value if it's unused
	Parent         uint32                     `json:"parent"`
	// Set either Event or RefContainer to zero for creating a container without 
	// an action, used for creating a new random / sequence container under a 
	// switch container or a layer container
	Event          uint32                     `json:"Event"` 
	RefContainer   uint32                     `json:"refContainer"`
	RefAction      uint32                     `json:"refAction"`
}

// Handle validation 
func ParseImportAsRanSeqCntrScript(s *ImportAsRanSeqCntrScript, script string) (map[string]uint8, error) {
	data, err := os.ReadFile(script)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(data, s); err != nil {
		return nil, err
	}
	
	{
		if !filepath.IsAbs(s.Workspace) {
			return nil, fmt.Errorf("Workspace path %s is not an absolute path.", s.Workspace)
		}
		stat, err := os.Lstat(s.Workspace)
		if err != nil {
			if os.IsNotExist(err) {
				slog.Error(fmt.Sprintf("Workspace %s does not exist", s.Workspace))
			} else {
				slog.Error(fmt.Sprintf("Failed to obtain information of workspace %s", s.Workspace))
			}
			workspace := filepath.Dir(script)
			slog.Warn(fmt.Sprintf("Using default workspace %s", workspace))
			s.Workspace = workspace
		}
		if !stat.IsDir() {
			slog.Error(fmt.Sprintf("Workspace path %s is not a directory", s.Workspace))
			workspace := filepath.Dir(script)
			slog.Warn(fmt.Sprintf("Using default workspace %s", workspace))
			s.Workspace = workspace
		}
	}

	if s.Format < waapi.ConversionFormatTypePCM || s.Format > waapi.ConversionFormatTypeWEMOpus {
		return nil, fmt.Errorf("Invalid conversion format type %d", s.Format)
	}

	{
		if s.PlaybackLimit > 1000 {
			return nil, fmt.Errorf("Invalid playback limit value %d is not in between 0 and 1000", s.PlaybackLimit)
		}
		c, in := wwise.BasePropChecker[wwise.TMakeUpGain]
		if !in {
			panic("Make up gain property checker does not exist.")
		}
		if err := c(s.MakeUpGain); err != nil {
			return nil, err
		}
		c, in = wwise.BasePropChecker[wwise.TInitialDelay]
		if !in {
			panic("Initial delay property checker does not exist.")
		}
		if err := c(s.InitialDelay); err != nil {
			return nil, err
		}
		if s.HDRActiveRange > 24.0 {
			return nil, fmt.Errorf("HDR active range value %f is not in between 0.0 and 24.0", s.HDRActiveRange)
		}
	}

	inputMap := make(map[string]uint8, len(s.Inputs))
	{
		i := uint8(0)
		for _, input := range s.Inputs {
			if !filepath.IsAbs(input) {
				input = filepath.Join(s.Workspace, input)
			}
			ext := filepath.Ext(input)
			if ext == "" {
				input += ".wav"
			} else if !strings.EqualFold(ext, ".wav") {
				slog.Error(fmt.Sprintf("%s is not in wav format.", input))
				continue
			}
			stat, err := os.Lstat(input)
			if err != nil {
				if os.IsNotExist(err) {
					slog.Error(fmt.Sprintf("%s does not exist.", input))
				} else {
					slog.Error(fmt.Sprintf("Failed to obtain information of %s", input), "error", err)
				}
				continue
			}
			if stat.IsDir() {
				slog.Error(fmt.Sprintf("%s is a directory.", input))
				continue
			}
			if _, in := inputMap[input]; !in {
				inputMap[input] = i 
				i += 1
			} else {
				slog.Warn(fmt.Sprintf("%s is duplicated.", input))
			}
		}
	}
	return inputMap, nil
}

func ImportAsRanSeqCntr(ctx context.Context, bnk *wwise.Bank, script string) error {
	if bnk.DIDX() == nil {
		return wwise.NoDIDX
	}
	if bnk.DATA() == nil {
		return wwise.NoDATA
	}
	h := bnk.HIRC()
	if h == nil {
		return wwise.NoHIRC
	}
	proj, err := waapi.GetProject()
	if err != nil {
		return err
	}

	s := ImportAsRanSeqCntrScript{}
	inputsMap, err := ParseImportAsRanSeqCntrScript(&s, script)
	if err != nil {
		return err
	}
	if len(inputsMap) <= 0 {
		slog.Warn("No input file is provided.")
		return nil
	}
	
	v, in := h.ActorMixerHirc.Load(s.Parent)
	if !in {
		return fmt.Errorf("No actor mixer hierarchy object has ID %d.", s.Parent)
	}
	parent := v.(wwise.HircObj)
	switch parent.(type) {
	case *wwise.ActorMixer:
	default:
		return fmt.Errorf("Parent actor mixer hierarchy type %s is yet supported.", wwise.HircTypeName[parent.HircType()])
	}

	_, in = h.Events.Load(s.Event)
	if !in {
		return fmt.Errorf("No event object has ID %d.", s.Event)
	}

	v, in = h.ActorMixerHirc.Load(s.RefContainer)
	if !in {
		return fmt.Errorf("No reference random / sequence container has ID %d.", s.RefContainer)
	}
	switch v.(wwise.HircObj).(type) {
	case *wwise.RanSeqCntr:
	default:
		return fmt.Errorf("Actor mixer hierarchy object %d is not type of random / sequence container.", s.RefContainer)
	}
	refCntr := v.(*wwise.RanSeqCntr)

	// Fetch one reference sound
	var refSound *wwise.Sound
	for i := 0; i < len(refCntr.Container.Children) && refSound == nil ; i++ {
		v, in := h.ActorMixerHirc.Load(refCntr.Container.Children[i])
		if !in { continue }
		o := v.(wwise.HircObj)
		switch h := o.(type) {
		case *wwise.Sound:
			refSound = h
		default:
		}
	}
	if refSound == nil {
		return fmt.Errorf("There's no reference sound in random / sequence container %d to create new sound objects", s.RefContainer)
	}
	v, in = h.Actions.Load(s.RefAction)
	if !in {
		return fmt.Errorf("No action has ID %d.", s.RefAction)
	}
	refAction := v.(*wwise.Action)

	wems := make([]string, len(inputsMap))
	wsource, err := waapi.CreateConversionListInOrder(ctx, inputsMap, wems, s.Conversion, false)
	if len(inputsMap) != len(wems) {
		panic("Length of input maps does not equal to the length of output array")
	}
	if err != nil {
		return err
	}

	if err := waapi.WwiseConversion(ctx, wsource, proj); err != nil {
		return err
	}

	newAudioDatas := make([][]byte, 0, len(wems))
	for _, wem := range wems {
		audioData, err := os.ReadFile(wem)
		if err != nil {
			slog.Error(fmt.Sprintf("Failed to read audio data from %s", wem))
			continue
		}
		newAudioDatas = append(newAudioDatas, audioData)
	}

	// Hierarchy IDs generation and Source IDs generation
	q, closeConn, commit, rollback, err := db.CreateConnWithTxQuery(ctx)
	if err != nil {
		return err
	}
	defer closeConn()
	newSoundIDs := make([]uint32, len(newAudioDatas))
	newSourceIDs := make([]uint32, len(newAudioDatas))
	for i := range newAudioDatas {
		newSoundIDs[i], err = db.TryHid(ctx, q)
		if err != nil {
			rollback()
			return err
		}
		newSourceIDs[i], err = db.TrySid(ctx, q)
		if err != nil {
			rollback()
			return err
		}
	}
	newCntrId, err := db.TryHid(ctx, q)
	if err != nil {
		rollback()
		return err
	}

	var newActionId = uint32(0)
	if s.RefAction != 0 && s.Event != 0 {
		newActionId, err = db.TryHid(ctx, q)
		if err != nil {
			rollback()
			return err
		}
	}

	for i, audioData := range newAudioDatas {
		if err := bnk.AppendAudio(audioData, newSourceIDs[i]); err != nil {
			rollback()
			return fmt.Errorf("Failed to add a new audio source file: %w.", err)
		}
	}
	bnk.ComputeDIDXOffset()
	if err := bnk.CheckDIDXDATA(); err != nil {
		rollback()
		return fmt.Errorf("Invalid Integrity appear in DIDX and DATA chunk: %w", err)
	}

	newCntr := refCntr.Clone(newCntrId, false)
	ImportRanSeqCntrModifer(&newCntr, &s, int(bnk.BKHD().BankGenerationVersion))
	if err := h.AppendNewRanSeqCntrToActorMixer(&newCntr, s.Parent, false); err != nil {
		rollback()
		return fmt.Errorf("Failed to add a new random / sequence container to actor mixer %d: %w", s.Parent, err)
	}

	var newSound *wwise.Sound
	var pluginID uint32
	for i := range newAudioDatas {
		switch s.Format {
		case waapi.ConversionFormatTypePCM:
			pluginID = wwise.PCM
		case waapi.ConversionFormatTypeADPCM:
			pluginID = wwise.ADPCM
		case waapi.ConversionFormatTypeVORBIS:
			pluginID = wwise.VORBIS
		case waapi.ConversionFormatTypeWEMOpus:
			pluginID = wwise.WEM_OPUS
		}
		newSound = &wwise.Sound{
			Id: newSoundIDs[i],
			BankSourceData: wwise.BankSourceData{
				PluginID: pluginID,
				StreamType: wwise.STREAM_TYPE_BNK,
				SourceID: newSourceIDs[i],
				InMemoryMediaSize: uint32(len(newAudioDatas[i])),
				SourceBits: 0,
			},
			BaseParam: refSound.BaseParam.Clone(false),
		}
		if err := h.AppendNewSoundToRanSeqContainer(newSound, newCntrId, false); err != nil {
			rollback()
			return fmt.Errorf("Failed to add a new sound object to random / sequence container %d: %w", newCntrId, err)
		}
	}
	
	if newActionId != 0 {
		newAction := refAction.Clone(newActionId, newCntrId)
		if err := h.AppendNewActionToEvent(&newAction, s.Event); err != nil {
			rollback()
			return fmt.Errorf("Failed to add a new action ot event %d: %w", s.Event, err)
		}
	}

	if err := commit(); err != nil {
		rollback()
		return err
	}

	return nil
}

func NewSoundToRanSeqCntr(ctx context.Context, bnk *wwise.Bank, script string) error {
	if bnk.DIDX() == nil {
		return wwise.NoDIDX
	}
	if bnk.DATA() == nil {
		return wwise.NoDATA
	}
	h := bnk.HIRC()
	if h == nil {
		return wwise.NoHIRC
	}
	proj, err := waapi.GetProject()
	if err != nil {
		return err
	}

	s := ImportAsRanSeqCntrScript{}
	inputsMap, err := ParseImportAsRanSeqCntrScript(&s, script)
	if err != nil {
		return err
	}
	if len(inputsMap) <= 0 {
		slog.Warn("No input file is provided.")
		return nil
	}

	v, in := h.ActorMixerHirc.Load(s.Parent)
	if !in {
		return fmt.Errorf("No actor mixer hierarchy object has ID %d.", s.Parent)
	}

	parent := v.(wwise.HircObj)
	switch parent.(type) {
	case *wwise.RanSeqCntr:
	default:
		return fmt.Errorf("Parent actor mixer hierarchy type %s does not support for this type of script", wwise.HircTypeName[parent.HircType()])
	}

	// Fetch one reference sound
	var refSound *wwise.Sound
	for lid := range parent.Leafs() {
		v, in := h.ActorMixerHirc.Load(lid)
		if !in {
			continue
		}
		o := v.(wwise.HircObj)
		switch h := o.(type) {
		case *wwise.Sound:
			refSound = h
		default:
		}
		if refSound != nil {
			break
		}
	}
	if refSound == nil {
		return fmt.Errorf("There's no reference sound in container %d to create new sound objects", s.Parent)
	}

	wems := make([]string, len(inputsMap))
	wsource, err := waapi.CreateConversionListInOrder(ctx, inputsMap, wems, s.Conversion, false)
	if len(inputsMap) != len(wems) {
		panic("Length of input maps does not equal to the length of output array")
	}
	if err != nil {
		return err
	}

	if err := waapi.WwiseConversion(ctx, wsource, proj); err != nil {
		return err
	}

	newAudioDatas := make([][]byte, 0, len(wems))
	for _, wem := range wems {
		audioData, err := os.ReadFile(wem)
		if err != nil {
			slog.Error(fmt.Sprintf("Failed to read audio data from %s", wem))
			continue
		}
		newAudioDatas = append(newAudioDatas, audioData)
	}

	// Hierarchy IDs generation and Source IDs generation
	q, closeConn, commit, rollback, err := db.CreateConnWithTxQuery(ctx)
	if err != nil {
		return err
	}
	defer closeConn()
	newSoundIDs := make([]uint32, len(newAudioDatas))
	newSourceIDs := make([]uint32, len(newAudioDatas))
	for i := range newAudioDatas {
		newSoundIDs[i], err = db.TryHid(ctx, q)
		if err != nil {
			rollback()
			return err
		}
		newSourceIDs[i], err = db.TrySid(ctx, q)
		if err != nil {
			rollback()
			return err
		}
	}

	for i, audioData := range newAudioDatas {
		if err := bnk.AppendAudio(audioData, newSourceIDs[i]); err != nil {
			rollback()
			return fmt.Errorf("Failed to add a new audio source file: %w.", err)
		}
	}
	bnk.ComputeDIDXOffset()
	if err := bnk.CheckDIDXDATA(); err != nil {
		rollback()
		return fmt.Errorf("Invalid Integrity appear in DIDX and DATA chunk: %w", err)
	}

	var newSound *wwise.Sound
	var pluginID uint32
	for i := range newAudioDatas {
		switch s.Format {
		case waapi.ConversionFormatTypePCM:
			pluginID = wwise.PCM
		case waapi.ConversionFormatTypeADPCM:
			pluginID = wwise.ADPCM
		case waapi.ConversionFormatTypeVORBIS:
			pluginID = wwise.VORBIS
		case waapi.ConversionFormatTypeWEMOpus:
			pluginID = wwise.WEM_OPUS
		}
		newSound = &wwise.Sound{
			Id: newSoundIDs[i],
			BankSourceData: wwise.BankSourceData{
				PluginID: pluginID,
				StreamType: wwise.STREAM_TYPE_BNK,
				SourceID: newSourceIDs[i],
				InMemoryMediaSize: uint32(len(newAudioDatas[i])),
				SourceBits: 0,
			},
			BaseParam: refSound.BaseParam.Clone(false),
		}
		if err := h.AppendNewSoundToRanSeqContainer(newSound, s.Parent, false); err != nil {
			rollback()
			return fmt.Errorf("Failed to add a new sound object to random / sequence container %d: %w", s.Parent, err)
		}
	}

	r := parent.(*wwise.RanSeqCntr)
	ImportRanSeqCntrModifer(r, &s, int(bnk.BKHD().BankGenerationVersion))

	if err := commit(); err != nil {
		rollback()
		return err
	}

	return nil
}

func ImportRanSeqCntrModifer(r *wwise.RanSeqCntr, s *ImportAsRanSeqCntrScript, ver int) {
	if s.Seq {
		r.PlayListSetting.Mode = wwise.ModeSequence
		r.PlayListSetting.SetResetPlayListAtEachPlay(false)
		r.ResetPlayListToLeafOrder()
	} else {
		r.PlayListSetting.Mode = wwise.ModeRandom
		r.PlayListSetting.SetResetPlayListAtEachPlay(true)
	}
	b := &r.BaseParam
	b.AdvanceSetting.SetIgnoreParentMaxNumInst(true)
	b.AdvanceSetting.MaxNumInstance = s.PlaybackLimit
	p := &b.PropBundle
	buf := make([]byte, 4) 

	if idx, in := p.HasPid(wwise.TMakeUpGain, ver); !in {
		binary.Encode(buf, wio.ByteOrder, s.MakeUpGain)
		p.AddWithVal(wwise.TMakeUpGain, [4]byte(buf), ver)
	} else {
		p.SetPropByIdxF32(idx, s.MakeUpGain)
	}
	if idx, in := p.HasPid(wwise.TInitialDelay, ver); !in {
		binary.Encode(buf, wio.ByteOrder, s.InitialDelay)
		p.AddWithVal(wwise.TInitialDelay, [4]byte(buf), ver)
	} else {
		p.SetPropByIdxF32(idx, s.InitialDelay)
	}
	if s.HDRActiveRange >= 0.0 {
		// Enable HDR Envelope and set HDR Active range
		b.SetEnableEnvelope(true, ver)
		i, _ := p.HDRActiveRange(ver)
		if i != -1 {
			p.SetPropByIdxF32(i, s.HDRActiveRange)
		}
	}
}
