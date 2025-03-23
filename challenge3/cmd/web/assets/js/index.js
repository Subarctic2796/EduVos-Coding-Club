'use strict';

/**
 * @template T
 * @param {T|null|undefined} t
 * @param {string} msg
 * @returns {T}
 * @throws
*/
function removeNullUndefined(t, msg) {
  if (!t) throw new Error(msg);
  return t;
}

/** @type Map<number, string> */
const TODOS = new Map();

function getTodos() {
  const ls = document.querySelectorAll("div#items > div");
  if (ls.length === 0) {
    console.log(ls.length);
  } else {
    for (let i = 0; i < ls.length; i++) {
      let t = ls[i];
      let a = t.children.item(0).textContent;
      console.log(i, a, t.id);
    }
  }
}
