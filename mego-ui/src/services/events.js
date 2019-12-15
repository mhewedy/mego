import axios from 'axios'
import errors from "./errors";

export default {
    search: function (input, successFn, failFn) {
        axios.post("/api/v1/events/search", input)
            .then(it => successFn(it.data))
            .catch(it => {
                errors.handle401(it);
                failFn && failFn(it.response.data.error)
            })
    },
    create: function (input, successFn, failFn) {
        axios.post("/api/v1/events/create", input)
            .then(() => successFn())
            .catch(it => {
                errors.handle401(it);
                failFn && failFn(it.response.data.error)
            })
    },
}
