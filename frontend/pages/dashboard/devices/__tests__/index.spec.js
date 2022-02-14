import { shallowMount } from "@vue/test-utils";
import Registry from "../index";

import Vue from "vue";
import Vuex from "vuex";
import Antd from "ant-design-vue";
import vuetify from "vuetify";

function factory(store) {
  return shallowMount(Registry, { store: store });
}

describe("Devices Registry", () => {
  let wrapper;

  let default_store = {
    commit: jest.fn((mutation, data) => {}),
    getters: {
      loggedInUser: {
        default_namespace: { id: "0x0" },
      },
      "devices/currentNamespace": jest.fn(() => "0x0"),
    },
    state: {
      devices: {
        namespace: "0x0",
        namespaces: [
          {
            id: "0x0",
            name: "test",
            markfordeletion: false,
            deleteinitiationtime: "0000-01-01T00:00:00Z",
          },
        ],
        pool: [],
      },
      window: {
        width: 1920,
        height: 990,
        gridSize: "xxl",
        menu: false,
        noAccessScopes: [],
        topAction: undefined,
        release: {
          html_url:
            "https://github.com/infinimesh/infinimesh/releases/tag/v0.1.5",
          tag_name: "v0.1.5",
        },
      },
      auth: {
        user: {
          uid: "0x1",
          name: "test",
          is_root: false,
          enabled: true,
          default_namespace: {
            id: "0x0",
            name: "test",
            markfordeletion: false,
            deleteinitiationtime: "",
          },
          password: "",
          is_admin: false,
          owner: "",
          username: "test",
        },
        loggedIn: true,
        strategy: "local",
      },
    },
  };

  beforeEach(() => {
    Vue.use(Vuex);
    Vue.use(Antd);
    Vue.use(vuetify);
  });

  it("mounts properly", () => {
    wrapper = factory(default_store);
    expect(wrapper.vm).toBeTruthy();
  });
});
