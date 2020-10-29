export const getters = {
  isAuthenticated(state) {
    return state.auth.loggedIn;
  },

  loggedInUser(state) {
    return state.auth.user;
  }
};

export const actions = {
  async logout() {
    await this.$auth.logout();
  }
};
