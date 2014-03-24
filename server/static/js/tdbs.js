/*
	LOL ripped of my other site :D
*/

function NewInput() {
	// Start by showing a loading message on our modal
	$("#modal-body").html("Loading...")
	$.post(
		"/process", 
		{
			data: $("#flags").val()
		},
			function(data) {
				$("#modal-body").html(data)
	}); // God JS is disgusting
}

/* JS warning, may need it later, then again we're using modals :P
$(document).ready(function() {
	$("#jswarn").remove();
	$("#no-js-placeholder").remove();
	$("#js-placeholder").show();
}); */

$(document).ready(function() {
	$("myModal").modal(); // Just to be safe it loads.
});