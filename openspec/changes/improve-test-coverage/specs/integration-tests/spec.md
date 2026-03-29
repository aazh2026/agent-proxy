## ADDED Requirements

### Requirement: 完整请求流程测试
系统 SHALL 提供端到端集成测试，验证从客户端请求到响应的完整流程。

#### Scenario: Chat Completions 端到端测试
- **WHEN** 运行 Chat Completions 端到端测试时
- **THEN** 系统验证完整请求流程：认证 → Token 解析 → 路由 → 厂商适配 → 响应

#### Scenario: Embeddings 端到端测试
- **WHEN** 运行 Embeddings 端到端测试时
- **THEN** 系统验证 Embeddings 请求的完整处理流程

### Requirement: 流式传输测试
系统 SHALL 提供流式传输测试，验证 SSE 流式响应功能。

#### Scenario: 流式响应测试
- **WHEN** 运行流式响应测试时
- **THEN** 系统验证 SSE 流式响应正确传输和结束标识

#### Scenario: 流式中断测试
- **WHEN** 运行流式中断测试时
- **THEN** 系统验证客户端断开时上游请求被取消

### Requirement: 认证流程测试
系统 SHALL 提供认证流程测试，验证各种认证方式的集成。

#### Scenario: X-User-ID 认证测试
- **WHEN** 运行 X-User-ID 认证测试时
- **THEN** 系统验证 X-User-ID 头正确识别用户并应用隔离

#### Scenario: 本地用户认证测试
- **WHEN** 运行本地用户认证测试时
- **THEN** 系统验证本地用户登录、会话创建、Token 使用流程

### Requirement: Token 生命周期测试
系统 SHALL 提供 Token 生命周期测试，验证 Token 管理的完整流程。

#### Scenario: Token 创建和使用测试
- **WHEN** 运行 Token 创建和使用测试时
- **THEN** 系统验证 Token 创建、存储、解析、使用的完整流程

#### Scenario: Token 自动刷新测试
- **WHEN** 运行 Token 自动刷新测试时
- **THEN** 系统验证 Token 过期前自动刷新功能