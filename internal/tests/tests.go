package tests

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/yyle88/must"
	"github.com/yyle88/rese"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewDBRun runs a function with an in-mem SQLite database
// Auto creates a temporary database connection and handles cleanup
// NewDBRun 在内存数据库中运行函数，用于测试目的
// 自动创建临时数据库连接并在函数执行后处理清理工作
func NewDBRun(t *testing.T, run func(db *gorm.DB)) {
	dsn := fmt.Sprintf("file:db-%s?mode=memory&cache=shared", uuid.New().String())
	db := rese.P1(gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}))
	defer rese.F0(rese.P1(db.DB()).Close)

	t.Log("--- DB BEGIN ---")
	run(db)
	t.Log("--- DB CLOSE ---")
}

// NewMemDB creates an in-mem SQLite database connection with auto cleanup
// Returns the database connection, cleanup is handled via t.Cleanup
// NewMemDB 创建内存 SQLite 数据库连接，自动清理
// 返回数据库连接，通过 t.Cleanup 处理清理工作
func NewMemDB(t *testing.T) *gorm.DB {
	dsn := fmt.Sprintf("file:db-%s?mode=memory&cache=shared", uuid.New().String())
	db := rese.P1(gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}))
	t.Cleanup(func() {
		must.Done(rese.P1(db.DB()).Close())
	})

	t.Log("--- DB BEGIN ---")
	t.Cleanup(func() {
		t.Log("--- DB CLOSE ---")
	})
	return db
}
