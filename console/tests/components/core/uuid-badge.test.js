import { describe, test, expect } from "vitest";
import { mount } from "@vue/test-utils";
import { h } from "vue";
import UuidBadge from "@/components/core/uuid-badge.vue";
import { NMessageProvider } from "naive-ui";

const getTempUuid = () => "550e8400-e29b-41d4-a716-446655440000";

const getComponent = (props) =>
  h(NMessageProvider, {}, () => h(UuidBadge, props));

describe("uuid-badge", () => {
  test("display short uuid", async () => {
    const props = {
      uuid: getTempUuid(),
    };

    const wrapper = mount(getComponent(props));

    expect(wrapper.find("#short_uuid").text()).toBe(props.uuid.slice(0, 8));
  });

  test("not display full uuid", async () => {
    const props = {
      uuid: getTempUuid(),
    };

    const wrapper = mount(getComponent(props));
    expect(wrapper.find("#full_uuid").exists()).toBe(false);
  });

  test("click event", async () => {
    const props = {
      uuid: getTempUuid(),
    };

    const wrapper = mount(getComponent(props));
    await wrapper.find("button").click;
  });
});
