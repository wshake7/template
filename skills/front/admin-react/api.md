# API 请求层

## 场景与目标
- 适用场景：新增/修改前端 API 请求时
- 目标：统一使用 Alova 请求库，保持 API 定义风格一致

## 目录/文件位置
- API 实例配置：`src/api/index.ts`
- 按资源拆分：`src/api/account.ts`、`dict.ts`、`encrypt.ts`、`operationLog.ts`、`role.ts`

## 核心依赖
- `alova`：类 Axios 的跨平台请求库
- `API` 实例（`src/api/index.ts`）：预配置了 Token 认证、AES+RSA 加密、响应拦截、NProgress

## API 实例配置

位置：`src/api/index.ts`

已配置的能力：
- Token 认证（从 AccountStore 读取）
- AES + RSA 混合加密
- 响应拦截器
- NProgress 进度条

## 定义新的 API

一个资源一个文件，通过 `DictApi` 对象聚合导出：

```typescript
// src/api/dict.ts
import type { Res, PagingRequest, PagingResult } from '~/domains/page'
import API from '~/api'

async function typeList(req: PagingRequest) {
  return await API.Post<Res<PagingResult<DictType>>>('/api/sys/dict/type/list', req, {
    cacheFor: 0,
  }).send()
}

export const DictApi = {
  typeList, typeCreate, typeUpdate, typeSwitch, typeDel,
}
```

## 请求方法规范

- 统一使用 `API.Post` 发起请求
- 分页列表：`cacheFor: 0` 禁用缓存
- 请求/响应类型就近定义，分页基础结构复用 `src/domains/page.ts`
- 接口命名语义化：`typeList`、`typeCreate`、`typeUpdate`、`typeSwitch`、`typeDel`

## 现有 API 文件一览

| 文件 | 职责 |
|------|------|
| `src/api/account.ts` | 账户登录/登出/信息 |
| `src/api/dict.ts` | 数据字典 CRUD |
| `src/api/encrypt.ts` | 加密公钥获取 |
| `src/api/operationLog.ts` | 操作日志查询 |
| `src/api/role.ts` | 角色管理 CRUD |

## 注意事项

1. 不直接使用 fetch/axios，统一走 Alova
2. 响应类型使用 `Res<T>` 泛型，错误由拦截器统一处理
3. 分页参数复用 `PagingRequest`，避免重复定义分页 DTO
4. API 文件只做请求定义，不做业务编排
