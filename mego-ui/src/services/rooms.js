import axios from 'axios'
import errors from "./errors";

export default {
    tree: function (successFn, failFn) {
        axios.get("/api/v1/rooms/tree")
            .then(it => successFn(it.data))
            .catch(it => {
                errors.handle401(it);
                failFn && failFn(it.response.data.error)
            })
    },
    list: function (successFn, failFn) {
        axios.get("/api/v1/rooms")
            .then(it => successFn(it.data))
            .catch(it => {
                errors.handle401(it);
                failFn && failFn(it.response.data.error)
            })
    }
}
