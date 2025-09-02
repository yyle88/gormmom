package example3

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/done"
	"github.com/yyle88/gormmom/internal/examples/example3/internal/models"
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

	done.Done(db.AutoMigrate(&models.T사용자{}))

	caseDB = db
	m.Run()
}

func TestKoreanUserExample(t *testing.T) {
	// Test Korean fields using gormrepo enterprise repository pattern
	repo := gormrepo.NewRepo(gormclass.Use(&models.T사용자{}))
	ctx := context.Background()

	// Create Korean test data
	users := []*models.T사용자{
		{
			U사용자명: "김민수",
			E이메일:  "kim@example.kr",
			N나이:   27,
			J전화:   "010-1234-5678",
			J주소:   "서울특별시 강남구",
			S상태:   "활성",
		},
		{
			U사용자명: "박지은",
			E이메일:  "park@example.kr",
			N나이:   32,
			J전화:   "010-2345-6789",
			J주소:   "부산광역시 해운대구",
			S상태:   "활성",
		},
		{
			U사용자명: "이철호",
			E이메일:  "lee@example.kr",
			N나이:   29,
			J전화:   "010-3456-7890",
			J주소:   "대구광역시 중구",
			S상태:   "비활성",
		},
	}

	// Batch insert Korean test data
	for _, user := range users {
		require.NoError(t, caseDB.Create(user).Error)
	}

	// Test 1: Find single record using Korean fields
	t.Run("Korean fields single find", func(t *testing.T) {
		result, err := repo.With(ctx, caseDB).First(func(db *gorm.DB, cls *models.T사용자Columns) *gorm.DB {
			return db.Where(cls.U사용자명.Eq("김민수"))
		})
		require.NoError(t, err)
		require.Equal(t, "kim@example.kr", result.E이메일)
		require.Equal(t, 27, result.N나이)
		t.Log("Korean record find result:", neatjsons.S(result))
	})

	// Test 2: Find multiple records using Korean fields
	t.Run("Korean fields multiple find", func(t *testing.T) {
		results, err := repo.With(ctx, caseDB).Find(func(db *gorm.DB, cls *models.T사용자Columns) *gorm.DB {
			return db.Where(cls.S상태.Eq("활성")).Order(cls.N나이.Ob("asc").Ox())
		})
		require.NoError(t, err)
		require.Len(t, results, 2)
		require.Equal(t, "김민수", results[0].U사용자명) // younger user first
		require.Equal(t, "박지은", results[1].U사용자명)
		t.Log("Active records find result:", len(results), "items")
	})

	// Test 3: Count statistics using Korean fields
	t.Run("Korean fields count stats", func(t *testing.T) {
		count, err := repo.With(ctx, caseDB).Count(func(db *gorm.DB, cls *models.T사용자Columns) *gorm.DB {
			return db.Where(cls.N나이.Gte(28))
		})
		require.NoError(t, err)
		require.Equal(t, int64(2), count) // 2 users with age >= 28
		t.Log("Records with age > 28:", count)
	})

	// Test 4: Complex condition search using Korean fields
	t.Run("Korean fields complex search", func(t *testing.T) {
		results, err := repo.With(ctx, caseDB).Find(func(db *gorm.DB, cls *models.T사용자Columns) *gorm.DB {
			return db.Where(cls.J주소.Like("%시%")).
				Where(cls.N나이.Between(25, 35)).
				Where(cls.E이메일.Like("%@example.kr"))
		})
		require.NoError(t, err)
		require.Len(t, results, 3) // all users match conditions
		t.Log("Complex condition search result:", len(results), "items")
	})

	// Test 5: Update operation using Korean fields
	t.Run("Korean fields update operation", func(t *testing.T) {
		// Update Kim Min-su's age
		err := repo.With(ctx, caseDB).Update(
			func(db *gorm.DB, cls *models.T사용자Columns) *gorm.DB {
				return db.Where(cls.U사용자명.Eq("김민수"))
			},
			func(cls *models.T사용자Columns) (string, interface{}) {
				return cls.N나이.Kv(28)
			},
		)
		require.NoError(t, err)
		t.Log("Successfully updated Kim Min-su's age")

		// Verify update result
		updated, err := repo.With(ctx, caseDB).First(func(db *gorm.DB, cls *models.T사용자Columns) *gorm.DB {
			return db.Where(cls.U사용자명.Eq("김민수"))
		})
		require.NoError(t, err)
		require.Equal(t, 28, updated.N나이)
		t.Log("Updated Kim Min-su info:", neatjsons.S(updated))
	})

	// Test 6: Batch update using Korean fields
	t.Run("Korean fields batch update", func(t *testing.T) {
		// Batch update inactive record status
		err := repo.With(ctx, caseDB).Updates(
			func(db *gorm.DB, cls *models.T사용자Columns) *gorm.DB {
				return db.Where(cls.S상태.Eq("비활성"))
			},
			func(cls *models.T사용자Columns) map[string]interface{} {
				return cls.
					Kw(cls.S상태.Kv("휴면")).
					Kw(cls.J주소.Kv("주소 미확인")).
					AsMap()
			},
		)
		require.NoError(t, err)
		t.Log("Successfully batch updated inactive users")
	})

	// Test 7: Existence check using Korean fields
	t.Run("Korean fields existence check", func(t *testing.T) {
		// Check if records with age > 40 exist
		exists, err := repo.With(ctx, caseDB).Exist(func(db *gorm.DB, cls *models.T사용자Columns) *gorm.DB {
			return db.Where(cls.N나이.Gt(40))
		})
		require.NoError(t, err)
		require.False(t, exists) // should not exist
		t.Log("Records with age > 40 exist:", exists)

		// Check if Seoul records exist
		exists, err = repo.With(ctx, caseDB).Exist(func(db *gorm.DB, cls *models.T사용자Columns) *gorm.DB {
			return db.Where(cls.J주소.Like("서울%"))
		})
		require.NoError(t, err)
		require.True(t, exists) // should exist
		t.Log("Seoul records exist:", exists)
	})

	// Test 8: Validate documentation examples
	t.Run("Validate documentation Korean examples", func(t *testing.T) {
		// Verify Korean field names can correctly map to database columns
		cols := (&models.T사용자{}).Columns()

		t.Log("Korean field mapping validation:")
		t.Logf("U사용자명 -> %s", cols.U사용자명.Name())
		t.Logf("E이메일 -> %s", cols.E이메일.Name())
		t.Logf("N나이 -> %s", cols.N나이.Name())
		t.Logf("S상태 -> %s", cols.S상태.Name())

		// Ensure all fields have correct column name mapping
		require.NotEmpty(t, cols.U사용자명.Name())
		require.NotEmpty(t, cols.E이메일.Name())
		require.NotEmpty(t, cols.N나이.Name())
		require.NotEmpty(t, cols.S상태.Name())

		t.Log("✅ Documentation Korean examples validated successfully!")
	})
}
