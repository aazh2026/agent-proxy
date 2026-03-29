## Why

需要系统性地识别当前 agent-proxy 代码实现与 PRD 需求之间的差距，为后续迭代提供清晰的功能缺失清单和优先级指导。当前代码已实现基础功能，但缺乏与 PRD 的正式对标分析，无法准确评估完成度和规划下一步工作。

## What Changes

- 创建 PRD 需求与代码实现的映射矩阵
- 识别已实现、部分实现、未实现的功能点
- 标注实现质量与 PRD 要求的差距（性能、安全、功能完整性）
- 生成可执行的差距报告，包含优先级和建议

## Capabilities

### New Capabilities
- `prd-gap-analysis`: 全面分析 PRD 需求与当前代码实现的差距
- `implementation-mapping`: 建立需求到代码的映射关系
- `gap-prioritization`: 基于 PRD 优先级对差距进行排序

### Modified Capabilities
（无现有 spec 文件，所有能力均为新建）

## Impact

- **代码分析范围**: 所有 internal/ 子包，重点关注 auth、token、routing、api、provider 模块
- **文档依赖**: docs/PRD.md 作为需求基准
- **输出产物**: 差距分析报告、实现状态矩阵、优先级清单
- **团队价值**: 为 V1.1/V1.2 迭代提供明确的功能缺口和开发重点