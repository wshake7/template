# Admin Playwright 测试

## 场景与目标

- 适用场景：新增或维护 Admin E2E、组件测试、测试配置。
- 目标：用现有 Playwright 配置验证关键页面和组件行为。

## 目录/文件位置

- E2E 测试：`front/apps/admin-react/tests/`
- 组件测试入口：`front/apps/admin-react/playwright/`
- E2E 配置：`front/apps/admin-react/playwright.config.ts`
- 组件测试配置：`front/apps/admin-react/playwright-ct.config.ts`

## 当前脚本

```bash
vp run admin-react#e2e:test
vp run admin-react#e2e:test-ui
vp run admin-react#e2e:show
vp run admin-react#e2e:codegen
```

## 编写规则

1. 测试文件放在现有测试目录中，命名为 `*.spec.ts`。
2. 优先测试用户可观察行为，不绑定实现细节。
3. 页面测试用 role、label、text 等稳定定位。
4. 组件测试使用 `@playwright/experimental-ct-react`。
5. 需要 Mock 时复用项目 MSW handlers。

## 示例

```typescript
import { expect, test } from '@playwright/test'

test('dashboard renders', async ({ page }) => {
  await page.goto('/dashboard')
  await expect(page.getByText('仪表盘')).toBeVisible()
})
```

## 注意事项

1. Playwright API 用法问题必须按 `AGENTS.md` 使用 Context7。
2. 修改前端代码后仍需执行 `vp staged`。
3. 测试失败时先保存失败原因；是否继续大范围修复按 `AGENTS.md` 询问用户。
4. 报告目录和测试输出是生成产物，不手动维护。
