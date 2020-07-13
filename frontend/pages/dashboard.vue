<template>
  <div id="components-layout-demo-basic">
    <a-layout>
      <a-layout-header class="wide-header">
        <Header v-model="menu_collapsed" />
      </a-layout-header>
      <a-layout>
        <a-layout-sider>Sider</a-layout-sider>
        <a-layout>
          <a-layout-content>
            <nuxt-child />
          </a-layout-content>
          <a-layout-footer>Footer</a-layout-footer>
        </a-layout>
      </a-layout>
    </a-layout>
  </div>
</template>

<script>
import Header from "@/components/Header.vue";

export default {
  components: {
    Header
  },
  data() {
    return {
      menu_collapsed: false
    };
  },
  fetch({ store, params }) {
    store.dispatch("getNamespaces");
  }
};
</script>

<style>
.wide-header {
  padding: 0 !important;
}

#components-layout-demo-basic {
  text-align: center;
}

#components-layout-demo-basic .ant-layout-footer {
  background: #7dbcea;
  color: #fff;
}
#components-layout-demo-basic .ant-layout-footer {
  line-height: 1.5;
}
#components-layout-demo-basic .ant-layout-sider {
  background: #3ba0e9;
  color: #fff;
  line-height: 120px;
}
#components-layout-demo-basic .ant-layout-content {
  background: rgba(16, 142, 233, 1);
  color: #fff;
  min-height: 120px;
  line-height: 120px;
}
#components-layout-demo-basic > .ant-layout {
  margin-bottom: 48px;
}
#components-layout-demo-basic > .ant-layout:last-child {
  margin: 0;
}
</style>
