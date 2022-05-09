package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContains(t *testing.T) {
	assert.Equal(t, "/%E4%BD%A0%E5%A5%BD%E4%B8%96%E7%95%8C", EncodeUrl("/你好世界"))

	assert.Equal(t, "%2F%E4%BD%A0%E5%A5%BD%E4%B8%96%E7%95%8C", EncodeUrlComponent("/你好世界"))
}
