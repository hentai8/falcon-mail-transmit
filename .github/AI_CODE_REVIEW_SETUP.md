# 🤖 AI 代码审核配置指南

本项目已配置 GitHub Actions 自动使用 Claude AI 对代码变更进行全面审核。

## 📋 功能特性

- ✅ 自动审核 main 分支的每次 push
- ✅ 仅审核变更的 Go 文件（增量审核）
- ✅ 使用 Claude 3.5 Sonnet 模型进行智能分析
- ✅ 在 commit 上自动发布审核评论
- ✅ 不阻塞代码合并，仅提供建议

## 🔧 配置步骤

### 1. 获取 Anthropic API Key

1. 访问 [Anthropic Console](https://console.anthropic.com/)
2. 注册/登录账号
3. 进入 API Keys 页面
4. 创建新的 API Key 并复制保存

### 2. 配置 GitHub Secrets

1. 进入 GitHub 仓库的 **Settings** > **Secrets and variables** > **Actions**
2. 点击 **New repository secret**
3. 添加以下 Secret：
   - **Name**: `ANTHROPIC_API_KEY`
   - **Value**: 你的 Anthropic API Key

### 3. 启用 GitHub Actions

1. 进入仓库的 **Actions** 标签
2. 如果 Actions 被禁用，点击启用
3. 确保工作流有以下权限：
   - Settings > Actions > General > Workflow permissions
   - 选择 "Read and write permissions"
   - 勾选 "Allow GitHub Actions to create and approve pull requests"

## 🎯 审核范围

AI 将从以下几个维度进行全面审核：

### 1. 代码质量
- Go 代码规范和命名约定
- 代码复杂度和可读性
- 重复代码识别

### 2. 安全性
- 安全漏洞扫描
- 敏感信息泄露检测
- 注入攻击风险评估
- 输入验证检查

### 3. 性能
- 并发安全分析（goroutine、channel）
- 资源泄露检测
- 死锁和竞态条件识别
- 性能优化建议

### 4. Go 最佳实践
- Context 使用规范
- Error 处理模式
- defer、panic、recover 正确使用
- Interface 设计建议

### 5. 可维护性
- 代码注释完整性
- 函数复杂度评估
- 测试覆盖率建议

## 🚀 使用方式

### 自动触发

每当你向 `main` 分支推送代码时，GitHub Actions 会自动：

1. 检测变更的 Go 文件
2. 调用 Claude AI 进行审核
3. 在对应的 commit 上发布审核报告

```bash
# 示例：推送代码到 main 分支
git add .
git commit -m "feat: 添加新功能"
git push origin main

# GitHub Actions 会自动运行审核
```

### 查看审核结果

审核完成后，你可以在以下位置查看结果：

1. **Commit 页面**: `https://github.com/[your-repo]/commit/[commit-sha]`
   - 审核评论会自动出现在 commit 下方

2. **Actions 标签**: 查看详细的执行日志
   - `https://github.com/[your-repo]/actions`

## 📊 审核报告示例

审核报告会包含：

```markdown
# 🤖 AI 代码审核报告

## 📄 文件: cmd/server.go

### 🔴 严重问题
- [行 XX] 潜在的并发安全问题...

### 🟡 警告
- [行 XX] 错误处理不够完善...

### 🔵 建议
- [行 XX] 可以使用更简洁的写法...

## 📋 总体评价
...

## ✅ 优点
...
```

## ⚙️ 自定义配置

### 修改触发条件

编辑 [.github/workflows/ai-code-review.yml](.github/workflows/ai-code-review.yml)：

```yaml
on:
  push:
    branches:
      - main
      - develop  # 添加其他分支
```

### 修改审核标准

编辑 [.github/scripts/ai_review.py](.github/scripts/ai_review.py) 中的 `_get_system_prompt()` 方法。

### 调整 Claude 模型

在 `ai_review.py` 中修改：

```python
self.model = "claude-3-5-sonnet-20241022"  # 可改为其他模型
```

可用模型：
- `claude-3-5-sonnet-20241022` - 最新版本（推荐）
- `claude-3-opus-20240229` - 更强大但更慢
- `claude-3-sonnet-20240229` - 平衡选择

## 🐛 故障排查

### 问题 1: API Key 无效

**错误信息**: `Claude API 调用失败: authentication error`

**解决方案**:
1. 检查 ANTHROPIC_API_KEY 是否正确设置
2. 确认 API Key 没有过期
3. 验证 API Key 权限

### 问题 2: 权限不足

**错误信息**: `failed to post comment: 403 Forbidden`

**解决方案**:
1. 检查 Workflow permissions 设置
2. 确保选择了 "Read and write permissions"

### 问题 3: 没有审核评论

**解决方案**:
1. 检查 Actions 标签中的执行日志
2. 确认有 Go 文件变更
3. 查看是否有 API 调用错误

## 💰 费用说明

### Anthropic API 定价（2026年3月）

Claude 3.5 Sonnet:
- 输入: $3 / 百万 tokens
- 输出: $15 / 百万 tokens

### 预估成本

- 每次审核（~1000 行代码）: 约 $0.02 - $0.05
- 每月 100 次提交: 约 $2 - $5

💡 **建议**: 设置 API 使用限额以避免意外费用

## 📚 相关文档

- [Anthropic API 文档](https://docs.anthropic.com/)
- [GitHub Actions 文档](https://docs.github.com/actions)
- [Claude 模型说明](https://docs.anthropic.com/claude/docs/models-overview)

## 🔒 安全注意事项

1. ⚠️ **永远不要**提交 API Key 到代码库
2. ⚠️ 使用 GitHub Secrets 安全存储敏感信息
3. ⚠️ 定期轮换 API Keys
4. ⚠️ 监控 API 使用量和费用

## 🤝 贡献

如果你想改进 AI 审核的质量或添加新功能：

1. 修改 `.github/scripts/ai_review.py`
2. 调整提示词以适应项目需求
3. 测试并提交 PR

## 📞 支持

如有问题或建议，请：
- 创建 GitHub Issue
- 联系项目维护者

---

**最后更新**: 2026年3月11日
