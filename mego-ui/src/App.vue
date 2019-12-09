<template>
  <div id="app">
    <Search></Search>
    <Result></Result>
  </div>
</template>

<script>
import Search from './components/Search.vue'
import Result from "./components/Result";

export default {
    data() {
        return {
            message: null,
            text: null
        }
    },
    methods: {
        greet() {
            this.$toast.add({severity: 'info', summary: 'Hello '  + this.text});
            this.message = 'Hello ' + this.text;
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
