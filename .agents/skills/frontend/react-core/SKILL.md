# Skill: Frontend Package React Core

## 何时使用

当任务涉及 `front/packages/react-core` 的 React 应用工厂、env、i18n、字典 hooks、Mock helper、Zustand store factory 时使用。

## 模块职责

`@vp/react-core` 是 React 层共享包。

- `src/env/createVpEnv.ts`：应用环境变量工厂。
- `src/i18n/createVpI18n.ts`：i18n 初始化工厂。
- `src/dict/createDictMatchHooks.tsx`：字典匹配 hooks 工厂。
- `src/stores/factories.ts`：通用 store factory。
- `src/mock/helpers.ts`：Mock helper。

## 调用关系

```text
@vp/react-core
├── 调用 @vp/core
├── 调用 @vp/utils
├── 被 admin-react 调用
├── 被 app-react 调用
└── 被 app-react-ssr 调用
```

`@vp/react-core` 不应调用 app，不应包含 Ant Design、后台业务页面或具体路由实现。React、react-i18next、zustand 通过 peer dependencies 由应用侧提供。

## 放置规则

- 多个 React app 都能用的 hooks、工厂、store helper 放这里。
- 纯契约先放 `@vp/core`。
- 请求流程放 `@vp/request`。
- 具体 UI 组件、Ant Design 适配、页面逻辑留在 app 内。
- SSR 和 CSR 都会调用的逻辑必须避免模块顶层浏览器全局对象访问。

## 验证

```bash
vp staged
vp run @vp/react-core#build
vp run app-react#build
vp run app-react-ssr#build
vp run admin-react#build
```

