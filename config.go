// Package ingot implements a simple parser and writer for ini-style configuration files.
// Structs can be used to generate configurations, and configurations can be used to populate structs.
package ingot

import (
	"strings"
)

// A Config is a representation of configuration settings.
type Config struct {
	sections []string
	options  map[string][]string
	data     map[string]map[string]string
}

// AddSection adds a new section to the configuration.
// It returns true if the new section was added, and false if the section already exists.
func (c *Config) AddSection(section string) bool {
	section = strings.Title(section)

	if _, ok := c.data[section]; ok {
		return false
	}
	c.data[section] = make(map[string]string)
	c.sections = append(c.sections, section)

	return true
}

// RemoveSection removes a section from the configuration.
// It returns true if the section was removed, and false if the section did not exist.
func (c *Config) RemoveSection(section string) bool {
	section = strings.Title(section)

	if _, ok := c.data[section]; !ok {
		return false
	}
	delete(c.data, section)

	for i, s := range c.sections {
		if s == section {
			c.sections = append(c.sections[:i], c.sections[i+1:]...)
			break
		}
	}
	delete(c.options, section)

	return true
}

// AddOption adds a new option and value to the specified section of the configuration.
// It returns true if the option and value were added, and false if the value was overwritten.
// If the section does not exist, it is created.
func (c *Config) AddOption(section, option, value string) bool {
	c.AddSection(section)

	section = strings.Title(section)
	option = strings.Title(option)

	_, ok := c.data[section][option]
	c.data[section][option] = value

	if !ok {
		c.options[section] = append(c.options[section], option)
	}

	return !ok
}

// RemoveOption removes an option from the specified section of the configuration.
// It returns true if the option was removed, and false if the option (or the section) did not exist.
func (c *Config) RemoveOption(section, option string) bool {
	option = strings.Title(option)

	if _, ok := c.data[section][option]; !ok {
		return false
	}
	delete(c.data[section], option)

	for i, opt := range c.options[section] {
		if opt == option {
			c.options[section] = append(c.options[section][:i], c.options[section][i+1:]...)
			break
		}
	}

	return true
}

// New creates an empty configuration representation.
func New() *Config {
	c := new(Config)
	c.sections = []string{}
	c.options = map[string][]string{}
	c.data = make(map[string]map[string]string)

	return c
}
