import { shallowMount } from "@vue/test-utils";
import Logout from "../logout";

import Vue from "vue";
import Antd from "ant-design-vue";
import vuetify from "vuetify";

describe("Logout", () => {
  let wrapper = shallowMount(Logout);

  it("mounts properly", () => {
    expect(wrapper.vm).toBeTruthy();
  });

  it("renders properly", () => {
    expect(wrapper.html()).toBe("<h1>Logging out...</h1>");
  });
});
