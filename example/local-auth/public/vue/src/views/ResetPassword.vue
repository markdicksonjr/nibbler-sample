<template>
  <Layout>
    <Header></Header>
    <Content>
      <i-row>
        <i-col :span="24">
          <div id="reset">
            <Card v-if="status.token && !status.changed">
              <Alert type="error" show-icon v-if="status.error">{{status.error}}</Alert>

              <h1>{{$t("message.forgot_password")}}</h1>

              <Form inline>
                <FormItem>
                  <i-input type="password" name="password" v-model="input.password" placeholder="Password">
                    <Icon type="ios-lock-outline" slot="prepend"></Icon>
                  </i-input>
                </FormItem>

                <FormItem>
                  <i-input type="password" name="confirm" v-model="input.confirm" placeholder="Confirm Password">
                    <Icon type="ios-lock-outline" slot="prepend"></Icon>
                  </i-input>
                </FormItem>

                <FormItem>
                  <Button type="primary" v-on:click="changePassword()">{{$t("message.change")}}</Button>
                </FormItem>
              </Form>

              <div style="text-align: center;">
                <router-link to="/login">{{$t("message.login")}}</router-link>
              </div>
            </Card>

            <Card v-if="status.changed">
              You've successfully changed your password, you may proceed to <router-link to="/login">{{$t("message.login_verb")}}</router-link>
            </Card>

            <Card v-if="status.error && !status.token">
              <Alert type="error" show-icon>{{status.error}}</Alert>
            </Card>
          </div>
        </i-col>
      </i-row>
    </Content>
    <Footer></Footer>
</template>

<script>
  export default {
    name: 'ResetPassword',
    data() {
      return {
        status: {
          error: "",
          token: "",
          changed: false
        },
        input: {
          password: "",
          confirm: ""
        }
      }
    },
    mounted: function() {
      this.status.token = this.$route.query.token
      if(!this.status.token) {
        this.status.error = "no token was provided";
      }
    },
    methods: {
      changePassword() {
        if(this.input.password === "" || this.input.confirm === "") {
          this.status.error = "a username and password must be present"
          return
        }

        if(this.input.password !== this.input.confirm) {
          this.status.error = "passwords do not match"
          return
        }

        // clear error, show spinner
        this.status.changed = false
        this.status.error = ""
        this.$Spin.show()

        this.$http.post('/api/password', "token=" + this.status.token + "&password=" + this.input.password, {
          headers: {
            'Content-Type': 'application/x-www-form-urlencoded'
          }
        }).then(response => {
          this.$Spin.hide()

          if(response.status === 200) {
            this.status.changed = true
          } else {
            this.status.error = "please try again"
          }
        }).catch(err => {
          this.$Spin.hide()
          this.status.error = "please try again"
        });
      }
    }
  }
</script>

<style scoped>
  #reset {
    width: 640px;
    margin-top: 200px;
    padding: 20px;
    margin-left: auto;
    margin-right: auto;
  }
  h1 {
    margin-bottom: 1rem;
  }
</style>

<style>
  #reset .ivu-input-group .ivu-input {
    min-width: 200px;
  }
</style>