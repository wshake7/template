# Skill: Frontend Admin React

## 何时使用

当任务涉及 `front/apps/admin-react` 或 `front/packages/utils`，包括页面、路由、菜单、API、登录鉴权、主题、Mock、表单、表格、测试和构建时使用。特别适合开发系统管理功能（如角色、用户、菜单、API 资源、字典、语言、API 日志等）。

## 核心路径

text
front/apps/admin-react/
├── src/main.tsx
├── src/router.ts
├── src/routeTree.gen.ts
├── src/routes/
│   ├── _app.tsx
│   ├── _app/account/
│   │   ├── role.tsx
│   │   └── user.tsx
│   ├── _app/system/
│   │   ├── language.tsx
│   │   ├── dict.tsx
│   │   ├── logger/
│   │   │   └── api.log.tsx
│   │   ├── resource.api.tsx
│   │   └── resource.menu.tsx
│   └── _app/dashboard.tsx
├── src/api/
│   ├── index.ts              # Alova 主实例
│   └── business/
│       ├── account.ts
│       ├── sysResourceApi.ts
│       ├── sysResourceMenu.ts
│       ├── sysRole.ts
│       ├── sysUser.ts
│       ├── sysDict.ts
│       ├── sysLanguage.ts
│       └── sysApiLog.ts
├── src/domains/               # 领域类型、HTTP 状态码
├── src/stores/
│   ├── account.ts             # 账号状态
│   ├── resourceMenu.ts        # 动态菜单树（持久化）
│   └── ...
├── src/components/
│   ├── business/
│   │   └── system/
│   │       ├── dictPanels.tsx
│   │       └── languagePanels.tsx
│   └── common/
│       └── antIconPicker.tsx
├── src/hooks/
│   └── useDictMatch.tsx       # 字典批量匹配
├── src/utils/
│   ├── antIcons.tsx           # 可用图标列表
│   ├── date.ts                # 日期格式化工具
│   ├── zod.ts
│   └── ...
└── locales/, tests/, playwright*.ts


## 关键架构

- 路由使用 TanStack Router 文件路由，`src/routes/**` 通过插件生成 `src/routeTree.gen.ts`。
- 根路由 `src/routes/__root.tsx` 根据 token 做登录跳转。
- 应用壳 `src/routes/_app.tsx` 使用 `ProLayout`、`PageContainer`，侧边栏菜单来源于后端 `sys_resource_menu` 表，通过 `useResourceMenuStore` 拉取并持久化缓存。该 store 使用 zustand persist，刷新页面后快速恢复菜单树。
- 系统管理模块下的页面路由统一放在 `src/routes/_app/system/` 目录中；账号管理相关页面（如用户、角色）放在 `src/routes/_app/account/` 目录中。
- HTTP 层在 `src/api/index.ts`，用 Alova + token auth + 请求加密/响应解密 + NProgress。
- 字典数据使用 `useDictMatch` 钩子批量获取并缓存启用字典项，返回 `{ value, label }` 格式的映射，适用于下拉选项、表格列渲染。
- 图标选择使用 `AntIconPicker` 组件，可选图标列表由 `src/utils/antIcons.tsx` 提供。
- 表单校验优先使用 Zod 与 `src/utils/zod.ts` 的 `useZodForm`。
- 页签缓存配置在 `src/config/tabs.ts`，`_app.tsx` 维护页签行为。

## 常见任务流程

### 新增系统管理页面

1. 在 `src/routes/_app/system/`（系统资源）或 `src/routes/_app/account/`（账号管理）下创建文件路由，使用 `createFileRoute`。
2. 在文件路由的 `staticData.menu` 中配置菜单：`name` 为页面标题，`menuType: 'menu'`；父级目录已在对应 `system.tsx` 或 `account.tsx` 中定义为 `catalog`。
3. 页面组件优先使用 Ant Design Pro Components，参照 `resource.menu.tsx` 等实现。
4. 在 `src/api/business/<resource>.ts` 中定义 API 方法和类型，URL 路径使用 `/api/sys/<resource>/*` 格式（与 Swagger 一致）。
5. 若列表需要状态下拉，使用 `useDictMatch('system:is_enabled')` 获取选项。
6. 若表单需要图标字段，使用 `AntIconPicker` 组件。
7. 实现完成后，启动开发服务验证路由、菜单展现和数据交互。

### 使用字典匹配

tsx
import { useDictMatch } from '@/hooks/useDictMatch';

function MyPage() {
  const { dictMap } = useDictMatch(['system:is_enabled']);

  const statusOptions = dictMap['system:is_enabled'] || [];
  // statusOptions => [{ value: '1', label: '启用' }, { value: '0', label: '停用' }]
}


`useDictMatch` 内部会批量请求 `/api/sys/dict/entry/match`，并缓存结果。

### 使用图标选择器

tsx
import AntIconPicker from '@/components/common/antIconPicker';

<Form.Item name="icon" label="图标" valuePropName="value">
  <AntIconPicker />
</Form.Item>


需确保 `src/utils/antIcons.tsx` 已导出所有可选图标。

### 格式化日期和时间

表格列或详情中展示时间时，使用全局可用的 `formatDateYYYYMMDDHHmm(value)`。该函数按 `YYYY-MM-DD HH:mm` 格式化日期字符串，无效输入返回原值或 `-`。已通过自动导入全局注册，无需显式 `import`。

tsx
// 表格列中使用
{
  title: '创建时间',
  dataIndex: 'createdAt',
  render: (_, record) => formatDateYYYYMMDDHHmm(record.createdAt),
}


如果需要其他日期格式，可以使用已安装的 `dayjs` 库直接调用，`dayjs` 已在 `package.json` 中按 `catalog:build` 引入。

### 修改菜单资源后刷新侧边栏

- 在 `resource.menu.tsx` 页面增删改菜单后，调用 `useResourceMenuStore.getState().refresh()` 强制刷新侧边栏缓存。

### 修改 API 或鉴权

- 前端 API 路径变更后，同步更新 `src/api/business/` 对应文件的 URL。
- Header 名称保持与后端 `admin/internal/domains/headers.go` 对齐。
- 业务码处理集中在 `HttpCodeCheck`，新增业务码同步更新前后端常量。

## 列表分页与表单约定

### Alova 分页

管理后台列表页优先使用 `alova/client` 的 `usePagination`，不要把分页请求写在 `ProTable.request` 里。

标准写法：

tsx
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


`ProTable` 对接方式：

tsx
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


注意：

- 列表搜索、状态筛选等外部状态放入 `watchingStates`；文本搜索通常配 `debounce: [500]`。
- 刷新当前页用 `send()`；保存、删除后通常 `await send()`。
- 如果表格分页数据之外还需要完整选项树，例如父级菜单选择，单独发 `noPaging: true` 请求加载完整数据，不要复用当前分页页数据。
- 表格序号使用 `(page - 1) * pageSize + index + 1`。

### Zod 与 useZodForm

后台表单优先使用 `zod` schema + `src/utils/zod.ts` 的 `useZodForm`，不要在每个字段上散落重复的 Ant Design required 规则。

标准写法：

tsx
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


`Form` 对接方式：

tsx
<Form<XxxFormValues> form={form} onFinish={onFinish}>
  <ProFormText name="name" label="名称" rules={rules} />
  <ProFormDigit name="sortOrder" label="排序" rules={rules} />
  <Form.Item name="isEnabled" label="状态" valuePropName="checked" rules={rules}>
    <Switch />
  </Form.Item>
</Form>


注意：

- 每个参与校验的字段复用同一个 `rules`，字段级错误由 `useZodForm` 按 zod issue path 回填。
- 不要为表单值手写一份重复的 `interface`；使用 `type XxxFormValues = z.infer<typeof XxxFormSchema>`。
- 默认值和编辑回显转换优先通过 `XxxFormSchema.partial().parse(...)` 过滤/规范化外部数据，再与默认值合并，最终用 `XxxFormSchema.parse(...)` 得到完整表单值。
- 如果默认表单值需要空字符串，基础 `XxxFormSchema` 不要直接写 `.min(1)` 这类会让默认值 parse 失败的规则；把必填、按类型必填等提交校验放到单独的 `XxxSubmitSchema = XxxFormSchema.superRefine(...)` 中，再传给 `useZodForm`。
- 条件字段使用 `superRefine` 做类型相关校验，例如只有菜单类型才要求 `component`。
- 抽屉或弹窗保存按钮用 `form.submit()` 触发表单提交，不要绕过 `onFinish` 手动 `validateFields()`。
- 条件渲染字段需要跨类型切换保留回显时，不要在整个 `Form` 上设置 `preserve={false}`；提交 payload 时按类型过滤无效字段。

## 命令

在仓库根目录执行：

bash
vp run admin-react#dev
vp run admin-react#build
vp staged


在应用目录也可执行：

bash
pnpm --filter admin-react dev
pnpm --filter admin-react build
pnpm --filter admin-react e2e:test


## 验证

- 修改前端代码后，提交前执行 `vp staged`。
- 页面/路由/API 行为变化优先执行 `vp run admin-react#build`。
- 交互复杂或布局敏感时，启动开发服务并用浏览器检查桌面和移动视口。
- 只改纯文档或技能文件时不需要运行前端验证。
