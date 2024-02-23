import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, test, vi } from 'vitest'
import { flushPromises } from '@vue/test-utils'
import { Access, Node, Nodes } from 'infinimesh-proto/build/es/node/access/access_pb'
import { Device, Devices } from 'infinimesh-proto/build/es/node/devices/devices_pb'
import { JoinGeneralRequest, JoinRequest } from 'infinimesh-proto/build/es/node/node_pb'
import { useDevicesStore } from '@/store/devices.js'

describe('devices store', () => {
  const mockBar = { start: vi.fn(), finish: vi.fn(), error: vi.fn() }

  beforeEach(() => {
    setActivePinia(createPinia())
  })

  test.concurrent('fetch devices', async () => {
    const store = useDevicesStore()

    await store.fetchDevices(false, true)
    expect({ devices: Object.values(store.devices), total: 0 }).toEqual(
      new Devices({ devices: Object.values(store.devices), total: 0 })
    )
  })

  test.concurrent('make devices token with read level', async () => {
    const store = useDevicesStore()

    await store.fetchDevices(false, true)
    const token = await store.makeDevicesToken(Object.keys(store.devices))

    expect(token.length).toBeGreaterThan(0)
  })

  test.concurrent('make devices token with mgmt level', async () => {
    const store = useDevicesStore()

    await store.fetchDevices(false, true)
    const token = await store.makeDevicesToken(Object.keys(store.devices), true)

    expect(token).not.toHaveLength(0)
  })

  test.concurrent('get devices state', async () => {
    const store = useDevicesStore()

    await store.fetchDevices(false, true)
    await store.getDevicesState(Object.keys(store.devices))

    for (const item of store.reported.values()) {
      expect(item).toBeInstanceOf(Array)
    }

    for (const item of store.desired.values()) {
      expect(item).toBeInstanceOf(Array)
    }

    for (const item of store.connection.values()) {
      expect(item).toBeInstanceOf(Array)
    }
  })

  test.concurrent('move device', async () => {
    const store = useDevicesStore()

    await store.fetchDevices(false, true)
    const device = Object.values(store.devices)[0]
    const namespace = 'infinimesh'

    await store.moveDevice(device.uuid, namespace)

    expect(store.devices[device.uuid].access).toEqual({
      ...device.access, namespace
    })
  })

  test('update device', async () => {
    const store = useDevicesStore()

    await store.fetchDevices(false, true)
    const device = Object.values(store.devices)[0]

    await store.updateDevice(device.uuid, {
      ...device, title: 'test title'
    })

    expect(store.devices[device.uuid].title).toBe('test title')
  })

  test('update device config', async () => {
    const store = useDevicesStore()

    await store.fetchDevices(false, true)
    const device = Object.values(store.devices)[0]

    await store.updateDeviceConfig(device.uuid, { string: 'value' })

    expect(store.devices[device.uuid].config.toJSON())
      .toEqual({ string: 'value' })
  })

  test('delete device', async () => {
    const store = useDevicesStore()
    const startSpy = vi.spyOn(mockBar, 'start')
    const finishSpy = vi.spyOn(mockBar, 'finish')

    await store.fetchDevices(false, true)
    const device = Object.values(store.devices)[0]

    await store.deleteDevice(device.uuid, mockBar)
    await new Promise((resolve) => setTimeout(resolve, 300))
    await flushPromises()
    await store.fetchDevices(false, true)

    expect(store.devices[device.uuid]).toBeUndefined()
    expect(startSpy).toHaveBeenCalledOnce()
    expect(finishSpy).toHaveBeenCalledOnce()
  })

  test('create device', async () => {
    const store = useDevicesStore()
    const startSpy = vi.spyOn(mockBar, 'start')
    const finishSpy = vi.spyOn(mockBar, 'finish')

    const value = Math.random()
    const role = Math.floor(Math.random() * 2 + 1)
    const level = Math.floor(Math.random() * 4 + 1)

    const device = new Device({
      uuid: value.toString(16).slice(2),
      token: value.toString(8).slice(2),
      title: 'New Device',
      tags: ['tag #1'],
      access: new Access({ namespace: 'infinimesh', level, role }),
      enabled: true,
      config: {}
    })

    await store.createDevice({ device, namespace: 'infinimesh' }, mockBar)
    await new Promise((resolve) => setTimeout(resolve, 300))
    await flushPromises()

    expect(store.devices[device.uuid]).toEqual(device)
    expect(startSpy).toHaveBeenCalledOnce()
    expect(finishSpy).toHaveBeenCalledOnce()
  })

  test('toggle device', async () => {
    const store = useDevicesStore()
    const startSpy = vi.spyOn(mockBar, 'start')
    const finishSpy = vi.spyOn(mockBar, 'finish')

    await store.fetchDevices()
    const { uuid, enabled } = Object.values(store.devices)[0]

    await store.toggle(uuid, mockBar)

    expect(store.devices[uuid].enabled).toBe(!enabled)
    expect(startSpy).toHaveBeenCalledOnce()
    expect(finishSpy).toHaveBeenCalledOnce()
  })

  test('toggle device basic', async () => {
    const store = useDevicesStore()
    const startSpy = vi.spyOn(mockBar, 'start')
    const finishSpy = vi.spyOn(mockBar, 'finish')

    await store.fetchDevices()
    const { uuid, basicEnabled } = Object.values(store.devices)[0]

    await store.toggleBasic(uuid, mockBar)

    expect(store.devices[uuid].basicEnabled).toBe(!basicEnabled)
    expect(startSpy).toHaveBeenCalledOnce()
    expect(finishSpy).toHaveBeenCalledOnce()
  })

  test('join device', async () => {
    const store = useDevicesStore()

    await store.fetchDevices(false, true)
    const device = Object.values(store.devices)[0]
    const node = await store.join(new JoinRequest({
      access: device.access.level, join: device.uuid
    }))

    expect(node).toEqual(new Node(node))
  })

  test('update node of device', async () => {
    const store = useDevicesStore()

    await store.fetchDevices(false, true)
    const device = Object.values(store.devices)[0]
    const { nodes } = await store.fetchJoins(device)

    const node = await store.join(new JoinGeneralRequest({
      access: device.access.level,
      join: device.uuid,
      node: nodes[0]
    }))

    expect(node).toEqual(new Node(node))
  })

  test('fetch joins', async () => {
    const store = useDevicesStore()

    await store.fetchDevices(false, true)
    const device = Object.values(store.devices)[0]
    const { nodes } = await store.fetchJoins(device)

    expect({ nodes }).toEqual(new Nodes({ nodes }))
  })
})
