package test

import (
	"testing"

	"github.com/neutrinoapp/neutrino/test/common"
	"github.com/stretchr/testify/assert"
)

func TestApiRegisterCreateAppLogin(t *testing.T) {
	assert.NotEqual(t, common.GetClient().Token, "")
}
