# 意图识别Tagger设计方案

## 概述

基于对RedWolf和Zen MCP Server项目的研究，为Claude Code Companion的tagger系统设计意图识别功能，实现智能路由和成本优化。

## 设计目标

1. **智能意图识别** - 自动识别用户请求的意图类型
2. **成本优化** - 根据意图复杂度选择合适的模型
3. **灵活配置** - 支持自定义意图场景和路由规则
4. **高性能** - 轻量级实现，不影响系统响应速度
5. **可扩展** - 支持多种识别策略和模型集成

## 核心架构

### 1. 意图识别Tagger

```
IntentTagger
├── 意图分析引擎
│   ├── 关键词匹配器
│   ├── 规则引擎
│   └── AI模型接口（可选）
├── 场景配置管理
│   ├── 预定义场景
│   ├── 自定义场景
│   └── 动态场景学习
└── 标签生成器
    ├── 意图标签
    ├── 复杂度标签
    └── 优先级标签
```

### 2. 意图场景分类

#### 代码相关意图
- **代码创建** (`code-creation`)
  - 关键词：写、创建、生成、实现、开发
  - 复杂度：中等
  - 推荐模型：Qwen Coder、Claude Sonnet

- **代码修改** (`code-modification`)
  - 关键词：修改、优化、重构、调整、改进
  - 复杂度：高
  - 推荐模型：Claude-4、GPT-4

- **代码调试** (`code-debugging`)
  - 关键词：调试、修复、解决、错误、bug
  - 复杂度：高
  - 推荐模型：Claude-4、O3

- **代码解释** (`code-explanation`)
  - 关键词：解释、说明、分析、理解
  - 复杂度：低-中等
  - 推荐模型：Qwen Plus、Claude Haiku

#### 对话相关意图
- **简单问答** (`simple-qa`)
  - 关键词：你好、什么是、如何、为什么
  - 复杂度：低
  - 推荐模型：本地小模型、Qwen Plus

- **技术咨询** (`technical-consultation`)
  - 关键词：建议、推荐、最佳实践、架构
  - 复杂度：中等-高
  - 推荐模型：Claude Sonnet、GPT-4

- **学习指导** (`learning-guidance`)
  - 关键词：学习、教程、入门、指南
  - 复杂度：中等
  - 推荐模型：Claude Haiku、Gemini Pro

## 实现策略

### 策略1：关键词匹配 + 规则引擎（推荐）

**优势：**
- 响应速度快
- 资源消耗低
- 可控性强
- 易于调试和维护

**实现方式：**
1. 预定义意图关键词库
2. 使用TF-IDF或词频分析
3. 结合上下文规则判断
4. 支持正则表达式匹配

**配置示例：**
```yaml
intent_detection:
  strategy: "keyword_rules"
  scenes:
    - name: "代码创建"
      tag: "code-creation"
      keywords:
        primary: ["写", "创建", "生成", "实现"]
        secondary: ["函数", "类", "模块", "组件"]
      rules:
        - "包含代码相关词汇"
        - "请求长度 > 10字符"
      weight: 0.8
```

### 策略2：AI驱动意图识别

**优势：**
- 识别准确度高
- 能理解复杂语义
- 自适应学习能力
- 处理歧义能力强

**实现方式：**
1. 集成轻量级NLP模型
2. 使用预训练意图分类模型
3. 支持在线API调用
4. 本地模型推理

**配置示例：**
```yaml
intent_detection:
  strategy: "ai_model"
  model_config:
    provider: "local"  # local, openai, gemini
    model: "intent-classifier-v1"
    confidence_threshold: 0.7
  fallback_strategy: "keyword_rules"
```

### 策略3：混合模式（最佳实践）

**实现流程：**
1. 首先使用关键词匹配快速筛选
2. 对于模糊情况使用AI模型精确识别
3. 结合历史数据进行学习优化
4. 支持用户反馈修正

## 配置系统设计

### 1. 意图场景配置

```yaml
tagging:
  taggers:
    - name: intent-detector
      type: builtin
      builtin_type: intent
      tag: intent
      enabled: true
      priority: 1
      config:
        strategy: "hybrid"
        confidence_threshold: 0.8
        scenes:
          - name: "代码创建"
            tag: "code-creation"
            complexity: "medium"
            keywords:
              primary: ["写", "创建", "生成", "实现", "开发"]
              secondary: ["函数", "类", "方法", "模块"]
            patterns:
              - "写.*函数"
              - "创建.*类"
              - "生成.*代码"
            weight: 1.0
            
          - name: "简单问答"
            tag: "simple-qa"
            complexity: "low"
            keywords:
              primary: ["你好", "什么是", "如何", "为什么"]
            patterns:
              - "^(你好|hi|hello)"
              - "什么是.*"
            weight: 0.5
```

### 2. 路由规则配置

```yaml
routing:
  intent_based_routing:
    enabled: true
    rules:
      - intent: "code-creation"
        preferred_models: ["qwen-coder", "claude-sonnet"]
        fallback_models: ["claude-haiku"]
        
      - intent: "code-debugging"
        preferred_models: ["claude-4", "gpt-4"]
        fallback_models: ["claude-sonnet"]
        
      - intent: "simple-qa"
        preferred_models: ["qwen-plus", "local-llama"]
        fallback_models: ["claude-haiku"]
```

## 技术实现细节

### 1. IntentTagger结构

```go
type IntentTagger struct {
    BaseTagger
    strategy        IntentStrategy
    scenes          []IntentScene
    confidenceThreshold float64
    modelClient     AIModelClient  // 可选AI模型
}

type IntentScene struct {
    Name        string            `json:"name"`
    Tag         string            `json:"tag"`
    Complexity  string            `json:"complexity"`
    Keywords    IntentKeywords    `json:"keywords"`
    Patterns    []string          `json:"patterns"`
    Weight      float64           `json:"weight"`
}

type IntentKeywords struct {
    Primary   []string `json:"primary"`
    Secondary []string `json:"secondary"`
}
```

### 2. 识别流程

```
用户请求 → 文本预处理 → 意图分析
    ↓
关键词匹配 → 规则评估 → 置信度计算
    ↓
AI模型验证（可选） → 最终意图确定 → 标签生成
```

### 3. 性能优化

- **缓存机制** - 缓存常见请求的意图识别结果
- **异步处理** - AI模型调用使用异步方式
- **降级策略** - AI模型不可用时自动降级到规则引擎
- **批量处理** - 支持批量意图识别请求

## 集成方案

### 1. 与现有Tagger系统集成

- 作为新的内置tagger类型添加到factory
- 与其他tagger并行执行
- 支持意图标签与其他标签组合

### 2. 与路由系统集成

- 意图标签作为路由决策依据
- 支持基于意图的端点选择
- 与现有优先级系统兼容

### 3. 与统计系统集成

- 记录意图识别准确率
- 统计不同意图的请求分布
- 支持意图识别效果分析

## 监控与优化

### 1. 性能指标

- **识别准确率** - 意图识别的正确率
- **响应时间** - 意图识别的处理时间
- **覆盖率** - 能够识别的请求比例
- **置信度分布** - 识别结果的置信度统计

### 2. 持续优化

- **用户反馈** - 支持用户标记错误识别
- **模型更新** - 定期更新意图识别模型
- **规则优化** - 基于统计数据优化匹配规则
- **场景扩展** - 根据使用情况添加新场景

## 部署建议

### 1. 分阶段实施

**阶段1：基础实现**
- 实现关键词匹配策略
- 添加基本意图场景
- 集成到现有tagger系统

**阶段2：AI增强**
- 集成AI模型接口
- 实现混合识别策略
- 添加学习优化功能

**阶段3：智能路由**
- 实现基于意图的路由
- 添加成本优化逻辑
- 完善监控和分析

### 2. 配置建议

- 初期使用保守的置信度阈值
- 逐步扩展意图场景覆盖
- 定期分析和调整配置参数

## 总结

通过借鉴RedWolf的智能路由思想和Zen MCP的多模型协作模式，为Claude Code Companion设计的意图识别系统将显著提升系统的智能化水平和成本效益。该设计方案具有良好的可扩展性和实用性，能够满足不同用户的需求。
