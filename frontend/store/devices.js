export const state = () => ({
  namespace: "",
  namespaces: []
});

export const mutations = {
  namespaces(state, val) {
    state.namespaces = val;
  }
};
