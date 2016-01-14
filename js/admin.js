function changePoints(p) {
  var type;
  if (p.id.startsWith("text_")) {
    type = "text";
  } else if (p.id.startsWith("area_")) {
    type = "area";
  } else {
    return;
  }

  $.post("points.html",
  {
    operation: "update",
    task: document.getElementById("task").value,
    type: type,
    id: p.id.substring(5),
    value: p.value
  },
  function(data, status){
  });
}
