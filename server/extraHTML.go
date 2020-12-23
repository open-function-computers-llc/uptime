package server

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
`
