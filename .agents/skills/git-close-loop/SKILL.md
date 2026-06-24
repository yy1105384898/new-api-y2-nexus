---
name: git-close-loop
description: Git 闭环提交与文档同步：单仓一个 feature 分支、多仓同名分支与配对 commit；对照 AGENTS.md 维护文档、跑 verify、中文提交，验证通过后 --no-ff 合并 main 并 push。
---

# Git 闭环提交与文档同步

将开发过程转化为「闭环」：**在独立分支开发 → 验证 → 合并主分支**。

**仓库差异**（doc 路径、verify 命令）以当前仓库 **`AGENTS.md` 的「Git 闭环」** 为准；本文件规定**通用 Git 行为**，各仓内容须保持一致。

## 分支与提交粒度

### 单仓库（默认）

| 规则 | 说明 |
|------|------|
| 一个功能 / 修复 | **一个** feature 分支（如 `feat/简短描述`、`fix/简短描述`） |
| 分支上的 commit | 通常 **一个** feature commit；仅用户要求或天然独立子任务时才拆多个 |
| 禁止 | 同一功能开多个 feature 分支；在 `main` 上直接 `commit` + `push` 功能/修复 |
| 合并 | 验证通过后 `git merge --no-ff <branch>`，禁止默认 fast-forward / squash |

### 多仓工作区（如 `new-api` + `infinite-canvas`）

各目录是**独立 Git 仓库**，无法共用一个分支，但须当作**同一逻辑变更**对齐：

| 对齐项 | 要求 |
|--------|------|
| 分支名 | **完全相同**，如 `feat/model-vendor-display-name` |
| feature commit header | **相同**：`<type>(<scope>): <同一简要描述>` |
| commit body | 写本仓变更；跨仓时加 `配合：<other-repo> …` |
| merge commit | 同一句话，如 `merge: 合并 feat/xxx（模型渠道别名与展示名）` |
| 完成条件 | **每个有变更的仓库** main 均已 push，再汇报闭环完成 |

**禁止**：多仓各用不同分支 slug（如一个 `feat/pricing-*`、另一个 `feat/model-*`）却声称同一功能已闭环。

```bash
# 单仓标准起手
git checkout main && git pull --rebase origin main
git checkout -b feat/your-topic
# … 开发、文档、verify、commit …
git push -u origin feat/your-topic   # 可选

git checkout main && git pull --rebase origin main
git merge --no-ff feat/your-topic -m "merge: 合并 feat/your-topic（一句话说明）"
git push origin main
```

多仓时：**先在所有涉及仓库确定统一分支名与 commit 文案**，再逐仓执行上述流程。

---

## 提交阶段

### 1. 影响面分析

- 确认当前在 **feature 分支**（非 main，除非用户已授权例外）。
- `git status` / `git diff`，列出代码与文档变更。
- 必读当前仓 `AGENTS.md` **文档影响面**，判断须同步哪些文件。
- 无契约/架构变化：提交 body 写「文档：无」。

### 2. 文档维护

- 按 `AGENTS.md` 文档影响面更新真值层；`AGENTS.md` 只更新指针，不堆长文。
- 合并前执行当前仓 `AGENTS.md` 规定的 **verify**。

### 3. 规范化中文提交

- **Header**：`<type>(<scope>): <简要描述>`（Angular，中文）
- **Body**：变更背景；**文档**：路径列表，无则「文档：无」；跨仓写 `配合：…`
- **Footer**（可选）：`Issue: #123`

### 4. 推送与合并

1. feature 分支：`git status -sb` 干净 → 可选 `git push -u origin HEAD`
2. `git checkout main && git pull --rebase origin main`
3. `git merge --no-ff <branch> -m "merge: …"`
4. `git push origin main`

推送/合并失败须说明原因，**不要**假装闭环完成。

**例外**：用户明确要求「直接提交 main」「hotfix 上 main」或「不要开分支」时可跳过开分支，汇报中说明原因。

---

## 产出

1. 单仓：一个 feature 分支 + feature commit + merge commit
2. 多仓：同名分支 + 配对 commit，各仓 main 已 push
3. 文档与实现一致，verify 已通过

## TL;DR

**单仓一个 feature 分支 → AGENTS 文档/verify → 中文 commit → merge --no-ff → push**；**多仓同名分支、配对 commit、全部 push 后再报完成**。
