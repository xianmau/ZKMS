 // 加载全局脚本
 $(function() {
   SetLayout();
   $(window).resize(function() {
     SetLayout();
   });
   SetNavActive();
 });

 // 按需加载脚本
 $(function() {
   var url = window.location.pathname;
   var p = url.split("/")
   for (var i = 0; i < p.length; i++) {
     if (p[i].length == 0) p.splice(i, 1);
   }
   if (p.length == 2) {
     if (p[1] == "zktree") {
       ZonesSetDragDiv()
       ZonesLoadZkTree();
     }
   } else if (p.length == 3) {
     if (p[1] == "reports" && p[2] == "brokerdetail") {
       ReportsLoadBrokerChart()
     } else if (p[1] == "reports" && p[2] == "loggerdetail") {
       ReportsLoadLoggerChart()
     } else if (p[1] == "reports" && (p[2] == "brokers" || p[2] == "loggers")) {
       ReportsDealFloatData()
     } else if (p[1] == "reports" && p[2] == "topics") {
       ReportsDealSegment()
     }
   }

 });

 /*----------------------------------
  * Layout scripts
  *----------------------------------*/
 function SetLayout() {
   var wh = $(window).height();
   var ww = $(window).width();
   $('#leftside').height(wh);
   $('#rightside').height(wh);
   $('#rightside').width(ww - 280);
 }

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
  * Zones scripts
  *----------------------------------*/
 function ZonesSetDragDiv() {
   // 设置结点信息层可拖曳
   $(".zktree-right").drag({
     handler: $('.dragdiv-title'),
     opacity: 0.7
   });
 }

 function ZonesLoadZkTree() {
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
           code += '<h4>PATH: </h4>';
           code += '<p class="znodeinfo-path">' + current_node + '</p>';
           code += '<h4>DATA: </h4>';
           code += '<p class="znodeinfo-data">';
           code += JSON.stringify(o, null, 4);
           code += '</p>';
           code += '';
           code += '';

           $('.node_show').html(code);
           $(".nano").nanoScroller();
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
 }


 /*----------------------------------
  * reports scripts
  *----------------------------------*/
 function ReportsDealFloatData() {
   $('.float-data').each(function(index) {
     var data = $(this).text()
     if (data == 0) {
       $(this).html("--")
     } else {

       $(this).html(Math.round((data * 100)) + "%")
     }
   })
 }

 function ReportsDealSegment() {
   // 处理段数据
   $('.segments-data').each(function() {
     var data = $(this).text();
     var o = jQuery.parseJSON(data);
     data = "";
     for (var i = 0; i < o.length; i++) {
       var tmp = "<strong>" + o[i].last_confirm_entry + "</strong><br />";
       tmp += "<ul>";
       for (var j = 0; j < o[i].loggers.length; j++) {
         tmp += "<li>" + o[i].loggers[j] + "</li>";
       }
       tmp += "</ul>";
       data += tmp;
     }
     var s = '<span class="label label-info" data-toggle="popover" title data-html="true" data-content="' + data + '" data-original-title="Segments Details">Segments Details</span>';
     $(this).html(s);
   });
   $('.segments-data .label').popover({
     placement: 'left'
   })
   $('.segments-data .label').popover('hide');
   $('.segments-data .label').click(function() {
     $('.segments-data .label').not(this).popover('hide');
   });
 }

 function ReportsLoadBrokerChart() {
   Highcharts.setOptions({
     global: {
       useUTC: true
     }
   });

   // Create the chart
   $('#container').highcharts('StockChart', {
     chart: {
       events: {
         load: function() {
           // set up the updating of the chart each second
           var series0 = this.series[0];
           var series1 = this.series[1];
           var series2 = this.series[2];
           setInterval(function() {
             $.get('/webmaster/reports/brokerdetail/getlatestdata?zoneid=' + $('#zoneid').val() + '&brokerid=' + $('#brokerid').val(), function(data) {
               series0.addPoint(jQuery.parseJSON(data)[0], true, true);
               series1.addPoint(jQuery.parseJSON(data)[1], true, true);
               series2.addPoint(jQuery.parseJSON(data)[2], true, true);
             });
           }, 60000);
         }
       }
     },

     rangeSelector: {
       buttons: [{
         count: 30,
         type: 'minute',
         text: '30M'
       }, {
         count: 1,
         type: 'hours',
         text: '1H'
       }, {
         type: 'all',
         text: 'All'
       }],
       inputEnabled: false,
       selected: 0
     },

     title: {
       text: 'Performance Evaluation of Broker: ' + $('#brokerid').val()
     },

     subtitle: {
       text: 'from zone: ' + $('#zoneid').val()
     },

     yAxis: {
       title: {
         text: 'percentage (%)'
       }
     },

     exporting: {
       enabled: false
     },

     tooltip: {
       shared: true
     },

     legend: {
       enabled: true,
       align: 'right',
       backgroundColor: '#FCFFC5',
       borderColor: 'black',
       borderWidth: 0,
       layout: 'vertical',
       verticalAlign: 'top',
       y: 0,
       shadow: true,
       floating: true
     },

     series: [{
       name: 'Cpu Rate(%)',
       data: jQuery.parseJSON($('#cpuData').val())
     }, {
       name: 'Net Rate(%)',
       data: jQuery.parseJSON($('#netData').val())
     }, {
       name: 'Disk Rate(%)',
       data: jQuery.parseJSON($('#diskData').val())
     }]
   });
 }

 function ReportsLoadLoggerChart() {
   Highcharts.setOptions({
     global: {
       useUTC: true
     }
   });

   // Create the chart
   $('#container').highcharts('StockChart', {
     chart: {
       events: {
         load: function() {
           // set up the updating of the chart each second
           var series0 = this.series[0];
           var series1 = this.series[1];
           var series2 = this.series[2];
           setInterval(function() {
             $.get('/webmaster/reports/loggerdetail/getlatestdata?zoneid=' + $('#zoneid').val() + '&loggerid=' + $('#loggerid').val(), function(data) {
               series0.addPoint(jQuery.parseJSON(data)[0], true, true);
               series1.addPoint(jQuery.parseJSON(data)[1], true, true);
               series2.addPoint(jQuery.parseJSON(data)[2], true, true);
             });
           }, 60000);
         }
       }
     },

     rangeSelector: {
       buttons: [{
         count: 30,
         type: 'minute',
         text: '30M'
       }, {
         count: 1,
         type: 'hours',
         text: '1H'
       }, {
         type: 'all',
         text: 'All'
       }],
       inputEnabled: false,
       selected: 0
     },

     title: {
       text: 'Performance Evaluation of Logger: ' + $('#loggerid').val()
     },

     subtitle: {
       text: 'from zone: ' + $('#zoneid').val()
     },

     yAxis: {
       title: {
         text: 'percentage (%)'
       }
     },

     exporting: {
       enabled: false
     },

     tooltip: {
       shared: true
     },

     legend: {
       enabled: true,
       align: 'right',
       backgroundColor: '#FCFFC5',
       borderColor: 'black',
       borderWidth: 0,
       layout: 'vertical',
       verticalAlign: 'top',
       y: 0,
       shadow: true,
       floating: true
     },

     series: [{
       name: 'Cpu Rate(%)',
       data: jQuery.parseJSON($('#cpuData').val())
     }, {
       name: 'Net Rate(%)',
       data: jQuery.parseJSON($('#netData').val())
     }, {
       name: 'Disk Rate(%)',
       data: jQuery.parseJSON($('#diskData').val())
     }]
   });
 }