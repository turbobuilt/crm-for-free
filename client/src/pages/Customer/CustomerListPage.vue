<template>
    <div class="customer-list-page">
        <div class="container">
            <div class="row">
                <div class="col-xs-12">
                    <input class="main-search" placeholder="Search..." :modelValue="SearchText" @input="searchTextChanged"/>
                </div>
            </div>
            <div class="row add-new-button-container">
                <div class="col-xs-12 justify-content-center d-flex">
                    <router-link class="btn btn-large btn-secondary new-button" :to="`/customer/new`">New</router-link>
                </div>
            </div>
            <div class="row">
                <div  class="col-xl-3 col-lg-4 col-md-6 col-sm-12" v-for="(customer,index) in customers" :key="index">
                    <CustomerListComponent :customer="customer" :searchText="SearchTextStatic"/>
                </div>
            </div>
        </div>
    </div>
</template>
<script lang="ts">
import {Vue, Options} from "vue-class-component"
import axios from "axios"
import CustomerListComponent from "./CustomerListComponent.vue";

@Options({
    components:{CustomerListComponent}
})
export default class CustomerListPage extends Vue {
    Loading = false;
    customers = null;
    SearchText = ""
    SearchTextStatic = ""
    SearchTimeout: any;

    created(){
        this.getItems();
    }

    async searchTextChanged(event: any){
        this.SearchText = event.target.value;
        if(this.SearchTimeout) {
            clearTimeout(this.SearchTimeout)
            this.SearchTimeout = null;
        }
        this.SearchTimeout = setTimeout(() => {
            this.getItems()
        }, 500)
    }

    async getItems(){
        try{
            if(this.Loading)
                return;
            this.Loading = true;
            var res = await axios.get("/api/v1.0/customer",{params: {
                SearchText: this.SearchText
            }})
            this.SearchTextStatic = this.SearchText;
            this.customers = res.data;
        }catch(err){
            console.error(err);
            if(err?.response?.data?.Message) {
                toastr.error(err.response.data.Message)
            }else{
                toastr.error("Unexpected error.  Please contact support.")
            }
        }finally{
            this.Loading = false;
        }
    }
}
</script>
<style lang="scss">
.customer-list-page {
    .card-item {
        margin-bottom: 20px;
    }
    .main-search {
        padding: 4px 2px;
        border-radius: 2px solid silver;
        width: 100%;
        margin: 20px 0;
    }
    .add-new-button-container {
        margin-bottom: 20px;
        .new-button {
            width: 100%;
        }
    }
}
</style>