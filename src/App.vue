<template>
  <v-app>
    <v-app-bar app color="primary" dark>
      <div class="d-flex align-center">
        <span
          class="shrink mt-1"
          contain
          min-width="100"
          style="font-size:2em"
          width="100"
        >Radius Manager</span>
      </div>
      <v-spacer></v-spacer>
    </v-app-bar>
    <v-main>
      <router-view></router-view>
    </v-main>
  </v-app>
</template>

<script>
export default {
  name: "App",

  data: () => ({
    //
  }),
  created() {
    localStorage.getItem("state") &&
      this.$store.replaceState(
        Object.assign(
          this.$store.state,
          JSON.parse(localStorage.getItem("state"))
        )
      );
    window.addEventListener("beforeunload", () => {
      localStorage.setItem("state", JSON.stringify(this.$store.state));
    });
  }
};
</script>
