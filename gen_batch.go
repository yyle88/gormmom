package gormmom

import (
	"os"

	"github.com/emirpasic/gods/v2/maps/linkedhashmap"
	"github.com/yyle88/formatgo"
	"github.com/yyle88/must"
	"github.com/yyle88/rese"
	"github.com/yyle88/tern"
)

// Configs represents a collection of configuration instances for batch processing
// Enables processing multiple GORM structures in a single operation
// Provides batch generation and file modification capabilities
//
// Configs 代表一集配置实例，用于批量处理
// 允许在单个操作中处理多个 GORM 结构
// 提供批量生成和文件修改功能
type Configs []*Config

// NewConfigs creates a batch configuration from GORM structures and options
// Processes multiple GORM structures with consistent options settings
// Returns configured batch prepared for tag generation and file modification
//
// NewConfigs 从 GORM 结构和选项创建批量配置
// 使用一致的选项设置处理多个 GORM 结构
// 返回配置好的批量操作，准备进行标签生成和文件修改
func NewConfigs(structs []*GormStruct, options *Options) Configs {
	var configs = make([]*Config, 0, len(structs))
	for _, structI := range structs {
		configs = append(configs, NewConfig(structI, options))
	}
	return configs
}

// CodeResults contains the result of batch code generation and replacement
// Provides statistics about changed lines and files during processing
// Includes detailed results for each processed file with change tracking
//
// CodeResults 包含批量代码生成和替换的结果
// 提供处理过程中更改的行数和文件统计
// 包含每个处理文件的详细结果和更改跟踪
type CodeResults struct {
	Items            []*CodeResult // Detailed results for each file // 每个文件的详细结果
	ChangedLineCount int           // Total number of changed lines // 更改的总行数
	ChangedFileCount int           // Total number of changed files // 更改的文件总数
}

// HasChange checks if any changes were made during batch processing
// Returns true if any lines or files were modified
//
// HasChange 检查批量处理过程中是否有任何更改
// 如果任何行或文件被修改则返回 true
func (R *CodeResults) HasChange() bool {
	return R.ChangedLineCount > 0 || R.ChangedFileCount > 0
}

// CodeResult contains the result of code generation for a single file
// Includes the generated code, source path, and change statistics
// Used for tracking modifications during batch processing operations
//
// CodeResult 包含单个文件的代码生成结果
// 包括生成的代码、源路径和更改统计
// 用于在批量处理操作期间跟踪修改
type CodeResult struct {
	OutputCode       []byte // Generated code content // 生成的代码内容
	SourcePath       string // Source file path // 源文件路径
	ChangedLineCount int    // Number of lines changed // 更改的行数
}

// HasChange checks if this file result contains any changes
// Returns true if any lines were modified in this file
//
// HasChange 检查此文件结果是否包含任何更改
// 如果此文件中有任何行被修改则返回 true
func (R *CodeResult) HasChange() bool {
	return R.ChangedLineCount > 0
}

// Generate performs batch code generation and file replacement operations
// Processes all configurations and applies changes to source files with formatting
// Returns comprehensive results with statistics about modifications made
//
// Generate 执行批量代码生成和文件替换操作
// 处理所有配置并将更改带格式化地应用到源文件
// 返回包含修改统计的全面结果
func (configs Configs) Generate() *CodeResults {
	results := configs.Preview()
	for _, item := range results.Items {
		if item.HasChange() {
			srcCode := must.Have(rese.A1(formatgo.FormatBytes(item.OutputCode)))
			must.Done(os.WriteFile(item.SourcePath, srcCode, 0644))
			item.OutputCode = srcCode
		}
	}
	return results
}

// Preview performs batch code generation without modifying source files
// Processes all configurations and returns generated code without writing to disk
// Returns comprehensive results with statistics about potential modifications
//
// Preview 执行批量代码生成但不修改源文件
// 处理所有配置并返回生成的代码而不写入磁盘
// 返回包含潜在修改统计的全面结果
func (configs Configs) Preview() *CodeResults {
	hashMap := linkedhashmap.New[string, *CodeResult]()
	for _, config := range configs {
		srcPath := config.structI.sourcePath

		previous, exist := hashMap.Get(srcPath)
		if exist {
			must.Full(previous)
			must.Same(previous.SourcePath, srcPath)
		}

		sourceCode := tern.BFF(exist, func() []byte {
			return previous.OutputCode
		}, func() []byte {
			return rese.A1(os.ReadFile(srcPath))
		})

		newCode := config.makeNewCode(sourceCode)
		must.Same(newCode.SourcePath, srcPath)
		if exist {
			must.Same(previous.SourcePath, newCode.SourcePath)
			previous.OutputCode = newCode.OutputCode
			previous.ChangedLineCount += newCode.ChangedLineCount
		} else {
			hashMap.Put(newCode.SourcePath, newCode)
		}
	}
	changedFileCount := 0
	changedLineCount := 0
	hashMap.Each(func(srcPath string, newCode *CodeResult) {
		must.Same(newCode.SourcePath, srcPath)
		if newCode.HasChange() {
			changedLineCount += newCode.ChangedLineCount
			changedFileCount++
		}
	})
	return &CodeResults{
		Items:            hashMap.Values(),
		ChangedLineCount: changedLineCount,
		ChangedFileCount: changedFileCount,
	}
}
