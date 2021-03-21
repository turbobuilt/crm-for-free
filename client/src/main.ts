import { createApp } from 'vue'
// TypeScript error? Run VSCode command
// TypeScript: Select TypeScript version - > Use Workspace Version
import App from './App.vue'
import * as VueRouter from "vue-router"
import routes from "./routes.vue";
import "toastr/build/toastr.min.css"

import axios from 'axios';
axios.defaults.baseURL = window.location.host === "erp.turbobuilt.com" ? "http://erp.turbobuilt.com/" : "http://localhost:8080"


const router = VueRouter.createRouter({
    history: VueRouter.createWebHashHistory(),
    routes, // short for `routes: routes`
})
  
const app = createApp(App)
app.use(router)

import TextEditor from "./components/TextEditor.vue";
app.component("TextEditor",TextEditor)

app.mount('#app')
