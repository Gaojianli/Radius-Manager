import Vue from 'vue'
import App from './App.vue'
import vuetify from './plugins/vuetify';
import router from './router'
import '@mdi/font/css/materialdesignicons.css'
import store from './store/index'
import VueCookies from 'vue-cookies'
Vue.use(VueCookies)
Vue.config.productionTip = false
router.beforeEach((to, _, next) => {
  if (to.meta.title) {
    document.title = `${to.meta.title} - Radius Manager`;
  }
  next()
})
new Vue({
  icons: {
    iconfont: 'mdi', // default - only for display purposes
  },
  vuetify,
  store,
  render: h => h(App),
  router: router
}).$mount('#app')