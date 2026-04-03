# Fuwari For AcoFork

### 有疑问？尝试 [![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/afoim/fuwari)

> [!CAUTION]
> 本仓库为 AcoFork 深度定制版，包含最新文章与定制功能。若作为模板进行二次开发，建议具备一定 Astro / 前端工程经验。

一个基于 Astro 构建的现代化技术博客系统，面向内容创作、展示与长期维护场景。

## ✨ 核心能力

- 高性能静态站点（Astro 5）
- 响应式布局与全局深色主题（dark-only）/ 彩虹模式
- Markdown / MDX 内容发布与增强渲染
- 全文检索、文章目录、阅读时长、更新提醒
- 代码块增强（行号、折叠、复制按钮扩展）
- 多模块页面：归档、友链、赞助、画廊、文件索引、封面生成
- SEO 与分发：RSS、Sitemap、Robots、重定向
- 统计能力：页面访问统计与可视化入口
- 完整维护脚本：新建文章、图片清理、命名规范化、AI 摘要、差异更新等

## 🛠️ 技术栈

- **框架**：Astro 5.x
- **交互**：Svelte 5
- **样式**：Tailwind CSS + Stylus
- **内容处理**：Remark / Rehype 扩展链路
- **代码高亮**：Expressive Code
- **页面过渡**：Swup
- **包管理**：pnpm
- **代码规范**：Biome

## 🚀 快速开始

### 环境要求

- Node.js 18+
- pnpm 9.x

### 安装依赖

```bash
pnpm install
```

### 本地开发

```bash
pnpm dev
```

### 构建生产版本

```bash
pnpm build
```

### 预览构建结果

```bash
pnpm preview
```

## 📌 常用命令

| 命令 | 作用 | 典型场景 |
| --- | --- | --- |
| `pnpm dev` | 启动开发服务器 | 日常开发 |
| `pnpm build` | 执行图片 CDN 处理并构建静态站点 | 发布前验证 |
| `pnpm preview` | 预览构建产物 | 上线前检查 |
| `pnpm type-check` | TypeScript 类型检查 | 提交前质量把关 |
| `pnpm lint` | Biome 检查并修复 | 统一代码风格 |
| `pnpm format` | 代码格式化 | 批量整理代码 |
| `pnpm new-post <slug>` | 快速创建新文章 | 内容生产 |
| `pnpm clean` | 清理未引用图片 | 资源治理 |
| `pnpm del-space` | 规范图片命名并更新引用 | 跨平台兼容 |
| `pnpm imgf` | 修复相邻图片空行 | 提升渲染与 diff 稳定性 |
| `pnpm update-diff` | 更新文章变更信息 | 更新通知 / 差异展示 |
| `pnpm ai-summary` | 生成或更新文章 AI 摘要 | 内容增强 |
| `pnpm watermark` | 批量添加图片水印 | 资源保护 |
| `pnpm cdnify` | 将图片引用指向 CDN 资源 | 加速访问 |

## 🧩 功能模块

### 1) 内容与文章系统

- Markdown 内容发布（含 frontmatter）
- 文章阅读时长、摘要、目录（TOC）
- 外链处理、标题锚点、数学公式支持
- 更新通知与文章差异展示组件

### 2) 页面模块

- **首页 / 分页**：文章流与筛选展示
- **归档页**：按时间聚合文章
- **友链页**：友链展示与提交入口
- **赞助页**：赞助信息与支持方式
- **画廊页**：Masonry 瀑布流 + 灯箱浏览
- **文件页**：本地公共资源索引 + OneDrive 浏览器
- **封面页**：在线封面生成工具

### 3) 交互与体验

- Swup 页面切换与平滑滚动
- 本地搜索与导航增强
- 主题色切换、彩虹模式、显示设置
- BackToTop、分页、文章排序等控件能力

## 🧠 内容增强能力

项目集成了较完整的内容处理链路：

- Remark：阅读时长、摘要、指令扩展、分节处理等
- Rehype：标题锚点、外链策略、数学渲染、组件化渲染
- 自定义增强：GitHub 风格提示块、URL/GitHub 卡片、AI 提示样式等
- Expressive Code：行号、折叠、复制按钮扩展与主题定制

## ⚙️ 配置指南

主要配置入口：`src/config.ts`

可配置内容包括：

- 站点标题、描述、关键词、语言
- 主题色、Banner、背景图
- 导航菜单与个人信息
- 目录深度、许可证、统计平台（Umami）
- GitHub 编辑入口、图像回退策略

如需调整构建行为与 Markdown 管线，请查看：`astro.config.mjs`。

## 🌐 发布与 SEO

项目已内置常见发布能力：

- `rss.xml`：内容订阅
- `sitemap.xml`：站点地图
- `robots.txt`：抓取策略
- Redirects：内置多条重定向规则
- 静态输出：可部署到任意静态托管平台

推荐平台：Vercel / Cloudflare Pages / Netlify / EdgeOne Pages。

## 📁 项目结构

```text
├── public/                 # 静态资源
├── src/
│   ├── components/         # 页面与功能组件
│   ├── content/            # 文章与内容资源
│   ├── data/               # 友链/赞助等数据
│   ├── layouts/            # 页面布局
│   ├── pages/              # 路由页面
│   ├── plugins/            # Remark/Rehype/代码块插件
│   ├── styles/             # 全局样式
│   ├── types/              # 类型定义
│   └── config.ts           # 核心配置
├── scripts/                # 维护与自动化脚本
└── package.json
```

## 🤝 贡献

欢迎提交 Issue 与 Pull Request。提交前建议执行：

```bash
pnpm lint && pnpm type-check
```

详见：[CONTRIBUTING.md](CONTRIBUTING.md)

## 📄 许可证

[MIT License](LICENSE)

## 🙏 致谢

感谢所有贡献者，特别感谢上游仓库对基础设施的贡献：[saicaca/fuwari](https://github.com/saicaca/fuwari)
