<template>
  <Layout>
    <Header>Header</Header>
    <Content>
      <i-row v-if="!status.registered">
        <i-col :span="24">
          <div id="register">
            <Card>
            <Alert type="error" show-icon v-if="status.error">{{status.error}}</Alert>

            <h1>{{$t("message.register")}}</h1>
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
                <Button type="primary" v-on:click="register()">{{$t("message.register")}}</Button>
              </FormItem>
            </Form>

            <div style="text-align: center;">
              <router-link to="/login">{{$t("message.login")}}</router-link>
            </div>
            </Card>
          </div>
        </i-col>
      </i-row>

      <i-row v-if="status.registered">
        <i-col :span="24">
          <Card>
            {{$t("message.registration_complete")}}
          </Card>
        </i-col>
      </i-row>
    </Content>
    <Footer></Footer>
  </Layout>
</template>

<script>
  export default {
    name: 'Register',
    data() {
      return {
        status: {
          registered: false,
          error: ''
        },
        input: {
          username: "",
          password: ""
        }
      }
    },
    methods: {
      register() {
        if(this.input.username !== "" && this.input.password !== "") {

          // clear error, show spinner
          this.status.error = "";
          this.$Spin.show();

          this.$http.post('/api/register', "email=" + this.input.username + "&password=" + this.input.password, {
            headers: {
              'Content-Type': 'application/x-www-form-urlencoded'
            }
          }).then(response => {
            this.$Spin.hide();

            if(response.status === 200) {
              localStorage.setItem('user', response.bodyText);
              this.status.registered = true;
            } else {
              this.status.error = "please try again"
            }
          }).catch(err => {
            this.$Spin.hide();
            this.status.error = "please try again, err = " + err;
          });
        } else {
          this.status.error = "An email and password must be present"
        }
      }
    }
  }
</script>

<style scoped>
  #register {
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
  #register .ivu-input-group .ivu-input {
    min-width: 190px;
  }
</style>