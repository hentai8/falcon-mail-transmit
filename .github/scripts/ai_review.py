#!/usr/bin/env python3
"""
AI Code Review Script using Claude API
用于审核 Go 代码变更的 AI 脚本
"""

import os
import sys
import json
import argparse
import subprocess
from typing import List, Dict
import anthropic


class CodeReviewer:
    """使用 Claude AI 进行代码审核"""
    
    def __init__(self, api_key: str):
        self.client = anthropic.Anthropic(api_key=api_key)
        self.model = "claude-sonnet-4-6"  # Claude 4.6 Sonnet (稳定版本)
    
    def get_file_diff(self, file_path: str, before_sha: str, current_sha: str) -> str:
        """获取文件的 diff"""
        try:
            if before_sha == "0000000000000000000000000000000000000000":
                # 新文件，显示全部内容
                result = subprocess.run(
                    ["git", "show", f"{current_sha}:{file_path}"],
                    capture_output=True,
                    text=True,
                    check=True
                )
                return f"新文件:\n{result.stdout}"
            else:
                # 现有文件的变更
                result = subprocess.run(
                    ["git", "diff", before_sha, current_sha, "--", file_path],
                    capture_output=True,
                    text=True,
                    check=True
                )
                return result.stdout
        except subprocess.CalledProcessError as e:
            return f"无法获取 diff: {e}"
    
    def get_file_content(self, file_path: str) -> str:
        """获取文件的完整内容"""
        try:
            with open(file_path, 'r', encoding='utf-8') as f:
                return f.read()
        except Exception as e:
            return f"无法读取文件: {e}"
    
    def review_code(self, files: List[str], before_sha: str, current_sha: str) -> str:
        """使用 Claude AI 审核代码"""
        
        # 构建审核提示词
        review_prompt = self._build_review_prompt(files, before_sha, current_sha)
        
        print("📝 正在调用 Claude API 进行代码审核...")
        print(f"📊 审核文件数量: {len(files)}")
        
        try:
            message = self.client.messages.create(
                model=self.model,
                max_tokens=8192,
                temperature=0,
                system=self._get_system_prompt(),
                messages=[
                    {
                        "role": "user",
                        "content": review_prompt
                    }
                ]
            )
            
            review_result = message.content[0].text
            print("✅ AI 审核完成")
            return review_result
            
        except Exception as e:
            error_msg = f"❌ Claude API 调用失败: {str(e)}"
            print(error_msg, file=sys.stderr)
            return error_msg
    
    def _get_system_prompt(self) -> str:
        """获取系统提示词 - 定义 AI 审核的规则和标准"""
        return """你是一个专业的 Go 语言代码审核专家。你的任务是对提交的代码进行全面审核。

审核范围应包括：

1. **代码质量**
   - Go 代码规范和最佳实践
   - 命名约定（变量、函数、类型）
   - 代码复杂度和可读性
   - 重复代码和代码组织

2. **安全性**
   - 潜在的安全漏洞
   - 敏感信息泄露（密钥、密码等）
   - SQL/NoSQL 注入风险
   - XSS、CSRF 等 Web 安全问题
   - 输入验证和错误处理

3. **性能**
   - 并发安全问题（goroutine、channel 使用）
   - 资源泄露（连接、文件句柄、goroutine）
   - 潜在的死锁或竞态条件
   - 内存和 CPU 性能问题
   - 不必要的内存分配

4. **Go 最佳实践**
   - context 的正确使用
   - error 处理模式
   - defer、panic、recover 的使用
   - interface 设计
   - 包的依赖关系

5. **可维护性**
   - 代码注释的完整性
   - 函数和方法的复杂度
   - 测试覆盖率相关建议
   - 文档和示例

审核输出格式要求：
- 使用中文输出
- 使用 Markdown 格式
- 为每个问题标注严重级别：🔴 严重、🟡 警告、🔵 建议
- 提供具体的代码位置和改进建议
- 如果代码质量很好，也要给予正面反馈
- 总结部分要简洁明了

请专业、客观、建设性地进行审核。"""
    
    def _build_review_prompt(self, files: List[str], before_sha: str, current_sha: str) -> str:
        """构建发送给 AI 的审核提示"""
        
        prompt = "# 代码审核请求\n\n"
        prompt += f"请审核以下 Go 代码变更（共 {len(files)} 个文件）：\n\n"
        
        for file_path in files:
            if not os.path.exists(file_path):
                continue
            
            prompt += f"## 📄 文件: `{file_path}`\n\n"
            
            # 获取 diff
            diff = self.get_file_diff(file_path, before_sha, current_sha)
            
            if diff and len(diff.strip()) > 0:
                prompt += "### 变更内容:\n\n"
                prompt += "```diff\n"
                prompt += diff[:5000]  # 限制 diff 长度
                if len(diff) > 5000:
                    prompt += "\n... (diff 太长，已截断)\n"
                prompt += "\n```\n\n"
            
            # 如果是新文件或 diff 较小，也包含当前完整内容（限制大小）
            if before_sha == "0000000000000000000000000000000000000000" or len(diff) < 1000:
                content = self.get_file_content(file_path)
                if content and len(content) < 3000:
                    prompt += "### 当前完整内容:\n\n"
                    prompt += "```go\n"
                    prompt += content
                    prompt += "\n```\n\n"
            
            prompt += "---\n\n"
        
        return prompt


class GitHubCommentPoster:
    """发布审核结果到 GitHub"""
    
    def __init__(self, token: str, repo: str):
        self.token = token
        self.repo = repo
    
    def post_commit_comment(self, commit_sha: str, review_result: str):
        """在 commit 上发布评论"""
        
        # 使用 GitHub CLI 发布评论
        comment_body = self._format_comment(review_result)
        
        # 保存到文件作为备份
        with open('/tmp/review_comment.md', 'w', encoding='utf-8') as f:
            f.write(comment_body)
        
        # 使用 curl 发布评论
        try:
            result = subprocess.run(
                [
                    "curl", "-X", "POST",
                    "-w", "\n%{http_code}",  # 输出 HTTP 状态码
                    f"https://api.github.com/repos/{self.repo}/commits/{commit_sha}/comments",
                    "-H", f"Authorization: token {self.token}",
                    "-H", "Accept: application/vnd.github.v3+json",
                    "-H", "Content-Type: application/json",
                    "-d", json.dumps({"body": comment_body})
                ],
                capture_output=True,
                text=True,
                check=False  # 不自动抛出异常，手动检查
            )
            
            # 解析响应
            output_lines = result.stdout.strip().split('\n')
            http_code = output_lines[-1] if output_lines else "000"
            response_body = '\n'.join(output_lines[:-1]) if len(output_lines) > 1 else ""
            
            # 检查 HTTP 状态码
            if http_code.startswith('2'):  # 2xx 成功
                print(f"✅ 审核结果已发布到 commit {commit_sha[:7]}")
                print(f"   查看: https://github.com/{self.repo}/commit/{commit_sha}")
            else:
                # 发布失败，显示详细错误
                print(f"⚠️  发布评论失败 (HTTP {http_code})", file=sys.stderr)
                print(f"   仓库: {self.repo}", file=sys.stderr)
                print(f"   Commit: {commit_sha}", file=sys.stderr)
                
                # 尝试解析错误信息
                try:
                    error_data = json.loads(response_body)
                    if 'message' in error_data:
                        print(f"   错误信息: {error_data['message']}", file=sys.stderr)
                    if 'documentation_url' in error_data:
                        print(f"   文档: {error_data['documentation_url']}", file=sys.stderr)
                except:
                    if response_body:
                        print(f"   响应: {response_body[:200]}", file=sys.stderr)
                
                # 降级方案：输出到标准输出
                print("\n" + "="*80)
                print("📋 AI 审核结果（评论发布失败，以下是审核内容）:")
                print("="*80)
                print(comment_body)
                print("="*80)
                
        except Exception as e:
            print(f"⚠️  发布评论时发生异常: {e}", file=sys.stderr)
            # 降级方案：输出到标准输出
            print("\n" + "="*80)
            print("📋 AI 审核结果（评论发布失败，以下是审核内容）:")
            print("="*80)
            print(comment_body)
            print("="*80)
    
    def _format_comment(self, review_result: str) -> str:
        """格式化评论内容"""
        header = """# 🤖 AI 代码审核报告

> 此报告由 Claude AI 自动生成，用于辅助代码审核

"""
        footer = """

---

<sub>💡 这是一个自动化的建议，请结合实际情况判断。如有疑问，欢迎讨论。</sub>
"""
        return header + review_result + footer


def main():
    parser = argparse.ArgumentParser(description='AI Code Review with Claude')
    parser.add_argument('--files', help='Changed files (newline separated)')
    parser.add_argument('--files-file', help='File containing list of changed files')
    parser.add_argument('--commit-sha', required=True, help='Current commit SHA')
    parser.add_argument('--repo', required=True, help='Repository name (owner/repo)')
    parser.add_argument('--before-sha', required=True, help='Previous commit SHA')
    
    args = parser.parse_args()
    
    # Get API keys
    api_key = os.environ.get('ANTHROPIC_API_KEY')
    github_token = os.environ.get('GITHUB_TOKEN')
    
    if not api_key:
        print("❌ Error: ANTHROPIC_API_KEY environment variable not set", file=sys.stderr)
        sys.exit(1)
    
    if not github_token:
        print("❌ Error: GITHUB_TOKEN environment variable not set", file=sys.stderr)
        sys.exit(1)
    
    # Parse file list
    if args.files_file:
        # Read from file
        try:
            with open(args.files_file, 'r') as f:
                files = [line.strip() for line in f if line.strip()]
        except Exception as e:
            print(f"❌ Error reading file list: {e}", file=sys.stderr)
            sys.exit(1)
    elif args.files:
        # Parse from argument
        files = [f.strip() for f in args.files.strip().split('\n') if f.strip()]
    else:
        print("❌ Error: Either --files or --files-file must be provided", file=sys.stderr)
        sys.exit(1)
    
    if not files:
        print("ℹ️  No files to review")
        return
    
    print(f"🔍 开始审核 {len(files)} 个文件...")
    for f in files:
        print(f"  - {f}")
    
    # 执行审核
    reviewer = CodeReviewer(api_key)
    review_result = reviewer.review_code(files, args.before_sha, args.commit_sha)
    
    # 发布结果
    poster = GitHubCommentPoster(github_token, args.repo)
    poster.post_commit_comment(args.commit_sha, review_result)
    
    print("\n✨ 代码审核流程完成！")


if __name__ == "__main__":
    main()
