<template>
  <v-container class="fill-height" fluid>
    <v-row align="center" justify="center">
      <v-col cols="12" sm="8" md="4">
        <v-card class="elevation-12">
          <v-toolbar color="primary" dark flat>
            <v-toolbar-title>登录</v-toolbar-title>
            <v-spacer></v-spacer>
          </v-toolbar>
          <v-card-text>
            <v-form>
              <v-text-field
                label="Server"
                autocomplete
                name="Server"
                prepend-icon="mdi-web"
                v-model="server"
                type="text"
              ></v-text-field>
              <v-text-field
                id="adminKey"
                v-bind:label="label"
                name="adminKey"
                autocomplete
                prepend-icon="mdi-lock"
                v-model="key"
                type="password"
              ></v-text-field>
            </v-form>
          </v-card-text>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn color="primary" @click="login">登录</v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>
<script>
export default {
  data() {
    return {
      key: "",
      cached: false,
      server: ""
    };
  },
  methods: {
    async login() {
      localStorage.setItem("server", this.server);
      let key = this.key;
      if (key == "" && this.cached) key = localStorage.getItem("adminKey");
      localStorage.setItem("adminKey", key);
      const result = await this.$store.dispatch("login", {
        server: this.server,
        key: key
      });
      if (result.code == 200) this.$router.push({ path: "/main" });
    }
  },
  mounted() {
    let adminKey = localStorage.getItem("adminKey");
    this.adminKey = adminKey ?? "";
    let server = localStorage.getItem("server");
    if (adminKey) this.cached = true;
    this.server = server ?? "";
  },
  computed: {
    label() {
      return `Key${this.cached ? "(已缓存)" : ""}`;
    }
  }
};
</script>