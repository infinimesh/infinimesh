<template>
  <div id="device">
    <a-row style="padding-top: 10px">
      <a-col :xxl="{ span: 4, offset: 1 }">
        <transition name="fade">
          <h1 class="lead" v-if="device">
            {{ device.name }} <span class="muted">{{ device.id }}</span>
          </h1>
        </transition>
      </a-col>
    </a-row>
    <a-row>
      <a-col :xxl="{ span: 12, offset: 1 }">
        <transition-group name="slide">
          <a-card title="Details" key="details" v-if="device"></a-card>
        </transition-group>
      </a-col>
    </a-row>
  </div>
</template>

<script>
export default {
  computed: {
    device() {
      return this.$store.getters["devices/get"](this.$route.params.id);
    }
  },
  mounted() {
    console.log("device");
    console.log(this.device);
  },
  validate({ params }) {
    return /0[xX][0-9a-fA-F]+/.test(params.id);
  }
};
</script>

<style>
#device {
  overflow: hidden;
  font-family: Exo;
  font-weight: 500;
}
.lead {
  font-size: 2em;
  color: #fff;
}
.muted {
  opacity: 0.7;
}

.slide-leave-active,
.slide-enter-active {
  transition: 1s;
}
.slide-enter {
  transform: translate(100%, 0);
}
.slide-leave-to {
  transform: translate(-100%, 0);
}
</style>
