<template>
  <v-select v-on:change="onChanged" :items="namespaces" label="Namespace" v-model="selected" item-text="name" item-value="name">
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
          console.log("Got ns, ns menu");
          let namespaces = this.$store.getters.getNamespaces;
          console.log("namespaces", namespaces)
          this.namespaces = namespaces.map(namespace => {
            return namespace.name;
          });
        })
        .catch(e => console.log(e));
    },
    methods: {
      onChanged(a) {
        console.log("changed", a);
          this.$store.dispatch("setNamespace", a);
          console.log(a);
      }
    }
  };
</script>

<style lang="css" scoped>

</style>
