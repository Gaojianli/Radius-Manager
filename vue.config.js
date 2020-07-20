module.exports = {
  "transpileDependencies": [
    "vuetify"
  ],
  pluginOptions:{
    electronBuilder:{
      builderOptions:{
        "appId":"me.gaojianli.radiusmanager",
        "productName":"Radius Manager",
        "copyright":"Copyright © Gaojianli 2020",
        "directories":{
          "output":"./build"
        },
        "win":{
          "target":{
            "target":"nsis",
            "arch":[
              "x64"
            ]
          },
          "icon":'resources/icon.ico'
        },
        "nsis":{
          "oneClick":false,
          "allowToChangeInstallationDirectory": true,
        }
      }
    }
  }
}