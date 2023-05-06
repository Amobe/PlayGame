package database

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Amobe/PlayGame/server/internal/domain/vo"
	"github.com/Amobe/PlayGame/server/internal/utils"
)

func TestMinionsMarshaller(t *testing.T) {
	origin, err := vo.NewCamp([]vo.Character{
		vo.NewCharacter(1, vo.NewAttributeMap(vo.NewAttribute(vo.AttributeTypeHP, decimal.NewFromInt(100)))),
		vo.NewCharacter(2, vo.NewAttributeMap(vo.NewAttribute(vo.AttributeTypeAGI, decimal.NewFromInt(100)))),
		vo.NewCharacter(3, vo.NewAttributeMap(vo.NewAttribute(vo.AttributeTypeATK, decimal.NewFromInt(100)))),
		vo.NewCharacter(4, vo.NewAttributeMap(vo.NewAttribute(vo.AttributeTypeDEF, decimal.NewFromInt(100)))),
		vo.NewCharacter(5, vo.NewAttributeMap(vo.NewAttribute(vo.AttributeTypeDodge, decimal.NewFromInt(100)))),
		vo.NewCharacter(6, vo.NewAttributeMap(vo.NewAttribute(vo.AttributeTypeHit, decimal.NewFromInt(100)))),
	})
	require.NoError(t, err)
	got := vo.Camp{}
	want, err := vo.NewCamp([]vo.Character{
		vo.NewCharacter(1, vo.NewAttributeMap(vo.NewAttribute(vo.AttributeTypeHP, decimal.NewFromInt(100)))),
		vo.NewCharacter(2, vo.NewAttributeMap(vo.NewAttribute(vo.AttributeTypeAGI, decimal.NewFromInt(100)))),
		vo.NewCharacter(3, vo.NewAttributeMap(vo.NewAttribute(vo.AttributeTypeATK, decimal.NewFromInt(100)))),
		vo.NewCharacter(4, vo.NewAttributeMap(vo.NewAttribute(vo.AttributeTypeDEF, decimal.NewFromInt(100)))),
		vo.NewCharacter(5, vo.NewAttributeMap(vo.NewAttribute(vo.AttributeTypeDodge, decimal.NewFromInt(100)))),
		vo.NewCharacter(6, vo.NewAttributeMap(vo.NewAttribute(vo.AttributeTypeHit, decimal.NewFromInt(100)))),
	})
	require.NoError(t, err)
	marshalled, err := utils.MarshalToJSON(origin)
	require.NoError(t, err)
	err = utils.UnmarshalFromJSON(marshalled, &got)
	require.NoError(t, err)

	assert.Equal(t, want, got)
}
