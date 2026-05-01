import type { ProColumns } from '@ant-design/pro-components'
import type { DictEntry, DictType } from '~/api/dict'
import { ModalForm, ProFormDigit, ProFormSwitch, ProFormText, ProTable } from '@ant-design/pro-components'
import { createFileRoute } from '@tanstack/react-router'
import { usePagination } from 'alova/client'
import {
  Button,
  Form,
  Input,
  Popconfirm,
  Select,
  Space,
  Splitter,
  Tag,
} from 'antd'
import { useCallback, useEffect, useMemo, useRef, useState } from 'react'
import z from 'zod'
import { DictApi } from '~/api/dict'
import { gMessage } from '~/utils/antd'
import { useZodForm } from '~/utils/zod'

export const Route = createFileRoute('/_app/system/dict')({
  staticData: {
    menu: {
      name: '数据字典',
      menuType: 'menu',
    },
  },
  staleTime: 1000 * 60 * 2,
  component: RouteComponent,
})

function statusTag(isEnabled: boolean) {
  if (isEnabled) {
    return <Tag color="success">启用</Tag>
  }
  return <Tag color="default">停用</Tag>
}

const ENTRY_LABEL_PLACEHOLDER = '$' + '{EntryLabel}'
const LABEL_COMPONENT_CUSTOM = '__custom__'
const CUSTOM_LABEL_COMPONENT_DEFAULT = `<Tag color="blue">${ENTRY_LABEL_PLACEHOLDER}</Tag>`

const labelComponentPresets = [
  { label: '纯文本', value: '' },
  { label: '默认', value: `<Tag>${ENTRY_LABEL_PLACEHOLDER}</Tag>` },
  { label: '成功', value: `<Tag color="success">${ENTRY_LABEL_PLACEHOLDER}</Tag>` },
  { label: '处理中', value: `<Tag color="processing">${ENTRY_LABEL_PLACEHOLDER}</Tag>` },
  { label: '警告', value: `<Tag color="warning">${ENTRY_LABEL_PLACEHOLDER}</Tag>` },
  { label: '错误', value: `<Tag color="error">${ENTRY_LABEL_PLACEHOLDER}</Tag>` },
  { label: '蓝色', value: `<Tag color="blue">${ENTRY_LABEL_PLACEHOLDER}</Tag>` },
  { label: '紫色', value: `<Tag color="purple">${ENTRY_LABEL_PLACEHOLDER}</Tag>` },
]

const labelComponentPresetValues = new Set(labelComponentPresets.map(item => item.value))

interface ParsedTagTemplate {
  content: string
  color?: string
  bordered?: boolean
}

const blockedHtmlTags = new Set(['script', 'style', 'iframe', 'object', 'embed', 'link', 'meta', 'base'])

function parseTagTemplate(template?: string): ParsedTagTemplate | null {
  const source = template?.trim()
  if (!source) {
    return { content: ENTRY_LABEL_PLACEHOLDER }
  }

  const openEnd = source.indexOf('>')
  if (openEnd <= 1) {
    return null
  }
  const openTag = source.slice(1, openEnd).trim()
  const attrStart = openTag.search(/\s/)
  const tag = attrStart === -1 ? openTag : openTag.slice(0, attrStart)
  if (tag !== 'Tag') {
    return null
  }
  const closeTag = `</${tag}>`
  if (!source.endsWith(closeTag)) {
    return null
  }
  const attrs = attrStart === -1 ? '' : openTag.slice(attrStart)
  const content = source.slice(openEnd + 1, -closeTag.length)

  const parsed: ParsedTagTemplate = {
    content: content.trim(),
  }
  let consumed = ''
  const attrRegex = /\s+([a-zA-Z]+)=(?:"([^"]*)"|'([^']*)'|\{(true|false)\})/g
  for (let attrMatch = attrRegex.exec(attrs); attrMatch; attrMatch = attrRegex.exec(attrs)) {
    consumed += attrMatch[0]
    const [, name, doubleQuotedValue, singleQuotedValue, boolValue] = attrMatch
    const stringValue = doubleQuotedValue ?? singleQuotedValue
    if (name === 'color' && stringValue && /^[\w#-]+$/.test(stringValue)) {
      parsed.color = stringValue
      continue
    }
    if (name === 'bordered' && boolValue) {
      parsed.bordered = boolValue === 'true'
      continue
    }
    return null
  }

  if (consumed.trim() !== attrs.trim()) {
    return null
  }
  return parsed
}

function sanitizeLabelHtml(labelComponent: string, entryLabel: string) {
  if (typeof document === 'undefined') {
    return ''
  }
  const template = document.createElement('template')
  template.innerHTML = labelComponent.replaceAll(ENTRY_LABEL_PLACEHOLDER, entryLabel)

  const walk = document.createTreeWalker(template.content, NodeFilter.SHOW_ELEMENT)
  const blockedNodes: Element[] = []
  while (walk.nextNode()) {
    const element = walk.currentNode as Element
    const tagName = element.tagName.toLowerCase()
    if (blockedHtmlTags.has(tagName)) {
      blockedNodes.push(element)
      continue
    }
    for (const attr of Array.from(element.attributes)) {
      const attrName = attr.name.toLowerCase()
      const attrValue = attr.value.trim().toLowerCase()
      if (attrName.startsWith('on')) {
        element.removeAttribute(attr.name)
        continue
      }
      if ((attrName === 'href' || attrName === 'src') && /^(?:javascript|data):/.test(attrValue)) {
        element.removeAttribute(attr.name)
      }
    }
  }
  for (const node of blockedNodes) {
    node.remove()
  }
  return template.innerHTML
}

function renderDictEntryLabel(labelComponent: string | undefined, entryLabel: string) {
  if (!labelComponent?.trim()) {
    return entryLabel
  }
  const parsed = parseTagTemplate(labelComponent)
  if (!parsed) {
    const html = sanitizeLabelHtml(labelComponent, entryLabel)
    if (!html) {
      return entryLabel
    }
    return (
      // eslint-disable-next-line react-dom/no-dangerously-set-innerhtml
      <span dangerouslySetInnerHTML={{ __html: html }} />
    )
  }
  const children = parsed.content.includes(ENTRY_LABEL_PLACEHOLDER)
    ? parsed.content.replaceAll(ENTRY_LABEL_PLACEHOLDER, entryLabel)
    : parsed.content

  return (
    <Tag color={parsed.color} bordered={parsed.bordered}>
      {children}
    </Tag>
  )
}

interface ShikiHighlighter {
  codeToHtml: (code: string, options: { lang: string, theme: string }) => string
}

let shikiHighlighterPromise: Promise<ShikiHighlighter> | undefined

function getShikiHighlighter() {
  shikiHighlighterPromise ??= Promise
    .all([
      import('shiki/core'),
      import('shiki/engine/javascript'),
      import('shiki/langs/tsx.mjs'),
      import('shiki/themes/github-light.mjs'),
    ])
    .then(([{ createHighlighterCore }, { createJavaScriptRegexEngine }, tsx, githubLight]) => createHighlighterCore({
      themes: [githubLight.default],
      langs: [tsx.default],
      engine: createJavaScriptRegexEngine({ forgiving: true }),
    }))
  return shikiHighlighterPromise
}

function ShikiCodeEditorBlock({
  value = '',
  onChange,
}: {
  value?: string
  onChange?: (value: string) => void
}) {
  const [highlightValue, setHighlightValue] = useState(value)
  const [html, setHtml] = useState('')
  const highlightRef = useRef<HTMLDivElement>(null)
  const syncTimerRef = useRef<number | undefined>(undefined)

  useEffect(() => {
    return () => {
      if (syncTimerRef.current !== undefined) {
        window.clearTimeout(syncTimerRef.current)
      }
    }
  }, [])

  useEffect(() => {
    let cancelled = false
    const timer = window.setTimeout(async () => {
      try {
        const highlighter = await getShikiHighlighter()
        const nextHtml = highlighter.codeToHtml(highlightValue || ' ', {
          lang: 'tsx',
          theme: 'github-light',
        })
        if (!cancelled) {
          setHtml(nextHtml)
        }
      }
      catch {
        if (!cancelled) {
          setHtml('')
        }
      }
    }, 600)

    return () => {
      cancelled = true
      window.clearTimeout(timer)
    }
  }, [highlightValue])

  return (
    <div
      className="dict-shiki-editor"
      style={{
        position: 'relative',
        minHeight: 116,
        border: '1px solid var(--ant-color-border)',
        borderRadius: 6,
        background: 'var(--ant-color-bg-container)',
        overflow: 'hidden',
      }}
    >
      <style>
        {`
          .dict-shiki-editor .shiki {
            margin: 0;
            min-height: 114px;
            padding: 10px 11px;
            background: transparent !important;
            font-family: ui-monospace, SFMono-Regular, Consolas, monospace;
            font-size: 13px;
            line-height: 1.6;
            white-space: pre-wrap;
            word-break: break-word;
          }
        `}
      </style>
      {/* eslint-disable react-dom/no-dangerously-set-innerhtml */}
      <div
        ref={highlightRef}
        style={{
          position: 'absolute',
          inset: 0,
          overflow: 'hidden',
          pointerEvents: 'none',
        }}
        dangerouslySetInnerHTML={{ __html: html }}
      />
      {/* eslint-enable react-dom/no-dangerously-set-innerhtml */}
      <Input.TextArea
        defaultValue={value}
        rows={4}
        spellCheck={false}
        onBlur={(event) => {
          const nextValue = event.currentTarget.value
          if (syncTimerRef.current !== undefined) {
            window.clearTimeout(syncTimerRef.current)
            syncTimerRef.current = undefined
          }
          onChange?.(nextValue)
          setHighlightValue(nextValue)
        }}
        onChange={(event) => {
          const nextValue = event.target.value
          if (syncTimerRef.current !== undefined) {
            window.clearTimeout(syncTimerRef.current)
          }
          syncTimerRef.current = window.setTimeout(() => {
            onChange?.(nextValue)
            setHighlightValue(nextValue)
            syncTimerRef.current = undefined
          }, 220)
        }}
        onScroll={(event) => {
          if (highlightRef.current) {
            highlightRef.current.scrollTop = event.currentTarget.scrollTop
            highlightRef.current.scrollLeft = event.currentTarget.scrollLeft
          }
        }}
        style={{
          position: 'relative',
          zIndex: 1,
          minHeight: 114,
          padding: '10px 11px',
          border: 0,
          boxShadow: 'none',
          resize: 'vertical',
          background: 'transparent',
          color: 'transparent',
          caretColor: 'var(--ant-color-text)',
          fontFamily: 'ui-monospace, SFMono-Regular, Consolas, monospace',
          fontSize: 13,
          lineHeight: 1.6,
        }}
      />
    </div>
  )
}

const DictTypeSchema = z.object({
  typeCode: z.string('请输入类型编码').min(1, '请输入类型编码'),
  typeName: z.string('请输入类型名称').min(1, '请输入类型名称'),
  isEnabled: z.boolean().default(true),
  sortOrder: z.number().default(0),
  remark: z.string().default(''),
})

const dictTypeDefaults = DictTypeSchema.partial().parse({})

export type DictTypeFormValues = z.infer<typeof DictTypeSchema>

const DictEntrySchema = z.object({
  labelComponent: z.string().max(255, '显示标签组件最多255位').default(''),
  entryLabel: z.string('请输入显示标签').min(1, '请输入显示标签'),
  entryValue: z.string('请输入数据值').min(1, '请输入数据值'),
  languageCode: z.string().default(''),
  sortOrder: z.number().default(0),
  isEnabled: z.boolean().default(true),
  remark: z.string().default(''),
})

const dictEntryDefaults = DictEntrySchema.partial().parse({})

export type DictEntryFormValues = z.infer<typeof DictEntrySchema>

function DictTypePanel({
  selectedType,
  onSelectType,
  onDeleteSelectedType,
  onUpdateSelectedType,
  onBatchCopyEntries,
}: {
  selectedType: DictType | undefined
  onSelectType: (record: DictType) => void
  onDeleteSelectedType: () => void
  onUpdateSelectedType: (record: DictType) => void
  onBatchCopyEntries: (entryIds: number[], targetTypeId: number) => void
}) {
  const [formOpen, setFormOpen] = useState(false)
  const [editing, setEditing] = useState<DictType>()
  const [hoveredDropTypeId, setHoveredDropTypeId] = useState<number | undefined>()
  const [selectedTypeIds, setSelectedTypeIds] = useState<number[]>([])
  const [searchText, setSearchText] = useState('')
  const {
    data,
    total,
    page,
    pageSize,
    loading,
    update,
    send,
  } = usePagination(
    (nextPage, nextPageSize) => {
      const params: Record<string, unknown> = {
        page: nextPage,
        pageSize: nextPageSize,
        orderBy: 'sort_order asc,id desc',
      }
      if (searchText.trim()) {
        params.query = JSON.stringify({
          $or: [
            { typeCode__icontains: searchText.trim() },
            { typeName__icontains: searchText.trim() },
            { remark__icontains: searchText.trim() },
          ],
        })
      }
      return DictApi.typeList(params)
    },
    {
      initialData: {
        total: 0,
        items: [],
      },
      initialPage: 1,
      initialPageSize: DEFAULT_PAGE_SIZE,
      total: response => response.data?.total ?? 0,
      data: response => response.data?.items ?? [],
      watchingStates: [searchText],
      debounce: [500],
    },
  )

  const { form, rules, onFinish } = useZodForm<DictTypeFormValues>({
    schema: DictTypeSchema,
    async onSubmit(values) {
      if (!values) {
        gMessage.error('请填写完整信息')
        return
      }

      const payload = {
        typeCode: values.typeCode,
        typeName: values.typeName,
        isEnabled: values.isEnabled,
        sortOrder: values.sortOrder,
        remark: values.remark ?? '',
      }

      if (editing) {
        await DictApi.typeUpdate({
          id: editing.id,
          ...payload,
        })
        if (selectedType?.id === editing.id) {
          onUpdateSelectedType({
            ...selectedType,
            ...payload,
          })
        }
      }
      else {
        await DictApi.typeCreate(payload)
      }

      gMessage.success('保存成功')
      setEditing(undefined)
      form.resetFields()
      setFormOpen(false)
      await send()
    },
  })

  const openCreate = () => {
    setEditing(undefined)
    form.resetFields()
    form.setFieldsValue(dictTypeDefaults as DictTypeFormValues)
    setFormOpen(true)
  }

  const openEdit = useCallback((record: DictType) => {
    setEditing(record)
    form.setFieldsValue({
      typeCode: record.typeCode,
      typeName: record.typeName,
      isEnabled: record.isEnabled,
      sortOrder: record.sortOrder,
      remark: record.remark,
    })
    setFormOpen(true)
  }, [form])

  const columns = useMemo<ProColumns<DictType>[]>(() => [
    {
      title: '序号',
      dataIndex: 'index',
      width: 60,
      render: (_, __, index) => (page - 1) * pageSize + index + 1,
    },
    {
      title: '类型编码',
      dataIndex: 'typeCode',
      width: 140,
    },
    {
      title: '类型名称',
      dataIndex: 'typeName',
      width: 160,
    },
    {
      title: '状态',
      dataIndex: 'isEnabled',
      width: 90,
      render: (_, record) => statusTag(record.isEnabled),
    },
    {
      title: '排序',
      dataIndex: 'sortOrder',
      width: 80,
    },
    {
      title: '备注',
      dataIndex: 'remark',
      ellipsis: true,
    },
    {
      title: '操作',
      valueType: 'option',
      width: 220,
      fixed: 'right',
      render: (_, record) => {
        const canWrite = record.canWrite !== false
        const canDelete = record.canDelete !== false
        return [
          <Button
            key="edit"
            type="link"
            size="small"
            disabled={!canWrite}
            style={{ padding: 0 }}
            onClick={(event) => {
              event.stopPropagation()
              openEdit(record)
            }}
          >
            编辑
          </Button>,
          <Button
            key="switch"
            type="link"
            size="small"
            disabled={!canWrite}
            style={{ padding: 0 }}
            onClick={async (event) => {
              event.stopPropagation()
              await DictApi.typeUpdate({
                id: record.id,
                isEnabled: !record.isEnabled,
              })
              gMessage.success('操作成功')
              await send()
            }}
          >
            {record.isEnabled ? '停用' : '启用'}
          </Button>,
          <Popconfirm
            key="del"
            title="确认删除该字典类型吗？"
            disabled={!canDelete}
            onConfirm={async (event) => {
              event?.stopPropagation()
              await DictApi.typeDel({ ids: [record.id] })
              gMessage.success('删除成功')
              if (selectedType?.id === record.id) {
                onDeleteSelectedType()
              }
              await send()
            }}
          >
            <Button
              type="link"
              size="small"
              disabled={!canDelete}
              style={{ padding: 0 }}
              onClick={(event) => {
                event.stopPropagation()
              }}
            >
              删除
            </Button>
          </Popconfirm>,
        ]
      },
    },
  ], [openEdit, selectedType?.id, onDeleteSelectedType, send, page, pageSize])

  return (
    <>
      <ProTable<DictType>
        rowKey="id"
        search={false}
        columns={columns}
        dataSource={data}
        loading={loading}
        headerTitle={(
          <Space>
            字典类型
          </Space>
        )}
        pagination={{
          showSizeChanger: true,
          current: page,
          pageSize,
          total,
          onChange: (nextPage, nextPageSize) => {
            update({
              page: nextPage,
              pageSize: nextPageSize,
            })
          },
        }}
        options={{
          reload: () => send(),
        }}
        toolBarRender={() => [
          <Button key="add" type="primary" onClick={openCreate}>
            新增类型
          </Button>,
          <Input.Search
            key="search"
            placeholder="搜索类型编码、名称、备注"
            allowClear
            value={searchText}
            onChange={(e) => {
              setSearchText(e.target.value)
            }}
            onSearch={(value) => {
              setSearchText(value)
            }}
            style={{ width: 280 }}
          />,
          selectedTypeIds.length > 0
            ? (
                <Popconfirm
                  key="batchDel"
                  title={`确认批量删除选中的 ${selectedTypeIds.length} 个字典类型吗？（将同时删除其下的所有字典项）`}
                  onConfirm={async () => {
                    await DictApi.typeDel({ ids: selectedTypeIds })
                    gMessage.success(`成功删除 ${selectedTypeIds.length} 项`)
                    setSelectedTypeIds([])
                    if (selectedType && selectedTypeIds.includes(selectedType.id)) {
                      onDeleteSelectedType()
                    }
                    await send()
                  }}
                >
                  <Button key="batchDel" danger>
                    批量删除
                    {selectedTypeIds.length > 0 && ` (${selectedTypeIds.length})`}
                  </Button>
                </Popconfirm>
              )
            : null,
        ]}
        rowSelection={{
          selectedRowKeys: selectedTypeIds,
          onChange: (keys) => {
            setSelectedTypeIds(keys as number[])
          },
          getCheckboxProps: record => ({
            disabled: record.canDelete === false,
          }),
        }}
        rowClassName={(record) => {
          const classes: string[] = []
          if (record.id === selectedType?.id) {
            classes.push('ant-table-row-selected')
          }
          if (record.id === hoveredDropTypeId && record.id !== selectedType?.id) {
            classes.push('ant-table-row-drop-target')
          }
          return classes.join(' ')
        }}
        onRow={record => ({
          onClick: () => {
            onSelectType(record)
          },
          onDragOver: (e) => {
            e.preventDefault()
            e.dataTransfer.dropEffect = record.canWrite === false ? 'none' : 'copy'
          },
          onDragEnter: (e) => {
            e.preventDefault()
            if (record.canWrite === false) { return }
            setHoveredDropTypeId(record.id)
          },
          onDragLeave: () => {
            setHoveredDropTypeId(undefined)
          },
          onDrop: (e) => {
            e.preventDefault()
            setHoveredDropTypeId(undefined)
            if (record.canWrite === false) { return }
            const raw = e.dataTransfer.getData('text/plain')
            if (!raw) { return }
            try {
              const entryIds: number[] = JSON.parse(raw)
              if (entryIds.length > 0) {
                onBatchCopyEntries(entryIds, record.id)
              }
            }
            catch {
              // ignore parse errors
            }
          },
        })}
      />
      <ModalForm
        title={editing ? '编辑字典类型' : '新增字典类型'}
        open={formOpen}
        onOpenChange={(open) => {
          if (!open) {
            setFormOpen(false)
            setEditing(undefined)
          }
        }}
        width={500}
        form={form}
        onFinish={onFinish}
        submitTimeout={2000}
      >
        <ProFormText required name="typeCode" label="类型编码" rules={rules} placeholder="请输入类型编码" />
        <ProFormText required name="typeName" label="类型名称" rules={rules} placeholder="请输入类型名称" />
        <ProFormDigit name="sortOrder" label="排序" fieldProps={{ precision: 0 }} />
        <ProFormSwitch name="isEnabled" label="启用状态" />
        <ProFormText name="remark" label="备注" placeholder="请输入备注" />
      </ModalForm>
    </>
  )
}

function DictEntryPanel({
  selectedType,
  selectedEntryIds,
  onSelectionChange,
  onClearType,
  refreshKey,
}: {
  selectedType: DictType | undefined
  selectedEntryIds: number[]
  onSelectionChange: (ids: number[]) => void
  onClearType: () => void
  refreshKey: number
}) {
  const [formOpen, setFormOpen] = useState(false)
  const [editing, setEditing] = useState<DictEntry>()
  const [searchText, setSearchText] = useState('')
  const [labelComponentMode, setLabelComponentMode] = useState<string>('')
  const {
    data,
    total,
    page,
    pageSize,
    loading,
    update,
    send,
  } = usePagination(
    (nextPage, nextPageSize) => {
      const conditions: Record<string, unknown>[] = []
      if (selectedType) {
        conditions.push({ sysDictTypeId: String(selectedType.id) })
      }
      if (searchText.trim()) {
        conditions.push({
          $or: [
            { entryLabel__icontains: searchText.trim() },
            { entryValue__icontains: searchText.trim() },
          ],
        })
      }
      return DictApi.entryList({
        page: nextPage,
        pageSize: nextPageSize,
        orderBy: 'sort_order asc,id desc',
        query: conditions.length > 0 ? JSON.stringify({ $and: conditions }) : undefined,
      })
    },
    {
      initialData: {
        total: 0,
        items: [],
      },
      initialPage: 1,
      initialPageSize: DEFAULT_PAGE_SIZE,
      watchingStates: [selectedType?.id, searchText, refreshKey],
      data: response => response.data?.items?.map(item => ({
        ...item,
        labelComponent: item.labelComponent ?? '',
        sortOrder: item.sortOrder ?? 0,
      })) ?? [],
      total: response => response.data?.total ?? 0,
    },
  )

  const { form, rules, onFinish } = useZodForm<DictEntryFormValues>({
    schema: DictEntrySchema,
    async onSubmit(values) {
      if (!values) {
        gMessage.error('请填写完整信息')
        return
      }

      const typeId = editing?.sysDictTypeId ?? selectedType?.id
      if (!typeId) {
        gMessage.error('请先选择字典类型')
        return
      }

      const payload = {
        labelComponent: values.labelComponent ?? '',
        entryLabel: values.entryLabel,
        entryValue: values.entryValue,
        languageCode: values.languageCode,
        sysDictTypeId: typeId,
        sortOrder: values.sortOrder,
        isEnabled: values.isEnabled,
        remark: values.remark ?? '',
      }

      if (editing) {
        await DictApi.entryUpdate({
          id: editing.id,
          ...payload,
        })
      }
      else {
        await DictApi.entryCreate(payload)
      }

      gMessage.success('保存成功')
      setEditing(undefined)
      form.resetFields()
      setFormOpen(false)
      await send()
    },
  })

  const watchedEntryLabel = Form.useWatch('entryLabel', form)
  const watchedLabelComponent = Form.useWatch('labelComponent', form)
  const previewEntryLabel = typeof watchedEntryLabel === 'string' && watchedEntryLabel.trim()
    ? watchedEntryLabel
    : '显示标签'
  const previewLabelComponent = typeof watchedLabelComponent === 'string' ? watchedLabelComponent : ''
  const labelComponentOptions = useMemo(() => [
    ...labelComponentPresets.map(item => ({
      label: (
        <Space>
          <span style={{ display: 'inline-block', minWidth: 56 }}>{item.label}</span>
          {renderDictEntryLabel(item.value, '示例')}
        </Space>
      ),
      value: item.value,
    })),
    {
      label: '自定义',
      value: LABEL_COMPONENT_CUSTOM,
    },
  ], [])

  const openCreate = () => {
    setEditing(undefined)
    setLabelComponentMode('')
    form.resetFields()
    form.setFieldsValue(dictEntryDefaults as DictEntryFormValues)
    setFormOpen(true)
  }

  const openEdit = useCallback((record: DictEntry) => {
    setEditing(record)
    setLabelComponentMode(labelComponentPresetValues.has(record.labelComponent ?? '') ? record.labelComponent ?? '' : LABEL_COMPONENT_CUSTOM)
    form.setFieldsValue({
      ...dictEntryDefaults,
      ...record,
      labelComponent: record.labelComponent ?? '',
    } as DictEntryFormValues)
    setFormOpen(true)
  }, [form])

  const columns = useMemo<ProColumns<DictEntry>[]>(() => [
    {
      title: '序号',
      dataIndex: 'index',
      width: 50,
      render: (_, __, index) => (page - 1) * pageSize + index + 1,
    },
    {
      title: '类型名称',
      dataIndex: ['sysDictType', 'typeName'],
      width: 80,
      render: (_, record) => record.sysDictType?.typeName ?? '-',
    },
    {
      title: '类型编码',
      dataIndex: ['sysDictType', 'typeCode'],
      width: 160,
      render: (_, record) => record.sysDictType?.typeCode ?? '-',
    },
    {
      title: '显示标签',
      dataIndex: 'entryLabel',
      width: 120,
      render: (_, record) => renderDictEntryLabel(record.labelComponent, record.entryLabel),
    },
    {
      title: '数据值',
      dataIndex: 'entryValue',
      width: 120,
    },
    {
      title: '语言',
      dataIndex: 'languageCode',
      width: 100,
    },
    {
      title: '状态',
      dataIndex: 'isEnabled',
      width: 90,
      render: (_, record) => statusTag(record.isEnabled),
    },
    {
      title: '排序',
      dataIndex: 'sortOrder',
      width: 80,
    },
    {
      title: '备注',
      dataIndex: 'remark',
      width: 140,
      ellipsis: true,
    },
    {
      title: '操作',
      valueType: 'option',
      width: 220,
      fixed: 'right',
      render: (_, record) => [
        <a
          key="edit"
          onClick={() => {
            openEdit(record)
          }}
        >
          编辑
        </a>,
        <a
          key="switch"
          onClick={async () => {
            await DictApi.entryUpdate({
              id: record.id,
              isEnabled: !record.isEnabled,
            })
            gMessage.success('操作成功')
            await send()
          }}
        >
          {record.isEnabled ? '停用' : '启用'}
        </a>,
        <Popconfirm
          key="del"
          title="确认删除该字典项吗？"
          onConfirm={async () => {
            await DictApi.entryDel({ ids: [record.id] })
            gMessage.success('删除成功')
            await send()
          }}
        >
          <a>删除</a>
        </Popconfirm>,
      ],
    },
  ], [openEdit, send, page, pageSize])

  return (
    <>
      <ProTable<DictEntry>
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
            update({
              page: nextPage,
              pageSize: nextPageSize,
            })
          },
        }}
        options={{
          reload: () => send(),
        }}
        rowSelection={{
          selectedRowKeys: selectedEntryIds,
          onChange: (keys) => {
            onSelectionChange(keys as number[])
          },
        }}
        onRow={record => ({
          draggable: true,
          onDragStart: (e) => {
            const ids = selectedEntryIds.length > 0 ? selectedEntryIds : [record.id]
            e.dataTransfer.setData('text/plain', JSON.stringify(ids))
            e.dataTransfer.effectAllowed = 'copy'
          },
        })}
        headerTitle={(
          <Space>
            {selectedType
              ? (
                  <>
                    {`字典项 - ${selectedType.typeName}`}
                    <Tag color="error" style={{ cursor: 'pointer' }} onClick={onClearType}>
                      清除筛选
                    </Tag>
                  </>
                )
              : '字典项'}
          </Space>
        )}
        toolBarRender={() => [
          <Space key="tips" size="middle">
            {selectedType
              ? (
                  <>
                    <Tag color="processing">
                      当前编码:
                      {selectedType.typeCode}
                    </Tag>
                  </>
                )
              : null}
            {selectedEntryIds.length > 0
              ? (
                  <Popconfirm
                    key="batchDel"
                    title={`确认批量删除选中的 ${selectedEntryIds.length} 个字典项吗？`}
                    onConfirm={async () => {
                      await DictApi.entryDel({ ids: selectedEntryIds })
                      gMessage.success(`成功删除 ${selectedEntryIds.length} 项`)
                      onSelectionChange([])
                      await send()
                    }}
                  >
                    <Button key="batchDel" danger>
                      批量删除
                      {selectedEntryIds.length > 0 && ` (${selectedEntryIds.length})`}
                    </Button>
                  </Popconfirm>
                )
              : null}
            <Button type="primary" disabled={!selectedType || selectedType.canWrite === false} onClick={openCreate}>
              新增字典项
            </Button>
          </Space>,
          <Input.Search
            key="search"
            placeholder="搜索显示标签、数据值"
            allowClear
            value={searchText}
            onChange={(e) => {
              setSearchText(e.target.value)
            }}
            onSearch={(value) => {
              setSearchText(value)
            }}
            style={{ width: 280 }}
          />,
        ]}
      />
      <ModalForm
        title={editing ? '编辑字典项' : '新增字典项'}
        open={formOpen}
        onOpenChange={(open) => {
          if (!open) {
            setFormOpen(false)
            setEditing(undefined)
            setLabelComponentMode('')
          }
        }}
        form={form}
        onFinish={onFinish}
        width={500}
        submitTimeout={2000}
      >
        <ProFormText required name="entryLabel" label="显示标签" rules={rules} placeholder="请输入显示标签" />
        <Form.Item label="标签样式">
          <Select
            value={labelComponentMode}
            options={labelComponentOptions}
            onChange={(nextMode) => {
              setLabelComponentMode(nextMode)
              if (nextMode === LABEL_COMPONENT_CUSTOM) {
                const current = form.getFieldValue('labelComponent') as string | undefined
                if (!current || labelComponentPresetValues.has(current)) {
                  form.setFieldValue('labelComponent', CUSTOM_LABEL_COMPONENT_DEFAULT)
                }
                return
              }
              form.setFieldValue('labelComponent', nextMode)
            }}
          />
        </Form.Item>
        {labelComponentMode === LABEL_COMPONENT_CUSTOM
          ? (
              <Form.Item
                name="labelComponent"
                label="自定义模板"
                rules={rules}
                extra={(
                  <span>
                    支持
                    {' '}
                    <code>{CUSTOM_LABEL_COMPONENT_DEFAULT.replace('>', ' bordered={false}>')}</code>
                  </span>
                )}
              >
                <ShikiCodeEditorBlock key={`${editing?.id ?? 'create'}:${formOpen}`} />
              </Form.Item>
            )
          : (
              <Form.Item name="labelComponent" hidden>
                <Input />
              </Form.Item>
            )}
        <Form.Item label="效果预览">
          {renderDictEntryLabel(previewLabelComponent, previewEntryLabel)}
        </Form.Item>
        <ProFormText required name="entryValue" label="数据值" rules={rules} placeholder="请输入数据值" />
        <ProFormText name="languageCode" label="语言代码" placeholder="请输入语言代码" />
        <ProFormDigit name="sortOrder" label="排序" fieldProps={{ precision: 0 }} />
        <ProFormSwitch name="isEnabled" label="启用状态" />
        <ProFormText name="remark" label="备注" placeholder="请输入备注" />
      </ModalForm>
    </>
  )
}

function RouteComponent() {
  const [selectedType, setSelectedType] = useState<DictType>()
  const [selectedEntryIds, setSelectedEntryIds] = useState<number[]>([])
  const [refreshKey, setRefreshKey] = useState(0)

  const handleSelectType = useCallback((record: DictType) => {
    setSelectedType(record)
    setSelectedEntryIds([])
  }, [])

  const handleDeleteSelectedType = useCallback(() => {
    setSelectedType(undefined)
    setSelectedEntryIds([])
  }, [])

  const handleUpdateSelectedType = useCallback((record: DictType) => {
    setSelectedType(record)
  }, [])

  const handleBatchCopyEntries = useCallback(async (entryIds: number[], targetTypeId: number) => {
    try {
      await DictApi.entryBatchCopy({ entryIds, targetTypeId })
      gMessage.success(`成功复制 ${entryIds.length} 项`)
      setSelectedEntryIds([])
      setRefreshKey(k => k + 1)
    }
    catch {
      gMessage.error('复制失败')
    }
  }, [])

  return (
    <Splitter>
      <Splitter.Panel defaultSize="40%" min="25%" max="75%">
        <DictTypePanel
          selectedType={selectedType}
          onSelectType={handleSelectType}
          onDeleteSelectedType={handleDeleteSelectedType}
          onUpdateSelectedType={handleUpdateSelectedType}
          onBatchCopyEntries={handleBatchCopyEntries}
        />
      </Splitter.Panel>
      <Splitter.Panel>
        <DictEntryPanel
          selectedType={selectedType}
          selectedEntryIds={selectedEntryIds}
          onSelectionChange={setSelectedEntryIds}
          onClearType={handleDeleteSelectedType}
          refreshKey={refreshKey}
        />
      </Splitter.Panel>
    </Splitter>
  )
}
