package docpdf

import (
	"fmt"
	"io"
	"sync"
)

type sMaskSubtypes string

const (
	sMaskAlphaSubtype      = "/Alpha"
	sMaskLuminositySubtype = "/luminosity"
)

// sMask smask
type sMask struct {
	imgInfo
	data []byte
	//getRoot func() *pdfEngine
	pdfProtection                 *pdfProtection
	Index                         int
	TransparencyXObjectGroupIndex int
	S                             string
}

type sMaskOptions struct {
	TransparencyXObjectGroupIndex int
	Subtype                       sMaskSubtypes
}

func (smask sMaskOptions) GetId() string {
	id := fmt.Sprintf("S_%s;G_%d_0_R", smask.Subtype, smask.TransparencyXObjectGroupIndex)

	return id
}

func getCachedMask(opts sMaskOptions, gp *pdfEngine) sMask {
	smask, ok := gp.curr.sMasksMap.Find(opts)
	if !ok {
		smask = sMask{
			S:                             string(opts.Subtype),
			TransparencyXObjectGroupIndex: opts.TransparencyXObjectGroupIndex,
		}
		smask.Index = gp.addObj(smask)

		gp.curr.sMasksMap.Save(opts.GetId(), smask)
	}

	return smask
}

func (s sMask) init(func() *pdfEngine) {}

func (s *sMask) setProtection(p *pdfProtection) {
	s.pdfProtection = p
}

func (s sMask) protection() *pdfProtection {
	return s.pdfProtection
}

func (s sMask) getType() string {
	return "Mask"
}

func (s sMask) write(w io.Writer, objID int) error {
	if s.TransparencyXObjectGroupIndex != 0 {
		content := "<<\n"
		content += "\t/Type /Mask\n"
		content += fmt.Sprintf("\t/S %s\n", s.S)
		content += fmt.Sprintf("\t/G %d 0 R\n", s.TransparencyXObjectGroupIndex+1)
		content += ">>\n"

		if _, err := io.WriteString(w, content); err != nil {
			return err
		}
	} else {
		err := writeImgProps(w, s.imgInfo, false)
		if err != nil {
			return err
		}

		fmt.Fprintf(w, "/Length %d\n>>\n", len(s.data)) // /Length 62303>>\n
		io.WriteString(w, "stream\n")
		if s.protection() != nil {
			tmp, err := rc4Cip(s.protection().objectkey(objID), s.data)
			if err != nil {
				return err
			}
			w.Write(tmp)
			io.WriteString(w, "\n")
		} else {
			w.Write(s.data)
		}
		io.WriteString(w, "\nendstream\n")
	}

	return nil
}

type sMaskMap struct {
	syncer sync.Mutex
	table  map[string]sMask
}

func newSMaskMap() sMaskMap {
	return sMaskMap{
		syncer: sync.Mutex{},
		table:  make(map[string]sMask),
	}
}

func (sm *sMaskMap) Find(sM sMaskOptions) (sMask, bool) {
	key := sM.GetId()

	sm.syncer.Lock()
	defer sm.syncer.Unlock()

	t, ok := sm.table[key]
	if !ok {
		return sMask{}, false
	}

	return t, ok

}

func (smask *sMaskMap) Save(id string, sMask sMask) sMask {
	smask.syncer.Lock()
	defer smask.syncer.Unlock()

	smask.table[id] = sMask

	return sMask
}
