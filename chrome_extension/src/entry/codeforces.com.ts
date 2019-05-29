// import { augment } from '../client';

// const site = 'codeforces';
// const id = location.pathname.split('/').slice(3).join('_');

// const parser = () => {
//     const timeLimit = parseInt(document.querySelector('.time-limit').textContent.match(/(\d+)\s*seconds/)[1], 10) * 1e9;
//     const testCases = Array.from(document.querySelectorAll('.sample-test > .input')).reduce((a, el, i) => {

//         const input = el.querySelector('pre');
//         const output = el.nextElementSibling.querySelector('pre');
//         a.push({
//             title: `Example ${i}`,
//             input: input.innerText,
//             output: input.innerText,
//         })
//         return a;
//     }, []);

//     return {
//         site: 'codeforces',
//         id,
//         restriction: {
//             timeLimit,
//         },
//         testCases,
//     };
// }

// const submit = (body: string) => {
//     console.log(body);
// }

// augment(document, site, id, parser, submit);
