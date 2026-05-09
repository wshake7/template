# Skill: Frontend Admin React

## 何时使用

当任务涉及 `front/apps/admin-react` 或 `front/packages/utils`，包括页面、路由、菜单、API、登录鉴权、主题、Mock、表单、表格、测试和构建时使用。

## 核心路径

```text
front/apps/admin-react/
├── src/main.tsx
├── src/router.ts
├── src/routeTree.gen.ts
├── src/routes/              # TanStack file routes
├── src/api/                 # Alova API 封装
├── src/domains/             # 领域类型、HTTP 状态码、Header
├── src/stores/              # Zustand store
├── src/components/          # 通用/业务组件
├── src/config/themes/       # Ant Design 主题
├── src/mocks/               # MSW mock
├── src/utils/               # 表单、加密、菜单、antd helpers
├── locales/                 # i18n 资源
└── tests/, playwright*.ts    # E2E / component test
```

## 关键架构

- 路由使用 TanStack Router 文件路由，`src/routes/**` 通过插件生成 `src/routeTree.gen.ts`。
- 根路由 `src/routes/__root.tsx` 根据 `useAccountStore` 的 token 做登录跳转。
- 应用壳在 `src/routes/_app.tsx`，使用 `ProLayout`、`PageContainer`、菜单树和页签。
- 菜单数据来自路由 `staticData.menu`，由 `src/utils/menu.ts` 转为 ProLayout 菜单。
- HTTP 层在 `src/api/index.ts`，用 Alova + token auth + 请求加密/响应解密 + NProgress。
- 业务 API 按资源拆在 `src/api/*.ts`，领域常量与响应码在 `src/domains/*.ts`。
- 状态使用 Zustand，账号状态持久化在 `src/stores/account.ts`。
- 表单校验优先使用 Zod 与 `src/utils/zod.ts` 的 `useZodForm`。
- 页签缓存配置在 `src/config/tabs.ts`，应用壳在 `_app.tsx` 中维护页签显示、隐藏和刷新行为。

## 常见任务流程

### 新增管理页面

1. 在 `src/routes/_app/**` 新增文件路由，使用 `createFileRoute`。
2. 写入 `staticData.menu`，至少包含 `name` 和 `menuType: 'menu'`；目录型节点使用 `menuType: 'catalog'`。
3. 页面组件优先沿用现有 Ant Design Pro Components 模式，例如 `ProTable`、`ModalForm`、`ProFormText`。
4. 如需后端数据，在 `src/api/<resource>.ts` 增加 API 方法和类型，在页面里通过 Alova hooks 调用。
5. 如果列表项需要行级权限控制，API 响应类型应包含 `canWrite`、`canDelete` 等 boolean 字段，并在 `ProTable` 操作列中据此控制按钮禁用或隐藏。
6. 如需 Mock，在 `src/mocks/handlers/**` 增加处理器，并在聚合文件中导出。

### 修改 API 或鉴权

1. 先读 `src/api/index.ts`、`src/domains/http.ts`、`src/utils/encrypt.ts`。
2. 保持 Header 名称与后端 `admin/internal/domains/headers.go` 对齐。
3. 登录成功路径会写入 token、publicKey、Cookie，并更新 router context。
4. 状态码处理集中在 `HttpCodeCheck`；新增业务码时同步更新前后端常量。
5. 后端新增请求或响应字段时，同步更新对应 `src/api/<resource>.ts` 类型和页面表单控件，保持字段名与 Swagger 定义一致。

### 修改页签缓存

1. 页签刷新间隔统一放在 `src/config/tabs.ts` 的 `TAB_REFRESH_INTERVAL`。
2. `_app.tsx` 中的页签缓存条目需要记录路径、隐藏时间和刷新版本，返回已隐藏页签时按间隔决定是否递增版本。
3. 渲染所有缓存页签但隐藏非活跃页签；通过 `key={path:version}` 触发过期页签重新挂载。
4. 关闭页签时同步删除缓存条目，避免长期保留不再使用的页面状态。

### 修改主题和布局

1. 主题入口在 `src/config/themes/*`，布局入口在 `_app.tsx`。
2. 优先用 Ant Design token、`antd-style` 和现有 CSS 变量，不新增孤立样式体系。
3. 管理后台界面应保持高信息密度、克制、可扫描；避免营销页式的大 hero 和装饰性布局。
4. 侧边栏菜单如需缓存后端动态菜单树，使用 zustand `persist` store 保存原始菜单节点；页面组件只通过 store 读写，不直接访问 `localStorage`。

## 命令

在仓库根目录执行：

```bash
vp run admin-react#dev
vp run admin-react#build
vp staged
```

在应用目录也可执行：

```bash
pnpm --filter admin-react dev
pnpm --filter admin-react build
pnpm --filter admin-react e2e:test
```

## 验证

- 修改前端代码后，提交前执行 `vp staged`。
- 页面/路由/API 行为变化优先执行 `vp run admin-react#build`。
- 交互复杂或布局敏感时，启动开发服务并用浏览器检查桌面和移动视口。
- 只改纯文档或技能文件时不需要运行前端验证。

## 列表分页与表单约定

### Alova 分页

管理后台列表页优先使用 `alova/client` 的 `usePagination`，不要把分页请求写在 `ProTable.request` 里。

标准写法：

```tsx
const {
  data,
  total,
  page,
  pageSize,
  loading,
  update,
  send,
} = usePagination(
  (nextPage, nextPageSize) =>
    XxxApi.list({
      page: nextPage,
      pageSize: nextPageSize,
      orderBy: 'id desc',
      query,
    }),
  {
    initialData: {
      total: 0,
      items: [],
    },
    initialPage: 1,
    initialPageSize: DEFAULT_PAGE_SIZE,
    total: response => response.data?.total ?? 0,
    data: response => response.data?.items ?? [],
    watchingStates: [searchText, enabledFilter],
    debounce: [500, 0],
  },
)
```

`ProTable` 对接方式：

```tsx
<ProTable<Xxx>
  rowKey="id"
  search={false}
  columns={columns}
  dataSource={data}
  loading={loading}
  pagination={{
    showSizeChanger: true,
    current: page,
    pageSize,
    total,
    onChange: (nextPage, nextPageSize) => {
      update({ page: nextPage, pageSize: nextPageSize })
    },
  }}
  options={{
    reload: () => send(),
  }}
/>
```

注意：

- 列表搜索、状态筛选等外部状态放入 `watchingStates`；文本搜索通常配 `debounce: [500]`。
- 刷新当前页用 `send()`；保存、删除后通常 `await send()`。
- 如果表格分页数据之外还需要完整选项树，例如父级菜单选择，单独发 `noPaging: true` 请求加载完整数据，不要复用当前分页页数据。
- 表格序号使用 `(page - 1) * pageSize + index + 1`。

### Zod 与 useZodForm

后台表单优先使用 `zod` schema + `src/utils/zod.ts` 的 `useZodForm`，不要在每个字段上散落重复的 Ant Design required 规则。

标准写法：

```tsx
const XxxFormSchema = z.object({
  name: z.string().trim().min(1, '请输入名称'),
  sortOrder: z.number().min(0, '排序不能小于 0'),
  isEnabled: z.boolean(),
  remark: z.string().optional(),
})

type XxxFormValues = z.infer<typeof XxxFormSchema>

const xxxFormDefaults = XxxFormSchema.parse({
  ...XxxFormSchema.partial().parse({}),
  name: '',
  sortOrder: 0,
  isEnabled: true,
  remark: '',
})

function toXxxFormValues(record: Xxx): XxxFormValues {
  return XxxFormSchema.parse({
    ...xxxFormDefaults,
    ...XxxFormSchema.partial().parse(record),
  })
}

const { form, rules, onFinish } = useZodForm<XxxFormValues>({
  schema: XxxFormSchema,
  async onSubmit(values) {
    if (!values) {
      gMessage.error('请填写完整信息')
      return
    }
    await XxxApi.save(values)
    gMessage.success('保存成功')
    form.resetFields()
    await send()
  },
})
```

`Form` 对接方式：

```tsx
<Form<XxxFormValues> form={form} onFinish={onFinish}>
  <ProFormText name="name" label="名称" rules={rules} />
  <ProFormDigit name="sortOrder" label="排序" rules={rules} />
  <Form.Item name="isEnabled" label="状态" valuePropName="checked" rules={rules}>
    <Switch />
  </Form.Item>
</Form>
```

注意：

- 每个参与校验的字段复用同一个 `rules`，字段级错误由 `useZodForm` 按 zod issue path 回填。
- 不要为表单值手写一份重复的 `interface`；使用 `type XxxFormValues = z.infer<typeof XxxFormSchema>`。
- 默认值和编辑回显转换优先通过 `XxxFormSchema.partial().parse(...)` 过滤/规范化外部数据，再与默认值合并，最终用 `XxxFormSchema.parse(...)` 得到完整表单值。
- 如果默认表单值需要空字符串，基础 `XxxFormSchema` 不要直接写 `.min(1)` 这类会让默认值 parse 失败的规则；把必填、按类型必填等提交校验放到单独的 `XxxSubmitSchema = XxxFormSchema.superRefine(...)` 中，再传给 `useZodForm`。
- 条件字段使用 `superRefine` 做类型相关校验，例如只有菜单类型才要求 `component`。
- 抽屉或弹窗保存按钮用 `form.submit()` 触发表单提交，不要绕过 `onFinish` 手动 `validateFields()`。
- 条件渲染字段需要跨类型切换保留回显时，不要在整个 `Form` 上设置 `preserve={false}`；提交 payload 时按类型过滤无效字段。

