(function(g) {
	g.error = function error(text) {
		$(".modal").modal("hide");
		$("body").append(
			$("<div>").addClass("alert alert-danger fade in").append(
				$("<a>").attr({
					"href": "#",
					"class": "close",
					"data-dismiss": "alert"
				}).html("&times;"),
				$("<strong>").text("Error:"),
				text
			)
		);
	};
})(this);
