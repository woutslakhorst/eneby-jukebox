package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseRFIDBytes(t *testing.T) {
	rfidResponse := "\x02\x33\x41\x30\x30\x31\x39\x45\x41\x32\x45\x45\x37\x03"
	rfidBytes := []byte(rfidResponse)

	t.Run("headByte", func(t *testing.T) {
		assert.True(t, isHeadByte(string(rfidBytes[0])))
	})

	t.Run("tailByte", func(t *testing.T) {
		assert.True(t, isTailByte(string(rfidBytes[13])))
	})

	t.Run("checksum", func(t *testing.T) {
		assert.True(t, checksum(rfidBytes[11:13], rfidBytes[1:11]))
	})

	t.Run("example", func(t *testing.T) {
		n, err := parseRFIDBytes(rfidBytes)
		if assert.NoError(t, err) {
			assert.Equal(t, uint32(1698350), n)
		}
	})
}
