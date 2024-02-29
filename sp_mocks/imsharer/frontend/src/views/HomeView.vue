<script setup>
import { onMounted, ref } from 'vue';
import layer8 from "layer8_interceptor";
import emitter from '@/plugins/mitt';

const isLoaded = ref(false);
const images = ref([]);
const modalImage = ref(null);

const fetchImages = () => {
  isLoaded.value = false;

  layer8.fetch('http://localhost:6001/api/gallery')
    .then((response) => response.json())
    .then(async (data) => {
      var imgs = []; 
      for (var i = 0; i < data.data.length; i++) {
        const image = data.data[i];
        const url = await layer8.static(image.url);
        imgs.push({
          id: image.id,
          name: image.name,
          url: url
        });
      }
      images.value = imgs;
      isLoaded.value = true;
    });
}

onMounted(async () => {
  await layer8.initEncryptedTunnel({
    providers: ["http://localhost:6001"],
    proxy: "http://localhost:5001" // works
  }, "dev")

  fetchImages();
});

emitter.on('reload_images', () => {
  fetchImages();
});

window.addEventListener('click', (e) => {
  if (e.target === document.querySelector('.modal') && modalImage.value) {
    modalImage.value = null;
  }
});
</script>

<template>
  <main>
    <section v-if="isLoaded">
      <section v-if="images.length === 0" class="notif">
        <p>No Images Found</p>
      </section>
      <section v-else class="gallery">
        <article v-for="image in images" :key="image.id">
          <video v-if="image.name.endsWith('.mp4')" controls class="item">
            <source :src="image.url" />
          </video>
          <img v-else class="item" :src="image.url" :alt="image.name" @click="modalImage = image" />
        </article>
      </section>
    </section>
    <section v-else class="loader">
      <p>Loading Images</p>
    </section>
    <section v-if="modalImage" class="modal">
      <img :src="modalImage.url" :alt="modalImage.name" />
    </section>
  </main>
</template>

<style>
main {
  padding: 0.5rem 5rem;
}

.gallery {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
  grid-template-rows: repeat(auto-fit, 260px);
  grid-auto-flow: dense;
  grid-gap: 0.3rem;
}

.gallery article {
  grid-column-end: span 1;
  grid-row-end: span 1;
  cursor: pointer;
}

.gallery .item {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.gallery img:hover {
  opacity: 0.7;
}

.modal {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.8);
  display: flex;
  justify-content: center;
  align-items: center;
}

.modal img {
  max-width: 80%;
  max-height: 80%;
}
</style>
