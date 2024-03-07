<!-- ItemDetail.vue -->
<template>
  <div>
    <h2>Item Detail</h2>
    <p>ID: {{ item.id }}</p>
    <p>Name: {{ item.name }}</p>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  data() {
    return {
      item: null,
    };
  },
  mounted() {
    this.fetchItem();
  },
  methods: {
    async fetchItem() {
      const itemId = this.$route.params.id;
      try {
        const response = await axios.get(`http://localhost:8080/api/items/${itemId}`);
        this.item = response.data;
      } catch (error) {
        console.error('Error fetching item:', error);
      }
    },
  },
};
</script>

<style scoped>
/* Styling for the component */
</style>
