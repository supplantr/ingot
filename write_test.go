package ingot

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestWriteBytes(t *testing.T) {
	Convey("Given an existing *Config", t, func() {
		Convey("After bytes are written", func() {
			in := []byte(testConfig)
			c, _ := ReadBytes(in)
			out, err := c.WriteBytes()
			So(err, ShouldBeNil)
			Convey("Output bytes should equal input bytes", func() {
				So(out, ShouldResemble, in)
			})
		})
	})
}
