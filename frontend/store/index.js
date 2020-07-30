import { state } from "./devices";

export const getters = {
  isAuthenticated(state) {
    return state.auth.loggedIn;
  },

  loggedInUser(state) {
    return state.auth.user;
  }
};

export const mutations = {
  logout(state) {
    state.auth.loggedIn = false;
  }
};

export const actions = {
  async logout() {
    await this.$auth.logout();
  }
};
