import axios from 'axios'

export default {
    search: function (q, exclude, successFn, failFn) {
        axios.post("/api/v1/attendees/search?q=" + q, exclude)
            .then(it => successFn(it.data))
            .catch(it => failFn(it))
    },
    getPhoto: function (email, successFn, failFn) {
        axios.get("/api/v1/attendees/" + email + "/photo")
            .then(it => successFn(it.data))
            .catch(it => failFn(it))
    }
}
