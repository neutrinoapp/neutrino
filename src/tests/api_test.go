package tests

import (
	"testing"

	"github.com/neutrinoapp/neutrino/src/tests/common"
	"github.com/stretchr/testify/assert"
)

func TestApiRegisterCreateAppLogin(t *testing.T) {
	assert.NotEqual(t, common.GetClient().Token, "")
}
