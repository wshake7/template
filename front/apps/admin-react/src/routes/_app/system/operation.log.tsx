import type { ProColumns } from '@ant-design/pro-components'
import { ProTable } from '@ant-design/pro-components'
import { createFileRoute } from '@tanstack/react-router'
import { usePagination } from 'alova/client'
import {
  Descriptions,
  Modal,
  Tag,
} from 'antd'
import { useCallback, useState } from 'react'
import API from '~/api/index'

export const Route = createFileRoute('/_app/system/operation/log')({
  staticData: {
    menu: {
      name: '操作日志',
      menuType: 'menu',
    },
  },
  staleTime: 1000 * 60 * 2,
  component: RouteComponent,
})

function successTag(success: boolean) {
  if (success) {
    return <Tag color="success">成功</Tag>
  }
  return <Tag color="error">失败</Tag>
}

function methodTag(method: string) {
  const colorMap: Record<string, string> = {
    GET: 'cyan',
    POST: 'blue',
    PUT: 'orange',
    DELETE: 'red',
    PATCH: 'green',
  }
  return <Tag color={colorMap[method] || 'default'}>{method}</Tag>
}

function costTimeDisplay(costTime: number) {
  if (costTime < 1000) {
    return `${costTime}ms`
  }
  return `${(costTime / 1000).toFixed(2)}s`
}

function RouteComponent() {
  const [detailOpen, setDetailOpen] = useState(false)
  const [detailData, setDetailData] = useState<SysOperationLog | null>(null)

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
      API.Post<Res<PagingResult<SysOperationLog>>>('/api/sys/operation/log/list', {
        page: nextPage,
        pageSize: nextPageSize,
        orderBy: 'id desc',
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
      total: response => response.data?.total ?? 0,
      data: response => response.data?.items ?? [],
    },
  )

  const openDetail = useCallback(async (record: SysOperationLog) => {
    try {
      const res = await OperationLogApi.detail({ id: record.id })
      if (res.data) {
        setDetailData(res.data)
        setDetailOpen(true)
      }
    }
    catch {
      // ignore
    }
  }, [])

  const columns: ProColumns<SysOperationLog>[] = [
    {
      title: 'ID',
      dataIndex: 'id',
      width: 70,
    },
    {
      title: '模块',
      dataIndex: 'module',
      width: 100,
    },
    {
      title: '方法',
      dataIndex: 'method',
      width: 80,
      render: (_, record) => methodTag(record.method),
    },
    {
      title: '请求路径',
      dataIndex: 'path',
      ellipsis: true,
    },
    {
      title: '操作者',
      dataIndex: 'username',
      width: 120,
    },
    {
      title: '客户端IP',
      dataIndex: 'clientIP',
      width: 130,
    },
    {
      title: '状态码',
      dataIndex: 'statusCode',
      width: 80,
    },
    {
      title: '结果',
      dataIndex: 'success',
      width: 70,
      render: (_, record) => successTag(record.success),
    },
    {
      title: '耗时',
      dataIndex: 'costTime',
      width: 80,
      render: (_, record) => costTimeDisplay(record.costTime),
    },
    {
      title: '操作时间',
      dataIndex: 'createdAt',
      width: 180,
    },
    {
      title: '操作',
      valueType: 'option',
      width: 80,
      render: (_, record) => [
        <a
          key="detail"
          onClick={() => {
            openDetail(record)
          }}
        >
          详情
        </a>,
      ],
    },
  ]

  return (
    <>
      <ProTable<SysOperationLog>
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
      />
      <Modal
        title="操作日志详情"
        open={detailOpen}
        onCancel={() => {
          setDetailOpen(false)
        }}
        footer={null}
        width={800}
      >
        {detailData && (
          <Descriptions column={2} bordered size="small">
            <Descriptions.Item label="ID">{detailData.id}</Descriptions.Item>
            <Descriptions.Item label="请求ID">{detailData.requestID}</Descriptions.Item>
            <Descriptions.Item label="模块">{detailData.module}</Descriptions.Item>
            <Descriptions.Item label="方法">{detailData.method}</Descriptions.Item>
            <Descriptions.Item label="请求路径" span={2}>{detailData.path}</Descriptions.Item>
            <Descriptions.Item label="请求URI" span={2}>{detailData.requestURI}</Descriptions.Item>
            <Descriptions.Item label="操作者">{detailData.username}</Descriptions.Item>
            <Descriptions.Item label="客户端IP">{detailData.clientIP}</Descriptions.Item>
            <Descriptions.Item label="状态码">{detailData.statusCode}</Descriptions.Item>
            <Descriptions.Item label="结果">{successTag(detailData.success)}</Descriptions.Item>
            <Descriptions.Item label="耗时">{costTimeDisplay(detailData.costTime)}</Descriptions.Item>
            <Descriptions.Item label="操作时间">{detailData.createdAt}</Descriptions.Item>
            <Descriptions.Item label="失败原因" span={2}>{detailData.reason || '-'}</Descriptions.Item>
            <Descriptions.Item label="地理位置" span={2}>{detailData.location || '-'}</Descriptions.Item>
            <Descriptions.Item label="来源" span={2}>{detailData.referer || '-'}</Descriptions.Item>
            <Descriptions.Item label="浏览器" span={2}>
              {detailData.browserName}
              {' '}
              {detailData.browserVersion}
            </Descriptions.Item>
            <Descriptions.Item label="操作系统" span={2}>
              {detailData.osName}
              {' '}
              {detailData.osVersion}
            </Descriptions.Item>
            <Descriptions.Item label="客户端" span={2}>
              {detailData.clientName}
              {' '}
              (
              {detailData.clientID}
              )
            </Descriptions.Item>
            <Descriptions.Item label="User-Agent" span={2}>{detailData.userAgent || '-'}</Descriptions.Item>
            <Descriptions.Item label="请求体" span={2}>
              <pre style={{ maxHeight: 200, overflow: 'auto', margin: 0, fontSize: 12 }}>{detailData.requestBody || '-'}</pre>
            </Descriptions.Item>
            <Descriptions.Item label="请求头" span={2}>
              <pre style={{ maxHeight: 200, overflow: 'auto', margin: 0, fontSize: 12 }}>{detailData.requestHeader || '-'}</pre>
            </Descriptions.Item>
            <Descriptions.Item label="响应信息" span={2}>
              <pre style={{ maxHeight: 200, overflow: 'auto', margin: 0, fontSize: 12 }}>{detailData.response || '-'}</pre>
            </Descriptions.Item>
          </Descriptions>
        )}
      </Modal>
    </>
  )
}
