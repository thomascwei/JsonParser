package internal

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGenerateHeadAndStructs(t *testing.T) {
	_, err := GenerateAllParseFuncString()
	require.NoError(t, err)

}
