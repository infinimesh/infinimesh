export default {
  computed: {
    drawerSize() {
      switch (this.$store.state.window.gridSize) {
        case "xxl":
          return "30%";
          break;
        case "xl":
          return "50%";
          break;
        case "lg":
          return "60%";
          break;
        case "md":
          return "75%";
          break;
        case "sm":
          return "90%";
          break;
        case "xs":
          return "100%";
          break;
      }
      return "0%";
    }
  }
};
