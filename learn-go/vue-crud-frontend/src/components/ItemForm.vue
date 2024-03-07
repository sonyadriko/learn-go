<template>
  <form @submit.prevent="submitForm">
    <!-- Form fields -->
    <button type="submit">Submit</button>
  </form>
</template>

<script>
import axios from 'axios';

export default {
  data() {
    return {
      // Form data
    };
  },
  methods: {
    submitForm() {
      axios.post('http://localhost:8080/api/items', this.formData)
        .then(response => {
          // Handle successful submission
        })
        .catch(error => {
          console.error('Error submitting form:', error);
        });
    }
  }
};
</script>
