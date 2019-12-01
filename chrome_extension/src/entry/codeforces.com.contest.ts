import { augment } from '../client';
import { query } from "../codeforces";

const site = 'codeforces';
const match = location.pathname.match(/\/contest\/(.*?)\/problem\/(.+)/)
const id = `${match[1]}_${match[2]}`;

const { timeLimit, testCases } = query(document);

augment(document, {
    site,
    id,
    restriction: {
        timeLimit,
    },
    testCases,
});
