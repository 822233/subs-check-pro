<!-- 项目logo -->
<p align="center">
<a href="https://sinspired.github.io/subs-check-pro/"><img src="https://raw.githubusercontent.com/sinspired/subs-check-pro/main/doc/logo/logo.png" alt="Project Logo" width="200"/></a>
</p>
<h1 align="center">Subs-Check⁺ PRO</h1>
<p align="center" color="#6a737d">
High-performance proxy subscription checker.
</p>

> **搜索关键词 / Search keywords：** Sub-Store 公网入口无法访问、Sub-Store 分享链接地址错误、Subs-Check Pro 自定义后端、Sub-Store URL 修复、Linux amd64 订阅检查、proxy subscription checker、Sub-Store public URL、self-hosted Sub-Store。

> [!IMPORTANT]
> **本仓库是基于上游 `sinspired/subs-check-pro` 的独立修复维护版本，不是上游官方发布。**
>
> 本仓库重点解决特定部署环境下，WebUI 首页的 Sub-Store 入口无法正确打开、分享链接仍指向错误地址的问题，并提供已编译的 Linux amd64 程序。原项目的其他功能和常规部署方式仍以本文后续说明及上游文档为准。

## 🔧 本仓库解决了什么问题？

### 问题背景

在部分自建部署环境中，Subs-Check Pro 主程序本身可以通过公网正常访问，但 WebUI 中集成的 Sub-Store 前端入口不一定能正确识别实际可用的 Sub-Store 服务地址。

这会出现几类现象：

- 访问 Subs-Check Pro 的 WebUI 没有问题，但从首页进入 Sub-Store 时，页面可能无法打开。
- 前端自动探测得到的地址与实际能够从公网访问的地址不一致。
- 即使手动打开了 Sub-Store，生成的分享链接仍可能沿用错误的基础地址，导致分享给其他设备后无法访问。
- 某些部署环境并不需要再增加一层 Cloudflare Tunnel、额外反向代理或新的域名入口；真正需要解决的是 WebUI 如何选择并使用正确的现有服务地址。

### 本仓库的修复方案

本次修复针对已确认的部署入口 `cheap-host1.cheapyun.com`：

| 项目 | 修复行为 |
| --- | --- |
| 指定公网主机 | 使用现有可访问入口 `http://cheap-host1.cheapyun.com:18299` |
| 首页 Sub-Store 入口 | 使用修正后的基础地址生成入口 |
| 分享链接 | 与首页入口共用同一套地址构建逻辑，避免入口和分享地址不一致 |
| 随机访问路径 | 保留原有随机路径，不将其丢弃或替换 |
| 首次访问参数 | 保留原有首次访问所需参数 |
| 自动探测 | 对这个已明确配置的主机跳过不适用的 `scp-store` 探测流程 |
| Cloudflare 提示 | 对该主机不再显示由错误探测流程引起的 Cloudflare 访问警告 |
| 其他域名及局域网 | 继续使用原有自动识别逻辑，不统一强制改成上述公网地址 |

**核心区别：** 修复的是 WebUI 对 Sub-Store 服务基础地址的判断和复用方式，而不是简单替换一个网页链接。首页入口和分享链接使用一致的地址构建逻辑，同时保留原有随机路径及必要参数，避免只修复首页、分享链接却仍然失效的问题。

### 技术实现

本次修改涉及两个相互配合的部分：

1. **WebUI 前端修复**：调整 `getBaseUrl()` / `buildSubStoreUrl()` 相关地址处理逻辑。识别指定公网主机后，使用明确的可访问入口；首页 Sub-Store 入口与分享链接复用该地址构建结果。
2. **主程序构建依赖固定**：主仓库固定引用配套 WebUI Fork 的指定提交，并调整 GoReleaser 构建流程，避免正式构建时重新获取上游前端内容、覆盖已经修复的版本。

这样可以确保源码中的修复不只是临时存在于某次本地构建中，而是随主程序版本一起进入正式编译产物。

### 适用范围与注意事项

- 如果你的部署环境与本次修复场景相同，使用的就是 `cheap-host1.cheapyun.com:18299` 这个公网入口，可以直接使用本仓库对应 Release。
- 如果你使用的是其他域名、其他端口或局域网地址，程序仍保留原有自动探测逻辑；本仓库并不意味着所有第三方 Sub-Store 地址都已自动适配。
- 上述固定入口是本次特定环境的修复配置，不是通用公共服务地址。使用者需要确认该地址从自己的网络环境可以访问。
- 本次修复不等于替用户部署 Sub-Store，也不会自动解决服务器防火墙、端口未开放、DNS 解析错误或后端服务未启动等问题。
- 本仓库是个人维护的修复分支。上游后续更新时，需要在新上游版本上重新检查并移植本次修改；不能假定上游版本已经包含这些修复。

## 📥 修复版下载

请优先从本仓库的 [Releases](https://github.com/822233/subs-check-pro/releases) 页面下载，不要误下上游仓库的 Release。

### 当前发布版本：v3.6.3

| 项目 | 信息 |
| --- | --- |
| 版本 | `v3.6.3` |
| 平台 | Linux |
| CPU 架构 | amd64 / x86_64 |
| 文件大小 | 约 55.1 MB |
| 校验 | 提供 `SHA256SUMS` |
| 构建信息 | 提供 `BUILD-INFO.txt` |

下载：**[进入 v3.6.3 Release](https://github.com/822233/subs-check-pro/releases)**

Linux amd64 用户请根据 Release 资产名称下载对应的二进制文件及校验文件。不要将 amd64 文件用于 ARM64、ARM 或其他架构的设备。

### 完整性校验

下载二进制和 `SHA256SUMS` 后，在 Linux 服务器执行：

```bash
sha256sum -c SHA256SUMS
```

校验结果应显示对应文件为 `OK`。如果校验失败，请重新下载，不要运行文件。

本次 v3.6.3 发布前已执行 Linux amd64 构建、产物架构与内容检查、UPX 压缩完整性检查及解压复核，并上传 Release 资产。SHA256：

```text
e79b4e58d005a83cc78136c501f052a10eaf3616c8d2e2c49f6263fb95691eef
```

> **部署状态说明：** Release 产物构建与上传成功，不代表所有用户环境均已完成部署验证。请按下方步骤在自己的服务器上备份现有配置、替换程序并检查服务日志及 WebUI 功能。

## 🚀 Linux amd64 快速部署

以下为通用替换流程。请先确认服务器架构，并结合你现有的 systemd 服务名称、启动参数和程序目录进行操作；不要直接覆盖未知路径，也不要删除原有配置目录。

### 1. 确认服务器架构

```bash
uname -m
```

输出为 `x86_64` 时，对应本 Release 的 Linux amd64 版本。

### 2. 下载并校验

从 [v3.6.3 Release](https://github.com/822233/subs-check-pro/releases) 下载二进制文件及 `SHA256SUMS`，放到服务器同一目录后执行：

```bash
sha256sum -c SHA256SUMS
```

### 3. 替换前备份

先确认当前服务实际使用的可执行文件路径、启动参数和配置文件路径，再备份现有程序。保留原有 `config`、数据目录、数据库及其他用户配置。

### 4. 替换并重启

将下载并校验通过的二进制文件放到服务实际使用的程序路径，保持原有文件名、执行权限、配置路径和启动参数不变，然后按当前服务器的服务管理方式重启。

### 5. 检查修复效果

建议依次验证：

1. 主程序服务能够正常启动，日志没有新的启动错误。
2. Subs-Check Pro WebUI 可以正常访问。
3. 从 WebUI 首页打开 Sub-Store，确认实际进入的是预期公网入口。
4. 生成一条分享链接，检查链接使用的基础地址是否与首页入口一致。
5. 在另一台设备或不同网络环境中打开分享链接，确认随机路径和首次访问参数仍然有效。
6. 其他原有订阅检测、配置管理及相关功能没有受到影响。

如果入口仍无法访问，应继续检查公网端口连通性、服务监听地址、防火墙规则及 Sub-Store 后端状态；不要仅凭首页能打开就判断 Sub-Store 已修复。

---

<p align="center">
<a href="https://github.com/sinspired/subs-check-pro/releases"><img src="https://img.shields.io/github/v/release/sinspired/subs-check-pro?logo=github&include_prereleases"/></a>
<a href="https://hub.docker.com/r/sinspired/subs-check-pro"><img src="https://img.shields.io/docker/pulls/sinspired/subs-check-pro?logo=docker" /></a>
<a href="https://ghcr.io/sinspired/subs-check-pro"><img src="https://img.shields.io/badge/ghcr.io-online-blue?logo=github"/></a>
<a href="https://sinspired.github.io/subs-check-pro/"><img alt="doc-site" src="https://img.shields.io/website?url=https%3A%2F%2Fsinspired.github.io%2Fsubs-check-pro%2F&up_message=online&down_message=offline&down_color=red&logo=docusaurus&label=%E6%96%87%E6%A1%A3" /></a>
<a href="https://github.com/sinspired/subs-check-pro/wiki"><img alt="Wiki" src="https://img.shields.io/badge/Wiki-docs-blue?logo=github" /></a>
<a href="https://github.com/sinspired/subs-free"><img alt="GitHub Release" src="https://img.shields.io/github/v/release/sinspired/subs-free?include_prereleases&display_name=release&style=flat&logo=wails&logoColor=E2434B&label=%E8%B7%A8%E5%B9%B3%E5%8F%B0%E7%8E%B0%E4%BB%A3GUI" /></a>
</p>

# 🚀 网络代理质量检测工具

**测活、测速、媒体解锁**，网络质量检测工具。采用全新设计，适配 PC 和手机设备的现代 WebUI 配置管理界面，自动生成 `mihomo` 和 `sing-box` 订阅，集成 `sub-store` 前端和后端，支持一键复制分享，可在高并发运行时保持低内存占用；支持自动无缝版本更新。

同时也提供了基于 **Wails v3** 现代化框架开发的跨平台本地图形界面客户端 [subs-free](https://github.com/sinspired/subs-free)，已支持 `Android 手机` 和 `Windows` `Linux` `Mac`。

![preview](./doc/images/Subs-Check-PRO_OG.png)

## ✨ 特性

- [x] ⚡ [自适应流水线高并发模式](https://github.com/sinspired/subs-check-pro/wiki/Features-Details)
- [x] 🗺️ [增强位置标签](https://github.com/sinspired/subs-check-pro/wiki/Features-Details)
- [x] 🎲 [智能节点乱序](https://github.com/sinspired/subs-check-pro/wiki/Features-Details)
- [x] 🕒 [保存并加载历次可用节点](https://github.com/sinspired/subs-check-pro/wiki/Features-Details)
- [x] 📊 [统计订阅链接总数、可用节点数量、成功率](https://github.com/sinspired/subs-check-pro/wiki/Features-Details)
- [x] 🚦 [自动检测本地代理环境](https://github.com/sinspired/subs-check-pro/wiki/System-Proxy)
- [x] 🎁 [自动检查更新，无缝升级新版本](https://github.com/sinspired/subs-check-pro/wiki/Features-Details)
- [x] 📦 [自动生成开箱即用的 sing-box 配置](https://github.com/sinspired/subs-check-pro/wiki/Subscriptions)
- [x] 🔒 [优化文件分享，支持分享码](https://github.com/sinspired/subs-check-pro/wiki/File-Service)
- [x] 📣 [消息通知：节点状态/版本更新/GeoDB 更新](https://github.com/sinspired/subs-check-pro/wiki/Notifications)
- [x] 🔋 [优化内存占用](https://github.com/sinspired/subs-check-pro/wiki/Features-Details)
- [x] ♾️ 支持检测百万-千万量级的节点库，依然保持较低的内存占用
- [x] 💻 支持 `Windows` `Linux` `macOS` 多平台部署
- [x] 📱 支持 `Android` 手机 APP
- [x] 🐳 支持 `docker` 部署
- [x] 🖥️ **全新跨平台现代桌面客户端**：可搭配基于 Wails v3 开发的 [subs-free](https://github.com/sinspired/subs-free) 完美实现本地图形化免命令操控（支持 Windows, macOS, Linux）
- [x] 📱 全新设计的 WebUI 管理界面，优化小屏设备访问体验
- [x] ✏️ 升级配置编辑器，支持自动补全与高亮，内置预览与配置分析，GUI 式编辑
- [x] 🧩 集成 `sub-store` 前端，WebUI 一键管理
- [x] 6️⃣ 支持 `IPv6` 代理节点
- [x] 🔗 适配多种非标订阅格式，提高获取订阅成功率
- [x] 📡 支持检测 `isp` 类型、`原生/广播IP`、住宅/机房
- [x] 📂 内置文件管理页面（`/files`），方便查看与下载生成文件
- [x] 📈 检测结果分析报告（`/analysis`），含地理位置与协议分布可视化
- [ ] 🚧 本项目现接受 issue 反馈

### 📖 教程（Wiki）

- 🐳 [Docker 部署](https://github.com/sinspired/subs-check-pro/wiki/Deployment#docker-运行)
- 🔁 [WatchTower 自动更新与通知](https://github.com/sinspired/subs-check-pro/wiki/Deployment#使用-watchtower-自动更新并通知)
- 📘 [Cloudflare Tunnel 外网访问](https://github.com/sinspired/subs-check-pro/wiki/Cloudflare-Tunnel)
- 📗 [通知渠道配置（Apprise）](https://github.com/sinspired/subs-check-pro/wiki/Notifications)
- 📙 [订阅使用方法](https://github.com/sinspired/subs-check-pro/wiki/Subscriptions)
- 📕 [内置文件服务](https://github.com/sinspired/subs-check-pro/wiki/File-Service)
- 📖 [自建GitHub代理（支持提高api速率和gist加速）](https://github.com/sinspired/CF-Proxy)

### 📣 使用交流，功能讨论，issue 反馈，新版本通知

- > Telegram 群组：[subs-check性能版](https://t.me/subs_check_pro)⁠
- > Telegram 频道：[关注频道](https://t.me/sinspired_ai)⁠

> [!TIP]
>
> 功能更新频繁，请务必查看最新的 [配置文件示例](https://github.com/sinspired/subs-check-pro/blob/main/config/config.yaml.example) 以获取最新功能支持。

> [!NOTE]
> 查看新增功能及设置方法： [新增功能与性能优化详情](https://github.com/sinspired/subs-check-pro/wiki/Features-Details)

## 📸 预览

![preview](./doc/images/preview.png)
![result](./doc/images/results.png)

### ✨ 现代 WebUI 管理界面

`http://localhost:8199/admin`

![login](./doc/images/login.png)
![WebUI](./doc/images/WebUI.png)

### 🖥️ [subs-free](https://github.com/sinspired/subs-free) 跨平台客户端

已支持 `Android 手机` 和 `Windows` `Linux` `Mac`。

#### 💻 桌面端

![preview](./doc/images/GUI.png)

#### 📱 手机 APP

| 正常版                               | Lite 版                                        |
| ------------------------------------ | ---------------------------------------------- |
| ![Android](./doc/images/Android.png) | ![Android_Lite](./doc/images/Android_Lite.png) |

➡️ 立即前往体验：[subs-free 仓库](https://github.com/sinspired/subs-free)

### 📊 检测结果分析报告

`http://localhost:8199/analysis`

![analysis](./doc/images/analysis.png)

### ⚡ 自动生成 singbox 订阅，支持一键分享

|                                          |                                           |                                       |
| ---------------------------------------- | ----------------------------------------- | ------------------------------------- |
| ![sub-store](./doc/images/sub-store.png) | ![shareMenu](./doc/images/share-menu.png) | ![singbox](./doc/images/sing-box.png) |

|                                  |                                        |
| -------------------------------- | -------------------------------------- |
| ![tgram](./doc/images/tgram.png) | ![dingtalk](./doc/images/dingtalk.png) |

## 🛠️ 快速开始

> 首次运行会在当前目录生成默认配置文件。完整安装与部署见 Wiki：
> [安装与部署](https://github.com/sinspired/subs-check-pro/wiki/Deployment)

### 💻 桌面端 UI 运行（推荐）

直接前往 [subs-free](https://github.com/sinspired/subs-free) 下载适配你当前系统（Windows / macOS / Linux）的安装包，开箱即用，无须配置复杂的命令行参数。

### 🌏 WebUI 控制面板

WebUI 集成了配置编辑、订阅分享、订阅管理，内置文件服务，检测结果分析报告，日志查看等功能，请务必使用 WebUI 作为主要管理入口

请主动修改 `config.yaml` `api-key` 作为 WebUI 访问密码，如未设置，请查看终端日志获取 `api-key`

浏览器输入 `http://localhost:8199/admin` 或 `http://127.0.0.1:8199/admin` 访问 WebUI

注意 `8199` 为默认监听端口，如已修改，请替换为实际端口

### 📦 二进制文件运行

下载 Releases 中适合的版本，解压后直接运行：

```powershell
./subs-check-pro.exe -f ./config/config.yaml
```

记录 `api-key`，访问 WebUi

`http://localhost:8199/admin`

### 🐳 Docker 运行（最简）

```bash
docker run -d \
  --name subs-check-pro \
  -p 8299:8299 \
  -p 8199:8199 \
  -v ./config:/app/config \
  -v ./output:/app/output \
  --restart always \
  ghcr.io/sinspired/subs-check-pro:latest
```

💡 记录 `api-key`，访问 WebUi

`http://localhost:8199/admin`

> 代理设置见 Wiki：[系统与 GitHub 代理](https://github.com/sinspired/subs-check-pro/wiki/System-Proxy)
>
> 自建测速地址：[Speedtest](https://github.com/sinspired/subs-check-pro/wiki/Speedtest)

### 🪜 优化系统代理和 GitHub 代理设置（可选）

可使用网上公开的github proxy 或使用 GitHub 项目 [CF-Proxy](https://github.com/sinspired/CF-Proxy) 自建 GitHub 代理

<details>
  <summary>展开查看</summary>

```yaml
# 优先级 1.system-proxy;2.github-proxy;3.ghproxy-group
# 即使未设置,也会检测常见端口(v2ray\clash)的系统代理自动设置

# 系统代理设置: 适用于拉取订阅、消息推送、文件上传等等。
# 写法跟环境变量一样，修改需重启生效
# system-proxy: "http://username:password@192.168.1.1:7890"
# system-proxy: "socks5://username:password@192.168.1.1:7890"
system-proxy: ''
# Github 代理：获取订阅使用
# github-proxy: "https://ghfast.top/"
github-proxy: ''
# GitHub 代理列表：程序会自动筛选可用的 GitHub 代理
ghproxy-group:
# - https://ghp.yeye.f5.si/
# - https://git.llvho.com/
# - https://hub.885666.xyz/
# - https://p.jackyu.cn/
# - https://github.cnxiaobai.com/
```

如果拉取非Github订阅速度慢，可使用通用的 HTTP_PROXY HTTPS_PROXY 环境变量加快速度；此变量不会影响节点测试速度

```bash
# HTTP 代理示例
export HTTP_PROXY=http://username:password@192.168.1.1:7890
export HTTPS_PROXY=http://username:password@192.168.1.1:7890

# SOCKS5 代理示例
export HTTP_PROXY=socks5://username:password@192.168.1.1:7890
export HTTPS_PROXY=socks5://username:password@192.168.1.1:7890

# SOCKS5H 代理示例
export HTTP_PROXY=socks5h://username:password@192.168.1.1:7890
export HTTPS_PROXY=socks5h://username:password@192.168.1.1:7890
```

如果想加速github的链接，可使用网上公开的github proxy

可使用 GitHub 项目 [CF-Proxy](https://github.com/sinspired/CF-Proxy) 自建 GitHub 代理，可直接复制 [worker.js](https://github.com/sinspired/CF-Proxy/blob/main/worker.js) 到 [Cloudflare](https://dash.cloudflare.com/) 的 workers。

```yaml
# Github Proxy，获取订阅使用，结尾要带的 /
# github-proxy: "https://ghfast.top/"
github-proxy: 'https://proxy.custom-domain/'
```

</details>

### 🌐 自建测速地址（可选）

<details>
  <summary>展开查看</summary>

> **⚠️ 注意：** 避免使用 Speedtest 或 Cloudflare 下载链接，因为部分节点会屏蔽测速网站。

1. 将 [worker.js](./doc/cloudflare/worker.js) 部署到 Cloudflare Workers。
2. 绑定自定义域名（避免被节点屏蔽）。
3. 在配置文件中设置 `speed-test-url` 为你的 Workers 地址。
4. 建议使用 `speed-test-url: random` 使用随机测速地址，避免对单个测速源造成负担。

```yaml
# 100MB
speed-test-url: https://custom-domain/speedtest?bytes=104857600
# 1GB
speed-test-url: https://custom-domain/speedtest?bytes=1073741824
```

</details>

## ⚙️ Cloudflare Tunnel 外网访问

完整配置请见 Wiki：[Cloudflare Tunnel](https://github.com/sinspired/subs-check-pro/wiki/Cloudflare-Tunnel)

## 🔔 通知渠道配置

完整文档请见 Wiki：[通知渠道](https://github.com/sinspired/subs-check-pro/wiki/Notifications)

## 🌍 域名相关

> [!TIP]
> Cloudflare Tunnel 外网访问和通知渠道部署需要域名，可自行注册或点击下方链接注册免费域名
>
> [![域名注册](https://img.shields.io/badge/DigitalPlat-注册免费域名-2563eb?style=flat-square&logo=databricks&logoColor=ffffff)](https://dashboard.digitalplat.org/signup?ref=HZcosTVlmQ)

> [!CAUTION]
> 请勿用于非法用途！

## 💾 保存方法

见 Wiki：[保存方法](https://github.com/sinspired/subs-check-pro/wiki/Storage)

## 📲 订阅使用方法

完整文档请见 Wiki：[订阅使用方法](https://github.com/sinspired/subs-check-pro/wiki/Subscriptions)

## 🌐 内置文件服务

完整文档请见 Wiki：[内置文件服务](https://github.com/sinspired/subs-check-pro/wiki/File-Service)

## 📱 安卓手机命令行运行

完整文档请见 Wiki：[安卓手机命令行运行](https://github.com/sinspired/subs-check-pro/wiki/android)

> [!TIP]
> 🎉 Android 手机应用已发布，前往 [Subs Free](https://github.com/sinspired/subs-free) 体验

## ✨ 新增功能与性能优化详情

完整说明请见 Wiki：[新增功能与性能优化详情](https://github.com/sinspired/subs-check-pro/wiki/Features-Details)

## 🤝 贡献与开发

- 欢迎提交 PR 与 Issue，一起完善项目。

```bash
git clone https://github.com/sinspired/subs-check-pro
cd subs-check-pro
```

## 🙏 鸣谢

[beck-8](https://github.com/beck-8)、[Sub-Store](https://github.com/sub-store-org/Sub-Store)、[bestruirui](https://github.com/bestruirui/BestSub)

<!-- ## ⭐ Star History

[![Stargazers over time](https://starchart.cc/sinspired/subs-check-pro.svg?variant=adaptive)](https://starchart.cc/sinspired/subs-check-pro) -->

## ⚖️ 免责声明

本工具仅供学习和研究使用，不鼓励或支持任何非法活动。使用者应自行承担风险并遵守相关法律法规。
