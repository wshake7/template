import { createAlova } from 'alova'
import { useSSE } from 'alova/client'
import adapterFetch from 'alova/fetch'
import reactHook from 'alova/react'
import { useEffect, useRef } from 'react'
import { XHeader } from '~/domains/http'

const SSE_API = createAlova({
  baseURL: '',
  statesHook: reactHook,
  requestAdapter: adapterFetch(),
  cacheFor: null,
  shareRequest: false,
})

export function useEventStream() {
  const token = useAccountStore(state => state.token)
  const connectedRef = useRef(false)
  const tokenRef = useRef(token)
  tokenRef.current = token

  const eventStream = useSSE(
    () => SSE_API.Get('/api/events', {
      headers: {
        [XHeader.Token]: useAccountStore.getState().token,
      },
    }),
    {
      immediate: false,
      withCredentials: true,
      reconnectionTime: 10000,
      interceptByGlobalResponded: false,
      responseType: 'json',
    },
  )

  const eventBoundRef = useRef(false)
  if (!eventBoundRef.current) {
    eventBoundRef.current = true
    eventStream
      .on<{ count: number }>('count', (event) => {
        console.log('[SSE count]', event.data.count)
      })
      .onError((event) => {
        console.error('[SSE error]', event.error)
        if (tokenRef.current === '') {
          eventStream.close()
          connectedRef.current = false
        }
      })
  }

  useEffect(() => {
    if (token === '') {
      eventStream.close()
      connectedRef.current = false
      return
    }
    if (connectedRef.current) {
      return
    }

    connectedRef.current = true
    eventStream.send().catch((error) => {
      connectedRef.current = false
      console.error('[SSE connect error]', error)
    })

    return () => {
      eventStream.close()
      connectedRef.current = false
    }
  }, [eventStream, token])
}
