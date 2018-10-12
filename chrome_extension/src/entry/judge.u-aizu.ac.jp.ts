import { augment } from '../client';

const next = (el: Element, selector?: string): Element | null => {
    let e = el.nextElementSibling;
    if (!selector) return e;
    while (e !== null && !e.matches(selector)) e = e.nextElementSibling;
    return e
}

augment(document, () => {
    const id = location.search.match(/id=(.*?)(?:&|$)/m)[1];
    const timeLimit = parseInt(document.querySelector('#pageinfo').textContent.match(/(\d+) sec/)[1], 10) * 1e9;
    const testCases = Array.from(document.querySelectorAll("h2")).reduce((a, el) => {
        if (!el.textContent.match(/^Sample Input/)) return a;
        const input = next(el, 'pre');
        const output = next(input, 'pre');
        a.push({
            title: el.textContent,
            input: input.textContent,
            output: output.textContent
        });
        return a;
    }, []);
    return {
        site: 'aoj',
        id,
        testCases,
        restriction: {
            timeLimit,
        },
    };
})
