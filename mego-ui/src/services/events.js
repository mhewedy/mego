import axios from 'axios'

export default {
    search: function (input, successFn, failFn) {
        axios.post("/api/v1/events/search", input)
            .then(it => successFn(it.data))
            .catch(it => failFn(it))
    }
}
