export const getters = {
  isAuthenticated(state) {
    return state.auth.loggedIn;
  },

  loggedInUser(state) {
    return state.auth.user;
  }
};

export const actions = {
  async getNamespaces({ commit }) {
    const namespaces = await this.$axios.$get("/namespaces");
    commit("devices/namespaces", namespaces.namespaces);
  },
  setNamespace({ commit, dispatch }, ns) {
    commit("devices/namespace", ns);
    dispatch("getDevices");
  },
  async getDevices({ commit, state }) {
    console.log("getDevices");
    let ns = "";
    if (state.devices.namespaces.length) {
      ns = state.devices.namespaces.filter(
        el => el.id == state.devices.namespace
      )[0].name;
    } else {
      ns = state.auth.user.default_namespace.name;
    }

    console.log(ns);
    const devices = await this.$axios.$get("/devices", {
      params: {
        namespace: ns
      }
    });
    commit("devices/pool", devices.devices);
  }
};
