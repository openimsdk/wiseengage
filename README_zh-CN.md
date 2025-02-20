<p align="center">
    <a href="https://openkf.github.io/website" target="_blank">
        <img src="assets/logo-gif/openkf-logo.gif" width="60%" height="30%"/>
    </a>
</p>
<h3 align="center" style="border-bottom: none">
    ⭐️  OpenKF（开放知识流）是一个在线智能客服系统。 ⭐️ <br>
<h3>

<p align=center>
<a href="https://goreportcard.com/report/github.com/openimsdk/openkf"><img src="https://goreportcard.com/badge/github.com/openimsdk/openkf" alt="A+"></a>
<a href="https://github.com/openimsdk/openkf/issues?q=is%3Aissue+is%3Aopen+sort%3Aupdated-desc+label%3A%22good+first+issue%22"><img src="https://img.shields.io/github/issues/OpenIMSDK/OpenKF/good%20first%20issue?logo=%22github%22" alt="good first"></a>
<a href="https://github.com/openimsdk/openkf"><img src="https://img.shields.io/github/stars/OpenIMSDK/OpenKF.svg?style=flat&logo=github&colorB=deeppink&label=stars"></a>
<a href="https://join.slack.com/t/openimsdk/shared_invite/zt-1tmoj26uf-_FDy3dowVHBiGvLk9e5Xkg"><img src="https://img.shields.io/badge/Slack-100%2B-blueviolet?logo=slack&amp;logoColor=white"></a>
<a href="https://github.com/openimsdk/openkf/blob/main/LICENSE"><img src="https://img.shields.io/badge/license-Apache--2.0-green"></a>
<a href="https://golang.org/"><img src="https://img.shields.io/badge/Language-Go-blue.svg"></a>
</p>

</p>

<p align="center">
    <a href="./README.md"><b>英文</b></a> •
    <a href="./README_zh-CN.md"><b>中文</b></a>
</p>

</p>
<br>


## 🧩 特性介绍

1. [OpenKF](https://github.com/openimsdk/openkf) 是一个基于 [OpenIM](https://github.com/OpenIMSDK) 的开源客服系统。
2. 支持 LLM（本地知识库）客服。
3. 支持多渠道客服，并且易于与第三方系统集成。
4. 易于部署和二次开发。

## 📺 系统预览

<details open> 
    <summary>登录注册</summary>
    <img src="assets/images/login_page.png" />
</details>
<details> 
    <summary>系统配置</summary>
    <img src="assets/images/config_page.png" />
</details>
<details> 
    <summary>工作台监控</summary>
    <img src="assets/images/dashboard_page.png" />
</details>
<details> 
    <summary>平台对接</summary>
    <img src="assets/images/platform_page.png" />
</details>
<details> 
    <summary>客服会话</summary>
    <img src="assets/images/session_page.png" />
</details>

## 🛫 快速开始 

> **注意**：你可以快速开始使用 OpenKF。

### 📦 安装

```bash
git clone https://github.com/openimsdk/openkf openkf && export openkf=$(pwd)/openkf && cd $openkf && make
```

### 🚀 运行

> **注意**： 我们需要先运行后端服务器

```bash
make build
```

> 打开另一个终端并运行以下命令

```bash
# make dev
cd web
npm run dev
```

### 📖 贡献者快速入门

善用 Makefile，它可以确保你的项目的质量。

```bash
Usage: make <TARGETS> ...

Targets:
  all                          Build all the necessary targets. 🏗️
  build                        Build binaries by default. 🛠️
  go.build                     Build the binary file of the specified platform. 👨‍💻
  build-multiarch              Build binaries for multiple platforms. 🌍
  tidy                         tidy go.mod 📦
  style                        Code style -> fmt,vet,lint 🎨
  fmt                          Run go fmt against code. ✨
  vet                          Run go vet against code. 🔍
  generate                     Run go generate against code and docs. ✅
  lint                         Run go lint against code. 🔎
  test                         Run unit test ✔️
  cover                        Run unit test with coverage. 🧪
  docker-build                 Build docker image with the manager. 🐳
  docker-push                  Push docker image with the manager. 🔝
  docker-buildx-push           Push docker image with the manager using buildx. 🚢
  copyright-verify             Validate boilerplate headers for assign files. 📄
  copyright-add                Add the boilerplate headers for all files. 📝
  swagger                      Generate swagger document. 📚
  serve-swagger                Serve swagger spec and docs. 🌐
  clean                        Clean all builds. 🧹
  help                         Show this help info. ℹ️
```

> **注意**： 我们强烈推荐你在提交代码之前运行 `make all`。🚀

```bash
make all
```

## 🕋 架构图

![Architecture](./assets/images/architecture.png)

**MVC 架构设计：**

![MVC](./assets/images/mvc.png)

## 🤖 文件目录描述

目录标准化设计结构：

```bash
.
├── assets
│   └── images
├── build
├── deploy
├── docs
├── kf_plugins # 本地知识库和LLM
│   ├── chat
│   ├── config
│   ├── data
│   ├── logs
│   ├── model
│   └── utils
├── scripts
│   ├── githooks
│   └── LICENSE
├── server # OpenKF 后端
│   ├── cmd
│   ├── data
│   ├── docs
│   ├── examples
│   ├── internal
│   ├── logs
│   ├── pkg
│   ├── test
│   └── tools
└── web # OpenKF 前端
    ├── public
    ├── scripts
    └── src
```

##  🗓️ 社区会议

我们欢迎任何人参与我们的社区，我们提供礼品和奖励，每周四晚上欢迎你加入我们。

我们的会议在[OpenIM Slack](https://join.slack.com/t/openimsdk/shared_invite/zt-1tmoj26uf-_FDy3dowVHBiGvLk9e5Xkg)🎯 `openkf` 管道，然后你可以搜索 openkf 管道加入。

我们在 [GitHub 讨论](https://github.com/openimsdk/openkf/discussions/categories/meeting)中记录每次[双周会议](https://github.com/orgs/OpenIMSDK/discussions/categories/meeting)的内容，会议记录我们使用 [Google Docs](https://docs.google.com/document/d/1nx8MDpuG74NASx081JcCpxPgDITNTpIIos0DS6Vr9GU/edit?usp=sharing)编写。

## 🤼‍ 贡献与开发

OpenIMSDK 的目标是构建一个顶级的开源社区。我们有一套标准，在[社区仓库](https://github.com/OpenIMSDK/community)中。

如果你想为 OpenKF 仓库做出贡献，请阅读我们的[贡献者文档](https://github.com/openimsdk/openkf/blob/main/CONTRIBUTING.md)。

开始前，请确保你的更改是需要的。最好的方式是创建一个[新的讨论](https://github.com/openimsdk/openkf/discussions/new/choose)或者[Slack 交流](https://join.slack.com/t/openimsdk/shared_invite/zt-1tmoj26uf-_FDy3dowVHBiGvLk9e5Xkg)，或者如果你发现了问题，先[报告它](https://github.com/openimsdk/openkf/issues/new/choose)。

**微信群：**

> 如果链接失效了，可以 [提一个 issue 说明](https://github.com/openimsdk/openkf/issues/new/choose)

![024ef07a8c5e56d914a677bb54cf587](http://sm.nsddd.top/sm202311271344480.jpg)


## 🚨 许可证

OpenKF 根据 Apache 2.0 许可证授权。请查看[LICENSE](https://github.com/openimsdk/openkf/tree/main/LICENSE)获取完整的许可证文本。

此存储库中显示的 OpenKF 标志，包括其变体和动态版本，在 [OpenKF](https://github.com/openimsdk/openkf)下的 [assets/logo](./assets/logo/) 和 [assets/logo-gif](./assets/logo-gif/) 目录下，受版权法保护。

## 🔮 感谢我们的贡献者！

<a href="https://github.com/openimsdk/openkf/graphs/contributors">   <img src="https://contrib.rocks/image?repo=OpenIMSDK/OpenKF" /> </a>
