import { ProTable } from '@ant-design/pro-components'
import { createFileRoute } from '@tanstack/react-router'
import { Splitter } from 'antd'

export const Route = createFileRoute('/_app/system/dict')({
  staticData: {
    menu: {
      name: '数据字典',
      menuType: 'menu',
    },
  },
  component: RouteComponent,
})

const DictTypeComponent = () => {
  return (
    <>
      <ProTable>

      </ProTable>
    </>
  )
}

const DictEntryComponent = () => {
  return (
    <>
      <ProTable>

      </ProTable>
    </>
  )
}

function RouteComponent() {
  return (
    <Splitter>
      <Splitter.Panel defaultSize="50%" min="20%" max="70%">
        <DictTypeComponent />
      </Splitter.Panel>
      <Splitter.Panel>
        <DictEntryComponent />
      </Splitter.Panel>
    </Splitter>
  )
}
