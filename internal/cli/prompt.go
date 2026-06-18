// Package cli, interaktif terminal girdisi için yardımcılar sağlar.
package cli

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Prompter, test edilebilir bir soru-cevap arayüzüdür.
type Prompter struct {
	r *bufio.Reader
	w io.Writer
}

// New, yeni bir Prompter oluşturur.
func New(in io.Reader, out io.Writer) *Prompter {
	return &Prompter{r: bufio.NewReader(in), w: out}
}

// AskInt, varsayılanlı bir tam sayı sorar. Boş girdi → def.
func (p *Prompter) AskInt(question string, def int) (int, error) {
	fmt.Fprintf(p.w, "%s [%d]: ", question, def)
	line, err := p.r.ReadString('\n')
	if err != nil && err != io.EOF {
		return 0, err
	}
	line = strings.TrimSpace(line)
	if line == "" {
		return def, nil
	}
	n, err := strconv.Atoi(line)
	if err != nil {
		return 0, fmt.Errorf("geçersiz sayı: %q", line)
	}
	return n, nil
}

// AskYesNo, evet/hayır sorar. Boş girdi → def.
func (p *Prompter) AskYesNo(question string, def bool) (bool, error) {
	d := "E/h"
	if !def {
		d = "e/H"
	}
	fmt.Fprintf(p.w, "%s (%s): ", question, d)
	line, err := p.r.ReadString('\n')
	if err != nil && err != io.EOF {
		return false, err
	}
	line = strings.ToLower(strings.TrimSpace(line))
	switch line {
	case "":
		return def, nil
	case "e", "evet", "y", "yes":
		return true, nil
	case "h", "hayir", "hayır", "n", "no":
		return false, nil
	default:
		return def, nil
	}
}

// AskLine, serbest metin sorar (varsayılansız).
func (p *Prompter) AskLine(question string) (string, error) {
	fmt.Fprintf(p.w, "%s: ", question)
	line, err := p.r.ReadString('\n')
	if err != nil && err != io.EOF {
		return "", err
	}
	return strings.TrimSpace(line), nil
}
