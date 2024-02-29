<script setup>
import { ref } from "vue";
import Navbar from "../components/Navbar.vue";
const BackendURL = "http://localhost:5001";
import layer8_interceptor from 'layer8_interceptor'

const BACKEND_URL =  import.meta.env.VITE_BACKEND_URL
const requestsSent = ref(0);
const totalTimeSpent = ref(0);
const numberOfRequest = ref(0)

console.log("verdict 1: ", layer8_interceptor.checkEncryptedTunnel())
setTimeout(()=>{
  console.log("verdict 2: ", layer8_interceptor.checkEncryptedTunnel())
}, 1000)

async function testWASMHandler() {
  const startTime = performance.now();
  for (let i = 0; i < numberOfRequest.value; i++) {
    const res = await layer8_interceptor.testWASM(i, "42");
    requestsSent.value++;
    console.log(res);
  }
  const endTime = performance.now();
  totalTimeSpent.value = endTime - startTime;
  console.log("Total request sent: ", requestsSent.value)
  console.log("Total time spent: ", totalTimeSpent.value, "ms")
}

async function getError(){
  try {
    console.log("Error Test")
    await layer8_interceptor.fetch(BackendURL + "/error", {
      method: "POST",
      headers: {
        "Content-Type": "Application/Json",
      },
      body: JSON.stringify({}),
    });
  } catch (error) {
    console.log(error);
    alert("Registration failed!");
    isRegister.value = true;
  }
}

let x = 0
async function getNextPicture(){
  let idx = x%2
  const pictureURLs = [
    'http://localhost:8000/media/boy.png',
    'http://localhost:8000/media/girl.png',
  ]
  let url = await layer8_interceptor.static(pictureURLs[idx]);
  const element = document.getElementById("pictureBox");
  element.src = url;
  x++
}


</script>

<template>
  <Navbar></Navbar>
  <div class="greetings">
    <div>
      <label for="">Number of request</label>
      <input type="text" v-model="numberOfRequest">
      <button @click="testWASMHandler" class="text-green-500">Execute</button>

      <div>
        Total request sent: {{ requestsSent }} times
      </div>

      <div>
        Total time spent: {{ totalTimeSpent }} ms
      </div> 

      <div>
        <button @click="getError"> Get Error from "/error"</button>
      </div>

      <div>
        <button @click="getNextPicture"> Get Next Picture</button>
      </div>
      <hr>
      <img id="pictureBox">

    </div>
  </div>
</template>

<style scoped></style>
