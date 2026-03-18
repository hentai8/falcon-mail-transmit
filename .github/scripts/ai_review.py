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
from typing import List, Dict, Optional
from pathlib import Path
import fnmatch
try:
    import yaml
except ImportError:
    print("⚠️  警告: PyYAML 未安装，将使用默认配置", file=sys.stderr)
    yaml = None
import anthropic


class CodeReviewer:
    """使用 Claude AI 进行代码审核"""
    
    def __init__(self, api_key: str, config_path: Optional[str] = None):
        self.client = anthropic.Anthropic(api_key=api_key)
        self.config = self._load_config(config_path)
        self.model = self.config.get('model', {}).get('name', 'claude-sonnet-4-6')
    
    def _load_config(self, config_path: Optional[str] = None) -> Dict:
        """加载配置文件"""
        if config_path is None:
            config_path = '.github/ai-review-config.yml'
        
        if yaml and os.path.exists(config_path):
            try:
                with open(config_path, 'r', encoding='utf-8') as f:
                    config = yaml.safe_load(f)
                print(f"✅ 已加载配置文件: {config_path}")
                return config
            except Exception as e:
                print(f"⚠️  加载配置文件失败: {e}，使用默认配置", file=sys.stderr)
        else:
            if config_path != '.github/ai-review-config.yml':
                print(f"⚠️  配置文件不存在: {config_path}，使用默认配置", file=sys.stderr)
        
        return self._get_default_config()
    
    def _get_default_config(self) -> Dict:
        """获取默认配置（向后兼容）"""
        return {
            'model': {'name': 'claude-sonnet-4-6', 'max_tokens': 8192, 'temperature': 0},
            'output': {
                'language': 'zh-CN',
                'format': 'markdown',
                'severity_labels': {'critical': '🔴 严重', 'warning': '🟡 警告', 'suggestion': '🔵 建议'}
            },
            'ignore': {'files': [], 'checks': []},
            'severity_thresholds': {'include_positive_feedback': True}
        }
    
    def should_ignore_file(self, file_path: str) -> bool:
        """检查文件是否应该被忽略"""
        ignore_patterns = self.config.get('ignore', {}).get('files', [])
        for pattern in ignore_patterns:
            if fnmatch.fnmatch(file_path, pattern):
                return True
        return False
    
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
    
    def get_changed_functions(self, diff: str) -> List[str]:
        """从 diff 中提取被修改的函数名"""
        functions = []
        lines = diff.split('\n')
        for line in lines:
            # 匹配 Go 函数定义: func FunctionName( 或 func (receiver) FunctionName(
            if line.startswith('+') or line.startswith('-'):
                line = line[1:].strip()
                if line.startswith('func '):
                    # 提取函数名
                    parts = line.split('(')
                    if len(parts) >= 2:
                        func_part = parts[0].replace('func', '').strip()
                        # 处理方法接收者: (r *Receiver) MethodName
                        if func_part.startswith('('):
                            func_part = func_part.split(')')[-1].strip()
                        if func_part:
                            functions.append(func_part)
        return list(set(functions))  # 去重
    
    def find_function_references(self, function_name: str, exclude_file: str) -> List[Dict[str, str]]:
        """在项目中查找函数的引用位置"""
        references = []
        context_config = self.config.get('context_analysis', {})
        max_refs = context_config.get('max_related_files', 3)
        
        if not context_config.get('find_references', True):
            return references
        
        try:
            # 使用 grep 查找函数引用（排除 vendor 目录和当前文件）
            result = subprocess.run(
                ['grep', '-r', '-n', '--include=*.go', '--exclude-dir=vendor',
                 function_name, '.'],
                capture_output=True,
                text=True,
                cwd=os.getcwd()
            )
            
            if result.returncode == 0:
                lines = result.stdout.strip().split('\n')[:max_refs * 3]  # 限制结果数
                for line in lines:
                    if ':' in line:
                        parts = line.split(':', 2)
                        if len(parts) >= 3:
                            file_path = parts[0].lstrip('./')
                            # 排除当前修改的文件和测试文件
                            if file_path != exclude_file and not file_path.endswith('_test.go'):
                                line_num = parts[1]
                                code = parts[2].strip()
                                references.append({
                                    'file': file_path,
                                    'line': line_num,
                                    'code': code
                                })
                                if len(references) >= max_refs:
                                    break
        except Exception as e:
            print(f"⚠️  查找引用失败: {e}", file=sys.stderr)
        
        return references
    
    def get_function_context(self, file_path: str, function_name: str) -> str:
        """获取特定函数的完整代码"""
        try:
            with open(file_path, 'r', encoding='utf-8') as f:
                lines = f.readlines()
            
            # 查找函数定义
            in_function = False
            func_lines = []
            brace_count = 0
            
            for i, line in enumerate(lines):
                # 开始匹配函数
                if not in_function and f'func ' in line and function_name in line:
                    in_function = True
                    func_lines.append(f"// Line {i+1}\n")
                
                if in_function:
                    func_lines.append(line)
                    # 计算大括号来确定函数结束
                    brace_count += line.count('{') - line.count('}')
                    if brace_count == 0 and '{' in ''.join(func_lines):
                        break
            
            return ''.join(func_lines) if func_lines else ""
        except Exception as e:
            return f"无法获取函数上下文: {e}"
    
    def analyze_imports_and_types(self, file_path: str) -> Dict[str, List[str]]:
        """分析文件的 imports 和类型定义"""
        analysis = {'imports': [], 'types': [], 'interfaces': []}
        try:
            with open(file_path, 'r', encoding='utf-8') as f:
                content = f.read()
            
            # 提取 import
            import_block = False
            for line in content.split('\n'):
                line = line.strip()
                if line.startswith('import ('):
                    import_block = True
                    continue
                if import_block:
                    if line == ')':
                        break
                    if line and not line.startswith('//'):
                        analysis['imports'].append(line)
                elif line.startswith('import '):
                    analysis['imports'].append(line.replace('import ', ''))
            
            # 提取类型定义
            for line in content.split('\n'):
                line = line.strip()
                if line.startswith('type ') and ' struct' in line:
                    type_name = line.split()[1]
                    analysis['types'].append(type_name)
                elif line.startswith('type ') and ' interface' in line:
                    type_name = line.split()[1]
                    analysis['interfaces'].append(type_name)
        except Exception as e:
            print(f"⚠️  分析文件结构失败: {e}", file=sys.stderr)
        
        return analysis
    
    def review_code(self, files: List[str], before_sha: str, current_sha: str) -> str:
        """使用 Claude AI 审核代码"""
        
        # 过滤忽略的文件
        filtered_files = [f for f in files if not self.should_ignore_file(f)]
        if len(filtered_files) < len(files):
            ignored_count = len(files) - len(filtered_files)
            print(f"ℹ️  已忽略 {ignored_count} 个文件（根据配置规则）")
        
        if not filtered_files:
            return "所有文件都被忽略，无需审核。"
        
        # 构建审核提示词
        review_prompt = self._build_review_prompt(filtered_files, before_sha, current_sha)
        
        print("📝 正在调用 Claude API 进行代码审核...")
        print(f"📊 审核文件数量: {len(filtered_files)}")
        
        model_config = self.config.get('model', {})
        try:
            message = self.client.messages.create(
                model=self.model,
                max_tokens=model_config.get('max_tokens', 8192),
                temperature=model_config.get('temperature', 0),
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
        """获取系统提示词 - 从配置文件构建审核规则"""
        output_config = self.config.get('output', {})
        language = output_config.get('language', 'zh-CN')
        format_type = output_config.get('format', 'markdown')
        severity_labels = output_config.get('severity_labels', {})
        
        # 构建严重级别标签文本
        severity_text = f"{severity_labels.get('critical', '🔴 严重')}、" \
                       f"{severity_labels.get('warning', '🟡 警告')}、" \
                       f"{severity_labels.get('suggestion', '🔵 建议')}"
        
        # 构建审核范围
        review_scope_text = self._build_review_scope_text()
        
        # 使用配置的模板或默认模板
        template = self.config.get('prompt_template', {}).get('system', '')
        if template:
            return template.format(
                language=language,
                format=format_type,
                severity_labels=severity_text
            ) + "\n\n" + review_scope_text
        
        # 默认提示词
        return f"""你是一个专业的 Go 语言代码审核专家。你的任务是对提交的代码进行全面审核。

{review_scope_text}

审核输出格式要求：
- 使用{language}输出
- 使用 {format_type} 格式
- 为每个问题标注严重级别：{severity_text}
- 提供具体的代码位置和改进建议
- {'如果代码质量很好，也要给予正面反馈' if self.config.get('severity_thresholds', {}).get('include_positive_feedback', True) else ''}
- 总结部分要简洁明了

请专业、客观、建设性地进行审核。"""
    
    def _build_review_scope_text(self) -> str:
        """从配置构建审核范围文本"""
        review_scope = self.config.get('review_scope', {})
        if not review_scope:
            return "请对代码进行全面审核。"
        
        scope_text = "审核范围应包括：\n\n"
        scope_index = 1
        
        scope_mapping = {
            'code_quality': '代码质量',
            'security': '安全性',
            'performance': '性能',
            'go_best_practices': 'Go 最佳实践',
            'maintainability': '可维护性'
        }
        
        for scope_key, scope_name in scope_mapping.items():
            scope_config = review_scope.get(scope_key, {})
            if scope_config.get('enabled', True):
                scope_text += f"{scope_index}. **{scope_name}**\n"
                checks = scope_config.get('checks', [])
                for check in checks:
                    check_name = check.get('name', '')
                    check_desc = check.get('description', '')
                    scope_text += f"   - {check_name}：{check_desc}\n"
                scope_text += "\n"
                scope_index += 1
        
        return scope_text
    
    def _build_review_prompt(self, files: List[str], before_sha: str, current_sha: str) -> str:
        """构建发送给 AI 的审核提示（增强版：包含完整上下文）"""
        
        context_config = self.config.get('context_analysis', {})
        context_enabled = context_config.get('enabled', True)
        max_file_size = context_config.get('max_file_size', 10000)
        
        # 使用配置的用户提示词前缀
        user_prefix = self.config.get('prompt_template', {}).get('user_prefix', '')
        if user_prefix:
            prompt = user_prefix + "\n\n"
        else:
            prompt = "# 代码审核请求\n\n"
        
        # 添加上下文分析说明
        if context_enabled:
            context_instructions = context_config.get('context_instructions', '')
            if context_instructions:
                prompt += f"{context_instructions}\n\n"
        
        prompt += f"请审核以下 Go 代码变更（共 {len(files)} 个文件）：\n\n"
        
        for file_path in files:
            if not os.path.exists(file_path):
                continue
            
            prompt += f"## 📄 文件: `{file_path}`\n\n"
            
            # 获取文件分析
            file_analysis = self.analyze_imports_and_types(file_path)
            if file_analysis['imports'] or file_analysis['types']:
                prompt += "### 文件结构：\n\n"
                if file_analysis['imports']:
                    prompt += f"**依赖包**: {len(file_analysis['imports'])} 个\n"
                if file_analysis['types']:
                    prompt += f"**类型定义**: {', '.join(file_analysis['types'][:5])}\n"
                if file_analysis['interfaces']:
                    prompt += f"**接口定义**: {', '.join(file_analysis['interfaces'][:5])}\n"
                prompt += "\n"
            
            # 获取 diff
            diff = self.get_file_diff(file_path, before_sha, current_sha)
            
            if diff and len(diff.strip()) > 0:
                prompt += "### 变更内容:\n\n"
                prompt += "```diff\n"
                prompt += diff[:5000]  # 限制 diff 长度
                if len(diff) > 5000:
                    prompt += "\n... (diff 太长，已截断)\n"
                prompt += "\n```\n\n"
                
                # 提取被修改的函数
                if context_config.get('analyze_function_context', True):
                    changed_functions = self.get_changed_functions(diff)
                    if changed_functions:
                        prompt += f"### 修改的函数：{', '.join(changed_functions)}\n\n"
                        
                        # 查找这些函数的引用
                        if context_config.get('find_references', True):
                            for func_name in changed_functions[:3]:  # 限制分析的函数数
                                refs = self.find_function_references(func_name, file_path)
                                if refs:
                                    prompt += f"#### 🔗 函数 `{func_name}` 的引用位置：\n\n"
                                    for ref in refs:
                                        prompt += f"- **{ref['file']}:{ref['line']}** - `{ref['code'][:80]}`\n"
                                    prompt += "\n"
                                    
                                    # 包含部分引用文件的上下文
                                    if context_config.get('include_related_files', True):
                                        for ref in refs[:2]:  # 只包含前2个引用文件
                                            ref_content = self.get_file_content(ref['file'])
                                            if ref_content and len(ref_content) < 3000:
                                                prompt += f"**相关文件 `{ref['file']}` 片段：**\n\n"
                                                # 显示引用行附近的上下文（前后各5行）
                                                try:
                                                    ref_lines = ref_content.split('\n')
                                                    line_num = int(ref['line']) - 1
                                                    start = max(0, line_num - 5)
                                                    end = min(len(ref_lines), line_num + 6)
                                                    context_lines = ref_lines[start:end]
                                                    prompt += "```go\n"
                                                    for i, line in enumerate(context_lines, start=start+1):
                                                        marker = '→' if i == line_num + 1 else ' '
                                                        prompt += f"{marker} {i:4d} | {line}\n"
                                                    prompt += "```\n\n"
                                                except:
                                                    pass
            
            # 包含完整文件内容（如果启用且文件不大）
            content = self.get_file_content(file_path)
            if context_config.get('include_full_file', True) and content:
                if len(content) < max_file_size:
                    prompt += "### 当前完整文件内容:\n\n"
                    prompt += "```go\n"
                    prompt += content
                    prompt += "\n```\n\n"
                elif len(diff) < 1000:  # diff 较小但文件较大时，至少包含修改的函数
                    if context_config.get('analyze_function_context', True):
                        changed_functions = self.get_changed_functions(diff)
                        for func_name in changed_functions[:2]:
                            func_context = self.get_function_context(file_path, func_name)
                            if func_context:
                                prompt += f"### 函数 `{func_name}` 完整代码:\n\n"
                                prompt += "```go\n"
                                prompt += func_context
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
    config_path = os.path.join(os.path.dirname(os.path.dirname(__file__)), 'ai-review-config.yml')
    reviewer = CodeReviewer(api_key, config_path)
    review_result = reviewer.review_code(files, args.before_sha, args.commit_sha)
    
    # 发布结果
    poster = GitHubCommentPoster(github_token, args.repo)
    poster.post_commit_comment(args.commit_sha, review_result)
    
    print("\n✨ 代码审核流程完成！")


if __name__ == "__main__":
    main()
