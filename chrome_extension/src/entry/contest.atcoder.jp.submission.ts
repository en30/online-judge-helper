import { inject } from '../client';

const site = 'atcoder';
const contest = location.hostname.split('.')[0];
const task = document.querySelector('#submit-task-selector option:checked').textContent.split('-')[0].trim().toLowerCase();
const id = `${contest}_${task}`;
inject((v) => {
    const el = document.querySelector('textarea[name="source_code"]') as HTMLTextAreaElement | null;
    if (el) el.value = v;
}, site, id);
