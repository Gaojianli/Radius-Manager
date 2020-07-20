import Vue from 'vue';
import Vuex from 'vuex';

Vue.use(Vuex);

export default new Vuex.Store({
    state: {
        UserList: [],
        server: '',
        key: ''
    },
    mutations: {
        setServer(state, value) {
            state.server = value;
        },
        setKey(state, value) {
            state.key = value;
        },
        setUserList(state, value) {
            state.UserList = value;
        },
        removeUser(state, toRemove) {
            state.UserList = state.UserList.filter(user => user.username != toRemove.username);
        },
        addUser(state, newUser) {
            state.UserList.push(newUser);
        },
        changeUser(state, newUser) {
            const toChange = state.UserList.find(it => it.username === newUser.username);
            if(toChange)
                toChange.type = newUser.type;
        }
    },
    actions: {
        async getUserList({ state, commit }) {
            const result = await fetch(`${state.server}/user`, { credentials: 'include' })
                .then(response => response.json())
                .catch(error => console.log('error', error));
            commit("setUserList", result.data);
            return result.data;
        },
        async login({ commit }, { server, key }) {
            commit('setServer', server);
            commit('setKey', key);
            const urlencoded = new URLSearchParams();
            urlencoded.append("adminKey", key);
            const requestOptions = {
                method: 'POST',
                body: urlencoded,
                redirect: 'follow',
                credentials: 'include'
            };
            return await fetch(`${server}/login`, requestOptions)
                .then(response => response.json())
                .catch(error => console.log('error', error));
        },
        async getPassword({ state }, { username }) {
            const result = await fetch(`${state.server}/user/${username}/password`, { credentials: 'include' })
                .then(res => res.json())
                .catch(error => console.log('error', error));
            return result.code === 200 ? result.data : '';
        },
        async deleteUser({ state, commit }, user) {
            const result = await fetch(`${state.server}/user/${user.username}`, { credentials: 'include', method: 'DELETE' })
                .catch(error => console.log('error', error));
            if (result.status === 204)
                commit('removeUser', user);
            return result.status === 204;
        },
        async addUser({ state, commit }, toAdd) {
            const headers = new Headers();
            headers.append("Content-Type", "application/x-www-form-urlencoded");

            const urlencoded = new URLSearchParams();
            urlencoded.append("username", toAdd.username);
            urlencoded.append("password", toAdd.password);
            urlencoded.append("type", toAdd.type);

            const requestOptions = {
                method: 'PUT',
                headers: headers,
                body: urlencoded,
                redirect: 'follow',
                credentials: 'include'
            };

            const result = await fetch(`${state.server}/user`, requestOptions)
                .then(response => response.json())
                .catch(error => console.log('error', error));
            if (result.code === 201) {
                commit('addUser', {
                    username: toAdd.username,
                    type: toAdd.type
                })
            }
            return result;
        },
        async changePassword({ state, commit }, { user, newPass }) {
            const headers = new Headers();
            headers.append("Content-Type", "application/x-www-form-urlencoded");

            const urlencoded = new URLSearchParams();
            urlencoded.append("password", newPass);
            urlencoded.append("type", user.type);

            var requestOptions = {
                method: 'POST',
                headers: headers,
                body: urlencoded,
                redirect: 'follow',
                credentials: 'include'
            };
            const result = await fetch(`${state.server}/user/${user.username}/changepassword`, requestOptions)
                .then(response => response.json());
            if (result.code === 200) {
                commit('changeUser',user);
            }
            return result;
        }
    }
})