package tests_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yyle88/gormmom/internal/tests"
	"github.com/yyle88/rese"
	"gorm.io/gorm"
)

func TestNewDBRun(t *testing.T) {
	tests.NewDBRun(t, func(db *gorm.DB) {
		var result int
		require.NoError(t, db.Raw("SELECT 1").Scan(&result).Error)
		require.Equal(t, 1, result)
	})
}

func TestNewMemDB(t *testing.T) {
	db := tests.NewMemDB(t)
	defer rese.F0(rese.P1(db.DB()).Close)

	var result int
	require.NoError(t, db.Raw("SELECT 1").Scan(&result).Error)
	require.Equal(t, 1, result)
}
