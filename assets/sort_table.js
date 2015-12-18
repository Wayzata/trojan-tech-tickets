$(function() {
	$("th").click(function() {
		$("th").removeClass("sortingBy");
		this.classList.add("sortingBy");
		var cellIndex = this.cellIndex;
		var e = $("tbody>tr");
		e.sort(function(a, b) {
			a = $("td:nth-child(" + (cellIndex+1) + ")", a).text();
			b = $("td:nth-child(" + (cellIndex+1) + ")", b).text();
			return a.localeCompare(b);
		});
		e.detach().appendTo($("tbody"));
	});
});
