package xbacklight // import "mrogalski.eu/go/xbacklight"

import (
	"encoding/binary"
	"math"

	"github.com/BurntSushi/xgb/randr"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xprop"
)

type Backlighter interface {
	// Return current backlight as a number from 0 to 1 (both ends inclusive).
	Get() (float64, error)

	// Set the backlight to a value between 0 and 1.
	// 0 disables backlight completely.
	// Non-0 always produces some backlight.
	Set(float64) error
}

func NewBacklighterPrimaryScreen() (Backlighter, error) {
	x, err := xgbutil.NewConn()
	if err != nil {
		return nil, err
	}
	if err := randr.Init(x.Conn()); err != nil {
		return nil, err
	}
	output, err := primaryOutput(x)
	if err != nil {
		return nil, err
	}
	return NewBacklighter(x, output)
}

func NewBacklighter(x *xgbutil.XUtil, output randr.Output) (Backlighter, error) {
	min, max, err := backlightRange(x, output)
	if err != nil {
		return nil, err
	}
	atom, err := xprop.Atom(x, "Backlight", false)
	if err != nil {
		return nil, err
	}
	return &bundle{x, output, atom, min, max}, nil
}

type bundle struct {
	x        *xgbutil.XUtil
	output   randr.Output
	atom     xproto.Atom
	min, max int
}

func (b *bundle) Get() (float64, error) {
	prop, err := randr.GetOutputProperty(b.x.Conn(), b.output, b.atom, xproto.AtomNone, 0, 4, false, false).Reply()
	if err != nil {
		return 0, err
	}
	rawBacklight := int(binary.LittleEndian.Uint32(prop.Data))
	return float64(rawBacklight-b.min) / float64(b.max-b.min), nil

}

func (b *bundle) Set(backlight float64) error {
	data := make([]byte, 4)
	binary.LittleEndian.PutUint32(data, uint32(math.Ceil(backlight*float64(b.max-b.min)))+uint32(b.min))
	_, err := randr.ChangeOutputProperty(b.x.Conn(), b.output, b.atom, xproto.AtomInteger, 32, xproto.PropModeReplace, 1, data).Reply()
	return err
}

func primaryOutput(x *xgbutil.XUtil) (randr.Output, error) {
	primary, err := randr.GetOutputPrimary(x.Conn(), x.RootWin()).Reply()
	if err != nil {
		return 0, err
	}
	return primary.Output, nil
}

func backlightRange(x *xgbutil.XUtil, output randr.Output) (int, int, error) {
	atom, err := xprop.Atom(x, "Backlight", false)
	if err != nil {
		return 0, 0, err
	}
	query, err := randr.QueryOutputProperty(x.Conn(), output, atom).Reply()
	if err != nil {
		return 0, 0, err
	}
	return int(query.ValidValues[0]), int(query.ValidValues[1]), nil
}
