<template>
  <v-select persistent-hint return-object label="Project" class="custom"  v-on:change="onChanged" :items="namespaces" v-model="selected" item-text="name" item-value="name" >
  </v-select>
</template>

<script>
import { mapGetters } from "vuex";

export default {
  data() {
    return {
      selected: "hanswurst",
      currentRoute: this.$route.name,
      namespaces: ["abc", "def"]
    };
  },
  computed: {
    ...mapGetters({ namespace: "getNamespace" })
  },
  created() {
    this.$store
      .dispatch("fetchNamespaces")
      .then(() => {
        let namespaces = this.$store.getters.getNamespaces;
        this.namespaces = namespaces.map(namespace => {
          return namespace.name;
        });
        this.selected = this.$store.getters.getNamespace;
        console.log("set selected ns to", this.selected);
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
