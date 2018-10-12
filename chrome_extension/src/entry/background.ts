import { createGraph } from '../client';

const graphGenerator = (directed: boolean) =>
    (info: chrome.contextMenus.OnClickData, tab: chrome.tabs.Tab) => {
        // object.selectedText does not include new line characters
        chrome.tabs.executeScript(tab.id, {
            code: "window.getSelection().toString()"
        }, createGraph(directed));
    };

const rootId = chrome.contextMenus.create({
    title: 'Online Judge Helper',
    contexts: ['selection']
});

chrome.contextMenus.create({
    title: 'Undirected Graph',
    contexts: ['selection'],
    parentId: rootId,
    onclick: graphGenerator(false)
});

chrome.contextMenus.create({
    title: 'Directed Graph',
    contexts: ['selection'],
    parentId: rootId,
    onclick: graphGenerator(true)
});
