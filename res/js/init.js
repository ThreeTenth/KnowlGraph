
const i18n = Vue.observable({
	...defaultLang,
})

let i18ns = new Map();
i18ns.set(defaultLang.language, defaultLang)