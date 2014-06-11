<!DOCTYPE html>

<html>
  <head>
    <title>Zookeeper Management System</title>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8">

    <link href="http://cdn.bootcss.com/bootstrap/3.1.1/css/bootstrap.min.css" rel="stylesheet">
    <link href="http://cdn.bootcss.com/jstree/3.0.0/package/dist/themes/default/style.css" rel="stylesheet">
    <link href="http://cdn.bootcss.com/jquery.nanoscroller/0.8.0/nanoscroller.min.css" rel="stylesheet">

    <link href="/static/css/webmaster.css" rel="stylesheet">
	</head>

  <body style="overflow-y:auto;overflow-x:hidden;">
    <div id="leftside">
      <div id="logo"></div>
      {{if .IsLogin}}
        <span class="login-stat">
          <a>Hi, {{.LoginName}}! </a>
          <a href="#">&lt;Exit&gt;</a>
        </span>
      {{else}}
        <span class="login-stat"><a href="/login">Log on</a></span>
      {{end}}
      <ul class="nav">
        <li><a href="/webmaster" class="">Dashboard</a></li>
        <li><a href="/webmaster/reports" class="">Reports</a></li>
        <li><a href="/webmaster/zones" class="">Zones</a></li>
        <li><a href="/webmaster/accounts" class="">Accounts</a></li>
        <li><a href="/webmaster/settings" class="">Settings</a></li>
      </ul>
      <span class="copyright">&copy 2014 YY.COM SDS GROUP</span>
    </div>

    <div id="rightside">
      {{.LayoutContent}}
    </div>

    <div class="clear"></div>

    <script src="https://code.jquery.com/jquery-1.11.1.min.js"></script>
    <script src="http://cdn.bootcss.com/bootstrap/3.1.1/js/bootstrap.min.js"></script>
    <script src="http://cdn.bootcss.com/jstree/3.0.0/jstree.min.js"></script>
    <script src="/static/js/jstree.data.json"></script>
    <script src="http://cdn.bootcss.com/jquery.nanoscroller/0.8.0/jquery.nanoscroller.min.js"></script>
    <script src="http://code.highcharts.com/stock/highstock.js"></script>
    <script src="http://code.highcharts.com/stock/modules/exporting.js"></script>

    <script src="/static/js/dragdiv.js"></script>
    
    <script src="/static/js/webmaster.js"></script>
  </body>
</html>
