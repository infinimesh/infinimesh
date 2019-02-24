<template>
  <div>
    <div style="overflow: hidden; clear: both">
    </div>
    <div style="overflow: hidden, clear: both">
      <Box name="Glass" />
      <Box name="Banana" />
      <Box name="Paper" />
    </div>
  </div>
</template>

<script>
import Box from "./Box";
import { DragDropContext } from "vue-react-dnd";
import HTML5Backend from "react-dnd-html5-backend";

export default {
  name: "container",
  components: {
    Box
  },
  mixins: [DragDropContext(HTML5Backend)]
};
</script>

<style>
</style>
