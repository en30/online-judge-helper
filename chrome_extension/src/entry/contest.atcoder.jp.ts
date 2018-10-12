import { augment } from '../client';

augment(document, () => {
    const timeLimit = parseInt(document.querySelector('#task-statement').previousElementSibling.textContent.match(/(\d+)sec/)[1], 10) * 1e9;
    const id = location.pathname.split('/').pop();

    const testCases = Array.from(document.querySelectorAll("h3")).reduce((a, el) => {
        if (!el.textContent.match(/入力例/)) return a;
        const title = el.textContent;
        const input = el.parentElement.querySelector('pre:last-of-type');
        const output = input.closest('.part').nextElementSibling.querySelector('pre:last-of-type');
        a.push({
            title,
            input: input.textContent,
            output: output.textContent,
        });
        return a;
    }, []);
    return {
        site: 'atcoder',
        id,
        restriction: {
            timeLimit,
        },
        testCases,
    };
});
