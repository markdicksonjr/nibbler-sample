<template>
  <Layout>
    <Header><TopNav></TopNav></Header>
    <Content>
      <h1>{{$t("message.welcome")}}</h1>
    </Content>
</template>

<style scoped>
  .ivu-layout-content {
    margin: 1rem;
  }
</style>

<script>
import TopNav from "../components/TopNav";

export default {
  name: 'home',
  components: {
    TopNav
  }
}
</script>
