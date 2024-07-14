<!-- This complete code (CallBack View) is a part of Layer8 Component -->
<script setup>
import { computed, ref } from "vue";
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import layer8_interceptor from 'layer8_interceptor'

const router = useRouter()
const code = ref(new URLSearchParams(window.location.search).get("code"))
const token = ref(localStorage.getItem("token") || null)
const BACKEND_URL =  import.meta.env.VITE_BACKEND_URL


onMounted(() => {
    setTimeout(() => {
        layer8_interceptor.fetch(BACKEND_URL + "/api/login/layer8/auth", {
            method: "POST",
            headers: {
                "Content-Type": "Application/Json"
            },
            body: JSON.stringify({
                callback_url: window.location.href,
            })
        })
            .then(res => res.json())
            .then(data => {
                localStorage.setItem("L8_TOKEN", data.token)
                router.push({ name: 'home' })
            })
            .catch(err => console.log(err))
    }, 1000);
})
</script>

<template>
    <div>
        <h1>Login with layer8...</h1>
    </div>
</template>

