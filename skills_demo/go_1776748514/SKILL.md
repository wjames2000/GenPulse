# Go

**版本**: 1.0.0

**描述**: 创建Go项目目录结构并初始化git仓库

**作者**: system

**创建时间**: 2026-04-21 13:15:14

**最后更新**: 2026-04-21 13:15:14

## 基本信息

- **分类**: filesystem
- **复杂度**: medium
- **标签**: file_operation, medium, fs_write, git_init, git_add, fs_create, BackendDevAgent, DevOpsAgent, success
- **使用次数**: 0
- **成功率**: 0.0%
- **状态**: 已启用（未验证）

## 执行步骤

### 步骤 1: 创建项目目录结构

**工具**: fs_create

**参数**:

- `path` (string) (必填): 参数 path
- `type` (string) (必填): 参数 type

**预期结果**: 操作成功完成

### 步骤 2: 创建go.mod文件

**工具**: fs_write

**参数**:

- `path` (string) (必填): 参数 path
- `content` (string) (必填): 参数 content

**执行条件**:

- previous_step_success: step_1_success

**预期结果**: 操作成功完成

### 步骤 3: 初始化git仓库

**工具**: git_init

**参数**:

- `path` (string) (必填): 参数 path

**执行条件**:

- previous_step_success: step_2_success

**预期结果**: 操作成功完成


## 使用示例

### 示例 1

Go项目初始化技能|创建Go项目目录结构并初始化git仓库


## 使用技巧

- 确保所有前置条件满足后再执行
- 仔细检查参数格式和类型
- 写入文件前检查目录权限
- 写入文件前检查目录权限

## 注意事项

- ⚠️ 此技能执行时间较长，请确保有足够时间

## 关联工具

- git_add
- fs_create
- fs_write
- git_init

## 适用Agent类型

- BackendDevAgent
- DevOpsAgent

## 性能指标

- **平均执行时间**: 10m0.000004099s
- **Token估算**: 619
