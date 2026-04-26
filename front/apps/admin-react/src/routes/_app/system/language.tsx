import type { ProColumns } from '@ant-design/pro-components'
import type { RecordKey } from '@ant-design/pro-components/es/utils/useEditableArray'
import type { LanguageEntry, LanguageType } from '~/api/language'
import { EditableProTable, ModalForm, ProFormDigit, ProFormSwitch, ProFormText, ProTable } from '@ant-design/pro-components'
import { createFileRoute } from '@tanstack/react-router'
import { usePagination } from 'alova/client'
import {
  Button,
  Form,
  Input,
  InputNumber,
  Modal,
  Popconfirm,
  Space,
  Splitter,
  Switch,
  Tag,
} from 'antd'
import { useCallback, useMemo, useState } from 'react'
import z from 'zod'
import { LangApi } from '~/api/language'
import { gMessage } from '~/utils/antd'
import { useZodForm } from '~/utils/zod'

export const Route = createFileRoute('/_app/system/language')({
  staticData: {
    menu: {
      name: '语言管理',
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

const LangTypeSchema = z.object({
  typeCode: z.string('请输入语言编码').min(1, '请输入语言编码'),
  typeName: z.string('请输入语言名称').min(1, '请输入语言名称'),
  isDefault: z.boolean().default(false),
  isEnabled: z.boolean().default(true),
  sortOrder: z.number().default(0),
})

const langTypeDefaults = LangTypeSchema.partial().parse({})

export type LangTypeFormValues = z.infer<typeof LangTypeSchema>

function LanguageTypePanel({
  selectedType,
  onSelectType,
  onDeleteSelectedType,
  onUpdateSelectedType,
}: {
  selectedType: LanguageType | undefined
  onSelectType: (record: LanguageType) => void
  onDeleteSelectedType: () => void
  onUpdateSelectedType: (record: LanguageType) => void
}) {
  const [formOpen, setFormOpen] = useState(false)
  const [editing, setEditing] = useState<LanguageType>()
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
          ],
        })
      }
      return LangApi.typeList(params)
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

  const { form, rules, onFinish } = useZodForm<LangTypeFormValues>({
    schema: LangTypeSchema,
    async onSubmit(values) {
      if (!values) {
        gMessage.error('请填写完整信息')
        return
      }

      const payload = {
        typeCode: values.typeCode,
        typeName: values.typeName,
        isDefault: values.isDefault,
        isEnabled: values.isEnabled,
        sortOrder: values.sortOrder,
      }

      if (editing) {
        await LangApi.typeUpdate({
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
        await LangApi.typeCreate(payload)
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
    form.setFieldsValue(langTypeDefaults as LangTypeFormValues)
    setFormOpen(true)
  }

  const openEdit = useCallback((record: LanguageType) => {
    setEditing(record)
    form.setFieldsValue({
      typeCode: record.typeCode,
      typeName: record.typeName,
      isDefault: record.isDefault,
      isEnabled: record.isEnabled,
      sortOrder: record.sortOrder,
    })
    setFormOpen(true)
  }, [form])

  const columns = useMemo<ProColumns<LanguageType>[]>(() => [
    {
      title: '序号',
      dataIndex: 'index',
      width: 60,
      render: (_, __, index) => (page - 1) * pageSize + index + 1,
    },
    {
      title: '类型编码',
      dataIndex: 'typeCode',
      width: 100,
    },
    {
      title: '语言名称',
      dataIndex: 'typeName',
      width: 140,
    },
    {
      title: '默认',
      dataIndex: 'isDefault',
      width: 70,
      render: (_, record) => (record.isDefault ? <Tag color="processing">默认</Tag> : '-'),
    },
    {
      title: '状态',
      dataIndex: 'isEnabled',
      width: 80,
      render: (_, record) => statusTag(record.isEnabled),
    },
    {
      title: '排序',
      dataIndex: 'sortOrder',
      width: 70,
    },
    {
      title: '操作',
      valueType: 'option',
      width: 240,
      render: (_: unknown, record: LanguageType) => [
        !record.isDefault
          ? (
              <a
                key="setDefault"
                onClick={async (event) => {
                  event.stopPropagation()
                  await LangApi.typeUpdate({ id: record.id, isDefault: true })
                  gMessage.success('设置成功')
                  if (selectedType?.id === record.id) {
                    onUpdateSelectedType({ ...selectedType, isDefault: true })
                  }
                  await send()
                }}
              >
                设为默认
              </a>
            )
          : null,
        record.isDefault
          ? null
          : (
              <a
                key="switch"
                onClick={async (event) => {
                  event.stopPropagation()
                  await LangApi.typeUpdate({ id: record.id, isEnabled: !record.isEnabled })
                  gMessage.success('操作成功')
                  if (selectedType?.id === record.id) {
                    onUpdateSelectedType({ ...selectedType, isEnabled: !record.isEnabled })
                  }
                  await send()
                }}
              >
                {record.isEnabled ? '停用' : '启用'}
              </a>
            ),
        <a
          key="edit"
          onClick={(event) => {
            event.stopPropagation()
            openEdit(record)
          }}
        >
          编辑
        </a>,
        record.isDefault
          ? null
          : (
              <Popconfirm
                key="del"
                title="确认删除该语言类型吗？"
                onConfirm={async (event) => {
                  event?.stopPropagation()
                  await LangApi.typeDel({ ids: [record.id] })
                  gMessage.success('删除成功')
                  if (selectedType?.id === record.id) {
                    onDeleteSelectedType()
                  }
                  await send()
                }}
              >
                <a
                  onClick={(event) => {
                    event.stopPropagation()
                  }}
                >
                  删除
                </a>
              </Popconfirm>
            ),
      ],
    },
  ], [openEdit, selectedType?.id, onDeleteSelectedType, send, page, pageSize])

  return (
    <>
      <ProTable<LanguageType>
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
        toolBarRender={() => [
          <Input.Search
            key="search"
            placeholder="搜索语言编码、名称"
            allowClear
            value={searchText}
            onChange={(e) => {
              setSearchText(e.target.value)
            }}
            onSearch={(value) => {
              setSearchText(value)
            }}
            style={{ width: 260 }}
          />,
          <Button key="add" type="primary" onClick={openCreate}>
            新增语言
          </Button>,
        ]}
        rowClassName={(record) => {
          if (record.id === selectedType?.id) {
            return 'ant-table-row-selected'
          }
          return ''
        }}
        onRow={record => ({
          onClick: () => {
            onSelectType(record)
          },
        })}
      />
      <ModalForm
        title={editing ? '编辑语言类型' : '新增语言类型'}
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
        <ProFormText required name="typeCode" label="语言编码" rules={rules} placeholder="如: zh-CN" />
        <ProFormText required name="typeName" label="语言名称" rules={rules} placeholder="如: 简体中文" />
        <ProFormDigit name="sortOrder" label="排序" fieldProps={{ precision: 0 }} />
        <ProFormSwitch name="isDefault" label="是否默认" />
        <ProFormSwitch name="isEnabled" label="启用状态" />
      </ModalForm>
    </>
  )
}

function LanguageEntryPanel({
  selectedType,
  refreshKey,
}: {
  selectedType: LanguageType | undefined
  refreshKey: number
}) {
  const [editableRowKeys, setEditableRowKeys] = useState<React.Key[]>([])
  const [searchText, setSearchText] = useState('')
  const [editModalOpen, setEditModalOpen] = useState(false)
  const [editModalData, setEditModalData] = useState<LanguageEntry[]>([])
  const [editModalForm] = Form.useForm()
  const [batchCreateOpen, setBatchCreateOpen] = useState(false)
  const [batchCreateForm] = Form.useForm()
  const [allTypes, setAllTypes] = useState<LanguageType[]>([])

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
        conditions.push({ sysLanguageTypeId: String(selectedType.id) })
      }
      if (searchText.trim()) {
        conditions.push({
          $or: [
            { entryCode__icontains: searchText.trim() },
            { entryValue__icontains: searchText.trim() },
          ],
        })
      }
      return LangApi.entryList({
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
      data: response => response.data?.items ?? [],
      total: response => response.data?.total ?? 0,
    },
  )

  const loadAllTypes = useCallback(async () => {
    const res = await LangApi.typeList({ noPaging: true, pageSize: 1000 }).send()
    if (res.data?.items) {
      setAllTypes(res.data.items)
    }
  }, [])

  const openEditModal = useCallback(async (record: LanguageEntry) => {
    await loadAllTypes()
    const res = await LangApi.entryList({ noPaging: true, pageSize: 1000, query: JSON.stringify({ entryCode: record.entryCode }) }).send()
    const entries = res.data?.items ?? []
    setEditModalData(entries)
    const fields: Record<string, string> = {}
    for (const entry of entries) {
      fields[entry.sysLanguageType?.typeCode ?? ''] = entry.entryValue
    }
    editModalForm.setFieldsValue(fields)
    setEditModalOpen(true)
  }, [loadAllTypes, editModalForm])

  const handleEditModalSave = useCallback(async () => {
    try {
      const values = await editModalForm.validateFields()
      const updatePromises = editModalData.map(async (entry) => {
        const newValue = values[entry.sysLanguageType?.typeCode ?? ''] ?? ''
        if (newValue !== entry.entryValue) {
          await LangApi.entryUpdate({ id: entry.id, entryValue: newValue })
        }
      })
      await Promise.all(updatePromises)
      gMessage.success('保存成功')
      setEditModalOpen(false)
      await send()
    }
    catch {
      // validation error
    }
  }, [editModalForm, editModalData, send])

  const handleSave = useCallback(async (_key: RecordKey, record: LanguageEntry, _: LanguageEntry) => {
    await LangApi.entryUpdate({
      id: record.id,
      entryValue: record.entryValue,
      sortOrder: record.sortOrder,
      isEnabled: record.isEnabled,
      remark: record.remark,
    })
    await send()
  }, [send])

  const handleDelete = useCallback(async (_key: RecordKey, record: LanguageEntry) => {
    await LangApi.entryDel({ ids: [record.id] })
    gMessage.success('删除成功')
    await send()
  }, [send])

  const openBatchCreate = useCallback(async () => {
    await loadAllTypes()
    batchCreateForm.resetFields()
    batchCreateForm.setFieldsValue({ entryCode: '', sortOrder: 0, isEnabled: true })
    setBatchCreateOpen(true)
  }, [loadAllTypes, batchCreateForm])

  const handleBatchCreateSave = useCallback(async () => {
    try {
      const values = await batchCreateForm.validateFields()
      const valuesMap: Record<string, string> = {}
      for (const type of allTypes) {
        valuesMap[type.typeCode] = values[type.typeCode] ?? ''
      }
      await LangApi.entryBatchCreate({
        entryCode: values.entryCode,
        values: valuesMap,
        sortOrder: values.sortOrder ?? 0,
        isEnabled: values.isEnabled ?? true,
      })
      gMessage.success('创建成功')
      setBatchCreateOpen(false)
      await send()
    }
    catch {
      // validation error
    }
  }, [batchCreateForm, allTypes, send])

  const columns = useMemo<ProColumns<LanguageEntry>[]>(() => [
    {
      title: '序号',
      dataIndex: 'index',
      width: 55,
      editable: false,
      render: (_: unknown, _record: LanguageEntry, index: number) => (page - 1) * pageSize + index + 1,
    },
    {
      title: 'Code',
      dataIndex: 'entryCode',
      width: 160,
      ellipsis: true,
      editable: false,
    },
    {
      title: '值',
      dataIndex: 'entryValue',
      fieldProps: {
        placeholder: '请输入值',
      },
      render: (_: unknown, record: LanguageEntry) =>
        record.entryValue || <span className="text-gray-300">(空)</span>,
    },
    {
      title: '排序',
      dataIndex: 'sortOrder',
      width: 65,
      valueType: 'digit',
      editable: false,
    },
    {
      title: '状态',
      dataIndex: 'isEnabled',
      width: 75,
      valueType: 'switch',
      editable: false,
      render: (_: unknown, record: LanguageEntry) => statusTag(record.isEnabled),
    },
    {
      title: '备注',
      dataIndex: 'remark',
      ellipsis: true,
      editable: false,
    },
    {
      title: '操作',
      valueType: 'option',
      width: 180,
      render: (_, record: LanguageEntry, _index) => [
        <Popconfirm
          key="del"
          title="确认删除该语言条目吗？"
          onConfirm={async () => {
            await LangApi.entryDel({ ids: [record.id] })
            gMessage.success('删除成功')
            await send()
          }}
        >
          <a onClick={e => e.stopPropagation()}>删除</a>
        </Popconfirm>,
        <a
          key="modalEdit"
          onClick={(e: React.MouseEvent) => {
            e.stopPropagation()
            openEditModal(record)
          }}
        >
          跨语言编辑
        </a>,
      ],
    },
  ], [page, pageSize, send, openEditModal])

  return (
    <>
      <EditableProTable<LanguageEntry>
        rowKey="id"
        search={false}
        columns={columns}
        value={data}
        loading={loading}
        onRow={record => ({
          onClick: () => {
            // 如果已经在编辑其他行，先不切换（或者可以直接切换）
            if (!editableRowKeys.includes(record.id)) {
              setEditableRowKeys([record.id])
            }
          },
          style: { cursor: 'pointer' },
        })}
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
        editable={{
          type: 'single',
          editableKeys: editableRowKeys,
          onChange: setEditableRowKeys,
          onSave: handleSave,
          onDelete: handleDelete,
          deletePopconfirmMessage: '确认删除该语言条目吗？',
          actionRender: (row, _config, doms) => [
            doms.save,
            doms.cancel,
          ],
        }}
        headerTitle={(
          <Space>
            {selectedType
              ? (
                  <>
                    {`语言条目 - ${selectedType.typeName}`}
                    <Tag color="processing">{selectedType.typeCode}</Tag>
                  </>
                )
              : '语言条目'}
          </Space>
        )}
        toolBarRender={() => [
          <Button key="batchAdd" type="primary" onClick={openBatchCreate}>
            新增
          </Button>,
          <Input.Search
            key="search"
            placeholder="搜索 Code、值"
            allowClear
            value={searchText}
            onChange={(e) => {
              setSearchText(e.target.value)
            }}
            onSearch={(value) => {
              setSearchText(value)
            }}
            style={{ width: 260 }}
          />,
        ]}
      />

      {/* Edit Modal - shows entry values across all language types */}
      <Modal
        title={`编辑条目 - ${editModalData[0]?.entryCode ?? ''}`}
        open={editModalOpen}
        onCancel={() => setEditModalOpen(false)}
        onOk={handleEditModalSave}
        width={600}
      >
        <Form form={editModalForm} layout="vertical">
          {allTypes.map(type => (
            <Form.Item
              key={type.typeCode}
              name={type.typeCode}
              label={`${type.typeName} (${type.typeCode})`}
            >
              <Input
                placeholder={`请输入${type.typeName}翻译`}
                disabled={!type.isEnabled}
              />
            </Form.Item>
          ))}
        </Form>
      </Modal>

      {/* Batch Create Modal */}
      <Modal
        title="新增条目"
        open={batchCreateOpen}
        onCancel={() => setBatchCreateOpen(false)}
        onOk={handleBatchCreateSave}
        width={600}
      >
        <Form form={batchCreateForm} layout="vertical">
          <Form.Item
            name="entryCode"
            label="条目编码 (Code)"
            rules={[{ required: true, message: '请输入条目编码' }]}
          >
            <Input placeholder="如: common.submit" />
          </Form.Item>
          <Form.Item name="sortOrder" label="排序" initialValue={0}>
            <InputNumber style={{ width: '100%' }} precision={0} />
          </Form.Item>
          <Form.Item name="isEnabled" label="启用状态" initialValue={true} valuePropName="checked">
            <Switch />
          </Form.Item>
          <div className="text-gray-500 mb-2 text-sm">各语言翻译值（不填默认为空字符串）</div>
          {allTypes.map(type => (
            <Form.Item
              key={type.typeCode}
              name={type.typeCode}
              label={`${type.typeName} (${type.typeCode})`}
            >
              <Input
                placeholder={type.isDefault ? '默认语言，建议填写' : '可选'}
                disabled={!type.isEnabled}
              />
            </Form.Item>
          ))}
        </Form>
      </Modal>

    </>
  )
}

function RouteComponent() {
  const [selectedType, setSelectedType] = useState<LanguageType>()
  const [refreshKey] = useState(0)

  const handleSelectType = useCallback((record: LanguageType) => {
    setSelectedType(record)
  }, [])

  const handleDeleteSelectedType = useCallback(() => {
    setSelectedType(undefined)
  }, [])

  const handleUpdateSelectedType = useCallback((record: LanguageType) => {
    setSelectedType(record)
  }, [])

  return (
    <Splitter>
      <Splitter.Panel defaultSize="40%" min="25%" max="60%">
        <LanguageTypePanel
          selectedType={selectedType}
          onSelectType={handleSelectType}
          onDeleteSelectedType={handleDeleteSelectedType}
          onUpdateSelectedType={handleUpdateSelectedType}
        />
      </Splitter.Panel>
      <Splitter.Panel>
        <LanguageEntryPanel
          selectedType={selectedType}
          refreshKey={refreshKey}
        />
      </Splitter.Panel>
    </Splitter>
  )
}
