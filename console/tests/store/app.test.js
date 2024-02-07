import { describe, test, expect, vi } from "vitest";
import { mount } from "@vue/test-utils";
import { useAppStore } from "@/store/app";
import { defineComponent, markRaw } from "vue";
import { createPinia, setActivePinia } from "pinia";

const TestComponent = defineComponent({
  template: "<div></div>",
});

const mockRouter = {
  push: vi.fn(),
};

const pinia = createPinia();
pinia.use(({ store }) => {
  store.$router = markRaw(mockRouter);
});

describe("store/app", () => {
  beforeEach(() => {
    setActivePinia(createPinia());

    const wrapper = mount(TestComponent, {
      global: {
        plugins: [pinia],
      },
    });
  });

  test("offline redirect", async () => {
    const routerSpy = vi.spyOn(mockRouter, "push");

    const appStore = useAppStore();

    appStore.offline();

    expect(routerSpy).toHaveBeenCalledOnce();
    expect(routerSpy).toHaveBeenCalledWith({ name: "Offline" });
  });

  test("logout redirect", async () => {
    const routerSpy = vi.spyOn(mockRouter, "push");

    const appStore = useAppStore();

    appStore.logout();

    expect(routerSpy).toHaveBeenCalledOnce();
    expect(routerSpy).toHaveBeenCalledWith({ name: "Login", query: {} });
  });

  test("logout redirect with message", async () => {
    const msg = "test";

    const routerSpy = vi.spyOn(mockRouter, "push");

    const appStore = useAppStore();

    appStore.logout(msg);

    expect(routerSpy).toHaveBeenCalledOnce();
    expect(routerSpy).toHaveBeenCalledWith({
      name: "Login",
      query: { msg: btoa(JSON.stringify(msg)) },
    });
  });

  test("logout reset state", async () => {
    const appStore = useAppStore();
    const resetSpy = vi.spyOn(appStore, "$reset");

    appStore.token = "sometoken";
    expect(appStore.token).not.toBe("");

    appStore.logout();

    expect(appStore.token).toBe("");
    expect(resetSpy).toHaveBeenCalledOnce();
  });
});
