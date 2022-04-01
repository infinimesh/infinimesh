import { createApp } from 'vue'
import App from './App.vue'

import router from './router/router'

import axios from 'axios'
import VueAxios from 'vue-axios'

import { createPinia } from 'pinia'
const pinia = createPinia()

const app = createApp(App)
  
app
  .use(pinia)
  .use(router)
  .use(VueAxios, axios)
  .provide('axios', app.config.globalProperties.axios)
  .mount('#app')
