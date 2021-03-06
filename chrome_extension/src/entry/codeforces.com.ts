import { augment } from '../client';
import { query } from "../codeforces";

const site = 'codeforces';
const id = location.pathname.split('/').slice(3).join('_');

const { timeLimit, testCases } = query(document);


augment(document, {
    site,
    id,
    restriction: {
        timeLimit,
    },
    testCases,
});
