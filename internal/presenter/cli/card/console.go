package card

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/rs/zerolog/log"
)

const selectLinesCount = 10

var selectTemplate = &promptui.SelectTemplates{
	Label:    ">> {{ . }}",
	Active:   promptui.IconSelect + " {{ . | cyan }}",
	Inactive: "  {{ . | cyan }} ({{ .HeatUnit | red }})",
	Selected: promptui.IconGood + " {{ . | green }}",
}

var (
	errInterrupt = fmt.Errorf("%w", promptui.ErrInterrupt)
	errEOF       = fmt.Errorf("%w", promptui.ErrEOF)
	errAbort     = fmt.Errorf("%w", promptui.ErrAbort)
)

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
	searcher := func(input string, index int) bool {
		item := items[index]
		name := strings.ToLower(item)
		input = strings.ToLower(input)

		return strings.Contains(name, input)
	}

	prompt := promptui.Select{
		Label:     label,
		Templates: selectTemplate,
		HideHelp:  true,
		Items:     items,
		Size:      selectLinesCount,
		Searcher:  searcher,
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
	ticker := time.NewTicker(1 * time.Second)

	go func() {
		defer close(ret)
		defer ticker.Stop()

		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				_, err := c.writer.WriteString(".")
				if err != nil {
					log.Error().Err(err).Msg("failed to write waiting string on console")
				}
			}
		}
	}()

	return ret
}

func mapError(err error) error {
	if errors.Is(err, promptui.ErrInterrupt) {
		return errInterrupt
	}

	if errors.Is(err, promptui.ErrAbort) {
		return errAbort
	}

	if errors.Is(err, promptui.ErrEOF) {
		return errEOF
	}

	return err
}
