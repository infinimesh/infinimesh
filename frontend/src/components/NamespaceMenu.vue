<template>
  <v-select :items="namespaces" label="Namespace" v-model="selected" item-text="name" item-value="name">
  </v-select>
</template>

<script>
  import {
    mapGetters
  } from 'vuex'
  
  
  export default {
    data: () => ({
      selected: "hanswurst",
      namespaces: [],
    }),
    created() {
      console.log("BF");
      this.$store
        .dispatch("fetchNamespaces")
        .then(() => {
          let namespaces = this.$store.getters.getNamespaces;
  
          this.namespaces = namespaces.map(namespace => {
            return namespace.name;
          });
        })
        .catch(e => console.log(e));
    },
  };
</script>

<style lang="css" scoped>
  
</style>