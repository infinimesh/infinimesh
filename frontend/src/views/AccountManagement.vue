<template>
  <div>
    <v-card-title>
      <h1 class="mb-3">Account Management</h1>
    </v-card-title>
    <v-card-text>
      <v-layout row wrap>
        <v-autocomplete
          v-model="model"
          :items="accounts"
          :loading="isLoading"
          :search-input.sync="search"
          color="white"
          hide-no-data
          hide-selected
          item-text="name"
          label="Registered accounts"
          placeholder="Start typing to Search"
          prepend-icon="mdi-database-search"
          return-object
        ></v-autocomplete>
        <v-spacer></v-spacer>
        <v-btn
          color="primary lighten-1"
          round
          @click="creatingUser = true; editingUser = false; model = false"
        >
          <v-icon>
            add
          </v-icon>
        </v-btn>
      </v-layout>
    </v-card-text>
    <v-divider></v-divider>
  <v-card-text>
    <v-expand-transition>
      <v-list
        v-if="model"
        class="grey lighten-4 indigo--text"
      >
        <v-list-tile
          v-for="(field, i) in fields"
          :key="i"
        >
          <v-list-tile-content>
            <v-list-tile-title v-text="field.value"></v-list-tile-title>
            <v-list-tile-sub-title v-text="field.key"></v-list-tile-sub-title>
          </v-list-tile-content>
        </v-list-tile>
      </v-list>
    </v-expand-transition>
    <v-expand-transition>
      <new-user
       v-if="creatingUser"
       @newUser="newUser=$event"
      >
      </new-user>
    </v-expand-transition>
    <v-alert :value="saving.value" type="success" icon="check_circle">
      {{ saving.message }}
    </v-alert>
    <v-alert :value="saving.value" type="error" icon="error">
      {{ saving.message }}
    </v-alert>
  </v-card-text>
  <v-card-actions>
    <v-btn
      v-if="!editingUser && model"
      @click="editingUser = true; creatingUser = false"
      round
      class="mr-4"
    >
      Edit
    </v-btn>
    <v-btn
      v-if="editingUser || creatingUser"
      :disabled="!newUser.passwordTwo"
      @click="save"
      round
      class="mr-4"
    >
      Save
    </v-btn>
    <v-btn
      v-if="model || creatingUser"
      @click="model = null; editingUser = false; creatingUser = false"
      round
    >
      Close
      <v-icon right>mdi-close-circle</v-icon>
    </v-btn>
  </v-card-actions>
</div>
</template>

<script>
import NewUser from "../components/NewUser.vue";

export default {
  data: () => ({
    descriptionLimit: 60,
    isLoading: false,
    model: null,
    search: null,
    editingUser: false,
    creatingUser: false,
    newUser: {},
    saving: {
      value: false,
      message: ""
    }
  }),
  computed: {
    fields() {
      if (!this.model) return [];

      return Object.keys(this.model).map(key => {
        return {
          key,
          value: this.model[key] || "n/a"
        };
      });
    },
    accounts() {
      let accounts = this.$store.getters.getAccounts;

      return accounts.map(account => {
        let name = account.name;
        return Object.assign({}, account, { name });
      });
    }
  },
  watch: {
    search() {
      // Items have already been loaded
      if (this.accounts.length > 0) return;

      // Items have already been requested
      if (this.isLoading) return;

      this.isLoading = true;

      // Lazily load input items
      this.$store
        .dispatch("fetchAccounts")
        .finally(() => (this.isLoading = false));
    }
  },
  methods: {
    save() {
      this.isLoading = true;

      this.$http
        .post("http://localhost:8081/accounts/users", {
          name: this.newUser.name,
          password: this.newUser.passwordOne,
          is_root: false
        })
        .then(res => console.log(res))
        .catch(err => {
          console.log(err);
          this.saving.value = true;
          this.saving.message = "Failure to create new account";
          setTimeout(() => (this.saving.value = false), 2000);
        })
        .finally(() => {
          this.isLoading = false;
          this.creatingUser = false;
          this.saving.value = true;
          this.saving.message = "Account created";
          setTimeout(() => (this.saving.value = false), 2000);
        });
    }
  },
  components: {
    NewUser
  }
};
</script>

<style lang="css" scoped></style>
