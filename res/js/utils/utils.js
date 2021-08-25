// utils.js

function canUseWenAuthn() {
  if (window.PublicKeyCredential === undefined ||
    typeof window.PublicKeyCredential !== "function") {
    return false
  }

  return true
}

var lookup = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/'

function uint8ToBase64(uint8) {
  var i
  var extraBytes = uint8.length % 3 // if we have 1 byte left, pad 2 bytes
  var output = ''
  var temp, length

  function encode(num) {
    return lookup.charAt(num)
  }

  function tripletToBase64(num) {
    return encode(num >> 18 & 0x3F) + encode(num >> 12 & 0x3F) + encode(num >> 6 & 0x3F) + encode(num & 0x3F)
  }

  // go through the array every three bytes, we'll deal with trailing stuff later
  for (i = 0, length = uint8.length - extraBytes; i < length; i += 3) {
    temp = (uint8[i] << 16) + (uint8[i + 1] << 8) + (uint8[i + 2])
    output += tripletToBase64(temp)
  }

  // pad the end with zeros, but make sure to not forget the extra bytes
  switch (extraBytes) {
    case 1:
      temp = uint8[uint8.length - 1]
      output += encode(temp >> 2)
      output += encode((temp << 4) & 0x3F)
      output += '=='
      break
    case 2:
      temp = (uint8[uint8.length - 2] << 8) + (uint8[uint8.length - 1])
      output += encode(temp >> 10)
      output += encode((temp >> 4) & 0x3F)
      output += encode((temp << 2) & 0x3F)
      output += '='
      break
    default:
      break
  }

  return output
}

// Encode an ArrayBuffer into a base64 string.
function bufferEncode(value) {
  return uint8ToBase64(value)
    .replace(/\+/g, "-")
    .replace(/\//g, "_")
    .replace(/=/g, "");
}

function calcElementLength(element) {
  var text = element.innerText
  if ('' === text || '\n' === text) {
    return text.length
  }
  return text.length + 1
}

function isChildAt(parent, child) {
  if (!child) {
    return false
  }

  if (child == parent) {
    return true
  }

  if (child == document.body) {
    return false
  }

  return isChildAt(parent, child.parentElement)
}

function seo(title, gist) {
  let vt = title
  if (!vt) {
    vt = gist
    if (60 < gist.length) {
      let idx = vt.search(/[.?!。？！]/g)
      if (1 < idx) {
        vt = gist.substr(0, idx)
      }
    }
  }

  return vt, gist, encodeURLTitle(vt)
}

function encodeURLTitle(title) {
  return title.trim()
    .replace(/[;,/?:@&=+$_.!~*'()# \n]+/g, '-')
    .replace(/-$/g, '')
    .toLowerCase()
}

// by @Stian Grytøyr's https://github.com/stiang/remove-markdown
function removeMarkdown(md, options) {
  options = options || {};
  options.listUnicodeChar = options.hasOwnProperty('listUnicodeChar') ? options.listUnicodeChar : false;
  options.stripListLeaders = options.hasOwnProperty('stripListLeaders') ? options.stripListLeaders : true;
  options.gfm = options.hasOwnProperty('gfm') ? options.gfm : true;
  options.useImgAltText = options.hasOwnProperty('useImgAltText') ? options.useImgAltText : true;

  var output = md || '';

  // Remove horizontal rules (stripListHeaders conflict with this rule, which is why it has been moved to the top)
  output = output.replace(/^(-\s*?|\*\s*?|_\s*?){3,}\s*$/gm, '');

  try {
    if (options.stripListLeaders) {
      if (options.listUnicodeChar)
        output = output.replace(/^([\s\t]*)([\*\-\+]|\d+\.)\s+/gm, options.listUnicodeChar + ' $1');
      else
        output = output.replace(/^([\s\t]*)([\*\-\+]|\d+\.)\s+/gm, '$1');
    }
    if (options.gfm) {
      output = output
        // Header
        .replace(/\n={2,}/g, '\n')
        // Fenced codeblocks
        .replace(/~{3}.*\n/g, '')
        // Strikethrough
        .replace(/~~/g, '')
        // Fenced codeblocks
        .replace(/`{3}.*\n/g, '');
    }
    output = output
      // Remove HTML tags
      .replace(/<[^>]*>/g, '')
      // Remove setext-style headers
      .replace(/^[=\-]{2,}\s*$/g, '')
      // Remove footnotes?
      .replace(/\[\^.+?\](\: .*?$)?/g, '')
      .replace(/\s{0,2}\[.*?\]: .*?$/g, '')
      // Remove images
      .replace(/\!\[(.*?)\][\[\(].*?[\]\)]/g, options.useImgAltText ? '$1' : '')
      // Remove inline links
      .replace(/\[(.*?)\][\[\(].*?[\]\)]/g, '$1')
      // Remove blockquotes
      .replace(/^\s{0,3}>\s?/g, '')
      // Remove reference-style links?
      .replace(/^\s{1,2}\[(.*?)\]: (\S+)( ".*?")?\s*$/g, '')
      // Remove atx-style headers
      .replace(/^(\n)?\s{0,}#{1,6}\s+| {0,}(\n)?\s{0,}#{0,} {0,}(\n)?\s{0,}$/gm, '$1$2$3')
      // Remove emphasis (repeat the line to remove double emphasis)
      .replace(/([\*_]{1,3})(\S.*?\S{0,1})\1/g, '$2')
      .replace(/([\*_]{1,3})(\S.*?\S{0,1})\1/g, '$2')
      // Remove code blocks
      .replace(/(`{3,})(.*?)\1/gm, '$2')
      // Remove inline code
      .replace(/`(.+?)`/g, '$1')
      // Replace two or more newlines with exactly two? Not entirely sure this belongs here...
      .replace(/\n{2,}/g, '\n\n');
  } catch (e) {
    console.error(e);
    return md;
  }
  return output;
};

function getUserLang() {
  var language = Cookies.get("user-lang")
  if (language != undefined)
    return language

  if (navigator.languages != undefined)
    language = navigator.languages[0]
  else
    language = navigator.language || window.navigator.userLanguage

  var localLang = language.split("-")
  if (2 == localLang.length)
    language = localLang[0]

  return language
}

function setUserLang(lang) {
  Cookies.set("user-lang", lang)
}

function getMeta(metaName) {
  const metas = document.getElementsByTagName('meta')

  for (let i = 0; i < metas.length; i++) {
    if (metas[i].getAttribute('name') === metaName) {
      return metas[i].getAttribute('content')
    }
  }

  return ''
}

function getLink(linkRel) {
  const metas = document.getElementsByTagName('link')

  for (let i = 0; i < metas.length; i++) {
    if (metas[i].getAttribute('rel') === linkRel) {
      return metas[i].getAttribute('href')
    }
  }

  return ''
}

function setI18nStrings(lang, strings) {
  localStorage.setItem('strings-' + lang, JSON.stringify(strings));
}

function getI18nStrings(lang) {
  let strings = localStorage.getItem('strings-' + lang);
  return JSON.parse(strings)
}

function removeI18nStrings(lang) {
  localStorage.removeItem('strings-' + lang);
}