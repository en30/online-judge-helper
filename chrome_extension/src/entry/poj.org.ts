import { augment } from '../client';

const site = 'poj';
const id = location.search.match(/id=(.*?)(?:&|$)/m)[1];
const timeLimit = parseInt(document.querySelector('.plm').textContent.match(/(\d+)MS/)[1], 10) * 1e6;

const sio = Array.from(document.querySelectorAll(".sio"));
const testCases = [
    {
        title: 'Sample',
        input: sio[0].textContent,
        output: sio[1].textContent
    }
]

augment(document, {
    site,
    id,
    restriction: {
        timeLimit,
    },
    testCases,
});
