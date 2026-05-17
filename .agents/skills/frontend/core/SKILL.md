# Skill: Frontend Package Core

## 何时使用

当任务涉及 `front/packages/core` 的共享前端契约、领域类型、HTTP 常量、通知器类型、页面类型、账号/字典/加密类型时使用。

## 模块职责

`@vp/core` 是框架无关的前端基础契约包。

- 存放类型、常量、接口契约。
- 不放 React 组件、hooks、DOM 操作、请求实现或 app 业务逻辑。
- 当前导出：`account`、`dict`、`encrypt`、`http`、`notifier`、`page` 相关类型和常量。

## 调用关系

```text
@vp/core
├── 被 @vp/react-core 调用
├── 被 @vp/request 调用
├── 被 admin-react 调用
├── 被 app-react 调用
└── 被 app-react-ssr 调用
```

`@vp/core` 处在共享包依赖底层，不应调用其他前端共享包，也不应调用任何 app。

## 放置规则

- 纯类型、枚举、常量、契约放这里。
- React hook、store factory、i18n/env 工厂放 `@vp/react-core`。
- 请求 client、token/cookie、加密请求流程放 `@vp/request`。
- 日期、className、Web Crypto 等工具放 `@vp/utils`。

## 验证

```bash
vp staged
vp run @vp/core#build
```

如果修改 public API，需要追加构建依赖它的包和 app。

