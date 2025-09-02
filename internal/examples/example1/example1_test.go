package example1

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/done"
	"github.com/yyle88/gormmom/internal/examples/example1/internal/models"
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

	done.Done(db.AutoMigrate(&models.T用户{}))

	caseDB = db
	m.Run()
}

func TestChineseUserExample(t *testing.T) {
	// Test Chinese fields using gormrepo enterprise repository pattern
	repo := gormrepo.NewRepo(gormclass.Use(&models.T用户{}))
	ctx := context.Background()

	// Create test data
	users := []*models.T用户{
		{
			U用户名: "张三",
			E邮箱:  "zhang@example.com",
			A年龄:  25,
			D电话:  "13800138000",
			J住所:  "北京市海淀区",
			S状态:  "活跃",
		},
		{
			U用户名: "李四",
			E邮箱:  "li@example.com",
			A年龄:  30,
			D电话:  "13800138001",
			J住所:  "上海市浦东区",
			S状态:  "活跃",
		},
		{
			U用户名: "王五",
			E邮箱:  "wang@example.com",
			A年龄:  28,
			D电话:  "13800138002",
			J住所:  "深圳市南山区",
			S状态:  "非活跃",
		},
	}

	// Batch insert data
	for _, user := range users {
		require.NoError(t, caseDB.Create(user).Error)
	}

	// Test 1: Use First to find single record
	t.Run("First Find", func(t *testing.T) {
		result, err := repo.With(ctx, caseDB).First(func(db *gorm.DB, cls *models.T用户Columns) *gorm.DB {
			return db.Where(cls.U用户名.Eq("张三"))
		})
		require.NoError(t, err)
		require.Equal(t, "zhang@example.com", result.E邮箱)
		require.Equal(t, 25, result.A年龄)
		t.Log("Find single record:", neatjsons.S(result))
	})

	// Test 2: Use Find to get multiple records
	t.Run("Find Multiple", func(t *testing.T) {
		results, err := repo.With(ctx, caseDB).Find(func(db *gorm.DB, cls *models.T用户Columns) *gorm.DB {
			return db.Where(cls.S状态.Eq("活跃")).Order(cls.A年龄.Ob("asc").Ox())
		})
		require.NoError(t, err)
		require.Len(t, results, 2)
		require.Equal(t, "张三", results[0].U用户名) // younger user first
		require.Equal(t, "李四", results[1].U用户名)
		t.Log("Find multiple records:", neatjsons.S(results))
	})

	// Test 3: Use Count to get statistics
	t.Run("Count Records", func(t *testing.T) {
		count, err := repo.With(ctx, caseDB).Count(func(db *gorm.DB, cls *models.T用户Columns) *gorm.DB {
			return db.Where(cls.A年龄.Gte(26))
		})
		require.NoError(t, err)
		require.Equal(t, int64(2), count) // 2 users with age >= 26
		t.Log("Records with age >= 26:", count)
	})

	// Test 4: Complex condition search
	t.Run("Complex Search", func(t *testing.T) {
		results, err := repo.With(ctx, caseDB).Find(func(db *gorm.DB, cls *models.T用户Columns) *gorm.DB {
			return db.Where(cls.A年龄.Between(25, 30)).Where(cls.J住所.Like("%市%"))
		})
		require.NoError(t, err)
		require.Len(t, results, 3) // all users match conditions
		t.Log("Complex search result:", len(results), "records")
	})

	// Test 5: Use Update to modify data
	t.Run("Update Record", func(t *testing.T) {
		// Update Zhang San's age
		err := repo.With(ctx, caseDB).Update(
			func(db *gorm.DB, cls *models.T用户Columns) *gorm.DB {
				return db.Where(cls.U用户名.Eq("张三"))
			},
			func(cls *models.T用户Columns) (string, interface{}) {
				return cls.A年龄.Kv(26)
			},
		)
		require.NoError(t, err)
		t.Log("Successfully updated Zhang San's age")

		// Verify update result
		updated, err := repo.With(ctx, caseDB).First(func(db *gorm.DB, cls *models.T用户Columns) *gorm.DB {
			return db.Where(cls.U用户名.Eq("张三"))
		})
		require.NoError(t, err)
		require.Equal(t, 26, updated.A年龄)
		t.Log("Updated Zhang San:", neatjsons.S(updated))
	})

	// Test 6: Use Updates for batch modification
	t.Run("Updates Batch", func(t *testing.T) {
		// Batch update inactive record information
		err := repo.With(ctx, caseDB).Updates(
			func(db *gorm.DB, cls *models.T用户Columns) *gorm.DB {
				return db.Where(cls.S状态.Eq("非活跃"))
			},
			func(cls *models.T用户Columns) map[string]interface{} {
				return cls.
					Kw(cls.A年龄.Kv(32)).
					Kw(cls.J住所.Kv("广州市天河区")).
					AsMap()
			},
		)
		require.NoError(t, err)
		t.Log("Successfully batch updated inactive record information")
	})

	// Test 7: Use Exist to check record existence
	t.Run("Exist Check", func(t *testing.T) {
		// Check if records with age > 35 exist
		exists, err := repo.With(ctx, caseDB).Exist(func(db *gorm.DB, cls *models.T用户Columns) *gorm.DB {
			return db.Where(cls.A年龄.Gt(35))
		})
		require.NoError(t, err)
		require.False(t, exists) // should not exist
		t.Log("Records with age > 35 exist:", exists)

		// Check if records with age = 26 exist
		exists, err = repo.With(ctx, caseDB).Exist(func(db *gorm.DB, cls *models.T用户Columns) *gorm.DB {
			return db.Where(cls.A年龄.Eq(26))
		})
		require.NoError(t, err)
		require.True(t, exists) // should exist (Zhang San)
		t.Log("Records with age = 26 exist:", exists)
	})

	// Test 8: Advanced enterprise search combinations
	t.Run("Advanced Enterprise Search", func(t *testing.T) {
		// Use multiple conditions and IN filters
		results, err := repo.With(ctx, caseDB).Find(func(db *gorm.DB, cls *models.T用户Columns) *gorm.DB {
			return db.Where(cls.U用户名.In([]string{"张三", "李四", "王五"})).
				Where(cls.A年龄.Gte(25)).
				Where(cls.E邮箱.NotLike("%不存在%")).
				Order(cls.A年龄.Ob("desc").Ox()).
				Order(cls.U用户名.Ob("asc").Ox())
		})
		require.NoError(t, err)
		require.Len(t, results, 3)
		t.Log("Advanced enterprise search result:", len(results), "records")

		// Verify sorting result: by age desc, then by name asc
		ages := make([]int, len(results))
		for i, user := range results {
			ages[i] = user.A年龄
		}
		t.Log("Age sorting result:", ages)
	})
}
