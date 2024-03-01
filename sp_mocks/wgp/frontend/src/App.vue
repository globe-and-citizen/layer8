<script setup>
import { RouterView } from "vue-router";
import layer8_interceptor from 'layer8_interceptor'

try{
  layer8_interceptor.initEncryptedTunnel({
    providers: ["http://localhost:8000"],
    proxy: "http://localhost:5001" // works
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


