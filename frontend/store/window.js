export const state = () => ({
  width: 0,
  gridSize: "xs",
  menu: true
});

export const mutations = {
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
  set(state, width) {
    state.commit("setWidth", width);
    let size = "xs";
    for (let [k, v] of Object.entries(gridSizes)) {
      if (width >= v) {
        size = k;
      } else {
        break;
      }
    }
    state.commit("setGridSize", size);
  },
  toggleMenu(state, val) {
    state.commit("setMenu", val);
  }
};

class WrongGridSizeException {
  constructor(given) {
    this.name = "WrongGridSizeException";
    this.message = `Wrong grid size name, can be ['xs', 'sm', 'md', 'lg', 'xl', 'xxl'], but "${given} were given."`;
  }
}
