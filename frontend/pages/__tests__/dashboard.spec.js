import { shallowMount } from "@vue/test-utils";
import Dashboard from "../dashboard";

import Vue from "vue";
import Vuex from "vuex";
import Antd from "ant-design-vue";
import vuetify from "vuetify";

describe("Dashboard", () => {
  let wrapper;
  let store;
  let actions;

  beforeEach(() => {
    Vue.use(Vuex);
    Vue.use(Antd);
    Vue.use(vuetify);

    actions = {
      "devices/getNamespaces": jest.fn()
    };
    store = new Vuex.Store({
      namespaced: true,
      actions
    });

    wrapper = shallowMount(Dashboard, { store });
  });

  it("mounts properly", () => {
    expect(wrapper.vm).toBeTruthy();
  });

  it("renders properly", () => {
    expect(wrapper.html()).toContain('<v-app-stub id="dashboard">');
  });

  it("obtains namespaces", () => {
    expect(actions["devices/getNamespaces"]).toHaveBeenCalled();
  });
});
