<template>
    <Layout>
        <Header>Header</Header>
        <Content>
            <i-row>
                <i-col :span="24">
                    <Alert type="error" show-icon v-if="status.error">{{status.error}}</Alert>
                    <div v-if="status.verified">
                        <Card>
                            {{$t("message.verification_complete_prefix")}} <router-link to="/login">{{$t("message.login_verb")}}</router-link>
                        </Card>
                    </div>
                </i-col>
            </i-row>
        </Content>
    <Footer></Footer>
</template>

<script>
    export default {
        name: "VerifyEmail",
        data() {
            return {
                status: {
                    error: '',
                    verified: false
                }
            }
        },
        mounted: function() {
            if(!this.$route.query.token) {
                this.status.error = "no token was provided";
                return;
            }

            // clear error, show spinner
            this.status.error = "";
            this.$Spin.show();

            this.status.error = "";
            this.$http.post('/api/email/validate', "token=" + this.$route.query.token, {
                headers: {
                    'Content-Type': 'application/x-www-form-urlencoded'
                }
            }).then(response => {
                this.$Spin.hide();

                if(response.status === 200) {
                    this.status.verified = true;
                } else {
                    this.status.error = "please try again"
                }
            }).catch(err => {
                this.$Spin.hide();
                this.status.error = "please try again"
            });
        }
    }
</script>

<style scoped>

</style>