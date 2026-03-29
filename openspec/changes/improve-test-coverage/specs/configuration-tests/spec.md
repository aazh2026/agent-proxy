## ADDED Requirements

### Requirement: 配置加载测试
系统 SHALL 提供配置加载测试，验证 YAML 配置文件加载和解析。

#### Scenario: YAML 配置加载测试
- **WHEN** 运行 YAML 配置加载测试时
- **THEN** 系统验证 agent-proxy.yaml 配置文件正确加载

#### Scenario: 配置验证测试
- **WHEN** 运行配置验证测试时
- **THEN** 系统验证无效配置返回错误信息

### Requirement: 环境变量覆盖测试
系统 SHALL 提供环境变量覆盖测试，验证环境变量优先级高于配置文件。

#### Scenario: 环境变量覆盖测试
- **WHEN** 运行环境变量覆盖测试时
- **THEN** 系统验证 AGENT_PROXY_* 环境变量覆盖配置文件设置

#### Scenario: 优先级测试
- **WHEN** 运行优先级测试时
- **THEN** 系统验证命令行参数 > 环境变量 > 配置文件的优先级顺序

### Requirement: 热重载测试
系统 SHALL 提供热重载测试，验证配置文件变更时自动重新加载。

#### Scenario: 配置热重载测试
- **WHEN** 运行配置热重载测试时
- **THEN** 系统验证配置文件变更后自动重新加载

#### Scenario: 热重载回调测试
- **WHEN** 运行热重载回调测试时
- **THEN** 系统验证配置变更触发回调函数执行