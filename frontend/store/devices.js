import Vue from "vue";

export const state = () => ({
  namespace: "",
  namespaces: [],
  pool: [],
  states_pool: {},
});

export const mutations = {
  namespace(state, val) {
    state.namespace = val;
  },
  namespaces(state, val) {
    state.namespaces = val;
  },
  update_namespace(state, ns) {
    Vue.set(
      state.namespaces,
      state.namespaces.findIndex((el) => el.id === ns.id),
      ns
    );
  },
  pool(state, val) {
    state.pool = val;
  },
  states_pool(state, val) {
    state.states_pool = { ...state.states_pool, ...val };
  },
};

export const actions = {
  async get({ commit, state, rootState }) {
    if (window.nuxt) window.$nuxt.$loading.start();

    let ns = "";
    if (state.namespaces.length) {
      ns = state.namespaces.filter((el) => el.id == state.namespace)[0].id;
    } else {
      ns = rootState.auth.user.default_namespace.id;
    }

    const devices = await this.$axios.$get("/api/devices", {
      params: {
        namespaceid: ns,
      },
    });
    commit("pool", devices.devices);

    if (window.$nuxt) window.$nuxt.$loading.finish();
  },
  async state({ commit, state, rootState }) {
    if (window.nuxt) window.$nuxt.$loading.start();

    let ns = "";
    if (state.namespaces.length) {
      ns = state.namespaces.filter((el) => el.id == state.namespace)[0].id;
    } else {
      ns = rootState.auth.user.default_namespace.id;
    }

    const states = await this.$axios.$get("/api/devices/states/all", {
      params: {
        id: ns,
      },
    });

    commit("states_pool", states.pool);

    if (window.$nuxt) window.$nuxt.$loading.finish();
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
        device: device,
      })
      .then((res) => {
        dispatch("get");
        if (success) success(res);
      })
      .catch((e) => {
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
  },
  async getNamespacePermissions({ commit }, ns) {
    commit("update_namespace", { ...ns, loading: true });
    this.$axios
      .$get(`/api/namespaces/${ns.id}/permissions`)
      .then((permissions) => {
        ns = { ...ns, ...permissions };
        commit("update_namespace", ns);
      })
      .catch(() => {
        commit("update_namespace", ns);
      });
  },
};

export const getters = {
  get: (state) => (id) => {
    return state.pool.filter((el) => el.id == id)[0];
  },
  get_state: (state) => (id) => {
    return state.states_pool[id];
  },
  currentNamespace: (state) => state.namespace,
};
