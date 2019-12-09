<template>
  <div id="app">

    <div v-show="messages">
      <Message v-for="msg of messages" :severity="msg.severity" :key="msg.key" :sticky="false">{{msg.content}}</Message>
    </div>

    <div>
      <Search @searched="sentSearchInput"></Search>
      <Result v-if="searchInput" :search-input="searchInput"></Result>
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
            messages: MessageService.messages
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
