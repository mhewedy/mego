import axios from 'axios'

export default {
    list: function (successFn, failFn) {
        axios.get("/api/v1/rooms/tree")
            .then(it => successFn(it.data))
            .catch(it => failFn(it))
    }
}
