import axios from 'axios';

interface TestCase {
    title: string,
    input: string,
    output: string,
}

interface Restriction {
    timeLimit: number
}

interface Problem {
    site: string,
    id: string,
    testCases: Array<TestCase>,
    restriction: Restriction
}

const client = axios.create({
    baseURL: 'http://localhost:4567'
});

const appendButton = (document: Document) => {
    const e = document.createElement('button');
    e.textContent = 'solve';
    e.style.cssText = `
        cursor: pointer;
        position: fixed;
        left: 0;
        bottom: 0;
        padding: 1em 2em;
        color: white;
        background-color: rgba(0, 0, 255, .4);
        z-index: 2147483647;
    `;
    document.body.appendChild(e);
    return e
};

const solve = async (data: Problem) => {
    try {
        if (process.env.NODE_ENV !== 'production') console.log(data);
        const res = await client.post('/problem', data);
        if (res.status == 0) {
            alert('It seems that background server is not working');
        }
    } catch (err) {
        alert(err);
    }
};

type Parser = () => Problem;
type Submitter = (submission: string) => void;
export const augment = (document: Document, problem: Problem) => {
    appendButton(document).addEventListener('click', () => {
        solve(problem);
    });
}

export const inject = (callback: (value: string) => void, site: string, id: string) => {
    const socket = new WebSocket(`ws://localhost:4567/submission?id=${id}&site=${site}`);
    socket.addEventListener("message", function(event) {
        callback(event.data);
    });
    socket.addEventListener("close", function(event) {
        setTimeout(() => {
            inject(callback, site, id);
        }, 1000)
    });
    socket.addEventListener("error", function(event) {
        console.error(event);
        socket.close();
    });
}

export const createGraph = (directed: boolean) => async function(selection: Array<string>) {
    try {
        const res = await client.post(
            '/graph',
            { directed, adjacentList: selection[0] },
            { responseType: 'blob' });
        window.open(URL.createObjectURL(res.data));
    } catch (err) {
        alert(err);
    }
}
