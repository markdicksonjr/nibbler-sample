import Vue from 'vue'
import VueI18n from 'vue-i18n'
import App from './App.vue'
import messages from './i18n'
import router from './router'
import './plugins/iview.js'

Vue.use(VueI18n)

Vue.config.productionTip = false

const i18n = new VueI18n({
  locale: 'en',
  messages,
})

new Vue({
  i18n,
  router,
  render: h => h(App)
}).$mount('#app')
