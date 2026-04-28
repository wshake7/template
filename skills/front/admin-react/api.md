# Admin API 请求层

## 场景与目标

- 适用场景：新增或修改 `front/apps/admin-react/src/api/*.ts`。
- 目标：统一 Alova 请求定义、响应类型、分页列表和写操作风格。

## 目录/文件位置

- API 实例：`front/apps/admin-react/src/api/index.ts`
- 资源 API：`front/apps/admin-react/src/api/*.ts`
- 分页类型：`front/apps/admin-react/src/domains/page.ts`
- HTTP 类型和响应码：`front/apps/admin-react/src/domains/http.ts`

## 当前约定

- 一个资源一个 API 文件，如 `role.ts`、`resource.ts`、`dict.ts`。
- 默认导入 API 实例：`import API from './index'`。
- 列表函数返回 Alova Method，由页面的 `usePagination` 负责 `.send()`。
- 创建、更新、删除函数内部调用 `.send()`。
- 分页列表设置 `cacheFor: 0`。
- 请求和响应类型在资源 API 文件内就近定义。

## API 模板

```typescript
import type { PagingRequest, PagingResult } from '~/domains/page'
import API from './index'

export interface Resource {
  id: number
  code: string
  name: string
}

export interface ReqResourceCreate {
  code: string
  name: string
}

export interface ReqResourceUpdate extends Partial<ReqResourceCreate> {
  id: number
}

export interface ReqResourceBatchDelete {
  ids: number[]
}

function resourceList(req: PagingRequest) {
  return API.Post<Res<PagingResult<Resource>>>('/api/sys/resource/list', req, {
    cacheFor: 0,
  })
}

async function resourceCreate(req: ReqResourceCreate) {
  await API.Post<Res>('/api/sys/resource/create', req, { cacheFor: 0 }).send()
}

export const ResourceApi = {
  resourceList,
  resourceCreate,
}
```

## 命名规则

- API 对象：`<Resource>Api`
- 列表：`resourceList`
- 创建：`resourceCreate`
- 更新：`resourceUpdate`
- 删除：`resourceDel`
- 批量操作：`resourceBatchXxx`
- 请求类型：`Req<Resource><Action>`
- 响应数据类型：业务名，如 `Resource`

## 与页面联动

1. `usePagination` 中传入 `ResourceApi.resourceList(params)`。
2. 写操作成功后调用页面的 `send()` 刷新列表。
3. 搜索过滤通过 `PagingRequest.query` 传 JSON 字符串。
4. 后端新增字段后，同步检查 API 类型、表格列和表单 schema。

## 注意事项

1. 不在页面里直接使用 `fetch` 或 `axios`。
2. API 文件只定义请求和类型，不写组件状态。
3. 加密、token、错误提示等横切逻辑留在 `src/api/index.ts`。
4. `PagingRequest.query` 语法见 `skills/backend/orm-query.md`。
