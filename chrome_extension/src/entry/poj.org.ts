import { augment } from '../client';

augment(document, () => {
    const id = location.search.match(/id=(.*?)(?:&|$)/m)[1];
    const timeLimit = parseInt(document.querySelector('.plm').textContent.match(/(\d+)MS/)[1], 10) * 1e6;
    const sio = Array.from(document.querySelectorAll(".sio"));
    return {
        site: 'poj',
        id,
        restriction: {
            timeLimit
        },
        testCases: [
            {
                title: 'Sample',
                input: sio[0].textContent,
                output: sio[1].textContent
            }
        ]
    };
})
