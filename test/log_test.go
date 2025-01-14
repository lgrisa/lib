package test

import (
	"github.com/lgrisa/lib/utils/logutil"
	. "github.com/onsi/gomega"
	"testing"
)

func TestLog(t *testing.T) {
	RegisterTestingT(t)

	logutil.InitLog(-1)

	//Ω

	logutil.LogDebugF("test")
}
