<template>
  <v-navigation-drawer
    app
    clipped
    permanent
    expand-on-hover
    mini-variant-width="64"
  >
    <v-list nav dark>
      <v-list-item
        link
        :key="page.link"
        v-for="page in pagesFiltered"
        :nuxt="true"
        :to="page.link"
      >
        <v-list-item-icon>
          <v-icon>{{ page.icon }}</v-icon>
        </v-list-item-icon>
        <v-list-item-title>{{ page.title }}</v-list-item-title>
      </v-list-item>
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
            <a @click="generateTokenVisible = true">Generate Token</a>
          </v-list-item-title>
          <account-generate-token
            :active="generateTokenVisible"
            @cancel="generateTokenVisible = false"
          />
        </v-list-item>
        <v-list-item>
          <v-list-item-title>
            <a @click="resetAccountPasswordVisible = true">Reset password</a>
          </v-list-item-title>
          <account-reset-password
            :active="resetAccountPasswordVisible"
            :account="user"
            @cancel="resetAccountPasswordVisible = false"
            @reset="handleResetAccountPassword"
          />
        </v-list-item>
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
import AccountGenerateToken from "@/components/account/GenerateToken.vue";
import AccountResetPassword from "@/components/account/ResetPassword.vue";
import AccountControlMixin from "@/mixins/account-control";

export default {
  name: "layout-sider",
  mixins: [AccountControlMixin],
  components: {
    AccountGenerateToken,
    AccountResetPassword,
  },
  data() {
    return {
      generateTokenVisible: false,
      resetAccountPasswordVisible: false,

      pages: [
        {
          title: "Device Registry",
          icon: "mdi-cloud-outline",
          link: "/dashboard/devices",
          admin: false,
        },
        {
          title: "Accounts",
          icon: "mdi-account-group",
          link: "/dashboard/accounts",
          admin: true,
        },
        {
          title: "Namespaces",
          icon: "mdi-folder-multiple-outline",
          link: "/dashboard/namespaces",
          admin: true,
        },
      ],
    };
  },
  computed: {
    pagesFiltered() {
      return this.pages.filter(
        (page) =>
          this.allowedScope(page.link) &&
          (!page.admin || this.user.is_admin || this.user.is_root)
      );
    },
    user() {
      return this.$store.getters.loggedInUser;
    },
    route: {
      get() {
        return [this.$route.name];
      },
      set(val) {
        this.$router.push(val[0]);
      },
    },
  },
  methods: {
    allowedScope(scope) {
      return this.$store.getters["window/hasAccess"](scope);
    },
    handleResetAccountPassword(password) {
      this.resetAccountPasswordVisible = false;
      this.updateAccount(
        this.user.uid,
        { password: password },
        "",
        "Reset password failed",
        () => {
          this.$message.success("Password changed successfuly");
        }
      );
    },
  },
};
</script>
