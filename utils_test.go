package cloudtoolkit

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestUtilSpec(t *testing.T) {

	Convey("Given you are on a machine with a network interface", t, func() {
		Convey("Get the local IP address", func() {
			ip := GetLocalIP()
			Convey("The IP should start with 192.168", func() {
				So(ip, ShouldStartWith, "192.168")
			})
		})
	})
}
