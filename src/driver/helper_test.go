package driver

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestCleanSecret(t *testing.T) {
	Convey("give some secret string", t, func() {
		secret1 := "123"
		secret2 := "12345678"

		Convey("test it", func() {
			So(cleanSecret(secret1), ShouldEqual, "***")
			So(cleanSecret(secret2), ShouldEqual, "123*****")
		})
	})
}
