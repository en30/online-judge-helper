document.addEventListener('CodeMirror', (e: CustomEvent) => {
    const editor = (<any>document.querySelector('.CodeMirror')).CodeMirror;
    editor.getDoc().setValue(e.detail);
});
