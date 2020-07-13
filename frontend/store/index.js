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
    commit("devices/namespaces", namespaces);
  }
};
