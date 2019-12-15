import axios from 'axios'
import errors from "./errors";

export default {
    search: function (q, exclude, successFn, failFn) {
        axios.post("/api/v1/attendees/search?q=" + q, exclude)
            .then(it => successFn(it.data))
            .catch(it => {
                errors.handle401(it);
                failFn && failFn(it.response.data.error)
            })
	
    },
    getPhoto: function (email, successFn, failFn) {
        /*axios.get("/api/v1/attendees/" + email + "/photo")
            .then(it => successFn(it.data))
            .catch(it => {
                errors.handle401(it);
                failFn && failFn(it.response.data.error)
            })*/
    }
}
