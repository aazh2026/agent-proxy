## Why

基于 analyze-prd-gaps-v2 差距分析，发现以下关键差距需要解决：
1. Token 自动刷新逻辑存在但未集成到运行时
2. 管理面板 UI 存在但未实现认证保护
3. 测试覆盖率不足，缺少端到端测试

这些差距影响用户体验、安全性和代码质量，需要在 V1.1 版本中解决。

## What Changes

- 集成 Token 自动刷新功能到 main.go，实现后台 Token 刷新
- 为管理面板 Web UI 添加认证保护
- 补充端到端测试，覆盖 Token 自动刷新和管理面板认证流程
- 完善测试覆盖率，确保核心功能有充分测试

## Capabilities

### New Capabilities
- `token-auto-refresh`: Token 自动刷新功能，支持后台 Token 刷新和生命周期管理
- `admin-panel-auth`: 管理面板认证功能，保护 Web UI 访问
- `integration-tests`: 端到端集成测试，覆盖核心功能流程

### Modified Capabilities
（无现有 spec 文件，所有能力均为新建）

## Impact

- **代码变更**: main.go, token/refresh.go, observability/webui.go, 新增测试文件
- **依赖**: 现有认证系统（auth/）、Token 系统（token/）、可观测性系统（observability/）
- **输出产物**: 集成 Token 自动刷新、管理面板认证保护、完整测试覆盖
- **团队价值**: 提升用户体验、增强安全性、提高代码质量