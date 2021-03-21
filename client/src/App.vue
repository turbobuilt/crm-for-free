<template>
	<div>
		<router-view></router-view>
	</div>
</template>

<script lang="ts">
// import { defineComponent } from 'vue'
import {Vue, Options} from "vue-class-component";
import VueRouter from 'vue-router'
import axios from "axios"



export default class App extends Vue {
	created(){
		var AuthorizationToken = localStorage.getItem("AuthorizationToken") || ""
		var LoginExpiration = localStorage.getItem("AuthorizationExpiration") || ""
		if(AuthorizationToken && parseInt(LoginExpiration) > Date.now()){
            axios.defaults.headers.common['Authorization'] = AuthorizationToken;
			this.$router.push("/customer")
		}else {
			this.$router.push("/login")
		}
	}
}
</script>

<style lang="scss">

</style>