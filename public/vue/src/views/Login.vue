<template>
  <Layout>
    <Header></Header>
    <Content>
      <i-row>
        <i-col :span="24">
          <div id="login">
            <Card>
              <Alert type="error" show-icon v-if="status.error">{{status.error}}</Alert>

              <h1>{{$t("message.login")}}</h1>

              <Form inline>
                <FormItem>
                  <i-input type="text" name="username" v-model="input.username" placeholder="Username">
                    <Icon type="ios-person-outline" slot="prepend"></Icon>
                  </i-input>
                </FormItem>

                <FormItem>
                  <i-input type="password" name="password" v-model="input.password" placeholder="Password">
                    <Icon type="ios-lock-outline" slot="prepend"></Icon>
                  </i-input>
                </FormItem>

                <FormItem>
                  <Button type="primary" v-on:click="login()">{{$t("message.login")}}</Button>
                </FormItem>
              </Form>

              <div style="text-align: center;">
                <router-link to="/forgot-password">{{$t("message.forgot_password")}}</router-link>&nbsp;|&nbsp;
                <router-link to="/register">{{$t("message.register")}}</router-link>
              </div>
            </Card>
          </div>
        </i-col>
      </i-row>
    </Content>
    <Footer></Footer>
</template>

<script>
  export default {
    name: 'Login',
    data() {
      return {
        status: {
          error: ""
        },
        input: {
          username: "",
          password: ""
        }
      }
    },
    methods: {
      login() {
        if(this.input.username !== "" && this.input.password !== "") {

          // clear error, show spinner
          this.status.error = "";
          this.$Spin.show();

          this.$http.post('/api/login', "email=" + this.input.username + "&password=" + this.input.password, {
            headers: {
              'Content-Type': 'application/x-www-form-urlencoded'
            }
          }).then(response => {
            this.$Spin.hide();

            if(response.status === 200) {
              localStorage.setItem('user', response.bodyText);
              this.$emit("authenticated", true);
              this.$router.replace({ name: "home" });
            } else {
              this.status.error = "please try again"
            }
          }).catch(err => {
            this.$Spin.hide();
            this.status.error = "please try again"
          });
        } else {
          this.status.error = "a username and password must be present"
        }
      }
    }
  }
</script>

<style scoped>
  #login {
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
  #login .ivu-input-group .ivu-input {
    min-width: 200px;
  }
</style>