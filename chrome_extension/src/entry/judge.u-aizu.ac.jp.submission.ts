import { inject } from '../client';

const site = 'aoj';
const id = location.hash.replace(/^#?submit\//, '').replace('/', '_');
inject(document.getElementById('submit_source') as HTMLInputElement, site, id);
