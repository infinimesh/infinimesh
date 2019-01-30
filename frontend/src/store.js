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
    nodeTree: {},
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
    getShadowMessages: state => {
      return state.shadow.messages;
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
    addShadowMessages: (state, messages) => {
      state.shadow.messages = messages;
    },
    storeNodeTree: (state, tree) => {
      state.nodeTree = transform(tree);
    },
    addChildNode: (state, payload) => {
      addNode(state.nodeTree, payload.id, payload.node);
    },
    deleteNode: (state, id) => {
      deleteNode(state.nodeTree, id);
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
    connectToShadow: ({ commit }, id) => {
      let xhr = new XMLHttpRequest();

      setTimeout(() => {
        xhr.open(
          "GET",
          Vue.http.options.root + `/devices/${id}/shadow/reported`,
          true
        );
        xhr.onprogress = function() {
          let jsonObjects = [];
          let obj = "";
          let messages = [];

          jsonObjects = xhr.responseText.replace(/\n$/, "").split(/\n/);
          for (obj of jsonObjects) {
            messages.splice(0, 0, JSON.parse(obj));
          }
          commit("addShadowMessages", messages);
        };
        xhr.send();
      }, 1000);
    },
    fetchNodeTree: ({ commit }) => {
      return new Promise((resolve, reject) => {
        commit("apiRequestPending", true);
        return Vue.http
          .get("objects")
          .then(response => {
            commit("apiRequestPending", false);
            commit("storeNodeTree", response.body);
            resolve();
          })
          .catch(error => {
            commit("apiRequestPending", false);
            commit("apiDataFailure", error);
            reject(error);
          });
      });
    },
    addChildNode: ({ commit }, payload) => {
      return new Promise((resolve, reject) => {
        commit("apiRequestPending", true);
        return Vue.http
          .post(`objects/${payload.parent}/children`, {
            name: payload.name
          })
          .then(response => {
            if (response.status === 200) {
              commit("apiRequestPending", false);
              let node = {
                name: payload.name,
                id: response.body.uid,
                type: "node",
                children: []
              };
              let obj = {
                id: payload.parent,
                node
              };
              commit("addChildNode", obj);
              resolve();
            }
          })
          .catch(error => {
            commit("apiRequestPending", false);
            commit("apiDataFailure", error);
            reject(error);
          });
      });
    },
    deleteNode: ({ commit }, id) => {
      commit("deleteNode", id);
    },
    updateDevice: ({ commit }, properties) => {
      commit("updateDevice", properties);
      return properties;
    },
    unRegisterDevice: ({ commit }, id) => {
      commit("unRegisterDevice", id);
      return id;
    }
  }
});

const transformObject = input => {
  let res = {};

  res.id = input.uid;
  res.name = input.name;
  res.type = input.type;
  res.children = [];
  if (input.devices) {
    for (let device of input.devices) {
      device.type = "device";
      res.children.push(transformObject(device));
    }
  }
  if (input.objects) {
    for (let object of input.objects) {
      object.type = "node";
      res.children.push(transformObject(object));
    }
  }
  return res;
};

const transform = input => {
  if (isEmpty(input)) {
    return;
  } else {
    let res = [];
    for (let value of input.objects) {
      value.type = "node";
      let el = transformObject(value);
      el.type = "node";
      res.push(el);
    }
    for (let value of input.devices) {
      value.type = "device";
      let el = transformObject(value);
      el.type = "device";
      res.push(el);
    }
    return res;
  }
};

const addNode = (input, id, node) => {
  for (let element of input) {
    if (element.id === id) {
      let newArr = element.children;
      newArr.push(node);
      element.children = newArr;
      return node.id;
    } else if (element.children) {
      addNode(element.children, id, node);
    }
  }
};

const deleteNode = (input, id) => {
  for (let element of input) {
    if (element.id === id) {
      console.log("input", input.indexOf(element));
      input.splice(input.indexOf(element), 1);
    } else if (element.children) {
      deleteNode(element.children, id);
    }
  }
};

const isEmpty = obj => {
  for (var key in obj) {
    if (obj.hasOwnProperty(key)) return false;
  }
  return true;
};
