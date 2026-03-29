## ADDED Requirements

### Requirement: 健康检查测试
系统 SHALL 提供健康检查测试，验证 /health、/health/ready、/health/live 端点。

#### Scenario: 健康检查端点测试
- **WHEN** 运行健康检查端点测试时
- **THEN** 系统验证 /health 端点返回正确状态

#### Scenario: 就绪检查端点测试
- **WHEN** 运行就绪检查端点测试时
- **THEN** 系统验证 /health/ready 端点检查数据库连接

#### Scenario: 存活检查端点测试
- **WHEN** 运行存活检查端点测试时
- **THEN** 系统验证 /health/live 端点返回存活状态

### Requirement: 指标收集测试
系统 SHALL 提供指标收集测试，验证 /metrics 端点和指标格式。

#### Scenario: 指标端点测试
- **WHEN** 运行指标端点测试时
- **THEN** 系统验证 /metrics 端点返回 Prometheus 格式指标

#### Scenario: 指标格式测试
- **WHEN** 运行指标格式测试时
- **THEN** 系统验证指标包含请求数、延迟、成功率等关键指标

### Requirement: 日志记录测试
系统 SHALL 提供日志记录测试，验证日志收集和查询功能。

#### Scenario: 日志端点测试
- **WHEN** 运行日志端点测试时
- **THEN** 系统验证 /logs 端点返回请求日志

#### Scenario: 日志脱敏测试
- **WHEN** 运行日志脱敏测试时
- **THEN** 系统验证 Token 在日志中只显示前4位和后4位

### Requirement: Web UI 测试
系统 SHALL 提供 Web UI 测试，验证管理面板功能和认证保护。

#### Scenario: Web UI 访问测试
- **WHEN** 运行 Web UI 访问测试时
- **THEN** 系统验证管理面板需要认证才能访问

#### Scenario: Web UI 功能测试
- **WHEN** 运行 Web UI 功能测试时
- **THEN** 系统验证管理面板显示实时指标、日志、Token 管理功能