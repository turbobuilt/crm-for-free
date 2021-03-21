<template>
    <div class="login-page">
        <div class="container">
            <br><br>
            <div class="row justify-content-md-center">
                <div class="card col-lg-4 col-md-6 col-sm-12">
                    <div class="card-body">
                        <div class="mb-3">
                            <label class="form-label">Email</label>
                            <input v-model="Email" class="form-control" required="true" minlength="4" maxlength="32"/>
                        </div>
                        <div class="mb-3" v-if="!ResettingPassword">
                            <label class="form-label">Password</label>
                            <input v-model="Password" class="form-control" type="password" name="password" required="true" minlength="8" maxlength="32"/>
                        </div>
                        <div class="d-flex justify-content-between" v-if="!ResettingPassword">
                            <button class="btn btn-primary" @click="login">Login</button>
                            <button class="btn btn-secondary" @click="resetPassword">Password Help</button>
                        </div>
                        <div class="d-flex justify-content-between" v-else>
                            <button class="btn btn-primary" @click="sendResetEmail">Send Reset Email</button>
                            <button class="btn btn-secondary" @click="stopResettingPassword">Cancel</button>
                        </div>
                        <br>
                        <div class="d-flex flex-column justify-content-center align-items-start">
                            <div>No Account yet? <router-link to="/register" class="">Sign Up Here</router-link></div>
                            <!-- <router-link to="/register" class="btn btn-secondary">Register</router-link> -->
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
<script lang="ts">
import {Vue, Options} from "vue-class-component";
import axios from "axios"
import toastr from "toastr"
import {store} from "../../store"

export default class LoginPage extends Vue {
    Email = null;
    Password = "";
    ErrorMessage = "";
    Loading = false;
    ResettingPassword = false;

    async login(){
        try{
            if(this.Loading)
                return;
            this.Loading = true;
            var res = await axios.post("/api/v1.0/user/Login", {
                Email: this.Email,
                Password: this.Password
            })
            var data = res.data;
            Object.assign(store,data);
            store.AuthorizationToken = `${store.UserKey}_${store.LoginExpiration}_${store.Verification}`;
            axios.defaults.headers.common['Authorization'] = store.AuthorizationToken;
            localStorage.setItem("AuthorizationToken", store.AuthorizationToken)
            localStorage.setItem("AuthorizationExpiration", store.LoginExpiration)
            this.$router.push("/customer")
        }catch(err){
            console.error(err);
            if(err?.response?.data?.Message) {
                toastr.error(err.response.data.Message)
            }else{
                toastr.error("Unexpected error.  Please contact support.")
            }
        } finally {
            this.Loading = false;
        }
    }

    async sendResetEmail(){
        try{
            if(this.Loading)
                return;
            this.Loading = true;
            var res = await axios.post("/api/v1.0/user/SendResetPasswordEmail", {
                Email: this.Email
            })
            var data = res.data;
            toastr.info("A password reset email has been sent to " + this.Email + "!")
            this.ResettingPassword = false;
        } catch(err) {
            console.error(err);
            if(err?.response?.data?.Message) {
                toastr.error(err.response.data.Message)
            }else{
                toastr.error("Unexpected error.  Please contact support.")
            }
        } finally {
            this.Loading = false;
        }
    }

    async resetPassword(){
        this.Password = ""
        this.ResettingPassword = true;
    }

    async stopResettingPassword(){
        this.ResettingPassword = false;
    }

}


</script>
<style lang="scss">

</style>