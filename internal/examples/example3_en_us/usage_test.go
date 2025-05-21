package example3_en_us

import (
	"math/rand/v2"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/done"
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
	db := rese.P1(gorm.Open(sqlite.Open("file::memory:?cache=private"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}))
	defer rese.F0(rese.P1(db.DB()).Close)

	done.Done(db.AutoMigrate(&Example{}))

	caseDB = db
	m.Run()
}

func TestUsage(t *testing.T) {
	examples := make([]*Example, 0, 10)
	for i := 0; i < 10; i++ {
		one := newFakeExample(t)
		t.Log(neatjsons.S(one))
		require.NoError(t, caseDB.Create(one).Error)
		examples = append(examples, one)
	}

	repo := gormrepo.NewRepo(gormclass.Use(&Example{}))

	t.Run("select-first", func(t *testing.T) {
		name := examples[rand.IntN(len(examples))].V名称

		example, err := repo.Repo(caseDB).First(func(db *gorm.DB, cls *ExampleColumns) *gorm.DB {
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

		results, err := repo.Repo(caseDB).Find(func(db *gorm.DB, cls *ExampleColumns) *gorm.DB {
			return db.Where(cls.V名称.In(names))
		})
		require.NoError(t, err)
		t.Log(neatjsons.S(results))
	})
}

func newFakeExample(t *testing.T) *Example {
	a := &Example{}
	require.NoError(t, gofakeit.Struct(a))
	a.ID = 0 // 设置为0以便于使用 gorm 创建数据
	return a
}
