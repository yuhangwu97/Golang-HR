import Vue from 'vue'
import Vuex from 'vuex'
import auth from './modules/auth'
import employee from './modules/employee'
import department from './modules/department'

Vue.use(Vuex)

export default new Vuex.Store({
  modules: {
    auth,
    employee,
    department
  }
})