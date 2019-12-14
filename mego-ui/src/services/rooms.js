import axios from 'axios'

export default {
    tree: function (successFn, failFn) {
        axios.get("/api/v1/rooms/tree")
            .then(it => successFn(it.data))
            .catch(it => failFn && failFn(it.response.data.error))
    },
    list: function (successFn, failFn) {
        axios.get("/api/v1/rooms")
            .then(it => successFn(it.data))
            .catch(it => failFn && failFn(it.response.data.error))
    }
}
