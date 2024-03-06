<template>
  <div>
    <h2>Item Detail</h2>
    <p>ID: {{ item.id }}</p>
    <p>Name: {{ item.name }}</p>
    <p>Price: {{ item.price }}</p>
  </div>
</template>

<script>
export default {
  data() {
    return {
      item: {}
    };
  },
  created() {
    // Fetch item details from the backend based on the route parameter
    const itemId = this.$route.params.id;
    fetch(`http://localhost:8080/api/items/${itemId}`)
      .then(response => response.json())
      .then(data => {
        this.item = data;
      });
  }
};
</script>
