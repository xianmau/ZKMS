<!DOCTYPE html>

<html>
  <head>
  	<title>Zookeeper Management System</title>
  	<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
    
  	<link href="http://cdn.bootcss.com/bootstrap/3.1.1/css/bootstrap.min.css" rel="stylesheet">

  	<style type="text/css">
  	.clear{clear: both;}
  	#leftside{position: relative; background: #FCFCFC;min-height: 500px;float: left;}
  	#rightside{position: relative; background: #333333;min-height: 500px;float: left;}
  	#logo{position: absolute;left:20px;top: 20px;width:80px;height: 80px;background: #EEE;border-radius:40px;}

	.nav{position: absolute;left:0px;top: 150px;width: 280px;padding: 0;margin: 0;}
	.nav li{width: 100%;border-bottom:1px solid #EEE;line-height: 80px;padding-left: 20px;font-size: 18px;font-weight: bold;color: #CCC;}
	.nav li:hover{border-bottom:1px solid #EEE;font-size: 18px;font-weight: bold;color: #FFF;background: #blue;}

  	</style>
  	</head>
  <body>

  	<div id="leftside">
  		<div id="logo"></div>

		<ul class="nav">
			<li>Dashboard</li>
			<li>Reports</li>
			<li>Zones</li>
			<li>Accounts</li>
			<li>Settings</li>
		</ul>

  	</div>

  	<div id="rightside">

  	</div>

  	<div class="clear"></div>
    <script type="text/javascript" src="http://code.jquery.com/jquery-2.0.3.min.js"></script>
    <script type="text/javascript" src="http://cdn.bootcss.com/bootstrap/3.1.1/js/bootstrap.min.js"></script>
    <script type="text/javascript">
    $(function(){
    	var wh = $(window).height()
    	var ww = $(window).width()
    	$('#leftside').height(wh)
    	$('#leftside').width(280)
    	$('#rightside').height(wh)
    	$('#rightside').width(ww - 280)
    });
    </script>
	</body>
</html>
