package server

var introHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>OFC Uptime Monitor</title>
    <link rel="shortcut icon" href="data:image/x-icon;base64,AAABAAEAEBAAAAEAIABoBAAAFgAAACgAAAAQAAAAIAAAAAEAIAAAAAAAAAQAABILAAASCwAAAAAAAAAAAAAAAAAAAAAAACIiIgAiIiIHIiIiSSIiIqgiIiLkIiIi+yIiIvsiIiLkIiIiqCIiIkgiIiIHIiIiAAAAAAAAAAAAIiIiACIiIgAiIiIYIiIikiIiIu8iIiL/IiIi/yIiIv8iIiL/IiIi/yIiIv8iIiLuIiIikiIiIhgiIiIAISEhACIiIgAiIiIYIiIisSIiIv8gICD/ICAg/yIiIv8iIiL/IiIi/yIiIv8iIiL/IiIi/yIiIv8iIiKxIiIiGCIiIgAiIiIGIiIikiIiIv8kJCT/SEhI/0tLS/8jIyP/IiIi/yIiIv8iIiL/IiIi/yIiIv8iIiL/IiIi/yIiIpIiIiIGIiIiSSIiIu0fHx//Y2Nj/+jo6P+oqKj/JCQk/x8fH/8fHx//ICAg/yIiIv8fHx//ICAg/yIiIv8iIiLtIiIiSSIiIqciIiL/Hx8f/5WVlf/l5eX/RERE/0RERP+NjY3/kZGR/0tLS/8oKCj/e3t7/1xcXP8hISH/IiIi/yIiIqciIiLjIiIi/x4eHv+Xl5f/2tra/0lJSf/Ozs7/zc3N/8fHx//Y2Nj/UFBQ/9bW1v+YmJj/Hx8f/yIiIv8iIiLjIiIi+yAgIP8/Pz//zc3N/8TExP9hYWH/8PDw/2JiYv9TU1P/8PDw/3R0dP/V1dX/mJiY/yAgIP8iIiL/IiIi+yIiIvsgICD/SkpK/9/f3/+8vLz/QkJC/9PT0//Jycn/wsLC/9ra2v+NjY3/8fHx/9vb2/9hYWH/ICAg/yIiIvsiIiLjIiIi/yAgIP+ampr/3Nzc/zIyMv9ISEj/lJSU/5iYmP9PT0//U1NT/+Xl5f/U1NT/WVlZ/yAgIP8iIiLjIiIipyIiIv8fHx//l5eX/+Dg4P85OTn/Hx8f/x8fH/8fHx//ICAg/yIiIv+Li4v/3t7e/4SEhP8hISH/IiIipyIiIkkiIiLtHx8f/3R0dP/19fX/m5ub/yUlJf8iIiL/IiIi/yIiIv8iIiL/Jycn/z8/P/81NTX/IiIi7SIiIkkiIiIGIiIikiEhIf8rKyv/aGho/2dnZ/8kJCT/IiIi/yIiIv8iIiL/IiIi/yIiIv8gICD/ISEh/yIiIpIiIiIGIiIiACIiIhgiIiKxISEh/x8fH/8fHx//IiIi/yIiIv8iIiL/IiIi/yIiIv8iIiL/IiIi/yIiIrEiIiIYIiIiACIiIgAiIiIAIiIiGCIiIpIiIiLvIiIi/yIiIv8iIiL/IiIi/yIiIv8iIiL/IiIi7yIiIpIiIiIYIiIiACIiIgAAAAAAAAAAACIiIgAiIiIHIiIiSSIiIqgiIiLkIiIi/CIiIvwiIiLkIiIiqCIiIkgiIiIHIiIiAAAAAAAAAAAA4AcAAMADAACAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAIABAADAAwAA4AcAAA==" />
</head>
<body>
`

var extraHTML = `
<style>
body {
    font-family: sans-serif;
}

ul {
	list-style-type: none;
    padding: 0;
}

li {
	margin-bottom: 10px;
    display: flex;
    justify-content: space-between;
	align-items: center;
	padding: 0 5px;
    background: lightgray;
}

li.online {
    background: palegreen;
}

li.down {
    background: orangered;
}

a.button {
    display: inline-block;
    border: 2px solid black;
    padding: 3px 10px;
    text-decoration: none;
    color: black;
    background: white;
    margin-left: 5px;
    text-transform: uppercase;
}

li span a:hover {
    color: white;
    background: black;
}
</style>
</body>
</html>
`

func htmlWrap(html string) string {
	return introHTML + html + extraHTML
}
