import Vue from 'vue';

import InputText from 'primevue/inputtext';
import Button from 'primevue/button';
import Toast from 'primevue/toast';
import AutoComplete from 'primevue/autocomplete';
import ToastService from 'primevue/toastservice';
import Tree from 'primevue/tree';
import Calendar from 'primevue/calendar';
import Spinner from 'primevue/spinner';
import ProgressSpinner from 'primevue/progressspinner';
import Message from 'primevue/message';
import Dialog from 'primevue/dialog';
import Textarea from 'primevue/textarea';

Vue.use(ToastService);

Vue.component('InputText', InputText);
Vue.component('Button', Button);
Vue.component('Toast', Toast);
Vue.component('AutoComplete', AutoComplete);
Vue.component('Tree', Tree);
Vue.component('Calendar', Calendar);
Vue.component('Spinner', Spinner);
Vue.component('ProgressSpinner', ProgressSpinner);
Vue.component('Message', Message);
Vue.component('Dialog', Dialog);
Vue.component('Textarea', Textarea);

import 'primevue/resources/themes/nova-light/theme.css';
import 'primevue/resources/primevue.min.css';
import 'primeicons/primeicons.css';
import 'primeflex/primeflex.css';

