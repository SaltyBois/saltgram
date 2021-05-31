import Vue from 'vue'
import VueRouter from 'vue-router'
import Home from '../views/Home.vue'
import axios from 'axios'
// import store from '../main.js'

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home,
    // NOTE(Jovan): Logged out if switching to home, just as protonmail
    // beforeEnter: (to, from, next) => {
    //   if(this.$store.state.jws) 
    //     next("/user");
    //   else
    //     next();
    // },
  },
 {
   path: '/user',
   name: 'User',
   component: () => import(/* webpackChunkName: "user" */ '../views/User.vue')
 },
 {
   // NOTE(Jovan): Confirm reset
  path: '/email/reset/:token',
  name: 'PasswordReset',
  beforeEnter: (to, from, next) => {
    let token = to.params["token"]
    axios.put("email/reset/" + token, {withCredentials: true})
      .then(r => {
        console.log(r);
        next();
      })
      .catch(r => {
        console.log(r);
        next({name: "Home"});
      });
  },
  component: () => import(/* webpackChunkName: "passwordReset" */ '../views/PasswordReset.vue')
 },
 {
   path: '/forgotpassword',
   name: 'ForgotPassword',
   component: () => import(/* webpackChunkName: "forgotPassword" */ '../views/ForgotPassword.vue')
 },
 {
   path: '/email/activate/:token',
   name: 'ActivateEmail',
   beforeEnter: (to, from, next) => {
    let token = to.params["token"]
    axios.put("email/activate/" + token)
      .finally(function(){
        next({ name: "Home"});
      })
   },
 },
{
    path: '/user/settings',
    name: 'Settings',
    component: () => import(/* webpackChunkName: "userSettings" */ '../views/UserSettings.vue')
},
{
    path: '/main', // TODO(Mile): This needs to be set on path '/' when auth successful
    name: 'Main',
    component: () => import(/* webpackChunkName: "mainPage" */ '../views/MainPage.vue')
},
{
    path: '/inbox', // TODO(Mile): This needs to be set on path '/' when auth successful
    name: 'Inbox',
    component: () => import(/* webpackChunkName: "inbox" */ '../views/Inbox.vue')
},
{
    path: '/notifications', // TODO(Mile): This needs to be set on path '/' when auth successful
    name: 'Notifications',
    component: () => import(/* webpackChunkName: "notifications" */ '../views/Notifications.vue')
},
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

// TODO(Jovan): Authentication
// router.beforeEach((to, from, next) => {
//   next(to);          // TODO(MILE): COMMENT THIS AFTER DEVELOPMENT PHASE AND UNCOMMENT BELOW
//   // console.log("Looking for jwt: ", store.state["jws"])
//   // let jws = store.state["jws"];
//   // axios.put("auth", to.path, {headers: {"Authorization" : "Bearer " + jws}})
//   //   .then(r => {
//   //     console.log(r);
//   //     next();
//   //   })
//   //   .catch(r => {
//   //     console.log(r);
//   //     next({name: "Home"})
//   //   })
// });

export default router
