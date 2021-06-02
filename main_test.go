package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_extractBvidFromUrl(t *testing.T) {
	bvid, err := extractBvidFromUrl("https://www.bilibili.com/video/BV1Yi4y1P7Bd")
	assert.Nil(t, err)
	assert.Equal(t, bvid, "BV1Yi4y1P7Bd")

	bvid, err = extractBvidFromUrl("https://www.bilibili.com/video/BV1Yi4y1P7Bd?spm_id_from=333.851.b_62696c695f7265706f72745f646f756761.30")
	assert.Nil(t, err)
	assert.Equal(t, bvid, "BV1Yi4y1P7Bd")
}
