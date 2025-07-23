import Vue from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'

// Element UI 2.x
import ElementUI from 'element-ui'
import 'element-ui/lib/theme-chalk/index.css'

// Ant Design Vue 1.x (只导入需要的组件，避免CSS冲突)
import { Layout, Menu, Avatar, Dropdown } from 'ant-design-vue'

Vue.config.productionTip = false

Vue.use(ElementUI)
Vue.use(Layout)
Vue.use(Menu)
Vue.use(Avatar)
Vue.use(Dropdown)

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount('#app')