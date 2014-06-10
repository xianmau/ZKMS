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
 		} else if (p[1] == "zones" || p[1] == "zktree") {
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



 /*----------------------------------
  * ZkTree scripts
  *----------------------------------*/
 $(function() {
 	$('#zktree_div')
 		.on('init.jstree', function() {})
 		.on('select_node.jstree', function(e, data) {
 			var current_node = "/" + data.instance.get_path(data.selected[0], '/', false);
 			$('.node_path').html(current_node);
 			$.ajax({
 				url: '/webmaster/zktree/getdata',
 				data: {
 					'zoneid': $('#cur_zone').val(),
 					'znode': current_node
 				},
 				success: function(data) {
 					var code = "";
 					var o = data;
 					if (o == null) {
 						o = [];
 					}
 					code += '<h5>PATH: </h5>';
 					code += '<p style="color:green;">' + current_node + '</p>';
 					code += '<h5>DATA: </h5>';
 					code += '<p style="white-space:pre;color:yellow;">';
 					code += JSON.stringify(o, null, 4);
 					code += '</p>';
 					code += '';
 					code += '';

 					$('.node_show').html(code);

 				},
 				dataType: 'json'
 			});
 		})
 		.jstree({
 			'core': {
 				'data': {
 					'dataType': 'json',
 					'url': function(node) {
 						return node.id === '#' ?
 							'/static/js/jstree.data.json' :
 							'/webmaster/zktree/getchildren';
 					},
 					'data': function(node) {
 						if (node.id == '#') {
 							return;
 						}
 						var nd = $.jstree.reference('#zktree_div');
 						return {
 							'zoneid': $('#cur_zone').val(),
 							'znode': nd.get_path(node, '/', false)
 						};
 					}
 				}
 			},
 			// 其它一些参数设置
 			'plugins': ["sort"]
 		});
 });