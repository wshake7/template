# 工具库开发 (front/packages/utils)

本项目包含一个可发布的 TypeScript 工具库，使用 Vite+ 的 `vp pack` 命令构建。

## 目录结构

```
front/packages/utils/
├── src/
│   └── index.ts        # 库入口文件
├── tests/
│   └── index.test.ts   # 测试文件
├── dist/               # 构建输出目录 (自动生成)
├── package.json
├── tsconfig.json
└── vite.config.ts
```

## 构建命令

```bash
# 在 front/packages/utils 目录下执行
vp dev    # 开发模式 (watch)
vp build  # 构建发布版本
vp test   # 运行测试
vp check  # 类型检查 + lint
```

## package.json 组件清单

位置: `front/packages/utils/package.json`

- 包信息:
  - 包名：`@vp/utils`
  - 模块类型：`type: module`
  - 导出：`./dist/index.mjs`
- 脚本:
  - `vp pack`：构建产物
  - `vp pack --watch`：监听构建
  - `vp test`：运行测试
  - `vp check`：类型与规范检查
  - `prepublishOnly -> vp run build`：发布前构建
- 开发依赖:
  - `typescript`
  - `@types/node`
  - `vite-plus`
  - `@typescript/native-preview`（TS 体验增强）
  - `bumpp`（版本管理）

## 添加新工具函数

1. 在 `src/index.ts` 中导出函数:

```typescript
export function myUtil() {
  // 实现
}
```

2. 添加测试 (可选但推荐):

```typescript
// tests/index.test.ts
import { describe, expect, test } from 'vite-plus/test'
import { myUtil } from '../src/index'

describe('myUtil', () => {
  test('should work', () => {
    expect(myUtil()).toBe('expected')
  })
})
```

3. 在 `front/packages/utils` 下运行 `vp build` 构建

## 发布配置

`package.json` 中的发布配置:

```json
{
  "name": "@vp/utils",
  "exports": {
    ".": "./dist/index.mjs",
    "./package.json": "./package.json"
  },
  "publishConfig": {
    "access": "public"
  }
}
```

构建输出为 ESM 格式 (`.mjs`)。

## 工作区配置

根目录 `pnpm-workspace.yaml` 已配置:

```yaml
packages:
  - front/apps/*
  - front/packages/*
```

这使得 `@vp/utils` 可以被 `front/apps/admin` 直接引用。

## 依赖管理

- 使用 Vite+ 的 catalog 机制管理依赖版本
- 构建依赖位于 `catalog:dev`
- 生产依赖位于 `catalog:build`

## 注意事项

1. 工具库应保持零运行时依赖
2. 使用 TypeScript 严格模式
3. 确保测试覆盖率
4. 导出类型应单独使用 `export type` 以获得更好的 Tree-shaking
