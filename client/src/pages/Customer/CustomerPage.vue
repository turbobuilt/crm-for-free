<template>
    <div class="customer-page">
        <div class="container">
            <div class="row">
                <div class="col-xs-12">
                    <input class="Customer" placeholder="Name" v-model="Customer.Name"/>
                </div>
            </div>
            <div class="row">
                <div class="col-xs-12">
                    <TextEditor v-model="Customer.Notes" v-model:plainText="Customer.NotesSearch"/>
                </div>
                <!-- {{Customer.NotesSearch}} -->
            </div>
            <div class="row">
                <button class="btn btn-primary" @click="save">Save</button>
            </div>
        </div>
    </div>
</template>
<script lang="ts">
import {Vue, Options} from "vue-class-component"
import toastr from "toastr"
import axios from "axios"
import { store } from "../../store";

@Options({})
export default class CustomerPage extends Vue {
    Loading = false;
    Customer = {
        Name: "",
        Notes: "",
        NotesSearch: "",
    }

    created(){
        if(this.id != 'new'){
            this.getItem();
        }
    }

    get id(){
        return this.$route.params.id;
    }

    async getItem(){
        try{
            if(this.Loading)
                return;
            this.Loading = true;
            
            var res = await axios.get("/api/v1.0/customer/" + this.id)
            this.Customer = res.data;
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

    async save(){
        try{
            if(this.Loading)
                return;
            this.Loading = true;
            // this.Customer.NotesSearch = 
            console.log("Id is ", this.id)
            if(!this.id || this.id === 'new') {
                var res = await axios.post("/api/v1.0/customer", this.Customer)
                var data = res.data;
                this.$router.replace("/customer/" + data.id)
            } else {
                var res = await axios.put("/api/v1.0/customer/" + this.id, this.Customer)
                var data = res.data;
            }
            toastr.info("Saved")
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
.customer-page {
    
}
</style>