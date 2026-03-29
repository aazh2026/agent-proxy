## ADDED Requirements

### Requirement: Provider 路由测试
系统 SHALL 提供全面的 Provider 路由测试，验证模型路由、负载均衡、故障转移逻辑。

#### Scenario: 模型路由测试
- **WHEN** 运行模型路由测试时
- **THEN** 系统验证 gpt-* 模型路由到 OpenAI，claude-* 路由到 Anthropic，gemini-* 路由到 Google

#### Scenario: 自定义别名映射测试
- **WHEN** 运行自定义别名映射测试时
- **THEN** 系统验证自定义别名可以覆盖默认路由规则

#### Scenario: 未知模型测试
- **WHEN** 运行未知模型测试时
- **THEN** 系统验证未知模型返回 404 错误

### Requirement: 负载均衡测试
系统 SHALL 提供负载均衡测试，验证多 Token 轮询和加权分配。

#### Scenario: 轮询负载均衡测试
- **WHEN** 运行轮询负载均衡测试时
- **THEN** 系统验证多个 Token 之间均匀分配请求

#### Scenario: 加权负载均衡测试
- **WHEN** 运行加权负载均衡测试时
- **THEN** 系统验证 Token 按配置权重分配请求

### Requirement: 故障转移测试
系统 SHALL 提供故障转移测试，验证 Token 失效时的自动切换。

#### Scenario: Token 失效测试
- **WHEN** 运行 Token 失效测试时
- **THEN** 系统验证禁用或过期的 Token 被自动跳过

#### Scenario: 跨厂商故障转移测试
- **WHEN** 运行跨厂商故障转移测试时
- **THEN** 系统验证当主厂商失败时自动切换到备用厂商