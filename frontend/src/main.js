import 'core-js/stable';
const runtime = require('@wailsapp/runtime');

// Main entry point
function start() {
	let inners = document.querySelectorAll("div.cut-corner-inner")
	inners.forEach(e => {
		let p = e.parentElement;
		e.style.height = (p.clientHeight-20)+"px"
		e.style.width = (p.clientWidth-20)+"px"
		// console.log((p.clientWidth-e.clientWidth)/2)
		e.style.marginTop = ((p.clientHeight-e.clientHeight)/2)+"px"
		e.style.marginLeft = (((p.clientWidth-e.clientWidth)/2)-1)+"px"
	})

	let headers = document.querySelectorAll("div.cut-corner-inner > h1")
	headers.forEach(e => {
		let p = e.parentNode;
		p.style.borderColor = "black";
		e.style.top = ((p.clientHeight/2)-(e.clientHeight/2))+"px"
		e.style.left = ((p.clientWidth/2)-(e.clientWidth/2))+"px"
	})

	let images = document.querySelectorAll("div.cut-corner-inner > img")
	images.forEach(e => {
		let p = e.parentNode;
		if (e.clientHeight > e.clientWidth) {
			e.style.height = p.clientHeight
		} else if (e.clientHeight < e.clientWidth) {
			e.style.width = p.clientWidth
		} else {
			e.style.height = p.clientHeight
			e.style.width = p.clientWidth
		}
	})
};

// We provide our entrypoint as a callback for runtime.Init
runtime.Init(start);