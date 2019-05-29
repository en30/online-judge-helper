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

const getSubmission = async (site: string, id: string): Promise<string> => {
    const res = await client.get(`/submission?id=${id}&site=${site}`);
    return res.data;
}

export const inject = async (input: HTMLInputElement, site: string, id: string) => {
    const source = await getSubmission(site, id);
    input.value = source;
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
