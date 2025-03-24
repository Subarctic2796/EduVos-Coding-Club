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

function saveLocalTodos() {
  for (const k of TODOS.keys()) {
    localStorage.setItem(`todos-${k}`, TODOS.get(k));
  }
}

/** @param {string} k */
function removeTodo(k) {
  TODOS.delete(k);
  localStorage.removeItem(`todos-${k}`);
}

function getLocalTodos() {
  const itemsdiv = removeNullUndefined(document.querySelector("div#items"), "something went wrong");
  for (let i = 0; i < localStorage.length; i++) {
    let key = localStorage.key(i);
    let id = key.substring(6);
    let value = localStorage.getItem(key);
    let el = new HTMLDivElement();
    el.id = id;
    el.className = "my-0.5 flex flex-row gap-1";
    el.innerHTML = `<p class="w-8/10 text-lg text-center bg-slate-900 border-2 border-slate-800 rounded-lg">${value}</p>
		<button
			class="w-2/10 bg-red-500 hover:bg-red-600 border-2 border-red-600 rounded-lg"
			hx-delete="/front/delete/${id}"
			hx-on::before-request="removeTodo(this.parentElement.id)"
		>Done</button>`;
    itemsdiv.append(el);
  }
}

function getLastTodo() {
  const last = removeNullUndefined(document.querySelector("div#items > div:last-of-type"), "couldn't get last");
  const id = last.id;
  const msg = removeNullUndefined(last.firstChild.textContent, "something very bad has happened");
  TODOS.set(id, msg);
}
