import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'

import App from './App.vue'
import SvgIcon from './components/SvgIcon.vue'
import router from './router'
import './styles/index.css'

createApp(App)
  .component('SvgIcon', SvgIcon)
  .use(createPinia())
  .use(router)
  .use(ElementPlus)
  .mount('#app')
