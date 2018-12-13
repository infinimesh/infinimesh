import Vue from "vue";
import Vuex from "vuex";

Vue.use(Vuex);

export default new Vuex.Store({
  state: {
    devices: [
      {
        status: "inactive",
        deviceId: 25,
        name: "Device 6",
        location: "Düsseldorf",
        tags: "test"
      },
      {
        status: "inactive",
        deviceId: 6,
        name: "Device 6",
        location: "Düsseldorf",
        tags: "test"
      }
    ]
  },
  getters: {
    getDevice: (state, deviceId) => {
      return state.devices.find(device => device.deviceId === deviceId);
    },
    getAllDevices: state => {
      return state.devices;
    }
  },
  mutations: {
    addDevice: (state, device) => {
      state.devices.push(device);
    },
    deleteDevice: (state, deviceId) => {
      let deviceIndex;

      deviceIndex = state.devices.findIndex(
        device => device.deviceId === deviceId
      );
      state.devices.splice(deviceIndex, 1);
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
