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

    <h2>Search Items</h2>
    <input v-model="searchTerm" placeholder="Search items..." />
    <button @click="searchItems">Search</button>

    <!-- Display search results -->
    <h2>Search Results</h2>
    <ul>
      <li v-for="result in searchResults" :key="result.id">
        {{ result.name }}
      </li>
    </ul>
      
    
    
    
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
      searchTerm: '', // Tambahkan variabel searchTerm
      searchResults: [], // Tambahkan variabel searchResults
    };
  },
  mounted() {
    this.fetchItems();
  },
  methods: {
    fetchItems() {
      axios.get('http://localhost:8080/api/items')
        .then(response => {
          this.items = response.data;
          // console.log(this.items);
        })
        .catch(error => {
          console.error('Error fetching items:', error);
        });
    },
    searchItems() {
      axios.get(`http://localhost:8080/api/items?search=${this.searchTerm}`)
        .then(response => {
          this.searchResults = response.data;
        })
        .catch(error => {
          console.error('Error searching items:', error);
        });
    },
    createItem() {
      axios.post('http://localhost:8080/api/items', {
        name: this.newItemName,
      })
      .then(response => {
        this.items.push(response.data);
        this.newItemName = '';
      })
      .catch(error => {
        console.error('Error creating item:', error);
      });
    },
    editItem(item) {
      this.selectedItem = { ...item };
    },
    updateItem() {
      axios.put(`http://localhost:8080/api/items/${this.selectedItem.id}`, this.selectedItem)
        .then(response => {
          const index = this.items.findIndex(item => item.id === response.data.id);
          if (index !== -1) {
            this.items, index, response.data;
          window.location.reload();

          }
          this.selectedItem = null;

        })
        .catch(error => {
          console.error('Error updating item:', error);
        });
    },
    deleteItem(itemId) {
      axios.delete(`http://localhost:8080/api/items/${itemId}`)
        .then(() => {
          this.items = this.items.filter(item => item.id !== itemId);
        })
        .catch(error => {
          console.error('Error deleting item:', error);
        });
    },
  },
};
</script>

<style scoped>
/* Styling for the component */
</style>
