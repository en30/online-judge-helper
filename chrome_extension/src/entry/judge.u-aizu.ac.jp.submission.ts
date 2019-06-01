import { inject } from '../client';

const site = 'aoj';
const id = location.hash.replace(/^#?submit\//, '').replace('/', '_');
inject((v) => {
    const el = document.getElementById('submit_source') as HTMLInputElement | null;
    if (el) el.value = v;
}, site, id);
