 /*----------------------------------
  * Layout scripts
  *----------------------------------*/
 $(function() {
 	SetLayout();
 	$(window).resize(function() {
 		SetLayout();
 	});
 	SetNavActive();
 });

 // set layout
 function SetLayout() {
 	var wh = $(window).height();
 	var ww = $(window).width();
 	$('#leftside').height(wh);
 	$('#rightside').height(wh);
 	$('#rightside').width(ww - 280);
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
 		} else if (p[1] == "accounts") {
 			$('.nav li a').eq(3).addClass("active");
 		} else if (p[1] == "settings") {
 			$('.nav li a').eq(4).addClass("active");
 		}
 	}
 }



 /*----------------------------------
  * Settings scripts
  *----------------------------------*/
 function SettingsShowCreateForm() {
 	$('.create-form').toggle('slow');
 }

 function SettingsCancel() {
 	$('#createzoneid').val("");
 	$('#createhosts').val("");
 	$('.create-form').toggle('slow');
 }

 function SettingsCreateZone() {
 	var createzoneid = $('#createzoneid').val();
 	var createhosts = $('#createhosts').val();

 	$.ajax({
 		type: 'POST',
 		url: '/webmaster/settings/CreateZone',
 		data: {
 			createzoneid: createzoneid,
 			createhosts: createhosts
 		},
 		success: function(data) {
 			location.reload();
 		},
 		dataType: 'text'
 	});
 }

 function SettingsEditZone(zoneid) {
 	var edithosts = $('#edithosts_' + zoneid).val();

 	$.ajax({
 		type: 'POST',
 		url: '/webmaster/settings/EditZone',
 		data: {
 			editzoneid: zoneid,
 			edithosts: edithosts
 		},
 		success: function(data) {
 			location.reload();
 		},
 		dataType: 'text'
 	});
 }

 function SettingsDeleteZone(zoneid) {
 	if (!confirm("Are your sure?")) {
 		return;
 	}
 	$.ajax({
 		type: 'POST',
 		url: '/webmaster/settings/DeleteZone',
 		data: {
 			deletezoneid: zoneid
 		},
 		success: function(data) {
 			location.reload();
 		},
 		dataType: 'text'
 	});
 }

 /*----------------------------------
  * Accounts scripts
  *----------------------------------*/
 function AccountsShowCreateForm() {
 	$('.create-form').toggle('slow');
 }

 function AccountsCancel() {
 	$('#createname').val("");
 	$('#createpassword').val("");
 	$('#createremark').val("");
 	$('.create-form').toggle('slow');
 }

 function AccountsCreate() {
 	var createname = $('#createname').val();
 	var createpassword = $('#createpassword').val();
 	var createremark = $('#createremark').val();

 	$.ajax({
 		type: 'POST',
 		url: '/webmaster/accounts/Create',
 		data: {
 			createname: createname,
 			createpassword: createpassword,
 			createremark: createremark
 		},
 		success: function(data) {
 			location.reload();
 		},
 		dataType: 'text'
 	});
 }

 function AccountsEdit(name) {
 	var editpassword = $('#editpassword_' + name).val();
 	var editremark = $('#editremark_' + name).val();

 	$.ajax({
 		type: 'POST',
 		url: '/webmaster/accounts/Edit',
 		data: {
 			editname: name,
 			editpassword: editpassword,
 			editremark: editremark
 		},
 		success: function(data) {
 			location.reload();
 		},
 		dataType: 'text'
 	});
 }

 function AccountsDelete(name) {
 	if (!confirm("Are your sure?")) {
 		return;
 	}
 	$.ajax({
 		type: 'POST',
 		url: '/webmaster/accounts/Delete',
 		data: {
 			deletename: name
 		},
 		success: function(data) {
 			location.reload();
 		},
 		dataType: 'text'
 	});
 }