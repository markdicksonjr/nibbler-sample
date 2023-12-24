<template>
    <Alert type="error" show-icon v-if="status.error">{{status.error}}</Alert>
</template>

<script>
    export default {
        name: "Logout",
        data() {
            return {
                status: {
                    error: ""
                }
            }
        },
        mounted: function() {
            this.$Spin.show();
            this.$http.get('/api/logout').then(response => {
                this.$Spin.hide();

                if(response.status === 200) {
                    localStorage.removeItem('user');
                    this.$router.replace({ name: "login" });
                } else {
                    this.status.error = "please try again"
                }
            }).catch(err => {
                this.$Spin.hide();
                this.status.error = "please try again, err = " + err;
            });
        }
    }
</script>

<style scoped>

</style>