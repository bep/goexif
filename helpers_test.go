package imagemeta

import (
	"encoding"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestStringer(t *testing.T) {
	c := qt.New(t)
	c.Assert(exifTypeUnsignedByte1.String(), qt.Equals, "exifTypeUnsignedByte")

	var source Source
	c.Assert(EXIF.String(), qt.Equals, "EXIF")
	c.Assert(IPTC.String(), qt.Equals, "IPTC")
	c.Assert(XMP.String(), qt.Equals, "XMP")
	c.Assert(source.String(), qt.Equals, "Source(0)")

	var imageFormatAuto ImageFormat
	var imageFormat42 ImageFormat = 42
	c.Assert(JPEG.String(), qt.Equals, "JPEG")
	c.Assert(PNG.String(), qt.Equals, "PNG")
	c.Assert(TIFF.String(), qt.Equals, "TIFF")
	c.Assert(WebP.String(), qt.Equals, "WebP")
	c.Assert(imageFormatAuto.String(), qt.Equals, "ImageFormatAuto")
	c.Assert(imageFormat42.String(), qt.Equals, "ImageFormat(42)")
}

func BenchmarkPrintableString(b *testing.B) {
	runBench := func(b *testing.B, name, s string) {
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = printableString(s)
			}
		})
	}

	runBench(b, "ASCII", "Hello, World!")
	runBench(b, "ASCII with whitespace", "   Hello, World!   ")
	runBench(b, "UTF-8", "Hello, 世界!")
	runBench(b, "Mixed", "Hello, 世界! 🌍")
	runBench(b, "Unprintable", "Hello, \x00World!")
}

func TestRat(t *testing.T) {
	c := qt.New(t)

	c.Run("NewRat", func(c *qt.C) {
		ru := NewRat[uint32](1, 2)
		c.Assert(ru.Num(), qt.Equals, uint32(1))
		c.Assert(ru.Den(), qt.Equals, uint32(2))

		ri := NewRat[int32](1, 2)
		c.Assert(ri.Num(), qt.Equals, int32(1))
		c.Assert(ri.Den(), qt.Equals, int32(2))

		c.Assert(func() { NewRat[int32](10, 0) }, qt.PanicMatches, "division by zero")

		// Normalization
		// Denominator must be positive.
		ri = NewRat[int32](13, -3)
		c.Assert(ri.Num(), qt.Equals, int32(-13))
		c.Assert(ri.Den(), qt.Equals, int32(3))
		// Remove the greatest common divisor.
		ri = NewRat[int32](6, 9)
		c.Assert(ri.Num(), qt.Equals, int32(2))
		c.Assert(ri.Den(), qt.Equals, int32(3))
		ri = NewRat[int32](90, 600)
		c.Assert(ri.Num(), qt.Equals, int32(3))
		c.Assert(ri.Den(), qt.Equals, int32(20))
	})

	c.Run("MarshalText", func(c *qt.C) {
		ru := NewRat[uint32](1, 2)
		text, err := ru.(encoding.TextMarshaler).MarshalText()
		c.Assert(err, qt.Equals, nil)
		c.Assert(string(text), qt.Equals, "1/2")
	})

	c.Run("UnmarshalText", func(c *qt.C) {
		ru := NewRat[uint32](1, 2)
		err := ru.(encoding.TextUnmarshaler).UnmarshalText([]byte("3/4"))
		c.Assert(err, qt.Equals, nil)
		c.Assert(ru.Num(), qt.Equals, uint32(3))
		c.Assert(ru.Den(), qt.Equals, uint32(4))

		err = ru.(encoding.TextUnmarshaler).UnmarshalText([]byte("4"))
		c.Assert(err, qt.Equals, nil)
		c.Assert(ru.Num(), qt.Equals, uint32(4))
		c.Assert(ru.Den(), qt.Equals, uint32(1))
	})

	c.Run("String", func(c *qt.C) {
		ru := NewRat[uint32](1, 2)
		c.Assert(ru.String(), qt.Equals, "1/2")
		ru = NewRat[uint32](4, 1)
		c.Assert(ru.String(), qt.Equals, "4")
	})
}
