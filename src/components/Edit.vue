<template>
  <v-dialog
    v-model="this.editorSwitch"
    fullscreen
    hide-overlay
    transition="dialog-bottom-transition"
  >
    <v-card>
      <v-toolbar dark color="primary">
        <v-btn icon dark @click="$parent.editorDialog=false;">
          <v-icon>mdi-close</v-icon>
        </v-btn>
        <v-toolbar-title>{{isEdit?'编辑用户':'新建用户'}}</v-toolbar-title>
        <v-spacer></v-spacer>
        <v-toolbar-items>
          <v-btn dark text @click="save">Save</v-btn>
        </v-toolbar-items>
      </v-toolbar>
      <v-card flat class="mx-auto mt-6" max-width="600">
        <v-card-title>
          <span class="headline">编辑用户信息</span>
        </v-card-title>
        <v-card-text>
          <v-form ref="form">
            <v-text-field label="用户名" :disabled="isEdit" v-model="username" required></v-text-field>
            <v-text-field
              label="密码"
              v-model="password"
              :append-icon="show ? 'mdi-eye-off' : 'mdi-eye'"
              required
              :type="show?'text':'password'"
              @click:append="show = !show"
            ></v-text-field>
            <v-checkbox v-model="encrypted" label="使用MD5加密" required></v-checkbox>
            <v-alert v-if="encrypted" dense border="left" type="warning">
              开启加密后将无法查看密码，此外连接时
              <strong>EAP方法</strong>请选择
              <strong>TTLS</strong>，
              <strong>二阶段认证</strong>请选择
              <strong>PAP</strong>
            </v-alert>
          </v-form>
        </v-card-text>
      </v-card>
    </v-card>
  </v-dialog>
</template>
<script>
import CryptoJS from "crypto-js";
export default {
  name: "Editor",
  props: {
    isEdit: Boolean,
    editorSwitch: Boolean,
    toEdit: Object
  },
  data() {
    return {
      encrypted: false,
      password: "",
      realusername: "",
      show: false
    };
  },
  computed: {
    username: {
      get() {
        if (this.isEdit) return this.toEdit.username;
        else return this.realusername;
      },
      set(val) {
        if (!this.isEdit) this.realusername = val;
      }
    }
  },
  methods: {
    async save() {
      if (this.isEdit) {
        const result = await this.$store.dispatch("changePassword", {
          user: {
            username: this.toEdit.username,
            type: this.encrypted ? "MD5-Password" : "Cleartext-Password"
          },
          newPass: this.password
        });
        if (result.code == 200) {
          this.$parent.snackbarText = "修改成功";
          this.realusername = "";
          this.password = "";
          this.$parent.refresh();
          this.encrypted = false;
        } else this.$parent.snackbarText = "修改失败";
      } else {
        const result = await this.$store.dispatch("addUser", {
          username: this.realusername,
          password: this.encrypted
            ? CryptoJS.MD5(this.password).toString(CryptoJS.enc.Base64)
            : this.password,
          type: this.encrypted ? "MD5-Password" : "Cleartext-Password"
        });
        if (result.code == 201) {
          this.$parent.snackbarText = "添加成功";
          this.realusername = "";
          this.password = "";
          this.encrypted = false;
          this.$parent.refresh();
        } else this.$parent.snackbarText = "添加失败";
      }
      this.$parent.snackbar = true;
      this.$parent.editorDialog = false;
    }
  }
};
</script>