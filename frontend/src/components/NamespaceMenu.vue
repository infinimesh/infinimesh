<template>
  <v-select v-on:change="onChanged" :items="namespaces" label="Namespace" v-model="selected" item-text="name" item-value="name">
  </v-select>
</template>

<script>
import { mapGetters } from "vuex";

export default {
  data() {
    return {
      selected: "hanswurst",
      currentRoute: this.$route.name,
      namespaces: []
    };
  },
  computed: {
    ...mapGetters({ namespace: "getNamespace" })
  },
  created() {
    console.log("BF");
    this.$store
      .dispatch("fetchNamespaces")
      .then(() => {
        console.log("Got ns, ns menu");
        let namespaces = this.$store.getters.getNamespaces;
        console.log("namespaces", namespaces);
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
      this.navigateTo(this.namespace);
    },
    navigateTo(namespace) {
      this.$router.push({
        name: this.currentRoute,
        params: {
          namespace
        }
      });
    }
  }
};
</script>

<style lang="css" scoped>

</style>
