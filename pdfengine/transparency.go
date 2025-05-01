package pdfengine

import (
	"fmt"
	"sync"

	"github.com/cdvelop/docpdf/errs"
)

type blendModeType string

const (
	hue             blendModeType = "/hue"
	colors          blendModeType = "/Colors"
	normalBlendMode blendModeType = "/Normal"
	darken          blendModeType = "/darken"
	screen          blendModeType = "/screen"
	overlay         blendModeType = "/overlay"
	lighten         blendModeType = "/lighten"
	multiply        blendModeType = "/multiply"
	exclusion       blendModeType = "/exclusion"
	colorBurn       blendModeType = "/colorBurn"
	hardLight       blendModeType = "/hardLight"
	softLight       blendModeType = "/softLight"
	difference      blendModeType = "/difference"
	saturation      blendModeType = "/saturation"
	luminosity      blendModeType = "/luminosity"
	colorDodge      blendModeType = "/colorDodge"
)

const defaultAplhaValue = 1

// transparency defines an object alpha.
type transparency struct {
	extGStateIndex int
	Alpha          float64
	blendModeType  blendModeType
}

func newTransparency(alpha float64, blendModeType string) (transparency, error) {
	if alpha < 0.0 || alpha > 1.0 {
		return transparency{}, fmt.Errorf("alpha value is out of range (0.0 - 1.0): %.3f", alpha)
	}

	bmtType, err := defineBlendModeType(blendModeType)
	if err != nil {
		return transparency{}, err
	}

	return transparency{
		Alpha:         alpha,
		blendModeType: bmtType,
	}, nil
}

func (t transparency) GetId() string {
	keyStr := fmt.Sprintf("%.3f_%s", t.Alpha, t.blendModeType)

	return keyStr
}

type transparencyMap struct {
	syncer sync.Mutex
	table  map[string]transparency
}

func newTransparencyMap() transparencyMap {
	return transparencyMap{
		syncer: sync.Mutex{},
		table:  make(map[string]transparency),
	}
}

func (tm *transparencyMap) Find(transp transparency) (transparency, bool) {
	key := transp.GetId()

	tm.syncer.Lock()
	defer tm.syncer.Unlock()

	t, ok := tm.table[key]
	if !ok {
		return transparency{}, false
	}

	return t, ok

}

func (tm *transparencyMap) Save(transparency transparency) transparency {
	tm.syncer.Lock()
	defer tm.syncer.Unlock()

	key := transparency.GetId()
	tm.table[key] = transparency

	return transparency
}

func defineBlendModeType(bmType string) (blendModeType, error) {
	switch bmType {
	case string(hue):
		return hue, nil
	case string(colors):
		return colors, nil
	case "", string(normalBlendMode):
		return normalBlendMode, nil
	case string(darken):
		return darken, nil
	case string(screen):
		return screen, nil
	case string(overlay):
		return overlay, nil
	case string(lighten):
		return lighten, nil
	case string(multiply):
		return multiply, nil
	case string(exclusion):
		return exclusion, nil
	case string(colorBurn):
		return colorBurn, nil
	case string(hardLight):
		return hardLight, nil
	case string(softLight):
		return softLight, nil
	case string(difference):
		return difference, nil
	case string(saturation):
		return saturation, nil
	case string(luminosity):
		return luminosity, nil
	case string(colorDodge):
		return colorDodge, nil
	default:
		return "", errs.New("blend mode is unknown")
	}
}
