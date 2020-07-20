<template>
  <v-container>
    <v-card height="80%" width="80%" class="mx-auto">
      <v-list subheader>
        <v-toolbar dense flat>
          <v-toolbar-title>用户列表</v-toolbar-title>
          <v-spacer></v-spacer>
          <v-tooltip top>
            <template v-slot:activator="{ on, attrs }">
              <v-btn icon v-bind="attrs" v-on="on" @click="addUser">
                <v-icon color="grey lighten-1">mdi-plus</v-icon>
              </v-btn>
            </template>
            <span>添加用户</span>
          </v-tooltip>
        </v-toolbar>
        <v-list-item v-for="user in userList" :key="user.username">
          <v-list-item-avatar>
            <v-icon>mdi-account</v-icon>
          </v-list-item-avatar>
          <v-list-item-content>
            <v-list-item-title v-text="user.username"></v-list-item-title>
          </v-list-item-content>
          <v-list-item-action v-if="user.type=='Cleartext-Password'" class="mr-6">
            <v-tooltip top>
              <template v-slot:activator="{ on, attrs }">
                <v-btn icon @click="getPassword(user)" v-bind="attrs" v-on="on">
                  <v-icon color="grey lighten-1">mdi-eye</v-icon>
                </v-btn>
              </template>
              <span>显示密码</span>
            </v-tooltip>
          </v-list-item-action>
          <v-list-item-action>
            <v-tooltip top>
              <template v-slot:activator="{ on, attrs }">
                <v-btn icon color="cyan" v-bind="attrs" v-on="on" @click="editUser(user)">
                  <v-icon color="cyan lighten-1">mdi-pencil</v-icon>
                </v-btn>
              </template>
              <span>修改密码</span>
            </v-tooltip>
          </v-list-item-action>
          <v-list-item-action>
            <v-tooltip top>
              <template v-slot:activator="{ on, attrs }">
                <v-btn
                  icon
                  v-bind="attrs"
                  v-on="on"
                  color="red"
                  @click="deleteDialog=true;currentUser=user"
                >
                  <v-icon color="red lighten-1">mdi-delete</v-icon>
                </v-btn>
              </template>
              <span>删除用户</span>
            </v-tooltip>
          </v-list-item-action>
        </v-list-item>
        <v-dialog persistent v-model="passwordDialog" max-width="300px">
          <v-card>
            <v-card-title class="headline">查看密码</v-card-title>
            <v-card-text>密码为：{{currentUser.password}}</v-card-text>
            <v-card-actions>
              <v-spacer></v-spacer>
              <v-btn color="green darken-1" text @click="closeDialog">OK</v-btn>
            </v-card-actions>
          </v-card>
        </v-dialog>
        <v-dialog v-model="deleteDialog" max-width="290">
          <v-card>
            <v-card-title class="headline">确定要删除此用户？</v-card-title>
            <v-card-text>注意：此操作不可逆！</v-card-text>
            <v-card-actions>
              <v-spacer></v-spacer>
              <v-btn color="blue darken-1" text @click="deleteDialog = false">取消</v-btn>
              <v-btn color="red darken-1" text @click="deleteUser(currentUser)">删除</v-btn>
            </v-card-actions>
          </v-card>
        </v-dialog>
      </v-list>
      <v-snackbar v-model="snackbar" timeout="2000" light>
        {{ snackbarText }}
        <template v-slot:action="{ attrs }">
          <v-btn color="red" text v-bind="attrs" @click="snackbar = false">关闭</v-btn>
        </template>
      </v-snackbar>
    </v-card>
    <Edit :editorSwitch='editorDialog' :isEdit='isEdit' :toEdit='toEditUser'></Edit>
  </v-container>
</template>
<script>
import Edit from './Edit'
export default {
  data() {
    return {
      userList: [],
      passwordDialog: false,
      deleteDialog: false,
      currentUser: {
        username: "",
        password: ""
      },
      toEditUser:null,
      snackbar: false,
      snackbarText: "",
      editorDialog:false,
      isEdit:false
    };
  },
  components:{
    Edit
  }
  ,
  async mounted() {
    this.userList = await this.$store.dispatch("getUserList");
  },
  methods: {
    async getPassword(user) {
      if (user.username !== this.currentUser.username) {
        const currentPassword = await this.$store.dispatch("getPassword", {
          username: user.username
        });
        this.currentUser.password = currentPassword;
        this.currentUser.username = user.username;
      }
      if (!this.currentUser.password == "") this.passwordDialog = true;
    },
    closeDialog() {
      this.passwordDialog = false;
    },
    async deleteUser(user) {
      const result = await this.$store.dispatch("deleteUser", user);
      if (result) {
        this.refresh();
        this.snackbar = true;
        this.snackbarText = "删除成功！";
      }
      this.deleteDialog = false;
    },
    editUser(user){
      this.isEdit=true;
      this.toEditUser=user;
      this.editorDialog=true;
    },
    addUser(){
      this.isEdit=false;
      this.editorDialog=true;
      this.toEditUser={
        username: "",
        password: ""
      };
    },
    refresh(){
      this.userList=this.$store.state.UserList;
    }
  }
};
</script>