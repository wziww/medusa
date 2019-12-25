package encrpt

import (
	"github/wziww/medusa/config"
	"testing"
)

func TestMain(m *testing.M) {
	config.Init()
	m.Run()
}
