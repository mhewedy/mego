<template>

  <div>

    <div style="padding: 20px 0 20px 39%; font-size: 30px">
      <span style="font-weight: bold">MEGO</span> <span> The Meeting Organizer</span>
    </div>

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
            <Password placeholder="Password" v-model="password" :feedback="false"/>
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

    <div v-if="loadingResult" class="p-grid">
      <div class="p-col-5"></div>
      <div class="p-col-2" style="text-align: center;">
        <ProgressSpinner mode="indeterminate"/>
      </div>
      <div class="p-col-5"></div>
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
                loadingResult: false
            }
        },
        methods: {
            login: function () {
                let user = {
                    username: this.username,
                    password: this.password,
                };
                this.loadingResult = true;
                UsersService.login(user, (it) => {

                    localStorage.setItem("mego_token", it.id_token);
                    this.$parent.token = it.id_token;
                    this.$http.defaults.headers.common['Authorization'] = 'bearer ' + it.id_token;
                }, it => {
                    this.loadingResult = false;
                    MessagesService.error(it)
                })
            }
        }
    }
</script>

<style scoped>

</style>
