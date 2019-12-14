<template>

  <div>

    <div style="padding-left: 40%">
      <h3>Username</h3>
      <div class="p-grid p-fluid">
        <div class="p-col-12 p-md-4">
          <div class="p-inputgroup">
            <span class="p-inputgroup-addon">
                <i class="pi pi-user"></i>
            </span>
            <InputText placeholder="Username" v-model="username"/>
          </div>
        </div>

      </div>
    </div>

    <div style="padding-left: 40%">
      <h3>Password</h3>
      <div class="p-grid p-fluid">
        <div class="p-col-12 p-md-4">
          <div class="p-inputgroup">
            <span class="p-inputgroup-addon">
                <i class="pi pi-key"></i>
            </span>
            <Password placeholder="Password" v-model="password"/>
          </div>
        </div>
      </div>
    </div>

    <div style="padding: 20px 0 0 40%;">
      <div class="p-grid p-fluid">
        <div class="p-col-12 p-md-4">
          <Button label="Login" @click="login()"/>
        </div>
      </div>
    </div>

  </div>

</template>

<script>
    import UsersService from '../services/users'
    import MessagesService from '../services/messages'

    export default {
        name: "Login",
        data() {
            return {
                username: null,
                password: null,
            }
        },
        methods: {
            login: function () {
                let user = {
                    username: this.username,
                    password: this.password,
                };
                UsersService.login(user, (it) => {
                    localStorage.setItem("mego_token", it.id_token);
                    this.$parent.token = it.id_token;
                    this.$http.defaults.headers.common['Authorization'] = 'bearer ' + it.id_token;
                }, it => MessagesService.error(it))
            }
        }
    }
</script>

<style scoped>

</style>
