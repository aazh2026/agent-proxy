## Why

当前代码库的测试覆盖不完整，缺少多个关键模块的测试，特别是 provider routing、observability 和 configuration 模块。这导致无法验证 PRD 需求的实现质量，也无法保证代码变更不会引入回归问题。需要补充完整的测试覆盖以确保代码质量和可维护性。

## What Changes

- 为 provider routing 模块添加测试，验证模型路由、负载均衡、故障转移逻辑
- 为 observability 模块添加测试，验证健康检查、指标收集、日志记录功能
- 为 configuration 模块添加测试，验证配置加载、热重载、验证逻辑
- 添加端到端集成测试，验证完整请求流程
- 补充性能测试，验证 PRD 性能目标

## Capabilities

### New Capabilities
- `provider-routing-tests`: Provider 路由测试，覆盖模型路由、负载均衡、故障转移
- `observability-tests`: 可观测性测试，覆盖健康检查、指标、日志、Web UI
- `configuration-tests`: 配置管理测试，覆盖配置加载、热重载、验证
- `integration-tests`: 端到端集成测试，覆盖完整请求流程
- `performance-tests`: 性能测试，验证延迟、QPS、内存占用

### Modified Capabilities
（无现有 spec 文件，所有能力均为新建）

## Impact

- **代码变更**: 新增测试文件在 internal/provider、internal/observability、internal/config 目录
- **依赖**: 现有测试框架（testing、testify）
- **输出产物**: 完整测试覆盖，包括单元测试、集成测试、性能测试
- **团队价值**: 提高代码质量、降低回归风险、确保 PRD 需求实现正确性