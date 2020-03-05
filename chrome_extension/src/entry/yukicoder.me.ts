import { augment, inject } from '../client';

const site = 'yukicoder';
const id = location.pathname.split('/').pop();
const timeLimit = parseFloat(document.getElementById('content').textContent.match(/実行時間制限.*?(\d\.\d+)秒/)[1]) * 1e9;

const testCases = Array.from(document.querySelectorAll(".sample")).reduce((a, el) => {
  const title = el.querySelector('h5').textContent;
  const input = el.querySelector('pre:first-of-type').textContent;
  const output = el.querySelector('pre:last-of-type').textContent;
  a.push({
    title,
    input,
    output,
  });
  return a;
}, []);

augment(document, {
  site,
  id,
  restriction: {
    timeLimit,
  },
  testCases,
});


inject((val) => {
  document.dispatchEvent(new CustomEvent('Ace', {
      detail: val
  }));
}, site, id);

// inject a script to use Ace
// Ace ignores direct assignment to textarea
// and Content Script cannot access ace
const s = document.createElement('script');
s.src = chrome.runtime.getURL('ace.js');
s.addEventListener('load', function() {
    this.remove();
})
document.documentElement.appendChild(s);
