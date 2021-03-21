<template>
    <div class="signup-page">
        <div class="container">
            <br><br>
            <div class="row justify-content-md-center">
                <div class="card col-lg-4 col-md-6 col-sm-12">
                    <div class="card-body" v-if="!WaitingForConfirmation">
                        <div class="mb-3">
                            <label class="form-label">Email</label>
                            <input v-model="Email" class="form-control" required="true" minlength="4" maxlength="32"/>
                        </div>
                        <div class="d-flex flex-column justify-content-center">
                            <button @click="signUp" class="btn btn-primary">{{Loading ? "Working..." : "Sign Up"}}</button>
                        </div>
                        <br>
                        <div class="d-flex flex-column justify-content-center align-items-start">
                            <div>Already Have an Account? <router-link to="/login" class="">Login</router-link></div>
                        </div>
                    </div>
                    <div class="card-body" v-if="WaitingForConfirmation">
                        <p>Please check your email to confirm your account</p>
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
    WaitingForConfirmation = false;

    async signUp(){
        try{
            if(this.Loading)
                return;
            if(!this.Email || !this.Email.match(/.{2,}@.{2,}\..{2,}/)){
                toastr.error("You need to enter a valid email address")
                return 
            }
            this.Loading = true;
            var res = await axios.post("/api/v1.0/user", {
                Email: this.Email
            })
            this.WaitingForConfirmation = true;
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
}


</script>
<style lang="scss">

</style>