export const state = () => ({
  width: 0,
  height: 0,
  gridSize: "xs",
  menu: true,
  noAccessScopes: [],
  topAction: undefined
});

export const mutations = {
  setTopAction(state, action) {
    state.topAction = action;
  },
  unsetTopAction(state) {
    state.topAction = undefined;
  },
  setHeight(state, height) {
    state.height = height;
  },
  setWidth(state, width) {
    state.width = width;
  },
  setGridSize(state, size) {
    if (["xs", "sm", "md", "lg", "xl", "xxl"].includes(size)) {
      state.gridSize = size;
    } else {
      throw new WrongGridSizeException(size);
    }
  },
  setMenu(state, val) {
    state.menu = val;
  },
  noAccess(state, scope) {
    state.noAccessScopes.push(scope);
  }
};

const gridSizes = {
  sm: 576,
  md: 768,
  lg: 992,
  xl: 1200,
  xxl: 1600
};

export const actions = {
  set(state, win) {
    state.commit("setWidth", win.width);
    state.commit("setHeight", win.height);
    let size = "xs";
    for (let [k, v] of Object.entries(gridSizes)) {
      if (win.width >= v) {
        size = k;
      } else {
        break;
      }
    }
    state.dispatch("setGrid", size);
  },
  toggleMenu(state, val) {
    state.commit("setMenu", val);
  },
  setGrid(state, size) {
    state.commit("setGridSize", size);
    state.commit("setMenu", ["xs", "sm"].includes(size));
  }
};

class WrongGridSizeException {
  constructor(given) {
    this.name = "WrongGridSizeException";
    this.message = `Wrong grid size name, can be ['xs', 'sm', 'md', 'lg', 'xl', 'xxl'], but "${given} were given."`;
  }
}

export const getters = {
  menu(state) {
    return state.menu;
  },
  hasAccess: state => scope => {
    return !state.noAccessScopes.includes(scope);
  }
};
