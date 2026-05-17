# Skill: Frontend Package Request

## 何时使用

当任务涉及 `front/packages/request` 的 Alova client、请求加密、响应解密、token/cookie、NProgress、HTTP 错误处理、通知器注入时使用。

## 模块职责

`@vp/request` 是跨应用请求层共享包。

- `src/client/alova.ts`：Alova client / 请求流程工厂。
- `src/encryption/helpers.ts`：加密请求、响应解密、签名、公钥等辅助能力。
- `src/index.ts`：统一导出公共 API。

## 调用关系

```text
@vp/request
├── 调用 @vp/core
├── 调用 @vp/utils
├── 被 admin-react 调用
├── 被 app-react 调用
└── 被 app-react-ssr 调用
```

`@vp/request` 不应调用任何 app，也不应直接依赖 Ant Design。需要 UI 提示时通过 `@vp/core` 的通知器契约注入。

## 放置规则

- 跨应用通用请求能力放这里。
- 某个 app 独有的业务 API 定义放 app 的 `src/api/business/**`。
- 账号状态读取、token 提供、通知器实现由 app 注入或薄适配。
- 加密算法底层通用工具优先放 `@vp/utils`，请求流程编排放 `@vp/request`。

## 验证

```bash
vp staged
vp run @vp/request#build
vp run admin-react#build
vp run app-react#build
vp run app-react-ssr#build
```

