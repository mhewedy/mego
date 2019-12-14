import Vue from 'vue';
import './plugins/axios'
import App from './App.vue';

import './plugins/primevue'

import Axios from 'axios'

Vue.prototype.$http = Axios;
const token = localStorage.getItem('mego_token');
if (token) {
    Vue.prototype.$http.defaults.headers.common['Authorization'] = 'bearer ' + token;
}

Vue.config.productionTip = false;

new Vue({
    render: h => h(App),
}).$mount('#app');
