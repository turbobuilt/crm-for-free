<script lang="ts">
import LoginPage from "./pages/UserPages/LoginPage.vue";
import CustomerListPage from "./pages/Customer/CustomerListPage.vue";
import CustomerPage from "./pages/Customer/CustomerPage.vue";
import RegisterPage from "./pages/UserPages/RegisterPage.vue";

const routes = [
    { path: '/register', component: RegisterPage },
    { path: '/login', component: LoginPage },
    { path: '/customer', component: CustomerListPage },
    { path: '/customer/:id', component: CustomerPage}
]
export default routes;
</script>