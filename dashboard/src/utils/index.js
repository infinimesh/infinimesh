import { h } from "vue"
import { NIcon } from "naive-ui";

export const renderIcon = (icon) => {
  return () => {
    return h(NIcon, null, {
      default: () => h(icon)
    });
  };
};

export const renderIconColored = (icon, color) => {
  return () => {
    return h(NIcon, { color }, {
      default: () => h(icon)
    });
  };
}