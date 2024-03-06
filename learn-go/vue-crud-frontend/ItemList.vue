<template>
  <div>
    <h2>Item List</h2>
    <ul>
      <li v-for="item in items" :key="item.id">
        <router-link :to="{ path: '/item/' + item.id }">{{ item.name }}</router-link>
      </li>
    </ul>
  </div>
</template>

<script>
export default {
  data() {
    return {
      items: []
    };
  },
  created() {
    // Fetch items from the backend
    fetch('http://localhost:8080/api/items')
      .then(response => response.json())
      .then(data => {
        this.items = data;
      });
  }
};
</script>
