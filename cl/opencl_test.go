package cl_test

import (
	"testing"

	"github.com/gohxs/vec-benchmark/cl"
)

func TestInitial(t *testing.T) {
	t.Log("Platforms:", cl.PlatformIDCount())
}

func TestImpl(t *testing.T) {
	//cl.CTest()
}
