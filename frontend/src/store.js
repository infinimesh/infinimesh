import Vue from "vue";
import Vuex from "vuex";

Vue.use(Vuex);

export default new Vuex.Store({
  state: {
    devices: [
      {
        enabled: false,
        id: "25",
        tags: ["test", "bbc"],
        certificate: "abc"
      },
      {
        enabled: true,
        id: "6",
        tags: ["test"],
        certificate: "abdd"
      }
    ],
    model: {
      enabled: undefined,
      id: "",
      tags: [],
      certificate: {
        pem_data: "",
        algorithm: ""
      }
    }
  },
  getters: {
    getDevice: state => id => {
      let device;
      let key;
      device = state.devices.find(device => device.id === id);
      for (key in state.model) {
        if (!device[key])
          device[key] = state.model[key]
      }
      return device;
    },
    getAllDevices: state => {
      let device;
      let key;
      for (device of state.devices) {
        for (key in state.model) {
          if (!device[key]) {
            device[key] = state.model[key]
          }
        }
      }
      return state.devices;
    }
  },
  mutations: {
    addDevice: (state, device) => {
      let deviceExists;
      deviceExists = state.devices.find(item => item.id === device.id);
      if (!deviceExists) {
        state.devices.push(device);
      }
    },
    unRegisterDevice: (state, id) => {
      let deviceIndex;

      deviceIndex = state.devices.findIndex(device => device.id === id);
      state.devices.splice(deviceIndex, 1);
    }
  },
  actions: {
    addDevice: ({ commit }, device) => {
      commit("addDevice", device);
      return device;
    },
    unRegisterDevice: ({ commit }, id) => {
      commit("unRegisterDevice", id);
      return id;
    }
  }
});
