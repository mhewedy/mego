export default {
    messages: [],
    error: function (content) {
        this.add('error', content)
    },
    info: function (content) {
        this.add('info', content)
    },
    add: function (severity, content) {
        let m = {severity: severity, content: content, key: new Date().getMilliseconds()};
        this.messages.push(m)
    }
}
