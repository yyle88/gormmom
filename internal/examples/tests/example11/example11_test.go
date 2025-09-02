package example11

import (
	"fmt"
	"math/rand/v2"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/done"
	models2 "github.com/yyle88/gormmom/internal/examples/tests/example11/internal/models"
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

	done.Done(db.AutoMigrate(&models2.Example{}))

	caseDB = db
	m.Run()
}

func TestUsage(t *testing.T) {
	examples := make([]*models2.Example, 0, 10)
	for i := 0; i < 10; i++ {
		one := newFakeExample(t)
		t.Log(neatjsons.S(one))
		require.NoError(t, caseDB.Create(one).Error)
		examples = append(examples, one)
	}

	repo := gormrepo.NewRepo(gormclass.Use(&models2.Example{}))

	t.Run("select-first", func(t *testing.T) {
		name := examples[rand.IntN(len(examples))].V名称

		example, err := repo.Repo(caseDB).First(func(db *gorm.DB, cls *models2.ExampleColumns) *gorm.DB {
			return db.Where(cls.V名称.Eq(name))
		})
		require.NoError(t, err)
		t.Log(neatjsons.S(example))
	})
	t.Run("select-where-in", func(t *testing.T) {
		var names = make([]string, 0, 5)
		for idx, one := range examples {
			if idx%2 == 0 {
				continue
			}
			names = append(names, one.V名称)
		}

		results, err := repo.Repo(caseDB).Find(func(db *gorm.DB, cls *models2.ExampleColumns) *gorm.DB {
			return db.Where(cls.V名称.In(names))
		})
		require.NoError(t, err)
		t.Log(neatjsons.S(results))
	})
}

func newFakeExample(t *testing.T) *models2.Example {
	a := &models2.Example{}
	require.NoError(t, gofakeit.Struct(a))
	a.ID = 0 // 设置为0以便于使用 gorm 创建数据
	return a
}
