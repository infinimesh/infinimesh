import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, test } from 'vitest'
import { flushPromises } from '@vue/test-utils'

import { Access, Nodes } from 'infinimesh-proto/build/es/node/access/access_pb'
import { JoinRequest } from 'infinimesh-proto/build/es/node/node_pb'
import { Namespace, Namespaces, Plugin } from 'infinimesh-proto/build/es/node/namespaces/namespaces_pb'
import { Accounts } from 'infinimesh-proto/build/es/node/accounts/accounts_pb'
import { useNSStore } from '@/store/namespaces.js'

describe('namespaces store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  test.concurrent('fetch namespaces', async () => {
    const store = useNSStore()

    await store.fetchNamespaces(true)
    expect({ namespaces: Object.values(store.namespaces) }).toEqual(
      new Namespaces({ namespaces: Object.values(store.namespaces) })
    )
  })

  test.concurrent('fetch deletables', async () => {
    const store = useNSStore()

    await store.fetchNamespaces(true)
    const { uuid } = await store.namespaces_list[0]
    const { nodes } = await store.deletables(uuid)

    expect({nodes}).toEqual(new Nodes({ nodes }))
  })

  test('update namespace', async () => {
    const store = useNSStore()

    await store.fetchNamespaces(true)
    const namespace = store.namespaces_list[0]

    await store.update({ ...namespace, title: 'test title' })
    await store.fetchNamespaces(true)

    expect(store.namespaces[namespace.uuid].title).toBe('test title')
  })

  test('update namespace config', async () => {
    const store = useNSStore()

    await store.fetchNamespaces(true)
    const namespace = store.namespaces_list[0]

    await store.update({ ...namespace, config: { string: 'value' } })
    await store.fetchNamespaces(true)

    expect(store.namespaces[namespace.uuid].config.toJSON())
      .toEqual({ string: 'value' })
  })

  test('create namespace', async () => {
    const store = useNSStore()

    const namespace = '11'
    const uuid = Math.random().toString(16).slice(2)
    const role = Math.floor(Math.random() * 2 + 1)
    const level = Math.floor(Math.random() * 4 + 1)

    const namespaceItem = new Namespace({
      uuid: namespace,
      title: 'New Namespace',
      access: new Access({ namespace, level, role }),
      plugin: new Plugin({ uuid, vars: {} }),
      config: {}
    })

    await store.create(namespaceItem)
    await store.fetchNamespaces(true)
    await flushPromises()

    expect(store.namespaces[namespaceItem.uuid]).toEqual(new Namespace({
      uuid: namespace,
      title: 'New Namespace',
      access: new Access({ namespace, level, role }),
      plugin: new Plugin({ uuid, vars: {} }),
      config: {}
    }))
  })

  test('delete namespace', async () => {
    const store = useNSStore()

    await store.fetchNamespaces(true)
    const namespace = store.namespaces_list[0]

    await store.remove(namespace.uuid)
    await store.fetchNamespaces(true)
    await flushPromises()

    expect(store.namespaces[namespace.uuid]).toBeUndefined()
  })

  test('join namespace', async () => {
    const store = useNSStore()

    await store.fetchNamespaces(true)
    const namespace = store.namespaces_list[0]
    const { accounts } = await store.join(new JoinRequest({
      access: namespace.access.level, namespace: namespace.uuid
    }))

    expect({ accounts }).toEqual(new Accounts({ accounts }))
  })

  test('fetch joins', async () => {
    const store = useNSStore()

    await store.fetchNamespaces(true)
    const { uuid } = store.namespaces_list[0]
    const { accounts } = await store.loadJoins(uuid)

    expect({ accounts }).toEqual(new Accounts({ accounts }))
  })
})
