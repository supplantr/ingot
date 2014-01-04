package ingot

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

const (
	testConfig  = "[sectionOne]\none=true\ntwo=false\n[sectionTwo]\none=false\n"
	typesConfig = `
[types]
bool=true
int=-1
int8=-127
int16=-128
int32=-32768
int64=-2147483648
uint=1
uint8=255
uint16=256
uint32=65536
uint64=4294967296
float32=1.0
float64=1.0
string=test
`
)

type SectionOne struct {
	One bool
	Two bool
}

type SectionTwo struct {
	One bool
}

type testStruct struct {
	SectionOne
	SectionTwo
}

type typesStruct struct {
	Bool    bool
	Int     int
	Int8    int8
	Int16   int16
	Int32   int32
	Int64   int64
	Uint    uint
	Uint8   uint8
	Uint16  uint16
	Uint32  uint32
	Uint64  uint64
	Float32 float32
	Float64 float64
	String  string
}

func TestToTypes(t *testing.T) {
	Convey("Given a new struct \"test\"", t, func() {
		Convey("It should be populated by *Config.data[\"Types\"]", func() {
			test := new(typesStruct)
			c, _ := ReadBytes([]byte(typesConfig))
			err := c.SectionToStruct("Types", test, true)
			So(err, ShouldBeNil)
		})
	})
}

func TestFromTypes(t *testing.T) {
	Convey("Given a new *Config \"c\"", t, func() {
		Convey("c.data[\"Types\"] should be generated from a struct", func() {
			c := New()
			test := typesStruct{
				true,
				-1,
				-127,
				-128,
				-32768,
				-2147483648,
				1,
				255,
				256,
				65536,
				4294967296,
				1.0,
				1.0,
				"test",
			}
			err := c.SectionFromStruct("Types", test)
			So(err, ShouldBeNil)
		})
	})
}

func TestSectionToStruct(t *testing.T) {
	Convey("Given a new struct \"test\"", t, func() {
		Convey("It should be populated by *Config.data[\"SectionOne\"]", func() {
			test := new(SectionOne)
			c, _ := ReadBytes([]byte(testConfig))
			err := c.SectionToStruct("SectionOne", test, true)
			So(err, ShouldBeNil)
			Convey("Field One should be: true", func() {
				So(test.One, ShouldBeTrue)
			})
			Convey("Field Two should be: false", func() {
				So(test.Two, ShouldBeFalse)
			})
		})
	})
}

func TestToStruct(t *testing.T) {
	Convey("Given a new struct \"test\" with embedded structs", t, func() {
		Convey("It should be populated by *Config.data", func() {
			test := new(testStruct)
			c, _ := ReadBytes([]byte(testConfig))
			err := c.ToStruct(test, true)
			So(err, ShouldBeNil)
			Convey("For test.SectionOne", func() {
				Convey("Field One should be: true", func() {
					So(test.SectionOne.One, ShouldBeTrue)
				})
				Convey("Field Two should be: false", func() {
					So(test.SectionOne.Two, ShouldBeFalse)
				})
			})
			Convey("For test.SectionTwo", func() {
				Convey("Field One should be: false", func() {
					So(test.SectionTwo.One, ShouldBeFalse)
				})
			})
		})
	})
}

func TestSectionFromStruct(t *testing.T) {
	Convey("Given a new *Config \"c\"", t, func() {
		Convey("c.data[\"SectionOne\"] should be generated from a struct", func() {
			c := New()
			test := SectionOne{true, false}
			err := c.SectionFromStruct("SectionOne", test)
			So(err, ShouldBeNil)
			s := "[%q] should equal: %q"
			Convey(fmt.Sprintf(s, "One", "true"), func() {
				So(c.data["SectionOne"]["One"], ShouldEqual, "true")
			})
			Convey(fmt.Sprintf(s, "Two", "false"), func() {
				So(c.data["SectionOne"]["Two"], ShouldEqual, "false")
			})
		})
	})
}

func TestFromStruct(t *testing.T) {
	Convey("Given a new *Config \"c\"", t, func() {
		Convey("c.data should be generated from a struct with embedded structs", func() {
			c := New()
			test := testStruct{SectionOne{true, false}, SectionTwo{false}}
			err := c.FromStruct(test)
			So(err, ShouldBeNil)
			Convey("For c.data[\"SectionOne\"]", func() {
				s := "[%q] should equal: %q"
				Convey(fmt.Sprintf(s, "One", "true"), func() {
					So(c.data["SectionOne"]["One"], ShouldEqual, "true")
				})
				Convey(fmt.Sprintf(s, "Two", "false"), func() {
					So(c.data["SectionOne"]["Two"], ShouldEqual, "false")
				})
			})
			Convey("For c.data[\"SectionTwo\"]", func() {
				Convey("[\"One\"] should equal: \"false\"", func() {
					So(c.data["SectionTwo"]["One"], ShouldEqual, "false")
				})
			})
		})
	})
}
