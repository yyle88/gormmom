package gormmom

import (
	"os"

	"github.com/emirpasic/gods/v2/maps/linkedhashmap"
	"github.com/yyle88/formatgo"
	"github.com/yyle88/gormmom/internal/utils"
	"github.com/yyle88/must"
	"github.com/yyle88/rese"
	"github.com/yyle88/tern"
)

type Configs []*Config

func NewConfigs(gormStructs []*GormStruct, options *Options) Configs {
	var configs = make([]*Config, 0, len(gormStructs))
	for _, gormStruct := range gormStructs {
		configs = append(configs, NewConfig(gormStruct, options))
	}
	return configs
}

type GenReplacesResult struct {
	NewCodeResults   []*NewCodeResult
	ChangedLineCount int
	ChangedFileCount int
}

func (R *GenReplacesResult) HasChange() bool {
	return R.ChangedLineCount > 0 || R.ChangedFileCount > 0
}

type NewCodeResult struct {
	NewCode          []byte
	SrcPath          string
	ChangedLineCount int
}

func (R *NewCodeResult) HasChange() bool {
	return R.ChangedLineCount > 0
}

func (configs Configs) GenReplaces() *GenReplacesResult {
	hashMap := linkedhashmap.New[string, *NewCodeResult]()
	for _, config := range configs {
		srcPath := config.gormStruct.sourcePath

		previous, exist := hashMap.Get(srcPath)
		if exist {
			must.Full(previous)
			must.Same(previous.SrcPath, srcPath)
		}

		sourceCode := tern.BFF(exist, func() []byte {
			return previous.NewCode
		}, func() []byte {
			return rese.A1(os.ReadFile(srcPath))
		})

		newCode := config.makeNewCode(sourceCode)
		must.Same(newCode.SrcPath, srcPath)
		if exist {
			must.Same(previous.SrcPath, newCode.SrcPath)
			previous.NewCode = newCode.NewCode
			previous.ChangedLineCount += newCode.ChangedLineCount
		} else {
			hashMap.Put(newCode.SrcPath, newCode)
		}
	}
	changedFileCount := 0
	changedLineCount := 0
	hashMap.Each(func(srcPath string, newCode *NewCodeResult) {
		must.Same(newCode.SrcPath, srcPath)
		if newCode.HasChange() {
			srcCode := must.Have(rese.A1(formatgo.FormatBytes(newCode.NewCode)))
			utils.WriteFile(srcPath, srcCode)
			newCode.NewCode = srcCode
			changedLineCount += newCode.ChangedLineCount
			changedFileCount++
		}
	})
	return &GenReplacesResult{
		NewCodeResults:   hashMap.Values(),
		ChangedLineCount: changedLineCount,
		ChangedFileCount: changedFileCount,
	}
}
