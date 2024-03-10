<script setup>
import { RouterView } from "vue-router";
import layer8_interceptor from 'layer8_interceptor'
const BACKEND_URL =  import.meta.env.VITE_BACKEND_URL
const PROXY_URL =  import.meta.env.VITE_PROXY_URL

try{
  layer8_interceptor.initEncryptedTunnel({
    providers: [BACKEND_URL],
    proxy: PROXY_URL // works
  }, "dev")
}catch(err){
  console.log(".initEncryptedTunnel error: ", err)
}
</script>

<template>
  <div>
    <RouterView />
  </div>
</template>

<style scoped>
</style>


