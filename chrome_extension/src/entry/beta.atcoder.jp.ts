// import { augment } from '../client';

// const site = 'atcoder';
// const id = location.pathname.split('/').pop();

// const parser = () => {
//     const timeLimit = parseInt(document.querySelector('#task-statement').previousElementSibling.textContent.match(/(\d+)\s*sec/)[1], 10) * 1e9;

//     const testCases = Array.from(document.querySelectorAll("h3")).reduce((a, el) => {
//         if (!el.textContent.match(/入力例/)) return a;
//         const title = el.childNodes[0].textContent;
//         const input = el.parentElement.querySelector('pre:last-of-type');
//         const output = input.closest('.part').nextElementSibling.querySelector('pre:last-of-type');
//         a.push({
//             title,
//             input: input.textContent,
//             output: output.textContent,
//         });
//         return a;
//     }, []);
//     return {
//         site,
//         id,
//         restriction: {
//             timeLimit,
//         },
//         testCases,
//     };
// };

// const submit = (body: string) => {
//     console.log(body);
// }

// augment(document, site, id, parser, submit);
