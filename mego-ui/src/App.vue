<template>
  <div id="app">

    <div v-show="messages">
      <Message v-for="msg of messages" :severity="msg.severity"
               :key="msg.key" :sticky="false">{{msg.content}}
      </Message>
    </div>

    <Login v-if="token == null"></Login>

    <div v-if="token != null">

      <div>
        <Search @search="search" :isResultLoading="isResultLoading"></Search>
        <Result @resultLoad="resultLoad" v-if="searchInput" :search-input.sync="searchInput"></Result>
      </div>
    </div>

  </div>
</template>

<script>
    import Search from './components/Search.vue'
    import Result from "./components/Result";
    import Login from "./components/Login";

    import MessageService from './services/messages'

    export default {
        data() {
            return {
                searchInput: null,
                isResultLoading: false,
                messages: MessageService.messages,
                token: localStorage.getItem("mego_token"),
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
            Search, Result, Login
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
