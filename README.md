# caddy-csp

[English](#english) | [中文](#中文)

---

# English

A Content-Security-Policy (CSP) plugin for Caddy, supporting flexible CSP header configuration to enhance web application security.

## Installation

Build with xcaddy:
```
xcaddy build --with github.com/shacha086/caddy-csp
```

## Usage Example

Caddyfile example:

```Caddyfile
:8080 {
    route {
        csp {
            # Add default-src policy, allow http://example.com and all https sources
            add default-src http://example.com :https
            # Add img-src policy, allow all https sources
            add img-src :https
            # Remove 'unsafe-inline' from script-src
            remove script-src "'unsafe-inline'"
            # Add 'self' to all policies
            add !all "'self'"
            # Set script-src to only allow 'none'
            set script-src "'none'"
        }
        respond "Hello, CSP!"
    }
}
```

### Directive Reference
- `add <directive> <value>`: Add value to the specified policy.
- `remove <directive> <value>`: Remove value from the specified policy.
- `set <directive> <value>`: Set the specified policy to the value (overwrite).
- `add !all <value>`: Add value to all policies.

Examples:
- `add default-src http://example.com :https` allows http://example.com and all https sources for default-src.
- `add img-src :https` allows all https sources for img-src.
- `remove script-src "'unsafe-inline'"` removes 'unsafe-inline' from script-src.
- `add !all "'self'"` adds 'self' to all policies.
- `set script-src "'none'"` only allows 'none' for script-src.

## FAQ
- To allow multiple sources, use `add` multiple times.
- `:https` means all https sources, `:http` means all http sources.
- For more syntax and usage, refer to the official CSP documentation.

## Contributing
Feel free to submit issues or PRs to improve this project.

## Contact
- Author: shacha086
- GitHub: https://github.com/shacha086/caddy-csp

---

# 中文

一个用于 Caddy 的 Content-Security-Policy (CSP) 插件，支持灵活配置 CSP 头，提升 Web 应用安全性。

## 安装

使用 xcaddy 构建：
```
xcaddy build --with github.com/shacha086/caddy-csp
```

## 使用示例

以下为 Caddyfile 配置示例：

```Caddyfile
:8080 {
    route {
        csp {
            # 添加 default-src 策略，允许 http://example.com 和所有 https 源
            add default-src http://example.com :https
            # 添加 img-src 策略，允许所有 https 源
            add img-src :https
            # 从 script-src 策略中移除 'unsafe-inline'
            remove script-src "'unsafe-inline'"
            # 为所有策略添加 'self'
            add !all "'self'"
            # 设置 script-src 仅允许 'none'
            set script-src "'none'"
        }
        respond "Hello, CSP!"
    }
}
```

### 指令说明
- `add <directive> <value>`：为指定策略添加值。
- `remove <directive> <value>`：从指定策略移除值。
- `set <directive> <value>`：设置指定策略为指定值（覆盖原有值）。
- `add !all <value>`：为所有策略添加值。

示例：
- `add default-src http://example.com :https` 允许 default-src 策略下的 http://example.com 和所有 https 源。
- `add img-src :https` 允许 img-src 策略下的所有 https 源。
- `remove script-src "'unsafe-inline'"` 移除 script-src 策略下的 'unsafe-inline'。
- `add !all "'self'"` 为所有策略添加 'self'。
- `set script-src "'none'"` 仅允许 script-src 策略下的 'none'。

## 常见问题
- 如需允许多个源，可多次使用 `add` 指令。
- `:https` 表示所有 https 源，`:http` 表示所有 http 源。
- 详细语法和更多用法请参考 CSP 官方文档。

## 贡献
欢迎提交 issue 或 PR 改进本项目。

## 联系方式
- 作者: shacha086
- GitHub: https://github.com/shacha086/caddy-csp
