package ingot

import (
	"fmt"
	"reflect"
	"strconv"
)

const (
	// Struct Errors
	NonexistentField = iota
	FieldIsNotStruct
	CannotSet
	UnsupportedType
)

const (
	// Section Errors
	NonexistentSection = iota
	SectionExists
)

func bitSize(k reflect.Value) int {
	var size int

	switch k.Kind() {
	case reflect.Int, reflect.Uint:
		size = 0
	case reflect.Int8, reflect.Uint8:
		size = 8
	case reflect.Int16, reflect.Uint16:
		size = 16
	case reflect.Int32, reflect.Uint32, reflect.Float32:
		size = 32
	case reflect.Int64, reflect.Uint64, reflect.Float64:
		size = 64
	}

	return size
}

func setStructField(f reflect.Value, option, value string, strict bool) error {
	if !f.IsValid() {
		if strict {
			return StructError{NonexistentField, option}
		}
		return nil
	}
	if !f.CanSet() {
		if strict {
			return StructError{CannotSet, option}
		}
		return nil
	}
	switch f.Kind() {
	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			if strict {
				return FormatError{option, value, err}
			}
			break
		}
		f.SetBool(b)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(value, 0, bitSize(f))
		if err != nil {
			if strict {
				return FormatError{option, value, err}
			}
			break
		}
		f.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u, err := strconv.ParseUint(value, 0, bitSize(f))
		if err != nil {
			if strict {
				return FormatError{option, value, err}
			}
			break
		}
		f.SetUint(u)
	case reflect.Float32, reflect.Float64:
		a, err := strconv.ParseFloat(value, bitSize(f))
		if err != nil {
			if strict {
				return FormatError{option, value, err}
			}
			break
		}
		f.SetFloat(a)
	case reflect.String:
		f.SetString(value)
	}

	return nil
}

// SectionToStruct populates a struct with the specified section of the configuration.
// The struct should contain fields with types corresponding to the section's options and values.
// If strict is false, it will continue despite errors and return nil.
func (c *Config) SectionToStruct(section string, sptr interface{}, strict bool) error {
	opts, ok := c.data[section]
	if !ok {
		return SectionError{NonexistentSection, section}
	}
	e := reflect.ValueOf(sptr).Elem()
	for option, value := range opts {
		f := e.FieldByName(option)
		if err := setStructField(f, option, value, strict); err != nil {
			return err
		}
	}

	return nil
}

// ToStruct populates a struct with the configuration.
// The struct should contain embedded structs corresponding to the configuration's sections.
// These embedded structs should contain fields with types corresponding to the related section's options and values.
// If strict is false, it will continue despite errors and return nil.
func (c *Config) ToStruct(sptr interface{}, strict bool) error {
	e := reflect.ValueOf(sptr).Elem()
	for section, opts := range c.data {
		s := e.FieldByName(section)
		if !s.IsValid() {
			if strict {
				return StructError{NonexistentField, section}
			}
			continue
		}
		if s.Kind() != reflect.Struct {
			if strict {
				return StructError{FieldIsNotStruct, section}
			}
			continue
		}
		for option, value := range opts {
			f := s.FieldByName(option)
			if err := setStructField(f, option, value, strict); err != nil {
				return err
			}
		}
	}

	return nil
}

func getStructField(f reflect.Value) (string, bool) {
	var value string

	switch f.Kind() {
	case reflect.Bool:
		value = strconv.FormatBool(f.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value = strconv.FormatInt(f.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		value = strconv.FormatUint(f.Uint(), 10)
	case reflect.Float32:
		value = strconv.FormatFloat(f.Float(), 'f', -1, 32)
	case reflect.Float64:
		value = strconv.FormatFloat(f.Float(), 'f', -1, 64)
	case reflect.String:
		value = f.String()
	default:
		return "", false
	}

	return value, true
}

// SectionFromStruct generates a section in the configuration from the specified struct.
func (c *Config) SectionFromStruct(section string, s interface{}) error {
	if ok := c.AddSection(section); !ok {
		return SectionError{SectionExists, section}
	}
	v := reflect.ValueOf(s)
	for i := 0; i < v.NumField(); i++ {
		option := reflect.TypeOf(s).Field(i).Name
		value, ok := getStructField(v.Field(i))
		if !ok {
			return StructError{UnsupportedType, option}
		}
		c.AddOption(section, option, value)
	}

	return nil
}

// FromStruct generates the entire configuration from the specified struct.
func (c *Config) FromStruct(s interface{}) error {
	v := reflect.ValueOf(s)
	for i := 0; i < v.NumField(); i++ {
		section := v.Type().Field(i).Name
		if ok := c.AddSection(section); !ok {
			return SectionError{SectionExists, section}
		}
		f := v.Field(i)
		for j := 0; j < f.NumField(); j++ {
			option := v.Type().FieldByIndex([]int{i, j}).Name
			value, ok := getStructField(f.Field(j))
			if !ok {
				return StructError{UnsupportedType, option}
			}
			c.AddOption(section, option, value)
		}
	}

	return nil
}

type StructError struct {
	Reason int
	Field  string
}

func (e StructError) Error() string {
	switch e.Reason {
	case NonexistentField:
		return fmt.Sprintf("field %q does not exist", e.Field)
	case FieldIsNotStruct:
		return fmt.Sprintf("field %q is not a struct", e.Field)
	case CannotSet:
		return fmt.Sprintf("cannot set field %q", e.Field)
	case UnsupportedType:
		return fmt.Sprintf("type of %q is unsupported", e.Field)
	default:
		return "unknown struct error"
	}
}

type FormatError struct {
	Field, Value string
	Err          error
}

func (e FormatError) Error() string {
	return fmt.Sprintf("value %q for option %q is malformed: %s", e.Value, e.Field, e.Err)
}

type SectionError struct {
	Reason  int
	Section string
}

func (e SectionError) Error() string {
	switch e.Reason {
	case NonexistentSection:
		return fmt.Sprintf("section %q does not exist", e.Section)
	case SectionExists:
		return fmt.Sprintf("section %q already exists", e.Section)
	default:
		return "unknown section error"
	}
}
