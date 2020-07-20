import Vue from 'vue'
// 导入 vue-router包
import VueRouter from 'vue-router'
// 手动安装 vue-router
Vue.use(VueRouter)

//导入 account 和 goodslist 组件
import Login from './components/Login.vue'
import Main from './components/Main.vue'

// 创建路由对象
let router = new VueRouter({
    routes: [
        // account goodslist
        {
            path: '/',
            name: '登录',
            component: Login,
            meta:{
                title:'登录'
            }
        },
        {
            path: '/main',
            name: '用户管理',
            component: Main,
            meta:{
                title:'用户管理'
            }
        }
    ]
});

// 导出 router 对象
export default router;