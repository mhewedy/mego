<template>
  <div id="app">

    <div v-show="messages">
      <Message v-for="msg of messages" :severity="msg.severity" :key="msg.key" :sticky="false">{{msg.content}}</Message>
    </div>

    <div>
      <Search @search="search" :isResultLoading="isResultLoading"></Search>
      <Result @resultLoad="resultLoad" v-if="searchInput" :search-input.sync="searchInput"></Result>
    </div>
  </div>
</template>

<script>
import Search from './components/Search.vue'
import Result from "./components/Result";
import MessageService from './services/messages'

export default {
    data() {
        return {
            searchInput: null,
            isResultLoading: false,
            messages: MessageService.messages
        }
    },
    methods: {
        search: function (searchInput) {
            this.searchInput = searchInput
        },
        resultLoad: function (isResultLoading) {
            this.isResultLoading = isResultLoading
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
