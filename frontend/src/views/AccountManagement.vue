<template>
  <div>
  <v-card-title>
    <h1 class="mb-3">Account Management</h1>
  </v-card-title>
  <v-card-text>
    <v-layout
      row
      wrap
    >
      <v-autocomplete
        v-model="model"
        :items="items"
        :loading="isLoading"
        :search-input.sync="search"
        color="white"
        hide-no-data
        hide-selected
        item-text="Description"
        item-value="API"
        label="Registered accounts"
        placeholder="Start typing to Search"
        prepend-icon="mdi-database-search"
        return-object
      ></v-autocomplete>
      <v-spacer></v-spacer>
      <v-btn
        color="primary lighten-1"
        round
        @click="isEditing = true"
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
      <v-card
        v-if="isEditing"
        flat
      >
        <v-card-text>
          <v-text-field
            label="Account data"
            clearable
          ></v-text-field>
        </v-card-text>
      </v-card>
    </v-expand-transition>
  </v-card-text>
  <v-card-actions>
    <v-btn
      :disabled="!isEditing && !model"
      @click="isEditing = true"
      round
    >
      {{ (isEditing) ? "Update" : "Edit"}}
    </v-btn>
    <v-spacer></v-spacer>
    <v-btn
      :disabled="!model && !isEditing"
      @click="model = null; isEditing = false"
      round
    >
      Close
      <v-icon right>mdi-close-circle</v-icon>
    </v-btn>
  </v-card-actions>
</div>
</template>

<script>
export default {
  data: () => ({
    descriptionLimit: 60,
    entries: [],
    isLoading: false,
    model: null,
    search: null,
    isEditing: false
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
    items() {
      return this.entries.map(entry => {
        const Description =
          entry.Description.length > this.descriptionLimit
            ? entry.Description.slice(0, this.descriptionLimit) + "..."
            : entry.Description;

        return Object.assign({}, entry, { Description });
      });
    }
  },
  watch: {
    search() {
      // Items have already been loaded
      if (this.items.length > 0) return;

      // Items have already been requested
      if (this.isLoading) return;

      this.isLoading = true;

      // Lazily load input items
      fetch("https://api.publicapis.org/entries")
        .then(res => res.json())
        .then(res => {
          const { count, entries } = res;
          this.count = count;
          this.entries = entries;
        })
        .catch(err => {
          console.log(err);
        })
        .finally(() => (this.isLoading = false));
    }
  }
};
</script>

<style lang="css" scoped>
</style>
