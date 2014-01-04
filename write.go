package ingot

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

func lower(s string) string {
	return strings.ToLower(s[:1]) + s[1:]
}

// Write writes the configuration to an io.Writer.
func (c *Config) Write(writer io.Writer) error {
	for _, section := range c.sections {
		line := fmt.Sprintf("[%s]\n", lower(section))
		if _, err := writer.Write([]byte(line)); err != nil {
			return err
		}
		for _, option := range c.options[section] {
			value := c.data[section][option]
			line := fmt.Sprintf("%s=%s\n", lower(option), value)
			if _, err := writer.Write([]byte(line)); err != nil {
				return err
			}
		}
	}

	return nil
}

// WriteFile writes the configuration to a file.
func (c *Config) WriteFile(fname string) error {
	file, err := os.Create(fname)
	if err != nil {
		return err
	}
	if err := c.Write(file); err != nil {
		return err
	}

	return file.Close()
}

// WriteBytes writes the configuration to a slice of bytes.
func (c *Config) WriteBytes() (config []byte, err error) {
	buf := bytes.NewBuffer(nil)
	if err = c.Write(buf); err != nil {
		return buf.Bytes(), err
	}
	return buf.Bytes(), nil
}
