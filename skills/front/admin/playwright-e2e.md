# E2E 测试 (Playwright)

本项目使用 Playwright 进行端到端测试，支持组件测试和页面测试。

## 测试目录

```
front/apps/admin/
├── tests/                      # E2E 测试
│   ├── example.spec.ts
│   └── test.spec.ts
├── playwright/                 # Playwright 组件测试
│   ├── index.tsx               # 测试入口组件
│   └── index.html
├── playwright-report/         # 测试报告 (自动生成)
└── playwright.config.ts       # 测试配置
```

## 测试配置

位置: `front/apps/admin/playwright.config.ts`

主要配置:
- `testDir`: 测试文件目录
- `outputDir`: 报告输出目录
- `timeout`: 测试超时时间
- `retries`: 失败重试次数

## 测试命令

```bash
# 运行 E2E 测试
vp run admin#e2e:test

# UI 模式 (可视化)
vp run admin#e2e:test-ui

# 查看测试报告
vp run admin#e2e:show

# 代码生成 (录制模式)
vp run admin#e2e:codegen
```

## 组件测试配置

位置: `front/apps/admin/playwright-ct.config.ts`

用于测试 React 组件的独立渲染。

## 编写测试

### 页面测试示例

```typescript
// tests/example.spec.ts
import { test, expect } from '@playwright/test'

test('example', async ({ page }) => {
  await page.goto('/')
  await expect(page).toHaveTitle(/Admin/)
})
```

### 组件测试示例

```typescript
// playwright/index.tsx
import { test, expect } from '@playwright/experimental-ct-react'
import { MyComponent } from './MyComponent'

test.use({ viewport: { width: 1280, height: 720 } })

test('MyComponent renders', async ({ mount }) => {
  const component = await mount(<MyComponent />)
  await expect(component).toContainText('Hello')
})
```

## CI 集成

GitHub Actions 配置: `.github/workflows/playwright.yml`

```yaml
- uses: actions/checkout@v4
- uses: voidzero-dev/setup-vp@v1
  with:
    cache: true
- run: vp install
- run: vp run admin#build
- run: vp run admin#e2e:test
```

## 测试工具

- **@faker-js/faker**: 生成测试数据
- **@playwright/test**: Playwright 测试框架
- **@playwright/experimental-ct-react**: React 组件测试插件

## 调试技巧

1. 使用 `vp run admin#e2e:test-ui` 可视化调试
2. 使用 `page.pause()` 暂停测试
3. 使用 `page.screenshot()` 截图
4. 使用 `vp run admin#e2e:show` 查看详细报告
