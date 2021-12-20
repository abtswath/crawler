package javascript

/**
From https://github.com/Qianlitp/crawlergo/blob/master/pkg/js/javascript.go
*/
const (
	InjectionScript = `
(function() {
	Object.defineProperty(navigator, 'webdriver', {
		get: () => false
	});
	Object.defineProperty(navigator, 'plugins', {
		get: () => [1, 2, 3, 4, 5]
	});
	window.chrome = {
		runtime: {}
	};
	const originalQuery = window.navigator.permissions.query;
	window.navigator.permissions.query = (parameters) => (
    	parameters.name === 'notifications' ?
			Promise.resolve({ state: Notification.permission }) :
			originalQuery(parameters)
	);
	Object.defineProperty(navigator, 'platform', {
		get: () => 'win32'
	});
	Object.defineProperty(navigator, 'language', {
		get: () => 'zh-CN'
	});
	Object.defineProperty(navigator, 'languages', {
		get: () => ["zh-CN", "zh"]
	});
	window.history.pushState = (a, b, c) => {
		window.collectURL(c);
	}
	window.history.replaceState = (a, b, c) => {
		window.collectURL(c);
	}
	Object.defineProperty(window.history, 'pushState', {'writable': false, 'configurable': false});
	Object.defineProperty(window.history, 'replaceState', {'writable': false, 'configurable': false});
	window.addEventListener('hashchange', () => {
		window.collectURL(document.location.href);
	});
	window.open = (url) => {
		window.collectURL(url);
	}
	Object.defineProperty(window, 'open', {'writable': false, 'configurable': false});

	window.close = () => {
		console.log('Trying to close page...')
	}
	Object.defineProperty(window, 'close', {'writable': false, 'configurable': false});

	HTMLFormElement.prototype.reset = () => {console.log('Trying to reset form...')};
	Object.defineProperty(HTMLFormElement.prototype, 'reset',{'writable': false, 'configurable': false});

	const oldEventHandler = Element.prototype.addEventListener;
	Element.prototype.addEventListener = function(eventName, eventFunc, useCapture) {
		let events = [eventName];
		if (this.hasAttribute('auto-trigger-events')) {
			events = events.concat(this.getAttribute('auto-trigger-events').split(','));
		}
		this.setAttribute('auto-trigger-events', events.join(','));
		oldEventHandler.apply(this, arguments);
	}
	
	const dom0ListenerHook = (target, eventName) => {
		let events = [eventName];
		if (target.hasAttribute('auto-trigger-events')) {
			events = events.concat(target.getAttribute('auto-trigger-events').split(','));
		}
		target.setAttribute('auto-trigger-events', events.join(','));
	}
	
	Object.defineProperties(HTMLElement.prototype, {
		onclick: {set: function (newValue) {onclick = newValue;dom0ListenerHook(this, 'click');}},
		onchange: {set: function (newValue) {onchange = newValue;dom0ListenerHook(this, 'change');}},
		onblur: {set: function (newValue) {onblur = newValue;dom0ListenerHook(this, 'blur');}},
		ondblclick: {set: function (newValue) {ondblclick = newValue;dom0ListenerHook(this, 'dbclick');}},
		onfocus: {set: function (newValue) {onfocus = newValue;dom0ListenerHook(this, 'focus');}},
		onkeydown: {set: function (newValue) {onkeydown = newValue;dom0ListenerHook(this, 'keydown');}},
		onkeypress: {set: function (newValue) {onkeypress = newValue;dom0ListenerHook(this, 'keypress');}},
		onkeyup: {set: function (newValue) {onkeyup = newValue;dom0ListenerHook(this, 'keyup');}},
		onload: {set: function (newValue) {onload = newValue;dom0ListenerHook(this, 'load');}},
		onmousedown: {set: function (newValue) {onmousedown = newValue;dom0ListenerHook(this, 'mousedown');}},
		onmousemove: {set: function (newValue) {onmousemove = newValue;dom0ListenerHook(this, 'mousemove');}},
		onmouseout: {set: function (newValue) {onmouseout = newValue;dom0ListenerHook(this, 'mouseout');}},
		onmouseover: {set: function (newValue) {onmouseover = newValue;dom0ListenerHook(this, 'mouseover');}},
		onmouseup: {set: function (newValue) {onmouseup = newValue;dom0ListenerHook(this, 'mouseup');}},
		onreset: {set: function (newValue) {onreset = newValue;dom0ListenerHook(this, 'reset');}},
		onresize: {set: function (newValue) {onresize = newValue;dom0ListenerHook(this, 'resize');}},
		onselect: {set: function (newValue) {onselect = newValue;dom0ListenerHook(this, 'select');}},
		onsubmit: {set: function (newValue) {onsubmit = newValue;dom0ListenerHook(this, 'submit');}},
		onunload: {set: function (newValue) {onunload = newValue;dom0ListenerHook(this, 'unload');}},
		onabort: {set: function (newValue) {onabort = newValue;dom0ListenerHook(this, 'abort');}},
		onerror: {set: function (newValue) {onerror = newValue;dom0ListenerHook(this, 'error');}},
	});
})();
`

	AfterDOMLoadedScript = `
(function () {
	const sleep = delay => {
		return new Promise(resolve => {
			setTimeout(resolve, delay)
		});
	}
	function getNodeUniqueID (n) {
		const tagName = n.tagName.toLowerCase();
		const classes = n.className.trim();
		let classList = '';
		if (classes !== '') {
			classList = '.' + classes.split(' ').join('.')
		}
		let id = '';
		if (n.id !== '') {
			id = '#' + n.id;
		}
		return tagName + classList + id;
	}
	(async function triggerInlineEvents() {
		const events = ['abort', 'blur', 'change', 'click', 'dblclick', 'error', 'focus', 'keydown', 'keypress', 'keyup', 'load', 'mousedown', 'mousemove', 'mouseout', 'mouseover', 'mouseup', 'reset', 'resize', 'select', 'submit', 'unload'];
		for (const evt of events) {
			let nodes = document.querySelectorAll('[on' + evt + ']');
			if (nodes.length > 100) {
				nodes = nodes.slice(0, 100);
			}

			for (const node of nodes) {
				await sleep(100);
				const event = new CustomEvent(evt);
				evt.initCustomEvent(event, false, true, null);
				try {
					node.dispatchEvent(event);
					console.log('Node [' + getNodeUniqueID(node) + '] inline event [' + evt + '] has been triggered...');
				} catch {
				}
			}
		}
	})();
	(async function triggerDomEvents() {
		let nodes = document.querySelectorAll('[auto-trigger-events]');
		if (nodes.length > 200) {
			nodes = nodes.slice(0, 200);
		}
		for (let node of nodes) {
			await sleep(100);
			let events = new Set(node.getAttribute('auto-trigger-events').split(','));
			for (let eventName of events) {
				if ((node.className && node.className.includes("close")) || (node.id && node.id.includes("close"))) {
					continue;
				}
				const evt = new CustomEvent(eventName);
				try {
					node.dispatchEvent(evt);
					console.log('Node [' + getNodeUniqueID(node) + '] event [' + eventName + '] has been triggered...');
				} catch (e) {
				}
			}
		}
	})();
	(async function triggerTagAEvent() {
		const nodes = document.querySelectorAll('[href]');
		for (let node of nodes) {
			const href = node.getAttribute("href");
			if (href.toLocaleLowerCase().startsWith("javascript:")) {
				await sleep(100);
				try {
					eval(href.substring(11));
					console.log('Tag a [' + getNodeUniqueID(node) + '] has been clicked...');
				} catch {
				}
			}
		}
	})();
})();
`
)
