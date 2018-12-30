<template>
    <Menu mode="horizontal" theme="dark" active-name="1">
        <div class="layout-logo"></div>
        <div class="layout-nav">
            <MenuItem name="1">
                <router-link to="/"><Icon type="ios-home"></Icon>{{$t("message.home")}}</router-link>
            </MenuItem>
            <MenuItem name="2">
                <router-link to="/about"><Icon type="ios-help"></Icon>{{$t("message.about")}}</router-link>
            </MenuItem>
            <MenuItem name="3">
                <router-link to="/logout"><Icon type="ios-power"></Icon></router-link>
            </MenuItem>
        </div>
    </Menu>
</template>

<script>
    export default {
        name: "TopNav"
    }
</script>

<style scoped>
    a {
        color: white;
    }
    a:hover {
        color: lightgray;
    }
    .ivu-menu-item:last-child {
        float: right;
    }
    a i {
        margin-right: 0.25rem;
    }
</style>