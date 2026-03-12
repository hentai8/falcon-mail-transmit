# 🔍 评论发布失败诊断指南

## 问题现象

✅ AI 审核完成  
✅ 显示"审核结果已发布"  
❌ GitHub Commit 页面评论数为 0

## 📋 排查步骤

### 步骤 1: 检查 GitHub Actions 日志

1. 访问：`https://github.com/你的用户名/falcon-mail-transmit/actions`
2. 点击最新的 "AI Code Review with Claude" 工作流
3. 展开 **"Run AI Code Review"** 步骤
4. 查看完整日志，特别是最后的输出

**期望看到**：
```
✅ 审核结果已发布到 commit xxx
   查看: https://github.com/...
```

**如果看到错误**：
```
⚠️  发布评论失败 (HTTP 403)
   错误信息: Resource not accessible by integration
```

这说明权限不足，继续下一步。

### 步骤 2: 检查 Workflow 权限设置

**方法 A: 检查仓库级别权限**

1. 进入仓库 **Settings**
2. 左侧菜单选择 **Actions** → **General**
3. 滚动到 **Workflow permissions** 部分
4. 确保选择了：
   - ✅ **Read and write permissions**
   - ✅ **Allow GitHub Actions to create and approve pull requests**
5. 点击 **Save**

**方法 B: 检查工作流文件权限**

查看 `.github/workflows/ai-code-review.yml` 中的 permissions 设置：

```yaml
permissions:
  contents: read
  issues: write
  pull-requests: write
```

如果你看到这个，需要添加：

```yaml
permissions:
  contents: read
  issues: write
  pull-requests: write
  statuses: write      # 添加这行
  checks: write        # 添加这行（可选）
```

### 步骤 3: 手动测试 GitHub Token

在 Actions 中添加测试步骤，验证 Token 权限：

```bash
# 测试能否访问仓库
curl -H "Authorization: token $GITHUB_TOKEN" \
     https://api.github.com/repos/你的用户名/falcon-mail-transmit

# 测试能否创建评论（使用一个已存在的 commit SHA）
curl -X POST \
     -H "Authorization: token $GITHUB_TOKEN" \
     -H "Content-Type: application/json" \
     https://api.github.com/repos/你的用户名/falcon-mail-transmit/commits/COMMIT_SHA/comments \
     -d '{"body":"测试评论"}'
```

### 步骤 4: 查看实际的审核结果

即使评论发布失败，审核内容也会在 Actions 日志中显示。

在 "Run AI Code Review" 步骤的日志底部应该能看到：

```
================================================================================
📋 AI 审核结果（评论发布失败，以下是审核内容）:
================================================================================
# 🤖 AI 代码审核报告
...
（完整的审核内容）
...
================================================================================
```

## 🔧 解决方案

### 方案 1: 修改 Workflow 权限（推荐）

编辑 `.github/workflows/ai-code-review.yml`：

```yaml
jobs:
  ai-code-review:
    runs-on: ubuntu-latest
    permissions:
      contents: read
      pull-requests: write
      issues: write
      statuses: write       # 添加这行
```

提交并推送更改。

### 方案 2: 使用 Personal Access Token

如果仓库级别的 GITHUB_TOKEN 权限不足，可以创建 PAT：

1. 访问：https://github.com/settings/tokens
2. 点击 **Generate new token (classic)**
3. 勾选权限：
   - ✅ `repo` (完整仓库访问)
   - ✅ `workflow`
4. 生成并复制 Token
5. 在仓库 Settings → Secrets 中添加：
   - Name: `PAT_TOKEN`
   - Value: 你的 Personal Access Token

然后修改 workflow：

```yaml
- name: Run AI Code Review
  env:
    ANTHROPIC_API_KEY: ${{ secrets.ANTHROPIC_API_KEY }}
    GITHUB_TOKEN: ${{ secrets.PAT_TOKEN }}  # 使用 PAT 而不是默认 token
```

### 方案 3: 使用 Issues 代替 Commit Comments

如果实在无法发布到 commit，可以改为创建 Issue：

修改 `ai_review.py` 中的逻辑，将审核结果作为 Issue 发布：

```python
# POST https://api.github.com/repos/{owner}/{repo}/issues
{
  "title": f"AI Code Review - Commit {commit_sha[:7]}",
  "body": comment_body,
  "labels": ["ai-review", "automated"]
}
```

## 🎯 下一步

1. **先执行步骤 1 和 2**，查看日志和权限设置
2. **应用方案 1**（最简单）
3. 如果还不行，尝试方案 2
4. 将日志中的错误信息反馈，以便进一步诊断

## 💡 常见错误代码

| HTTP 代码 | 含义 | 解决方法 |
|----------|------|---------|
| 403 | 权限不足 | 检查 Workflow permissions |
| 401 | 认证失败 | 检查 Token 是否有效 |
| 404 | 资源不存在 | 检查仓库名和 commit SHA |
| 422 | 请求格式错误 | 检查请求体格式 |

---

需要帮助？将 Actions 日志中的错误信息发给我。
