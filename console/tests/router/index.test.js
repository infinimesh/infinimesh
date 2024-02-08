import { describe, test, expect, vi } from "vitest";
import { mount } from "@vue/test-utils";
import { router } from "@/router";
import { defineComponent, markRaw } from "vue";
import { createPinia } from "pinia";
import piniaPersist from "pinia-plugin-persist";

const TestComponent = defineComponent({
  template: `<div></div>`,
});

describe("router", () => {
  beforeEach(async () => {
    router.push({ name: "Login" });
    await router.isReady;

    const pinia = createPinia();
    pinia.use(piniaPersist);
    pinia.use(({ store }) => {
      store.$router = markRaw(router);
    });

    window.PLATFORM_NAME = "infinimesh";
    const wrapper = mount(TestComponent, {
      global: {
        plugins: [router, pinia],
      },
    });
  });

  test("push", async () => {
    const pushSpy = vi.spyOn(router, "push");

    router.push({ name: "Login" });
    await router.isReady();

    expect(router.currentRoute.value.name).toBe("Login");
    expect(pushSpy).toHaveBeenCalledOnce();
  });

  test("authorized route protected", async () => {
    router.push({ name: "Dashboard" });
    await router.isReady();

    expect(router.currentRoute.value.name).not.toBe("Dashboard");
    expect(router.currentRoute.value.name).toBe("Login");
  });

  test("document title", async () => {
    router.push({ name: "Login" });
    await router.isReady();

    expect(document.title).toContain("Login");
  });
});
