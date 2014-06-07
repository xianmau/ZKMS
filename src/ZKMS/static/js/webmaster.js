 $(function() {
 	SetLayout();
 	$(window).resize(function() {
 		SetLayout();
 	});
 	SetNavActive();
 });

 // set layout
 function SetLayout() {
 	var wh = $(window).height()
 	var ww = $(window).width()
 	$('#leftside').height(wh)
 	$('#leftside').width(280)
 	$('#rightside').height(wh)
 	$('#rightside').width(ww - 280)
 }

 // set nav active
 function SetNavActive() {
 	var url = window.location.pathname;
 	var p = url.split("/")
 	for (var i = 0; i < p.length; i++) {
 		if (p[i].length == 0) p.splice(i, 1);
 	}
 	if (p.length == 1) {
 		$('.nav li a').eq(0).addClass("active");
 	} else if (p.length > 1) {
 		if (p[1] == "reports") {
 			$('.nav li a').eq(1).addClass("active");
 		} else if (p[1] == "zones") {
 			$('.nav li a').eq(2).addClass("active");
 		}else if (p[1] == "accounts") {
 			$('.nav li a').eq(3).addClass("active");
 		}else if (p[1] == "settings") {
 			$('.nav li a').eq(4).addClass("active");
 		}
 	}
 }