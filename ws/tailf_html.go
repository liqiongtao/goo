package ws

import "text/template"

var (
	tailfHtml = `
	<!DOCTYPE html>
	<html lang="en">
		<head>
			<meta charset="utf-8">
			<title></title>
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<style>
				body {background:#2C3E50; color:#2ECC71;}
				pre {width:100%; word-break:break-all; white-space:pre-wrap; word-wrap:break-word;}
			</style>
		</head>
		<body>
			<pre id="message"></pre>
			<script type="text/javascript">
				(function() {
					var msg = document.getElementById('message');
					var ws = new WebSocket(window.location.href.replace('http', 'ws'));
					ws.onopen = function(evt) {
						console.log('Connection opened');
					}
					ws.onclose = function(evt) {
						console.log('Connection closed');
					}
					ws.onmessage = function(evt) {
						try {
							msg.innerText += decodeURIComponent(evt.data);
						} catch(e) {
						}
					}
				})();
			</script>
		</body>
	</html>
`
	tailfTmpl = template.Must(template.New("").Parse(tailfHtml))
)
