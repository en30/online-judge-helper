

interface Window {
  ace: any,
}

document.addEventListener('Ace', (e: CustomEvent) => {
  window.ace.edit(document.querySelector(".ace_editor")).setValue(e.detail);
});
