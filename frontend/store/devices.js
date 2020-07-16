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
