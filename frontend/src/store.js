import Vue from "vue";
import Vuex from "vuex";
import VueResource from "vue-resource";

Vue.use(Vuex);
Vue.use(VueResource);

export default new Vuex.Store({
  state: {
    apiDataPending: false,
    apiDataFailure: {
      status: false,
      error: ""
    },
    devices: [],
    shadow: {
      initialState: {
        data: "No data received",
        timestamp: "N/A"
      },
      messages: []
    },
    model: {
      enabled: undefined,
      id: "",
      tags: [],
      certificate: {
        pem_data: "",
        algorithm: ""
      }
    },
    nodeTree: {}
  },
  getters: {
    getDevice: state => id => {
      let device;
      let key;
      device = state.devices.find(device => device.id === id);
      if (device) {
        for (key in state.model) {
          if (!device[key]) device[key] = state.model[key];
        }
        return device;
      } else {
        return undefined;
      }
    },
    getInitialShadow: state => {
      return state.shadow.initialState;
    },
    getAllDevices: state => {
      if (state.devices) {
        let device;
        let key;
        for (device of state.devices) {
          for (key in state.model) {
            if (!device[key]) {
              device[key] = state.model[key];
            }
          }
        }
        return state.devices;
      } else {
        return undefined;
      }
    },
    getNodeTree: state => {
      if (state.nodeTree) {
        return state.nodeTree;
      } else {
        return undefined;
      }
    }
  },
  mutations: {
    apiRequestPending: (state, status) => {
      state.apiDataPending = status;
    },
    apiDataFailure: (state, error) => {
      state.apiDataFailure.status = true;
      state.apiDataFailure.error = error;
    },
    storeDevices: (state, devices) => {
      state.devices = devices;
    },
    storeShadow: (state, apiResponse) => {
      state.shadow.initialState.data = apiResponse.data;
      state.shadow.initialState.timestamp = apiResponse.timestamp;
    },
    updateDevice: (state, properties) => {
      let deviceIndex;
      let property;
      deviceIndex = state.devices.findIndex(
        device => device.id === properties.id
      );
      if (deviceIndex) {
        for (property in properties) {
          state.devices[deviceIndex][property] = properties[property];
        }
      }
    },
    unRegisterDevice: (state, id) => {
      let deviceIndex;
      deviceIndex = state.devices.findIndex(device => device.id === id);
      if (deviceIndex) {
        state.devices.splice(deviceIndex, 1);
      } else {
        return "Device Id doesn't exist";
      }
    },
    setNodeTree: (state, tree) => {
      state.nodeTree = tree;
    }
  },
  actions: {
    fetchDevices(store) {
      return new Promise((resolve, reject) => {
        store.commit("apiRequestPending", true);
        return Vue.http
          .get("devices")
          .then(response => {
            store.commit("apiRequestPending", false);
            store.commit("storeDevices", response.body.devices);
            resolve();
          })
          .catch(error => {
            store.commit("apiRequestPending", false);
            store.commit("apiDataFailure", error);
            reject(error);
          });
      });
    },
    fetchInitialShadow: ({ commit }, id) => {
      return new Promise((resolve, reject) => {
        commit("apiRequestPending", true);
        return Vue.http
          .get(`devices/${id}/shadow`)
          .then(response => {
            commit("apiRequestPending", false);
            commit("storeShadow", response.body.shadow.reported);
            resolve();
          })
          .catch(error => {
            commit("apiRequestPending", false);
            commit("apiDataFailure", error);
            reject(error);
          });
      });
    },
    updateDevice: ({ commit }, properties) => {
      commit("updateDevice", properties);
      return properties;
    },
    unRegisterDevice: ({ commit }, id) => {
      commit("unRegisterDevice", id);
      return id;
    },
    setNodeTree: ({ commit }, tree) => {
      commit("setNodeTree", tree);
    }
  }
});
