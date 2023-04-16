package database

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Amobe/PlayGame/server/internal/domain/battle"
	"github.com/Amobe/PlayGame/server/internal/domain/vo"
	"github.com/Amobe/PlayGame/server/internal/utils"
)

// fix marshaller error
func TestMinionsMarshaller(t *testing.T) {
	origin := battle.NewMinions(true, []vo.Character{
		vo.NewCharacter("c1", vo.NewAttributeMap(vo.NewAttribute(vo.AttributeTypeHP, decimal.NewFromInt(100)))),
		vo.NewCharacter("c2", vo.NewAttributeMap(vo.NewAttribute(vo.AttributeTypeAGI, decimal.NewFromInt(100)))),
		vo.NewCharacter("c3", vo.NewAttributeMap(vo.NewAttribute(vo.AttributeTypeATK, decimal.NewFromInt(100)))),
		vo.NewCharacter("c4", vo.NewAttributeMap(vo.NewAttribute(vo.AttributeTypeDEF, decimal.NewFromInt(100)))),
		vo.NewCharacter("c5", vo.NewAttributeMap(vo.NewAttribute(vo.AttributeTypeDodge, decimal.NewFromInt(100)))),
		vo.NewCharacter("c6", vo.NewAttributeMap(vo.NewAttribute(vo.AttributeTypeHit, decimal.NewFromInt(100)))),
	})
	got := &battle.Minions{}
	want := battle.NewMinions(true, []vo.Character{
		vo.NewCharacter("c1", vo.NewAttributeMap(vo.NewAttribute(vo.AttributeTypeHP, decimal.NewFromInt(100)))),
		vo.NewCharacter("c2", vo.NewAttributeMap(vo.NewAttribute(vo.AttributeTypeAGI, decimal.NewFromInt(100)))),
		vo.NewCharacter("c3", vo.NewAttributeMap(vo.NewAttribute(vo.AttributeTypeATK, decimal.NewFromInt(100)))),
		vo.NewCharacter("c4", vo.NewAttributeMap(vo.NewAttribute(vo.AttributeTypeDEF, decimal.NewFromInt(100)))),
		vo.NewCharacter("c5", vo.NewAttributeMap(vo.NewAttribute(vo.AttributeTypeDodge, decimal.NewFromInt(100)))),
		vo.NewCharacter("c6", vo.NewAttributeMap(vo.NewAttribute(vo.AttributeTypeHit, decimal.NewFromInt(100)))),
	})
	marshalled, err := utils.MarshalToJSON(origin)
	require.NoError(t, err)
	err = utils.UnmarshalFromJSON(marshalled, got)
	require.NoError(t, err)

	assert.Equal(t, want, got)
}
