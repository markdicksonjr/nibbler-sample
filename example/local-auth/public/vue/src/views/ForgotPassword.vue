<template>
  <Layout>
    <Header>

    </Header>
    <Content>
      <i-row v-if="!status.sent">
        <i-col :span="24">
          <div id="forgot">
            <Card>
              <Alert type="error" show-icon v-if="status.error">{{status.error}}</Alert>

              <h1>{{$t("message.forgot_password")}}</h1>

              <Form inline>
                <FormItem>
                  <i-input type="text" name="email" v-model="input.email" placeholder="Email">
                    <Icon type="ios-person-outline" slot="prepend"></Icon>
                  </i-input>
                </FormItem>

                <FormItem>
                  <Button type="primary" v-on:click="submit()">{{$t("message.submit")}}</Button>
                </FormItem>
              </Form>

              <div style="text-align: center;">
                <router-link to="/login">{{$t("message.login")}}</router-link>
              </div>
            </Card>
          </div>
        </i-col>
      </i-row>

      <i-row v-if="status.sent">
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
    name: 'ForgotPassword',
    data() {
      return {
        status: {
          sent: false,
          error: ""
        },
        input: {
          email: ""
        }
      }
    },
    methods: {
      submit() {
        if(this.input.email !== "") {

          // clear error and status, show spinner
          this.status.error = "";
          this.status.sent = false;
          this.$Spin.show();

          this.$http.post('/api/password/reset-token', "email=" + this.input.email, {
            headers: {
              'Content-Type': 'application/x-www-form-urlencoded'
            }
          }).then(response => {
            this.$Spin.hide();

            if(response.status === 200) {
              this.status.sent = true;
            } else {
              this.status.error = "please try again"
            }
          }).catch(err => {
            this.$Spin.hide();
            this.status.error = "please try again, err = " + err;
          });
        } else {
          this.status.error = "an email must be present"
        }
      }
    }
  }
</script>

<style scoped>
  #forgot {
    width: 400px;
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
  #forgot .ivu-input-group .ivu-input {
    min-width: 200px;
  }
</style>