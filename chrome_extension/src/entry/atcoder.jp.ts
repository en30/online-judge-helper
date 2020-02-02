import { augment, inject } from '../client';

const site = 'atcoder';
const id = location.pathname.split('/').pop();

const timeLimit = parseInt(document.querySelector('#task-statement').previousElementSibling.textContent.match(/(\d+)\s*sec/)[1], 10) * 1e9;

// ABC 009~?
const normalizeNewline = (str: string) => str.trimRight() + "\n";

const testCases = Array.from(document.querySelectorAll("h3")).reduce((a, el) => {
    if (!el.textContent.match(/入力例|Sample Input/)) return a;
    const title = el.childNodes[0].textContent;
    const input = el.parentElement.querySelector('pre:last-of-type');
    let e = input.closest('.part').nextElementSibling;
    while (!e.classList.contains('part')) e = e.nextElementSibling;
    const output = e.querySelector('pre[id^=for_copy],pre[id^=pre-sample],pre:last-of-type');
    a.push({
        title,
        input: input.textContent,
        output: normalizeNewline(output.textContent),
    });
    return a;
}, []);

augment(document, {
    site,
    id,
    restriction: {
        timeLimit,
    },
    testCases,
});

inject((val) => {
    document.dispatchEvent(new CustomEvent('CodeMirror', {
        detail: val
    }));
}, site, id);

// inject a script to use CodeMirror
// CodeMirror ignores direct assignment to textarea
// and Content Script cannot access CodeMirror
const s = document.createElement('script');
s.src = chrome.runtime.getURL('code_mirror.js');
s.addEventListener('load', function() {
    this.remove();
})
document.documentElement.appendChild(s);
