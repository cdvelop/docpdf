package pdfengine

import (
	"fmt"
	"io"
	"sync"

	"errors"
)

// TODO: add all fields https://www.adobe.com/content/dam/acom/en/devnet/acrobat/pdfs/PDF32000_2008.pdf 8.4.5 page 128
type extGState struct {
	Index      int
	ca         *float64
	CA         *float64
	BM         *blendModeType
	SMaskIndex *int
}

type extGStateOptions struct {
	StrokingCA    *float64
	NonStrokingCa *float64
	BlendMode     *blendModeType
	SMaskIndex    *int
}

func (extOpt extGStateOptions) GetId() string {
	id := ""
	if extOpt.StrokingCA != nil {
		id += fmt.Sprintf("CA_%.3f;", *extOpt.StrokingCA)
	}
	if extOpt.NonStrokingCa != nil {
		id += fmt.Sprintf("ca_%.3f;", *extOpt.NonStrokingCa)
	}
	if extOpt.BlendMode != nil {
		id += fmt.Sprintf("BM_%s;", *extOpt.BlendMode)
	}
	if extOpt.SMaskIndex != nil {
		id += fmt.Sprintf("SMask_%d_0_R;", *extOpt.SMaskIndex)
	}

	return id
}

func getCachedExtGState(opts extGStateOptions, gp *PdfEngine) (extGState, error) {
	state, ok := gp.curr.extGStatesMap.Find(opts)
	if !ok {
		state = extGState{
			BM:         opts.BlendMode,
			CA:         opts.StrokingCA,
			ca:         opts.NonStrokingCa,
			SMaskIndex: opts.SMaskIndex,
		}

		state.Index = gp.addObj(state)

		pdfObj := gp.pdfObjs[gp.indexOfProcSet]
		procSet, ok := pdfObj.(extGState)
		if !ok {
			return extGState{}, errors.New("invalid PDF object type")
		}

		procSet.Index = state.Index

		gp.curr.extGStatesMap.Save(opts.GetId(), state)
	}

	return state, nil
}

func (egs extGState) Init(func() *PdfEngine) {}

func (egs extGState) GetType() string {
	return "extGState"
}

func (egs extGState) Write(w Writer, objID int) error {
	content := "<<\n"
	content += "\t/Type /extGState\n"

	if egs.ca != nil {
		content += fmt.Sprintf("\t/ca %.3F\n", *egs.ca)
	}
	if egs.CA != nil {
		content += fmt.Sprintf("\t/CA %.3F\n", *egs.CA)
	}
	if egs.BM != nil {
		content += fmt.Sprintf("\t/BM %s\n", *egs.BM)
	}

	if egs.SMaskIndex != nil {
		content += fmt.Sprintf("\t/sMask %d 0 R\n", *egs.SMaskIndex+1)
	}

	content += ">>\n"

	if _, err := io.WriteString(w, content); err != nil {
		return err
	}

	return nil
}

type extGStatesMap struct {
	syncer sync.Mutex
	table  map[string]extGState
}

func newExtGStatesMap() extGStatesMap {
	return extGStatesMap{
		syncer: sync.Mutex{},
		table:  make(map[string]extGState),
	}
}

func (m *extGStatesMap) Find(opts extGStateOptions) (extGState, bool) {
	key := opts.GetId()

	m.syncer.Lock()
	defer m.syncer.Unlock()

	t, ok := m.table[key]
	if !ok {
		return extGState{}, false
	}

	return t, ok
}

func (tm *extGStatesMap) Save(id string, extGState extGState) extGState {
	tm.syncer.Lock()
	defer tm.syncer.Unlock()

	tm.table[id] = extGState

	return extGState
}
