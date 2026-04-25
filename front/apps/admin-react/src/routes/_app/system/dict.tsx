import type { ProColumns } from '@ant-design/pro-components'
import type { DictEntry, DictType } from '~/api/dict'
import { ModalForm, ProFormDigit, ProFormSwitch, ProFormText, ProTable } from '@ant-design/pro-components'
import { createFileRoute } from '@tanstack/react-router'
import { usePagination } from 'alova/client'
import {
  Button,
  Popconfirm,
  Space,
  Splitter,
  Tag,
} from 'antd'
import { useCallback, useMemo, useState } from 'react'
import z from 'zod'
import { DictApi } from '~/api/dict'
import API from '~/api/index'
import { gMessage } from '~/utils/antd'
import { useZodForm } from '~/utils/zod'

export const Route = createFileRoute('/_app/system/dict')({
  staticData: {
    menu: {
      name: '数据字典',
      menuType: 'menu',
    },
  },
  component: RouteComponent,
})

function statusTag(isEnabled: boolean) {
  if (isEnabled) {
    return <Tag color="success">启用</Tag>
  }
  return <Tag color="default">停用</Tag>
}

const DictTypeSchema = z.object({
  typeCode: z.string('请输入类型编码').min(1, '请输入类型编码'),
  typeName: z.string('请输入类型名称').min(1, '请输入类型名称'),
  isEnabled: z.boolean(),
  sortOrder: z.number(),
  description: z.string(),
})

export type DictTypeFormValues = z.infer<typeof DictTypeSchema>

const DictEntrySchema = z.object({
  entryLabel: z.string('请输入显示标签').min(1, '请输入显示标签'),
  entryValue: z.string('请输入数据值').min(1, '请输入数据值'),
  numericValue: z.number(),
  languageCode: z.string(),
  sortOrder: z.number(),
  isEnabled: z.boolean(),
  remark: z.string(),
})

export type DictEntryFormValues = z.infer<typeof DictEntrySchema>

function DictTypePanel({
  selectedType,
  onSelectType,
  onDeleteSelectedType,
  onUpdateSelectedType,
}: {
  selectedType: DictType | undefined
  onSelectType: (record: DictType) => void
  onDeleteSelectedType: () => void
  onUpdateSelectedType: (record: DictType) => void
}) {
  const [formOpen, setFormOpen] = useState(false)
  const [editing, setEditing] = useState<DictType>()
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
      API.Post<Res<PagingResult<DictType>>>('/api/dict/type/list', {
        page: nextPage,
        pageSize: nextPageSize,
        orderBy: 'sort_order asc,id desc',
      }, {
        cacheFor: 5 * 60 * 1000,
      }),
    {
      initialData: {
        total: 0,
        items: [],
      },
      initialPage: 1,
      initialPageSize: 10,
      total: response => response.data?.total ?? 0,
      data: response => response.data?.items ?? [],
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
        description: values.description,
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
      setFormOpen(false)
      setEditing(undefined)
      form.resetFields()
      await send()
    },
  })

  const openCreate = () => {
    setEditing(undefined)
    form.resetFields()
    setFormOpen(true)
  }

  const openEdit = useCallback((record: DictType) => {
    setEditing(record)
    form.setFieldsValue({
      typeCode: record.typeCode,
      typeName: record.typeName,
      isEnabled: record.isEnabled,
      sortOrder: record.sortOrder,
      description: record.description,
    })
    setFormOpen(true)
  }, [form])

  const columns = useMemo<ProColumns<DictType>[]>(() => [
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
      title: '描述',
      dataIndex: 'description',
      ellipsis: true,
    },
    {
      title: '操作',
      valueType: 'option',
      width: 220,
      render: (_, record) => [
        <a
          key="edit"
          onClick={(event) => {
            event.stopPropagation()
            openEdit(record)
          }}
        >
          编辑
        </a>,
        <a
          key="switch"
          onClick={async (event) => {
            event.stopPropagation()
            await DictApi.typeSwitch({
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
          title="确认删除该字典类型吗？"
          onConfirm={async (event) => {
            event?.stopPropagation()
            await DictApi.typeDel({ id: record.id })
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
        </Popconfirm>,
      ],
    },
  ], [openEdit, selectedType?.id, onDeleteSelectedType, send])

  return (
    <>
      <ProTable<DictType>
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
        toolBarRender={() => [
          <Button key="add" type="primary" onClick={openCreate}>
            新增类型
          </Button>,
        ]}
        rowClassName={record => (record.id === selectedType?.id ? 'ant-table-row-selected' : '')}
        onRow={record => ({
          onClick: () => {
            onSelectType(record)
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
        <ProFormText name="description" label="描述" placeholder="请输入描述" />
      </ModalForm>
    </>
  )
}

function DictEntryPanel({
  selectedType,
}: {
  selectedType: DictType | undefined
}) {
  const [formOpen, setFormOpen] = useState(false)
  const [editing, setEditing] = useState<DictEntry>()
  const {
    data,
    total,
    page,
    pageSize,
    loading,
    update,
    send,
  } = usePagination(
    () =>
      API.Post<Res<PagingResult<DictEntry>>>('/api/dict/entry/list', {
        noPaging: true,
        orderBy: 'sort_order asc,id desc',
      }, {
        cacheFor: 0,
      }),
    {
      initialData: {
        total: 0,
        items: [],
      },
      initialPage: 1,
      initialPageSize: 10,
      immediate: false,
      watchingStates: [selectedType?.id],
      data: (response) => {
        const allItems = response.data?.items ?? []
        if (!selectedType) {
          return []
        }
        return allItems.filter(item => item.sysDictTypeId === selectedType.id)
      },
      total: (response) => {
        const allItems = response.data?.items ?? []
        if (!selectedType) {
          return 0
        }
        return allItems.filter(item => item.sysDictTypeId === selectedType.id).length
      },
    },
  )
  const pagedData = useMemo(() => {
    const start = (page - 1) * pageSize
    return data.slice(start, start + pageSize)
  }, [data, page, pageSize])

  const { form, rules, onFinish } = useZodForm<DictEntryFormValues>({
    schema: DictEntrySchema,
    async onSubmit(values) {
      if (!values) {
        gMessage.error('请填写完整信息')
        return
      }

      if (!selectedType) {
        gMessage.error('请先选择字典类型')
        return
      }

      const payload = {
        entryLabel: values.entryLabel,
        entryValue: values.entryValue,
        numericValue: values.numericValue,
        languageCode: values.languageCode,
        sysDictTypeId: selectedType.id,
        sortOrder: values.sortOrder,
        isEnabled: values.isEnabled,
        remark: values.remark,
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
      setFormOpen(false)
      setEditing(undefined)
      form.resetFields()
      await send()
    },
  })

  const openCreate = () => {
    setEditing(undefined)
    form.resetFields()
    setFormOpen(true)
  }

  const openEdit = useCallback((record: DictEntry) => {
    setEditing(record)
    form.setFieldsValue({
      entryLabel: record.entryLabel,
      entryValue: record.entryValue,
      numericValue: record.numericValue,
      languageCode: record.languageCode,
      sortOrder: record.sortOrder,
      isEnabled: record.isEnabled,
      remark: record.remark,
    })
    setFormOpen(true)
  }, [form])

  const columns = useMemo<ProColumns<DictEntry>[]>(() => [
    {
      title: '显示标签',
      dataIndex: 'entryLabel',
      width: 160,
    },
    {
      title: '数据值',
      dataIndex: 'entryValue',
      width: 160,
    },
    {
      title: '数值',
      dataIndex: 'numericValue',
      width: 90,
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
      ellipsis: true,
    },
    {
      title: '操作',
      valueType: 'option',
      width: 220,
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
            await DictApi.entrySwitch({
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
            await DictApi.entryDel({ id: record.id })
            gMessage.success('删除成功')
            await send()
          }}
        >
          <a>删除</a>
        </Popconfirm>,
      ],
    },
  ], [openEdit, send])

  return (
    <>
      <ProTable<DictEntry>
        rowKey="id"
        search={false}
        columns={columns}
        dataSource={pagedData}
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
        headerTitle={selectedType ? `字典项 - ${selectedType.typeName}` : '字典项'}
        toolBarRender={() => [
          <Space key="tips" size="middle">
            {selectedType
              ? (
                  <Tag color="processing">
                    当前编码:
                    {selectedType.typeCode}
                  </Tag>
                )
              : <Tag>请先点击左侧字典类型</Tag>}
            <Button type="primary" disabled={!selectedType} onClick={openCreate}>
              新增字典项
            </Button>
          </Space>,
        ]}
      />
      <ModalForm
        title={editing ? '编辑字典项' : '新增字典项'}
        open={formOpen}
        onOpenChange={(open) => {
          if (!open) {
            setFormOpen(false)
            setEditing(undefined)
          }
        }}
        form={form}
        onFinish={onFinish}
        width={500}
        submitTimeout={2000}
      >
        <ProFormText required name="entryLabel" label="显示标签" rules={rules} placeholder="请输入显示标签" />
        <ProFormText required name="entryValue" label="数据值" rules={rules} placeholder="请输入数据值" />
        <ProFormDigit name="numericValue" label="数值" fieldProps={{ precision: 0 }} />
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

  const handleSelectType = useCallback((record: DictType) => {
    setSelectedType(record)
  }, [])

  const handleDeleteSelectedType = useCallback(() => {
    setSelectedType(undefined)
  }, [])

  const handleUpdateSelectedType = useCallback((record: DictType) => {
    setSelectedType(record)
  }, [])

  return (
    <Splitter>
      <Splitter.Panel defaultSize="50%" min="25%" max="75%">
        <DictTypePanel
          selectedType={selectedType}
          onSelectType={handleSelectType}
          onDeleteSelectedType={handleDeleteSelectedType}
          onUpdateSelectedType={handleUpdateSelectedType}
        />
      </Splitter.Panel>
      <Splitter.Panel>
        <DictEntryPanel
          selectedType={selectedType}
        />
      </Splitter.Panel>
    </Splitter>
  )
}
