import { describe, test, expect, vi } from "vitest";
import { mount } from "@vue/test-utils";
import { usePluginsStore } from "@/store/plugins";
import { defineComponent } from "vue";
import { createPinia, setActivePinia } from "pinia";

const TestComponent = defineComponent({
  template: "<div></div>",
});

const pinia = createPinia();

describe("store/sessions", () => {
  beforeEach(() => {
    const wrapper = mount(TestComponent, {
      global: {
        plugins: [pinia],
      },
    });

    setActivePinia(pinia);
  });

  test("fetchPlugins", async () => {
    const pluginsStore = usePluginsStore();

    const promise = pluginsStore.fetchPlugins();
    expect(pluginsStore.loading).toBe(true);

    await promise;

    expect(pluginsStore.plugins).not.toBe([]);
    expect(pluginsStore.loading).toBe(false);
  });

  test("create", async () => {
    const pluginsStore = usePluginsStore();

    const data = await pluginsStore.create({
      title: "plugin",
      description: "plugin desk",
    });

    expect(data).not.toBe(null);
    expect(data.uuid).toBeTypeOf("string");
  });

  test("create", async () => {
    const pluginsStore = usePluginsStore();

    const data = await pluginsStore.create({
      title: "plugin",
      description: "plugin desk",
    });

    expect(data).not.toBe(null);
    expect(data.uuid).toBeTypeOf("string");
  });

  test("update", async () => {
    const pluginsStore = usePluginsStore();

    const device = {
      uuid: "someuuid",
      title: "plugin",
      description: "plugin desk",
    };
    const data = await pluginsStore.update(device.uuid, device);

    expect(data).not.toBe(null);
    expect(data.uuid).toBe(device.uuid);
  });

  test("height with value", async () => {
    const pluginsStore = usePluginsStore();

    const device = {
      uuid: "someuuid",
      title: "plugin",
      description: "plugin desk",
    };

    pluginsStore.heights.set(device.uuid, 200);
    const height = pluginsStore.height(device.uuid);

    expect(height).toBe("200px");
  });

  test("height default", async () => {
    const pluginsStore = usePluginsStore();

    const height = pluginsStore.height("123");
    expect(height).toBe("10vh");
  });
});
