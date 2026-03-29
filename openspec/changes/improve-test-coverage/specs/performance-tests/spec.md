## ADDED Requirements

### Requirement: 延迟性能测试
系统 SHALL 提供延迟性能测试，验证代理延迟增量符合 PRD 要求。

#### Scenario: 单请求延迟测试
- **WHEN** 运行单请求延迟测试时
- **THEN** 系统验证代理延迟增量平均 <5ms，P99 <10ms

#### Scenario: 并发延迟测试
- **WHEN** 运行并发延迟测试时
- **THEN** 系统验证 1000 并发下的延迟表现

### Requirement: 吞吐量性能测试
系统 SHALL 提供吞吐量性能测试，验证 QPS 符合 PRD 要求。

#### Scenario: QPS 测试
- **WHEN** 运行 QPS 测试时
- **THEN** 系统验证单机 8 核 16G 环境下 QPS ≥10000

#### Scenario: 混合请求测试
- **WHEN** 运行混合请求测试时
- **THEN** 系统验证流式/非流式混合请求的吞吐量

### Requirement: 内存占用测试
系统 SHALL 提供内存占用测试，验证内存使用符合 PRD 要求。

#### Scenario: 初始内存测试
- **WHEN** 运行初始内存测试时
- **THEN** 系统验证 Go 实现初始内存 <30MB

#### Scenario: 峰值内存测试
- **WHEN** 运行峰值内存测试时
- **THEN** 系统验证 1000 用户规模下峰值内存 <50MB

### Requirement: 启动时间测试
系统 SHALL 提供启动时间测试，验证服务启动时间符合 PRD 要求。

#### Scenario: 启动时间测试
- **WHEN** 运行启动时间测试时
- **THEN** 系统验证服务启动时间 <100ms

### Requirement: 流式传输性能测试
系统 SHALL 提供流式传输性能测试，验证流式响应性能符合 PRD 要求。

#### Scenario: 流式延迟测试
- **WHEN** 运行流式延迟测试时
- **THEN** 系统验证 SSE 流式响应 chunk 端到端延迟 <1ms

#### Scenario: 零缓冲测试
- **WHEN** 运行零缓冲测试时
- **THEN** 系统验证流式传输零缓冲、零阻塞