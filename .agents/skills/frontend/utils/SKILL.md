# Skill: Frontend Package Utils

## 何时使用

当任务涉及 `front/packages/utils` 的浏览器安全、框架无关工具时使用，包括 className 合并、日期格式化、Web Crypto 等。

## 模块职责

`@vp/utils` 是前端通用工具包。

- `src/common/cn.ts`：className 合并。
- `src/crypto/webCrypto.ts`：浏览器 Web Crypto 工具。
- `src/date/format.ts`：日期格式化工具。
- `tests/**`：工具包测试。

## 调用关系

```text
@vp/utils
├── 被 @vp/react-core 调用
├── 被 @vp/request 调用
├── 被 admin-react 调用
├── 被 app-react 调用
└── 被 app-react-ssr 调用
```

`@vp/utils` 不应调用任何 app，也不应依赖 React、Alova、Zustand、Ant Design 或业务类型。需要类型契约时优先评估是否应该放在 `@vp/core`。

## 放置规则

- 框架无关、无业务状态、无副作用的工具放这里。
- 纯类型/常量放 `@vp/core`。
- React hooks/factory 放 `@vp/react-core`。
- 请求流程放 `@vp/request`。
- 不要在这里读取应用环境变量、路由、store、cookie 或 UI 消息组件。

## 验证

```bash
vp staged
vp run @vp/utils#build
vp run @vp/utils#test
```

