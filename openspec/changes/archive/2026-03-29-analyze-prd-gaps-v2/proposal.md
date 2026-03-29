## Why

基于前期分析（已归档的 analyze-prd-gaps 变更），需要持续跟踪 PRD 需求与代码实现的差距，确保产品迭代方向正确。随着代码库演进和需求变更，差距分析需要定期更新，为团队提供最新的实现状态和优先级指导。

## What Changes

- 重新评估当前代码库与最新 PRD 需求的符合度
- 识别新增或变更的功能需求
- 评估前期分析中发现的差距是否已解决
- 更新优先级和实现建议

## Capabilities

### New Capabilities
- `prd-gap-analysis-v2`: 基于最新代码库的全面 PRD 差距分析
- `implementation-status-tracking`: 跟踪前期发现的差距解决进度
- `incremental-gap-identification`: 识别新增或变更的需求

### Modified Capabilities
（无现有 spec 文件，所有能力均为新建）

## Impact

- **代码分析范围**: 所有 internal/ 子包，重点关注前期发现的差距领域
- **文档依赖**: docs/PRD.md 作为需求基准
- **前期参考**: 已归档的 analyze-prd-gaps 变更结果
- **输出产物**: 更新后的差距分析报告、实现状态矩阵、优先级清单
- **团队价值**: 为持续迭代提供最新的功能缺口和开发重点