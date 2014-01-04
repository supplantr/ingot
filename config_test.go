package ingot

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestRemoveSection(t *testing.T) {
	Convey("Given an existing *Config \"c\"", t, func() {
		Convey("After c.data[\"SectionOne\"] is removed", func() {
			c, _ := ReadBytes([]byte(testConfig))
			ok := c.RemoveSection("SectionOne")
			So(ok, ShouldBeTrue)
			Convey("c.data[\"SectionOne\"] should not exist", func() {
				_, ok := c.data["SectionOne"]
				So(ok, ShouldBeFalse)
			})
			Convey("\"SectionOne\" should not exist in c.sections", func() {
				ok := true
				for _, section := range c.sections {
					if section == "SectionOne" {
						ok = false
					}
				}
				So(ok, ShouldBeTrue)
			})
			Convey("c.options[\"SectionOne\"] should not exist", func() {
				_, ok := c.options["SectionOne"]
				So(ok, ShouldBeFalse)
			})
		})
	})
}

func TestRemoveOption(t *testing.T) {
	Convey("Given an existing *Config \"c\"", t, func() {
		Convey("After c.data[\"SectionOne\"][\"One\"] is removed", func() {
			c, _ := ReadBytes([]byte(testConfig))
			ok := c.RemoveOption("SectionOne", "One")
			So(ok, ShouldBeTrue)
			Convey("c.data[\"SectionOne\"][\"One\"] should not exist", func() {
				_, ok := c.data["SectionOne"]["One"]
				So(ok, ShouldBeFalse)
			})
			Convey("\"One\" should not exist in c.options[\"SectionOne\"]", func() {
				ok := true
				for _, option := range c.options["SectionOne"] {
					if option == "One" {
						ok = false
					}
				}
				So(ok, ShouldBeTrue)
			})
		})
	})
}

func TestConfigOrder(t *testing.T) {
	Convey("Given an existing *Config", t, func() {
		c, _ := ReadBytes([]byte(testConfig))
		s := "The %s %s should equal: %q"
		Convey(fmt.Sprintf(s, "first", "section", "SectionOne"), func() {
			section := c.sections[0]
			So(section, ShouldEqual, "SectionOne")
			options := c.options[section]
			Convey(fmt.Sprintf(s, "first", "option", "One"), func() {
				So(options[0], ShouldEqual, "One")
			})
			Convey(fmt.Sprintf(s, "second", "option", "Two"), func() {
				So(options[1], ShouldEqual, "Two")
			})
		})
		Convey(fmt.Sprintf(s, "second", "section", "SectionTwo"), func() {
			section := c.sections[1]
			So(section, ShouldEqual, "SectionTwo")
			options := c.options[section]
			Convey(fmt.Sprintf(s, "first", "option", "One"), func() {
				So(options[0], ShouldEqual, "One")
			})
		})
	})
}
