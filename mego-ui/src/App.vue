<template>
  <div id="app">
    <Search @searched="sentSearchInput"></Search>
    <Result v-if="searchInput" :search-input="searchInput"></Result>
  </div>
</template>

<script>
import Search from './components/Search.vue'
import Result from "./components/Result";

export default {
    data() {
        return {
            searchInput: null
        }
    },
    methods: {
        sentSearchInput: function (searchInput) {
            this.searchInput = searchInput
        }
    },
    components: {
      Search, Result
    }
}
</script>

<style scoped>
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
  margin-top: 60px;
}

.app-container {
  text-align: center;
}

body #app .p-button {
  margin-left: .2em;
}

form {
  margin-top: 2em;
}
</style>
