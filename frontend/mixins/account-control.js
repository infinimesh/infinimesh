export default {
  methods: {
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
