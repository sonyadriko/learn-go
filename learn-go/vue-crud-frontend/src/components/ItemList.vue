<!-- ItemList.vue -->
<template>
  <div>
    <h2>Item List</h2>
    <ul>
      <li v-for="item in items" :key="item.id">
        {{ item.name }}
        <button @click="editItem(item)">Edit</button>
        <button @click="deleteItem(item.id)">Delete</button>
      </li>
    </ul>

    <h2>Create New Item</h2>
    <form @submit.prevent="createItem">
      <label>Name:</label>
      <input v-model="newItemName" required />
      <button type="submit">Create</button>
    </form>

    <h2>Edit Item</h2>
    <form v-if="selectedItem" @submit.prevent="updateItem">
      <label>Name:</label>
      <input v-model="selectedItem.name" required />
      <button type="submit">Update</button>
    </form>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  data() {
    return {
      items: [],
      newItemName: '',
      selectedItem: null,
    };
  },
  mounted() {
    this.fetchItems();
  },
  methods: {
    async fetchItems() {
      try {
        const response = await axios.get('http://localhost:8080/api/items');
        this.items = response.data;
      } catch (error) {
        console.error('Error fetching items:', error);
      }
    },
    async createItem() {
      try {
        const response = await axios.post('http://localhost:8080/api/items', {
          name: this.newItemName,
        });
        this.items.push(response.data);
        this.newItemName = '';
      } catch (error) {
        console.error('Error creating item:', error);
      }
    },
    editItem(item) {
      this.selectedItem = { ...item };
    },
    async updateItem() {
      try {
        const response = await axios.put(
          `http://localhost:8080/api/items/${this.selectedItem.id}`,
          this.selectedItem
        );
        const index = this.items.findIndex((item) => item.id === response.data.id);
        if (index !== -1) {
          this.$set(this.items, index, response.data);
        }
        this.selectedItem = null;
      } catch (error) {
        console.error('Error updating item:', error);
      }
    },
    async deleteItem(itemId) {
      try {
        await axios.delete(`http://localhost:8080/api/items/${itemId}`);
        this.items = this.items.filter((item) => item.id !== itemId);
      } catch (error) {
        console.error('Error deleting item:', error);
      }
    },
  },
};
</script>

<style scoped>
/* Styling for the component */
</style>
