// import { augment } from '../client';

// const site = 'poj';
// const id = location.search.match(/id=(.*?)(?:&|$)/m)[1];

// const parser = () => {
//     const timeLimit = parseInt(document.querySelector('.plm').textContent.match(/(\d+)MS/)[1], 10) * 1e6;
//     const sio = Array.from(document.querySelectorAll(".sio"));
//     return {
//         site,
//         id,
//         restriction: {
//             timeLimit
//         },
//         testCases: [
//             {
//                 title: 'Sample',
//                 input: sio[0].textContent,
//                 output: sio[1].textContent
//             }
//         ]
//     };
// };

// const submit = (body: string) => {
//     console.log(body);
// }

// augment(document, site, id, parser, submit);
