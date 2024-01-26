import { describe, test, vi, expect } from "vitest";
import { mount } from "@vue/test-utils";
import AccessBadge from "@/components/core/access-badge";
import { Level } from "infinimesh-proto/build/es/node/access/access_pb";

describe("access-badge", () => {
  test("display admin level (string)", async () => {
    const props = {
      access: "ADMIN",
    };

    const wrapper = mount(AccessBadge, {
      props,
    });

    expect(wrapper.text()).toContain("Admin");
  });

  test("display none level (number)", async () => {
    const props = { access: Level.NONE };

    const wrapper = mount(AccessBadge, {
      props,
    });

    expect(wrapper.text()).toContain("None");
  });

  test("display not existed level (number) throw error", async () => {
    const props = { access: 666 };

    expect(() =>
      mount(AccessBadge, {
        props,
      })
    ).toThrowError();
  });

  test("call cb", async () => {
    const props = {
      access: Level.NONE,
      cb: vi.fn(),
    };

    const cbSpy = vi.spyOn(props, "cb");

    const wrapper = mount(AccessBadge, {
      props,
    });

    await wrapper.find("button").trigger("click");

    expect(cbSpy).toHaveBeenCalledOnce();
    expect(cbSpy).toHaveBeenCalledWith(Level.NONE);
  });

  test("call disabled cb", async () => {
    const props = {
      access: Level.NONE,
      cb: vi.fn(),
      disabled: true,
    };

    const cbSpy = vi.spyOn(props, "cb");

    const wrapper = mount(AccessBadge, {
      props,
    });

    await wrapper.find("button").trigger("click");

    expect(cbSpy).not.toHaveBeenCalledOnce();
  });
});
