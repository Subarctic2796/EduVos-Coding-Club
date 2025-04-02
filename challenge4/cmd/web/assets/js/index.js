document.documentElement.classList.toggle(
  "dark",
  localStorage.theme === "dark" ||
  (!("theme" in localStorage) && window.matchMedia("(prefers-color-scheme: dark)").matches)
);
document.documentElement.classList.toggle("dark");

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

function save() {
  /** @type HTMLTextAreaElement */
  const note = removeNullUndefined(document.querySelector("#note"), "no note");
  const key = `note-${name.value}`;
  localStorage.setItem(key, note.value);
  note.value = "";
}
