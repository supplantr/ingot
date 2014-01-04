package ingot

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	// Parse Errors
	NoSection = iota
	CouldNotParse
)

// Read reads an io.Reader and generates the configuration.
func (c *Config) Read(reader io.Reader) error {
	var section, option string

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		switch {
		case len(line) == 0, line[0] == '#', line[0] == ';':
			continue
		case line[0] == '[' && line[len(line)-1] == ']':
			if s := strings.TrimSpace(line[1 : len(line)-1]); s != "" {
				section = s
				c.AddSection(section)
			}
		case section == "":
			return ParseError{NoSection, line}
		default:
			i := strings.Index(line, "=")
			if i > 0 {
				option = strings.TrimSpace(line[:i])
				value := strings.TrimSpace(line[i+1:])
				c.AddOption(section, option, value)
			} else {
				return ParseError{CouldNotParse, line}
			}
		}
	}

	return nil
}

// ReadBytes reads a slice of bytes and returns a configuration representation.
func ReadBytes(config []byte) (c *Config, err error) {
	buf := bytes.NewBuffer(config)
	c = New()
	if err = c.Read(buf); err != nil {
		return nil, err
	}

	return c, nil
}

// ReadFile reads a file and returns a configuration representation.
func ReadFile(fname string) (c *Config, err error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	c = New()
	if err = c.Read(file); err != nil {
		return nil, err
	}

	return c, file.Close()
}

type ParseError struct {
	Reason int
	Line   string
}

func (e ParseError) Error() string {
	switch e.Reason {
	case NoSection:
		return fmt.Sprintf("no sections present before: %q", e.Line)
	case CouldNotParse:
		return fmt.Sprintf("could not parse text: %q", e.Line)
	default:
		return "unknown parse error"
	}
}
