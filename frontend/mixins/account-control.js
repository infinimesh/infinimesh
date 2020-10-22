export default {
  methods: {
    getAccountsPool() {
      const vm = this;
      vm.loading = true;
      vm.$axios
        .get("/api/accounts")
        .then(res => (vm.accounts = res.data.accounts))
        .catch(e => {
          if (e.response.status == 403) {
            vm.$notification.error({
              message: "Oops",
              description: e.response.data.message
            });
            vm.$store.commit("window/noAccess", "dashboard-accounts");
            vm.$router.push({ name: "dashboard-devices" });
          }
        })
        .then(() => (vm.loading = false));
    },
    deleteAccount(account) {
      const vm = this;
      this.$axios({
        url: `/api/accounts/${account.uid}`,
        method: "delete"
      })
        .then(() => {
          vm.$message.success("Account successfuly deleted!");
          vm.getAccountsPool();
        })
        .catch(e => {
          vm.$notification.error({
            message: "Error deleting account " + account.name,
            description: e.response.data.message
          });
        });
    },
    toogleAccount(account) {
      this.updateAccount(
        account.uid,
        {
          enabled: !account.enabled
        },
        `Account successfuly ${account.enabled ? "disabled" : "enabled"}!`,
        `Error ${account.enabled ? "disabling" : "enabling"} account`
      );
    },
    handleAccountAdd(account) {
      const vm = this;
      account.account.owner = vm.user.uid;
      vm.$axios({
        method: "post",
        url: "/api/accounts",
        data: account
      })
        .then(() => {
          vm.$notification.success({
            message: "Account created successfuly"
          });
          vm.createAccountDrawerVisible = false;
          vm.getAccountsPool();
        })
        .catch(err => {
          this.$notification.error({
            message: "Failed to create an account",
            description: `Response: ${err.response.data.message}`,
            duration: 10
          });
        });
    },
    updateAccount(id, data, success, error, success_callback, error_callback) {
      const vm = this;
      vm.loading = true;
      vm.$axios({
        method: "patch",
        url: `/api/accounts/${id}`,
        data: data
      })
        .then(
          success_callback
            ? success_callback
            : () => {
                vm.$message.success(success);
                vm.getAccountsPool();
              }
        )
        .catch(
          error_callback
            ? error_callback(e)
            : e => {
                vm.$notification.error({
                  message: error,
                  description: e.response.data.message
                });
              }
        )
        .then(() => (vm.loading = false));
    }
  }
};
