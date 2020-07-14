export const state = () => ({
  namespace: "",
  namespaces: []
});

export const mutations = {
  namespace(state, val) {
    state.namespace = val;
  },
  namespaces(state, val) {
    state.namespaces = val;
  }
};
