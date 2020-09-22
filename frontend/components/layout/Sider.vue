<template>
  <v-navigation-drawer
    app
    clipped
    permanent
    expand-on-hover
    color="#104e83"
    mini-variant-width="64"
  >
    <v-list nav v-model="route" dark>
      <v-list-item link :key="page.link" v-for="page in pagesFiltered">
        <v-list-item-icon>
          <v-icon>{{ page.icon }}</v-icon>
        </v-list-item-icon>
        <v-list-item-title>{{ page.title }}</v-list-item-title>
      </v-list-item>
      <!-- <v-list-item no-action>
        <v-list-item-icon>
          <v-icon>{{ page.icon }}</v-icon>
        </v-list-item-icon>
        <v-list-item-title>{{ page.title }}</v-list-item-title>
      </v-list-item> -->
      <v-list-group no-action prepend-icon="mdi-account-circle-outline">
        <template v-slot:activator>
          <v-list-item-title
            ><span>
              <b>{{ user.name }}</b>
            </span></v-list-item-title
          >
        </template>
        <v-list-item>
          <v-list-item-title>
            <a @click="$store.dispatch('logout')">Log Out</a>
          </v-list-item-title>
        </v-list-item>
      </v-list-group>
    </v-list>
  </v-navigation-drawer>
</template>

<script>
export default {
  data() {
    return {
      pages: [
        {
          title: "Device Registry",
          icon: "mdi-cloud-outline",
          link: "dashboard-devices"
        },
        {
          title: "Accounts",
          icon: "mdi-account-group",
          link: "dashboard-accounts"
        },
        {
          title: "Namespaces",
          icon: "mdi-folder-multiple-outline",
          link: "dashboard-namespaces"
        }
      ]
    };
  },
  computed: {
    pagesFiltered() {
      return this.pages.filter(page => this.allowedScope(page.link));
    },
    user() {
      return this.$store.getters.loggedInUser;
    },
    route: {
      get() {
        return [this.$route.name];
      },
      set(val) {
        this.$router.push({ name: val[0] });
      }
    }
  },
  methods: {
    allowedScope(scope) {
      return this.$store.getters["window/hasAccess"](scope);
    }
  }
};
</script>
