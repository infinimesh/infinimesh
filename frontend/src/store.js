import Vue from "vue";
import Vuex from "vuex";

Vue.use(Vuex);

export default new Vuex.Store({
  state: {
    devices: {
      testid123: {
        status: "active",
        name: "Device 3",
        location: "Berlin",
        tags: "test"
      },
      testid124: {
        status: "active",
        name: "Device 4",
        location: "Ddorf",
        tags: "test"
      }
    }
  },
  getters: {
    getDevice: (state, deviceId) => {
      return state.devices.deviceId;
    },
    getAllDevices: (state) => {
      return state.devices;
    }
  },
  mutations: {
    addDevice: (state, device) => {
      state.devices[Object.keys(device)[0]] = device.deviceId;
      console.log(device)
    },
    deleteDevice: (state, deviceId) => {
      delete state.devices.deviceId;
    }
  },
  actions: {
    addDevice: ({ commit }, device) => {
      commit("addDevice", device);
      return device;
    },
    deleteDevice: ({ commit }, deviceId) => {
      commit("deleteDevice", deviceId);
      return deviceId;
    }
  }
});
