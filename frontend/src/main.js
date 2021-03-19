import 'core-js/stable';
const runtime = require('@wailsapp/runtime');

// Main entry point
function start() {
	const app = document.getElementById("app")
	const loader = document.querySelector(".wails-reconnect-overlay")
	const loaderMessage = document.getElementById("wails-reconnect-overlay-message")

	let buttons = runtime.Store.New("Buttons")
	buttons.subscribe( (btns) => {
		renderButtons(btns)
		updateButtons()
	})

	runtime.Events.On("ready", ready => {
		if (!ready) {
			loaderMessage.innerHTML = "Waiting for server connection"
			app.innerHTML = ""
			loader.style.display = "block"
		} else {
			window.backend.Client.UpdateButtons()
			loader.style.display = "none"
		}
	})
};

function renderButtons(buttons) {
	app.innerHTML = ""
	buttons.forEach(button => {
		let typeRender

		if (button.type == "text")
			typeRender = `<h1>${button.name}</h1>`
		else if (button.type == "image")
			typeRender = `<img alt="${button.name}" src="${button.img}" />`
		else
			return

		app.innerHTML += `
			<div class="grid-item">
				<div class="cut-corner">
					<div class="cut-corner-inner" data-name="${button.name}" data-plugin="${button.plugin}" data-action="${button.action}">
						${typeRender}
					</div>
				</div>
			</div>
		`
	})
}

function updateButtons() {
	let inners = document.querySelectorAll("div.cut-corner-inner")
	inners.forEach(e => {
		let p = e.parentElement;
		e.style.height = (p.clientHeight-20)+"px"
		e.style.width = (p.clientWidth-20)+"px"
		e.style.marginTop = ((p.clientHeight-e.clientHeight)/2)+"px"
		e.style.marginLeft = (((p.clientWidth-e.clientWidth)/2)-1)+"px"

		e.addEventListener("click", (e) => {
			e.currentTarget.dataset.plugin
			window.backend.Client.PressButton(e.currentTarget.dataset.plugin, e.currentTarget.dataset.action)
		})
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
}

// We provide our entrypoint as a callback for runtime.Init
runtime.Init(start);