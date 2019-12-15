export default {
    handle401: function (error) {
        if (error.response.status === 401) {
            localStorage.removeItem("mego_token");
            window.location.reload();
        }
    }
}
