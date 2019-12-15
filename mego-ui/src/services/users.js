import axios from 'axios'

export default {
    login: function (user, successFn, failFn) {
        axios.post("/api/v1/login", user)
            .then(it => successFn(it.data))
            .catch(it => failFn && failFn(it.response.data.error))
    },
    logout: function (successFn, failFn) {
        axios.post("/api/v1/logout")
            .then(it => successFn(it.data))
            .catch(it => failFn && failFn(it.response.data.error))
    }
}
