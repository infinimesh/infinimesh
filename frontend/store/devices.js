export const state = () => ({
  namespace: "",
  namespaces: [],
  pool: []
});

export const mutations = {
  namespace(state, val) {
    state.namespace = val;
  },
  namespaces(state, val) {
    state.namespaces = val;
  },
  pool(state, val) {
    state.pool = val;
  }
};

export const actions = {
  async get({ commit, state, rootState }) {
    let ns = "";
    if (state.namespaces.length) {
      ns = state.namespaces.filter(el => el.id == state.namespace)[0].id;
    } else {
      ns = rootState.auth.user.default_namespace.id;
    }

    const devices = await this.$axios.$get("/api/devices", {
      params: {
        namespace: ns
      }
    });
    commit("pool", devices.devices);
  },
  /**
   *
   * @param {object} device - Device object for creation
   * @param {Function} success - Callback on device creation success
   * @param {Function} error - Callback on device creation error
   * @param {Function} always - Callback invoked after all previous callbacks no matter if there is an error or not
   */
  add({ dispatch }, { device, success, error, always }) {
    if (!device) return;
    this.$axios
      .$post("/api/devices", {
        device: device
      })
      .then(res => {
        dispatch("get");
        if (success) success(res);
      })
      .catch(e => {
        if (error) error(e);
      })
      .then(() => {
        if (always) always();
      });
  },
  async getNamespaces({ commit }) {
    const namespaces = await this.$axios.$get("/api/namespaces");
    commit("namespaces", namespaces.namespaces);
  },
  setNamespace({ commit, dispatch }, ns) {
    commit("namespace", ns);
    dispatch("get");
    $nuxt.$router.push({ name: "dashboard" });
  }
};

export const getters = {
  get: state => id => {
    return state.pool.filter(el => el.id == id)[0];
  }
};
