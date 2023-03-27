# XTool: 基于 Go 编写一键式 Xray 安装脚本

## Notice
由于官方安装脚本尚未支持 `Xray 1.8.0`, 待官方更新后支持 Reality

## Preview
![](https://github.com/LordPenguin666/XTool/blob/main/img/img.png?raw=true)

## Functions

- 从官方下载并替换带有插件的 Caddy
- Xray 安装, 并自动配置

## Usage

### 已测试系统
- Debian 11
- Rocky Linux 9 

### 下载运行

- x86_64 (amd64)

```shell
wget https://github.com/LordPenguin666/XTool/releases/download/latest/xtool
chmod +x xtool
./xtool
```

- aarch64 (arm64)

```shell
wget https://github.com/LordPenguin666/XTool/releases/download/latest/xtool-arm
chmod +x xtool-arm
./xtool-arm
```