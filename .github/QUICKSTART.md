# 🚀 快速配置指南

## 步骤 1: 获取 Anthropic API Key

1. 访问: https://console.anthropic.com/
2. 注册/登录账号
3. 点击 **API Keys** 
4. 创建新的 API Key
5. 复制并保存好（只显示一次）

## 步骤 2: 配置 GitHub Secret

1. 访问你的 GitHub 仓库
2. 进入: **Settings** → **Secrets and variables** → **Actions**
3. 点击 **New repository secret**
4. 填写：
   ```
   Name: ANTHROPIC_API_KEY
   Value: sk-ant-apixx-xxxxxxxxxxxx  (你的 API Key)
   ```
5. 点击 **Add secret**

## 步骤 3: 启用 GitHub Actions

1. 进入仓库的 **Settings** → **Actions** → **General**
2. 在 **Workflow permissions** 部分:
   - ✅ 选择 **Read and write permissions**
   - ✅ 勾选 **Allow GitHub Actions to create and approve pull requests**
3. 点击 **Save**

## 步骤 4: 测试

推送代码到 main 分支：

```bash
# 提交当前的 GitHub Actions 配置
git add .github/
git commit -m "feat: 添加 AI 代码审核功能"
git push origin main

# 或者如果你在其他分支，先合并到 main
git checkout main
git merge your-branch
git push origin main
```

## 步骤 5: 查看结果

1. **查看 Actions 执行**:
   - 访问: `https://github.com/你的用户名/falcon-mail-transmit/actions`
   - 点击最新的 "AI Code Review with Claude" 工作流
   - 查看执行日志

2. **查看审核评论**:
   - 访问你的 commit 页面
   - AI 审核评论会出现在页面底部

## ✅ 配置完成！

现在每次推送到 main 分支时，AI 都会自动审核你的代码变更。

## 🐛 遇到问题？

查看完整的故障排查指南：[AI_CODE_REVIEW_SETUP.md](AI_CODE_REVIEW_SETUP.md#-故障排查)

## 💡 提示

- 第一次运行可能需要等待几分钟
- 如果没有 Go 文件变更，审核会自动跳过
- 审核结果仅供参考，不会阻塞代码合并
