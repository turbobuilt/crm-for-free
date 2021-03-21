<template>
    <div class="card card-item">
        <div class="card-body">
            <router-link :to="`/customer/${customer.Id}`" class="card-title">
                {{customer.Name}}
            </router-link>
            <div>
                {{preview}}
            </div>
        </div>
    </div>
</template>
<script lang="ts">
import {Vue, Options} from "vue-class-component";

@Options({
    props: ["customer", "searchText"]
})
export default class CustomerListComponent extends Vue {
    readonly searchText: any;
    readonly customer: any;

    get preview(){
        if(!this.searchText){
            return;
        }
        // var words = this.customer.NotesSearch.split(" ")
        // var index = words.indexOf(this.searchText);
        // var result = words.slice(Math.max(index - 5,0), index + 5).join(" ")
        var text = this.customer.NotesSearch.slice(0, 200)
        if(text.length < this.customer.NotesSearch.length)
            text += "..."
        return text
    }   

}

</script>
<style lang="scss">
.customer-list-component {

}
</style>