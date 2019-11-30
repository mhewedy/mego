var data =
    '[{ "value": 1, "text": "Task 1", "continent": "Task" }, { "value": 2, "text": "Task 2", "continent": "Task" }, { "value": 3, "text": "Task 3", "continent": "Task" }, { "value": 4, "text": "Task 4", "continent": "Task" }, { "value": 5, "text": "Task 5", "continent": "Task" }, { "value": 6, "text": "Task 6", "continent": "Task" } ]';

//get data pass to json
var task = new Bloodhound({
    datumTokenizer: Bloodhound.tokenizers.obj.whitespace("text"),
    queryTokenizer: Bloodhound.tokenizers.whitespace,
    local: jQuery.parseJSON(data) //your can use json type
});

task.initialize();

var elt = $("#tag1");
elt.tagsinput({
    itemValue: "value",
    itemText: "text",
    typeaheadjs: {
        name: "task",
        displayKey: "text",
        source: task.ttAdapter()
    }
});

//insert data to input in load page
elt.tagsinput("add", {
    value: 1,
    text: "task 1",
    continent: "task"
});
