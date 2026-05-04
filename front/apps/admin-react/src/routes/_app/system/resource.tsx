import type { ProColumns } from '@ant-design/pro-components'
import type { Resource } from '~/api/sysResource'
import { ModalForm, ProFormSelect, ProFormSwitch, ProFormText, ProTable } from '@ant-design/pro-components'
import { createFileRoute } from '@tanstack/react-router'
import { usePagination } from 'alova/client'
import {
  Button,
  Input,
  Popconfirm,
  Space,
  Tag,
} from 'antd'
import { useCallback, useMemo, useState } from 'react'
import z from 'zod'
import { ResourceApi } from '~/api/sysResource'
import { useDictMatch } from '~/hooks/useDictMatch'
import { gMessage } from '~/utils/antd'
import { useZodForm } from '~/utils/zod'

export const Route = createFileRoute('/_app/system/resource')({
  staticData: {
    menu: {
      name: '资源管理',
      menuType: 'menu',
    },
  },
  staleTime: 1000 * 60 * 2,
  component: ResourceManagement,
})

const resourceTypeOptions = [
  { label: 'api', value: 'api' },
  { label: 'data', value: 'data' },
  { label: 'menu', value: 'menu' },
  { label: 'component', value: 'component' },
]

const enabledStatusValue = (isEnabled: boolean) => isEnabled ? '1' : '0'
const fallbackEnabledStatusLabel = (isEnabled: boolean) => isEnabled ? '启用' : '停用'

const ResourceSchema = z.object({
  type: z.string('请选择资源类型').min(1, '请选择资源类型'),
  code: z.string('请输入资源编码').min(1, '请输入资源编码'),
  name: z.string('请输入资源名称').min(1, '请输入资源名称'),
  isEnabled: z.boolean().default(true),
  remark: z.string().default(''),
})

const resourceDefaults = ResourceSchema.partial().parse({})

export type ResourceFormValues = z.infer<typeof ResourceSchema>

function ResourceManagement() {
  const [formOpen, setFormOpen] = useState(false)
  const [editing, setEditing] = useState<Resource>()
  const [selectedIds, setSelectedIds] = useState<number[]>([])
  const [searchText, setSearchText] = useState('')
  const enabledStatus = useDictMatch(DictCode.SYS_ENABLED_STATUS_DICT_CODE)
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
        orderBy: 'id desc',
      }
      if (searchText.trim()) {
        params.query = JSON.stringify({
          $or: [
            { code__icontains: searchText.trim() },
            { name__icontains: searchText.trim() },
            { remark__icontains: searchText.trim() },
          ],
        })
      }
      return ResourceApi.resourceList(params)
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

  const { form, rules, onFinish } = useZodForm<ResourceFormValues>({
    schema: ResourceSchema,
    async onSubmit(values) {
      if (!values) {
        gMessage.error('请填写完整信息')
        return
      }

      const payload = {
        type: values.type,
        code: values.code,
        name: values.name,
        isEnabled: values.isEnabled,
        remark: values.remark ?? '',
      }

      if (editing) {
        await ResourceApi.resourceUpdate({
          id: editing.id,
          ...payload,
        })
      }
      else {
        await ResourceApi.resourceCreate(payload)
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
    form.setFieldsValue(resourceDefaults as ResourceFormValues)
    setFormOpen(true)
  }

  const openEdit = useCallback((record: Resource) => {
    setEditing(record)
    form.setFieldsValue({
      type: record.type,
      code: record.code,
      name: record.name,
      isEnabled: record.isEnabled,
      remark: record.remark,
    })
    setFormOpen(true)
  }, [form])

  const columns = useMemo<ProColumns<Resource>[]>(() => [
    {
      title: '序号',
      dataIndex: 'index',
      width: 60,
      render: (_, __, index) => (page - 1) * pageSize + index + 1,
    },
    {
      title: '资源类型',
      dataIndex: 'type',
      width: 120,
      render: (_, record) => {
        const colorMap: Record<string, string> = {
          api: 'blue',
          data: 'green',
          menu: 'orange',
          component: 'purple',
        }
        return <Tag color={colorMap[record.type] || 'default'}>{record.type}</Tag>
      },
    },
    {
      title: '资源编码',
      dataIndex: 'code',
      width: 200,
    },
    {
      title: '资源名称',
      dataIndex: 'name',
      width: 200,
    },
    {
      title: '状态',
      dataIndex: 'isEnabled',
      width: 90,
      render: (_, record) => enabledStatus.renderLabel(enabledStatusValue(record.isEnabled), fallbackEnabledStatusLabel(record.isEnabled)),
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
            await ResourceApi.resourceUpdate({
              id: record.id,
              isEnabled: !record.isEnabled,
            })
            gMessage.success('操作成功')
            await send()
          }}
        >
          {enabledStatus.getLabel(enabledStatusValue(!record.isEnabled), fallbackEnabledStatusLabel(!record.isEnabled))}
        </a>,
        <Popconfirm
          key="del"
          title="确认删除该资源吗？"
          onConfirm={async (event) => {
            event?.stopPropagation()
            await ResourceApi.resourceDel({ ids: [record.id] })
            gMessage.success('删除成功')
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
  ], [enabledStatus, openEdit, send, page, pageSize])

  return (
    <>
      <ProTable<Resource>
        rowKey="id"
        search={false}
        columns={columns}
        dataSource={data}
        loading={loading}
        headerTitle={(
          <Space>
            资源列表
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
            新增资源
          </Button>,
          <Input.Search
            key="search"
            placeholder="搜索资源编码、名称、备注"
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
          selectedIds.length > 0
            ? (
                <Popconfirm
                  key="batchDel"
                  title={`确认批量删除选中的 ${selectedIds.length} 个资源吗？`}
                  onConfirm={async () => {
                    await ResourceApi.resourceDel({ ids: selectedIds })
                    gMessage.success(`成功删除 ${selectedIds.length} 项`)
                    setSelectedIds([])
                    await send()
                  }}
                >
                  <Button key="batchDel" danger>
                    批量删除
                    {selectedIds.length > 0 && ` (${selectedIds.length})`}
                  </Button>
                </Popconfirm>
              )
            : null,
        ]}
        rowSelection={{
          selectedRowKeys: selectedIds,
          onChange: (keys) => {
            setSelectedIds(keys as number[])
          },
        }}
      />
      <ModalForm
        title={editing ? '编辑资源' : '新增资源'}
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
        <ProFormSelect required name="type" label="资源类型" rules={rules} options={resourceTypeOptions} placeholder="请选择资源类型" />
        <ProFormText required name="code" label="资源编码" rules={rules} placeholder="请输入资源编码" />
        <ProFormText required name="name" label="资源名称" rules={rules} placeholder="请输入资源名称" />
        <ProFormSwitch name="isEnabled" label="启用状态" />
        <ProFormText name="remark" label="备注" placeholder="请输入备注" />
      </ModalForm>
    </>
  )
}
