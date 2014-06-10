<!DOCTYPE html>

<html>
  <head>
    <title>Zookeeper Management System</title>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8">

    <link href="http://cdn.bootcss.com/bootstrap/3.1.1/css/bootstrap.min.css" rel="stylesheet">
    <link href="http://cdn.bootcss.com/jstree/3.0.0/package/dist/themes/default/style.css" rel="stylesheet">
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

    <script type="text/javascript" src="http://code.jquery.com/jquery-2.0.3.min.js"></script>
    <script type="text/javascript" src="http://cdn.bootcss.com/bootstrap/3.1.1/js/bootstrap.min.js"></script>
    <script type="text/javascript" src="/static/js/webmaster.js"></script>
    <script src="http://cdn.bootcss.com/jstree/3.0.0/jstree.min.js"></script>
    <script src="/static/js/jstree.data.json"></script>
  </body>
</html>
