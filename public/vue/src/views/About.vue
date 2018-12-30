<template>
  <Layout>
    <Header><TopNav></TopNav></Header>
    <Content>
      <h1>{{$t("message.about")}}</h1>
    </Content>
  </Layout>
</template>

<style scoped>
  .ivu-layout-content {
    margin: 1rem;
  }
</style>

<script>
import TopNav from "../components/TopNav";

export default {
  name: 'about',
  components: {
    TopNav
  }
}
</script>
