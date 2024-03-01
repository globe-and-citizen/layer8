<script setup>
import { computed, onMounted, ref, watch } from "vue";
import Navbar from "../components/Navbar.vue";
import { useRouter } from "vue-router";
import layer8_interceptor from 'layer8_interceptor'
 
const BACKEND_URL =  import.meta.env.VITE_BACKEND_URL
const router = useRouter();
const SpToken = ref(localStorage.getItem("SP_TOKEN") || null);
const L8Token = ref(localStorage.getItem("L8_TOKEN") || null);
const isLoggedIn = computed(() => SpToken.value !== null);

let nextPoem = ref({});

const userName = computed(() => {
  if (SpToken.value && SpToken.value.split(".").length > 1) {
    const payload = JSON.parse(atob(SpToken.value.split(".")[1]));
    return payload.email;
  }
  return "";
});

// const token = jwt.sign({ isEmailVerified, displayName, countryName }, SECRET_KEY);
const metaData = computed(() => {
  if (L8Token.value && L8Token.value.split(".").length > 1) {
    const payload = JSON.parse(atob(L8Token.value.split(".")[1]));
    return payload;
  }
  return "";
});

const getPoem = async () => {
  try {
    const resp = await layer8_interceptor.fetch( BACKEND_URL + "/nextpoem");
    let poemObj = await resp.json();

    if (poemObj.title) {
      nextPoem.value = poemObj;
    } else {
      console.error("The poemObj is malformed or other error....");
    }
  } catch (error) {
    console.error("error:", error);
  }
};

const logoutUser = () => {
  SpToken.value = null;
  localStorage.removeItem("SP_TOKEN");
  localStorage.clear();
  router.push({ name: "loginRegister" });
};

onMounted(async()=>{
  let user = localStorage.getItem("_user") ? JSON.parse(localStorage.getItem("_user")) : null
  let url = await layer8_interceptor.static(user.profile_image);
  const pictureBox = document.getElementById("pictureBox");
  pictureBox.src = url;
})

</script>

<template>
  <div class="h-screen bg-primary flex flex-col">
    <Navbar></Navbar>
    <div
      class="bg-primary-content w-full flex justify-center items-center p-4 flex-1"
    >
      <div
        class="card w-auto bg-base-100 shadow-xl p-8 h-min prose"
        v-if="isLoggedIn"
      >
        <h1>Welcome {{ userName }}!</h1>
        <!-- <div v-if="user?.profile_image"> -->
         <div> 
          <img id="pictureBox">
        </div>
        <h3>Your MetaData:</h3>
        <h4 v-if="metaData.displayName">
          Username: {{ metaData.displayName }}
        </h4>
        <h4 v-if="metaData.countryName">Country: {{ metaData.countryName }}</h4>
        <!-- Change here -->
        <h4>
          Email Verified: Email is {{ metaData.isEmailVerified ? "" : "not" }} verified!
        </h4>


        <br />
        <div class="flex gap-6">
          <button class="btn" @click="getPoem">Get Next Poem</button>
          <button class="btn btn-secondary" @click="logoutUser">Logout</button>
        </div>
        <div>
          <br />
          <div id="newPoem">
            <h3>Next Poem goes here:</h3>
            <table class="table">
              <tr>
                <th>Title:</th>
                <td>{{ nextPoem.title }}</td>
              </tr>
              <tr>
                <th>Author:</th>
                <td>{{ nextPoem.author }}</td>
              </tr>
              <tr>
                <th>Body:</th>
                <td>{{ nextPoem.body }}</td>
              </tr>
            </table>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
