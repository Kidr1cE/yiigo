package yiigo

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDate(t *testing.T) {
	assert.Equal(t, "2016-03-19 15:03:19", Date(1458370999))
}

func TestStrToTime(t *testing.T) {
	assert.Equal(t, int64(1562910319), StrToTime("2019-07-12 13:45:19"))
}

func TestWeekAround(t *testing.T) {
	monday, sunday := WeekAround(1562909685)

	assert.Equal(t, "2019-07-08", monday)
	assert.Equal(t, "2019-07-14", sunday)
}

func TestIP2Long(t *testing.T) {
	assert.Equal(t, uint32(3221234342), IP2Long("192.0.34.166"))
}

func TestLong2IP(t *testing.T) {
	assert.Equal(t, "192.0.34.166", Long2IP(uint32(3221234342)))
}

func TestAddSlashes(t *testing.T) {
	assert.Equal(t, `Is your name O\'Reilly?`, AddSlashes("Is your name O'Reilly?"))
}

func TestStripSlashes(t *testing.T) {
	assert.Equal(t, "Is your name O'Reilly?", StripSlashes(`Is your name O\'Reilly?`))
}

func TestQuoteMeta(t *testing.T) {
	assert.Equal(t, `Hello world\. \(can you hear me\?\)`, QuoteMeta("Hello world. (can you hear me?)"))
}

func TestOpenFile(t *testing.T) {
	_, err := OpenFile("app.log", os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0775)

	assert.Nil(t, err)
}

func TestVersionCompare(t *testing.T) {
	ok, err := VersionCompare("1.0.0", "1.0.0")
	assert.Nil(t, err)
	assert.True(t, ok)

	ok, err = VersionCompare("1.0.0", "1.0.1")
	assert.Nil(t, err)
	assert.False(t, ok)

	ok, err = VersionCompare("=1.0.0", "1.0.0")
	assert.Nil(t, err)
	assert.True(t, ok)

	ok, err = VersionCompare("=1.0.0", "1.0.1")
	assert.Nil(t, err)
	assert.False(t, ok)

	ok, err = VersionCompare("!=4.0.4", "4.0.0")
	assert.Nil(t, err)
	assert.True(t, ok)

	ok, err = VersionCompare("!=4.0.4", "4.0.4")
	assert.Nil(t, err)
	assert.False(t, ok)

	ok, err = VersionCompare(">2.0.0", "2.0.1")
	assert.Nil(t, err)
	assert.True(t, ok)

	ok, err = VersionCompare(">2.0.0", "1.0.1")
	assert.Nil(t, err)
	assert.False(t, ok)

	ok, err = VersionCompare(">=1.0.0&<2.0.0", "1.0.2")
	assert.Nil(t, err)
	assert.True(t, ok)

	ok, err = VersionCompare(">=1.0.0&<2.0.0", "2.0.1")
	assert.Nil(t, err)
	assert.False(t, ok)

	ok, err = VersionCompare("<2.0.0|>3.0.0", "1.0.2")
	assert.Nil(t, err)
	assert.True(t, ok)

	ok, err = VersionCompare("<2.0.0|>3.0.0", "3.0.1")
	assert.Nil(t, err)
	assert.True(t, ok)

	ok, err = VersionCompare("<2.0.0|>3.0.0", "2.0.1")
	assert.Nil(t, err)
	assert.False(t, ok)
}
