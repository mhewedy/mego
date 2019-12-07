import Vue from 'vue';
import './plugins/axios'
import App from './App.vue';

import './plugins/primevue'

Vue.config.productionTip = false;

new Vue({
  render: h => h(App),
}).$mount('#app')
