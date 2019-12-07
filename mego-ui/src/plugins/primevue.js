import Vue from 'vue';

import InputText from 'primevue/inputtext';
import Button from 'primevue/button';
import Toast from 'primevue/toast';
import AutoComplete from 'primevue/autocomplete';
import ToastService from 'primevue/toastservice';
import Tree from 'primevue/tree';

Vue.use(ToastService);

Vue.component('InputText', InputText);
Vue.component('Button', Button);
Vue.component('Toast', Toast);
Vue.component('AutoComplete', AutoComplete);
Vue.component('Tree', Tree);

import 'primevue/resources/themes/nova-light/theme.css';
import 'primevue/resources/primevue.min.css';
import 'primeicons/primeicons.css';

