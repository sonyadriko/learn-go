<!-- ItemList.vue -->
<template>
    <div>
        <h2>Item List</h2>
        <ul>
          <li v-for="item in filteredItems" :key="item.id">
            {{ item.name }}
            <button @click="editItem(item)">Edit</button>
            <button @click="updateDelete(item)">Delete</button>
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
                searchTerm: '', // Tambahkan variabel searchTerm
                searchResults: [],
                deleteItemId: null, // Tambahkan variabel searchResults
            };
        },
        computed: {
          filteredItems() {
            // Filter items yang memiliki deleted_at null
            return this.items.filter(item => item.deleted_at === null);
          },
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
                this.selectedItem = {
                    ...item
                };
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
            updateDelete(item) {
                this.deleteItemId = item
                console.log(this.deleteItemId)

                axios.put(`http://localhost:8080/api/itemss/`+this.deleteItemId)
                    .then(response => {
                        const index = this.items.findIndex(i => i.id === response.data.id);
                        if (index !== -1) {
                            this.$set(this.items, index, response.data);
                        }
                        this.selectedItem = null;
                    })
                    .catch(error => {
                        console.error('Error updating item:', error);
                        // Handle error accordingly, for example, show a notification to the user
                    });
            },

            // deleteItem(itemId) {
            //   axios.delete(`http://localhost:8080/api/items/${itemId}`)
            //     .then(() => {
            //       this.items = this.items.filter(item => item.id !== itemId);
            //     })
            //     .catch(error => {
            //       console.error('Error deleting item:', error);
            //     });
            // },
        },
    };
</script>

<style scoped>
    /* Styling for the component */
</style>
