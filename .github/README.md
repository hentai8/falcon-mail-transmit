# GitHub Actions & Scripts

这个目录包含项目的 GitHub Actions 工作流和相关脚本。

## 📁 目录结构

```
.github/
├── workflows/
│   └── ai-code-review.yml      # AI 代码审核工作流
├── scripts/
│   └── ai_review.py            # Claude AI 审核脚本
├── AI_CODE_REVIEW_SETUP.md     # AI 审核配置指南
└── README.md                    # 本文件
```

## 🤖 AI 代码审核

本项目使用 Claude AI 自动审核代码变更。

### 快速开始

1. **配置 API Key**
   ```bash
   # 在 GitHub 仓库设置中添加 Secret
   # Name: ANTHROPIC_API_KEY
   # Value: 你的 Claude API Key
   ```

2. **推送代码到 main 分支**
   ```bash
   git push origin main
   ```

3. **查看审核结果**
   - 访问 commit 页面查看 AI 评论
   - 或在 Actions 标签查看执行日志

### 详细文档

请查看 [AI_CODE_REVIEW_SETUP.md](AI_CODE_REVIEW_SETUP.md) 获取完整的配置说明。

## 🔧 工作流说明

### ai-code-review.yml

**触发条件**: Push 到 main 分支

**执行步骤**:
1. 检出代码
2. 识别变更的 Go 文件
3. 获取文件 diff
4. 调用 Claude API 审核
5. 发布审核结果到 commit

**审核内容**:
- 代码质量
- 安全性
- 性能
- Go 最佳实践
- 可维护性

## 🛠️ 脚本说明

### ai_review.py

Python 脚本，负责：
- 调用 Anthropic Claude API
- 构建审核提示词
- 解析审核结果
- 发布到 GitHub

**依赖**:
- `anthropic` - Claude API 客户端
- Python 3.11+

## 📊 使用统计

查看 Actions 页面了解：
- 审核执行次数
- 成功/失败率
- 执行时间

## 🤝 维护指南

### 修改审核标准

编辑 `scripts/ai_review.py` 中的 `_get_system_prompt()` 方法。

### 添加新的工作流

在 `workflows/` 目录下创建新的 `.yml` 文件。

### 测试脚本

本地测试审核脚本：

```bash
cd .github/scripts
python3 ai_review.py --help
```

## 📞 获取帮助

- 查看 [AI_CODE_REVIEW_SETUP.md](AI_CODE_REVIEW_SETUP.md)
- 创建 GitHub Issue
- 联系项目维护者

---

**最后更新**: 2026年3月11日
