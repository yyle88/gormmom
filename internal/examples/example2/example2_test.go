package example2

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/done"
	"github.com/yyle88/gormmom/internal/examples/example2/internal/models"
	"github.com/yyle88/gormrepo"
	"github.com/yyle88/gormrepo/gormclass"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/rese"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var caseDB *gorm.DB

func TestMain(m *testing.M) {
	dsn := fmt.Sprintf("file:db-%s?mode=memory&cache=shared", uuid.New().String())
	db := rese.P1(gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}))
	defer rese.F0(rese.P1(db.DB()).Close)

	done.Done(db.AutoMigrate(&models.Tユーザー{}))

	caseDB = db
	m.Run()
}

func TestJapaneseUserExample(t *testing.T) {
	// Test Japanese fields using gormrepo enterprise repository pattern
	repo := gormrepo.NewRepo(gormclass.Use(&models.Tユーザー{}))
	ctx := context.Background()

	// Create Japanese test data
	users := []*models.Tユーザー{
		{
			Uユーザー名: "田中太郎",
			Eメール:   "tanaka@example.jp",
			A年齢:    25,
			D電話:    "090-1234-5678",
			J住所:    "東京都渋谷区",
			Sステータス: "アクティブ",
		},
		{
			Uユーザー名: "佐藤花子",
			Eメール:   "sato@example.jp",
			A年齢:    30,
			D電話:    "090-2345-6789",
			J住所:    "大阪市中央区",
			Sステータス: "アクティブ",
		},
		{
			Uユーザー名: "鈴木一郎",
			Eメール:   "suzuki@example.jp",
			A年齢:    28,
			D電話:    "090-3456-7890",
			J住所:    "名古屋市中区",
			Sステータス: "非アクティブ",
		},
	}

	// Batch insert Japanese test data
	for _, user := range users {
		require.NoError(t, caseDB.Create(user).Error)
	}

	// Test 1: Find single record using Japanese fields
	t.Run("Japanese fields single find", func(t *testing.T) {
		result, err := repo.With(ctx, caseDB).First(func(db *gorm.DB, cls *models.TユーザーColumns) *gorm.DB {
			return db.Where(cls.Uユーザー名.Eq("田中太郎"))
		})
		require.NoError(t, err)
		require.Equal(t, "tanaka@example.jp", result.Eメール)
		require.Equal(t, 25, result.A年齢)
		t.Log("Japanese record find result:", neatjsons.S(result))
	})

	// Test 2: Find multiple records using Japanese fields
	t.Run("Japanese fields multiple find", func(t *testing.T) {
		results, err := repo.With(ctx, caseDB).Find(func(db *gorm.DB, cls *models.TユーザーColumns) *gorm.DB {
			return db.Where(cls.Sステータス.Eq("アクティブ")).Order(cls.A年齢.Ob("asc").Ox())
		})
		require.NoError(t, err)
		require.Len(t, results, 2)
		require.Equal(t, "田中太郎", results[0].Uユーザー名) // younger user first
		require.Equal(t, "佐藤花子", results[1].Uユーザー名)
		t.Log("Active records find result:", len(results), "items")
	})

	// Test 3: Count statistics using Japanese fields
	t.Run("Japanese fields count stats", func(t *testing.T) {
		count, err := repo.With(ctx, caseDB).Count(func(db *gorm.DB, cls *models.TユーザーColumns) *gorm.DB {
			return db.Where(cls.A年齢.Gte(26))
		})
		require.NoError(t, err)
		require.Equal(t, int64(2), count) // 2 users with age >= 26
		t.Log("Records with age > 26:", count)
	})

	// Test 4: Complex condition search using Japanese fields
	t.Run("Japanese fields complex search", func(t *testing.T) {
		results, err := repo.With(ctx, caseDB).Find(func(db *gorm.DB, cls *models.TユーザーColumns) *gorm.DB {
			return db.Where(cls.J住所.Like("%区%")).
				Where(cls.A年齢.Between(25, 30)).
				Where(cls.Eメール.Like("%@example.jp"))
		})
		require.NoError(t, err)
		require.Len(t, results, 3) // all users match conditions
		t.Log("Complex search result:", len(results), "items")
	})

	// Test 5: Update operation using Japanese fields
	t.Run("Japanese fields update operation", func(t *testing.T) {
		// Update Tanaka Taro's age
		err := repo.With(ctx, caseDB).Update(
			func(db *gorm.DB, cls *models.TユーザーColumns) *gorm.DB {
				return db.Where(cls.Uユーザー名.Eq("田中太郎"))
			},
			func(cls *models.TユーザーColumns) (string, interface{}) {
				return cls.A年齢.Kv(26)
			},
		)
		require.NoError(t, err)
		t.Log("Successfully updated Tanaka Taro's age")

		// Verify update result
		updated, err := repo.With(ctx, caseDB).First(func(db *gorm.DB, cls *models.TユーザーColumns) *gorm.DB {
			return db.Where(cls.Uユーザー名.Eq("田中太郎"))
		})
		require.NoError(t, err)
		require.Equal(t, 26, updated.A年齢)
		t.Log("Updated Tanaka Taro info:", neatjsons.S(updated))
	})

	// Test 6: Batch update using Japanese fields
	t.Run("Japanese fields batch update", func(t *testing.T) {
		// Batch update inactive record status
		err := repo.With(ctx, caseDB).Updates(
			func(db *gorm.DB, cls *models.TユーザーColumns) *gorm.DB {
				return db.Where(cls.Sステータス.Eq("非アクティブ"))
			},
			func(cls *models.TユーザーColumns) map[string]interface{} {
				return cls.
					Kw(cls.Sステータス.Kv("休眠")).
					Kw(cls.J住所.Kv("未確認")).
					AsMap()
			},
		)
		require.NoError(t, err)
		t.Log("Successfully batch updated inactive record status")
	})

	// Test 7: Existence check using Japanese fields
	t.Run("Japanese fields existence check", func(t *testing.T) {
		// Check if records with age > 35 exist
		exists, err := repo.With(ctx, caseDB).Exist(func(db *gorm.DB, cls *models.TユーザーColumns) *gorm.DB {
			return db.Where(cls.A年齢.Gt(35))
		})
		require.NoError(t, err)
		require.False(t, exists) // should not exist
		t.Log("Records with age > 35 exist:", exists)

		// Check if Tokyo records exist
		exists, err = repo.With(ctx, caseDB).Exist(func(db *gorm.DB, cls *models.TユーザーColumns) *gorm.DB {
			return db.Where(cls.J住所.Like("東京%"))
		})
		require.NoError(t, err)
		require.True(t, exists) // should exist
		t.Log("Tokyo records exist:", exists)
	})

	// Test 8: Validate documentation examples
	t.Run("Validate documentation Japanese examples", func(t *testing.T) {
		// Verify Japanese field names can correctly map to database columns
		cols := (&models.Tユーザー{}).Columns()

		t.Log("Japanese field mapping validation:")
		t.Logf("Uユーザー名 -> %s", cols.Uユーザー名.Name())
		t.Logf("Eメール -> %s", cols.Eメール.Name())
		t.Logf("A年齢 -> %s", cols.A年齢.Name())
		t.Logf("Sステータス -> %s", cols.Sステータス.Name())

		// Ensure all fields have correct column name mapping
		require.NotEmpty(t, cols.Uユーザー名.Name())
		require.NotEmpty(t, cols.Eメール.Name())
		require.NotEmpty(t, cols.A年齢.Name())
		require.NotEmpty(t, cols.Sステータス.Name())

		t.Log("✅ Documentation Japanese examples validated successfully!")
	})
}
