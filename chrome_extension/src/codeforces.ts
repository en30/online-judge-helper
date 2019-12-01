export const query = (document: HTMLDocument) => {
    const timeLimit = parseInt(document.querySelector('.time-limit').textContent.match(/(\d+)\s*seconds?/)[1], 10) * 1e9;
    const testCases = Array.from(document.querySelectorAll('.sample-tests .input')).reduce((a, el, i) => {

        const input = el.querySelector('pre');
        const output = el.nextElementSibling.querySelector('pre');
        a.push({
            title: `Example ${i}`,
            input: input.textContent.replace(/\n?$/, "\n"),
            output: output.textContent.replace(/\n?$/, "\n"),
        })
        return a;
    }, []);

    return { timeLimit, testCases };
}
