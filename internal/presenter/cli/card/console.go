package card

import (
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/manifoldco/promptui"
)

var selectTemplate = &promptui.SelectTemplates{
	Label:    ">> {{ . }}",
	Active:   promptui.IconSelect + " {{ . | cyan }}",
	Inactive: "  {{ . | cyan }} ({{ .HeatUnit | red }})",
	Selected: promptui.IconGood + " {{ . | green }}",
}

var errInterrupt = fmt.Errorf("%w", promptui.ErrInterrupt)
var errEOF = fmt.Errorf("%w", promptui.ErrEOF)
var errAbort = fmt.Errorf("%w", promptui.ErrAbort)

type console struct {
	writer io.StringWriter
}

func newConsole(writer io.StringWriter) console {
	return console{
		writer: writer,
	}
}

func (c console) WriteString(data string) (int, error) {
	return c.writer.WriteString(data)
}

func (c console) Select(label string, items []string) (int, error) {
	prompt := promptui.Select{
		Label:     label,
		Templates: selectTemplate,
		HideHelp:  true,
		Items:     items,
	}

	resultIdx, _, err := prompt.Run()

	return resultIdx, mapError(err)
}

func (c console) Prompt(label, defaultVal string) (string, error) {
	prompt := promptui.Prompt{
		Label:   label,
		Default: defaultVal,
	}

	res, err := prompt.Run()

	return res, mapError(err)
}

func (c console) Busy(done <-chan struct{}) <-chan struct{} {
	ret := make(chan struct{})
	go func() {
		defer close(ret)
		for {
			c.writer.WriteString(".")
			select {
			case <-done:
				return
			default:
			}
			time.Sleep(1 * time.Second)
		}
	}()

	return ret
}

func mapError(err error) error {
	if errors.Is(err, promptui.ErrInterrupt) {
		return errInterrupt
	} else if errors.Is(err, promptui.ErrAbort) {
		return errAbort
	} else if errors.Is(err, promptui.ErrEOF) {
		return errEOF
	}

	return err
}
